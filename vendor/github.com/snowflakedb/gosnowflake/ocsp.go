// Copyright (c) 2017-2019 Snowflake Computing Inc. All right reserved.

package gosnowflake

import (
	"bufio"
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"golang.org/x/crypto/ocsp"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// caRoot includes the CA certificates.
var caRoot map[string]*x509.Certificate

// certPOol includes the CA certificates.
var certPool *x509.CertPool

// cacheDir is the location of OCSP response cache file
var cacheDir = ""

// cacheFileName is the file name of OCSP response cache file
var cacheFileName = ""

// cacheUpdated is true if the memory cache is updated
var cacheUpdated = true

// OCSPFailOpenMode is OCSP fail open mode. OCSPFailOpenTrue by default and may set to ocspModeFailClosed for fail closed mode
type OCSPFailOpenMode int

const (
	ocspFailOpenNotSet OCSPFailOpenMode = 0
	// OCSPFailOpenTrue represents OCSP fail open mode.
	OCSPFailOpenTrue   OCSPFailOpenMode = 1
	// OCSPFailOpenFalse represents OCSP fail closed mode.
	OCSPFailOpenFalse  OCSPFailOpenMode = 2
)
const (
	ocspModeFailOpen   = "FAIL_OPEN"
	ocspModeFailClosed = "FAIL_CLOSED"
	ocspModeInsecure   = "INSECURE"
)

// OCSP fail open mode
var ocspFailOpen = OCSPFailOpenTrue

const (
	// retryOCSPTimeout is the total timeout for OCSP checks.
	retryOCSPTimeout     = 60 * time.Second
	retryOCSPHTTPTimeout = 20 * time.Second
)

const (
	cacheFileBaseName = "ocsp_response_cache.json"
	// cacheExpire specifies cache data expiration time in seconds.
	cacheExpire           = float64(24 * 60 * 60)
	cacheServerURL        = "http://ocsp.snowflakecomputing.com"
	cacheServerEnabledEnv = "SF_OCSP_RESPONSE_CACHE_SERVER_ENABLED"
)

const (
	tolerableValidityRatio = 100               // buffer for certificate revocation update time
	maxClockSkew           = 900 * time.Second // buffer for clock skew
)

type ocspStatusCode int

type ocspStatus struct {
	code ocspStatusCode
	err  error
}

const (
	ocspSuccess                ocspStatusCode = 0
	ocspStatusGood             ocspStatusCode = -1
	ocspStatusRevoked          ocspStatusCode = -2
	ocspStatusUnknown          ocspStatusCode = -3
	ocspStatusOthers           ocspStatusCode = -4
	ocspNoServer               ocspStatusCode = -5
	ocspFailedParseOCSPHost    ocspStatusCode = -6
	ocspFailedComposeRequest   ocspStatusCode = -7
	ocspFailedDecomposeRequest ocspStatusCode = -8
	ocspFailedSubmit           ocspStatusCode = -9
	ocspFailedResponse         ocspStatusCode = -10
	ocspFailedExtractResponse  ocspStatusCode = -11
	ocspFailedParseResponse    ocspStatusCode = -12
	ocspInvalidValidity        ocspStatusCode = -13
	ocspMissedCache            ocspStatusCode = -14
	ocspCacheExpired           ocspStatusCode = -15
	ocspFailedDecodeResponse   ocspStatusCode = -16
)

// copied from crypto/ocsp.go
type certID struct {
	HashAlgorithm pkix.AlgorithmIdentifier
	NameHash      []byte
	IssuerKeyHash []byte
	SerialNumber  *big.Int
}

// cache key
type certIDKey struct {
	HashAlgorithm crypto.Hash
	NameHash      string
	IssuerKeyHash string
	SerialNumber  string
}

var (
	ocspResponseCache     map[certIDKey][]interface{}
	ocspResponseCacheLock *sync.RWMutex
)

// copied from crypto/ocsp
var hashOIDs = map[crypto.Hash]asn1.ObjectIdentifier{
	crypto.SHA1:   asn1.ObjectIdentifier([]int{1, 3, 14, 3, 2, 26}),
	crypto.SHA256: asn1.ObjectIdentifier([]int{2, 16, 840, 1, 101, 3, 4, 2, 1}),
	crypto.SHA384: asn1.ObjectIdentifier([]int{2, 16, 840, 1, 101, 3, 4, 2, 2}),
	crypto.SHA512: asn1.ObjectIdentifier([]int{2, 16, 840, 1, 101, 3, 4, 2, 3}),
}

// copied from crypto/ocsp
func getOIDFromHashAlgorithm(target crypto.Hash) asn1.ObjectIdentifier {
	for hash, oid := range hashOIDs {
		if hash == target {
			return oid
		}
	}
	glog.V(0).Infof("no valid OID is found for the hash algorithm. %#v", target)
	return nil
}

func getHashAlgorithmFromOID(target pkix.AlgorithmIdentifier) crypto.Hash {
	for hash, oid := range hashOIDs {
		if oid.Equal(target.Algorithm) {
			return hash
		}
	}
	glog.V(0).Infof("no valid hash algorithm is found for the oid. Falling back to SHA1: %#v", target)
	return crypto.SHA1
}

// calcTolerableValidity returns the maximum validity buffer
func calcTolerableValidity(thisUpdate, nextUpdate time.Time) time.Duration {
	return durationMax(time.Duration(nextUpdate.Sub(thisUpdate)/tolerableValidityRatio), maxClockSkew)
}

// isInValidityRange checks the validity
func isInValidityRange(currTime, thisUpdate, nextUpdate time.Time) bool {
	if currTime.Sub(thisUpdate.Add(-maxClockSkew)) < 0 {
		return false
	}
	if nextUpdate.Add(calcTolerableValidity(thisUpdate, nextUpdate)).Sub(currTime) < 0 {
		return false
	}
	return true
}

func checkTotalTimeout(totalTimeout *time.Duration, sleepTime time.Duration) (ok bool) {
	if *totalTimeout > 0 {
		*totalTimeout -= sleepTime
	}
	if *totalTimeout <= 0 {
		return false
	}
	glog.V(2).Infof("sleeping %v for retryOCSP. to timeout: %v. retrying", sleepTime, *totalTimeout)
	time.Sleep(sleepTime)
	return true
}

func extractCertIDKeyFromRequest(ocspReq []byte) (*certIDKey, *ocspStatus) {
	r, err := ocsp.ParseRequest(ocspReq)
	if err != nil {
		return nil, &ocspStatus{
			code: ocspFailedDecomposeRequest,
			err:  err,
		}
	}

	// encode CertID, used as a key in the cache
	encodedCertID := &certIDKey{
		r.HashAlgorithm,
		base64.StdEncoding.EncodeToString(r.IssuerNameHash),
		base64.StdEncoding.EncodeToString(r.IssuerKeyHash),
		r.SerialNumber.String(),
	}
	return encodedCertID, &ocspStatus{
		code: ocspSuccess,
	}
}

func encodeCertIDKey(certIDKeyBase64 string) *certIDKey {
	r, err := base64.StdEncoding.DecodeString(certIDKeyBase64)
	if err != nil {
		return nil
	}
	var c certID
	rest, err := asn1.Unmarshal(r, &c)
	if err != nil {
		// error in parsing
		return nil
	}
	if len(rest) > 0 {
		// extra bytes to the end
		return nil
	}
	return &certIDKey{
		getHashAlgorithmFromOID(c.HashAlgorithm),
		base64.StdEncoding.EncodeToString(c.NameHash),
		base64.StdEncoding.EncodeToString(c.IssuerKeyHash),
		c.SerialNumber.String(),
	}
}

func decodeCertIDKey(k *certIDKey) string {
	serialNumber := new(big.Int)
	serialNumber.SetString(k.SerialNumber, 10)
	nameHash, err := base64.StdEncoding.DecodeString(k.NameHash)
	if err != nil {
		return ""
	}
	issuerKeyHash, err := base64.StdEncoding.DecodeString(k.IssuerKeyHash)
	if err != nil {
		return ""
	}
	encodedCertID, err := asn1.Marshal(certID{
		pkix.AlgorithmIdentifier{
			Algorithm:  getOIDFromHashAlgorithm(k.HashAlgorithm),
			Parameters: asn1.RawValue{Tag: 5 /* ASN.1 NULL */},
		},
		nameHash,
		issuerKeyHash,
		serialNumber,
	})
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(encodedCertID)
}

func checkOCSPResponseCache(encodedCertID *certIDKey, subject, issuer *x509.Certificate) *ocspStatus {
	ocspResponseCacheLock.RLock()
	gotValueFromCache := ocspResponseCache[*encodedCertID]
	ocspResponseCacheLock.RUnlock()

	status := extractOCSPCacheResponseValue(gotValueFromCache, subject, issuer)
	if !isValidOCSPStatus(status.code) {
		deleteOCSPCache(encodedCertID)
	}
	return status
}

func deleteOCSPCache(encodedCertID *certIDKey) {
	ocspResponseCacheLock.Lock()
	delete(ocspResponseCache, *encodedCertID)
	cacheUpdated = true
	ocspResponseCacheLock.Unlock()
}

func validateOCSP(ocspRes *ocsp.Response) *ocspStatus {
	curTime := time.Now()

	if ocspRes == nil {
		return &ocspStatus{
			code: ocspFailedDecomposeRequest,
			err:  errors.New("OCSP Response is nil"),
		}
	}
	if !isInValidityRange(curTime, ocspRes.ThisUpdate, ocspRes.NextUpdate) {
		return &ocspStatus{
			code: ocspInvalidValidity,
			err:  fmt.Errorf("invalid validity: producedAt: %v, thisUpdate: %v, nextUpdate: %v", ocspRes.ProducedAt, ocspRes.ThisUpdate, ocspRes.NextUpdate),
		}
	}
	switch ocspRes.Status {
	case ocsp.Good:
		return &ocspStatus{
			code: ocspStatusGood,
			err:  nil,
		}
	case ocsp.Revoked:
		return &ocspStatus{
			code: ocspStatusRevoked,
			err:  fmt.Errorf("OCSP revoked: reason:%v, at:%v", ocspRes.RevocationReason, ocspRes.RevokedAt),
		}
	case ocsp.Unknown:
		return &ocspStatus{
			code: ocspStatusUnknown,
			err:  fmt.Errorf("OCSP unknown"),
		}
	default:
		return &ocspStatus{
			code: ocspStatusOthers,
			err:  fmt.Errorf("OCSP others. %v", ocspRes.Status),
		}
	}
}

func retryOCSPCacheServer(
	client clientInterface,
	req requestFunc,
	ocspServerHost string,
	totalTimeout time.Duration,
	httpTimeout time.Duration) (
	cacheContent *map[string][]interface{},
	ocspS *ocspStatus) {
	var respd map[string][]interface{}
	retryCounter := 0
	sleepTime := time.Duration(0)
	headers := make(map[string]string)
	for {
		sleepTime = defaultWaitAlgo.decorr(retryCounter, sleepTime)
		res, err := retryHTTP(context.TODO(), client, req, "GET", ocspServerHost, headers, nil, httpTimeout, false)
		if err != nil {
			if ok := checkTotalTimeout(&totalTimeout, sleepTime); ok {
				retryCounter++
				continue
			}
			glog.V(2).Infof("failed to get OCSP cache from OCSP Cache Server. %v\n", err)
			return nil, &ocspStatus{
				code: ocspFailedSubmit,
				err:  err,
			}
		}
		defer res.Body.Close()
		glog.V(2).Infof("StatusCode from OCSP Cache Server: %v\n", res.StatusCode)
		if res.StatusCode != http.StatusOK {
			if ok := checkTotalTimeout(&totalTimeout, sleepTime); ok {
				retryCounter++
				continue
			}
			return nil, &ocspStatus{
				code: ocspFailedResponse,
				err:  fmt.Errorf("HTTP code is not OK. %v: %v", res.StatusCode, res.Status),
			}
		}
		glog.V(2).Info("reading contents")

		dec := json.NewDecoder(res.Body)
		for {
			if err := dec.Decode(&respd); err == io.EOF {
				break
			} else if err != nil {
				glog.V(2).Infof("failed to decode OCSP cache. %v\n", err)
				return nil, &ocspStatus{
					code: ocspFailedExtractResponse,
					err:  err,
				}
			}
		}
		break
	}
	return &respd, &ocspStatus{
		code: ocspSuccess,
	}
}

// retryOCSP is the second level of retry method if the returned contents are corrupted. It often happens with OCSP
// serer and retry helps.
func retryOCSP(
	client clientInterface,
	req requestFunc,
	ocspHost string,
	headers map[string]string,
	reqBody []byte,
	issuer *x509.Certificate,
	totalTimeout time.Duration,
	httpTimeout time.Duration) (
	ocspRes *ocsp.Response,
	ocspResBytes []byte,
	ocspS *ocspStatus) {
	retryCounter := 0
	sleepTime := time.Duration(0)
	for {
		sleepTime = defaultWaitAlgo.decorr(retryCounter, sleepTime)
		res, err := retryHTTP(context.TODO(), client, req, "POST", ocspHost, headers, reqBody, httpTimeout, false)
		if err != nil {
			if ok := checkTotalTimeout(&totalTimeout, sleepTime); ok {
				retryCounter++
				continue
			}
			return ocspRes, ocspResBytes, &ocspStatus{
				code: ocspFailedSubmit,
				err:  err,
			}
		}
		defer res.Body.Close()
		glog.V(2).Infof("StatusCode from OCSP Server: %v\n", res.StatusCode)
		if res.StatusCode != http.StatusOK {
			if ok := checkTotalTimeout(&totalTimeout, sleepTime); ok {
				retryCounter++
				continue
			}
			return ocspRes, ocspResBytes, &ocspStatus{
				code: ocspFailedResponse,
				err:  fmt.Errorf("HTTP code is not OK. %v: %v", res.StatusCode, res.Status),
			}
		}
		glog.V(2).Info("reading contents")
		ocspResBytes, err = ioutil.ReadAll(res.Body)
		if err != nil {
			if ok := checkTotalTimeout(&totalTimeout, sleepTime); ok {
				retryCounter++
				continue
			}
			return ocspRes, ocspResBytes, &ocspStatus{
				code: ocspFailedExtractResponse,
				err:  err,
			}
		}
		glog.V(2).Info("parsing OCSP response")
		ocspRes, err = ocsp.ParseResponse(ocspResBytes, issuer)
		if err != nil {
			if ok := checkTotalTimeout(&totalTimeout, sleepTime); ok {
				retryCounter++
				continue
			}
			return ocspRes, ocspResBytes, &ocspStatus{
				code: ocspFailedParseResponse,
				err:  err,
			}
		}
		break
	}
	return ocspRes, ocspResBytes, &ocspStatus{
		code: ocspSuccess,
	}
}

// getRevocationStatus checks the certificate revocation status for subject using issuer certificate.
func getRevocationStatus(subject, issuer *x509.Certificate) *ocspStatus {
	glog.V(2).Infof("Subject: %v, Issuer: %v\n", subject.Subject, issuer.Subject)

	status, ocspReq, encodedCertID := validateWithCache(subject, issuer)
	if isValidOCSPStatus(status) {
		return &ocspStatus{
			code: status,
			err:  nil,
		}
	}
	if ocspReq == nil || encodedCertID == nil {
		return &ocspStatus{
			code: status,
			err:  fmt.Errorf("failed to compose OCSP request.%v", ""),
		}
	}
	glog.V(2).Infof("cache missed\n")
	glog.V(2).Infof("OCSP Server: %v\n", subject.OCSPServer)
	if len(subject.OCSPServer) == 0 {
		return &ocspStatus{
			code: ocspNoServer,
			err:  fmt.Errorf("no OCSP server is attached to the certificate. %v", subject.Subject),
		}
	}
	ocspHost := subject.OCSPServer[0]
	u, err := url.Parse(ocspHost)
	if err != nil {
		return &ocspStatus{
			code: ocspFailedParseOCSPHost,
			err:  fmt.Errorf("failed to parse OCSP server host. %v", ocspHost),
		}
	}
	ocspClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: snowflakeInsecureTransport,
	}
	headers := make(map[string]string)
	headers["Content-Type"] = "application/ocsp-request"
	headers["Accept"] = "application/ocsp-response"
	headers["Content-Length"] = string(len(ocspReq))
	headers["Host"] = u.Hostname()
	ocspRes, ocspResBytes, ocspS := retryOCSP(ocspClient, http.NewRequest, ocspHost, headers, ocspReq, issuer, retryOCSPTimeout, retryOCSPHTTPTimeout)
	if ocspS.code != ocspSuccess {
		return ocspS
	}

	ret := validateOCSP(ocspRes)
	if !isValidOCSPStatus(ret.code) {
		return ret // return invalid
	}
	v := []interface{}{float64(time.Now().UTC().Unix()), base64.StdEncoding.EncodeToString(ocspResBytes)}
	ocspResponseCacheLock.Lock()
	ocspResponseCache[*encodedCertID] = v
	cacheUpdated = true
	ocspResponseCacheLock.Unlock()
	return ret
}

