package akash

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/openfaas/faas-provider/types"
	akashTypes "github.com/vitwit/faas-akash/types"
)

func Deploy(serviceMap akashTypes.ServiceMap) func(w http.ResponseWriter, r *http.Request) {
	return Update(serviceMap)
}

func Update(serviceMap akashTypes.ServiceMap) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var body types.FunctionDeployment
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		defer func() {
			_ = r.Body.Close()
			color.Red("request payload: %s", body)
		}()

		// lookup for akash client bin file
		// akash client MUST BE PRESENT in the $PATH
		akashCmd, err := exec.LookPath("akash")
		if err != nil {
			w.WriteHeader(http.StatusPreconditionFailed)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		d, _ := os.Getwd()

		dplFile := fmt.Sprintf("%s/x/riot.yaml", d)

		out, err := exec.Command(akashCmd, "deployment", "create", "-m", "json", dplFile).Output()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		akashOut, err := marshalAkashStdout(fmt.Sprintf("%s", bytes.TrimSpace(out)))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		serviceMap[body.Service] = &akashTypes.AkashDeployments{
			Name:    body.Service,
			URL:     akashOut.Leases[0].Services[0].Hosts,
			LeaseID: akashOut.Leases[0].LeaseID,
			IP:      akashOut.Leases[0].Services[1].Hosts,
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"success": string(out)})
	}
}

func akashDeploy() (io.Writer, error) {
	//check if the akash cli tool is installed
	akashBin, err := exec.LookPath("akash")
	if err != nil {
		return nil, fmt.Errorf("ERROR_AKASH_CLI_NOT_FOUND: %w", err)
	}

	cmdOut, err := exec.Command(akashBin, "deployment", "create", "-m", "json", "./x/riot.yaml").Output()
	if err != nil {
		return nil, fmt.Errorf("error creating deployment on akash network: %w", err)
	}

	return bytes.NewBuffer(cmdOut), nil
}
