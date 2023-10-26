package adapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"workshop_demo/internal/model"
)

type DNSQuotaAdapter struct {
}

type DBaaSQuotaAdapter struct {
}

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

func (_ *DNSQuotaAdapter) FetchDNSQuotas(token string) (model.Quota[model.DNSQuota], error) {
	var quota model.Quota[model.DNSQuota]

	request, err := http.NewRequest(http.MethodGet, "https://dns.de-fra.ionos.com/quota", nil)
	if err != nil {
		return quota, err
	}

	request.Header.Add("Authorization", token)

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return quota, err
	}

	if res.StatusCode != 200 {
		return quota, errors.New("Error while getting DNS quotas")
	}

	var clientObject DNSResponse
	body2, err := io.ReadAll(res.Body)
	if err != nil {
		return quota, err
	}

	err = json.Unmarshal(body2, &clientObject)
	if err != nil {
		return quota, err
	}

	quota = model.Quota[model.DNSQuota]{
		Limit: model.DNSQuota{
			Records:        clientObject.QuotaLimits.Records,
			SecondaryZones: clientObject.QuotaLimits.SecondaryZones,
			Zones:          clientObject.QuotaLimits.Zones,
		},
		Usage: model.DNSQuota{
			Records:        clientObject.QuotaUsage.Records,
			SecondaryZones: clientObject.QuotaUsage.SecondaryZones,
			Zones:          clientObject.QuotaUsage.Zones,
		},
	}

	return quota, nil
}

func (_ *DBaaSQuotaAdapter) FetchDBaaSQuotas(token string) (model.Quota[model.DatabaseQuota], error) {
	var quota model.Quota[model.DatabaseQuota]

	request, err := http.NewRequest(
		http.MethodGet,
		"https://api.ionos.com/databases/quota",
		nil,
	)
	if err != nil {
		return quota, err
	}
	request.Header.Add("Authorization", token)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return quota, err
	}

	if response.StatusCode != 200 {
		return quota, errors.New("Error while getting DBaaS quotas " + fmt.Sprint(response.StatusCode))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return quota, err
	}

	var responseObject DBaaSResponse
	_ = json.Unmarshal(body, &responseObject)

	quota = model.Quota[model.DatabaseQuota]{
		Limit: model.DatabaseQuota{
			MongoDB:  responseObject.QuotaLimits.CountMongoclustersDbaasIonosCom,
			Postgres: responseObject.QuotaLimits.CountPostgresclustersDbaasIonosCom,
			Cpu:      responseObject.QuotaLimits.Cpu,
			Memory:   responseObject.QuotaLimits.Memory,
			Storage:  responseObject.QuotaLimits.Storage,
		},
		Usage: model.DatabaseQuota{
			MongoDB:  responseObject.QuotaUsage.CountMongoclustersDbaasIonosCom,
			Postgres: responseObject.QuotaUsage.CountPostgresclustersDbaasIonosCom,
			Cpu:      responseObject.QuotaUsage.Cpu,
			Memory:   responseObject.QuotaUsage.Memory,
			Storage:  responseObject.QuotaUsage.Storage,
		},
	}
	return quota, nil
}