func isValidOCSPStatus(status ocspStatusCode) bool {
	return status == ocspStatusGood || status == ocspStatusRevoked || status == ocspStatusUnknown
}

// verifyPeerCertificate verifies all of certificate revocation status
func verifyPeerCertificate(verifiedChains [][]*x509.Certificate) (err error) {
	for i := 0; i < len(verifiedChains); i++ {
		n := len(verifiedChains[i]) - 1
		if !verifiedChains[i][n].IsCA || string(verifiedChains[i][n].RawIssuer) != string(verifiedChains[i][n].RawSubject) {
			// if the last certificate is not root CA, add it to the list
			rca := caRoot[string(verifiedChains[i][n].RawIssuer)]
			if rca == nil {
				return fmt.Errorf("failed to find root CA. pkix.name: %v", verifiedChains[i][n].Issuer)
			}
			verifiedChains[i] = append(verifiedChains[i], rca)
			n++
		}
		results := getAllRevocationStatus(verifiedChains[i])
		if r := canEarlyExitForOCSP(results, len(verifiedChains[i])); r != nil {
			return r.err
		}
	}

	ocspResponseCacheLock.Lock()
	if cacheUpdated {
		writeOCSPCacheFile()
	}
	cacheUpdated = false
	ocspResponseCacheLock.Unlock()
	return nil
}

