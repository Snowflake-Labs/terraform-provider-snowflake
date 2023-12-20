## Start Docker-Compose

```
docker-compose up
```

the keycloak server will be running on localhost:8080, and the admin username password is admin//admin

## Create a Test User
First create a test user in the cloudeng3 account and grant it SECURITYADMIN:

```
CREATE OR REPLACE USER "test" LOGIN_NAME = 'test@snowflake.com' EMAIL = 'test@snowflake.com';
GRANT ROLE SECURITYADMIN TO USER "test";
```

This will be used to test SSO using SAML and OAuth later. It is very important that the LOGIN_NAME be the email, as this is how SAML respose NameID is used to match with the user.

## Create a SAML Integration

In the cloudeng3 account, run the following commands to create the SAML Security Integration:

```
CREATE OR REPLACE SECURITY INTEGRATION "keycloak_saml"
    TYPE = SAML2
    ENABLED = TRUE
    SAML2_ENABLE_SP_INITIATED = TRUE
    SAML2_ISSUER = 'http://localhost:8080/realms/snowflake'
    SAML2_SSO_URL = 'http://localhost:8080/realms/snowflake/protocol/saml'
    SAML2_PROVIDER = 'Custom'
    SAML2_X509_CERT = 'MIICoTCCAYkCBgGMRnSjkTANBgkqhkiG9w0BAQsFADAUMRIwEAYDVQQDDAlzbm93Zmxha2UwHhcNMjMxMjA3MjI0MzE4WhcNMzMxMjA3MjI0NDU4WjAUMRIwEAYDVQQDDAlzbm93Zmxha2UwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCCtHH9Ws9xhYddJmcwhsgj7hxgj5iL8q+1NpsdyPYYTM8lfnfjUEGzfYsZzFpxqMr5l/m9fhvkmmoOlwdo2d8BDIxQNlAZ/ChjthFiaxCr8SiRosk5lK7riylSA6po6iToq9fi4ehV0j66ulFfLcZqeTDIXzO9eLq9YpAmlTaBMr6tmOSlkCCHr8cpDqJLPnN3Vb4mVsHOu5RXVKauqDt7nN1TuO0ZvultIFPHnk7o4Yv83kyegTyNXhO/kXN44mmufpG6kg+h8FbOscp+fAJQto91r42HtsEG2X+qkzzDqpzOxZf7reFtOn6KyTFIFo0N987N5srIva4G0F7kbqHHAgMBAAEwDQYJKoZIhvcNAQELBQADggEBAAiwFjINNHDOkqusHl82+a/XaMP4o6mKBiMdYrodeKpE193NhaHCh2AVlVCxbOcTtaYvBauF1a8q9I59CDn/hCgXV+dtjA3fRgdUXJMnk8Wf81RRPjvLb1VJxgekFwczOChXE5bDmJ7hPyPA7mjrbmJd4q88yL2UucL5meO/Hhyw4ZvJ5+8DkWv9YL+cLZlBAmPNw8CAzK0AQ8pNodPMwrbka88eBE3e4tnWBmrpD8/hWGYgjYssXELYP3zYYEQPoMvvf2NtlYkZr4Wy0A8C36WqeM8GbrRavS3K3T791htuauAxlTnWjT9cTEBViCLHAQdre17aonXkOM0y24RYUMA='
    SAML2_SNOWFLAKE_ACS_URL = 'https://iya62698.snowflakecomputing.com/fed/login'
    SAML2_SNOWFLAKE_ISSUER_URL = 'https://iya62698.snowflakecomputing.com'
    SAML2_SP_INITIATED_LOGIN_PAGE_LABEL = 'Keycloak SSO'
```

Note: if you get a 403 bad request in the next step, it is likely because the X509 cert needs to be updated. The X509 cert comes from http://localhost:8080/realms/snowflake/protocol/saml/descriptor
under the tag: <ds:X509Certificate>. It appears to not be constant, and changes each time docker-compose goes up or down.


