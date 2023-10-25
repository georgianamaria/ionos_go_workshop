package service

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"workshop-day-2a/internal/api/quotav1"
	"workshop-day-2a/internal/controller"
	"workshop-day-2a/internal/model"
)

var _ quotav1.ServerInterface = (*Quota)(nil)

type Quota struct {
	CompileQuotaCtrl *controller.CompileQuota
}

// GetQuotas implements quotav1.ServerInterface.
func (s *Quota) GetQuotas(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	// controller
	ql, err := s.CompileQuotaCtrl.Do(r.Context(), token)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		slog.Error("Error while compiling: ", err)
		return
	}

	// convert model 2 response
	serverResponse := convertModels(ql)
	jsonBody, err := json.Marshal(serverResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Error while marshaling response: ", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBody)
}

func convertModels(ql model.QuotaLimitList) *quotav1.Quotas {
	return &quotav1.Quotas{
		DBaaS: quotav1.DBaaSUsageLimits{
			Limits: quotav1.DBaaSQuota{
				CPU:              ql.LimitFor("dbaas", "cpu"),
				Memory:           -1,
				MongoClusters:    ql.LimitFor("dbaas", "mongo-clusters"),
				PostgresClusters: -1,
				Storage:          -1,
			},
			Usage: quotav1.DBaaSQuota{
				CPU:              ql.UsageFor("dbaas", "cpu"),
				Memory:           -1,
				MongoClusters:    ql.UsageFor("dbaas", "mongo-clusters"),
				PostgresClusters: -1,
				Storage:          -1,
			},
		},
		DNS: quotav1.DNSUsageLimits{
			Limits: quotav1.DNSQuota{
				Records:        ql.LimitFor("dns", "records"),
				SecondaryZones: -1,
				Zones:          ql.LimitFor("dns", "zones"),
			},
			Usage: quotav1.DNSQuota{
				Records:        ql.UsageFor("dns", "records"),
				SecondaryZones: -1,
				Zones:          ql.UsageFor("dns", "zones"),
			},
		},
	}
}

// GetHealth implements quotav1.ServerInterface.
func (s *Quota) GetHealth(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func strToInt(s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return v
}
