package types

import (
	"fmt"
	"github.com/openfaas/faas-provider/types"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v3"
)

type (
	// Config to read from yaml/env
	Config struct {
		Port         int           `yaml:"port"`
		ReadTimeout  time.Duration `yaml:"readTimeout"`
		WriteTimeout time.Duration `yaml:"writeTimeout"`
	}
	AkashDeployments struct {
		Name    string
		URL     string
		LeaseID string
		IP      string
	}
	ServiceMap map[string]*AkashDeployments

	InvokeResolver struct {
		ServiceMap ServiceMap
	}
	AkashStdout struct {
		Leases []struct {
			LeaseID  string `json:"lease_id"`
			Services []struct {
				Available string `json:"available"`
				Hosts     string `json:"hosts"`
				Name      string `json:"name"`
				Total     string `json:"total"`
			} `json:"services"`
		} `json:"leases"`
	}

	FaasAkashDeployments map[string]*AkashDeployments

	FaasError map[string]string

	SDL struct {
		Version    string        `yaml:"version"`
		Services   SDLService    `yaml:"services"`
		Profiles   SDLProfile    `yaml:"profiles"`
		Deployment SDLDeployment `yaml:"deployment"`
	}
	SDLService struct {
		Web SDLWebService `yaml:"web"`
	}
	SDLProfile struct {
		Compute   SDLProfileCompute   `yaml:"compute"`
		Placement SDLProfilePlacement `yaml:"placement"`
	}
	SDLDeployment struct {
		Web SDLDeploymentWeb `yaml:"web"`
	}

	SDLDeploymentWeb struct {
		Global SDLDeploymentGlobal `yaml:"global"`
	}

	SDLDeploymentGlobal struct {
		Profile string `yaml:"profile"`
		Count   int    `yaml:"count"`
	}

	SDLProfileCompute struct {
		Web SDLWebCompute `yaml:"web"`
	}

	SDLWebCompute struct {
		CPU    string `yaml:"cpu"`
		Memory string `yaml:"memory"`
		Disk   string `yaml:"disk"`
	}
	SDLProfilePlacement struct {
		Global SDLGlobalPlacement `yaml:"global"`
	}

	SDLGlobalPlacement struct {
		Pricing SDLPricing `yaml:"pricing"`
	}
	SDLPricing struct {
		Web string `yaml:"web"`
	}
	SDLWebService struct {
		Image  string         `yaml:"image"`
		Expose []SDLWebExpose `yaml:"expose"`
	}
	SDLWebExpose struct {
		Port int                      `yaml:"port"`
		To   []map[string]interface{} `yaml:"to"`
	}

	FunctionDeployment struct {
		Price string
		// Service corresponds to a Service
		Service string `json:"service"`

		// Image corresponds to a Docker image
		Image string `json:"image"`

		// Network is specific to Docker Swarm - default overlay network is: func_functions
		Network string `json:"network"`

		// EnvProcess corresponds to the fprocess variable for your container watchdog.
		EnvProcess string `json:"envProcess"`

		// EnvVars provides overrides for functions.
		EnvVars map[string]string `json:"envVars"`

		// RegistryAuth is the registry authentication (optional)
		// in the same encoded format as Docker native credentials
		// (see ~/.docker/config.json)
		RegistryAuth string `json:"registryAuth,omitempty"`

		// Constraints are specific to back-end orchestration platform
		Constraints []string `json:"constraints"`

		// Secrets list of secrets to be made available to function
		Secrets []string `json:"secrets"`

		// Labels are metadata for functions which may be used by the
		// back-end for making scheduling or routing decisions
		Labels *map[string]string `json:"labels"`

		// Annotations are metadata for functions which may be used by the
		// back-end for management, orchestration, events and build tasks
		Annotations *map[string]string `json:"annotations"`

		// Limits for function
		Limits *types.FunctionResources `json:"limits"`

		// Requests of resources requested by function
		Requests *types.FunctionResources `json:"requests"`

		// ReadOnlyRootFilesystem removes write-access from the root filesystem
		// mount-point.
		ReadOnlyRootFilesystem bool `json:"readOnlyRootFilesystem"`

		// Namespace for the function to be deployed into
		Namespace string `json:"namespace,omitempty"`
	}
)

func (ad *AkashDeployments) ParseProviderID() (string, error) {

	// leaseID format deployment_id/group/2/provider_id
	if ad.LeaseID == "" {
		return "", fmt.Errorf("%s", "akash network lease id is empty")
	}

	parts := strings.Split(ad.LeaseID, "/")

	// Split on lease id should give four sub-strings
	if len(parts) != 4 {
		return "", nil
	}

	// parts[3] is the deploymentID
	return parts[3], nil
}

func (ad *AkashDeployments) ParseDeploymentID() (string, error) {
	// leaseID format deployment_id/group/2/provider_id
	if ad.LeaseID == "" {
		return "", fmt.Errorf("%s", "akash network lease id is empty")
	}

	parts := strings.Split(ad.LeaseID, "/")

	// Split on lease id should give four sub-strings
	if len(parts) != 4 {
		return "", nil
	}

	// parts[0] is the deploymentID
	return parts[0], nil
}

func (f *FaasError) Error() string {
	return fmt.Sprintf("%s", f)
}

// necessary checks for the faas-akash configurations
func (c *Config) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Port, validation.Required, validation.Max(9000), validation.Min(3000)),
		validation.Field(&c.ReadTimeout, validation.Required),
		validation.Field(&c.WriteTimeout, validation.Required),
	)
}
