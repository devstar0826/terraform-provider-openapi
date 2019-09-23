package openapi

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/terraform/helper/schema"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewProviderFactory(t *testing.T) {
	Convey("Given a provider name, an analyser and the service config", t, func() {
		providerName := "provider"
		analyser := &specAnalyserStub{}
		serviceConfig := &ServiceConfigV1{}
		Convey("When newProviderFactory is called ", func() {
			p, err := newProviderFactory(providerName, analyser, serviceConfig)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the provider returned should NOT be nil", func() {
				So(p, ShouldNotBeNil)
			})
		})
	})
	Convey("Given a provider name that is empty, an analyser and the service config", t, func() {
		providerName := ""
		analyser := &specAnalyserStub{}
		serviceConfig := &ServiceConfigV1{}
		Convey("When newProviderFactory is called ", func() {
			_, err := newProviderFactory(providerName, analyser, serviceConfig)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("And the provider returned should NOT be nil", func() {
				So(err.Error(), ShouldEqual, "provider name not specified")
			})
		})
	})
	Convey("Given a provider name that is not terraform compliant, an analyser and the service config", t, func() {
		providerName := "someNonTerraformCompliantName"
		analyser := &specAnalyserStub{}
		serviceConfig := &ServiceConfigV1{}
		Convey("When newProviderFactory is called ", func() {
			_, err := newProviderFactory(providerName, analyser, serviceConfig)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("And the provider returned should NOT be nil", func() {
				So(err.Error(), ShouldEqual, "provider name 'someNonTerraformCompliantName' not terraform name compliant, please consider renaming provider to 'some_non_terraform_compliant_name'")
			})
		})
	})
	Convey("Given a provider name, a nil analyser and the service config", t, func() {
		providerName := "provider"
		Convey("When newProviderFactory is called ", func() {
			_, err := newProviderFactory(providerName, nil, nil)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("And the provider returned should NOT be nil", func() {
				So(err.Error(), ShouldEqual, "provider missing an OpenAPI Spec Analyser")
			})
		})
	})

	Convey("Given a provider name, an analyser and a nil service config", t, func() {
		providerName := "provider"
		analyser := &specAnalyserStub{}
		Convey("When newProviderFactory is called ", func() {
			_, err := newProviderFactory(providerName, analyser, nil)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("And the provider returned should NOT be nil", func() {
				So(err.Error(), ShouldEqual, "provider missing the service configuration")
			})
		})
	})
}

