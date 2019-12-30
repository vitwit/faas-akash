package main

import (
	bootstrap "github.com/openfaas/faas-provider"
	"github.com/openfaas/faas-provider/proxy"
	"github.com/openfaas/faas-provider/types"
	"github.com/vitwit/faas-akash/akash"
	"github.com/vitwit/faas-akash/config"
	akashTypes "github.com/vitwit/faas-akash/types"
	"log"
)

func main() {
	//a non-nil map is required, write ops to a nil map can panic
	svcMap := akashTypes.ServiceMap{}

	// an optional config path can be provided
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	Start(cfg, svcMap)
}

// Start faas-akash
func Start(cfg *akashTypes.Config, serviceMap akashTypes.ServiceMap) {

	faasConfig := types.FaaSConfig{
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 1000,
		ReadTimeout:         cfg.ReadTimeout,
		WriteTimeout:        cfg.WriteTimeout,
	}

	resolver := akashTypes.InvokeResolver{
		ServiceMap: serviceMap,
	}

	bootstrapHandlers := types.FaaSHandlers{
		FunctionProxy:        proxy.NewHandlerFunc(faasConfig, resolver),
		DeleteHandler:        akash.DeleteHandler(),
		DeployHandler:        akash.Deploy(serviceMap),
		FunctionReader:       akash.ReadHandler(serviceMap),
		UpdateHandler:        akash.Update(serviceMap),
		ReplicaReader:        akash.ReplicaReader(serviceMap),
		ListNamespaceHandler: akash.ListNamespaces(),
		//ReplicaUpdater:       func(w http.ResponseWriter, r *http.Request) {},
		//HealthHandler:        func(w http.ResponseWriter, r *http.Request) {},
		//InfoHandler:          func(w http.ResponseWriter, r *http.Request) {},
	}

	port := 8081
	bootstrapConfig := types.FaaSConfig{
		ReadTimeout:     cfg.ReadTimeout,
		WriteTimeout:    cfg.WriteTimeout,
		TCPPort:         &port,
		EnableBasicAuth: false,
		EnableHealth:    true,
	}

	log.Printf("TCP port: %d\n", cfg.Port)

	bootstrap.Serve(&bootstrapHandlers, &bootstrapConfig)
}