```
ALTER SECURITY INTEGRATION "keycloak_saml"
  SET SAML2_X509_CERT = 'MIICoTCCAYkCBgGMgiaheTANBgkqhkiG9w0BAQsFADAUMRIwEAYDVQQDDAlzbm93Zmxha2UwHhcNMjMxMjE5MTI1NTE5WhcNMzMxMjE5MTI1NjU5WjAUMRIwEAYDVQQDDAlzbm93Zmxha2UwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCq+N/fg+2TFmmlh6YxGIdAqKNKxXKPGZPHSHdY1fOuhycn5zNBF1pN5Mv++x4P6hgvaXaijjOSGKz6vN3WIqfn79Jud92d+qUO2uC/2d0GDj47yEyDGpsGB7DE8dlKuxKg22OjufvmlVj8+xUe/k1kN+cX9vy0X+7kAX1rMTSu4m9Pdd/is8PXmim9rx5hyrbLHBeNUg+EBjnhzWqHk6JTIA4Bj62PK5BKJU57TSRhpqBb+eA+IZt1OrAskh7SSPACIkcotaeKfD5iPthMyfMgL/Nuh6GN1okFfThm7FLgZrApU7RcEH9Ztp3P5LOSSgxRwyOkMw6lo+KH4+wn3NVxAgMBAAEwDQYJKoZIhvcNAQELBQADggEBAFt6yJukZxlphADHGaK0vMcMP4CkA64/DVC7NIx00wA3JQKkM1f0LArdSBzbsL+MI2ASyC9QiNwGw+dhURsIzL8l8QX/CANpZNIj+8ZX0F1yOYqVr2MpTqLw6DxTy4s+efdEKqiOVNZnl56QmE23ma5nAennp+YPdMyW0vKPAkf9jhXkpKB0TjrLyhJ7yz7OrBEmfYSPG4hDMnNHQFa2BQkouK1LifU0NCJb2fYaJxF8P/rhVwvYpoBoK9qkFuL40JnSfWlldIjLVI5Uw8O99jpzjQcsjgf/XAj2AmHdvw3M6ek35m5vrUwkIrdeNgHBM4LT331QAJCgy4Iw89h2pNs='
```

## Test SAML Integration using Browser

In the browser navigate to https://sfdevrel-cloud_engineering3.snowflakecomputing.com/console/login#/ and click "Sign in using Keycloak SSO". This will redirect you to 

The username // password in keycloak (which is also in the realm-export.json) is test@snowflake.com // 1234 

## Test Snowflake Terraform Provider using Browser Auth

The following minimal Terraform configuration can be used to test Browser based authentication using the keycloak IdP
```
terraform {
  required_providers {
    snowflake = {
      source  = "Snowflake-Labs/snowflake"
      version = "0.75.0"
    }
  }
}

provider "snowflake" {
  user = "TEST"
  role = "SECURITYADMIN"
  account       = "sfdevrel-cloud_engineering3"
  authenticator = "ExternalBrowser"
}

data "snowflake_current_account" "current_account" {}

output "current_account" {
  value = data.snowflake_current_account.current_account
}

```
## Create OAuth Security Integration in Snowflake

In the cloudeng3 account, run the following commands to create the OAuth Security Integration:

```
CREATE OR REPLACE SECURITY INTEGRATION "keycloak"
    TYPE = EXTERNAL_OAUTH
    ENABLED = true
    EXTERNAL_OAUTH_TYPE = CUSTOM
    EXTERNAL_OAUTH_ISSUER = 'http://localhost:8080/realms/snowflake'
    EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM = 'upn'
    EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE = 'login_name'
    EXTERNAL_OAUTH_RSA_PUBLIC_KEY = 'MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAgrRx/VrPcYWHXSZnMIbII+4cYI+Yi/KvtTabHcj2GEzPJX5341BBs32LGcxacajK+Zf5vX4b5JpqDpcHaNnfAQyMUDZQGfwoY7YRYmsQq/EokaLJOZSu64spUgOqaOok6KvX4uHoVdI+urpRXy3GankwyF8zvXi6vWKQJpU2gTK+rZjkpZAgh6/HKQ6iSz5zd1W+JlbBzruUV1Smrqg7e5zdU7jtGb7pbSBTx55O6OGL/N5MnoE8jV4Tv5FzeOJprn6RupIPofBWzrHKfnwCULaPda+Nh7bBBtl/qpM8w6qczsWX+63hbTp+iskxSBaNDffOzebKyL2uBtBe5G6hxwIDAQAB'
    EXTERNAL_OAUTH_ALLOWED_ROLES_LIST = ('SECURITYADMIN')
    EXTERNAL_OAUTH_AUDIENCE_LIST = ('https://sfdevrel-cloud_engineering3.snowflakecomputing.com/');
```

## Generate an OAuth Token + Refresh Token

You can perform a curl request to get the OAuth token 

```
curl -k -v -X POST -H 'Content-type: application/x-www-form-urlencoded' \
-d "client_id=cloudeng3&grant_type=password&username=test@snowflake.com&password=1234" \
http://localhost:8080/realms/snowflake/protocol/openid-connect/token
```

