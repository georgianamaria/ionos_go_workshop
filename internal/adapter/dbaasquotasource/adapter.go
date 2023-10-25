package dbaasquotasource

import (
	"context"
	"strconv"
	"workshop-day-2a/internal/model"
	"workshop-day-2a/internal/port"
	"workshop-day-2a/pkg/quotaclient"
)

var _ port.QuotaSource = (*Adapter)(nil)

type Adapter struct {
}

// FetchQuotaLimits implements port.QuotaSource.
func (*Adapter) FetchQuotaLimits(ctx context.Context, token string) (model.QuotaLimitList, error) {
	resp, err := quotaclient.DBaaSQuotas(token)
	if err != nil {
		return nil, err
	}

	var ql model.QuotaLimitList

	ql = append(ql, &model.QuotaLimit{
		Name:     "mongo-clusters",
		Usage:    strToInt(resp.QuotaUsage.CountMongoclustersDbaasIonosCom),
		Limit:    strToInt(resp.QuotaLimits.CountMongoclustersDbaasIonosCom),
		Provider: "dbaas",
	})
	ql = append(ql, &model.QuotaLimit{
		Name:     "cpu",
		Usage:    strToInt(resp.QuotaUsage.Cpu),
		Limit:    strToInt(resp.QuotaLimits.Cpu),
		Provider: "dbaas",
	})

	return ql, nil
}

// SourceName implements port.QuotaSource.
func (*Adapter) SourceName() string {
	return "DBaaS"
}

func strToInt(s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return v
}
