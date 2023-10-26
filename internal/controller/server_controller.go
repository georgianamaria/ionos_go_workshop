package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"workshop_demo/internal/service"
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

	jsonBody, err := json.Marshal(quotas)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("Error while marshaling response: ", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBody)
}
