package message

import (
	"crypto/ed25519"
	"fmt"
	"github.com/Web3Auth/siws-go/pkg/types"
	"net/url"
	"strings"
	"time"

	"github.com/relvacode/iso8601"

	utils "github.com/Web3Auth/siws-go/pkg/utils"
)

func (m *Message) GetDomain() string {
	return m.Payload.Domain
}

func (m *Message) GetAddress() string {
	return m.Payload.Address
}

func (m *Message) GetURI() url.URL {
	return m.Payload.Uri
}

func (m *Message) GetVersion() string {
	return m.Payload.Version
}

func (m *Message) GetStatement() *string {
	if m.Payload.Statement != nil {
		ret := *m.Payload.Statement
		return &ret
	}
	return nil
}

func (m *Message) GetNonce() string {
	return m.Payload.Nonce
}

func (m *Message) GetChainID() int {
	return m.Payload.ChainID
}

func (m *Message) GetIssuedAt() string {
	return m.Payload.IssuedAt
}

func (m *Message) getExpirationTime() *time.Time {
	if !utils.IsEmpty(m.Payload.ExpirationTime) {
		ret, _ := iso8601.ParseString(*m.Payload.ExpirationTime)
		return &ret
	}
	return nil
}

func (m *Message) GetExpirationTime() *string {
	if m.Payload.ExpirationTime != nil {
		ret := *m.Payload.ExpirationTime
		return &ret
	}
	return nil
}

func (m *Message) getNotBefore() *time.Time {
	if !utils.IsEmpty(m.Payload.NotBefore) {
		ret, _ := iso8601.ParseString(*m.Payload.NotBefore)
		return &ret
	}
	return nil
}

func (m *Message) GetNotBefore() *string {
	if m.Payload.NotBefore != nil {
		ret := *m.Payload.NotBefore
		return &ret
	}
	return nil
}

func (m *Message) GetRequestID() *string {
	if m.Payload.RequestID != nil {
		ret := *m.Payload.RequestID
		return &ret
	}
	return nil
}

func (m *Message) GetResources() []string {
	return m.Payload.Resources
}

func (m *Message) PrepareMessage() string {
	greeting := fmt.Sprintf("%s wants you to sign in with your Solana account:", m.Payload.Domain)
	headerArr := []string{greeting, m.Payload.Address}
	if utils.IsEmpty(m.Payload.Statement) {
		headerArr = append(headerArr, "\n")
	} else {
		headerArr = append(headerArr, fmt.Sprintf("\n%s\n", *m.Payload.Statement))
	}

	header := strings.Join(headerArr, "\n")

	uri := fmt.Sprintf("URI: %s", m.Payload.Uri.String())
	version := fmt.Sprintf("Version: %s", m.Payload.Version)
	chainId := fmt.Sprintf("Chain ID: %d", m.Payload.ChainID)
	nonce := fmt.Sprintf("Nonce: %s", m.Payload.Nonce)
	issuedAt := fmt.Sprintf("Issued At: %s", m.Payload.IssuedAt)

	bodyArr := []string{uri, version, chainId, nonce, issuedAt}

	if !utils.IsEmpty(m.Payload.ExpirationTime) {
		value := fmt.Sprintf("Expiration Time: %s", *m.Payload.ExpirationTime)
		bodyArr = append(bodyArr, value)
	}

	if !utils.IsEmpty(m.Payload.NotBefore) {
		value := fmt.Sprintf("Not Before: %s", *m.Payload.NotBefore)
		bodyArr = append(bodyArr, value)
	}

	if !utils.IsEmpty(m.Payload.RequestID) {
		value := fmt.Sprintf("Request ID: %s", *m.Payload.RequestID)
		bodyArr = append(bodyArr, value)
	}

	if len(m.Payload.Resources) > 0 {
		resourcesArr := make([]string, len(m.Payload.Resources))
		for i, v := range m.Payload.Resources {
			resourcesArr[i] = fmt.Sprintf("- %s", v)
		}

		resources := strings.Join(resourcesArr, "\n")
		value := fmt.Sprintf("Resources:\n%s", resources)

		bodyArr = append(bodyArr, value)
	}

	body := strings.Join(bodyArr, "\n")

	return strings.Join([]string{header, body}, "\n")
}

func (m *Message) ValidNow() (bool, error) {
	return m.ValidAt(time.Now().UTC())
}

func (m *Message) ValidAt(when time.Time) (bool, error) {
	if m.Payload.ExpirationTime != nil {
		if time.Now().UTC().After(*m.getExpirationTime()) {
			return false, &types.ExpiredMessage{"Message expired"}
		}
	}

	if m.Payload.NotBefore != nil {
		if time.Now().UTC().Before(*m.getNotBefore()) {
			return false, &types.InvalidMessage{"Message not yet valid"}
		}
	}

	return true, nil
}

func (m *Message) Verify(signature string, nonce *string, timestamp *time.Time) (bool, error) {
	var err error

	if timestamp != nil {
		_, err = m.ValidAt(*timestamp)
	} else {
		_, err = m.ValidNow()
	}

	if err != nil {
		return false, err
	}

	if nonce != nil {
		if m.GetNonce() != *nonce {
			return false, &types.InvalidSignature{"Message nonce doesn't match"}
		}
	}
	resp := m.VerifySIP99(signature)
	if resp  {
		return true,nil
	}
	return false, &types.InvalidSignature{"Message signature invalid"}
}

func (m *Message) VerifySIP99(signature string) (bool) {
	if utils.IsEmpty(&signature) {
		return false
	}

	return ed25519.Verify(ed25519.PublicKey(m.Payload.Address),[]byte(m.String()), []byte(signature))
}

func (m *Message) String() string {
	return m.PrepareMessage()
}
