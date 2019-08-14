package openapi

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/dikhan/terraform-provider-openapi/openapi/openapiutils"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
)

const extTfResourceRegionsFmt = "x-terraform-resource-regions-%s"

// specV2Analyser defines an SpecAnalyser implementation for OpenAPI v2 specification
// Forcing creation of this object via constructor so proper input validation is performed before creating the struct
// instance
type specV2Analyser struct {
	openAPIDocumentURL string
	d                  *loads.Document
}

// newSpecAnalyserV2 creates an instance of specV2Analyser which implements the SpecAnalyser interface
// This implementation provides an analyser that understands an OpenAPI v2 document
func newSpecAnalyserV2(openAPIDocumentFilename string) (*specV2Analyser, error) {
	if openAPIDocumentFilename == "" {
		return nil, errors.New("open api document filename argument empty, please provide the url of the OpenAPI document")
	}
	apiSpec, err := loads.JSONSpec(openAPIDocumentFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the OpenAPI document from '%s' - error = %s", openAPIDocumentFilename, err)
	}
	apiSpec, err = apiSpec.Expanded()
	if err != nil {
		return nil, fmt.Errorf("failed to expand the OpenAPI document from '%s' - error = %s", openAPIDocumentFilename, err)
	}
	return &specV2Analyser{
		d:                  apiSpec,
		openAPIDocumentURL: openAPIDocumentFilename,
	}, nil
}

func (specAnalyser *specV2Analyser) createMultiRegionResources(regions []string, resourceRootPath string, resourceRoot, pathItem spec.PathItem, resourcePayloadSchemaDef *spec.Schema) ([]SpecResource, error) {
	var resources []SpecResource
	for _, regionName := range regions {
		r, err := newSpecV2ResourceWithRegion(regionName, resourceRootPath, *resourcePayloadSchemaDef, resourceRoot, pathItem, specAnalyser.d.Spec().Definitions, specAnalyser.d.Spec().Paths.Paths)
		if err != nil {
			return nil, fmt.Errorf("failed to create a resource with region: %s", err)
		}
		log.Printf("[INFO] multi region resource name = %s, region = '%s'", r.getResourceName(), regionName)
		resources = append(resources, r)
	}
	return resources, nil
}

func (specAnalyser *specV2Analyser) GetTerraformCompliantResources() ([]SpecResource, error) {
	var resources []SpecResource
	start := time.Now()
	spec := specAnalyser.d.Spec()
	paths := spec.Paths
	for resourcePath, pathItem := range paths.Paths {
		resourceRootPath, resourceRoot, resourcePayloadSchemaDef, err := specAnalyser.isEndPointFullyTerraformResourceCompliant(resourcePath)
		if err != nil {
			log.Printf("[DEBUG] resource path '%s' not terraform compliant: %s", resourcePath, err)
			continue
		}

		isMultiRegion, regions, err := specAnalyser.isMultiRegionResource(resourceRoot, specAnalyser.d.Spec().Extensions)
		if err != nil {
			log.Printf("multi region configuration for resource '%s' is not valid: ", err)
			continue
		}
		if isMultiRegion {
			log.Printf("[INFO] resource '%s' is configured with host override AND multi region; creating one reasource per region", resourceRootPath)
			multiRegionResources, err := specAnalyser.createMultiRegionResources(regions, resourceRootPath, *resourceRoot, pathItem, resourcePayloadSchemaDef)
			if err != nil {
				log.Printf("[WARN] ignoring multiregion resource '%s' due to an error: %s", resourceRootPath, err)
				continue
			}
			resources = append(resources, multiRegionResources...)
			continue
		}

		r, err := newSpecV2Resource(resourceRootPath, *resourcePayloadSchemaDef, *resourceRoot, pathItem, specAnalyser.d.Spec().Definitions, specAnalyser.d.Spec().Paths.Paths)
		if err != nil {
			log.Printf("[WARN] ignoring resource '%s' due to an error while creating a creating the SpecV2Resource: %s", resourceRootPath, err)
			continue
		}

		err = specAnalyser.validateSubResourceTerraformCompliance(*r)
		if err != nil {
			log.Printf("[WARN] ignoring subresource name='%s' with rootPath='%s' due to not meeting validation requirements: %s", r.getResourceName(), resourceRootPath, err)
			continue
		}

		log.Printf("[INFO] found terraform compliant resource [name='%s', rootPath='%s', instancePath='%s']", r.getResourceName(), resourceRootPath, resourcePath)
		resources = append(resources, r)
	}
	log.Printf("[INFO] found %d terraform compliant resources (time: %s)", len(resources), time.Since(start))
	return resources, nil
}

