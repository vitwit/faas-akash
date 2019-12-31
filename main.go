package main

import (
	"github.com/fatih/color"
	bootstrap "github.com/openfaas/faas-provider"
	"github.com/openfaas/faas-provider/types"
	"github.com/vitwit/faas-akash/akash"
	"github.com/vitwit/faas-akash/config"
	akashTypes "github.com/vitwit/faas-akash/types"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	dir, err := ioutil.TempDir("./", "tmp")
	if err != nil {
		color.Red("error creating temp directory: %s", err)
		return
	}

	//a non-nil map is required, write ops to a nil map can panic
	svcMap := akashTypes.ServiceMap{}

	// an optional config path can be provided
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-done
		os.RemoveAll(dir)
		os.Exit(1)
	}()

	defer func() {
		os.RemoveAll(dir)
	}()

	Start(cfg, svcMap, dir)
}

// Start faas-akash
func Start(cfg *akashTypes.Config, serviceMap akashTypes.ServiceMap, dir string) {

	//faasConfig := types.FaaSConfig{
	//	MaxIdleConns:        1000,
	//	MaxIdleConnsPerHost: 1000,
	//	ReadTimeout:         cfg.ReadTimeout,
	//	WriteTimeout:        cfg.WriteTimeout,
	//}

	resolver := akashTypes.InvokeResolver{
		ServiceMap: serviceMap,
	}

	bootstrapHandlers := types.FaaSHandlers{
		//FunctionProxy:        proxy.NewHandlerFunc(faasConfig, resolver),
		FunctionProxy:        akash.Proxy(&resolver),
		DeleteHandler:        akash.DeleteHandler(),
		DeployHandler:        akash.Deploy(serviceMap, dir),
		FunctionReader:       akash.ReadHandler(serviceMap),
		UpdateHandler:        akash.Update(serviceMap, dir),
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
