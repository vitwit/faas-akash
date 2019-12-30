package akash

import (
	"encoding/json"
	"net/http"
)

func ListNamespaces() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		list := []string{}
		out, _ := json.Marshal(list)
		_, _ = w.Write(out)
	}
}
