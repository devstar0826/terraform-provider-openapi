package main

import (
	"fmt"
	"net/http"
	"reflect"

	"io/ioutil"

	"strconv"

	"github.com/dikhan/http_goclient"
	"github.com/go-openapi/spec"
	"github.com/hashicorp/terraform/helper/schema"
)

type resourceFactory struct {
	httpClient   http_goclient.HttpClient
	ResourceInfo resourceInfo
}

func (r resourceFactory) createSchemaResource() (*schema.Resource, error) {
	s, err := r.ResourceInfo.createTerraformResourceSchema()
	if err != nil {
		return nil, err
	}
	return &schema.Resource{
		Schema: s,
		Create: r.create,
		Read:   r.read,
		Delete: r.delete,
		Update: r.update,
	}, nil
}

func (r resourceFactory) checkHTTPStatusCode(res *http.Response, expectedHTTPStatusCodes []int) error {
	if !responseContainsExpectedStatus(expectedHTTPStatusCodes, res.StatusCode) {
		switch res.StatusCode {
		case http.StatusUnauthorized:
			return fmt.Errorf("HTTP Reponse Status Code %d - Unauthorized: API access is denied due to invalid credentials", res.StatusCode)
		default:
			b, _ := ioutil.ReadAll(res.Body)
			if len(b) > 0 {
				return fmt.Errorf("HTTP Reponse Status Code %d not matching expected one %v. Response Body = %s", res.StatusCode, expectedHTTPStatusCodes, string(b))
			}
			return fmt.Errorf("HTTP Reponse Status Code %d not matching expected one %v", res.StatusCode, expectedHTTPStatusCodes)
		}

	}
	return nil
}

func (r resourceFactory) create(data *schema.ResourceData, i interface{}) error {
	input := r.getPayloadFromData(data)
	output := map[string]interface{}{}

	headers, url := r.prepareAPIKeyAuthentication(r.ResourceInfo.createPathInfo.Post, i.(providerConfig), r.ResourceInfo.getResourceURL())
	res, err := r.httpClient.PostJson(url, headers, input, &output)
	if err != nil {
		return err
	}
	if err := r.checkHTTPStatusCode(res, []int{http.StatusCreated, http.StatusAccepted}); err != nil {
		return fmt.Errorf("POST %s failed: %s", url, err)
	}

	if output["id"] == nil {
		return fmt.Errorf("object returned from api is missing mandatory property 'id'")
	}

	switch output["id"].(type) {
	case int:
		data.SetId(strconv.Itoa(output["id"].(int)))
	case float64:
		data.SetId(strconv.Itoa(int(output["id"].(float64))))
	default:
		data.SetId(output["id"].(string))
	}
	r.updateResourceState(output, data)

	return nil
}

func (r resourceFactory) read(data *schema.ResourceData, i interface{}) error {
	output, err := r.readRemote(data.Id(), i.(providerConfig))
	if err != nil {
		return err
	}
	return r.updateResourceState(output, data)
}

func (r resourceFactory) readRemote(id string, config providerConfig) (map[string]interface{}, error) {
	output := map[string]interface{}{}
	headers, url := r.prepareAPIKeyAuthentication(r.ResourceInfo.pathInfo.Get, config, r.ResourceInfo.getResourceIDURL(id))
	res, err := r.httpClient.Get(url, headers, &output)
	if err != nil {
		return nil, err
	}
	if err := r.checkHTTPStatusCode(res, []int{http.StatusOK}); err != nil {
		return nil, fmt.Errorf("GET %s failed: %s", url, err)
	}
	return output, nil
}

func (r resourceFactory) update(data *schema.ResourceData, i interface{}) error {
	if r.ResourceInfo.pathInfo.Put == nil {
		return fmt.Errorf("%s resource does not support PUT opperation, check the swagger file exposed on '%s'", r.ResourceInfo.name, r.ResourceInfo.host)
	}
	input := r.getPayloadFromData(data)
	output := map[string]interface{}{}

	if err := r.checkImmutableFields(data, i); err != nil {
		return err
	}

	headers, url := r.prepareAPIKeyAuthentication(r.ResourceInfo.pathInfo.Put, i.(providerConfig), r.ResourceInfo.getResourceIDURL(data.Id()))
	res, err := r.httpClient.PutJson(url, headers, input, &output)
	if err != nil {
		return err
	}
	if err := r.checkHTTPStatusCode(res, []int{http.StatusOK}); err != nil {
		return fmt.Errorf("UPDATE %s failed: %s", url, err)
	}
	return r.updateResourceState(output, data)
}

