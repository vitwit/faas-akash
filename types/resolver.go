package types

import (
	"fmt"
	"github.com/fatih/color"
	"net/url"
)

// Resolver is responsible for locating the function, when we access the faas-url
func (ir InvokeResolver) Resolve(functionName string) (url.URL, error) {
	fmt.Println("Resolve: ", functionName)
	service, ok := ir.ServiceMap[functionName]

	if !ok {
		color.Red("not found: %s", functionName)
		return url.URL{}, fmt.Errorf("not found")
	}

	fmt.Println(functionName, "=", service)

	//const watchdogPort = 8080

	urlStr := fmt.Sprintf("http://%s", service.URL)

	color.Green("url from akash: %s", urlStr)
	urlRes, err := url.Parse(urlStr)
	if err != nil {
		return url.URL{}, err
	}

	return *urlRes, nil
}