func canEarlyExitForOCSP(results []*ocspStatus, chainSize int) *ocspStatus {
	msg := ""
	if ocspFailOpen == OCSPFailOpenFalse {
		// Fail closed. any error is returned to stop connection
		for _, r := range results {
			if r.err != nil {
				return r
			}
		}
	} else {
		// Fail open and all results are valid.
		allValid := len(results) == chainSize
		for _, r := range results {
			if !isValidOCSPStatus(r.code) {
				allValid = false
				break
			}
		}
		for _, r := range results {
			if allValid && r.code == ocspStatusRevoked {
				return r
			}
			if r.code != ocspStatusGood && r.err != nil {
				msg += "\n" + r.err.Error()
			}
		}
	}
	if len(msg) > 0 {
		glog.V(1).Infof(
			"WARNING!!! Using fail-open to connect. Driver is connecting to an "+
				"HTTPS endpoint without OCSP based Certificate Revocation checking "+
				"as it could not obtain a valid OCSP Response to use from the CA OCSP "+
				"responder. Detail: %v", msg[1:])
	}
	return nil
}

func validateWithCacheForAllCertificates(verifiedChains []*x509.Certificate) bool {
	n := len(verifiedChains) - 1
	for j := 0; j < n; j++ {
		subject := verifiedChains[j]
		issuer := verifiedChains[j+1]
		status, _, _ := validateWithCache(subject, issuer)
		if !isValidOCSPStatus(status) {
			return false
		}
	}
	return true
}

