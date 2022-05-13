package runner

import (
	"github.com/Web3Auth/siws-go/pkg/message"
	"github.com/Web3Auth/siws-go/pkg/types"
	utils "github.com/Web3Auth/siws-go/pkg/utils"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func InitMessage(domain, address, uri, version string, options map[string]interface{}) (*message.Message, error) {
	validateURI, err := url.Parse(uri)
	if err != nil {
		return nil, &types.InvalidMessage{"Invalid format for field `uri`"}
	}

	var statement *string
	if val, ok := options["statement"]; ok {
		value := val.(string)
		statement = &value
	}

	var nonce string
	if val, ok := utils.IsStringAndNotEmpty(options, "nonce"); ok {
		nonce = *val
	} else {
		return nil, &types.InvalidMessage{"Missing or empty `nonce` property"}
	}

	var chainId int
	if val, ok := options["chainId"]; ok {
		switch val.(type) {
		case int:
			chainId = val.(int)
		case string:
			parsed, err := strconv.Atoi(val.(string))
			if err != nil {
				return nil, &types.InvalidMessage{"Invalid format for field `chainId`, must be an integer"}
			}
			chainId = parsed
		default:
			return nil, &types.InvalidMessage{"`chainId` must be a string or a integer"}
		}
	} else {
		chainId = 1
	}

	var issuedAt string
	timestamp, err := utils.ParseTimestamp(options, "issuedAt")
	if err != nil {
		return nil, err
	}

	if timestamp != nil {
		issuedAt = *timestamp
	} else {
		issuedAt = time.Now().UTC().Format(time.RFC3339)
	}

	var expirationTime *string
	timestamp, err = utils.ParseTimestamp(options, "expirationTime")
	if err != nil {
		return nil, err
	}

	if timestamp != nil {
		expirationTime = timestamp
	}

	var notBefore *string
	timestamp, err = utils.ParseTimestamp(options, "notBefore")
	if err != nil {
		return nil, err
	}

	if timestamp != nil {
		notBefore = timestamp
	}

	var requestID *string
	if val, ok := utils.IsStringAndNotEmpty(options, "requestId"); ok {
		requestID = val
	}

	var resources []string
	if val, ok := options["resources"]; ok {
		switch val.(type) {
		case []string:
			resources = val.([]string)
		default:
			return nil, &types.InvalidMessage{"`resources` must be a []string"}
		}
	}
	return &message.Message{
		Payload: message.Payload{
			Domain:  domain,
			Address: address,
			Uri:     *validateURI,
			Version: version,

			Statement: statement,
			Nonce:     nonce,
			ChainID:   chainId,

			IssuedAt:       issuedAt,
			ExpirationTime: expirationTime,
			NotBefore:      notBefore,

			RequestID: requestID,
			Resources: resources,
		},
	}, nil
}

func parseMessage(message string) (map[string]interface{}, error) {
	match := types.SIWS_MESSAGE.FindStringSubmatch(message)
	if match == nil {
		return nil, &types.InvalidMessage{"Message could not be parsed"}
	}

	result := make(map[string]interface{})
	for i, name := range types.SIWS_MESSAGE.SubexpNames() {
		if i != 0 && name != "" && match[i] != "" {
			result[name] = match[i]
		}
	}

	if val, ok := result["resources"]; ok {
		result["resources"] = strings.Split(val.(string), "\n- ")[1:]
	}

	return result, nil
}

func ParseMessage(message string) (*message.Message, error) {
	result, err := parseMessage(message)
	if err != nil {
		return nil, err
	}

	parsed, err := InitMessage(
		result["domain"].(string),
		result["address"].(string),
		result["uri"].(string),
		result["version"].(string),
		result,
	)

	if err != nil {
		return nil, err
	}

	return parsed, nil
}