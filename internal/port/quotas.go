package port

import "workshop_demo/internal/model"

type DNSQuotas interface {
	FetchDNSQuotas(token string) (model.Quota[model.DNSQuota], error)
}

type DBaaSQuotas interface {
	FetchDBaaSQuotas(token string) (model.Quota[model.DatabaseQuota], error)
}
