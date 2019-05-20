package instrumentedsql

import (
	"context"
	"database/sql/driver"
	"fmt"

	"github.com/kr/pretty"
	"github.com/pkg/errors"
)

type wrappedDriver struct {
	Logger
	Tracer
	parent driver.Driver
}

type wrappedConn struct {
	Logger
	Tracer
	parent driver.Conn
}

type wrappedTx struct {
	Logger
	Tracer
	ctx    context.Context
	parent driver.Tx
}

type wrappedStmt struct {
	Logger
	Tracer
	ctx    context.Context
	query  string
	parent driver.Stmt
}

type wrappedResult struct {
	Logger
	Tracer
	ctx    context.Context
	parent driver.Result
}

type wrappedRows struct {
	Logger
	Tracer
	ctx    context.Context
	parent driver.Rows
}

// WrapDriver will wrap the passed SQL driver and return a new sql driver that uses it and also logs and traces calls using the passed logger and tracer
// The returned driver will still have to be registered with the sql package before it can be used.
//
// Important note: Seeing as the context passed into the various instrumentation calls this package calls,
// Any call without a context passed will not be instrumented. Please be sure to use the ___Context() and BeginTx() function calls added in Go 1.8
// instead of the older calls which do not accept a context.
func WrapDriver(driver driver.Driver, opts ...Opt) driver.Driver {
	d := wrappedDriver{parent: driver}

	for _, opt := range opts {
		opt(&d)
	}

	if d.Logger == nil {
		d.Logger = nullLogger{}
	}
	if d.Tracer == nil {
		d.Tracer = nullTracer{}
	}

	return d
}

func (d wrappedDriver) Open(name string) (driver.Conn, error) {
	conn, err := d.parent.Open(name)
	if err != nil {
		return nil, err
	}

	return wrappedConn{Tracer: d.Tracer, Logger: d.Logger, parent: conn}, nil
}

func (c wrappedConn) Prepare(query string) (driver.Stmt, error) {
	parent, err := c.parent.Prepare(query)
	if err != nil {
		return nil, err
	}

	return wrappedStmt{Tracer: c.Tracer, Logger: c.Logger, query: query, parent: parent}, nil
}

func (c wrappedConn) Close() error {
	return c.parent.Close()
}

func (c wrappedConn) Begin() (driver.Tx, error) {
	tx, err := c.parent.Begin()
	if err != nil {
		return nil, err
	}

	return wrappedTx{Tracer: c.Tracer, Logger: c.Logger, parent: tx}, nil
}

func (c wrappedConn) BeginTx(ctx context.Context, opts driver.TxOptions) (tx driver.Tx, err error) {
	span := c.GetSpan(ctx).NewChild("sql-tx-begin")
	span.SetLabel("component", "database/sql")
	defer func() {
		span.SetLabel("err", fmt.Sprint(err))
		span.Finish()
		c.Log(ctx, "sql-tx-begin", "err", err)
	}()

	if connBeginTx, ok := c.parent.(driver.ConnBeginTx); ok {
		tx, err = connBeginTx.BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}

		return wrappedTx{Tracer: c.Tracer, Logger: c.Logger, ctx: ctx, parent: tx}, nil
	}

	tx, err = c.parent.Begin()
	if err != nil {
		return nil, err
	}

	return wrappedTx{Tracer: c.Tracer, Logger: c.Logger, ctx: ctx, parent: tx}, nil
}

func (c wrappedConn) PrepareContext(ctx context.Context, query string) (stmt driver.Stmt, err error) {
	span := c.GetSpan(ctx).NewChild("sql-prepare")
	span.SetLabel("component", "database/sql")
	defer func() {
		span.SetLabel("err", fmt.Sprint(err))
		span.Finish()
		c.Log(ctx, "sql-prepare", "err", err)
	}()

	if connPrepareCtx, ok := c.parent.(driver.ConnPrepareContext); ok {
		stmt, err := connPrepareCtx.PrepareContext(ctx, query)
		if err != nil {
			return nil, err
		}

		return wrappedStmt{Tracer: c.Tracer, Logger: c.Logger, ctx: ctx, parent: stmt}, nil
	}

	return c.Prepare(query)
}

func (c wrappedConn) Exec(query string, args []driver.Value) (driver.Result, error) {
	if execer, ok := c.parent.(driver.Execer); ok {
		res, err := execer.Exec(query, args)
		if err != nil {
			return nil, err
		}

		return wrappedResult{Tracer: c.Tracer, Logger: c.Logger, parent: res}, nil
	}

	return nil, driver.ErrSkip
}

