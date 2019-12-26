package akash

import (
	"net/http"
)

func copyHeaders(destination *http.Header, source *http.Header) {
	for k, v := range *source {
		vClone := make([]string, len(v))
		copy(vClone, v)
		(*destination)[k] = vClone
	}
}