func TestCreateProvider(t *testing.T) {
	Convey("Given a provider factory", t, func() {
		apiKeyAuthProperty := newStringSchemaDefinitionPropertyWithDefaults("apikey_auth", "", true, false, "someAuthValue")
		headerProperty := newStringSchemaDefinitionPropertyWithDefaults("header_name", "", true, false, "someHeaderValue")
		p := providerFactory{
			name: "provider",
			specAnalyser: &specAnalyserStub{
				headers: SpecHeaderParameters{
					SpecHeaderParam{
						Name: headerProperty.Name,
					},
				},
				security: &specSecurityStub{
					securityDefinitions: &SpecSecurityDefinitions{
						newAPIKeyHeaderSecurityDefinition(apiKeyAuthProperty.Name, authorization),
					},
					globalSecuritySchemes: createSecuritySchemes([]map[string][]string{
						{
							apiKeyAuthProperty.Name: []string{""},
						},
					}),
				},
				backendConfiguration: &specStubBackendConfiguration{},
			},
			serviceConfiguration: &ServiceConfigStub{},
		}
		Convey("When createProvider is called ", func() {
			p, err := p.createProvider()
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the provider returned should NOT be nil", func() {
				So(p, ShouldNotBeNil)
			})
			Convey("And the provider should have a property for the auth", func() {
				So(p.Schema[apiKeyAuthProperty.Name], ShouldNotBeNil)
			})
			Convey("And the provider should have a property for the header", func() {
				So(p.Schema[headerProperty.Name], ShouldNotBeNil)
			})
			Convey("And the provider should NOT have a property called region since the backend is NOT multi-region", func() {
				So(p.Schema["region"], ShouldBeNil)
			})
		})
	})

	Convey("Given a provider factory with multi-region backend configuration", t, func() {
		apiKeyAuthProperty := newStringSchemaDefinitionPropertyWithDefaults("apikey_auth", "", true, false, "someAuthValue")
		headerProperty := newStringSchemaDefinitionPropertyWithDefaults("header_name", "", true, false, "someHeaderValue")
		p := providerFactory{
			name: "provider",
			specAnalyser: &specAnalyserStub{
				headers: SpecHeaderParameters{
					SpecHeaderParam{
						Name: headerProperty.Name,
					},
				},
				security: &specSecurityStub{
					securityDefinitions: &SpecSecurityDefinitions{
						newAPIKeyHeaderSecurityDefinition(apiKeyAuthProperty.Name, authorization),
					},
					globalSecuritySchemes: createSecuritySchemes([]map[string][]string{
						{
							apiKeyAuthProperty.Name: []string{""},
						},
					}),
				},
				backendConfiguration: &specStubBackendConfiguration{
					host:    "some-service.${region}.api.com",
					regions: []string{"rst", "dub"},
				},
			},
			serviceConfiguration: &ServiceConfigStub{},
		}
		Convey("When createProvider is called ", func() {
			p, err := p.createProvider()
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the provider returned should NOT be nil", func() {
				So(p, ShouldNotBeNil)
			})
			Convey("And the provider should have a property called region since the backend is multi-region", func() {
				So(p.Schema["region"], ShouldNotBeNil)
			})
			Convey("And the provider region property default value should be the first item of the regions list", func() {
				value, err := p.Schema["region"].DefaultFunc()
				So(err, ShouldBeNil)
				So(value.(string), ShouldEqual, "rst")
			})
			Convey("And the provider validator function for the region property should exist and be functional (existing valid regions should NOT throw an error)", func() {
				warns, errors := p.Schema["region"].ValidateFunc("rst", "")
				So(warns, ShouldBeNil)
				So(errors, ShouldBeNil)
			})
			Convey("And the provider validator function for the region property should exist and be functional (non existing regions should throw an error)", func() {
				_, errors := p.Schema["region"].ValidateFunc("nonExisting", "region")
				So(errors, ShouldNotBeNil)
				So(errors[0].Error(), ShouldEqual, "property region value nonExisting is not valid, please make sure the value is one of [rst dub]")
			})
		})
	})

	Convey("Given a provider factory where the specAnalyser has an error", t, func() {
		expectedError := "specAnalyser has an error"
		p := providerFactory{
			name: "provider",
			specAnalyser: &specAnalyserStub{
				error: errors.New(expectedError),
			},
		}
		Convey("When createProvider is called ", func() {
			p, err := p.createProvider()
			Convey("Then the error returned should be as expected", func() {
				So(err.Error(), ShouldEqual, expectedError)
			})
			Convey("And the provider returned should be nil", func() {
				So(p, ShouldBeNil)
			})
		})
	})

	Convey("Given a provider factory where the specAnalyser has an error on the backendConfiguration", t, func() {
		expectedError := "backendConfiguration error"
		p := providerFactory{
			name: "provider",
			specAnalyser: &specAnalyserStub{
				backendConfiguration: &specStubBackendConfiguration{
					err: errors.New(expectedError),
				},
			},
		}
		Convey("When createProvider is called ", func() {
			p, err := p.createProvider()
			Convey("Then the error returned should be as expected", func() {
				So(err.Error(), ShouldEqual, expectedError)
			})
			Convey("And the provider returned should be nil", func() {
				So(p, ShouldBeNil)
			})
		})
	})

	Convey("Given a provider factory where createTerraformProviderResourceMap fails", t, func() {
		expectedError := "resource name can not be empty"
		apiKeyAuthProperty := newStringSchemaDefinitionPropertyWithDefaults("apikey_auth", "", true, false, "someAuthValue")
		headerProperty := newStringSchemaDefinitionPropertyWithDefaults("header_name", "", true, false, "someHeaderValue")
		p := providerFactory{
			name: "provider",
			specAnalyser: &specAnalyserStub{
				resources: []SpecResource{
					&specStubResource{},
				},
				headers: SpecHeaderParameters{
					SpecHeaderParam{Name: headerProperty.Name},
				},
				security: &specSecurityStub{
					securityDefinitions: &SpecSecurityDefinitions{
						newAPIKeyHeaderSecurityDefinition(apiKeyAuthProperty.Name, authorization),
					},
					globalSecuritySchemes: createSecuritySchemes([]map[string][]string{
						{
							apiKeyAuthProperty.Name: []string{""},
						},
					}),
				},
				backendConfiguration: &specStubBackendConfiguration{},
			},
			serviceConfiguration: &ServiceConfigStub{},
		}
		Convey("When createProvider is called ", func() {
			p, err := p.createProvider()
			Convey("Then the error returned should be as expected", func() {
				So(err.Error(), ShouldEqual, expectedError)
			})
			Convey("And the provider returned should be nil", func() {
				So(p, ShouldBeNil)
			})
		})
	})

	Convey("Given a provider factory where createTerraformProviderDataSourceMap fails", t, func() {
		expectedError := "resource name can not be empty"
		apiKeyAuthProperty := newStringSchemaDefinitionPropertyWithDefaults("apikey_auth", "", true, false, "someAuthValue")
		headerProperty := newStringSchemaDefinitionPropertyWithDefaults("header_name", "", true, false, "someHeaderValue")
		p := providerFactory{
			name: "provider",
			specAnalyser: &specAnalyserStub{
				dataSources: []SpecResource{
					&specStubResource{
						name:         "",
						path:         "/v1/resource",
						shouldIgnore: false,
						schemaDefinition: &specSchemaDefinition{
							Properties: specSchemaDefinitionProperties{},
						},
						resourceGetOperation: &specResourceOperation{},
						timeouts:             &specTimeouts{},
					},
				},
				headers: SpecHeaderParameters{
					SpecHeaderParam{Name: headerProperty.Name},
				},
				security: &specSecurityStub{
					securityDefinitions: &SpecSecurityDefinitions{
						newAPIKeyHeaderSecurityDefinition(apiKeyAuthProperty.Name, authorization),
					},
					globalSecuritySchemes: createSecuritySchemes([]map[string][]string{
						{
							apiKeyAuthProperty.Name: []string{""},
						},
					}),
				},
				backendConfiguration: &specStubBackendConfiguration{},
			},
			serviceConfiguration: &ServiceConfigStub{},
		}
		Convey("When createProvider is called ", func() {
			p, err := p.createProvider()
			Convey("Then the error returned should be as expected", func() {
				So(err.Error(), ShouldEqual, expectedError)
			})
			Convey("And the provider returned should be nil", func() {
				So(p, ShouldBeNil)
			})
		})
	})
}