func (r resourceFactory) delete(data *schema.ResourceData, i interface{}) error {
	if r.ResourceInfo.pathInfo.Delete == nil {
		return fmt.Errorf("%s resource does not support DELETE opperation, check the swagger file exposed on '%s'", r.ResourceInfo.name, r.ResourceInfo.host)
	}
	headers, url := r.prepareAPIKeyAuthentication(r.ResourceInfo.pathInfo.Delete, i.(providerConfig), r.ResourceInfo.getResourceIDURL(data.Id()))
	res, err := r.httpClient.Delete(url, headers)
	if err != nil {
		return err
	}
	if err := r.checkHTTPStatusCode(res, []int{http.StatusNoContent}); err != nil {
		return fmt.Errorf("DELETE %s failed: %s", url, err)
	}
	return nil
}

func (r resourceFactory) checkImmutableFields(updated *schema.ResourceData, i interface{}) error {
	var remoteData map[string]interface{}
	var err error
	if remoteData, err = r.readRemote(updated.Id(), i.(providerConfig)); err != nil {
		return err
	}
	for _, immutablePropertyName := range r.ResourceInfo.getImmutableProperties() {
		if updated.Get(immutablePropertyName) != remoteData[immutablePropertyName] {
			// Rolling back data so tf values are not stored in the state file; otherwise terraform would store the
			// data inside the updated (*schema.ResourceData) in the state file
			r.updateResourceState(remoteData, updated)
			return fmt.Errorf("property %s is immutable and therefore can not be updated. Update operation was aborted; no updates were performed", immutablePropertyName)
		}
	}
	return nil
}

func (r resourceFactory) updateResourceState(input map[string]interface{}, data *schema.ResourceData) error {
	for propertyName, propertyValue := range input {
		if propertyName == "id" {
			continue
		}
		if err := data.Set(propertyName, propertyValue); err != nil {
			return err
		}
	}
	return nil
}

func (r resourceFactory) getPayloadFromData(data *schema.ResourceData) map[string]interface{} {
	input := map[string]interface{}{}
	for propertyName, property := range r.ResourceInfo.schemaDefinition.Properties {
		// ReadOnly properties are not considered for the payload data
		if propertyName == "id" || property.ReadOnly {
			continue
		}
		switch reflect.TypeOf(data.Get(propertyName)).Kind() {
		case reflect.Slice:
			input[propertyName] = data.Get(propertyName).([]interface{})
		case reflect.String:
			input[propertyName] = data.Get(propertyName).(string)
		case reflect.Int:
			input[propertyName] = data.Get(propertyName).(int)
		case reflect.Float64:
			input[propertyName] = data.Get(propertyName).(float64)
		case reflect.Bool:
			input[propertyName] = data.Get(propertyName).(bool)
		}
	}
	return input
}

func (r resourceFactory) authRequired(operation *spec.Operation, providerConfig providerConfig) (bool, string) {
	for _, operationSecurityPolicy := range operation.Security {
		for operationSecurityDefName := range operationSecurityPolicy {
			for providerSecurityDefName := range providerConfig.SecuritySchemaDefinitions {
				if operationSecurityDefName == providerSecurityDefName {
					return true, operationSecurityDefName
				}
			}
		}
	}
	return false, ""
}

func (r resourceFactory) prepareAPIKeyAuthentication(operation *spec.Operation, providerConfig providerConfig, url string) (map[string]string, string) {
	if required, securityDefinitionName := r.authRequired(operation, providerConfig); required {
		headers := map[string]string{}
		if &providerConfig != nil {
			securitySchemaDefinition := providerConfig.SecuritySchemaDefinitions[securityDefinitionName]
			if &securitySchemaDefinition.apiKeyHeader != nil {
				headers[securitySchemaDefinition.apiKeyHeader.name] = securitySchemaDefinition.apiKeyHeader.value
			} else if &securitySchemaDefinition.apiKeyQuery != nil {
				url = fmt.Sprintf("%s?%s=%s", url, securitySchemaDefinition.apiKeyQuery.name, securitySchemaDefinition.apiKeyQuery.value)
			}
		}
		return headers, url
	}
	return nil, url
}
