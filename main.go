package main

import (
	"apiserver/config"
	"apiserver/logger"
	"apiserver/model"
	"apiserver/router"
	"apiserver/router/middleware"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfg = pflag.StringP("config", "c", "", "apiserver config file path")
)

func main() {
	pflag.Parse()

	// init config
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	// init db
    model.DB.Init()
    defer model.DB.Close()

	gin.SetMode(viper.GetString("runmode"))

	g := gin.New()

	router.Load(g, // Middlwares.
		middleware.Logging(),
		middleware.RequestId(),)

	// Ping the server to make sure the router is working.
	go func() {
		if err := pingServer(); err != nil {
			logger.Error("The router has no response, or it might took too long to start up", err)
		}
		logger.Info("The router has been deployed successfully")
	}()

	cert := viper.GetString("tls.cert")
	key := viper.GetString("tls.key")

	if key != "" && cert != "" {
		go func ()  {
			logger.Infof("Start to listening the incoming requests on https address: %s", viper.GetString("tls.addr"))
			logger.Infof(http.ListenAndServeTLS(viper.GetString("tls.addr"), cert, key, g).Error())
		}()
	}

	logger.Infof("Start to listening the incoming requests on http address: %s", viper.GetString("addr"))
	logger.Infof(http.ListenAndServe(viper.GetString("addr"), g).Error())
}

// pingServer pings the http server to make sure the router is working.
func pingServer() error {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get(viper.GetString("url") + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		logger.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("cannot connect to the router")
}
