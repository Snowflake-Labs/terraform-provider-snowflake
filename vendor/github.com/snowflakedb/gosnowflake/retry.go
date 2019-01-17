// Copyright (c) 2017-2018 Snowflake Computing Inc. All right reserved.

package gosnowflake

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"context"

	"sync"
)

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

const requestGUIDKey string = "request_guid"

// Format of "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
const uuidLen int = 36

// This class takes in an url during construction and replace the
// value of request_guid every time the replace() is called
// When the url does not contain request_guid, just return the original
// url
type requestGUIDReplacerI interface {
	// replace the url with new ID
	replace() string
}

// Make requestGUIDReplacer given a url string
func makeRequestGUIDReplacer(url string) requestGUIDReplacerI {
	startIndex := strings.Index(url, requestGUIDKey)
	if startIndex == -1 {
		return &transientReplacer{url}
	}
	replacer := &requestGUIDReplacer{}
	startIndex += len(requestGUIDKey) + 1
	replacer.prefix = url[:startIndex]

	startIndex += uuidLen
	replacer.suffix = url[startIndex:]
	return replacer
}

// this replacer does nothing but replace the url
type transientReplacer struct {
	url string
}

func (replacer *transientReplacer) replace() string {
	return replacer.url
}

/*
requestGUIDReplacer is a one-shot object that is created out of the retry loop and
called with replace to change the retry_guid's value upon every retry
*/
type requestGUIDReplacer struct {
	// cached prefix and suffix to avoid parsing same url again
	prefix string
	suffix string
}

/**
This function would replace they value of the requestGUIDKey in a url with a newly
generated uuid
*/
func (replacer *requestGUIDReplacer) replace() string {
	return replacer.prefix + uuid.New().String() + replacer.suffix
}

type waitAlgo struct {
	mutex *sync.Mutex   // required for random.Int63n
	base  time.Duration // base wait time
	cap   time.Duration // maximum wait time
}

func randSecondDuration(n time.Duration) time.Duration {
	return time.Duration(random.Int63n(int64(n/time.Second))) * time.Second
}

// decorrelated jitter backoff
func (w *waitAlgo) decorr(attempt int, sleep time.Duration) time.Duration {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	t := 3*sleep - w.base
	switch {
	case t > 0:
		return durationMin(w.cap, randSecondDuration(t)+w.base)
	case t < 0:
		return durationMin(w.cap, randSecondDuration(-t)+3*sleep)
	}
	return w.base
}

var defaultWaitAlgo = &waitAlgo{
	mutex: &sync.Mutex{},
	base:  5 * time.Second,
	cap:   160 * time.Second,
}

type requestFunc func(method, urlStr string, body io.Reader) (*http.Request, error)

type clientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

func retryHTTP(
	ctx context.Context,
	client clientInterface,
	req requestFunc,
	method string,
	fullURL string,
	headers map[string]string,
	body []byte,
	timeout time.Duration,
	raise4XX bool) (res *http.Response, err error) {

	totalTimeout := timeout
	glog.V(2).Infof("retryHTTP.totalTimeout: %v", totalTimeout)
	retryCounter := 0
	sleepTime := time.Duration(0)

	var rIDReplacer requestGUIDReplacerI

	for {
		req, err := req(method, fullURL, bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		if req != nil {
			// req can be nil in tests
			req = req.WithContext(ctx)
		}
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		res, err = client.Do(req)
		if err == nil && res.StatusCode == http.StatusOK {
			// exit if success
			break
		}
		if raise4XX && res != nil && res.StatusCode >= 400 && res.StatusCode < 500 {
			// abort connection if raise4XX flag is enabled and the range of HTTP status code are 4XX.
			// This is currently used for Snowflake login. The caller must generate an error object based on HTTP status.
			break
		}

		// context cancel or timeout
		if err != nil {
			urlError, isURLError := err.(*url.Error)
			if isURLError &&
				(urlError.Err == context.DeadlineExceeded || urlError.Err == context.Canceled) {
				return res, urlError.Err
			}
		}

		// cannot just return 4xx and 5xx status as the error can be sporadic. retry often helps.
		if err != nil {
			glog.V(2).Infof(
				"failed http connection. no response is returned. err: %v. retrying...\n", err)
		} else {
			glog.V(2).Infof(
				"failed http connection. HTTP Status: %v. retrying...\n", res.StatusCode)
		}
		// uses decorrelated jitter backoff
		sleepTime = defaultWaitAlgo.decorr(retryCounter, sleepTime)

		if totalTimeout > 0 {
			glog.V(2).Infof("to timeout: %v", totalTimeout)
			// if any timeout is set
			totalTimeout -= sleepTime
			if totalTimeout <= 0 {
				if err != nil {
					return nil, fmt.Errorf("timeout. err: %v. Hanging?", err)
				}
				if res != nil {
					return nil, fmt.Errorf("timeout. HTTP Status: %v. Hanging?", res.StatusCode)
				}
				return nil, errors.New("timeout. Hanging?")
			}
		}
		retryCounter++
		if rIDReplacer == nil {
			rIDReplacer = makeRequestGUIDReplacer(fullURL)
		}
		fullURL = rIDReplacer.replace()
		glog.V(2).Infof("sleeping %v. to timeout: %v. retrying", sleepTime, totalTimeout)

		await := time.NewTimer(sleepTime)
		select {
		case <-await.C:
			// retry the request
		case <-ctx.Done():
			await.Stop()
			return res, ctx.Err()
		}
	}
	return res, err
}
