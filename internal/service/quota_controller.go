package service

import (
	"workshop_demo/internal/model"
	"workshop_demo/internal/port"
)

type Quota struct {
	dbaasQuotas port.DBaaSQuotas
	dnsQuotas   port.DNSQuotas
}

func NewQuota(dbaasQuotas port.DBaaSQuotas, dnsQuotas port.DNSQuotas) Quota {
	return Quota{
		dbaasQuotas: dbaasQuotas,
		dnsQuotas:   dnsQuotas,
	}
}

func (q Quota) GetQuotas(token string) (model.Quotas, error) {
	dbaasQuotas, err := q.dbaasQuotas.FetchDBaaSQuotas(token)
	if err != nil {
		return model.Quotas{}, err
	}

	dnsQuotas, err := q.dnsQuotas.FetchDNSQuotas(token)
	if err != nil {
		return model.Quotas{}, err
	}

	return model.Quotas{
		DBaaSResponse: dbaasQuotas,
		DNSResponse:   dnsQuotas,
	}, nil
}
