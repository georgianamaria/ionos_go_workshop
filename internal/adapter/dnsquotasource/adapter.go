package dnsquotasource

import (
	"context"
	"workshop-day-2a/internal/model"
	"workshop-day-2a/internal/port"
	"workshop-day-2a/pkg/quotaclient"
)

var _ port.QuotaSource = (*Adapter)(nil)

type Adapter struct {
}

// FetchQuotaLimits implements port.QuotaSource.
func (*Adapter) FetchQuotaLimits(ctx context.Context, token string) (model.QuotaLimitList, error) {
	resp, err := quotaclient.DNSQuotas(token)
	if err != nil {
		return nil, err
	}

	var ql model.QuotaLimitList

	ql = append(ql, &model.QuotaLimit{
		Name:     "zones",
		Usage:    resp.QuotaUsage.Zones,
		Limit:    resp.QuotaLimits.Zones,
		Provider: "dns",
	})
	ql = append(ql, &model.QuotaLimit{
		Name:     "records",
		Usage:    resp.QuotaUsage.Records,
		Limit:    resp.QuotaLimits.Records,
		Provider: "dns",
	})

	return ql, nil
}

// SourceName implements port.QuotaSource.
func (*Adapter) SourceName() string {
	return "DNS"
}
