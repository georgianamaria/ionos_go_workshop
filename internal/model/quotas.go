package model

type DNSQuota struct {
	Records        int
	SecondaryZones int
	Zones          int
}

type DatabaseQuota struct {
	MongoDB  string
	Postgres string
	Cpu      string
	Memory   string
	Storage  string
}

type Quotas struct {
	DBaaSResponse Quota[DatabaseQuota]
	DNSResponse   Quota[DNSQuota]
}

type Quota[T any] struct {
	Limit T
	Usage T
}
