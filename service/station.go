package service

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"net/http"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const radioBrowserURL = "https://de1.api.radio-browser.info"

func getUserAgent() string {
	version := application.BuildInfo.Main.Version
	if version == "" {
		version = "0.0.0-dev"
	}
	return fmt.Sprintf(
		"Vlna/%s (https://github.com/SladkyCitron/vlna) Go/%s (%s/%s)",
		version,
		runtime.Version()[2:],
		runtime.GOOS,
		runtime.GOARCH,
	)
}

type Station struct {
	ChangeUUID                string  `json:"changeuuid"`
	StationUUID               string  `json:"stationuuid"`
	Name                      string  `json:"name"`
	URL                       string  `json:"url"`
	URLResolved               string  `json:"url_resolved"`
	Homepage                  string  `json:"homepage"`
	Favicon                   string  `json:"favicon"`
	Tags                      string  `json:"tags"`
	Country                   string  `json:"country"`
	CountryCode               string  `json:"countrycode"`
	State                     string  `json:"state"`
	ISO_3166_2                string  `json:"iso_3166_2"`
	Language                  string  `json:"language"`
	LanguageCodes             string  `json:"languagecodes"`
	Votes                     int     `json:"votes"`
	LastChangeTime            string  `json:"lastchangetime"`
	LastChangeTimeISO8601     string  `json:"lastchangetime_iso8601"`
	Codec                     string  `json:"codec"`
	Bitrate                   int     `json:"bitrate"`
	HLS                       int     `json:"hls"`
	LastCheckOK               int     `json:"lastcheckok"`
	LastCheckTime             string  `json:"lastchecktime"`
	LastCheckTimeISO8601      string  `json:"lastchecktime_iso8601"`
	LastCheckOKTime           string  `json:"lastcheckoktime"`
	LastCheckOKTimeISO8601    string  `json:"lastcheckoktime_iso8601"`
	LastLocalCheckTime        string  `json:"lastlocalchecktime"`
	LastLocalCheckTimeISO8601 string  `json:"lastlocalchecktime_iso8601"`
	ClickTimestamp            string  `json:"clicktimestamp"`
	ClickTimestampISO8601     string  `json:"clicktimestamp_iso8601"`
	ClickCount                int     `json:"clickcount"`
	ClickTrend                int     `json:"clicktrend"`
	SSLError                  int     `json:"ssl_error"`
	GeoLat                    float64 `json:"geo_lat"`
	GeoLong                   float64 `json:"geo_long"`
	GeoDistance               float64 `json:"geo_distance"`
	HasExtendedInfo           bool    `json:"has_extended_info"`
}

type Stations []Station

type StationService struct {
	ctx    context.Context
	client *http.Client
}

func NewStationService() *StationService {
	return &StationService{
		client: &http.Client{},
	}
}

func (s *StationService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	s.ctx = ctx
	return nil
}

// GetStationsByCountryCode fetches a list of stations filtered by the given country code,
// limited to a maximum of 100 stations.
//
// It is used for the Explore page to display stations from the user's region.
func (s *StationService) GetStationsByCountryCode(countryCode string) (Stations, error) {
	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodGet,
		radioBrowserURL+"/json/stations/bycountrycodeexact/"+countryCode+"?order=votes&limit=100",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create network request: %w", err)
	}
	req.Header.Set("User-Agent", getUserAgent())

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var stations Stations
	if err := json.UnmarshalRead(resp.Body, &stations); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return stations, nil
}
