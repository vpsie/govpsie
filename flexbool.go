package govpsie

import (
	"bytes"
	"fmt"
)

// FlexBool is a flag the API may encode as a JSON boolean, as the numbers 0 / 1,
// or as those values in a string (load balancer settings come back as 0 / 1).
// Use it through a pointer: a missing field or a JSON null leaves the pointer nil.
type FlexBool bool

// UnmarshalJSON implements json.Unmarshaler.
func (b *FlexBool) UnmarshalJSON(data []byte) error {
	switch string(bytes.Trim(bytes.TrimSpace(data), `"`)) {
	case "true", "1":
		*b = true
	case "false", "0":
		*b = false
	default:
		return fmt.Errorf("govpsie: expected a boolean or 0 / 1, got %s", data)
	}
	return nil
}

// Bool returns the value as a *bool, nil when the field was absent or null.
func (b *FlexBool) Bool() *bool {
	if b == nil {
		return nil
	}
	v := bool(*b)
	return &v
}
