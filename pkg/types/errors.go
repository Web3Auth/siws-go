package types

import (
	"fmt"
)

type ExpiredMessage struct{ String string }
type InvalidMessage struct{ String string }
type InvalidSignature struct{ String string }

func (m *ExpiredMessage) Error() string {
	return "Expired Message"
}

func (m *InvalidMessage) Error() string {
	return "Invalid Message"
}

func (m *InvalidSignature) Error() string {
	return fmt.Sprintf("Invalid Signature: %s", m.String)
}
