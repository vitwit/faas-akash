package types

import (
	"fmt"
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
