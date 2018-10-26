package sql_type_util

import (
	"database/sql"
	"encoding/json"
)

type NullInt64 struct {
	sql.NullInt64
}

func(this *NullInt64) MarshalJSON() ([]byte, error) {
	if this.Valid {
		return json.Marshal(this.Int64)
	} else {
		return json.Marshal(nil)
	}
}

func (v NullInt64) UnmarshalJSON(data []byte) error {
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

