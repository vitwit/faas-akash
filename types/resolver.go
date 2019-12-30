package types

import (
	"fmt"
	"github.com/fatih/color"
	"net/url"
)

func (ir InvokeResolver) Resolve(functionName string) (url.URL, error) {
	fmt.Println("Resolve: ", functionName)
	serviceIP, ok := ir.ServiceMap[functionName]

	if !ok {
		color.Red("not found: %s", functionName)
		return url.URL{}, fmt.Errorf("not found")
	}

	fmt.Println(functionName, "=", serviceIP)

	//const watchdogPort = 8080

	urlStr := fmt.Sprintf("http://%s", serviceIP.IP)

	urlRes, err := url.Parse(urlStr)
	if err != nil {
		return url.URL{}, err
	}

	return *urlRes, nil
}
