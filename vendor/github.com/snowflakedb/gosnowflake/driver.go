// Copyright (c) 2017-2018 Snowflake Computing Inc. All right reserved.

package gosnowflake

import (
	"database/sql"
	"database/sql/driver"
	"net/http"
	"strings"
)

// SnowflakeDriver is a context of Go Driver
type SnowflakeDriver struct{}

// Open creates a new connection.
func (d SnowflakeDriver) Open(dsn string) (driver.Conn, error) {
	glog.V(2).Info("Open")
	var err error
	sc := &snowflakeConn{
		SequeceCounter: 0,
	}
	sc.cfg, err = ParseDSN(dsn)
	if err != nil {
		sc.cleanup()
		return nil, err
	}
	st := SnowflakeTransport
	if sc.cfg.InsecureMode {
		// no revocation check with OCSP. Think twice when you want to enable this option.
		st = snowflakeInsecureTransport
	}
	if err != nil {
		return nil, err
	}
	// authenticate
	sc.rest = &snowflakeRestful{
		Host:     sc.cfg.Host,
		Port:     sc.cfg.Port,
		Protocol: sc.cfg.Protocol,
		Client: &http.Client{
			// request timeout including reading response body
			Timeout:   defaultClientTimeout,
			Transport: st,
		},
		Authenticator:       sc.cfg.Authenticator,
		LoginTimeout:        sc.cfg.LoginTimeout,
		RequestTimeout:      sc.cfg.RequestTimeout,
		FuncPost:            postRestful,
		FuncGet:             getRestful,
		FuncPostQuery:       postRestfulQuery,
		FuncPostQueryHelper: postRestfulQueryHelper,
		FuncRenewSession:    renewRestfulSession,
		FuncPostAuth:        postAuth,
		FuncCloseSession:    closeSession,
		FuncCancelQuery:     cancelQuery,
		FuncPostAuthSAML:    postAuthSAML,
		FuncPostAuthOKTA:    postAuthOKTA,
		FuncGetSSO:          getSSO,
	}
	var authData *authResponseMain
	var samlResponse []byte
	var proofKey []byte

	authenticator := strings.ToUpper(sc.cfg.Authenticator)
	glog.V(2).Infof("Authenticating via %v", authenticator)
	switch authenticator {
	case authenticatorExternalBrowser:
		samlResponse, proofKey, err = authenticateByExternalBrowser(
			sc.rest,
			sc.cfg.Authenticator,
			sc.cfg.Application,
			sc.cfg.Account,
			sc.cfg.User,
			sc.cfg.Password)
		if err != nil {
			sc.cleanup()
			return nil, err
		}
	case authenticatorOAuth:
	case authenticatorSnowflake:
	case authenticatorJWT:
		// Nothing to do, parameters needed for auth should be already set in sc.cfg
		break
	default:
		// this is actually okta, which is something misleading
		samlResponse, err = authenticateBySAML(
			sc.rest,
			sc.cfg.Authenticator,
			sc.cfg.Application,
			sc.cfg.Account,
			sc.cfg.User,
			sc.cfg.Password)
		if err != nil {
			sc.cleanup()
			return nil, err
		}
	}
	authData, err = authenticate(
		sc,
		samlResponse,
		proofKey)
	if err != nil {
		sc.cleanup()
		return nil, err
	}

	err = d.validateDefaultParameters(authData.SessionInfo.DatabaseName, &sc.cfg.Database)
	if err != nil {
		return nil, err
	}
	err = d.validateDefaultParameters(authData.SessionInfo.SchemaName, &sc.cfg.Schema)
	if err != nil {
		return nil, err
	}
	err = d.validateDefaultParameters(authData.SessionInfo.WarehouseName, &sc.cfg.Warehouse)
	if err != nil {
		return nil, err
	}
	err = d.validateDefaultParameters(authData.SessionInfo.RoleName, &sc.cfg.Role)
	if err != nil {
		return nil, err
	}
	sc.populateSessionParameters(authData.Parameters)
	sc.startHeartBeat()
	return sc, nil
}

func (d SnowflakeDriver) validateDefaultParameters(sessionValue string, defaultValue *string) error {
	if *defaultValue != "" && strings.ToLower(*defaultValue) != strings.ToLower(sessionValue) {
		return &SnowflakeError{
			Number:      ErrCodeObjectNotExists,
			SQLState:    SQLStateConnectionFailure,
			Message:     errMsgObjectNotExists,
			MessageArgs: []interface{}{*defaultValue},
		}
	}
	*defaultValue = sessionValue
	return nil
}

func init() {
	sql.Register("snowflake", &SnowflakeDriver{})
}
