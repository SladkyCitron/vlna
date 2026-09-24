package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type IPInfo struct {
	IPVersion       int      `json:"ipVersion"`
	IPAddress       string   `json:"ipAddress"`
	Latitude        float64  `json:"latitude"`
	Longitude       float64  `json:"longitude"`
	CountryName     string   `json:"countryName"`
	CountryCode     string   `json:"countryCode"`
	Capital         string   `json:"capital"`
	PhoneCodes      []int    `json:"phoneCodes"`
	TimeZones       []string `json:"timeZones"`
	ZipCode         string   `json:"zipCode"`
	CityName        string   `json:"cityName"`
	RegionName      string   `json:"regionName"`
	RegionCode      string   `json:"regionCode"`
	Continent       string   `json:"continent"`
	ContinentCode   string   `json:"continentCode"`
	Currencies      []string `json:"currencies"`
	Languages       []string `json:"languages"`
	ASN             string   `json:"asn"`
	ASNOrganization string   `json:"asnOrganization"`
	IsProxy         bool     `json:"isProxy"`
}

type IPInfoService struct {
	ctx    context.Context
	client *http.Client
}

func NewIPInfoService() *IPInfoService {
	return &IPInfoService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *IPInfoService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	s.ctx = ctx
	return nil
}

func (s *IPInfoService) Fetch() (*IPInfo, error) {
	req, err := http.NewRequestWithContext(s.ctx, http.MethodGet, "https://free.freeipapi.com/api/v1/json", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create network request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var ipInfo IPInfo
	if err := json.NewDecoder(resp.Body).Decode(&ipInfo); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &ipInfo, nil
}
