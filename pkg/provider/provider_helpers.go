package provider

import (
	"crypto/rsa"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/snowflakedb/gosnowflake"
	"github.com/youmark/pkcs8"
	"golang.org/x/crypto/ssh"
)

type authenticationType string

const (
	authenticationTypeSnowflake           authenticationType = "SNOWFLAKE"
	authenticationTypeOauth               authenticationType = "OAUTH"
	authenticationTypeExternalBrowser     authenticationType = "EXTERNALBROWSER"
	authenticationTypeOkta                authenticationType = "OKTA"
	authenticationTypeJwtLegacy           authenticationType = "JWT"
	authenticationTypeJwt                 authenticationType = "SNOWFLAKE_JWT"
	authenticationTypeTokenAccessor       authenticationType = "TOKENACCESSOR"
	authenticationTypeUsernamePasswordMfa authenticationType = "USERNAMEPASSWORDMFA"
)

var allAuthenticationTypes = []authenticationType{
	authenticationTypeSnowflake,
	authenticationTypeOauth,
	authenticationTypeExternalBrowser,
	authenticationTypeOkta,
	authenticationTypeJwt,
	authenticationTypeTokenAccessor,
	authenticationTypeUsernamePasswordMfa,
}

func toAuthenticatorType(s string) (gosnowflake.AuthType, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(authenticationTypeSnowflake):
		return gosnowflake.AuthTypeSnowflake, nil
	case string(authenticationTypeOauth):
		return gosnowflake.AuthTypeOAuth, nil
	case string(authenticationTypeExternalBrowser):
		return gosnowflake.AuthTypeExternalBrowser, nil
	case string(authenticationTypeOkta):
		return gosnowflake.AuthTypeOkta, nil
	case string(authenticationTypeJwt), string(authenticationTypeJwtLegacy):
		return gosnowflake.AuthTypeJwt, nil
	case string(authenticationTypeTokenAccessor):
		return gosnowflake.AuthTypeTokenAccessor, nil
	case string(authenticationTypeUsernamePasswordMfa):
		return gosnowflake.AuthTypeUsernamePasswordMFA, nil
	default:
		return gosnowflake.AuthType(0), fmt.Errorf("invalid authenticator type: %s", s)
	}
}

type protocol string

const (
	protocolHttp  protocol = "HTTP"
	protocolHttps protocol = "HTTPS"
)

var allProtocols = []protocol{
	protocolHttp,
	protocolHttps,
}

func toProtocol(s string) (protocol, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(protocolHttp):
		return protocolHttp, nil
	case string(protocolHttps):
		return protocolHttps, nil
	default:
		return "", fmt.Errorf("invalid protocol: %s", s)
	}
}

func getPrivateKey(privateKeyPath, privateKeyString, privateKeyPassphrase string) (*rsa.PrivateKey, error) {
	if privateKeyPath == "" && privateKeyString == "" {
		return nil, nil
	}
	privateKeyBytes := []byte(privateKeyString)
	var err error
	if len(privateKeyBytes) == 0 && privateKeyPath != "" {
		privateKeyBytes, err = readFile(privateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("private Key file could not be read err = %w", err)
		}
	}
	return parsePrivateKey(privateKeyBytes, []byte(privateKeyPassphrase))
}

func readFile(privateKeyPath string) ([]byte, error) {
	expandedPrivateKeyPath, err := homedir.Expand(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("invalid Path to private key err = %w", err)
	}

	privateKeyBytes, err := os.ReadFile(expandedPrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("could not read private key err = %w", err)
	}

	if len(privateKeyBytes) == 0 {
		return nil, errors.New("private key is empty")
	}

	return privateKeyBytes, nil
}

func parsePrivateKey(privateKeyBytes []byte, passhrase []byte) (*rsa.PrivateKey, error) {
	privateKeyBlock, _ := pem.Decode(privateKeyBytes)
	if privateKeyBlock == nil {
		return nil, fmt.Errorf("could not parse private key, key is not in PEM format")
	}

	if privateKeyBlock.Type == "ENCRYPTED PRIVATE KEY" {
		if len(passhrase) == 0 {
			return nil, fmt.Errorf("private key requires a passphrase, but private_key_passphrase was not supplied")
		}
		privateKey, err := pkcs8.ParsePKCS8PrivateKeyRSA(privateKeyBlock.Bytes, passhrase)
		if err != nil {
			return nil, fmt.Errorf("could not parse encrypted private key with passphrase, only ciphers aes-128-cbc, aes-128-gcm, aes-192-cbc, aes-192-gcm, aes-256-cbc, aes-256-gcm, and des-ede3-cbc are supported err = %w", err)
		}
		return privateKey, nil
	}

	privateKey, err := ssh.ParseRawPrivateKey(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse private key err = %w", err)
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("privateKey not of type RSA")
	}
	return rsaPrivateKey, nil
}

type GetRefreshTokenResponseBody struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func GetAccessTokenWithRefreshToken(
	tokenEndPoint string,
	clientID string,
	clientSecret string,
	refreshToken string,
	redirectURI string,
) (string, error) {
	client := &http.Client{}

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("redirect_uri", redirectURI)
	body := strings.NewReader(data.Encode())

	request, err := http.NewRequest("POST", tokenEndPoint, body)
	if err != nil {
		return "", fmt.Errorf("request to the endpoint could not be completed %w", err)
	}
	request.SetBasicAuth(clientID, clientSecret)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("response status returned an err = %w", err)
	}
	if response.StatusCode != 200 {
		return "", fmt.Errorf("response status code: %s: %s err = %w", strconv.Itoa(response.StatusCode), http.StatusText(response.StatusCode), err)
	}
	defer response.Body.Close()
	dat, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("response body was not able to be parsed err = %w", err)
	}
	var result GetRefreshTokenResponseBody
	err = json.Unmarshal(dat, &result)
	if err != nil {
		return "", fmt.Errorf("error parsing JSON from Snowflake err = %w", err)
	}
	return result.AccessToken, nil
}

func envNameFieldDescription(description, envName string) string {
	return fmt.Sprintf("%s Can also be sourced from the `%s` environment variable.", description, envName)
}
