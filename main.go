package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/openfaas/faas-provider/proxy"
	"github.com/vitwit/faas-akash/akash"
	"github.com/vitwit/faas-akash/types"

	"github.com/fatih/color"
	bootstrap "github.com/openfaas/faas-provider"
	bootstrapTypes "github.com/openfaas/faas-provider/types"
	"github.com/sirupsen/logrus"
	"github.com/vitwit/faas-akash/config"
)

func init() {
	setLogger()
}

func main() {
	// can take an optional config path
	cfg, err := config.Read()
	if err != nil {
		logrus.Printf("ERROR_READ_FAAS_AKASH_CONFIG: %s", err)
		return
	}

	bootstrapConfig := bootstrapTypes.FaaSConfig{
		ReadTimeout:         cfg.ReadTimeout,
		WriteTimeout:        cfg.WriteTimeout,
		MaxIdleConns:        1024,
		MaxIdleConnsPerHost: 1024,
		TCPPort:             &cfg.Port,
		EnableHealth:        true,
		EnableBasicAuth:     false,
	}

	bootstrapHandlers := &bootstrapTypes.FaaSHandlers{
		FunctionProxy:        proxy.NewHandlerFunc(bootstrapConfig, invokeResolver{}),
		DeleteHandler:        akash.DeleteHandler(),
		DeployHandler:        akash.DeployHandler(),
		FunctionReader:       akash.Reader(functions),
		ReplicaReader:        akash.ReplicaReader(functions),
		UpdateHandler:        akash.UpdateHandler(),
		ListNamespaceHandler: akash.ListNamespaces(),
		ReplicaUpdater:       func(w http.ResponseWriter, r *http.Request) {},
		HealthHandler:        func(w http.ResponseWriter, r *http.Request) {},
		InfoHandler:          func(w http.ResponseWriter, r *http.Request) {},
	}

	color.Green("started faas-akash provider on port: %d", cfg.Port)
	bootstrap.Serve(bootstrapHandlers, &bootstrapConfig)
}

func setLogger() {
	logFormat := os.Getenv("LOG_FORMAT")
	logLevel := os.Getenv("LOG_LEVEL")
	if strings.EqualFold(logFormat, "json") {
		logrus.SetFormatter(&logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyMsg:  "message",
				logrus.FieldKeyTime: "@timestamp",
			},
			TimestampFormat: time.RFC3339,
		})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	if level, err := logrus.ParseLevel(logLevel); err == nil {
		logrus.SetLevel(level)
	}
}

type invokeResolver struct{}

var functions map[string]*types.AkashDeployments

func (invokeResolver) Resolve(functionName string) (url.URL, error) {
	ad, ok := functions[functionName]
	if !ok {
		return url.URL{}, fmt.Errorf("not found")
	}

	//@TODO need to implement a lookup mechanism for function ip, url from akash network
	const watchdogPort = 8080

	urlStr := fmt.Sprintf("http://%s:%d", ad.IP, watchdogPort)

	urlRes, err := url.Parse(urlStr)
	if err != nil {
		return url.URL{}, err
	}

	return *urlRes, nil
}