func (specAnalyser *specV2Analyser) validateSubResourceTerraformCompliance(r SpecV2Resource) error {
	parentResourceInfo := r.getParentResourceInfo()
	if parentResourceInfo != nil {
		resourcePath := r.Path
		for _, parentInstanceURIs := range parentResourceInfo.parentInstanceURIs {
			if pathExists, _ := specAnalyser.pathExists(parentInstanceURIs); !pathExists {
				return fmt.Errorf("subresource with path '%s' is missing parent path instance definition '%s'", resourcePath, parentInstanceURIs)
			}
		}
		for _, parentURI := range parentResourceInfo.parentURIs {
			parentPathExists, parentPathItem := specAnalyser.pathExists(parentURI)
			if !parentPathExists {
				return fmt.Errorf("subresource with path '%s' is missing parent root path definition '%s'", resourcePath, parentURI)
			}
			parentResource := SpecV2Resource{RootPathItem: parentPathItem}
			if parentResource.shouldIgnoreResource() {
				return fmt.Errorf("subresource with path '%s' contains a parent %s that is marked as ignored, therefore ignoring the subresource too", resourcePath, parentURI)
			}
		}
	}
	return nil
}

func (specAnalyser *specV2Analyser) pathExists(path string) (bool, spec.PathItem) {
	p, exists := specAnalyser.d.Spec().Paths.Paths[path]
	if !exists {
		log.Printf("[WARN] path %s not found, falling back to checking if the path with trailing slash %s/ exists", path, path)
		p, exists = specAnalyser.d.Spec().Paths.Paths[path+"/"]
		if !exists {
			return false, spec.PathItem{}
		}
	}
	return true, p
}

// isMultiRegionResource returns true on ly if:
// - the value is parametrized following the pattern: some.subdomain.${keyword}.domain.com, where ${keyword} must be present in the string, otherwise the resource will not be considered multi region
// - there is a matching 'x-terraform-resource-regions-${keyword}' extension defined in the swagger root level (extensions passed in), where ${keyword} will be the value of the parameter in the above URL
// - and finally the value of the extension is an array of strings containing the different regions where the resource can be created
func (specAnalyser *specV2Analyser) isMultiRegionResource(resourceRoot *spec.PathItem, extensions spec.Extensions) (bool, []string, error) {
	overrideHost := getResourceOverrideHost(resourceRoot.Post)
	if overrideHost == "" {
		return false, nil, nil
	}
	isMultiRegionHost, regex := openapiutils.IsMultiRegionHost(overrideHost)
	if !isMultiRegionHost {
		return false, nil, nil
	}
	region := regex.FindStringSubmatch(overrideHost)
	if len(region) != 5 {
		return false, nil, fmt.Errorf("override host %s provided does not comply with expected regex format", overrideHost)
	}
	regionIdentifier := region[3]
	regionExtensionName := specAnalyser.getResourceRegionExtensionName(regionIdentifier)
	if resourceRegions, exists := openapiutils.StringExtensionExists(extensions, regionExtensionName); exists {
		resourceRegions = strings.Replace(resourceRegions, " ", "", -1)
		regions := strings.Split(resourceRegions, ",")
		if len(regions) < 1 {
			return false, nil, fmt.Errorf("could not find any region for '%s' matching region extension %s: '%s'", regionIdentifier, regionExtensionName, resourceRegions)
		}
		apiRegions := []string{}
		for _, region := range regions {
			apiRegions = append(apiRegions, region)
		}
		if len(apiRegions) < 1 {
			return false, nil, fmt.Errorf("could not build properly the resource region map for '%s' matching region extension %s: '%s'", regionIdentifier, regionExtensionName, resourceRegions)
		}
		return true, apiRegions, nil
	}
	return false, nil, fmt.Errorf("missing matching '%s' root level region extension '%s'", regionIdentifier, regionExtensionName)
}

func (specAnalyser *specV2Analyser) getResourceRegionExtensionName(regionIdentifier string) string {
	return fmt.Sprintf(extTfResourceRegionsFmt, regionIdentifier)
}

func (specAnalyser *specV2Analyser) GetSecurity() SpecSecurity {
	return &specV2Security{
		SecurityDefinitions: specAnalyser.d.Spec().SecurityDefinitions,
		GlobalSecurity:      specAnalyser.d.Spec().Security,
	}
}

// GetAllHeaderParameters gets all the parameters of type headers present in the swagger file and returns the header
// configurations. Currently only the following parameters are supported:
// - root level parameters (not supported)
// - path level parameters (not supported)
// - operation level parameters (supported)
func (specAnalyser *specV2Analyser) GetAllHeaderParameters() (SpecHeaderParameters, error) {
	return getAllHeaderParameters(specAnalyser.d.Spec().Paths.Paths), nil
}

