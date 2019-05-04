package visit

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/kandros/visits/internal/dynamo"
	"github.com/kandros/visits/internal/util"
)

type RequestBodyInput struct {
	Href     string `json:"href"`
	DeviceID string `json:"device_id"`
	VisitID  string `json:"visit_id"`
	Port     string `json:"port"`
	Hostname string `json:"hostname"`
}

type Visit struct {
	Href      string    `json:"href"`
	DeviceID  string    `json:"device_id"`
	VisitID   string    `json:"visit_id"`
	UserAgent string    `json:"user_agent"`
	Hostname  string    `json:"hostname"`
	URL       string    `json:"url"`
	Timestamp time.Time `json:"timestamp"`
	IP        string    `json:"ip"`
	CountryID string    `json:"country_id"`
	IsMobile  bool      `json:"is_mobile"`
	IsTablet  bool      `json:"is_tablet"`
}

func NewFromRequest(r *http.Request) Visit {
	ip, _ := util.MustGetIp(r.RemoteAddr)
	decoder := json.NewDecoder(r.Body)
	var b RequestBodyInput
	err := decoder.Decode(&b)
	if err != nil {
		panic(err)
	}

	timestamp := time.Now()
	userAgent := getUserAgent(r.Header)
	countryID, isMobile, isTablet := getCloudFrontHeaders(r.Header)

	return Visit{
		UserAgent: userAgent,
		Href:      b.Href,
		Hostname:  b.Hostname,
		VisitID:   b.VisitID,
		DeviceID:  b.DeviceID,
		URL:       "https://google.com",
		Timestamp: timestamp,
		IP:        ip,
		CountryID: countryID,
		IsMobile:  isMobile,
		IsTablet:  isTablet,
	}
}

func (v Visit) Persist() error {
	_, err := dynamo.Store(v)
	if err != nil {
		return err
	}

	return nil
}

func getUserAgent(headers http.Header) (useragent string) {
	if v, ok := headers["User-Agent"]; ok && len(v) > 0 {
		useragent = v[0]
	}
	return
}

func getCloudFrontHeaders(headers http.Header) (countryID string, isMobile, isTablet bool) {
	if v, ok := headers["Cloudfront-Viewer-Country"]; ok && len(v) > 0 {
		countryID = v[0]
	}
	if v, ok := headers["Cloudfront-Is-Mobile-Viewer"]; ok && len(v) > 0 {
		isMobile, _ = strconv.ParseBool(v[0])
	}
	if v, ok := headers["Cloudfront-Is-Tablet-Viewer"]; ok && len(v) > 0 {
		isTablet, _ = strconv.ParseBool(v[0])
	}
	return
}
