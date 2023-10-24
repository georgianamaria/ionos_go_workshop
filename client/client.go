package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type DBaaSResponse struct {
	QuotaUsage struct {
		CountMongoclustersDbaasIonosCom    string `json:"count/mongoclusters.dbaas.ionos.com"`
		CountPostgresclustersDbaasIonosCom string `json:"count/postgresclusters.dbaas.ionos.com"`
		Cpu                                string `json:"cpu"`
		Memory                             string `json:"memory"`
		Storage                            string `json:"storage"`
	} `json:"quotaUsage"`
	QuotaLimits struct {
		CountMongoclustersDbaasIonosCom    string `json:"count/mongoclusters.dbaas.ionos.com"`
		CountPostgresclustersDbaasIonosCom string `json:"count/postgresclusters.dbaas.ionos.com"`
		Cpu                                string `json:"cpu"`
		Memory                             string `json:"memory"`
		Storage                            string `json:"storage"`
	} `json:"quotaLimits"`
}

type DNSResponse struct {
	QuotaLimits struct {
		Records        int `json:"records"`
		SecondaryZones int `json:"secondaryZones"`
		Zones          int `json:"zones"`
	} `json:"quotaLimits"`
	QuotaUsage struct {
		Records        int `json:"records"`
		SecondaryZones int `json:"secondaryZones"`
		Zones          int `json:"zones"`
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

func DBaaSQuotas(token string) (*DBaaSResponse, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		"https://api.ionos.com/databases/quota",
		nil,
	)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", token)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New("Error while getting DBaaS quotas " + fmt.Sprint(response.StatusCode))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var responseObject DBaaSResponse
	_ = json.Unmarshal(body, &responseObject)
	return &responseObject, nil
}
