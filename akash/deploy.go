package akash

import (
	"encoding/json"
	"net/http"

	faasTypes "github.com/openfaas/faas-provider/types"
	log "github.com/sirupsen/logrus"
	"github.com/vitwit/faas-akash/types"
)

func DeployHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := r.Body.Close(); err != nil {
				log.Printf("error closing deploy func body: %s", err)
			}
		}()

		log.Info("deployment request")
		defer r.Body.Close()

		var request faasTypes.FunctionDeployment
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			_ = json.NewEncoder(w).Encode(types.FaasError{"error": err.Error()})
			return
		}

		//var out io.Writer
		akashDeploy()

		headers := w.Header()
		copyHeaders(&headers, &r.Header)
		w.WriteHeader(http.StatusOK)
		log.Printf("created deployment on akash network")
	}
}

//func createAkashYaml() error {
//
//}

func akashDeploy() error {
	//akashBin, err := exec.LookPath("akash")
	//if err != nil {
	//	return fmt.Errorf("ERROR_AKASH_CLI_NOT_FOUND: %w", err)
	//}
	//
	//cmd := exec.Command(akashBin, "deployment", "create", "./x/riot.yaml")
	//
	//cmd.Run()
	////go func(cmd *exec.Cmd) {
	////	cmd.Run()
	////}(cmd)
	//
	////out <- cmd.Stdout
	//src, _ := cmd.StdoutPipe()
	//_, _ = io.Copy(out, src)
	//io.t
	return nil
}
