package quotaclient

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
