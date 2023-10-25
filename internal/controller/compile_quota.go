package controller

import (
	"context"
	"fmt"

	"workshop-day-2a/internal/model"
	"workshop-day-2a/internal/port"
)

type CompileQuota struct {
	Sources []port.QuotaSource
}

func (c *CompileQuota) Do(ctx context.Context, token string) (model.QuotaLimitList, error) {
	var ql model.QuotaLimitList

	for _, source := range c.Sources {
		newQl, err := source.FetchQuotaLimits(ctx, token)
		if err != nil {
			return nil, fmt.Errorf("%w: error while getting %s quotas: %v",
				model.ErrQuotaUnavailable, source.SourceName(), err)

		}

		ql = append(ql, newQl...)
	}

	return ql, nil
}
