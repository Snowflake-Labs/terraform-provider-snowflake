package sdk

import (
	"context"
	"database/sql"
)

func (c *Client) ExecUnsafe(ctx context.Context, sql string) (sql.Result, error) {
	return c.exec(ctx, sql)
}

// QueryUnsafe for now only supports single query. For more queries we will have to adjust the behaviour. From the gosnowflake driver docs:
//
//	 (...) while using the multi-statement feature, pass a Context that specifies the number of statements in the string.
//		When multiple queries are executed by a single call to QueryContext(), multiple result sets are returned. After you process the first result set, get the next result set (for the next SQL statement) by calling NextResultSet().
//
// Therefore, only single resultSet is processed.
func (c *Client) QueryUnsafe(ctx context.Context, sql string) ([]map[string]*any, error) {
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

func unsafeExecuteProcessRows(rows *sql.Rows) ([]map[string]*any, error) {
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	processedRows := make([]map[string]*any, 0)
	for rows.Next() {
		row, err := unsafeExecuteProcessRow(rows, columnNames)
		if err != nil {
			return nil, err
		}
		processedRows = append(processedRows, row)
	}

	return processedRows, nil
}

func unsafeExecuteProcessRow(rows *sql.Rows, columnNames []string) (map[string]*any, error) {
	values := make([]any, len(columnNames))
	for i, _ := range values {
		values[i] = new(any)
	}

	err := rows.Scan(values...)
	if err != nil {
		return nil, err
	}

	row := make(map[string]*any)
	for i, col := range columnNames {
		row[col] = values[i].(*any)
	}
	return row, nil
}