func validateWithCache(subject, issuer *x509.Certificate) (ocspStatusCode, []byte, *certIDKey) {
	ocspReq, err := ocsp.CreateRequest(subject, issuer, &ocsp.RequestOptions{})
	if err != nil {
		glog.V(2).Infof("failed to create OCSP request from the certificates.\n")
		return ocspFailedComposeRequest, nil, nil
	}
	encodedCertID, ocspS := extractCertIDKeyFromRequest(ocspReq)
	if ocspS.code != ocspSuccess {
		glog.V(2).Infof("failed to extract CertID from OCSP Request.\n")
		return ocspFailedComposeRequest, ocspReq, nil
	}
	ocspValidatedWithCache := checkOCSPResponseCache(encodedCertID, subject, issuer)
	return ocspValidatedWithCache.code, ocspReq, encodedCertID
}

func downloadOCSPCacheServer() {
	if strings.EqualFold(os.Getenv(cacheServerEnabledEnv), "false") {
		glog.V(2).Infof("skipping downloading OCSP Cache.")
		return
	}
	ocspClient := &http.Client{
		Timeout:   30 * time.Second,
		Transport: snowflakeInsecureTransport,
	}
	ocspURL := fmt.Sprintf("%v/%v", cacheServerURL, cacheFileBaseName)
	glog.V(2).Infof("downloading OCSP Cache from server %v", ocspURL)
	ret, ocspStatus := retryOCSPCacheServer(ocspClient, http.NewRequest, ocspURL, retryOCSPTimeout, retryOCSPHTTPTimeout)
	if ocspStatus.code != ocspSuccess {
		return
	}

	ocspResponseCacheLock.Lock()
	for k, cacheValue := range *ret {
		status := extractOCSPCacheResponseValueWithoutSubject(cacheValue)
		if !isValidOCSPStatus(status.code) {
			continue
		}
		cacheKey := encodeCertIDKey(k)
		ocspResponseCache[*cacheKey] = cacheValue
	}
	cacheUpdated = true
	ocspResponseCacheLock.Unlock()
}

