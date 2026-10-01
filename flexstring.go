package govpsie

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// FlexString is a text value the API may encode either as a JSON string or as
// a JSON number. Monetary amounts such as invoice totals and monthly charges
// arrive in both forms depending on the endpoint and environment; FlexString
// accepts either and keeps the value's textual form, so no precision is lost.
// A JSON null decodes to the empty string.
type FlexString string

// UnmarshalJSON implements json.Unmarshaler.
func (s *FlexString) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)

	switch {
	case bytes.Equal(data, []byte("null")):
		*s = ""
		return nil
	case len(data) > 0 && data[0] == '"':
		var v string
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		*s = FlexString(v)
		return nil
	}

	var n json.Number
	if err := json.Unmarshal(data, &n); err != nil {
		return fmt.Errorf("govpsie: expected a string or number, got %s", data)
	}
	*s = FlexString(n.String())

	return nil
}

// String returns the value as text.
func (s FlexString) String() string {
	return string(s)
}
