package provider_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/provider"
	_ "github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/require"
)

func TestProvider(t *testing.T) {
	r := require.New(t)
	err := provider.Provider().InternalValidate()
	r.NoError(err)
}

// func TestConfigureProvider(t *testing.T) {
// 	// r := require.New(t)
// }

func TestDSN(t *testing.T) {
	type args struct {
		account,
		user,
		password string
		browserAuth bool
		region,
		role,
		host string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"simple", args{"acct", "user", "pass", false, "region", "role", ""},
			"user:pass@acct.region.snowflakecomputing.com:443?ocspFailOpen=true&region=region&role=role&validateDefaultParameters=true", false},
		{"us-west-2 special case", args{"acct2", "user2", "pass2", false, "us-west-2", "role2", ""},
			"user2:pass2@acct2.snowflakecomputing.com:443?ocspFailOpen=true&role=role2&validateDefaultParameters=true", false},
		{"customhostwregion", args{"acct3", "user3", "pass3", false, "", "role3", "zha123.us-east-1.privatelink.snowflakecomputing.com"},
			"user3:pass3@zha123.us-east-1.privatelink.snowflakecomputing.com:443?account=acct3&ocspFailOpen=true&role=role3&validateDefaultParameters=true", false},
		{"customhostignoreregion", args{"acct4", "user4", "pass4", false, "fakeregion", "role4", "zha1234.us-east-1.privatelink.snowflakecomputing.com"},
			"user4:pass4@zha1234.us-east-1.privatelink.snowflakecomputing.com:443?account=acct4&ocspFailOpen=true&role=role4&validateDefaultParameters=true", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := provider.DSN(tt.args.account, tt.args.user, tt.args.password, tt.args.browserAuth, "", "", "", "", tt.args.region, tt.args.role, tt.args.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("DSN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DSN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOAuthDSN(t *testing.T) {
	type args struct {
		account,
		user,
		oauthAccessToken,
		region,
		role string
	}
	pseudorandom_access_token := "ETMsjLOLvQ-C/bzGmmdvbEM/RSQFFX-a+sefbQeQoJqwdFNXZ+ftBIdwlasApA+/MItZLNRRW-rYJiEZMvAAdzpGLxaghIoww+vDOuIeAFBDUxTAY-I+qGbQOXipkNcmzwuAaugjYtlTjPXGjqKw-OSsVacQXzsQyAMnbMyUrbdhRQEETIqTAdMuDqJBeaSj+LMsKDXzLd-guSlm-mmv+="
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"simple_oauth", args{"acct", "user", pseudorandom_access_token, "region", "role"},
			"user:@acct.region.snowflakecomputing.com:443?authenticator=oauth&ocspFailOpen=true&region=region&role=role&token=ETMsjLOLvQ-C%2FbzGmmdvbEM%2FRSQFFX-a%2BsefbQeQoJqwdFNXZ%2BftBIdwlasApA%2B%2FMItZLNRRW-rYJiEZMvAAdzpGLxaghIoww%2BvDOuIeAFBDUxTAY-I%2BqGbQOXipkNcmzwuAaugjYtlTjPXGjqKw-OSsVacQXzsQyAMnbMyUrbdhRQEETIqTAdMuDqJBeaSj%2BLMsKDXzLd-guSlm-mmv%2B%3D&validateDefaultParameters=true", false},
		{"oauth_over_password", args{"acct", "user", pseudorandom_access_token, "region", "role"},
			"user:@acct.region.snowflakecomputing.com:443?authenticator=oauth&ocspFailOpen=true&region=region&role=role&token=ETMsjLOLvQ-C%2FbzGmmdvbEM%2FRSQFFX-a%2BsefbQeQoJqwdFNXZ%2BftBIdwlasApA%2B%2FMItZLNRRW-rYJiEZMvAAdzpGLxaghIoww%2BvDOuIeAFBDUxTAY-I%2BqGbQOXipkNcmzwuAaugjYtlTjPXGjqKw-OSsVacQXzsQyAMnbMyUrbdhRQEETIqTAdMuDqJBeaSj%2BLMsKDXzLd-guSlm-mmv%2B%3D&validateDefaultParameters=true", false},
		{"empty_token_no_password_errors_out", args{"acct", "user", "", "region", "role"},
			"", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := provider.DSN(tt.args.account, tt.args.user, "", false, "", "", "", tt.args.oauthAccessToken, tt.args.region, tt.args.role, "")

			if (err != nil) != tt.wantErr {
				t.Errorf("DSN() error = %v, dsn = %v, wantErr %v", err, got, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DSN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetOauthDATA(t *testing.T) {
	type param struct {
		refresh_token,
		redirect_url string
	}
	refresh_token := "ETMsDgAAAXdeJNwXABRBRVMvQ0JDL1BLQ1M1UGFwPu1hHM3UoUexZBtXW+0cE7KJx2yoUV0ysWu3HKwhJ1v/iEa1Np5EdjGDsBqedR15aFb8NstLTWDUoTJPuQNZRJTjJeuxrX/JUM3/wzcrKt2zDf6QIpkfLXuSlDH4VABeqsaRdl5z6bE9VJVgAUKgZwizwedHAt6pcJgFcQffYZPaY="
	redirect_url := "https://localhost.com"
	cases := []struct {
		name    string
		param   param
		want    string
		wantErr bool
	}{
		{"simpleData", param{refresh_token, redirect_url},
			"grant_type=refresh_token&redirect_uri=https%3A%2F%2Flocalhost.com&refresh_token=ETMsDgAAAXdeJNwXABRBRVMvQ0JDL1BLQ1M1UGFwPu1hHM3UoUexZBtXW%2B0cE7KJx2yoUV0ysWu3HKwhJ1v%2FiEa1Np5EdjGDsBqedR15aFb8NstLTWDUoTJPuQNZRJTjJeuxrX%2FJUM3%2FwzcrKt2zDf6QIpkfLXuSlDH4VABeqsaRdl5z6bE9VJVgAUKgZwizwedHAt6pcJgFcQffYZPaY%3D",
			false},
		{"errorData", param{"no_refresh_token", redirect_url},
			"grant_type=refresh_token&redirect_uri=https%3A%2F%2Flocalhost.com&refresh_token=no_refresh_token",
			false},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := provider.GetOauthData(tt.param.refresh_token, tt.param.redirect_url)
			want, err := url.ParseQuery(tt.want)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetData() error = %v, dsn = %v, wantErr %v", err, got, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("GetData() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestGetOauthResponse(t *testing.T) {
	type param struct {
		dataStuff,
		endpoint,
		clientid,
		clientscret string
	}
	dataStuff := "grant_type=refresh_token&redirect_uri=https%3A%2F%2Flocalhost.com&refresh_token=ETMsDgAAAXdeJNwXABRBRVMvQ0JDL1BLQ1M1UGFwPu1hHM3UoUexZBtXW%2B0cE7KJx2yoUV0ysWu3HKwhJ1v%2FiEa1Np5EdjGDsBqedR15aFb8NstLTWDUoTJPuQNZRJTjJeuxrX%2FJUM3%2FwzcrKt2zDf6QIpkfLXuSlDH4VABeqsaRdl5z6bE9VJVgAUKgZwizwedHAt6pcJgFcQffYZPaY%3D"
	endpoint := "https://example.snowflakecomputing.com/oauth/token-request"
	clientid := "nWsfd+gowithgoiwm1vJvGLckmLIMPS="
	clientsecret := "ThjKLFMD45wKIgVTecwVXguZrt+yHG1Ydth8eeQB34XU="
	cases := []struct {
		name    string
		param   param
		want    string
		wantErr bool
	}{
		{"simpleContent", param{dataStuff, endpoint, clientid, clientsecret},
			"application/x-www-form-urlencoded;charset=UTF-8",
			false},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := provider.GetOauthRequest(strings.NewReader(tt.param.dataStuff), tt.param.endpoint, tt.param.clientid, tt.param.clientscret)
			if err != nil {
				t.Errorf("GetOauthRequest() %v", err)
			}
			if !reflect.DeepEqual(got.Header.Get("Content-Type"), tt.want) {
				t.Errorf("GetResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestGetOauthAccessToken(t *testing.T) {
	type param struct {
		dataStuff,
		endpoint,
		clientid,
		clientsecret string
	}
	dataStuff := "grant_type=refresh_token&redirect_uri=https%3A%2F%2Flocalhost.com&refresh_token=ETMsDgAAAXdeJNwXABRBRVMvQ0JDL1BLQ1M1UGFwPu1hHM3UoUexZBtXW%2B0cE7KJx2yoUV0ysWu3HKwhJ1v%2FiEa1Np5EdjGDsBqedR15aFb8NstLTWDUoTJPuQNZRJTjJeuxrX%2FJUM3%2FwzcrKt2zDf6QIpkfLXuSlDH4VABeqsaRdl5z6bE9VJVgAUKgZwizwedHAt6pcJgFcQffYZPaY%3D"
	endpoint := "https://example.snowflakecomputing.com/oauth/token-request"
	clientid := "nWsfd+gowithgoiwm1vJvGLckmLIMPS="
	clientsecret := "ThjKLFMD45wKIgVTecwVXguZrt+yHG1Ydth8eeQB34XU="
	cases := []struct {
		name       string
		param      param
		want       string
		statuscode string
		wantTok    string
		wantErr    bool
	}{
		{"simpleAccessToken", param{dataStuff, endpoint, clientid, clientsecret},
			`{"access_token": "ABCDEFGHIabchefghiJKLMNOPQRjklmnopqrSTUVWXYZstuvwxyz","token_type": "Bearer","expires_in": 600}`,
			"200", "ABCDEFGHIabchefghiJKLMNOPQRjklmnopqrSTUVWXYZstuvwxyz", false},
		{"errorAccessToken", param{dataStuff, endpoint, clientid, clientsecret},
			"",
			"404", "", false},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				statusCODE, err := strconv.Atoi(tt.statuscode)
				if err != nil {
					t.Errorf("Invalid statuscode type %v", err)
				}
				return &http.Response{
					StatusCode: statusCODE,
					Body:       ioutil.NopCloser(bytes.NewBufferString(tt.want)),
					Header:     make(http.Header),
				}
			})
			req_got, err := provider.GetOauthRequest(strings.NewReader(tt.param.dataStuff), tt.param.endpoint, tt.param.clientid, tt.param.clientsecret)
			if err != nil {
				t.Errorf("GetOauthRequest() %v", err)
			}
			body, err := client.Do(req_got)
			if err != nil {
				t.Errorf("Body was not returned %v", err)
			}
			got, err := ioutil.ReadAll(body.Body)
			if err != nil {
				t.Errorf("Response body was not able to be parsed %v", err)
			}
			var result provider.Result
			json.Unmarshal([]byte(got), &result)
			if result.AccessToken != tt.wantTok {
				t.Errorf("TestGetAccessToken() = %v, want %v", result.AccessToken, tt.want)
			}
		})
	}
}