func TestCreateValidateFunc(t *testing.T) {
	Convey("Given a provider factory", t, func() {
		p := providerFactory{}
		allowedValues := []string{"allowed"}
		Convey("When createValidateFunc is called ", func() {
			validate := p.createValidateFunc(allowedValues)
			Convey("Then the validate function should not be nil", func() {
				So(validate, ShouldNotBeNil)
			})
			Convey("And when the function is called with a valid value it should return nil errors", func() {
				warns, errors := validate("allowed", "")
				So(warns, ShouldBeNil)
				So(errors, ShouldBeNil)
			})
			Convey("And when the function is called with an invalid value it should return nil errors", func() {
				warns, errors := validate("notAllowed", "region")
				So(warns, ShouldBeNil)
				So(errors, ShouldNotBeNil)
				So(errors[0].Error(), ShouldEqual, "property region value notAllowed is not valid, please make sure the value is one of [allowed]")
			})
		})
	})
}

func TestCreateTerraformProviderSchema(t *testing.T) {
	Convey("Given a provider factory containing couple properties with commands (that exit with no error) set up", t, func() {
		apiKeyAuthProperty := newStringSchemaDefinitionPropertyWithDefaults("apikey_auth", "", true, false, "someAuthValue")
		headerProperty := newStringSchemaDefinitionPropertyWithDefaults("header_name", "", true, false, "someHeaderValue")

		serviceConfig := &ServiceConfigStub{
			SchemaConfiguration: []*ServiceSchemaPropertyConfigurationStub{
				{
					SchemaPropertyName:   "apikey_auth",
					ExecuteCommandCalled: false,
				},
				{
					SchemaPropertyName:   "header_name",
					ExecuteCommandCalled: false,
				},
			},
		}
		p := providerFactory{
			name: "provider",
			specAnalyser: &specAnalyserStub{
				headers: SpecHeaderParameters{
					SpecHeaderParam{
						Name: headerProperty.Name,
					},
				},
				security: &specSecurityStub{
					securityDefinitions: &SpecSecurityDefinitions{
						newAPIKeyHeaderSecurityDefinition(apiKeyAuthProperty.Name, authorization),
					},
					globalSecuritySchemes: createSecuritySchemes([]map[string][]string{
						{
							apiKeyAuthProperty.Name: []string{""},
						},
					}),
				},
			},
			serviceConfiguration: serviceConfig,
		}
		Convey("When createTerraformProviderSchema is called with a backend configuration that is not multi-region", func() {
			backendConfig := &specStubBackendConfiguration{}
			providerSchema, err := p.createTerraformProviderSchema(backendConfig)
			Convey("Then the expectedValue returned should be true", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the provider schema for the resource should contain the expected attributes", func() {
				So(providerSchema, ShouldContainKey, apiKeyAuthProperty.Name)
				So(providerSchema, ShouldContainKey, headerProperty.Name)
			})
			Convey("And the provider schema default function should not be nil", func() {
				So(providerSchema[apiKeyAuthProperty.Name].DefaultFunc, ShouldNotBeNil)
			})
			Convey("And the provider schema properties commands should have been executed", func() {
				So(serviceConfig.SchemaConfiguration[0].ExecuteCommandCalled, ShouldBeTrue)
				So(serviceConfig.SchemaConfiguration[1].ExecuteCommandCalled, ShouldBeTrue)
			})
		})
	})

	Convey("Given a provider factory containing a property with command (that exit with error) set up", t, func() {
		apiKeyAuthProperty := newStringSchemaDefinitionPropertyWithDefaults("apikey_auth", "", true, false, "someAuthValue")
		expectedError := "some error executing the command"
		serviceConfig := &ServiceConfigStub{
			SchemaConfiguration: []*ServiceSchemaPropertyConfigurationStub{
				{
					SchemaPropertyName:   "apikey_auth",
					ExecuteCommandCalled: false,
					Err:                  fmt.Errorf(expectedError),
				},
			},
		}
		p := providerFactory{
			name: "provider",
			specAnalyser: &specAnalyserStub{
				security: &specSecurityStub{
					securityDefinitions: &SpecSecurityDefinitions{
						newAPIKeyHeaderSecurityDefinition(apiKeyAuthProperty.Name, authorization),
					},
					globalSecuritySchemes: createSecuritySchemes([]map[string][]string{
						{
							apiKeyAuthProperty.Name: []string{""},
						},
					}),
				},
			},
			serviceConfiguration: serviceConfig,
		}
		Convey("When createTerraformProviderSchema is called with a backend configuration that is not multi-region", func() {
			backendConfig := &specStubBackendConfiguration{}
			_, err := p.createTerraformProviderSchema(backendConfig)
			Convey("Then the expectedValue returned should NOT be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("And the error message should be the expected one", func() {
				So(err.Error(), ShouldEqual, expectedError)
			})
			Convey("And the provider schema properties commands should have been executed", func() {
				So(serviceConfig.SchemaConfiguration[0].ExecuteCommandCalled, ShouldBeTrue)
			})
		})
	})

	Convey("Given a provider factory that is configured with security definitions that are not all part of the global schemes", t, func() {
		var globalSecurityDefinitionName = "apiKeyAuth"
		var otherSecurityDefinitionName = "otherSecurityDefinitionName"
		p := providerFactory{
			name: "provider",
			specAnalyser: &specAnalyserStub{
				headers: SpecHeaderParameters{},
				security: &specSecurityStub{
					securityDefinitions: &SpecSecurityDefinitions{
						newAPIKeyHeaderSecurityDefinition(globalSecurityDefinitionName, authorization),
						newAPIKeyHeaderSecurityDefinition(otherSecurityDefinitionName, "Authorization2"),
					},
					globalSecuritySchemes: createSecuritySchemes([]map[string][]string{
						{
							globalSecurityDefinitionName: []string{""},
						},
					}),
				},
			},
			serviceConfiguration: &ServiceConfigStub{},
		}
		Convey("When createTerraformProviderSchema is called with a backend configuration that is not multi-region", func() {
			backendConfig := &specStubBackendConfiguration{}
			providerSchema, err := p.createTerraformProviderSchema(backendConfig)
			Convey("Then the expectedValue returned should be true", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the provider schema for the resource should contain the expected attributes with names automatically converted to be compliant", func() {
				So(providerSchema, ShouldContainKey, "api_key_auth")
				So(providerSchema, ShouldContainKey, "other_security_definition_name")
			})
			Convey("And the api_key_auth should be required as it's a global scheme", func() {
				So(providerSchema["api_key_auth"].Required, ShouldBeTrue)
			})
			Convey("And the other_security_definition_name should be optional as it's not referred in the global schemes", func() {
				So(providerSchema["other_security_definition_name"].Optional, ShouldBeTrue)
			})
			Convey("And the provider schema default function for all the properties", func() {
				So(providerSchema["api_key_auth"].DefaultFunc, ShouldNotBeNil)
				So(providerSchema["other_security_definition_name"].DefaultFunc, ShouldNotBeNil)
			})
		})
	})

	Convey("Given a provider factory", t, func() {
		p := providerFactory{
			name: "provider",
			specAnalyser: &specAnalyserStub{
				headers: SpecHeaderParameters{},
				security: &specSecurityStub{
					securityDefinitions:   &SpecSecurityDefinitions{},
					globalSecuritySchemes: createSecuritySchemes([]map[string][]string{}),
				},
			},
			serviceConfiguration: &ServiceConfigStub{},
		}
		Convey("When createTerraformProviderSchema is called with a backend configuration that IS multi-region", func() {
			multiRegionHost := "api.${region}.server.com"
			expectedDefaultRegion := "rst1"
			backendConfig := newStubBackendMultiRegionConfiguration(multiRegionHost, []string{expectedDefaultRegion})
			providerSchema, err := p.createTerraformProviderSchema(backendConfig)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the provider schema for the resource should contain the region property", func() {
				So(providerSchema, ShouldContainKey, providerPropertyRegion)
			})

			Convey("And the provider region property default function should not be nil", func() {
				So(providerSchema[providerPropertyRegion].DefaultFunc, ShouldNotBeNil)
			})
			Convey("And the provider schema region property should match the first element of the regions array", func() {
				value, err := providerSchema[providerPropertyRegion].DefaultFunc()
				So(err, ShouldBeNil)
				So(value, ShouldEqual, expectedDefaultRegion)
			})
		})
	})

	Convey("Given a provider factory with an spec analyser containing one resource (testing endpoints)", t, func() {
		resourceName := "resource_name_v1"
		resource := newSpecStubResource(resourceName, "", false, nil)
		p := providerFactory{
			name: "provider",
			specAnalyser: &specAnalyserStub{
				resources: []SpecResource{resource},
				headers:   SpecHeaderParameters{},
				security: &specSecurityStub{
					securityDefinitions:   &SpecSecurityDefinitions{},
					globalSecuritySchemes: createSecuritySchemes([]map[string][]string{}),
				},
			},
			serviceConfiguration: &ServiceConfigStub{},
		}
		Convey("When createTerraformProviderSchema is called with some backend configuration", func() {
			backendConfig := &specStubBackendConfiguration{}
			providerSchema, err := p.createTerraformProviderSchema(backendConfig)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey(fmt.Sprintf("And the provider schema should contain %s property", providerPropertyEndPoints), func() {
				So(providerSchema, ShouldContainKey, providerPropertyEndPoints)
			})
			Convey(fmt.Sprintf("And the provider schema %s property should not be nil", providerPropertyEndPoints), func() {
				So(providerSchema[providerPropertyEndPoints], ShouldNotBeNil)
			})
		})
	})
	Convey("Given a provider factory with an spec analyser with no resources (testing endpoints)", t, func() {
		p := providerFactory{
			name: "provider",
			specAnalyser: &specAnalyserStub{
				resources: []SpecResource{},
				headers:   SpecHeaderParameters{},
				security: &specSecurityStub{
					securityDefinitions:   &SpecSecurityDefinitions{},
					globalSecuritySchemes: createSecuritySchemes([]map[string][]string{}),
				},
			},
			serviceConfiguration: &ServiceConfigStub{},
		}
		Convey("When createTerraformProviderSchema is called with some backend configuration (non related)", func() {
			backendConfig := &specStubBackendConfiguration{}
			providerSchema, err := p.createTerraformProviderSchema(backendConfig)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey(fmt.Sprintf("And the provider schema should NOT contain the %s property", providerPropertyEndPoints), func() {
				So(providerSchema, ShouldNotContainKey, providerPropertyEndPoints)
			})
		})
	})
}

