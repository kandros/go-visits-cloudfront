package visit

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
	"time"
)

// req, _ := http.NewRequest(http.MethodGet, "/", nil)

func TestNewFromRequest(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		panic(err)
	}
	request.RemoteAddr = "2.39.43.60:443"
	request.Header = cloudFrontMockHeaders()
	v := NewFromRequest(request)

	stringNotEmpty(t, v.VisitID)
	stringsEqual(t, v.IP, "2.39.43.60")
	stringsEqual(t, v.URL, "https://google.com")
	stringsEqual(t, v.CountryID, "IT")
	boolEqual(t, v.IsMobile, true)
	boolEqual(t, v.IsTablet, false)

	if v.Timestamp.Unix() > time.Now().Unix() {
		t.Errorf("Timestamp cannot be after now, got %v", v.Timestamp)
	}
}

func stringNotEmpty(t *testing.T, got string) {
	if got == "" {
		t.Errorf("got empty string")
	}
}

func stringsEqual(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

func boolEqual(t *testing.T, got, want bool) {
	if got != want {
		t.Errorf("expected %s to be want '%v' got '%v' ", reflect.TypeOf(got).Name(), want, got)
	}
}

func TestGetCloudFrontHeaders(t *testing.T) {
	type expected struct {
		countryID string
		isMobile  bool
		isTablet  bool
	}
	checkHeaders := func(t *testing.T, headers http.Header, e expected) {
		t.Helper()
		countryID, isMobile, isTablet := getCloudFrontHeaders(headers)

		if countryID != e.countryID {
			t.Errorf("expected countryID to be '%s' got '%s'", e.countryID, countryID)
		}

		if isMobile != e.isMobile {
			t.Errorf("expected isMobile to be '%v' got '%v'", e.isMobile, isMobile)
		}

		if isTablet != e.isTablet {
			t.Errorf("expected isTablet to be '%v' got '%v'", e.isTablet, isTablet)
		}
	}
	t.Run("using cloudfront headers", func(t *testing.T) {
		headers := cloudFrontMockHeaders()
		checkHeaders(t, headers, expected{countryID: "IT", isMobile: true, isTablet: false})
	})

	t.Run("using empty headers", func(t *testing.T) {
		var headers = http.Header{}
		checkHeaders(t, headers, expected{countryID: "", isMobile: false, isTablet: false})
	})
}

func cloudFrontMockHeaders() http.Header {
	var headers = http.Header{}
	json.Unmarshal(jsonHeaders, &headers)
	return headers
}

// headers from a read lambda log
var jsonHeaders = []byte(`
{
	"Accept": [
		"*/*"
	],
	"Accept-Encoding": [
		"gzip"
	],
	"Cloudfront-Forwarded-Proto": [
		"https"
	],
	"Cloudfront-Is-Desktop-Viewer": [
		"true"
	],
	"Cloudfront-Is-Mobile-Viewer": [
		"true"
	],
	"Cloudfront-Is-Smarttv-Viewer": [
		"false"
	],
	"Cloudfront-Is-Tablet-Viewer": [
		"false"
	],
	"Cloudfront-Viewer-Country": [
		"IT"
	],
	"Connection": [
		"close"
	],
	"User-Agent": [
		"curl/7.54.0"
	],
	"Via": [
		"2.0 eb5be0dD626eaabd9fb97f4fb78fcb40.cloudfront.net (CloudFront)"
	],
	"X-Amz-Cf-Id": [
		"W0Jk9zVPoRwEIxA2yq7paxFdFCNuBUgW9ZsiXq1TPBNhcyJ9re78sg=="
	],
	"X-Amzn-Trace-Id": [
		"Root=1-6cc01e53-421b553f9a54abec794915b5"
	],
	"X-Context": [
		"{\"apiId\":\"u9i2g6007k\",\"resourceId\":\"o8fd8zemn8\",\"requestId\":\"06d5a812-656b-11e9-abec-f12c9e90f22b\",\"accountId\":\"054195889806\",\"stage\":\"staging\",\"identity\":{\"apiKey\":\"\",\"accountId\":\"\",\"userAgent\":\"curl/7.54.0\",\"sourceIp\":\"2.40.43.60\",\"accessKey\":\"\",\"caller\":\"\",\"user\":\"\",\"userARN\":\"\",\"cognitoIdentityId\":\"\",\"cognitoIdentityPoolId\":\"\",\"cognitoAuthenticationType\":\"\",\"cognitoAuthenticationProvider\":\"\"},\"authorizer\":null}"
	],
	"X-Forwarded-For": [
		"2.39.43.60, 70.132.17.146"
	],
	"X-Forwarded-Port": [
		"443"
	],
	"X-Forwarded-Proto": [
		"https"
	],
	"X-Request-Id": [
		"06d5a912-666b-11e9-abec-f12c9e90f22b"
	],
	"X-Stage": [
		"staging"
	]
}
`)
