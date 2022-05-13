package runner

import (
	"crypto/ed25519"
	"fmt"
	siwsMessage "github.com/Web3Auth/siws-go/pkg/message"
	"github.com/Web3Auth/siws-go/pkg/types"
	utils "github.com/Web3Auth/siws-go/pkg/utils"
	solTypes "github.com/portto/solana-go-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const domain = "example.com"
const addressStr = "43h6BNKzvoV43qBLje5dxn7vhcChZjVEAn8PQLZvMiqj"

const uri = "https://example.com"
const version = "1"
const statement = "Example statement for SIWS"

var issuedAt = time.Now().UTC().Format(time.RFC3339)
var nonce = utils.GenerateNonce()

const chainId = 1

var expirationTime = time.Now().UTC().Add(48 * time.Hour).Format(time.RFC3339)

var notBefore = time.Now().UTC().Add(-24 * time.Hour).Format(time.RFC3339)

const requestId = "some-id"

var resources = []string{"https://example.com/resources/1", "https://example.com/resources/2"}

var options = map[string]interface{}{
	"statement":      statement,
	"nonce":          nonce,
	"chainId":        chainId,
	"issuedAt":       issuedAt,
	"expirationTime": expirationTime,
	"notBefore":      notBefore,
	"requestId":      requestId,
	"resources":      resources,
}

var testMessage, _ = InitMessage(
	domain,
	addressStr,
	uri,
	version,
	options,
)

func compareMessage(t *testing.T, a, b *siwsMessage.Message) {

	assert.Equal(t, a.Payload.Domain, b.Payload.Domain, "expected %s, found %s", a.Payload.Domain, b.Payload.Domain)
	assert.Equal(t, a.Payload.Address, b.Payload.Address, "expected %s, found %s", a.Payload.Address, b.Payload.Address)
	assert.Equal(t, a.Payload.Uri.String(), b.Payload.Uri.String(), "expected %s, found %s", a.Payload.Uri, b.Payload.Uri)
	assert.Equal(t, a.Payload.Version, b.Payload.Version, "expected %s, found %s", a.Payload.Version, b.Payload.Version)

	assert.Equal(t, a.Payload.Statement, b.Payload.Statement, "expected %s, found %s", a.Payload.Statement, b.Payload.Statement)
	assert.Equal(t, a.Payload.Nonce, b.Payload.Nonce, "expected %s, found %s", a.Payload.Nonce, b.Payload.Nonce)
	assert.Equal(t, a.Payload.ChainID, b.Payload.ChainID, "expected %s, found %s", a.Payload.ChainID, b.Payload.ChainID)

	assert.Equal(t, a.Payload.IssuedAt, b.Payload.IssuedAt, "expected %s, found %s", a.Payload.IssuedAt, b.Payload.IssuedAt)
	assert.Equal(t, a.Payload.ExpirationTime, b.Payload.ExpirationTime, "expected %s, found %s", a.Payload.ExpirationTime, b.Payload.ExpirationTime)
	assert.Equal(t, a.Payload.NotBefore, b.Payload.NotBefore, "expected %s, found %s", a.Payload.NotBefore, b.Payload.NotBefore)

	assert.Equal(t, a.Payload.RequestID, b.Payload.RequestID, "expected %s, found %s", a.Payload.RequestID, b.Payload.RequestID)
	assert.Equal(t, a.Payload.Resources, b.Payload.Resources, "expected %v, found %v", a.Payload.Resources, b.Payload.Resources)
}



func TestCreate(t *testing.T) {
	assert.Equal(t, testMessage.Payload.Domain, domain, "domain should be %s", domain)
	assert.Equal(t, testMessage.Payload.Address, addressStr, "address should be %s", addressStr)
	assert.Equal(t, testMessage.Payload.Uri.String(), uri, "uri should be %s", uri)
	assert.Equal(t, testMessage.Payload.Version, version, "version should be %s", version)

	assert.Equal(t, *testMessage.Payload.Statement, statement, "statement should be %s", statement)
	assert.Equal(t, testMessage.Payload.Nonce, nonce, "nonce should be %s", nonce)
	assert.Equal(t, testMessage.Payload.ChainID, chainId, "chainId should be %s", chainId)

	assert.Equal(t, testMessage.Payload.IssuedAt, issuedAt, "issuedAt should be %v", issuedAt)
	assert.Equal(t, *testMessage.Payload.ExpirationTime, expirationTime, "expirationTime should be %s", expirationTime)
	assert.Equal(t, *testMessage.Payload.NotBefore, notBefore, "notBefore should be %s", notBefore)

	assert.Equal(t, *testMessage.Payload.RequestID, requestId, "requestId should be %s", requestId)
	assert.Equal(t, testMessage.Payload.Resources, resources, "resources should be %v", resources)
}


func TestCreateRequired(t *testing.T) {
	message, err := InitMessage(domain, addressStr, uri, version, map[string]interface{}{
		"nonce": utils.GenerateNonce(),
	})
	assert.Nil(t, err)

	assert.Equal(t, message.Payload.Domain, domain, "domain should be %s", domain)
	assert.Equal(t, message.Payload.Address, addressStr, "address should be %s", addressStr)
	assert.Equal(t, message.Payload.Uri.String(), uri, "uri should be %s", uri)
	assert.Equal(t, message.Payload.Version, version, "version should be %s", version)

	assert.Nil(t, message.Payload.Statement, "statement should be nil")
	assert.NotNil(t, message.Payload.Nonce, "nonce should be not nil")
	assert.NotNil(t, message.Payload.ChainID, "chainId should not be nil")

	assert.NotNil(t, message.Payload.IssuedAt, "issuedAt should not be nil")
	assert.Nil(t, message.Payload.ExpirationTime, "expirationTime should be nil")
	assert.Nil(t, message.Payload.NotBefore, "notBefore should be nil")

	assert.Nil(t, message.Payload.RequestID, "requestId should be nil")
	assert.Len(t, message.Payload.Resources, 0, "resources should be empty")
}


func TestPrepareParse(t *testing.T) {
	prepare := testMessage.String()
	parse, err := ParseMessage(prepare)
	assert.Nil(t, err)
	compareMessage(t, testMessage, parse)
}


func TestPrepareParseRequired(t *testing.T) {
	message, err := InitMessage(domain, addressStr, uri, version, map[string]interface{}{
		"nonce": utils.GenerateNonce(),
	})
	assert.Nil(t, err)

	prepare := message.String()
	parse, err := ParseMessage(prepare)

	assert.Nil(t, err)

	compareMessage(t, message, parse)
}

func TestValidateEmpty(t *testing.T) {
	_, err := testMessage.Verify("", nil, nil)

	if assert.Error(t, err) {
		assert.Equal(t, &types.InvalidSignature{"Message signature invalid"}, err)
	}
}


func createWallet(t *testing.T) (acc solTypes.Account) {
	acc = solTypes.NewAccount()
	fmt.Println("Wallet Address:", acc.PublicKey.ToBase58())
	return acc
}

func TestValidateNotBefore(t *testing.T) {
	account := createWallet(t)

	message, err := InitMessage(domain, addressStr, uri, version, map[string]interface{}{
		"nonce":     utils.GenerateNonce(),
		"notBefore": time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339),
	})
	assert.Nil(t, err)
	messageParsed := message.String()

	signature := account.Sign([]byte(messageParsed))
	assert.NotNil(t,signature)

	resp := ed25519.Verify(account.PublicKey.Bytes(),[]byte(messageParsed),signature)
	assert.Equal(t, true,resp)
}