func (c wrappedConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (r driver.Result, err error) {
	span := c.GetSpan(ctx).NewChild("sql-conn-exec")
	span.SetLabel("component", "database/sql")
	span.SetLabel("query", query)
	span.SetLabel("args", pretty.Sprint(args))
	defer func() {
		span.SetLabel("err", fmt.Sprint(err))
		span.Finish()
		c.Log(ctx, "sql-conn-exec", "query", query, "args", pretty.Sprint(args), "err", err)
	}()

	if execContext, ok := c.parent.(driver.ExecerContext); ok {
		res, err := execContext.ExecContext(ctx, query, args)
		if err != nil {
			return nil, err
		}

		return wrappedResult{Tracer: c.Tracer, Logger: c.Logger, ctx: ctx, parent: res}, nil
	}

	// Fallback implementation
	dargs, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}

	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return c.Exec(query, dargs)
}

func (c wrappedConn) Ping(ctx context.Context) (err error) {
	if pinger, ok := c.parent.(driver.Pinger); ok {
		span := c.GetSpan(ctx).NewChild("sql-ping")
		span.SetLabel("component", "database/sql")
		defer func() {
			span.SetLabel("err", fmt.Sprint(err))
			span.Finish()
			c.Log(ctx, "sql-ping", "err", err)
		}()

		return pinger.Ping(ctx)
	}

	c.Log(ctx, "sql-dummy-ping")

	return nil
}

func (c wrappedConn) Query(query string, args []driver.Value) (driver.Rows, error) {
	if queryer, ok := c.parent.(driver.Queryer); ok {
		rows, err := queryer.Query(query, args)
		if err != nil {
			return nil, err
		}

		return wrappedRows{Tracer: c.Tracer, Logger: c.Logger, parent: rows}, nil
	}

	return nil, driver.ErrSkip
}

func (c wrappedConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (rows driver.Rows, err error) {
	span := c.GetSpan(ctx).NewChild("sql-conn-query")
	span.SetLabel("component", "database/sql")
	span.SetLabel("query", query)
	span.SetLabel("args", pretty.Sprint(args))
	defer func() {
		span.SetLabel("err", fmt.Sprint(err))
		span.Finish()
		c.Log(ctx, "sql-conn-query", "query", query, "args", pretty.Sprint(args), "err", err)
	}()

	if queryerContext, ok := c.parent.(driver.QueryerContext); ok {
		rows, err := queryerContext.QueryContext(ctx, query, args)
		if err != nil {
			return nil, err
		}

		return wrappedRows{Tracer: c.Tracer, Logger: c.Logger, ctx: ctx, parent: rows}, nil
	}

	dargs, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}

	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return c.Query(query, dargs)
}

func (t wrappedTx) Commit() (err error) {
	span := t.GetSpan(t.ctx).NewChild("sql-tx-commit")
	span.SetLabel("component", "database/sql")
	defer func() {
		span.SetLabel("err", fmt.Sprint(err))
		span.Finish()
		t.Log(t.ctx, "sql-tx-commit", "err", err)
	}()

	return t.parent.Commit()
}

func (t wrappedTx) Rollback() (err error) {
	span := t.GetSpan(t.ctx).NewChild("sql-tx-rollback")
	span.SetLabel("component", "database/sql")
	defer func() {
		span.SetLabel("err", fmt.Sprint(err))
		span.Finish()
		t.Log(t.ctx, "sql-tx-rollback", "err", err)
	}()

	return t.parent.Rollback()
}

func (s wrappedStmt) Close() (err error) {
	span := s.GetSpan(s.ctx).NewChild("sql-stmt-close")
	span.SetLabel("component", "database/sql")
	defer func() {
		span.SetLabel("err", fmt.Sprint(err))
		span.Finish()
		s.Log(s.ctx, "sql-stmt-close", "err", err)
	}()

	return s.parent.Close()
}

func (s wrappedStmt) NumInput() int {
	return s.parent.NumInput()
}

func (s wrappedStmt) Exec(args []driver.Value) (res driver.Result, err error) {
	span := s.GetSpan(s.ctx).NewChild("sql-stmt-exec")
	span.SetLabel("component", "database/sql")
	span.SetLabel("query", s.query)
	span.SetLabel("args", pretty.Sprint(args))
	defer func() {
		span.SetLabel("err", fmt.Sprint(err))
		span.Finish()
		s.Log(s.ctx, "sql-stmt-exec", "query", s.query, "args", pretty.Sprint(args), "err", err)
	}()

	res, err = s.parent.Exec(args)
	if err != nil {
		return nil, err
	}

	return wrappedResult{Tracer: s.Tracer, Logger: s.Logger, ctx: s.ctx, parent: res}, nil
}

