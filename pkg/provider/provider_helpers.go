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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/go-homedir"
	"github.com/snowflakedb/gosnowflake"
	"github.com/youmark/pkcs8"
	"golang.org/x/crypto/ssh"
)

func mergeSchemas(schemaCollections ...map[string]*schema.Resource) map[string]*schema.Resource {
	out := map[string]*schema.Resource{}
	for _, schemaCollection := range schemaCollections {
		for name, s := range schemaCollection {
			out[name] = s
		}
	}
	return out
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

func toAuthenticatorType(authenticator string) gosnowflake.AuthType {
	switch authenticator {
	case "Snowflake":
		return gosnowflake.AuthTypeSnowflake
	case "OAuth":
		return gosnowflake.AuthTypeOAuth
	case "ExternalBrowser":
		return gosnowflake.AuthTypeExternalBrowser
	case "Okta":
		return gosnowflake.AuthTypeOkta
	case "JWT":
		return gosnowflake.AuthTypeJwt
	case "TokenAccessor":
		return gosnowflake.AuthTypeTokenAccessor
	case "UsernamePasswordMFA":
		return gosnowflake.AuthTypeUsernamePasswordMFA
	default:
		return gosnowflake.AuthTypeSnowflake
	}
}

func getInt64Env(key string, defaultValue int64) int64 {
	s := os.Getenv(key)
	if s == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return int64(i)
}

func getBoolEnv(key string, defaultValue bool) bool {
	s := strings.ToLower(os.Getenv(key))
	if s == "" {
		return defaultValue
	}
	switch s {
	case "true", "1":
		return true
	case "false", "0":
		return false
	default:
		return defaultValue
	}
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
