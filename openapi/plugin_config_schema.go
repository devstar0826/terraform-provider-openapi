package openapi

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"gopkg.in/yaml.v2"
	"os"
)

// ServiceConfigurations contains the map with all service configurations
type ServiceConfigurations map[string]ServiceConfiguration

// PluginConfigSchema defines the interface/expected behaviour for PluginConfigSchema implementations.
type PluginConfigSchema interface {
	// Validate performs a check to confirm that the schema content is correct
	Validate() error
	// GetServiceConfig returns the service configuration for a given provider name
	GetServiceConfig(providerName string) (ServiceConfiguration, error)
	// GetAllServiceConfigurations returns all the service configuration
	GetAllServiceConfigurations() (ServiceConfigurations, error)
	// GetVersion returns the plugin configuration version
	GetVersion() (string, error)
	// Marshal serializes the value provided into a YAML document
	Marshal() ([]byte, error)
}

// PluginConfigSchemaV1 defines PluginConfigSchema version 1
// Configuration example:
// version: '1'
// services:
//   monitor:
//     swagger-url: http://monitor-api.com/swagger.json
//     insecure_skip_verify: true
//   cdn:
//     swagger-url: https://cdn-api.com/swagger.json
//   vm:
//     swagger-url: http://vm-api.com/swagger.json
type PluginConfigSchemaV1 struct {
	Version  string                      `yaml:"version"`
	Services map[string]*ServiceConfigV1 `yaml:"services"`
}

// NewPluginConfigSchemaV1 creates a new PluginConfigSchemaV1 that implements PluginConfigSchema interface
func NewPluginConfigSchemaV1(services map[string]*ServiceConfigV1) *PluginConfigSchemaV1 {
	return &PluginConfigSchemaV1{
		Version:  "1",
		Services: services,
	}
}

// Validate makes sure that schema data is correct
func (p *PluginConfigSchemaV1) Validate() error {
	if p.Version != "1" {
		return fmt.Errorf("provider configuration version not matching current implementation, please use version '1' of provider configuration specification")
	}
	for k, v := range p.Services {
		if !govalidator.IsURL(v.SwaggerURL) {
			// fall back to try to load the swagger file from disk in case the path provided is a path to a file on disk
			if _, err := os.Stat(v.SwaggerURL); os.IsNotExist(err) {
				return fmt.Errorf("service '%s' found in the provider configuration does not contain a valid SwaggerURL value ('%s'). URL must be either a valid formed URL or a path to an existing swagger file stored in the disk", k, v.SwaggerURL)
			}
		}
	}
	return nil
}

// GetServiceConfig returns the configuration for the given provider name
func (p *PluginConfigSchemaV1) GetServiceConfig(providerName string) (ServiceConfiguration, error) {
	if providerName == "" {
		return nil, fmt.Errorf("providerName not specified")
	}
	serviceConfig := p.Services[providerName]
	if serviceConfig == nil {
		return nil, fmt.Errorf("'%s' not found in provider's services configuration", providerName)
	}
	return serviceConfig, nil
}

// GetVersion returns the plugin configuration version
func (p *PluginConfigSchemaV1) GetVersion() (string, error) {
	return p.Version, nil
}

// GetAllServiceConfigurations returns all the service configuration
func (p *PluginConfigSchemaV1) GetAllServiceConfigurations() (ServiceConfigurations, error) {
	serviceConfigurations := ServiceConfigurations{}
	for k, v := range p.Services {
		serviceConfigurations[k] = v
	}
	return serviceConfigurations, nil
}

// Marshal serializes the value provided into a YAML document
func (p *PluginConfigSchemaV1) Marshal() ([]byte, error) {
	out, err := yaml.Marshal(p)
	return out, err
}