func TestCreateTerraformProviderResourceMap(t *testing.T) {
	Convey("Given a provider factory", t, func() {
		p := providerFactory{
			name: "provider",
			specAnalyser: &specAnalyserStub{
				resources: []SpecResource{
					newSpecStubResource("resource", "/v1/resource", false, &specSchemaDefinition{
						Properties: specSchemaDefinitionProperties{
							idProperty,
							stringProperty,
							intProperty,
							numberProperty,
							boolProperty,
							slicePrimitiveProperty,
							computedProperty,
							optionalProperty,
							sensitiveProperty,
							forceNewProperty,
						},
					}),
				},
			},
		}
		Convey("When createTerraformProviderResourceMap is called ", func() {
			schemaResource, err := p.createTerraformProviderResourceMap()
			expectedResourceName := "provider_resource"
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Then the schema resource should contain the resource", func() {
				So(schemaResource, ShouldContainKey, expectedResourceName)
			})
			Convey("And the schema for the resource should contain the expected attributes", func() {
				So(schemaResource[expectedResourceName].Schema, ShouldContainKey, stringProperty.Name)
				So(schemaResource[expectedResourceName].Schema, ShouldContainKey, computedProperty.Name)
				So(schemaResource[expectedResourceName].Schema, ShouldContainKey, intProperty.Name)
				So(schemaResource[expectedResourceName].Schema, ShouldContainKey, numberProperty.Name)
				So(schemaResource[expectedResourceName].Schema, ShouldContainKey, boolProperty.Name)
				So(schemaResource[expectedResourceName].Schema, ShouldContainKey, slicePrimitiveProperty.Name)
				So(schemaResource[expectedResourceName].Schema, ShouldContainKey, optionalProperty.Name)
				So(schemaResource[expectedResourceName].Schema, ShouldContainKey, sensitiveProperty.Name)
				So(schemaResource[expectedResourceName].Schema, ShouldContainKey, forceNewProperty.Name)
			})
			Convey("And the schema property types should match the expected configuration", func() {
				So(schemaResource[expectedResourceName].Schema[stringProperty.Name].Type, ShouldEqual, schema.TypeString)
				So(schemaResource[expectedResourceName].Schema[intProperty.Name].Type, ShouldEqual, schema.TypeInt)
				So(schemaResource[expectedResourceName].Schema[numberProperty.Name].Type, ShouldEqual, schema.TypeFloat)
				So(schemaResource[expectedResourceName].Schema[boolProperty.Name].Type, ShouldEqual, schema.TypeBool)
				So(schemaResource[expectedResourceName].Schema[slicePrimitiveProperty.Name].Type, ShouldEqual, schema.TypeList)
			})
			Convey("And the schema property options should match the expected configuration", func() {
				So(schemaResource[expectedResourceName].Schema[computedProperty.Name].Computed, ShouldBeTrue)
				So(schemaResource[expectedResourceName].Schema[optionalProperty.Name].Optional, ShouldBeTrue)
				So(schemaResource[expectedResourceName].Schema[sensitiveProperty.Name].Sensitive, ShouldBeTrue)
				So(schemaResource[expectedResourceName].Schema[forceNewProperty.Name].ForceNew, ShouldBeTrue)
			})
		})
	})

	Convey("Given a provider factory with a factory loaded with a resource that should be ignored", t, func() {
		p := providerFactory{
			name: "provider",
			specAnalyser: &specAnalyserStub{
				resources: []SpecResource{
					newSpecStubResource("resource", "/v1/resource", true, nil),
				},
			},
		}
		Convey("When createTerraformProviderResourceMap is called ", func() {
			schemaResource, err := p.createTerraformProviderResourceMap()
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Then the schema resource should contain the resource", func() {
				So(schemaResource, ShouldBeEmpty)
			})
		})
	})
}

