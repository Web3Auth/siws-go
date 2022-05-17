package message

import (
	"net/url"
)

type Payload struct {
	Domain  string `json:"domain"`
	Address string `json:"address"`
	Uri     url.URL `json:"uri"`
	Version string `json:"version"`

	Statement *string `json:"statement"`
	Nonce     string `json:"nonce"`
	ChainID   int `json:"chain_id"`

	IssuedAt       string `json:"issued_at"`
	ExpirationTime *string `json:"expiration_time"`
	NotBefore      *string `json:"not_before"`

	RequestID *string `json:"request_id"`
	Resources []string `json:"resources"`
}

type Header struct {
	T string `json:"t"`
}

type Signature struct {
	T string `json:"t"`

	M SignatureMeta `json:"m"`

	S string `json:"s"`
}

type SignatureMeta struct {

}

type Message struct {
	Payload Payload `json:"payload,omitempty"`
	Header Header `json:"header,omitempty"`
	Signature Signature `json:"signature,omitempty"`
}
