package model

type DNSQuota struct {
	Records        int `json:"records"`
	SecondaryZones int `json:"secondaryZones"`
	Zones          int `json:"zones"`
}

type DatabaseQuota struct {
	CountMongoclustersDbaasIonosCom    string `json:"count/mongoclusters.dbaas.ionos.com"`
	CountPostgresclustersDbaasIonosCom string `json:"count/postgresclusters.dbaas.ionos.com"`
	Cpu                                string `json:"cpu"`
	Memory                             string `json:"memory"`
	Storage                            string `json:"storage"`
}

type ServerResponse struct {
	DBaaSResponse Quota[DatabaseQuota] `json:"dbaas"`
	DNSResponse   Quota[DNSQuota]      `json:"dns"`
}

type Quota[T any] struct {
	Limit T `json:"limit"`
	Usage T `json:"usage"`
}
