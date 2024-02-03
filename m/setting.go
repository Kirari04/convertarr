package m

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Setting struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Value     SettingValue
}

type SettingValue struct {
	HasBeenSetup bool

	// auth
	EnableAuthentication bool
	AuthenticationType   *string // nil | form | basic
	Username             string
	PwdHash              string

	// scanning
	EnableAutomaticScanns    bool
	AutomaticScannsInterval  time.Duration
	AutomaticScannsAtStartup bool
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
