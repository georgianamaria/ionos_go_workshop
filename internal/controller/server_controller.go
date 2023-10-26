package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"workshop_demo/internal/server"
	"workshop_demo/internal/service"
	"workshop_demo/internal/util"
)

type ServerController struct {
	quotaService service.Quota
}

func NewServer(quotaService service.Quota) ServerController {
	return ServerController{
		quotaService: quotaService,
	}
}

func (s *ServerController) GetHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *ServerController) GetQuotas(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	quotas, err := s.quotaService.GetQuotas(token)

	response := server.Quotas{
		DBaaS: &server.DBaaSQuota{
			Limits: server.DBaaSQuotaFields{
				CPU:              util.StringToInt64(quotas.DBaaSResponse.Limit.Cpu),
				Memory:           util.StringToInt64(quotas.DBaaSResponse.Limit.Memory),
				MongoClusters:    util.StringToInt64(quotas.DBaaSResponse.Limit.MongoDB),
				PostgresClusters: util.StringToInt64(quotas.DBaaSResponse.Limit.Postgres),
				Storage:          util.StringToInt64(quotas.DBaaSResponse.Limit.Storage),
			},
			Usage: server.DBaaSQuotaFields{
				CPU:              util.StringToInt64(quotas.DBaaSResponse.Usage.Cpu),
				Memory:           util.StringToInt64(quotas.DBaaSResponse.Usage.Memory),
				MongoClusters:    util.StringToInt64(quotas.DBaaSResponse.Usage.MongoDB),
				PostgresClusters: util.StringToInt64(quotas.DBaaSResponse.Usage.Postgres),
				Storage:          util.StringToInt64(quotas.DBaaSResponse.Usage.Storage),
			},
		},
		DNS: &server.DNSQuota{
			Limits: server.DNSQuotaFields{
				Records:        int64(quotas.DNSResponse.Limit.Records),
				SecondaryZones: int64(quotas.DNSResponse.Limit.SecondaryZones),
				Zones:          int64(quotas.DNSResponse.Limit.Zones),
			},
			Usage: server.DNSQuotaFields{
				Records:        int64(quotas.DNSResponse.Usage.Records),
				SecondaryZones: int64(quotas.DNSResponse.Usage.SecondaryZones),
				Zones:          int64(quotas.DNSResponse.Usage.Zones),
			},
		},
	}

	jsonBody, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Error while marshaling response: ", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBody)
}