func TestConfigureProvider(t *testing.T) {
	Convey("Given a provider factory", t, func() {
		apiKeyAuthProperty := newStringSchemaDefinitionPropertyWithDefaults("apikey_auth", "", true, false, "someAuthValue")
		headerProperty := newStringSchemaDefinitionPropertyWithDefaults("header_name", "", true, false, "someHeaderValue")
		p := providerFactory{
			name: "provider",
			specAnalyser: &specAnalyserStub{
				headers: SpecHeaderParameters{
					SpecHeaderParam{
						Name: headerProperty.Name,
					},
				},
				security: &specSecurityStub{
					securityDefinitions: &SpecSecurityDefinitions{
						newAPIKeyHeaderSecurityDefinition(apiKeyAuthProperty.Name, authorization),
					},
					globalSecuritySchemes: createSecuritySchemes([]map[string][]string{
						{
							apiKeyAuthProperty.Name: []string{""},
						},
					}),
				},
			},
		}
		testProviderSchema := newTestSchema(apiKeyAuthProperty, headerProperty)
		Convey("When configureProvider is called with a backend that is not multi-region and the returned configureFunc is invoked upon ", func() {
			backendConfig := &specStubBackendConfiguration{}
			configureFunc := p.configureProvider(backendConfig)
			client, err := configureFunc(testProviderSchema.getResourceData(t))
			providerClient := client.(*ProviderClient)
			Convey("Then error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the client should implement ClientOpenAPI interface", func() {
				var _ ClientOpenAPI = providerClient
			})
		})
	})
}