This will return a JSON response with the access_token and refresh_token. The access_token can be used to authenticate with the Snowflake Terraform Provider. Example is shown below.

```
{"access_token":"eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJfQTRjbzdoRzNIcWJhSWpvYmhlLUlna2RUSWFfejY0Vjc0a2d0VlA1a0tZIn0.eyJleHAiOjE3MDMwNzQwOTQsImlhdCI6MTcwMzA3Mzc5NCwianRpIjoiNjdjOTczYmItOTZiZS00OTc1LTg0ZjMtODA2NTk3YWExNzU4IiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL3JlYWxtcy9zbm93Zmxha2UiLCJhdWQiOiJodHRwczovL3NmZGV2cmVsLWNsb3VkX2VuZ2luZWVyaW5nMy5zbm93Zmxha2Vjb21wdXRpbmcuY29tLyIsInN1YiI6InRlc3QtdXNlciIsInR5cCI6IkJlYXJlciIsImF6cCI6ImNsb3VkZW5nMyIsInNlc3Npb25fc3RhdGUiOiIwMDBiODAxNy0xNzIzLTQxZjAtODA5ZC0xNGM5NDIyMThiODciLCJzY29wZSI6InNub3dmbGFrZSIsInNpZCI6IjAwMGI4MDE3LTE3MjMtNDFmMC04MDlkLTE0Yzk0MjIxOGI4NyIsInNjcCI6WyJzZXNzaW9uOnJvbGU6c2VjdXJpdHlhZG1pbiJdLCJ1cG4iOiJ0ZXN0QHNub3dmbGFrZS5jb20ifQ.a0blq7YrAl3Z0a1aKJ-SOZjaHcy_L-ELFapLdrgiup9w3ajC_GRtLm_FKCJAzC5WphH4MQm5eCmlRLfbkxDfnoTnfikQn-OXpHQZLkQanjWXetMRA74E8YtJSQO9BlSr8riEkCoHqt0E1j2zDpvlH97vLIZ-AxATp6eA7wMEnJvQVEODBcyECWd5u8yeZIO8OLvZ6eFV7iL5VMLA56vtWAcgo5gVRBmBDjngIs0v46uu4tGWi02iGlyEJ_bUiU4TmuXIoMb02XUvZhdzfomBAvBylAjCayJDa-XXWVNOisfvrEvTAP4k0KEn5FN7H6Qxi2eLzD5vkd-lV4kED2Z7ZA","expires_in":300,"refresh_expires_in":1800,"refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJmNjVmMTY3Ni03ZjUwLTRmNDEtOWY3Zi0wNjNiMDc1MTYzMWMifQ.eyJleHAiOjE3MDMwNzU1OTQsImlhdCI6MTcwMzA3Mzc5NCwianRpIjoiMjNlODI4MzktZjJiMS00NGM4LTg5NDEtZTQ2NTdkNWVhNGJkIiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL3JlYWxtcy9zbm93Zmxha2UiLCJhdWQiOiJodHRwOi8vbG9jYWxob3N0OjgwODAvcmVhbG1zL3Nub3dmbGFrZSIsInN1YiI6InRlc3QtdXNlciIsInR5cCI6IlJlZnJlc2giLCJhenAiOiJjbG91ZGVuZzMiLCJzZXNzaW9uX3N0YXRlIjoiMDAwYjgwMTctMTcyMy00MWYwLTgwOWQtMTRjOTQyMjE4Yjg3Iiwic2NvcGUiOiJzbm93Zmxha2UiLCJzaWQiOiIwMDBiODAxNy0xNzIzLTQxZjAtODA5ZC0xNGM5NDIyMThiODcifQ.gr9pEx7z3vU29P6cKx6UZdwxW1GRq0EEr85LBImNCBQ","token_type":"Bearer","not-before-policy":0,"session_state":"000b8017-1723-41f0-809d-14c942218b87","scope":"snowflake"}
```

The following Terraform configuration can be used to authenticate using the OAuth access token

```
```
terraform {
  required_providers {
    snowflake = {
      source  = "Snowflake-Labs/snowflake"
      version = "0.75.0"
    }
  }
}

provider "snowflake" {
  user = "TEST"
  role = "SECURITYADMIN"
  account       = "sfdevrel-cloud_engineering3"
  authenticator = "OAuth"
  token = "<oauth_access_token>"
}

data "snowflake_current_account" "current_account" {}

output "current_account" {
  value = data.snowflake_current_account.current_account
}
```
