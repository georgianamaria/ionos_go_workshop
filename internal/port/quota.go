package port

import (
	"context"
	"workshop-day-2a/internal/model"
)

type QuotaSource interface {
	SourceName() string
	FetchQuotaLimits(ctx context.Context, token string) (model.QuotaLimitList, error)
}
