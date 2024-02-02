package m

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Setting struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Value     SettingValue
}

type SettingValue struct {
	EnableAuthentication bool
	AuthenticationType   *string // nil | form | basic
}

func (j *SettingValue) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := SettingValue{}
	err := json.Unmarshal(bytes, &result)
	*j = SettingValue(result)
	return err

}
func (j SettingValue) Value() (driver.Value, error) {
	v, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(v).MarshalJSON()
}
