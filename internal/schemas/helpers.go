package schemas

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// Timestamp format used in /web/check-updates.php
const timestampFormat = "2006-01-02 15:04:05"

// Timestamp is like time.Time, but provides json marshalling.
// This was introduced for the `ReleaseFiles` schema specifically.
type Timestamp struct {
	time.Time
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), "\"")

	isNull := value == "null"
	isEmpty := value == ""
	if isNull || isEmpty {
		t.Time = time.Time{}
		return nil
	}

	parsed, err := time.Parse(timestampFormat, value)
	if err != nil {
		// TODO: Maybe have a list of timestamp formats to try instead of just this one
		return err
	}
	t.Time = parsed
	return nil
}

func (t Timestamp) Value() (driver.Value, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *Timestamp) Scan(value any) error {
	if value == nil {
		t.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		t.Time = v
		return nil
	case string:
		parsed, err := time.Parse(timestampFormat, v)
		if err != nil {
			return err
		}
		t.Time = parsed
		return nil
	case []byte:
		parsed, err := time.Parse(timestampFormat, string(v))
		if err != nil {
			return err
		}
		t.Time = parsed
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into Timestamp", value)
	}
}

func ResolveSafeName(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "_"))
}
