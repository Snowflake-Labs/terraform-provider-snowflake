package provider

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/go-homedir"
	"github.com/snowflakedb/gosnowflake"
)

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
	upperCase := strings.ToUpper(s)
	switch upperCase {
	case string(protocolHttp),
		string(protocolHttps):
		return protocol(upperCase), nil
	default:
		return "", fmt.Errorf("invalid protocol: %s", s)
	}
}

type driverLogLevel string

const (
	// these values
	logLevelTrace   driverLogLevel = "trace"
	logLevelDebug   driverLogLevel = "debug"
	logLevelInfo    driverLogLevel = "info"
	logLevelPrint   driverLogLevel = "print"
	logLevelWarning driverLogLevel = "warning"
	logLevelError   driverLogLevel = "error"
	logLevelFatal   driverLogLevel = "fatal"
	logLevelPanic   driverLogLevel = "panic"
)

var allDriverLogLevels = []driverLogLevel{
	logLevelTrace,
	logLevelDebug,
	logLevelInfo,
	logLevelPrint,
	logLevelWarning,
	logLevelError,
	logLevelFatal,
	logLevelPanic,
}

func toDriverLogLevel(s string) (driverLogLevel, error) {
	lowerCase := strings.ToLower(s)
	switch lowerCase {
	case string(logLevelTrace),
		string(logLevelDebug),
		string(logLevelInfo),
		string(logLevelPrint),
		string(logLevelWarning),
		string(logLevelError),
		string(logLevelFatal),
		string(logLevelPanic):
		return driverLogLevel(lowerCase), nil
	default:
		return "", fmt.Errorf("invalid driver log level: %s", s)
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
	return sdk.ParsePrivateKey(privateKeyBytes, []byte(privateKeyPassphrase))
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

func handleStringField(d *schema.ResourceData, key string, field *string) error {
	if v, ok := d.GetOk(key); ok {
		*field = v.(string)
	}
	return nil
}

func handleBoolField(d *schema.ResourceData, key string, field *bool) error {
	if v, ok := d.GetOk(key); ok {
		*field = v.(bool)
	}
	return nil
}

func handleDurationInSecondsAttribute(d *schema.ResourceData, key string, field *time.Duration) error {
	if v, ok := d.GetOk(key); ok {
		*field = time.Second * time.Duration(int64(v.(int)))
	}
	return nil
}

func handleIntAttribute(d *schema.ResourceData, key string, field *int) error {
	if v, ok := d.GetOk(key); ok {
		*field = v.(int)
	}
	return nil
}

func handleBooleanStringAttribute(d *schema.ResourceData, key string, field *gosnowflake.ConfigBool) error {
	if v := d.Get(key).(string); v != provider.BooleanDefault {
		parsed, err := provider.BooleanStringToBool(v)
		if err != nil {
			return err
		}
		if parsed {
			*field = gosnowflake.ConfigBoolTrue
		} else {
			*field = gosnowflake.ConfigBoolFalse
		}
	}
	return nil
}
