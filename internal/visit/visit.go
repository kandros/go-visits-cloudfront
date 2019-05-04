package visit

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/kandros/visits/internal/dynamo"
	"github.com/kandros/visits/internal/util"
)

type Visit struct {
	VisitID   string    `json:"visit_id"`
	URL       string    `json:"url"`
	Timestamp time.Time `json:"timestamp"`
	IP        string    `json:"ip"`
	CountryID string    `json:"country_id"`
	IsMobile  bool      `json:"is_mobile"`
	IsTablet  bool      `json:"is_tablet"`
}

func NewFromRequest(r *http.Request) Visit {
	ip, _ := util.MustGetIp(r.RemoteAddr)

	visitID := uuid.New().String()
	timestamp := time.Now()
	countryID, isMobile, isTablet := getCloudFrontHeaders(r.Header)

	return Visit{
		VisitID:   visitID,
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
