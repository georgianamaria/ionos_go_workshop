package model

type QuotaLimit struct {
	Name     string
	Usage    int64
	Limit    int64
	Provider string
}

type QuotaLimitList []*QuotaLimit

func (ql QuotaLimitList) UsageFor(provider, name string) int64 {
	for _, q := range ql {
		if q.Provider == provider && q.Name == name {
			return q.Usage
		}
	}
	return 0
}

func (ql QuotaLimitList) LimitFor(provider, name string) int64 {
	for _, q := range ql {
		if q.Provider == provider && q.Name == name {
			return q.Limit
		}
	}
	return 0
}