func getAllRevocationStatus(verifiedChains []*x509.Certificate) []*ocspStatus {
	cached := validateWithCacheForAllCertificates(verifiedChains)
	if !cached {
		downloadOCSPCacheServer()
	}
	n := len(verifiedChains) - 1
	results := make([]*ocspStatus, n)
	for j := 0; j < n; j++ {
		results[j] = getRevocationStatus(verifiedChains[j], verifiedChains[j+1])
		if !isValidOCSPStatus(results[j].code) {
			return results
		}
	}
	return results
}

// verifyPeerCertificateSerial verifies the certificate revocation status in serial.
func verifyPeerCertificateSerial(_ [][]byte, verifiedChains [][]*x509.Certificate) (err error) {
	return verifyPeerCertificate(verifiedChains)
}

// initOCSPCache initializes OCSP Response cache file.
func initOCSPCache() {
	ocspResponseCache = make(map[certIDKey][]interface{})
	ocspResponseCacheLock = &sync.RWMutex{}
	cacheFileName = filepath.Join(cacheDir, cacheFileBaseName)

	glog.V(2).Infof("reading OCSP Response cache file. %v\n", cacheFileName)
	f, err := os.Open(cacheFileName)
	if err != nil {
		glog.Infof("failed to open. Ignored. %v\n", err)
		return
	}
	defer f.Close()

	buf := make(map[string][]interface{})

	r := bufio.NewReader(f)
	dec := json.NewDecoder(r)
	for {
		if err := dec.Decode(&buf); err == io.EOF {
			break
		} else if err != nil {
			glog.V(2).Infof("failed to read. Ignored. %v\n", err)
			return
		}
	}
	for k, cacheValue := range buf {
		status := extractOCSPCacheResponseValueWithoutSubject(cacheValue)
		if !isValidOCSPStatus(status.code) {
			continue
		}
		cacheKey := encodeCertIDKey(k)
		ocspResponseCache[*cacheKey] = cacheValue

	}
	cacheUpdated = false
}
func extractOCSPCacheResponseValueWithoutSubject(cacheValue []interface{}) *ocspStatus {
	return extractOCSPCacheResponseValue(cacheValue, nil, nil)
}