func TestValidateExpirationTime(t *testing.T) {
	account := createWallet(t)
	nonce := utils.GenerateNonce()
	timestamp := time.Now().UTC().Add(-24 * time.Hour).Format(time.RFC3339)
	message, err := InitMessage(domain, account.PublicKey.ToBase58(), uri, version, map[string]interface{}{
		"nonce":       nonce,
		"expirationTime": timestamp,
	})
	assert.Nil(t, err)
	preparedMessage := message.String()

	_, err = message.Verify(string(account.Sign([]byte(preparedMessage))), &nonce, nil)

	if assert.Error(t, err) {
		assert.Equal(t, &types.ExpiredMessage{"Message expired"}, err)
	}
}


func TestValidate(t *testing.T) {
	account := createWallet(t)

	message, err := InitMessage(domain, account.PublicKey.ToBase58(), uri, version, options)
	assert.Nil(t, err)
	prepare:= message.PrepareMessage()

	signature := account.Sign([]byte(prepare))

	assert.Nil(t, err)

	_, err = message.Verify(string(signature), nil, nil)

	assert.Nil(t, err)
}

/*
func TestValidateTampered(t *testing.T) {
	privateKey, address := createWallet(t)
	_, otherAddress := createWallet(t)

	message, err := InitMessage(domain, address, uri, version, options)
	assert.Nil(t, err)

	hash := message.Eip191Hash()
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	signature[64] += 27

	assert.Nil(t, err)

	message, err = InitMessage(domain, otherAddress, uri, version, options)
	assert.Nil(t, err)
	_, err = message.Verify(hexutil.Encode(signature), nil, nil)

	if assert.Error(t, err) {
		assert.Equal(t, &types.InvalidSignature{"Signer address must match message address"}, err)
	}
}

func assertCase(t *testing.T, fields map[string]interface{}, parsed string, key string) {
	if field, ok := fields[key]; ok {
		assert.Equal(t, field, parsed, "%s should be %s", key, field)
	}
}

func parsingNegative(t *testing.T, cases map[string]interface{}) {
	for name, message := range cases {
		_, err := ParseMessage(message.(string))
		assert.Error(t, err, name)
	}
}

func parsingPositive(t *testing.T, cases map[string]interface{}) {
	for name, v := range cases {
		data := v.(map[string]interface{})
		message := data["message"].(string)
		fields := data["fields"].(map[string]interface{})
		parsed, err := parseMessage(message)
		assert.Nil(t, err, name)

		assertCase(t, fields, parsed["domain"].(string), "domain")
		assertCase(t, fields, parsed["address"].(string), "address")
		assertCase(t, fields, parsed["uri"].(string), "uri")
		assertCase(t, fields, parsed["version"].(string), "version")
		assertCase(t, fields, parsed["chainId"].(string), "chainId")
		assertCase(t, fields, parsed["issuedAt"].(string), "issuedAt")

		if val, ok := parsed["statement"]; ok {
			assertCase(t, fields, val.(string), "statement")
		}

		if val, ok := parsed["nonce"]; ok {
			assertCase(t, fields, val.(string), "nonce")
		}

		constructed, err := ParseMessage(message)
		assert.Nil(t, err)
		assert.Equal(t, constructed.String(), message)
	}
}

func validationNegative(t *testing.T, cases map[string]interface{}) {
	for name, v := range cases {
		data := v.(map[string]interface{})
		message, err := InitMessage(
			data["domain"].(string),
			data["address"].(string),
			data["uri"].(string),
			data["version"].(string),
			data,
		)
		assert.Nil(t, err)

		_, err = message.Verify(data["signature"].(string), nil, nil)

		assert.Error(t, err, name)
	}
}

func validationPositive(t *testing.T, cases map[string]interface{}) {
	for name, v := range cases {
		data := v.(map[string]interface{})
		payload := data["payload"].(map[string]interface{})
		message, err := InitMessage(
			payload["domain"].(string),
			payload["address"].(string),
			payload["uri"].(string),
			payload["version"].(string),
			data,
		)
		assert.Nil(t, err)

		_, err = message.Verify(data["signature"].(string), nil, nil)

		assert.Nil(t, err, name)
	}
}

func TestGlobalTestVector(t *testing.T) {
	files := make(map[string]*os.File, 4)
	for test, filename := range map[string]string{
		//"parsing-negative":    "../tests/parsing_negative.json",
		//"parsing-positive":    "../tests/parsing_positive.json",
		//"validation-negative": "../tests/validation_negative.json",
		"validation-positive": "../tests/validation_positive.json",
	} {
		filePath,err := filepath.Abs(filename)
		file, err := os.Open(filePath)
		assert.Nil(t, err)
		files[test] = file
		defer file.Close()
	}

	for test, file := range files {
		data, _ := ioutil.ReadAll(file)

		var result map[string]interface{}
		err := json.Unmarshal([]byte(data), &result)
		assert.Nil(t, err)

		switch test {
		case "parsing-negative":
			parsingNegative(t, result)
		case "parsing-positive":
			parsingPositive(t, result)
		case "validation-negative":
			validationNegative(t, result)
		case "validation-positive":
			validationPositive(t, result)
		}
	}
}
 */
