package date

import (
	"encoding/json"
	"strings"
	"time"
)

// Date custom date type to help in marshaling and unmarshaling the
// intended date layout. The layout used is `02/01/2006`, that is
// `dd/mm/yyyy`.
//
// This struct implements the [json.Marshaler] and the [json.Unmarshaler]
// interfaces.
type Date struct {
	time.Time
}

const dateLayout = "02/01/2006" // dd/mm/yyyy

// Today returns the [time.Now] date parsed using the `02/01/2006` date layout.
// If parsing returns an error, then the time is just the default [time.Now]
func Today() Date {
	now, err := time.Parse(dateLayout, time.Now().Format(dateLayout))
	if err != nil {
		return Date{Time: time.Now()}
	}
	return Date{Time: now}
}

// UnmarshalJSON implements the [json.Unmarshaler] interface.
func (d *Date) UnmarshalJSON(b []byte) error {
	if len(strings.Trim(string(b), `"`)) == 0 {
		return nil
	}
	if string(b) == "null" {
		return nil
	}

	// Remove quotes
	s := strings.Trim(string(b), `"`)

	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return err
	}

	d.Time = t
	return nil
}

// MarshalJSON implements the [json.Marshaler] interface.
func (d *Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Format(dateLayout))
}