func extractOCSPCacheResponseValue(cacheValue []interface{}, subject, issuer *x509.Certificate) *ocspStatus {
	subjectName := "Unknown"
	if subject != nil {
		subjectName = subject.Subject.CommonName
	}

	curTime := time.Now()
	if len(cacheValue) != 2 {
		return &ocspStatus{
			code: ocspMissedCache,
			err:  fmt.Errorf("miss cache data. subject: %v", subjectName),
		}
	}
	if ts, ok := cacheValue[0].(float64); ok {
		currentTime := float64(curTime.UTC().Unix())
		if currentTime-ts >= cacheExpire {
			return &ocspStatus{
				code: ocspCacheExpired,
				err: fmt.Errorf("cache expired. current: %v, cache: %v",
					time.Unix(int64(currentTime), 0).UTC(), time.Unix(int64(ts), 0).UTC()),
			}
		}
	} else {
		return &ocspStatus{
			code: ocspFailedDecodeResponse,
			err:  errors.New("the first cache element is not float64"),
		}
	}
	var err error
	var r *ocsp.Response
	if s, ok := cacheValue[1].(string); ok {
		var b []byte
		b, err = base64.StdEncoding.DecodeString(s)
		if err != nil {
			return &ocspStatus{
				code: ocspFailedDecodeResponse,
				err:  fmt.Errorf("failed to decode OCSP Response value in a cache. subject: %v, err: %v", subjectName, err),
			}
		}
		// check the revocation status here
		r, err = ocsp.ParseResponse(b, issuer)
		if err != nil {
			glog.V(2).Infof("the second cache element is not a valid OCSP Response. Ignored. subject: %v\n", subjectName)
			return &ocspStatus{
				code: ocspFailedParseResponse,
				err:  fmt.Errorf("failed to parse OCSP Respose. subject: %v, err: %v", subjectName, err),
			}
		}
	} else {
		return &ocspStatus{
			code: ocspFailedDecodeResponse,
			err:  errors.New("the second cache element is not string"),
		}

	}
	return validateOCSP(r)
}

