package akash

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os/exec"
	"time"

	akashTypes "github.com/vitwit/faas-akash/types"
)

func Deploy(serviceMap akashTypes.ServiceMap, dir string) func(w http.ResponseWriter, r *http.Request) {
	return Update(serviceMap, dir)
}

func Update(serviceMap akashTypes.ServiceMap, dir string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var body akashTypes.FunctionDeployment
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		prefix := fmt.Sprintf("akash-deployment-%s", body.Service)
		f, e := ioutil.TempFile(dir, prefix)
		if e != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, _ := yaml.Marshal(CreateDefaultDeploymentManifest(body.Image))

		f.Write(b)
		color.Red("request payload: %s", body.EnvVars)

		defer func() {
			_ = r.Body.Close()
			f.Close()
		}()

		time.Sleep(time.Second * 5)

		akashOut, err := akashDeploy(f.Name())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		serviceMap[body.Service] = &akashTypes.AkashDeployments{
			Name:    body.Service,
			URL:     akashOut.Leases[0].Services[0].Hosts,
			LeaseID: akashOut.Leases[0].LeaseID,
			IP:      akashOut.Leases[0].Services[1].Hosts,
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
