package akash

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/openfaas/faas-provider/types"
	akashTypes "github.com/vitwit/faas-akash/types"
)

func ReplicaReader(serviceMap akashTypes.ServiceMap) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		functionName := vars["name"]

		if _, ok := serviceMap[functionName]; ok {
			found := types.FunctionStatus{
				Name:              functionName,
				AvailableReplicas: 1,
				Replicas:          1,
			}

			functionBytes, _ := json.Marshal(found)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(functionBytes)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func ReadHandler(serviceMap akashTypes.ServiceMap) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		res := []types.FunctionStatus{}
		for k := range serviceMap {
			res = append(res, types.FunctionStatus{
				Name:     k,
				Replicas: 1,
			})
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(res)
	}
}