func (s wrappedStmt) Query(args []driver.Value) (rows driver.Rows, err error) {
	span := s.GetSpan(s.ctx).NewChild("sql-stmt-query")
	span.SetLabel("component", "database/sql")
	span.SetLabel("query", s.query)
	span.SetLabel("args", pretty.Sprint(args))
	defer func() {
		span.SetLabel("err", fmt.Sprint(err))
		span.Finish()
		s.Log(s.ctx, "sql-stmt-query", "query", s.query, "args", pretty.Sprint(args), "err", err)
	}()

	rows, err = s.parent.Query(args)
	if err != nil {
		return nil, err
	}

	return wrappedRows{Tracer: s.Tracer, Logger: s.Logger, ctx: s.ctx, parent: rows}, nil
}

func (s wrappedStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (res driver.Result, err error) {
	span := s.GetSpan(ctx).NewChild("sql-stmt-exec")
	span.SetLabel("component", "database/sql")
	span.SetLabel("query", s.query)
	span.SetLabel("args", pretty.Sprint(args))
	defer func() {
		span.SetLabel("err", fmt.Sprint(err))
		span.Finish()
		s.Log(ctx, "sql-stmt-exec", "query", s.query, "args", pretty.Sprint(args), "err", err)
	}()

	if stmtExecContext, ok := s.parent.(driver.StmtExecContext); ok {
		res, err := stmtExecContext.ExecContext(ctx, args)
		if err != nil {
			return nil, err
		}

		return wrappedResult{Tracer: s.Tracer, Logger: s.Logger, ctx: ctx, parent: res}, nil
	}

	// Fallback implementation
	dargs, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}

	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return s.Exec(dargs)
}

func (s wrappedStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (rows driver.Rows, err error) {
	span := s.GetSpan(ctx).NewChild("sql-stmt-query")
	span.SetLabel("component", "database/sql")
	span.SetLabel("query", s.query)
	span.SetLabel("args", pretty.Sprint(args))
	defer func() {
		span.SetLabel("err", fmt.Sprint(err))
		span.Finish()
		s.Log(ctx, "sql-stmt-query", "query", s.query, "args", pretty.Sprint(args), "err", err)
	}()

	if stmtQueryContext, ok := s.parent.(driver.StmtQueryContext); ok {
		rows, err := stmtQueryContext.QueryContext(ctx, args)
		if err != nil {
			return nil, err
		}

		return wrappedRows{Tracer: s.Tracer, Logger: s.Logger, ctx: ctx, parent: rows}, nil
	}

	dargs, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}

	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return s.Query(dargs)
}

func (r wrappedResult) LastInsertId() (id int64, err error) {
	span := r.GetSpan(r.ctx).NewChild("sql-res-lastInsertId")
	span.SetLabel("component", "database/sql")
	defer func() {
		span.SetLabel("err", fmt.Sprint(err))
		span.Finish()
		r.Log(r.ctx, "sql-res-lastInsertId", "err", err)
	}()

	return r.parent.LastInsertId()
}

func (r wrappedResult) RowsAffected() (num int64, err error) {
	span := r.GetSpan(r.ctx).NewChild("sql-res-rowsAffected")
	span.SetLabel("component", "database/sql")
	defer func() {
		span.SetLabel("err", fmt.Sprint(err))
		span.Finish()
		r.Log(r.ctx, "sql-res-rowsAffected", "err", err)
	}()

	return r.parent.RowsAffected()
}

func (r wrappedRows) Columns() []string {
	return r.parent.Columns()
}

func (r wrappedRows) Close() error {
	return r.parent.Close()
}

func (r wrappedRows) Next(dest []driver.Value) (err error) {
	span := r.GetSpan(r.ctx).NewChild("sql-rows-next")
	span.SetLabel("component", "database/sql")
	defer func() {
		span.SetLabel("err", fmt.Sprint(err))
		span.Finish()
		r.Log(r.ctx, "sql-rows-next", "err", err)
	}()

	return r.parent.Next(dest)
}

// namedValueToValue is a helper function copied from the database/sql package
func namedValueToValue(named []driver.NamedValue) ([]driver.Value, error) {
	dargs := make([]driver.Value, len(named))
	for n, param := range named {
		if len(param.Name) > 0 {
			return nil, errors.New("sql: driver does not support the use of Named Parameters")
		}
		dargs[n] = param.Value
	}
	return dargs, nil
}
