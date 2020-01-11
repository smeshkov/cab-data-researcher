package main

import (
	"flag"
	"net/http"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/smeshkov/cab-data-researcher/app"
	"github.com/smeshkov/cab-data-researcher/cfg"
	"github.com/smeshkov/cab-data-researcher/db"
)

var (
	env string
)

func main() {
	cfgFile := flag.String("config", "_resources/config.yml", "Configuration file")
	flag.Parse()

	// init logger
	err := cfg.LogSetup(env)
	if err != nil {
		zap.L().Warn("error in setting up logger", zap.Error(err))
	}
	l := zap.L().With(zap.String("cfg", *cfgFile))
	// flush log entries if any
	defer cfg.LogSync()

	// init configuration
	c, err := cfg.Load(*cfgFile)
	if err != nil {
		l.Fatal("failed to load configuration", zap.Error(err))
	}

	// init database
	l = l.With(zap.String("db_driver", c.DB.Driver))
	cabDB, err := db.NewCabDB(c.DB.Driver, c.DB.DataSource)
	if err != nil {
		l.Fatal("failed to open DB connection", zap.Error(err))
	}
	defer cabDB.Close()

	// todo: init cache

	// init server
	srv := &http.Server{
		ReadHeaderTimeout: c.Server.ReadTimeout,
		IdleTimeout:       c.Server.IdleTimeout,
		ReadTimeout:       c.Server.ReadTimeout,
		WriteTimeout:      c.Server.WriteTimeout,
		Addr:              c.Server.Addr,
		Handler:           app.CreateHandler(env, &c, cabDB),
	}

	l = l.With(
		zap.String("server_name", c.Server.Name), 
		zap.String("server_addr", c.Server.Addr))

	l.Info("starting server")
	if err := srv.ListenAndServe(); err != nil {
		l.Fatal("failed to start server", zap.Error(err))
	}
}
