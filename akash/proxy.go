package akash

import (
	"fmt"
	"github.com/gorilla/mux"
	akashTypes "github.com/vitwit/faas-akash/types"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func reverseProxy(proxyURL *url.URL, w http.ResponseWriter, r *http.Request) {

	proxy := httputil.NewSingleHostReverseProxy(proxyURL)

	r.URL.Host = proxyURL.Host
	r.URL.Scheme = proxyURL.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = proxyURL.Host

	//r.URL.ForceQuery = proxyURL.ForceQuery
	//r.URL.RawQuery = proxyURL.RawQuery

	proxy.ServeHTTP(w, r)
}

func Proxy(resolver *akashTypes.InvokeResolver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if resolver == nil {
			http.Error(w, errNilResolver, http.StatusPreconditionFailed)
			return
		}

		pathVars := mux.Vars(r)

		fnName := pathVars["name"]
		if fnName == "" {
			http.Error(w, errMissingFunctionName, http.StatusBadRequest)
			return
		}

		fnURL, err := resolver.Resolve(fnName)
		if err != nil {
			msg := fmt.Sprintf(errFunctionNotFound, strings.ToUpper(fnName))
			http.Error(w, msg, http.StatusNotFound)
			return
		}

		log.Printf("proxy request to url: %s", fnURL.String())
		reverseProxy(&fnURL, w, r)
	}
}