func (specAnalyser *specV2Analyser) GetAPIBackendConfiguration() (SpecBackendConfiguration, error) {
	return newOpenAPIBackendConfigurationV2(specAnalyser.d.Spec(), specAnalyser.openAPIDocumentURL)
}

// isEndPointFullyTerraformResourceCompliant returns true only if:
// - The path given 'resourcePath' is an instance path (e,g: "/users/{username}")
// - The path given has GET operation defined (required). PUT and DELETE are optional
// - The root path for the given path 'resourcePath' is found (e,g: "/users")
// - The root path for the given path 'resourcePath' has mandatory POST operation defined
// - The root path for the given path 'resourcePath' has a parameter of type 'body' with a schema property referencing to an existing definition object
// - The root path POST payload definition and the returned object in the response matches. Similarly, the GET operation should also have the same return object
// - The resource schema definition must contain a field that uniquelly identifies the resource or have a field with the 'x-terraform-id' extension set to true
// For instance, if resourcePath was "/users/{id}" and paths contained the following entries and implementations:
// paths:
//   /v1/users:
//     post:
//		 parameters:
//		 - in: "body"
//		   name: "body"
//		   description: "user to create"
//		   required: true
//		   schema:
//		     $ref: "#/definitions/User"
//		 responses:
//		   201:
//		     description: "successful operation"
//		     schema:
//		       $ref: "#/definitions/User"
//   /v1/users/{id}:
//	   get:
//	     parameters:
//	       - name: "id"
//	         in: "path"
//	         description: "The user id that needs to be fetched"
//	         required: true
//	         type: "string"
//	     responses:
//	       200:
//	      	 description: "successful operation"
//	         schema:
//	           $ref: "#/definitions/User"
// definitions:
//   Users:
//     type: "object"
//     required:
//       - name
//     properties:
//       id:
//         type: "string"
//         readOnly: true
//       name:
//         type: "string"
// then the expected returned value is true. Otherwise if the above criteria is not met, it is considered that
// the resourcePath provided is not terraform resource compliant.
func (specAnalyser *specV2Analyser) isEndPointFullyTerraformResourceCompliant(resourcePath string) (string, *spec.PathItem, *spec.Schema, error) {
	err := specAnalyser.validateInstancePath(resourcePath)
	if err != nil {
		return "", nil, nil, err
	}
	resourceRootPath, resourceRootPathItem, resourceRootPostSchemaDef, err := specAnalyser.validateRootPath(resourcePath)
	if err != nil {
		return "", nil, nil, err
	}
	err = specAnalyser.validateResourceSchemaDefinition(resourceRootPostSchemaDef)
	if err != nil {
		return "", nil, nil, err
	}
	return resourceRootPath, resourceRootPathItem, resourceRootPostSchemaDef, nil
}

func (specAnalyser *specV2Analyser) validateInstancePath(path string) error {
	isResourceInstance, err := specAnalyser.isResourceInstanceEndPoint(path)
	if err != nil {
		return fmt.Errorf("error occurred while checking if path '%s' is a resource instance path", path)
	}
	if !isResourceInstance {
		return fmt.Errorf("path '%s' is not a resource instance path", path)
	}
	endPoint := specAnalyser.d.Spec().Paths.Paths[path]
	if endPoint.Get == nil {
		return fmt.Errorf("resource instance path '%s' missing required GET operation", path)
	}
	return nil
}

func (specAnalyser *specV2Analyser) validateRootPath(resourcePath string) (string, *spec.PathItem, *spec.Schema, error) {
	resourceRootPath, err := specAnalyser.findMatchingResourceRootPath(resourcePath)
	if err != nil {
		return "", nil, nil, err
	}

	postExist := specAnalyser.postDefined(resourceRootPath)
	if !postExist {
		return "", nil, nil, fmt.Errorf("resource root path '%s' missing required POST operation", resourceRootPath)
	}

	resourceRootPathItem, _ := specAnalyser.d.Spec().Paths.Paths[resourceRootPath]
	resourceRootPostOperation := resourceRootPathItem.Post

	resourceRootPostSchemaDef, err := specAnalyser.getBodyParameterBodySchema(resourceRootPostOperation)
	if err != nil {
		return "", nil, nil, fmt.Errorf("resource root path '%s' POST operation validation error: %s", resourceRootPath, err)
	}

	return resourceRootPath, &resourceRootPathItem, resourceRootPostSchemaDef, nil
}

