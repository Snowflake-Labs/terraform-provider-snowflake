package sdk

import (
	"context"
	"database/sql"
	"fmt"
)

func (c *Client) ExecUnsafe(ctx context.Context, sql string) (sql.Result, error) {
	return c.exec(ctx, sql)
}

func (c *Client) QueryUnsafe(ctx context.Context, sql string) ([]map[string]string, error) {
	rows, err := c.db.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	allRows, err := unsafeExecuteProcessRows(rows)
	if err != nil {
		return nil, err
	}
	return allRows, nil
}

func unsafeExecuteProcessRows(rows *sql.Rows) ([]map[string]string, error) {
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	allRows := make([]map[string]string, 0)

	unsafeExecuteProcessResultSet := func(rows *sql.Rows, columnNames []string) error {
		for rows.Next() {
			row, err := unsafeExecuteProcessRow(rows, columnNames)
			if err != nil {
				return err
			}
			allRows = append(allRows, row)
		}
		return nil
	}

	err = unsafeExecuteProcessResultSet(rows, columnNames)
	if err != nil {
		return nil, err
	}
	for rows.NextResultSet() {
		err := unsafeExecuteProcessResultSet(rows, columnNames)
		if err != nil {
			return nil, err
		}
	}

	return allRows, nil
}

func unsafeExecuteProcessRow(rows *sql.Rows, columnNames []string) (map[string]string, error) {
	values := make([]any, len(columnNames))
	for i, _ := range values {
		values[i] = new(any)
	}

	err := rows.Scan(values...)
	if err != nil {
		return nil, err
	}

	row := make(map[string]string)
	for i, col := range columnNames {
		row[col] = fmt.Sprintf("%v", *values[i].(*interface{}))
	}
	return row, nil
}
