package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"workshop_demo/client"
	"workshop_demo/model"
)

func main() {

	http.HandleFunc("/quotas",
		func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")

			dbaasQuotas, err := client.DBaaSQuotas(token)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				slog.Error("Error while getting DBaaS quotas: ", err)
				return
			}

			dnsResponse, err := client.DNSQuotas(token)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				slog.Error("Error while getting DNS quotas: ", err)
				return
			}

			serverResponse := convertModels(*dnsResponse, *dbaasQuotas)

			jsonBody, err := json.Marshal(serverResponse)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				slog.Error("Error while marshaling response: ", err)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(jsonBody)
		},
	)

	println("starting server")
	http.ListenAndServe(":8080", nil)
}

func convertModels(
	dns client.DNSResponse,
	dbaas client.DBaaSResponse,
) model.ServerResponse {
	return model.ServerResponse{
		DNSResponse: model.Quota[model.DNSQuota]{
			Limit: dns.QuotaLimits,
			Usage: dns.QuotaUsage,
		},
		DBaaSResponse: model.Quota[model.DatabaseQuota]{
			Limit: dbaas.QuotaLimits,
			Usage: dbaas.QuotaUsage,
		},
	}
}