func (specAnalyser *specV2Analyser) validateResourceSchemaDefinition(schema *spec.Schema) error {
	identifier := ""
	for propertyName, property := range schema.Properties {
		if propertyName == "id" {
			identifier = propertyName
			continue
		}
		if exists, useAsIdentifier := property.Extensions.GetBool(extTfID); exists && useAsIdentifier {
			identifier = propertyName
			break
		}
	}
	if identifier == "" {
		return fmt.Errorf("resource schema is missing a property that uniquely identifies the resource, either a property named 'id' or a property with the extension '%s' set to true", extTfID)
	}
	return nil
}

// postIsPresent checks if the given resource has a POST implementation returning true if the path is found
// in paths and the path exposes a POST operation
func (specAnalyser *specV2Analyser) postDefined(resourceRootPath string) bool {
	b, exists := specAnalyser.d.Spec().Paths.Paths[resourceRootPath]
	if !exists || b.Post == nil {
		return false
	}
	return true
}

func (specAnalyser *specV2Analyser) getBodyParameterBodySchema(resourceRootPostOperation *spec.Operation) (*spec.Schema, error) {
	if resourceRootPostOperation == nil {
		return nil, fmt.Errorf("resource root operation does not have a POST operation")
	}
	bodyCounter := 0
	var bodyParameter spec.Parameter
	for _, parameter := range resourceRootPostOperation.Parameters {
		if parameter.In == "body" {
			bodyCounter = bodyCounter + 1
			bodyParameter = parameter
		}
	}

	if bodyCounter <= 0 {
		return nil, fmt.Errorf("resource root operation missing the body parameter")
	}

	if bodyCounter > 1 {
		return nil, fmt.Errorf("resource root operation contains multiple 'body' parameters")
	}

	if bodyParameter.Schema == nil {
		return nil, fmt.Errorf("resource root operation missing the schema for the POST operation body parameter")
	}

	if bodyParameter.Schema.Ref.String() != "" {
		return nil, fmt.Errorf("the operation ref was not expanded properly, check that the ref is valid (no cycles, bogus, etc)")
	}

	if len(bodyParameter.Schema.Properties) > 0 {
		return bodyParameter.Schema, nil
	}
	return nil, fmt.Errorf("POST operation contains an schema with no properties")
}

// resourceInstanceRegex loads up the regex specified in const resourceInstanceRegex
// If the regex is not able to compile the regular expression the function exists calling os.Exit(1) as
// there is the regex is completely busted
func (specAnalyser *specV2Analyser) resourceInstanceRegex() (*regexp.Regexp, error) {
	r, err := regexp.Compile(resourceInstanceRegex)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while compiling the resourceInstanceRegex regex '%s': %s", resourceInstanceRegex, err)
	}
	return r, nil
}

// isResourceInstanceEndPoint checks if the given path is of form /resource/{id}
func (specAnalyser *specV2Analyser) isResourceInstanceEndPoint(p string) (bool, error) {
	r, err := specAnalyser.resourceInstanceRegex()
	if err != nil {
		return false, err
	}
	return r.MatchString(p), nil
}

// findMatchingResourceRootPath returns the corresponding POST root and path for a given end point
// Example: Given 'resourcePath' being "/users/{username}" the result could be "/users" or "/users/" depending on
// how the POST operation (resourceRootPath) of the given resource is defined in swagger.
// If there is no match the returned string will be empty
func (specAnalyser *specV2Analyser) findMatchingResourceRootPath(resourcePath string) (string, error) {
	r, err := specAnalyser.resourceInstanceRegex()
	if err != nil {
		return "", err
	}
	result := r.FindStringSubmatch(resourcePath)
	log.Printf("[DEBUG] resource root path match result - %s", result)
	if len(result) != 2 {
		return "", fmt.Errorf("resource instance path '%s' missing valid resource root path, more than two results returned from match '%s'", resourcePath, result)
	}

	resourceRootPath := result[1] // e,g: /v1/cdns/{id} /v1/cdns/

	if _, exists := specAnalyser.d.Spec().Paths.Paths[resourceRootPath]; exists {
		log.Printf("[DEBUG] found resource root path with trailing '/' - %+s", resourceRootPath)
		return resourceRootPath, nil
	}

	// Handles the case where the swagger file root path does not have a trailing slash in the path
	resourceRootPath = strings.TrimRight(resourceRootPath, "/")
	if _, exists := specAnalyser.d.Spec().Paths.Paths[resourceRootPath]; exists {
		log.Printf("[DEBUG] found resource root path without trailing '/' - %+s", resourceRootPath)
		return resourceRootPath, nil
	}

	return "", fmt.Errorf("resource instance path '%s' missing resource root path", resourcePath)
}
