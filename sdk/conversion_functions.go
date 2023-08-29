package sdk

import (
	"context"
	"fmt"
	"time"
)

type ConversionFunctions interface {
	ToTimestampLTZ(ctx context.Context, t time.Time) (time.Time, error)
	ToTimestampNTZ(ctx context.Context, t time.Time) (time.Time, error)
}

type conversionFunctions struct {
	client *Client
}

func (v *conversionFunctions) ToTimestampLTZ(ctx context.Context, t time.Time) (time.Time, error) {
	s := &struct {
		ToTimestampLTZ time.Time `db:"TO_TIMESTAMP_LTZ"`
	}{}
	sql := fmt.Sprintf(`SELECT TO_TIMESTAMP_LTZ('%s') AS "TO_TIMESTAMP_LTZ"`, t.Format(time.RFC3339Nano))
	err := v.client.queryOne(ctx, s, sql)
	if err != nil {
		return time.Time{}, err
	}
	return s.ToTimestampLTZ, nil
}

func (v *conversionFunctions) ToTimestampNTZ(ctx context.Context, t time.Time) (time.Time, error) {
	s := &struct {
		ToTimestampNTZ time.Time `db:"TO_TIMESTAMP_NTZ"`
	}{}
	sql := fmt.Sprintf(`SELECT TO_TIMESTAMP_NTZ('%s') AS "TO_TIMESTAMP_NTZ"`, t.Format(time.RFC3339Nano))
	err := v.client.queryOne(ctx, s, sql)
	if err != nil {
		return time.Time{}, err
	}
	return s.ToTimestampNTZ, nil
}
