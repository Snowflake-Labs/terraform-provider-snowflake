package sdk

import (
	"context"
	"database/sql"
	"regexp"
	"strconv"
	"strings"
)

func (v *Share) ExternalID() ExternalObjectIdentifier {
	return NewExternalObjectIdentifier(NewAccountIdentifierFromFullyQualifiedName(v.OwnerAccount), NewAccountObjectIdentifier(v.Name))
}

// shareRecipientsTruncationRegexp matches the marker Snowflake appends to the "to" column of SHOW SHARES once a share
// has more than three recipients, e.g. "ORG.A,ORG.B,ORG.C,... (4 TOTAL)".
var shareRecipientsTruncationRegexp = regexp.MustCompile(`^\.\.\. \((\d+) TOTAL\)$`)

// additionalConvert maps the "to" column of SHOW SHARES manually, because the column is not a plain list of
// "<organization_name>.<account_name>" values: for more than three recipients Snowflake lists only the first three
// and appends a "... (N TOTAL)" marker. Parsing the whole column strictly used to fail on that marker, which left
// Share.To empty and made snowflake_share.accounts show up as a diff on every plan.
func (r shareRow) additionalConvert(s *Share) error {
	s.To, s.ToTotal = parseShareRecipients(r.To)
	return nil
}

// parseShareRecipients parses the "to" column of SHOW SHARES entry by entry, so one entry that is not in the
// "<organization_name>.<account_name>" form does not discard the others. The returned total is the number of
// recipients Snowflake reports when the column is truncated, and 0 when it lists all of them.
func parseShareRecipients(value string) ([]AccountIdentifier, int) {
	items := ParseCommaSeparatedStringArray(value, false)
	ids := make([]AccountIdentifier, 0, len(items))
	total := 0
	for _, item := range items {
		if item == "" {
			continue
		}
		if match := shareRecipientsTruncationRegexp.FindStringSubmatch(item); match != nil {
			total, _ = strconv.Atoi(match[1])
			continue
		}
		if id, err := ParseAccountIdentifier(item); err == nil {
			ids = append(ids, id)
			continue
		}
		parts := strings.SplitN(item, ".", 2)
		if len(parts) == 1 {
			ids = append(ids, NewAccountIdentifierFromAccountLocator(strings.Trim(parts[0], `"`)))
			continue
		}
		ids = append(ids, NewAccountIdentifier(parts[0], parts[1]))
	}
	return ids, total
}

func (r *shareDetailsRow) additionalConvert(s *ShareInfo) error {
	s.Kind = ObjectType(r.Kind)
	trimmedS := strings.Trim(r.Name, "\"")
	// TODO(SNOW-1229218): Use a common mapper to get object id.
	s.Name = s.Kind.GetObjectIdentifier(trimmedS)
	return nil
}

// showShareTargetAccountsOptions reads the full recipient list of an outbound share. SHOW SHARES and
// INFORMATION_SCHEMA.SHARES list at most three recipients; SNOWFLAKE.ACCOUNT_USAGE.SHARES lists all of them
// (with a latency of up to 3 hours), in the "<organization_name>.<account_name>" form.
type showShareTargetAccountsOptions struct {
	selectFrom bool   `ddl:"static" sql:"SELECT TARGET_ACCOUNTS FROM SNOWFLAKE.ACCOUNT_USAGE.SHARES WHERE DELETED_ON IS NULL AND NAME ="`
	name       string `ddl:"keyword,single_quotes"`
	orderBy    bool   `ddl:"static" sql:"ORDER BY MODIFIED_ON DESC LIMIT 1"`
}

type shareTargetAccountsRow struct {
	TargetAccounts sql.NullString `db:"TARGET_ACCOUNTS"`
}

func (s *shares) ShowTargetAccounts(ctx context.Context, id AccountObjectIdentifier) ([]AccountIdentifier, error) {
	query, err := structToSQL(&showShareTargetAccountsOptions{name: id.Name()})
	if err != nil {
		return nil, err
	}
	var row shareTargetAccountsRow
	if err := s.client.queryOne(ctx, &row, query); err != nil {
		return nil, err
	}
	if !row.TargetAccounts.Valid {
		return []AccountIdentifier{}, nil
	}
	ids, _ := parseShareRecipients(row.TargetAccounts.String)
	return ids, nil
}

func (s *shares) DescribeProvider(ctx context.Context, id AccountObjectIdentifier) ([]ShareInfo, error) {
	return s.Describe(ctx, id)
}

func (s *shares) DescribeConsumer(ctx context.Context, id ExternalObjectIdentifier) ([]ShareInfo, error) {
	opts := &DescribeShareOptions{name: id}
	dbRows, err := validateAndQuery[shareDetailsRow](s.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[shareDetailsRow, ShareInfo](dbRows)
}
