package main

import (
	"errors"
	"fmt"
	"log"

	"net/http"

	"github.com/go-openapi/loads"
	"github.com/hashicorp/terraform/helper/schema"
)

type ProviderFactory struct {
	Name            string
	DiscoveryApiUrl string
}

func (p ProviderFactory) createProvider() *schema.Provider {
	apiSpecAnalyser, err := p.createApiSpecAnalyser()
	if err != nil {
		log.Fatalf("error occurred while retrieving api specification. Error=%s", err)
	}
	provider, err := p.generateProviderFromApiSpec(apiSpecAnalyser)
	if err != nil {
		log.Fatalf("error occurred while creating schema provider. Error=%s", err)
	}
	return provider
}

func (p ProviderFactory) createApiSpecAnalyser() (*ApiSpecAnalyser, error) {
	if p.DiscoveryApiUrl == "" {
		return nil, errors.New("required param 'apiUrl' missing")
	}
	apiSpec, err := loads.JSONSpec(p.DiscoveryApiUrl)
	if err != nil {
		return nil, fmt.Errorf("error occurred when retrieving api spec from %s. Error=%s", p.DiscoveryApiUrl, err)
	}
	apiSpecAnalyser := &ApiSpecAnalyser{apiSpec}
	return apiSpecAnalyser, nil
}

func (p ProviderFactory) generateProviderFromApiSpec(apiSpecAnalyser *ApiSpecAnalyser) (*schema.Provider, error) {
	resourceMap := map[string]*schema.Resource{}
	for resourceName, resourceInfo := range apiSpecAnalyser.getCrudResources() {
		r := ResourceFactory{
			&http.Client{},
			resourceInfo,
		}
		resource := r.createSchemaResource()
		resourceName := p.getProviderResourceName(resourceName)
		resourceMap[resourceName] = resource
	}
	provider := &schema.Provider{
		ResourcesMap: resourceMap,
	}
	return provider, nil
}

func (p ProviderFactory) getProviderResourceName(resourceName string) string {
	fullResourceName := fmt.Sprintf("%s_%s", p.Name, resourceName)
	return fullResourceName
}
