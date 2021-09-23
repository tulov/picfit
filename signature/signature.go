package signature

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var signRegex = regexp.MustCompile("&?sig=[^&]*")

// VerifyParameters encodes map parameters with a key and returns if parameters match signature
func VerifyParameters(key string, qs map[string]interface{}) bool {
	params := url.Values{}

	for k, v := range qs {
		s, ok := v.(string)
		if ok {
			params.Set(k, s)
			continue
		}

		l, ok := v.([]string)
		if ok {
			for i := range l {
				params.Add(k, l[i])
			}
		}
	}

	return VerifySign(key, params.Encode())
}

func VerifyRequest(key string, method string, requestUrl string, bodyHash string, qs map[string]interface{}) bool {
	params := url.Values{}
	params.Set("_method", method)
	params.Set("_url", requestUrl)
	params.Set("_body_hash", bodyHash)
	for k, v := range qs {
		s, ok := v.(string)
		if ok {
			params.Set(k, s)
			continue
		}

		l, ok := v.([]string)
		if ok {
			for i := range l {
				params.Add(k, l[i])
			}
		}
	}

	return VerifySign(key, params.Encode())
}

// Sign encodes query string using a key
func Sign(key string, qs string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(qs))

	byteArray := mac.Sum(nil)

	return hex.EncodeToString(byteArray)
}

// SignRaw encodes raw query string (not sorted) using a key
func SignRaw(key string, method string, requestUrl string, filePath string, queryString string) (string, error) {
	params := url.Values{}
	params.Add("_method", method)
	params.Add("_url", requestUrl)
	var body io.Reader
	if filePath == "" {
		body = strings.NewReader("")
	} else {
		f, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		body = f
	}
	params.Add("_body_hash", SignBody(key, body))
	queryString = fmt.Sprintf("%s&%s", queryString, params.Encode())
	values, err := url.ParseQuery(queryString)
	if err != nil {
		return "", err
	}

	return Sign(key, values.Encode()), nil
}

// AppendSign appends the signature to query string
func AppendSign(qs string, sig string) string {
	params := url.Values{}
	params.Add("sig", sig)

	return fmt.Sprintf("%s&%s", qs, params.Encode())
}

// VerifySign extracts the signature and compare it with query string
func VerifySign(key string, qs string) bool {
	unsignedQueryString := signRegex.ReplaceAllString(qs, "")

	sign := Sign(key, unsignedQueryString)
	values, _ := url.ParseQuery(qs)

	return values.Get("sig") == sign
}

func SignBody(secretKey string, body io.Reader) string {
	b, err := ioutil.ReadAll(body)
	bodyHash := ""
	if err == nil {
		mac := hmac.New(sha1.New, []byte(secretKey))
		mac.Write(b)
		bodyHash = hex.EncodeToString(mac.Sum(nil))
	}
	return bodyHash
}
