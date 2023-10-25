package quotaclient

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type DNSResponse struct {
	QuotaLimits struct {
		Records        int64 `json:"records"`
		SecondaryZones int64 `json:"secondaryZones"`
		Zones          int64 `json:"zones"`
	} `json:"quotaLimits"`
	QuotaUsage struct {
		Records        int64 `json:"records"`
		SecondaryZones int64 `json:"secondaryZones"`
		Zones          int64 `json:"zones"`
	} `json:"quotaUsage"`
}

func DNSQuotas(token string) (*DNSResponse, error) {
	request2, err := http.NewRequest(http.MethodGet, "https://dns.de-fra.ionos.com/quota", nil)
	if err != nil {
		return nil, err
	}

	request2.Header.Add("Authorization", token)

	response2, err := http.DefaultClient.Do(request2)
	if err != nil {
		return nil, err
	}

	if response2.StatusCode != 200 {
		return nil, errors.New("Error while getting DNS quotas")
	}

	var responseObject2 DNSResponse
	body2, err := io.ReadAll(response2.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body2, &responseObject2)
	if err != nil {
		return nil, err
	}

	return &responseObject2, nil
}
