package akash

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	akashTypes "github.com/vitwit/faas-akash/types"
)

func reverseProxy(proxyURL *url.URL, fnName string, w http.ResponseWriter, r *http.Request) {

	proxy := httputil.NewSingleHostReverseProxy(proxyURL)

	r.URL.Host = proxyURL.Host
	r.URL.Scheme = proxyURL.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = proxyURL.Host
	prefixToTrim := fmt.Sprintf("/function/%s", fnName)

	// trim the prefix, you'll be left with the suffix
	suffix := strings.TrimPrefix(r.URL.Path, prefixToTrim)
	r.URL.Path = suffix
	//r.URL.ForceQuery = proxyURL.ForceQuery
	//r.URL.RawQuery = proxyURL.RawQuery

	log.Printf("proxy request to url: %s", r.URL.String())
	log.Printf("proxy request to url: %s", r.URL.Path)

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

		reverseProxy(&fnURL, fnName, w, r)
	}
}
