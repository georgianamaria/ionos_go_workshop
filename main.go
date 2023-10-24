package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
)

const token = ""

func main() {

	http.HandleFunc("/quotas",
		func(w http.ResponseWriter, r *http.Request) {
			dbaasQuotas := DBaaSQuotas()
			dnsResponse := DNSQuotas()

			serverResponse := ServerResponse{
				DBaaSResponse: dbaasQuotas,
				DNSResponse:   dnsResponse,
			}

			jsonBody, err := json.Marshal(serverResponse)
			if err != nil {
				panic(err)
			}

			w.WriteHeader(http.StatusOK)
			w.Write(jsonBody)
		},
	)

	http.ListenAndServe("localhost:8080", nil)
}

func DNSQuotas() DNSResponse {
	request2, err := http.NewRequest(http.MethodGet, "https://dns.de-fra.ionos.com/quota", nil)
	request2.Header.Add("Authorization", "Bearer "+os.Getenv("IONOS_TOKEN"))

	response2, err := http.DefaultClient.Do(request2)
	var responseObject2 DNSResponse
	body2, err := io.ReadAll(response2.Body)

	err = json.Unmarshal(body2, &responseObject2)
	if err != nil {
		panic(err)
	}
	return responseObject2
}

func DBaaSQuotas() DBaaSResponse {
	request, _ := http.NewRequest(
		http.MethodGet,
		"https://api.ionos.com/databases/quota",
		nil,
	)
	request.Header.Add("Authorization", "Bearer "+os.Getenv("IONOS_TOKEN"))

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		slog.Error("Error in response")
		panic(err)
	}
	body, err := io.ReadAll(response.Body)

	var responseObject DBaaSResponse
	_ = json.Unmarshal(body, &responseObject)
	return responseObject
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
	} `json "quotaLimits"`
	QuotaUsage struct {
		Records        int `json:"records"`
		SecondaryZones int `json:"secondaryZones"`
		Zones          int `json:"zones"`
	} `json: "quotaUsage"`
}

type ServerResponse struct {
	DBaaSResponse DBaaSResponse `json:"dbaas"`
	DNSResponse   DNSResponse   `json:"dns"`
}
