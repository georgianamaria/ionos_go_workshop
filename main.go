package main

import (
	"net/http"
	"workshop_demo/internal/adapter"
	"workshop_demo/internal/controller"
	"workshop_demo/internal/server"
	"workshop_demo/internal/service"
)

func main() {
	serverImplementation := controller.NewServer(service.NewQuota(&adapter.DBaaSQuotaAdapter{}, &adapter.DNSQuotaAdapter{}))
	handler := server.Handler(&serverImplementation)

	println("starting server")
	http.ListenAndServe("localhost:8080", handler)
}