func TestCreateProviderConfig(t *testing.T) {
	Convey("Given a provider factory configured with a global header and security scheme", t, func() {
		apiKeyAuthProperty := newStringSchemaDefinitionPropertyWithDefaults("apikey_auth", "", true, false, "someAuthValue")
		headerProperty := newStringSchemaDefinitionPropertyWithDefaults("header_name", "", true, false, "someHeaderValue")
		expectedSecurityDefinitions := SpecSecurityDefinitions{
			newAPIKeyHeaderSecurityDefinition(apiKeyAuthProperty.Name, authorization),
		}
		p := providerFactory{
			name: "provider",
			specAnalyser: &specAnalyserStub{
				headers: SpecHeaderParameters{
					SpecHeaderParam{
						Name: headerProperty.Name,
					},
				},
				security: &specSecurityStub{
					securityDefinitions: &expectedSecurityDefinitions,
					globalSecuritySchemes: createSecuritySchemes([]map[string][]string{
						{
							apiKeyAuthProperty.Name: []string{""},
						},
					}),
				},
			},
		}
		Convey(fmt.Sprintf("When createProviderConfig is called with a resource data containing the values for api key auth property (%s) and a header property (%s)", apiKeyAuthProperty.Default, headerProperty.Default), func() {
			testProviderSchema := newTestSchema(apiKeyAuthProperty, headerProperty)
			providerConfiguration, err := p.createProviderConfig(testProviderSchema.getResourceData(t))
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Then the provider configuration returned should contain the header with its value (coming from the resource schema)", func() {
				So(providerConfiguration.Headers[headerProperty.Name], ShouldEqual, headerProperty.Default)
			})
			Convey("And the provider configuration returned should contain the apiKey security with the correct apiKey header name and value (coming from the resource schema)", func() {
				So(providerConfiguration.SecuritySchemaDefinitions[apiKeyAuthProperty.Name].getContext().(apiKey).name, ShouldEqual, expectedSecurityDefinitions[0].getAPIKey().Name)
				So(providerConfiguration.SecuritySchemaDefinitions[apiKeyAuthProperty.Name].getContext().(apiKey).value, ShouldEqual, apiKeyAuthProperty.Default)
			})
		})
	})

	Convey("Given a provider factory configured with a global header and security scheme that use non terraform compliant names", t, func() {
		apiKeyAuthProperty := newStringSchemaDefinitionPropertyWithDefaults("apikey_auth", "", true, false, "someAuthValue")
		var headerProperty = newStringSchemaDefinitionPropertyWithDefaults("header_name", "", true, false, "someHeaderValue")
		var apiKeyAuthPreferredNonCompliantNameProperty = newStringSchemaDefinitionPropertyWithDefaults("apiKeyAuth", "", true, false, "someAuthValue")
		var headerNonCompliantNameProperty = newStringSchemaDefinitionPropertyWithDefaults("headerName", "", true, false, "someHeaderValue")

		expectedSecurityDefinitions := SpecSecurityDefinitions{
			newAPIKeyHeaderSecurityDefinition(apiKeyAuthPreferredNonCompliantNameProperty.Name, authorization),
		}
		p := providerFactory{
			name: "provider",
			specAnalyser: &specAnalyserStub{
				headers: SpecHeaderParameters{
					SpecHeaderParam{
						Name: headerNonCompliantNameProperty.Name,
					},
				},
				security: &specSecurityStub{
					securityDefinitions: &expectedSecurityDefinitions,
					globalSecuritySchemes: createSecuritySchemes([]map[string][]string{
						{
							apiKeyAuthPreferredNonCompliantNameProperty.Name: []string{""},
						},
					}),
				},
			},
		}
		Convey(fmt.Sprintf("When createProviderConfig is called with a resource data containing the values for api key auth property (%s) and a header property (%s)", apiKeyAuthProperty.Default, headerProperty.Default), func() {
			testProviderSchema := newTestSchema(apiKeyAuthPreferredNonCompliantNameProperty, headerNonCompliantNameProperty)
			providerConfiguration, err := p.createProviderConfig(testProviderSchema.getResourceData(t))
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Then the provider configuration returned should contain the header with its value (coming from the resource schema)", func() {
				// provider config keys are always snake_case
				So(providerConfiguration.Headers["header_name"], ShouldEqual, headerProperty.Default)
			})
			Convey("And the provider configuration returned should contain the apiKey security with the correct apiKey header name and value (coming from the resource schema)", func() {
				// The key values stored in the provider configuration are always terraform compliant names, hence querying 'apiKeyAuth' with its corresponding snake_case name
				So(providerConfiguration.SecuritySchemaDefinitions["api_key_auth"].getContext().(apiKey).name, ShouldEqual, expectedSecurityDefinitions[0].getAPIKey().Name)
				So(providerConfiguration.SecuritySchemaDefinitions["api_key_auth"].getContext().(apiKey).value, ShouldEqual, apiKeyAuthPreferredNonCompliantNameProperty.Default)
			})
		})
	})

	Convey("Given a provider factory configured with a global header that has a preferred terraform name configured", t, func() {
		headerProperty := newStringSchemaDefinitionPropertyWithDefaults("header_name", "", true, false, "someHeaderValue")
		headerPreferredNameProperty := newStringSchemaDefinitionPropertyWithDefaults("header_name", "preferred_header_name", true, false, "someHeaderValue")
		p := providerFactory{
			name: "provider",
			specAnalyser: &specAnalyserStub{
				headers: SpecHeaderParameters{
					SpecHeaderParam{
						Name:          headerPreferredNameProperty.Name,
						TerraformName: headerPreferredNameProperty.PreferredName,
					},
				},
				security: &specSecurityStub{},
			},
		}
		Convey(fmt.Sprintf("When createProviderConfig is called with a resource data containing the values for a header property (%s)", headerProperty.Default), func() {
			testProviderSchema := newTestSchema(headerPreferredNameProperty)
			providerConfiguration, err := p.createProviderConfig(testProviderSchema.getResourceData(t))
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Then the provider configuration returned should contain the header with its value (coming from the resource schema), the key used to look up the value is the actual header name", func() {
				So(providerConfiguration.Headers[headerPreferredNameProperty.PreferredName], ShouldEqual, headerProperty.Default)
			})
		})
	})

}

