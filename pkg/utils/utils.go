package siwe

import (
	"fmt"
	"github.com/Web3Auth/siws-go/pkg/types"
	"strings"
	"time"

	"github.com/dchest/uniuri"
	"github.com/relvacode/iso8601"
)

func ParseTimestamp(fields map[string]interface{}, key string) (*string, error) {
	var value string

	if val, ok := fields[key]; ok {
		switch val.(type) {
		case time.Time:
			value = val.(time.Time).UTC().Format(time.RFC3339)
		case string:
			_, err := iso8601.ParseString(val.(string))
			if err != nil {
				return nil, &types.InvalidMessage{fmt.Sprintf("Invalid format for field `%s`", key)}
			}
			value = val.(string)
		default:
			return nil, &types.InvalidMessage{fmt.Sprintf("`%s` must be either an ISO8601 formatted string or time.Time", key)}
		}
	}

	if value == "" {
		return nil, nil
	}

	return &value, nil
}

func GenerateNonce() string {
	return uniuri.NewLen(16)
}

func IsEmpty(str *string) bool {
	return str == nil || len(strings.TrimSpace(*str)) == 0
}

func IsStringAndNotEmpty(m map[string]interface{}, k string) (*string, bool) {
	if v, ok := m[k]; ok {
		switch v.(type) {
		case string:
			s := v.(string)
			if s != "" {
				return &s, true
			}
		}
	}
	return nil, false
}
