package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"workshop-day-2a/internal/adapter/dbaasquotasource"
	"workshop-day-2a/internal/adapter/dnsquotasource"
	"workshop-day-2a/internal/api/quotav1"
	"workshop-day-2a/internal/config"
	"workshop-day-2a/internal/controller"
	"workshop-day-2a/internal/port"
	"workshop-day-2a/internal/service"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	service := service.Quota{
		CompileQuotaCtrl: &controller.CompileQuota{
			Sources: []port.QuotaSource{
				&dnsquotasource.Adapter{},
				&dbaasquotasource.Adapter{},
			},
		},
	}

	var opts config.Options
	cmd := &cobra.Command{
		Use:   os.Args[0],
		Short: fmt.Sprintf("quota daemon"),
		Run: func(cmd *cobra.Command, args []string) {
			err := config.InitViperFlags(cmd, args)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}

			handler := quotav1.HandlerWithOptions(&service, quotav1.ChiServerOptions{})
			http.Handle("/", handler)

			fmt.Printf("starting server at :%d\n", opts.Port)
			http.ListenAndServe(fmt.Sprintf(":%d", opts.Port), nil)
		},
	}
	opts.AddFlags(cmd)

	cmd.Execute()
}
