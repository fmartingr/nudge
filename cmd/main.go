package main

import (
	"context"

	"github.com/fmartingr/nudge/internal/server"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("nudge")

	viper.SetDefault("LOG_LEVEL", server.DefaultLogLevel)
	loglevel := viper.GetString("LOG_LEVEL")

	viper.SetDefault("PORT", server.DefaultPort)
	port := viper.GetInt("PORT")

	viper.SetDefault("INTERVAL", server.DefaultInterval)
	interval := viper.GetInt("INTERVAL")

	viper.SetDefault("IPS", server.DefaultIPs)
	ips := viper.GetStringSlice("IPS")

	ctx := context.Background()
	conf := server.NewDefaultConfig()
	conf.LogLevel = loglevel
	conf.Interval = interval
	conf.Port = port
	conf.IPs = ips
	conf.Print()
	server := server.NewServer(conf)
	go server.Start(ctx)
	if err := server.WaitStop(); err != nil {
		panic(err)
	}
}
