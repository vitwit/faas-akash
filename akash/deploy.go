package akash

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"

	"github.com/fatih/color"
	"github.com/openfaas/faas-provider/types"
	akashTypes "github.com/vitwit/faas-akash/types"
	"gopkg.in/yaml.v2"
)

// @TODO Use update code as deploy code
func Deploy(serviceMap akashTypes.ServiceMap, dir string) func(w http.ResponseWriter, r *http.Request) {
	return Update(serviceMap, dir)
}

// @TODO Implement update
func Update(serviceMap akashTypes.ServiceMap, dir string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var body types.FunctionDeployment
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		prefix := fmt.Sprintf("akash-deployment-%s", body.Service)
		//@TODO removing this file in defer returns an error from akash client
		f, e := ioutil.TempFile(dir, prefix)
		if e != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// marshal the go-struct into yaml-compatible structure
		b, _ := yaml.Marshal(CreateDefaultDeploymentManifest(body.Image))

		// if b is empty, file will be empty too
		// @TODO add checks for proper data flush to file
		_, _ = f.Write(b)

		defer func() {
			_ = r.Body.Close()
			_ = f.Close()
		}()

		akashOut, err := akashDeploy(f.Name())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, out := range akashOut.Leases {
			serviceMap[body.Service] = &akashTypes.AkashDeployments{
				// out.Services[0].Hosts is the url for the akash service
				Name:    body.Service,
				URL:     out.Services[0].Hosts,
				LeaseID: out.LeaseID,
				IP:      out.Services[1].Hosts,
			}
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(serviceMap)
	}
}

func akashDeploy(manifestFile string) (*akashTypes.AkashStdout, error) {
	// lookup for akash client bin file
	// akash client MUST BE PRESENT in the $PATH
	color.Green("creating akash deployment from file: %s", manifestFile)
	akashBin, err := exec.LookPath("akash")
	if err != nil {
		return nil, fmt.Errorf("ERROR_AKASH_CLI_NOT_FOUND: %w", err)
	}

	cmdOut, err := exec.Command(akashBin, "deployment", "create", "-m", "json", manifestFile).Output()
	if err != nil {
		return nil, fmt.Errorf("error creating deployment on akash network: %w", err)
	}

	output, err := marshalAkashStdout(string(cmdOut))
	if err != nil {
		return nil, err
	}

	return output, nil
}
