package converter

import (
	"workshop_demo/client"
	"workshop_demo/model"
)

func ConvertModels(
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