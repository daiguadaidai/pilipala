package types

import (
	"database/sql"
	"encoding/json"
)

type NullInt64 struct {
	sql.NullInt64
}

func GetNullInt64(_data int64) NullInt64 {
	return NullInt64{sql.NullInt64{_data, true}}
}

func (v *NullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	} else {
		return json.Marshal(nil)
	}
}

func (v *NullInt64) UnmarshalJSON(data []byte) error {
	var s *int64
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.Int64 = *s
	} else {
		v.Valid = false
	}
	return nil
}
