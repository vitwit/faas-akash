package akash

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vitwit/faas-akash/types"
)

//func proxyClient() *http.Client {
//	return &http.Client{
//		Transport: &http.Transport{
//			Proxy: http.ProxyFromEnvironment,
//			DialContext: (&net.Dialer{
//				Timeout:   time.Second * 2,
//				KeepAlive: 1 * time.Second,
//			}).DialContext,
//			IdleConnTimeout:       120 * time.Millisecond,
//			ExpectContinueTimeout: 1500 * time.Millisecond,
//		},
//		Timeout: time.Second * 2,
//		CheckRedirect: func(req *http.Request, via []*http.Request) error {
//			return http.ErrUseLastResponse
//		},
//	}
//}

func Reader(functions types.FaasAkashDeployments) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var res []types.AkashDeployments
		for k, v := range functions {
			res = append(res, types.AkashDeployments{
				Name:    k,
				URL:     v.URL,
				LeaseID: v.LeaseID,
				IP:      v.IP,
			})
		}
		body, _ := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}
}

func ReplicaReader(functions types.FaasAkashDeployments) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		functionName := vars["name"]

		if _, ok := functions[functionName]; ok {
			found := types.AkashDeployments{
				Name: functionName,
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

func UpdateHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	}
}