// writeOCSPCacheFile writes a OCSP Response cache file. This is called if all revocation status is success.
// lock file is used to mitigate race condition with other process.
func writeOCSPCacheFile() {
	glog.V(2).Infof("writing OCSP Response cache file. %v\n", cacheFileName)
	cacheLockFileName := cacheFileName + ".lck"
	err := os.Mkdir(cacheLockFileName, 0600)
	switch {
	case os.IsExist(err):
		statinfo, err := os.Stat(cacheLockFileName)
		if err != nil {
			glog.V(2).Infof("failed to write OCSP response cache file. file: %v, err: %v. ignored.\n", cacheFileName, err)
			return
		}
		if time.Since(statinfo.ModTime()) < 15*time.Minute {
			glog.V(2).Infof("other process locks the cache file. %v. ignored.\n", cacheFileName)
			return
		}
		err = os.Remove(cacheLockFileName)
		if err != nil {
			glog.V(2).Infof("failed to delete lock file. file: %v, err: %v. ignored.\n", cacheLockFileName, err)
			return
		}
		err = os.Mkdir(cacheLockFileName, 0600)
		if err != nil {
			glog.V(2).Infof("failed to delete lock file. file: %v, err: %v. ignored.\n", cacheLockFileName, err)
			return
		}
	}
	defer os.RemoveAll(cacheLockFileName)

	buf := make(map[string][]interface{})
	for k, v := range ocspResponseCache {
		cacheKeyInBase64 := decodeCertIDKey(&k)
		buf[cacheKeyInBase64] = v
	}

	j, err := json.Marshal(buf)
	if err != nil {
		glog.V(2).Info("failed to convert OCSP Response cache to JSON. ignored.")
		return
	}
	err = ioutil.WriteFile(cacheFileName, j, 0644)
	if err != nil {
		glog.V(2).Infof("failed to write OCSP Response cache. err: %v. ignored.\n", err)
	}
}

