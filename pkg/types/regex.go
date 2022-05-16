package types

import (
	"fmt"
	"regexp"
)

const SIWS_DOMAIN = "^(?P<domain>([^?#]*)) wants you to sign in with your Solana account:\\n"
const SIWS_ADDRESS = "(?P<address>[a-zA-Z0-9]{32,44})\\n\\n"
const SIWS_STATEMENT = "((?P<statement>[^\\n]+)\\n)?\\n"
const SIWS_URI = "(([^:?#]+):)?(([^?#]*))?([^?#]*)(\\?([^#]*))?(#(.*))"

var SIWS_URI_LINE = fmt.Sprintf("URI: (?P<uri>%s?)\\n", SIWS_URI)

const SIWS_VERSION = "Version: (?P<version>1)\\n"
const SIWS_CHAIN_ID = "Chain ID: (?P<chainId>[0-9]+)\\n"
const SIWS_NONCE = "Nonce: (?P<nonce>[a-zA-Z0-9]{8,})\\n"
const SIWS_DATETIME = "([0-9]+)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60)(\\.[0-9]+)?(([Zz])|([\\+|\\-]([01][0-9]|2[0-3]):[0-5][0-9]))"

var SIWS_ISSUED_AT = fmt.Sprintf("Issued At: (?P<issuedAt>%s)", SIWS_DATETIME)
var SIWS_EXPIRATION_TIME = fmt.Sprintf("(\\nExpiration Time: (?P<expirationTime>%s))?", SIWS_DATETIME)
var SIWS_NOT_BEFORE = fmt.Sprintf("(\\nNot Before: (?P<notBefore>%s))?", SIWS_DATETIME)

const SIWS_REQUEST_ID = "(\\nRequest ID: (?P<requestId>[-._~!$&'()*+,;=:@%a-zA-Z0-9]*))?"

var SIWS_RESOURCES = fmt.Sprintf("(\\nResources:(?P<resources>(\\n- %s?)+))?$", SIWS_URI)

var SIWS_MESSAGE = regexp.MustCompile(fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%s%s",
	SIWS_DOMAIN,
	SIWS_ADDRESS,
	SIWS_STATEMENT,
	SIWS_URI_LINE,
	SIWS_VERSION,
	SIWS_CHAIN_ID,
	SIWS_NONCE,
	SIWS_ISSUED_AT,
	SIWS_EXPIRATION_TIME,
	SIWS_NOT_BEFORE,
	SIWS_REQUEST_ID,
	SIWS_RESOURCES))
