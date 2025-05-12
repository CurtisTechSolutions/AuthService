package db

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Birthday is a custom type for handling birthday dates in JSON.
// It uses the time.Time type internally and provides a custom UnmarshalJSON method
// to parse the date in the format "YYYY-MM-DD".
type Birthday struct {
	time.Time
}

func (d *Birthday) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}

	// take care of null..
	if len(b) == 0 || string(b) == "null" {
		d.Time = time.Time{}
		return
	}

	d.Time, err = time.Parse("2006-01-02", string(b))
	return
}

func (d *Birthday) Scan(b interface{}) (err error) {
	switch x := b.(type) {
	case time.Time:
		d.Time = x
	default:
		err = fmt.Errorf("unsupported scan type %T", b)
	}
	return
}

func (d Birthday) Value() (driver.Value, error) {
	// check if the date was not set..
	if d.Time.IsZero() {
		return nil, nil
	}
	return d.Time.Format("2006-01-02"), nil
}
