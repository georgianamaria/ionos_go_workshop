package config

import "github.com/spf13/cobra"

type Options struct {
	Port int
}

func (o *Options) AddFlags(cmd *cobra.Command) {
	cmd.Flags().IntVar(&o.Port, "port", 8080, "Port to start the API on.")
}
