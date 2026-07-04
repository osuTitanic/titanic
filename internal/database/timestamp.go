package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// TODO: We eventually want to migrate to timezone-aware timestamps for all schemas
// 	     Unfortunately, I made the stupid decision to use non-timezoned timestamps so far
// 	     So here's the current workaround: a custom codec that forces all timestamps to be
//       encoded as UTC

// The schema uses "timestamp without time zone" columns, which the original
// python services populated with naive UTC datetimes.
// pgx already scans these as UTC, but on encode it discards the timezone
// and stores the local wall clock.
// This codec pins encodes to UTC, so callers can pass any time.Time
// without worrying about its timezone.

// registerUTCTimestampCodec installs the UTC timestamp codecs
func registerUTCTimestampCodec(ctx context.Context, conn *pgx.Conn) error {
	timestampType := &pgtype.Type{
		Name:  "timestamp",
		OID:   pgtype.TimestampOID,
		Codec: &utcTimestampCodec{pgtype.TimestampCodec{ScanLocation: time.UTC}},
	}
	conn.TypeMap().RegisterType(timestampType)
	conn.TypeMap().RegisterType(&pgtype.Type{
		Name:  "_timestamp",
		OID:   pgtype.TimestampArrayOID,
		Codec: &pgtype.ArrayCodec{ElementType: timestampType},
	})
	return nil
}

type utcTimestampCodec struct {
	pgtype.TimestampCodec
}

func (c *utcTimestampCodec) PlanEncode(m *pgtype.Map, oid uint32, format int16, value any) pgtype.EncodePlan {
	plan := c.TimestampCodec.PlanEncode(m, oid, format, value)
	if plan == nil {
		return nil
	}
	// On encode, convert any time.Time to UTC before passing to the next plan
	return utcEncodePlan{plan}
}

type utcEncodePlan struct {
	next pgtype.EncodePlan
}

func (p utcEncodePlan) Encode(value any, buf []byte) ([]byte, error) {
	valuer, ok := value.(pgtype.TimestampValuer)
	if !ok {
		return p.next.Encode(value, buf)
	}

	ts, err := valuer.TimestampValue()
	if err != nil {
		return nil, err
	}

	// We don't want to change the value of infinite or invalid timestamps
	if ts.Valid && ts.InfinityModifier == pgtype.Finite {
		ts.Time = ts.Time.UTC()
	}
	return p.next.Encode(ts, buf)
}