// readCACerts read a set of root CAs
func readCACerts() {
	raw := []byte(caRootPEM)
	certPool = x509.NewCertPool()
	caRoot = make(map[string]*x509.Certificate)
	var p *pem.Block
	for {
		p, raw = pem.Decode(raw)
		if p == nil {
			break
		}
		if p.Type != "CERTIFICATE" {
			continue
		}
		c, err := x509.ParseCertificate(p.Bytes)
		if err != nil {
			panic("failed to parse CA certificate.")
		}
		certPool.AddCert(c)
		caRoot[string(c.RawSubject)] = c
	}
}

// createOCSPCacheDir creates OCSP response cache directory. If SNOWFLAKE_TEST_WORKSPACE is set,
func createOCSPCacheDir() {
	cacheDir = os.Getenv("SNOWFLAKE_TEST_WORKSPACE")
	if cacheDir == "" {
		switch runtime.GOOS {
		case "windows":
			cacheDir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Snowflake", "Caches")
		case "darwin":
			home := os.Getenv("HOME")
			if home == "" {
				glog.V(2).Info("HOME is blank.")
			}
			cacheDir = filepath.Join(home, "Library", "Caches", "Snowflake")
		default:
			home := os.Getenv("HOME")
			if home == "" {
				glog.V(2).Info("HOME is blank")
			}
			cacheDir = filepath.Join(home, ".cache", "snowflake")
		}
	}
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err := os.MkdirAll(cacheDir, os.ModePerm)
		if err != nil {
			glog.V(2).Infof("failed to create cache directory. %v, err: %v. ignored\n", cacheDir, err)
		}
	}
}

func init() {
	readCACerts()
	createOCSPCacheDir()
	initOCSPCache()
}

// snowflakeInsecureTransport is the transport object that doesn't do certificate revocation check.
var snowflakeInsecureTransport = &http.Transport{
	MaxIdleConns:    10,
	IdleConnTimeout: 30 * time.Minute,
	Proxy:           http.ProxyFromEnvironment,
}

// SnowflakeTransport includes the certificate revocation check with OCSP in sequential. By default, the driver uses
// this transport object.
var SnowflakeTransport = &http.Transport{
	TLSClientConfig: &tls.Config{
		RootCAs:               certPool,
		VerifyPeerCertificate: verifyPeerCertificateSerial,
	},
	MaxIdleConns:    10,
	IdleConnTimeout: 30 * time.Minute,
	Proxy:           http.ProxyFromEnvironment,
}

// SnowflakeTransportTest includes the certificate revocation check in parallel
var SnowflakeTransportTest = SnowflakeTransport
