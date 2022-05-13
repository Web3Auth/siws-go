package message

import (
	"net/url"
)

type Payload struct {
	Domain  string
	Address string
	Uri     url.URL
	Version string

	Statement *string
	Nonce     string
	ChainID   int

	IssuedAt       string
	ExpirationTime *string
	NotBefore      *string

	RequestID *string
	Resources []string
}

type Header struct {
	T string
}

type Signature struct {
	T string // signature scheme

	M SignatureMeta // signature related metadata (optional)

	S string // signature
}

type SignatureMeta struct {

}

type Message struct {
	Payload Payload
	Header Header
	Signature Signature
}