func TestGetProviderResourceName(t *testing.T) {
	Convey("Given a provider factory", t, func() {
		p := providerFactory{
			name: "provider",
		}
		Convey("When getProviderResourceName is called with a resource name", func() {
			expectedResourceName := "resource"
			providerResourceName, err := p.getProviderResourceName(expectedResourceName)
			Convey("Then the err returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Then the value returned should be", func() {
				So(providerResourceName, ShouldEqual, fmt.Sprintf("%s_%s", p.name, expectedResourceName))
			})
		})
		Convey("When getProviderResourceName is called with a resource name that has version", func() {
			expectedResourceName := "resource_v1"
			providerResourceName, err := p.getProviderResourceName(expectedResourceName)
			Convey("Then the err returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Then the value returned should be", func() {
				So(providerResourceName, ShouldEqual, fmt.Sprintf("%s_%s", p.name, expectedResourceName))
			})
		})
		Convey("When getProviderResourceName is called with an empty resource name", func() {
			expectedResourceName := ""
			_, err := p.getProviderResourceName(expectedResourceName)
			Convey("Then the err returned should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestCreateTerraformProviderDataSourceMap(t *testing.T) {

	testcases := []struct {
		name                 string
		specV2stub           SpecAnalyser
		expectedResourceName string
		expectedError        string
	}{
		{
			name: "happy path",
			specV2stub: &specAnalyserStub{
				dataSources: []SpecResource{newSpecStubResource("resource", "/v1/resource", false, &specSchemaDefinition{})},
			},
			expectedResourceName: "provider_resource",
		},
		{
			name: "getProviderResourceName fails ",
			specV2stub: &specAnalyserStub{
				dataSources: []SpecResource{newSpecStubResource("", "/v1/resource", false, &specSchemaDefinition{})},
			},
			expectedError: "resource name can not be empty",
		},
		{
			name: "createTerraformDataSource fails",
			specV2stub: &specAnalyserStub{
				dataSources: []SpecResource{&specStubResource{
					name: "hello",
					funcGetResourceSchema: func() (*specSchemaDefinition, error) {
						return nil, errors.New("createTerraformDataSource failed")
					},
				}},
			},
			expectedError: "createTerraformDataSource failed",
		},
	}

	for _, tc := range testcases {
		p := providerFactory{
			name:         "provider",
			specAnalyser: tc.specV2stub,
		}
		schemaResource, err := p.createTerraformProviderDataSourceMap()

		if tc.expectedError == "" {
			assert.Nil(t, err)
			assert.Contains(t, schemaResource, tc.expectedResourceName, tc.name)
		} else {
			assert.EqualError(t, err, tc.expectedError)
		}

	}

}
