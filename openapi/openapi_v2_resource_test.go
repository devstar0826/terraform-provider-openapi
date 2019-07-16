package openapi

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/go-openapi/jsonreference"
	"github.com/go-openapi/spec"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewSpecV2Resource(t *testing.T) {
	Convey("Given a root path /users/ containing a trailing slash and a root path item item", t, func() {
		path := "/users/"
		rootPathItem := spec.PathItem{
			PathItemProps: spec.PathItemProps{
				Post: &spec.Operation{},
			},
		}
		Convey("When getResourceName method is called", func() {
			schemaDefinitions := map[string]spec.Schema{}
			r, err := newSpecV2Resource(path, spec.Schema{}, rootPathItem, spec.PathItem{}, schemaDefinitions)
			Convey("The the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the value returned should be 'users'", func() {
				So(r.getResourceName(), ShouldEqual, "users")
			})
		})
	})
	Convey("Given a root path /users with no trailing slash and a root path item item", t, func() {
		path := "/users"
		rootPathItem := spec.PathItem{
			PathItemProps: spec.PathItemProps{
				Post: &spec.Operation{},
			},
		}
		Convey("When getResourceName method is called", func() {
			schemaDefinitions := map[string]spec.Schema{}
			r, err := newSpecV2Resource(path, spec.Schema{}, rootPathItem, spec.PathItem{}, schemaDefinitions)
			Convey("The the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the value returned should be 'users'", func() {
				So(r.getResourceName(), ShouldEqual, "users")
			})
		})
	})

	Convey("Given a root path that is versioned such as '/v1/users/' and a root path item item", t, func() {
		path := "/v1/users/"
		rootPathItem := spec.PathItem{
			PathItemProps: spec.PathItemProps{
				Post: &spec.Operation{},
			},
		}
		Convey("When getResourceName method is called", func() {
			schemaDefinitions := map[string]spec.Schema{}
			r, err := newSpecV2Resource(path, spec.Schema{}, rootPathItem, spec.PathItem{}, schemaDefinitions)
			Convey("The the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the value returned should be 'users_v1'", func() {
				So(r.getResourceName(), ShouldEqual, "users_v1")
			})
		})
	})

	Convey("Given a root path that is versioned with number higher than 9 such as '/v12/users/' and a root path item item", t, func() {
		path := "/v12/users/"
		rootPathItem := spec.PathItem{
			PathItemProps: spec.PathItemProps{
				Post: &spec.Operation{},
			},
		}
		Convey("When getResourceName method is called", func() {
			schemaDefinitions := map[string]spec.Schema{}
			r, err := newSpecV2Resource(path, spec.Schema{}, rootPathItem, spec.PathItem{}, schemaDefinitions)
			Convey("The the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the value returned should be 'users_v1'", func() {
				So(r.getResourceName(), ShouldEqual, "users_v12")
			})
		})
	})

	Convey("Given a root path such as '/v1/something/users' and a root path item item", t, func() {
		path := "/v1/something/users"
		rootPathItem := spec.PathItem{
			PathItemProps: spec.PathItemProps{
				Post: &spec.Operation{},
			},
		}
		Convey("When getResourceName method is called", func() {
			schemaDefinitions := map[string]spec.Schema{}
			r, err := newSpecV2Resource(path, spec.Schema{}, rootPathItem, spec.PathItem{}, schemaDefinitions)
			Convey("The the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the value returned should still be 'users'", func() {
				So(r.getResourceName(), ShouldEqual, "users")
			})
		})
	})

	Convey("Given a root path which has path parameters '/api/v1/nodes/{name}/proxy' and a root path item item", t, func() {
		path := "/api/v1/nodes/{name}/proxy"
		rootPathItem := spec.PathItem{
			PathItemProps: spec.PathItemProps{
				Post: &spec.Operation{},
			},
		}
		Convey("When getResourceName method is called", func() {
			schemaDefinitions := map[string]spec.Schema{}
			r, err := newSpecV2Resource(path, spec.Schema{}, rootPathItem, spec.PathItem{}, schemaDefinitions)
			Convey("The the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the value returned should still be 'nodes_v1_proxy'", func() {
				So(r.getResourceName(), ShouldEqual, "nodes_v1_proxy")
			})
		})
	})
	Convey("Given a root path '/users' and the create operation has the extension 'x-terraform-resource-name' and a root path item item", t, func() {
		path := "/v1/users"
		expectedResourceName := "user"
		rootPathItem := spec.PathItem{
			PathItemProps: spec.PathItemProps{
				Post: &spec.Operation{
					VendorExtensible: spec.VendorExtensible{
						Extensions: spec.Extensions{
							extTfResourceName: expectedResourceName,
						},
					},
				},
			},
		}
		Convey("When getResourceName method is called", func() {
			schemaDefinitions := map[string]spec.Schema{}
			r, err := newSpecV2Resource(path, spec.Schema{}, rootPathItem, spec.PathItem{}, schemaDefinitions)
			Convey("The the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			resourceName := r.getResourceName()
			expectedTerraformName := fmt.Sprintf("%s_v1", expectedResourceName)
			Convey(fmt.Sprintf("And the value returned should still be '%s'", expectedTerraformName), func() {
				So(resourceName, ShouldEqual, expectedTerraformName)
			})
		})
	})

	Convey("Given an invalid root path", t, func() {
		invalidRootPath := "/api/v1/nodes/{name}"
		rootPathItem := spec.PathItem{
			PathItemProps: spec.PathItemProps{
				Post: &spec.Operation{},
			},
		}
		Convey("When newSpecV2Resource method is called", func() {
			schemaDefinitions := map[string]spec.Schema{}
			_, err := newSpecV2Resource(invalidRootPath, spec.Schema{}, rootPathItem, spec.PathItem{}, schemaDefinitions)
			Convey("And the err returned should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Given an empty path", t, func() {
		path := ""
		Convey("When newSpecV2Resource method is called", func() {
			schemaDefinitions := map[string]spec.Schema{}
			_, err := newSpecV2Resource(path, spec.Schema{}, spec.PathItem{}, spec.PathItem{}, schemaDefinitions)
			Convey("And the err returned should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestBuildResourceName(t *testing.T) {

	testCases := []struct {
		path                 string
		expectedResourceName string
		expectedError        error
	}{
		{
			path:                 "/cdns",
			expectedResourceName: "cdns",
			expectedError:        nil,
		},
		{
			path:                 "/v1/cdns",
			expectedResourceName: "cdns_v1",
			expectedError:        nil,
		},
		{
			path:                 "/v1/cdns/",
			expectedResourceName: "cdns_v1",
			expectedError:        nil,
		},
		{
			path:                 "/cdns/{id}/firewalls",
			expectedResourceName: "cdns_firewalls",
			expectedError:        nil,
		},
		{
			path:                 "/v1/cdns/{id}/firewalls",
			expectedResourceName: "cdns_v1_firewalls",
			expectedError:        nil,
		},
		{
			path:                 "/cdns/{id}/v1/firewalls",
			expectedResourceName: "cdns_firewalls_v1",
			expectedError:        nil,
		},
		{
			path: "/v1/cdns/{id}/v2/firewalls",
			//expectedResourceName: "cdns_v1_cdn_v1_firewalls",
			expectedResourceName: "cdns_v1_firewalls_v2",
			expectedError:        nil,
		},
		{
			path:                 "/v1/cdns/{id}/v2/firewalls/{id}/v3/rules",
			expectedResourceName: "cdns_v1_firewalls_v2_rules_v3",
			expectedError:        nil,
		},
		{ // This is considered a wrongly structured path not following resful best practises for building subresource paths, however the plugin still supports it to not be so opinionated
			path:                 "/v1/cdns/{id}/firewalls/v3/rules",
			expectedResourceName: "cdns_v1_rules_v3",
			expectedError:        nil,
		},
	}

	for _, tc := range testCases {
		Convey("Given a SpecV2Resource with a sub-resource root path", t, func() {
			r := SpecV2Resource{
				Path: tc.path,
			}

			Convey("When buildResourceName is called", func() {
				resourceName, err := r.buildResourceName()
				if tc.expectedError != nil {
					Convey("Then the error returned should be the expected one", func() {
						So(err.Error(), ShouldEqual, tc.expectedError.Error())
					})
				}
				Convey("And the resource name should be the expected one", func() {
					So(resourceName, ShouldEqual, tc.expectedResourceName)
				})
			})
		})
	}

	// TODO: Add support for sub-resources names to honour preferred parent resource name as specified in with the x-terraform-resource-name in the parent path configuration.
	//  For instance, if /v1/cdns/{id}/v1/firewalls path had the x-terraform-resource-name defined and the value was "cdn_v1_firewall"
	// the expected returned name would be "cdn_v1_firewall_v1" (note the version should be automatically injected)
	Convey("Given a SpecV2Resource with a path and a preferred terraform resource name", t, func() {
		expectedResourceName := "user"
		r := SpecV2Resource{
			Path: "/users",
			RootPathItem: spec.PathItem{
				PathItemProps: spec.PathItemProps{
					Post: &spec.Operation{
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								extTfResourceName: expectedResourceName,
							},
						},
					},
				},
			},
		}
		Convey("When buildResourceName is called", func() {
			resourceName, err := r.buildResourceName()
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the resource name should be the expected one", func() {
				So(resourceName, ShouldEqual, expectedResourceName)
			})
		})
	})
}

func TestIsSubResource(t *testing.T) {
	Convey("Given a SpecV2Resource configured with a root path", t, func() {
		r := SpecV2Resource{
			Path: "/cdns",
		}
		Convey("When isSubResource is called", func() {
			isSubResource, parentResourceNames, fullParentResourceName := r.isSubResource()
			Convey("And the the bool returned should be false", func() {
				So(isSubResource, ShouldBeFalse)
			})
			Convey("And the parentResourceNames should be empty", func() {
				So(parentResourceNames, ShouldBeEmpty)
			})
			Convey("And the fullParentResourceName should be empty", func() {
				So(fullParentResourceName, ShouldBeEmpty)
			})
		})
	})
	Convey("Given a SpecV2Resource configured with a root path using versioning", t, func() {
		r := SpecV2Resource{
			Path: "/v1/cdns",
		}
		Convey("When isSubResource is called", func() {
			isSubResource, parentResourceNames, fullParentResourceName := r.isSubResource()
			Convey("And the the bool returned should be false", func() {
				So(isSubResource, ShouldBeFalse)
			})
			Convey("And the parentResourceNames should be empty", func() {
				So(parentResourceNames, ShouldBeEmpty)
			})
			Convey("And the fullParentResourceName should be empty", func() {
				So(fullParentResourceName, ShouldBeEmpty)
			})
		})
	})
	Convey("Given a SpecV2Resource configured with a path that is indeed a sub-resource (with parent using versioning)", t, func() {
		r := SpecV2Resource{
			Path: "/v1/cdns/{id}/firewalls",
		}
		Convey("When isSubResource is called", func() {
			isSubResource, parentResourceNames, fullParentResourceName := r.isSubResource()
			Convey("And the the bool returned should be true", func() {
				So(isSubResource, ShouldBeTrue)
			})
			Convey("And the parentResourceNames should not be empty and contain the right items", func() {
				So(len(parentResourceNames), ShouldEqual, 1)
				So(parentResourceNames[0], ShouldEqual, "cdns_v1")
			})
			Convey("And the fullParentResourceName should match the expected name", func() {
				So(fullParentResourceName, ShouldEqual, "cdns_v1")
			})
		})
	})
	Convey("Given a SpecV2Resource configured with a path that is indeed a sub-resource (with parent not using versioning)", t, func() {
		r := SpecV2Resource{
			Path: "/cdns/{id}/firewalls",
		}
		Convey("When isSubResource is called", func() {
			isSubResource, parentResourceNames, fullParentResourceName := r.isSubResource()
			Convey("And the the bool returned should be true", func() {
				So(isSubResource, ShouldBeTrue)
			})
			Convey("And the parentResourceNames should not be empty and contain the right items", func() {
				So(len(parentResourceNames), ShouldEqual, 1)
				So(parentResourceNames[0], ShouldEqual, "cdns")
			})
			Convey("And the fullParentResourceName should match the expected name", func() {
				So(fullParentResourceName, ShouldEqual, "cdns")
			})
		})
	})
	Convey("Given a SpecV2Resource configured with a path that is indeed a multiple level sub-resource", t, func() {
		r := SpecV2Resource{
			Path: "/cdns/{id}/firewalls/{id}/rules",
		}
		Convey("When isSubResource is called", func() {
			isSubResource, parentResourceNames, fullParentResourceName := r.isSubResource()
			Convey("And the the bool returned should be true", func() {
				So(isSubResource, ShouldBeTrue)
			})
			Convey("And the parentResourceNames should not be empty and contain the right items", func() {
				So(len(parentResourceNames), ShouldEqual, 2)
				So(parentResourceNames[0], ShouldEqual, "cdns")
				So(parentResourceNames[1], ShouldEqual, "firewalls")
			})
			Convey("And the fullParentResourceName should match the expected name", func() {
				So(fullParentResourceName, ShouldEqual, "cdns_firewalls")
			})
		})
	})
	Convey("Given a SpecV2Resource configured with a path that is indeed a multiple level sub-resource with versioning", t, func() {
		r := SpecV2Resource{
			Path: "/v1/cdns/{id}/v2/firewalls/{id}/v3/rules",
		}
		Convey("When isSubResource is called", func() {
			isSubResource, parentResourceNames, fullParentResourceName := r.isSubResource()
			Convey("And the the bool returned should be true", func() {
				So(isSubResource, ShouldBeTrue)
			})
			Convey("And the parentResourceNames should not be empty and contain the right items", func() {
				So(len(parentResourceNames), ShouldEqual, 2)
				So(parentResourceNames[0], ShouldEqual, "cdns_v1")
				So(parentResourceNames[1], ShouldEqual, "firewalls_v2")
			})
			Convey("And the fullParentResourceName should match the expected name", func() {
				So(fullParentResourceName, ShouldEqual, "cdns_v1_firewalls_v2")
			})
		})
	})
	Convey("Given a SpecV2Resource configured with a path that is a subresource but the path is wrongly structured not following best restful practises for building subresource paths (the 'firewalls' parent in the path is missing the id path param)", t, func() {
		r := SpecV2Resource{
			Path: "/v1/cdns/{id}/v2/firewalls/v3/rules",
		}
		Convey("When isSubResource is called", func() {
			isSubResource, parentResourceNames, fullParentResourceName := r.isSubResource()
			Convey("And then the resource should be considered a subresource and the output should match the expected output values", func() {
				So(isSubResource, ShouldBeTrue)
				So(len(parentResourceNames), ShouldEqual, 1)
				So(parentResourceNames[0], ShouldEqual, "cdns_v1")
				So(fullParentResourceName, ShouldEqual, "cdns_v1")
			})
		})
	})
}

func assertSchemaProperty(actualSpecSchemaDefinition *specSchemaDefinition, expectedName string, expectedType schemaDefinitionPropertyType, expectedRequired, expectedReadOnly, expectedComputed bool) {
	prop, err := actualSpecSchemaDefinition.getProperty(expectedName)
	So(err, ShouldBeNil)
	fmt.Printf(">>> Validating '%s' property settings\n", prop.Name)
	So(prop.Type, ShouldEqual, expectedType)
	So(prop.Required, ShouldEqual, expectedRequired)
	So(prop.ReadOnly, ShouldEqual, expectedReadOnly)
	So(prop.Computed, ShouldEqual, expectedComputed)
}

func TestGetResourceSchema(t *testing.T) {
	Convey("Given a SpecV2Resource containing a root path", t, func() {
		r := &SpecV2Resource{
			Path: "/cdns",
			SchemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Required: []string{"number_required_prop"},
					Properties: map[string]spec.Schema{
						"string_readonly_prop": {
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"string"},
							},
							SwaggerSchemaProps: spec.SwaggerSchemaProps{
								ReadOnly: true,
							},
						},
						"int_optional_computed_prop": {
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"integer"},
							},
							VendorExtensible: spec.VendorExtensible{
								Extensions: spec.Extensions{
									extTfComputed: true,
								},
							},
						},
						"number_required_prop": {
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"number"},
							},
						},
						"bool_prop": {
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"boolean"},
							},
						},
					},
				},
			},
		}
		Convey("When getSchemaDefinition is called with a schema containing various properties", func() {
			specSchemaDefinition, err := r.getResourceSchema()
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the specSchemaDefinition returned should be configured as expected", func() {
				So(len(specSchemaDefinition.Properties), ShouldEqual, len(r.SchemaDefinition.SchemaProps.Properties))
				assertSchemaProperty(specSchemaDefinition, "string_readonly_prop", typeString, false, true, true)
				assertSchemaProperty(specSchemaDefinition, "int_optional_computed_prop", typeInt, false, false, true)
				assertSchemaProperty(specSchemaDefinition, "number_required_prop", typeFloat, true, false, false)
				assertSchemaProperty(specSchemaDefinition, "bool_prop", typeBool, false, false, false)
			})
		})
	})
}

func TestGetSchemaDefinition(t *testing.T) {

	Convey("Given a blank SpecV2Resource", t, func() {
		r := &SpecV2Resource{}
		Convey("When getSchemaDefinition is called with a nil arg", func() {
			_, err := r.getSchemaDefinition(nil)
			Convey("Then the error returned matches the expected one", func() {
				So(err.Error(), ShouldEqual, "schema argument must not be nil")
			})
		})
		Convey("When getSchemaDefinition is called passing a blank schema", func() {
			d, e := r.getSchemaDefinition(&spec.Schema{})
			Convey("Then the error returned is not nil", func() {
				So(e, ShouldBeNil)
			})
			Convey("And the schema definition contains empty Properties", func() {
				So(d, ShouldNotBeNil)
				So(d.Properties, ShouldBeEmpty)
			})
		})
		Convey("When getSchemaDefinition is called passing a schema with a weird property type", func() {
			schema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"string_readonly_prop": {
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"something weird"},
							},
						},
					},
				},
			}
			d, e := r.getSchemaDefinition(&schema)
			Convey("Then the error returned is not nil", func() {
				So(e, ShouldNotBeNil)
			})
			Convey("And the schema definition returned is nil", func() {
				So(d, ShouldBeNil)
			})
		})

	})
	Convey("Given a SpecV2Resource containing a root path", t, func() {
		r := &SpecV2Resource{
			Path: "/cdns",
		}
		Convey("When getSchemaDefinition is called with a schema containing various properties", func() {
			s := &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Required: []string{"number_required_prop"},
					Properties: map[string]spec.Schema{
						"string_readonly_prop": {
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"string"},
							},
							SwaggerSchemaProps: spec.SwaggerSchemaProps{
								ReadOnly: true,
							},
						},
						"int_optional_computed_prop": {
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"integer"},
							},
							VendorExtensible: spec.VendorExtensible{
								Extensions: spec.Extensions{
									extTfComputed: true,
								},
							},
						},
						"number_required_prop": {
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"number"},
							},
						},
						"bool_prop": {
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"boolean"},
							},
						},
					},
				},
			}
			specSchemaDefinition, err := r.getSchemaDefinition(s)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the specSchemaDefinition returned should be configured as expected", func() {
				So(len(specSchemaDefinition.Properties), ShouldEqual, len(s.SchemaProps.Properties))
				assertSchemaProperty(specSchemaDefinition, "string_readonly_prop", typeString, false, true, true)
				assertSchemaProperty(specSchemaDefinition, "int_optional_computed_prop", typeInt, false, false, true)
				assertSchemaProperty(specSchemaDefinition, "number_required_prop", typeFloat, true, false, false)
				assertSchemaProperty(specSchemaDefinition, "bool_prop", typeBool, false, false, false)
			})
		})
	})

	Convey("Given a SpecV2Resource containing a subresource path (one level)", t, func() {
		r := &SpecV2Resource{
			Path: "/v1/cdns/{id}/firewalls",
		}
		Convey("When getSchemaDefinition is called with a schema containing various properties", func() {
			s := &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Required: []string{"number_required_prop"},
					Properties: map[string]spec.Schema{
						"string_readonly_prop": {
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"string"},
							},
							SwaggerSchemaProps: spec.SwaggerSchemaProps{
								ReadOnly: true,
							},
						},
						"int_optional_computed_prop": {
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"integer"},
							},
							VendorExtensible: spec.VendorExtensible{
								Extensions: spec.Extensions{
									extTfComputed: true,
								},
							},
						},
						"number_required_prop": {
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"number"},
							},
						},
						"bool_prop": {
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"boolean"},
							},
						},
					},
				},
			}
			specSchemaDefinition, err := r.getSchemaDefinition(s)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the specSchemaDefinition returned should be configured with the expected number of properties including the parent id one", func() {
				So(len(specSchemaDefinition.Properties), ShouldEqual, len(s.SchemaProps.Properties)+1)
			})
			Convey("And the specSchemaDefinition returned should be configured as expected", func() {
				assertSchemaProperty(specSchemaDefinition, "string_readonly_prop", typeString, false, true, true)
				assertSchemaProperty(specSchemaDefinition, "int_optional_computed_prop", typeInt, false, false, true)
				assertSchemaProperty(specSchemaDefinition, "number_required_prop", typeFloat, true, false, false)
				assertSchemaProperty(specSchemaDefinition, "bool_prop", typeBool, false, false, false)
			})
			Convey("And the specSchemaDefinition returned should be configured with the parent id property too", func() {
				assertSchemaProperty(specSchemaDefinition, "cdns_v1_id", typeString, false, false, true)
			})
		})
	})

	Convey("Given a SpecV2Resource containing a subresource path (one level) that has a non resftul subresource path", t, func() {
		r := &SpecV2Resource{
			Path: "/v1/cdns/{id}/firewalls/v1/rules",
		}
		Convey("When getSchemaDefinition is called with a schema containing various properties", func() {
			s := &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Required: []string{"number_required_prop"},
					Properties: map[string]spec.Schema{
						"string_readonly_prop": {
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"string"},
							},
							SwaggerSchemaProps: spec.SwaggerSchemaProps{
								ReadOnly: true,
							},
						},
					},
				},
			}
			specSchemaDefinition, err := r.getSchemaDefinition(s)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the specSchemaDefinition returned should be configured with the expected number of properties including the parent id one", func() {
				So(len(specSchemaDefinition.Properties), ShouldEqual, len(s.SchemaProps.Properties)+1)
			})
			Convey("And the specSchemaDefinition returned should be configured as expected", func() {
				assertSchemaProperty(specSchemaDefinition, "string_readonly_prop", typeString, false, true, true)
			})
			Convey("And the specSchemaDefinition returned should be configured with the parent id property too", func() {
				assertSchemaProperty(specSchemaDefinition, "cdns_v1_id", typeString, false, false, true)
			})
		})
	})

	Convey("Given a SpecV2Resource containing a subresource path (two level)", t, func() {
		r := &SpecV2Resource{
			Path: "/v1/cdns/{cdn_id}/v2/firewalls/{fw_id}/rules",
		}
		Convey("When getSchemaDefinition is called with a schema containing various properties", func() {
			s := &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Required: []string{"number_required_prop"},
					Properties: map[string]spec.Schema{
						"string_readonly_prop": {
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"string"},
							},
							SwaggerSchemaProps: spec.SwaggerSchemaProps{
								ReadOnly: true,
							},
						},
						"int_optional_computed_prop": {
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"integer"},
							},
							VendorExtensible: spec.VendorExtensible{
								Extensions: spec.Extensions{
									extTfComputed: true,
								},
							},
						},
						"number_required_prop": {
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"number"},
							},
						},
						"bool_prop": {
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"boolean"},
							},
						},
					},
				},
			}
			specSchemaDefinition, err := r.getSchemaDefinition(s)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the specSchemaDefinition returned should be configured with the expected number of properties including the parent id ones", func() {
				So(len(specSchemaDefinition.Properties), ShouldEqual, len(s.SchemaProps.Properties)+2)
			})
			Convey("And the specSchemaDefinition returned should be configured as expected", func() {
				assertSchemaProperty(specSchemaDefinition, "string_readonly_prop", typeString, false, true, true)
				assertSchemaProperty(specSchemaDefinition, "int_optional_computed_prop", typeInt, false, false, true)
				assertSchemaProperty(specSchemaDefinition, "number_required_prop", typeFloat, true, false, false)
				assertSchemaProperty(specSchemaDefinition, "bool_prop", typeBool, false, false, false)
			})
			Convey("And the specSchemaDefinition returned should be configured with the parent id property too", func() {
				assertSchemaProperty(specSchemaDefinition, "cdns_v1_id", typeString, false, false, true)
				assertSchemaProperty(specSchemaDefinition, "firewalls_v2_id", typeString, false, false, true)
			})
		})
	})

	Convey("Given a SpecV2Resource containing a root path (no subresource)", t, func() {
		r := &SpecV2Resource{
			Path: "/foo",
		}
		Convey("When getSchemaDefinition is called with a schema containing various properties", func() {
			s := &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"string_readonly_prop": {
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"string"},
							},
							SwaggerSchemaProps: spec.SwaggerSchemaProps{
								ReadOnly: true,
							},
						},
						"int_optional_computed_prop": {
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"integer"},
							},
							VendorExtensible: spec.VendorExtensible{
								Extensions: spec.Extensions{
									extTfComputed: true,
								},
							},
						},
					},
				},
			}
			specSchemaDefinition, err := r.getSchemaDefinition(s)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the specSchemaDefinition returned should be configured as expected", func() {
				So(len(specSchemaDefinition.Properties), ShouldEqual, 2)
				stringProp, err := specSchemaDefinition.getProperty("string_readonly_prop")
				So(err, ShouldBeNil)
				So(stringProp.Type, ShouldEqual, typeString)
				So(stringProp.ReadOnly, ShouldBeTrue)
				intProp, err := specSchemaDefinition.getProperty("int_optional_computed_prop")
				So(err, ShouldBeNil)
				So(intProp.Name, ShouldEqual, "int_optional_computed_prop")
				So(intProp.Type, ShouldEqual, typeInt)
				So(intProp.Computed, ShouldBeTrue)
			})
		})
	})

	Convey("Given a SpecV2Resource that is a subresource (one level parent)", t, func() {
		r := &SpecV2Resource{
			Path: "/parent/{parent_id}/child",
		}
		Convey("When getSchemaDefinition is called with a schema containing various properties", func() {
			s := &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"string_readonly_prop": {
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"string"},
							},
							SwaggerSchemaProps: spec.SwaggerSchemaProps{
								ReadOnly: true,
							},
						},
					},
				},
			}
			specSchemaDefinition, err := r.getSchemaDefinition(s)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the specSchemaDefinition returned should be configured as expected", func() {
				So(len(specSchemaDefinition.Properties), ShouldEqual, 2)
				So(specSchemaDefinition.Properties[0].Name, ShouldEqual, "string_readonly_prop")
				So(specSchemaDefinition.Properties[0].Type, ShouldEqual, typeString)
				So(specSchemaDefinition.Properties[0].ReadOnly, ShouldBeTrue)
			})
			Convey("And the specSchemaDefinition returned should also include the parent property", func() {
				So(specSchemaDefinition.Properties[1].Name, ShouldEqual, "parent_id")
				So(specSchemaDefinition.Properties[1].Type, ShouldEqual, typeString)
				So(specSchemaDefinition.Properties[1].Computed, ShouldBeTrue)
				So(specSchemaDefinition.Properties[1].Required, ShouldBeFalse)
			})
		})
	})

	Convey("Given a SpecV2Resource that is a subresource (two level parent)", t, func() {
		r := &SpecV2Resource{
			Path: "/parent/{parent_id}/subparent/{subparent_id}/child",
		}
		Convey("When getSchemaDefinition is called with a schema containing various properties", func() {
			s := &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"string_readonly_prop": {
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"string"},
							},
							SwaggerSchemaProps: spec.SwaggerSchemaProps{
								ReadOnly: true,
							},
						},
					},
				},
			}
			specSchemaDefinition, err := r.getSchemaDefinition(s)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the specSchemaDefinition returned should be configured as expected", func() {
				So(len(specSchemaDefinition.Properties), ShouldEqual, 3)
				So(specSchemaDefinition.Properties[0].Name, ShouldEqual, "string_readonly_prop")
				So(specSchemaDefinition.Properties[0].Type, ShouldEqual, typeString)
				So(specSchemaDefinition.Properties[0].ReadOnly, ShouldBeTrue)
			})
			Convey("And the specSchemaDefinition returned should also include the parents properties", func() {
				So(specSchemaDefinition.Properties[1].Name, ShouldEqual, "parent_id")
				So(specSchemaDefinition.Properties[1].Type, ShouldEqual, typeString)
				So(specSchemaDefinition.Properties[1].Computed, ShouldBeTrue)
				So(specSchemaDefinition.Properties[1].Required, ShouldBeFalse)
				So(specSchemaDefinition.Properties[2].Name, ShouldEqual, "subparent_id")
				So(specSchemaDefinition.Properties[2].Type, ShouldEqual, typeString)
				So(specSchemaDefinition.Properties[2].Computed, ShouldBeTrue)
				So(specSchemaDefinition.Properties[2].Required, ShouldBeFalse)
			})
		})
	})

	// TODO: Handle case where parent resources have preferred resource names as specified in with the x-terraform-resource-name in the parent path configuration. Hence, the resulting autogenered parent id property should honor the preferred name
	//Convey("Given a SpecV2Resource containing a sub-resource path (one level) and the parent resource using a preferred resource name", t, func() {
	//	// Specifying here how the parent resource will look like when the preferred name has been speficied
	//	//parentResource := SpecV2Resource{
	//	//	RootPathItem: spec.PathItem{
	//	//		PathItemProps: spec.PathItemProps{
	//	//			Post: &spec.Operation{
	//	//				VendorExtensible: spec.VendorExtensible{
	//	//					Extensions: spec.Extensions{
	//	//						extTfResourceName: "cdn",
	//	//					},
	//	//				},
	//	//			},
	//	//		},
	//	//	},
	//	//}
	//	r := &SpecV2Resource{
	//		Path: "/v1/cdns/{id}/firewalls",
	//	}
	//	Convey("When getSchemaDefinition is called with a schema containing various properties", func() {
	//		s := &spec.Schema{
	//			SchemaProps: spec.SchemaProps{
	//				Required: []string{"number_required_prop"},
	//				Properties: map[string]spec.Schema{
	//					"string_readonly_prop": {
	//						SchemaProps: spec.SchemaProps{
	//							Type: spec.StringOrArray{"string"},
	//						},
	//						SwaggerSchemaProps: spec.SwaggerSchemaProps{
	//							ReadOnly: true,
	//						},
	//					},
	//				},
	//			},
	//		}
	//		specSchemaDefinition, err := r.getSchemaDefinition(s)
	//		Convey("Then the error returned should be nil", func() {
	//			So(err, ShouldBeNil)
	//		})
	//		Convey("And the specSchemaDefinition returned should be configured with the expected number of properties including the parent id one", func() {
	//			So(len(specSchemaDefinition.Properties), ShouldEqual, len(s.SchemaProps.Properties) + 1)
	//		})
	//		Convey("And the specSchemaDefinition returned should be configured as expected", func() {
	//			assertSchemaProperty(specSchemaDefinition, "string_readonly_prop", typeString, false,true, true)
	//		})
	//		Convey("And the specSchemaDefinition returned should be configured with the parent id property named using the preferred parent name configured in the parent resource", func() {
	//			assertSchemaProperty(specSchemaDefinition, "cdn_v1_id", typeString, false,false, true)
	//		})
	//	})
	//})
}

func Test_getParentPropertiesNames(t *testing.T) {
	Convey("Given a SpecV2Resource with no Path", t, func() {
		r := &SpecV2Resource{}
		Convey("When the method getParentPropertiesNames is called", func() {
			p := r.getParentPropertiesNames()
			Convey("Then array returned should be empty", func() {
				So(p, ShouldBeEmpty)
			})
		})
	})

	Convey("Given a SpecV2Resource with some Path that is not a subresource", t, func() {
		r := &SpecV2Resource{Path: "/foo"}
		Convey("When the method getParentPropertiesNames is called", func() {
			p := r.getParentPropertiesNames()
			Convey("Then array returned should be empty", func() {
				So(p, ShouldBeEmpty)
			})
		})
	})

	Convey("Given a SpecV2Resource with some Path for a versioned parent resource where the path begins with the version", t, func() {
		r := &SpecV2Resource{
			Path: "/v2/cdns/{id}/v1/firewalls",
		}

		Convey("When the method getParentPropertiesNames is called", func() {
			p := r.getParentPropertiesNames()
			Convey("And the array returned should contain the expected parent name", func() {
				So(len(p), ShouldEqual, 1)
				So(p[0], ShouldEqual, "cdns_v2_id")
			})
		})
	})

	Convey("Given a SpecV2Resource with some Path for an unversioned parent resource", t, func() {
		r := &SpecV2Resource{
			Path: "/cdns/{id}/v1/firewalls",
		}
		Convey("When the method getParentPropertiesNames is called", func() {
			p := r.getParentPropertiesNames()
			Convey("And the array returned should contain the expected parent name", func() {
				So(len(p), ShouldEqual, 1)
				So(p[0], ShouldEqual, "cdns_id")
			})
		})
	})

	Convey("Given a SpecV2Resource with some Path that contains two parents, the first one using version and the second one without", t, func() {
		r := &SpecV2Resource{
			Path: "/v1/cdns/{id}/firewalls/{id}/rules",
		}
		Convey("When the method getParentPropertiesNames is called", func() {
			p := r.getParentPropertiesNames()
			Convey("And the array returned should contain the expected parent names", func() {
				So(len(p), ShouldEqual, 2)
				So(p[0], ShouldEqual, "cdns_v1_id")
				So(p[1], ShouldEqual, "firewalls_id")
			})
		})
	})

	Convey("Given a SpecV2Resource with some Path that contains two parents both using versioning", t, func() {
		r := &SpecV2Resource{
			Path: "/v1/cdns/{id}/v2/firewalls/{id}/rules",
		}
		Convey("When the method getParentPropertiesNames is called", func() {
			p := r.getParentPropertiesNames()

			Convey("And the array returned should contain the expected parent names", func() {
				So(len(p), ShouldEqual, 2)
				So(p[0], ShouldEqual, "cdns_v1_id")
				So(p[1], ShouldEqual, "firewalls_v2_id")
			})
		})
	})

	Convey("Given a SpecV2Resource with some Path that contains two parents; the first one with no version and the second one with it", t, func() {
		r := &SpecV2Resource{
			Path: "/cdns/{id}/v1/firewalls/{id}/rules",
		}
		Convey("When the method getParentPropertiesNames is called", func() {
			p := r.getParentPropertiesNames()
			Convey("And the array returned should contain the expected parent names", func() {
				So(len(p), ShouldEqual, 2)
				So(p[0], ShouldEqual, "cdns_id")
				So(p[1], ShouldEqual, "firewalls_v1_id")
			})
		})
	})

	Convey("Given a SpecV2Resource with some Path that contains two parents that do not use versioning", t, func() {
		r := &SpecV2Resource{
			Path: "/cdns/{id}/firewalls/{id}/rules",
		}
		Convey("When the method getParentPropertiesNames is called", func() {
			p := r.getParentPropertiesNames()
			Convey("And the array returned should contain the expected parent names", func() {
				So(len(p), ShouldEqual, 2)
				So(p[0], ShouldEqual, "cdns_id")
				So(p[1], ShouldEqual, "firewalls_id")
			})
		})
	})
}

func TestGetResourcePath(t *testing.T) {

	Convey("Given a SpecV2Resource with path resource that is not parameterised (root resource)", t, func() {
		r := SpecV2Resource{
			Path: "/v1/cdns",
		}
		Convey("When getResourcePath is called with an empty list of IDs", func() {
			resourcePath, err := r.getResourcePath([]string{})
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the returned resource path should match the expected one", func() {
				So(resourcePath, ShouldEqual, "/v1/cdns")
			})
		})
		Convey("When getResourcePath is called with a nil list of IDs", func() {
			resourcePath, err := r.getResourcePath(nil)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the returned resource path should match the expected one", func() {
				So(resourcePath, ShouldEqual, "/v1/cdns")
			})
		})
	})

	Convey("Given a SpecV2Resource with path resource that is parameterised (one level sub-resource)", t, func() {
		r := SpecV2Resource{
			Path: "/v1/cdns/{cdn_id}/v1/firewalls",
		}
		Convey("When getResourcePath is called with a list of IDs", func() {
			ids := []string{"parentID"}
			resourcePath, err := r.getResourcePath(ids)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the returned resource path should match the expected one", func() {
				So(resourcePath, ShouldEqual, "/v1/cdns/parentID/v1/firewalls")
			})
		})
		Convey("When getResourcePath is called with an empty list of IDs", func() {
			_, err := r.getResourcePath([]string{})
			Convey("Then the error returned should not be nil", func() {
				So(err.Error(), ShouldEqual, "could not resolve sub-resource path correctly '/v1/cdns/{cdn_id}/v1/firewalls' ([[/{cdn_id}/ {cdn_id}]]) with the given ids - missing ids to resolve the path params properly: []")
			})
		})
		Convey("When getResourcePath is called with an nil list of IDs", func() {
			_, err := r.getResourcePath(nil)
			Convey("Then the error returned should not be nil", func() {
				So(err.Error(), ShouldEqual, "could not resolve sub-resource path correctly '/v1/cdns/{cdn_id}/v1/firewalls' ([[/{cdn_id}/ {cdn_id}]]) with the given ids - missing ids to resolve the path params properly: []")
			})
		})
		Convey("When getResourcePath is called with a list of IDs that is bigger than the parameterised params in the path", func() {
			_, err := r.getResourcePath([]string{"cdnID", "somethingThatDoesNotBelongHere"})
			Convey("Then the error returned should not be nil", func() {
				So(err.Error(), ShouldEqual, "could not resolve sub-resource path correctly '/v1/cdns/{cdn_id}/v1/firewalls' ([[/{cdn_id}/ {cdn_id}]]) with the given ids - more ids than path params: [cdnID somethingThatDoesNotBelongHere]")
			})
		})
	})

	Convey("Given a SpecV2Resource with path resource that is parameterised (few levels sub-resource)", t, func() {
		r := SpecV2Resource{
			Path: "/v1/cdns/{cdn_id}/v1/firewalls/{fw_id}/rules",
		}
		Convey("When getResourcePath is called with a list of IDs", func() {
			ids := []string{"cdnID", "fwID"}
			resourcePath, err := r.getResourcePath(ids)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the returned resource path should match the expected one", func() {
				So(resourcePath, ShouldEqual, "/v1/cdns/cdnID/v1/firewalls/fwID/rules")
			})
		})
	})
}

func TestCreateSchemaDefinitionProperty(t *testing.T) {
	Convey("Given a SpecV2Resource", t, func() {
		r := SpecV2Resource{}

		//////////////////
		// Type checks
		//////////////////

		Convey("When createSchemaDefinitionProperty is called with a propertyName and a propertySchema of type string that is not required", func() {
			propertyName := "propertyName"
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"string"},
				},
			}
			requiredProperties := []string{}
			schemaDefinitionProperty, err := r.createSchemaDefinitionProperty(propertyName, propertySchema, requiredProperties)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the schema definition property should be configured with the right name and type", func() {
				So(schemaDefinitionProperty.Name, ShouldEqual, propertyName)
				So(schemaDefinitionProperty.Type, ShouldEqual, typeString)
			})
			Convey("And the schema definition property should not be readOnly", func() {
				So(schemaDefinitionProperty.isReadOnly(), ShouldBeFalse)
			})
			Convey("And the schema definition property should not be computed", func() {
				So(schemaDefinitionProperty.isComputed(), ShouldBeFalse)
			})
			Convey("And the schema definition property should not be required", func() {
				So(schemaDefinitionProperty.isRequired(), ShouldBeFalse)
			})
		})
		Convey("When createSchemaDefinitionProperty is called with a propertyName, propertySchema of type integer that is not required", func() {
			propertyName := "propertyName"
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"integer"},
				},
			}
			requiredProperties := []string{}
			schemaDefinitionProperty, err := r.createSchemaDefinitionProperty(propertyName, propertySchema, requiredProperties)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the schema definition property should be configured with the right name and type", func() {
				So(schemaDefinitionProperty.Name, ShouldEqual, propertyName)
				So(schemaDefinitionProperty.Type, ShouldEqual, typeInt)
			})
			Convey("And the schema definition property should not be readOnly", func() {
				So(schemaDefinitionProperty.isReadOnly(), ShouldBeFalse)
			})
			Convey("And the schema definition property should not be computed", func() {
				So(schemaDefinitionProperty.isComputed(), ShouldBeFalse)
			})
			Convey("And the schema definition property should not be required", func() {
				So(schemaDefinitionProperty.isRequired(), ShouldBeFalse)
			})
		})
		Convey("When createSchemaDefinitionProperty is called with a propertyName, propertySchema of type number that is not required", func() {
			propertyName := "propertyName"
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"number"},
				},
			}
			requiredProperties := []string{}
			schemaDefinitionProperty, err := r.createSchemaDefinitionProperty(propertyName, propertySchema, requiredProperties)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the schema definition property should be configured with the right name and type", func() {
				So(schemaDefinitionProperty.Name, ShouldEqual, propertyName)
				So(schemaDefinitionProperty.Type, ShouldEqual, typeFloat)
			})
			Convey("And the schema definition property should not be readOnly", func() {
				So(schemaDefinitionProperty.isReadOnly(), ShouldBeFalse)
			})
			Convey("And the schema definition property should not be computed", func() {
				So(schemaDefinitionProperty.isComputed(), ShouldBeFalse)
			})
			Convey("And the schema definition property should not be required", func() {
				So(schemaDefinitionProperty.isRequired(), ShouldBeFalse)
			})
		})
		Convey("When createSchemaDefinitionProperty is called with a propertyName, propertySchema of type boolean that is not required", func() {
			propertyName := "propertyName"
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"boolean"},
				},
			}
			requiredProperties := []string{}
			schemaDefinitionProperty, err := r.createSchemaDefinitionProperty(propertyName, propertySchema, requiredProperties)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the schema definition property should be configured with the right name and type", func() {
				So(schemaDefinitionProperty.Name, ShouldEqual, propertyName)
				So(schemaDefinitionProperty.Type, ShouldEqual, typeBool)
			})
			Convey("And the schema definition property should not be readOnly", func() {
				So(schemaDefinitionProperty.isReadOnly(), ShouldBeFalse)
			})
			Convey("And the schema definition property should not be computed", func() {
				So(schemaDefinitionProperty.isComputed(), ShouldBeFalse)
			})
			Convey("And the schema definition property should not be required", func() {
				So(schemaDefinitionProperty.isRequired(), ShouldBeFalse)
			})
		})

		Convey("When createSchemaDefinitionProperty is called with a propertyName, propertySchema of unknown type", func() {
			propertyName := "propertyName"
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					//Type: spec.StringOrArray{"boolean"}, NO TYPE ASSIGNED
				},
			}
			requiredProperties := []string{}
			_, err := r.createSchemaDefinitionProperty(propertyName, propertySchema, requiredProperties)
			Convey("Then the error returned should NOT be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("And the error message should equal", func() {
				So(err.Error(), ShouldEqual, "non supported '[]' type")
			})
		})

		Convey("When createSchemaDefinitionProperty is called with a propertyName, propertySchema of type object with nested properties that is not required", func() {
			propertyName := "propertyName"
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"object"},
					Properties: map[string]spec.Schema{
						"objectProperty": spec.Schema{
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"string"},
							},
						},
					},
				},
			}
			requiredProperties := []string{}
			schemaDefinitionProperty, err := r.createSchemaDefinitionProperty(propertyName, propertySchema, requiredProperties)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the schema definition property should be configured with the right name and type", func() {
				So(schemaDefinitionProperty.Name, ShouldEqual, propertyName)
				So(schemaDefinitionProperty.Type, ShouldEqual, typeObject)
			})
			Convey("And the schema definition property should not be readOnly", func() {
				So(schemaDefinitionProperty.isReadOnly(), ShouldBeFalse)
			})
			Convey("And the schema definition property should not be computed", func() {
				So(schemaDefinitionProperty.isComputed(), ShouldBeFalse)
			})
			Convey("And the schema definition property should not be required", func() {
				So(schemaDefinitionProperty.isRequired(), ShouldBeFalse)
			})
		})

		Convey("When createSchemaDefinitionProperty is called with a propertyName, propertySchema of type object with NO nested properties nor a REF", func() {
			propertyName := "propertyName"
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"object"},
					// Missing object schema information
				},
			}
			requiredProperties := []string{}
			_, err := r.createSchemaDefinitionProperty(propertyName, propertySchema, requiredProperties)
			Convey("Then the error returned should NOT be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("And the error message should equal", func() {
				So(err.Error(), ShouldEqual, "failed to process object type property 'propertyName': object is missing the nested schema definition or the ref is poitning to a non existing schema definition")
			})
		})

		Convey("When createSchemaDefinitionProperty is called with a propertyName and non required propertySchema of type array with items of type string", func() {
			propertyName := "propertyName"
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"array"},
					Items: &spec.SchemaOrArray{
						Schema: &spec.Schema{
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"string"},
							},
						},
					},
				},
			}
			schemaDefinitionProperty, err := r.createSchemaDefinitionProperty(propertyName, propertySchema, []string{})
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the schema definition property should be configured with the right name and type", func() {
				So(schemaDefinitionProperty.Name, ShouldEqual, propertyName)
				So(schemaDefinitionProperty.Type, ShouldEqual, typeList)
				So(schemaDefinitionProperty.ArrayItemsType, ShouldEqual, typeString)
				So(schemaDefinitionProperty.SpecSchemaDefinition, ShouldBeNil)
			})
			Convey("And the schema definition property should not be readOnly", func() {
				So(schemaDefinitionProperty.isReadOnly(), ShouldBeFalse)
			})
			Convey("And the schema definition property should not be computed", func() {
				So(schemaDefinitionProperty.isComputed(), ShouldBeFalse)
			})
			Convey("And the schema definition property should not be required", func() {
				So(schemaDefinitionProperty.isRequired(), ShouldBeFalse)
			})
		})

		Convey("When createSchemaDefinitionProperty is called with a propertyName and propertySchema non required of type array with items of type object (nested)", func() {
			propertyName := "propertyName"
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"array"},
					Items: &spec.SchemaOrArray{
						Schema: &spec.Schema{
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"object"},
								Properties: map[string]spec.Schema{
									"prop1": spec.Schema{
										SchemaProps: spec.SchemaProps{
											Type: spec.StringOrArray{"string"},
										},
									},
									"prop2": spec.Schema{
										SchemaProps: spec.SchemaProps{
											Type: spec.StringOrArray{"integer"},
										},
									},
								},
							},
						},
					},
				},
			}
			schemaDefinitionProperty, err := r.createSchemaDefinitionProperty(propertyName, propertySchema, []string{})
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the schema definition property should be configured with the right name, list type amd items type object", func() {
				So(schemaDefinitionProperty.Name, ShouldEqual, propertyName)
				So(schemaDefinitionProperty.Type, ShouldEqual, typeList)
				So(schemaDefinitionProperty.ArrayItemsType, ShouldEqual, typeObject)
			})
			Convey("And schema definition should contain the schema of the array items", func() {
				So(schemaDefinitionProperty.SpecSchemaDefinition, ShouldNotBeNil)
				exists, _ := assertPropertyExists(schemaDefinitionProperty.SpecSchemaDefinition.Properties, "prop1")
				So(exists, ShouldBeTrue)
				exists, _ = assertPropertyExists(schemaDefinitionProperty.SpecSchemaDefinition.Properties, "prop2")
				So(exists, ShouldBeTrue)

			})
			Convey("And the schema definition property should not be readOnly", func() {
				So(schemaDefinitionProperty.isReadOnly(), ShouldBeFalse)
			})
			Convey("And the schema definition property should not be computed", func() {
				So(schemaDefinitionProperty.isComputed(), ShouldBeFalse)
			})
			Convey("And the schema definition property should not be required", func() {
				So(schemaDefinitionProperty.isRequired(), ShouldBeFalse)
			})
		})

		Convey("When createSchemaDefinitionProperty is called with a propertyName and propertySchema non required of type array with items of type object (external ref definition)", func() {
			r := SpecV2Resource{
				SchemaDefinitions: map[string]spec.Schema{
					"Listeners": {
						SchemaProps: spec.SchemaProps{
							Type: spec.StringOrArray{"object"},
							Properties: map[string]spec.Schema{
								"protocol": {
									SchemaProps: spec.SchemaProps{
										Type: spec.StringOrArray{"string"},
									},
								},
							},
						},
					},
				},
			}
			propertyName := "propertyName"
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"array"},
					Items: &spec.SchemaOrArray{
						Schema: &spec.Schema{
							SchemaProps: spec.SchemaProps{
								Ref: spec.Ref{Ref: jsonreference.MustCreateRef("#/definitions/Listeners")},
							},
						},
					},
				},
			}
			schemaDefinitionProperty, err := r.createSchemaDefinitionProperty(propertyName, propertySchema, []string{})
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the schema definition property should be configured with the right name, list type amd items type object", func() {
				So(schemaDefinitionProperty.Name, ShouldEqual, propertyName)
				So(schemaDefinitionProperty.Type, ShouldEqual, typeList)
				So(schemaDefinitionProperty.ArrayItemsType, ShouldEqual, typeObject)
			})
			Convey("And schema definition should contain the schema of the array items", func() {
				So(schemaDefinitionProperty.SpecSchemaDefinition, ShouldNotBeNil)
				So(schemaDefinitionProperty.SpecSchemaDefinition.Properties, ShouldNotBeEmpty)
				So(schemaDefinitionProperty.SpecSchemaDefinition.Properties[0].Name, ShouldEqual, "protocol")
				So(schemaDefinitionProperty.SpecSchemaDefinition.Properties[0].Type, ShouldEqual, typeString)
			})
			Convey("And the schema definition property should not be readOnly", func() {
				So(schemaDefinitionProperty.isReadOnly(), ShouldBeFalse)
			})
			Convey("And the schema definition property should not be computed", func() {
				So(schemaDefinitionProperty.isComputed(), ShouldBeFalse)
			})
			Convey("And the schema definition property should not be required", func() {
				So(schemaDefinitionProperty.isRequired(), ShouldBeFalse)
			})
		})

		Convey("When createSchemaDefinitionProperty is called with a property schema that is required", func() {
			propertyName := "propertyName"
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"string"},
				},
			}
			requiredProperties := []string{"propertyName"}
			schemaDefinitionProperty, err := r.createSchemaDefinitionProperty(propertyName, propertySchema, requiredProperties)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the schema definition property should be required", func() {
				So(schemaDefinitionProperty.Required, ShouldBeTrue)
			})
			Convey("And the schema definition property should not be readOnly", func() {
				So(schemaDefinitionProperty.isReadOnly(), ShouldBeFalse)
			})
			Convey("And the schema definition property should not be computed", func() {
				So(schemaDefinitionProperty.isComputed(), ShouldBeFalse)
			})
		})

		Convey("When createSchemaDefinitionProperty is called with a property schema that is required AND the property is also set as readOnly (nonesense)", func() {
			propertyName := "propertyName"
			propertySchema := spec.Schema{
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					ReadOnly: true,
				},
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"string"},
				},
			}
			requiredProperties := []string{"propertyName"}
			schemaDefinitionProperty, err := r.createSchemaDefinitionProperty(propertyName, propertySchema, requiredProperties)
			Convey("Then the error returned should NOT be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("Then the error returned should match the expected one", func() {
				So(err.Error(), ShouldEqual, "failed to process property 'propertyName': a required property cannot be readOnly too")
			})
			Convey("Then the result should be nil", func() {
				So(schemaDefinitionProperty, ShouldBeNil)
			})
		})

		Convey("When createSchemaDefinitionProperty is called with an optional property schema that is computed (readOnly)", func() {
			propertyName := "propertyName"
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"string"},
				},
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					ReadOnly: true,
				},
			}
			requiredProperties := []string{}
			schemaDefinitionProperty, err := r.createSchemaDefinitionProperty(propertyName, propertySchema, requiredProperties)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the schema definition property should not be required", func() {
				So(schemaDefinitionProperty.isRequired(), ShouldBeFalse)
			})
			Convey("And the schema definition property should be readOnly", func() {
				So(schemaDefinitionProperty.isReadOnly(), ShouldBeTrue)
			})
			Convey("And the schema definition property should be computed", func() {
				So(schemaDefinitionProperty.isComputed(), ShouldBeTrue)
			})
		})

		Convey("When createSchemaDefinitionProperty is called with an optional property schema that is computed (readOnly) AND has a default value (meaning the computed value is known at runtime)", func() {
			expectedDefaultValue := "someDefaultValue"
			propertyName := "propertyName"
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type:    spec.StringOrArray{"string"},
					Default: expectedDefaultValue,
				},
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					ReadOnly: true,
				},
			}
			requiredProperties := []string{}
			schemaDefinitionProperty, err := r.createSchemaDefinitionProperty(propertyName, propertySchema, requiredProperties)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the schema definition property should not be required", func() {
				So(schemaDefinitionProperty.isRequired(), ShouldBeFalse)
			})
			Convey("And the schema definition property should be readOnly", func() {
				So(schemaDefinitionProperty.isReadOnly(), ShouldBeTrue)
			})
			Convey("And the schema definition property should be computed", func() {
				So(schemaDefinitionProperty.isComputed(), ShouldBeTrue)
			})
			Convey("And the schema definition property should have the right default value", func() {
				So(schemaDefinitionProperty.Default, ShouldEqual, expectedDefaultValue)
			})
		})

		Convey("When createSchemaDefinitionProperty is called with an optional property schema", func() {
			propertyName := "propertyWithNestedObj"

			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"object"},
					Properties: map[string]spec.Schema{
						"nested_obj": spec.Schema{
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"object"},
								Properties: map[string]spec.Schema{
									"nested_prop": spec.Schema{
										SchemaProps: spec.SchemaProps{
											Type: spec.StringOrArray{"string"},
										},
									},
								},
							},
						},
					},
				}}

			requiredProperties := []string{}
			schemaDefinitionProperty, err := r.createSchemaDefinitionProperty(propertyName, propertySchema, requiredProperties)

			fmt.Println(schemaDefinitionProperty, err)

			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the schema definition property type should be an object", func() {
				So(schemaDefinitionProperty.Type, ShouldEqual, typeObject)
			})

			Convey("And the schema definition property specs should contain only 1 item of type object", func() {
				So(len(schemaDefinitionProperty.SpecSchemaDefinition.Properties), ShouldEqual, 1)
				So(schemaDefinitionProperty.SpecSchemaDefinition.Properties[0].Type, ShouldEqual, typeObject)
			})

			Convey("And the nested object's property is a string", func() {
				nestedSpecSchema := *(schemaDefinitionProperty.SpecSchemaDefinition.Properties)[0]
				So(nestedSpecSchema.SpecSchemaDefinition.Properties[0].Type, ShouldEqual, typeString)
			})

		})

		Convey("When createSchemaDefinitionProperty is called with an optional property schema that has a default value (this means the property is optional-computed, since the API is expected to honour the default value (known at runtime) if input is not provided by the client)", func() {
			expectedDefaultValue := "someDefaultValue"
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type:    spec.StringOrArray{"string"},
					Default: expectedDefaultValue,
				},
			}
			schemaDefinitionProperty, err := r.createSchemaDefinitionProperty("propertyName", propertySchema, []string{})
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the schema definition property should be optional", func() {
				So(schemaDefinitionProperty.isRequired(), ShouldBeFalse)
			})
			Convey("And the schema definition property should not be readOnly", func() {
				So(schemaDefinitionProperty.isReadOnly(), ShouldBeFalse)
			})
			Convey("And the schema definition property should be computed", func() {
				So(schemaDefinitionProperty.isComputed(), ShouldBeFalse)
			})
			Convey("And the schema definition property should have the right default value", func() {
				So(schemaDefinitionProperty.Default, ShouldEqual, expectedDefaultValue)
			})
		})

		/////////////////////
		// Extension checks
		/////////////////////

		Convey("When createSchemaDefinitionProperty is called with a property schema that has the 'x-terraform-field-name' extension", func() {
			expectedTerraformName := "property_terraform_name"
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"string"},
				},
				VendorExtensible: spec.VendorExtensible{
					Extensions: spec.Extensions{
						extTfFieldName: expectedTerraformName,
					},
				},
			}
			schemaDefinitionProperty, err := r.createSchemaDefinitionProperty("propertyName", propertySchema, []string{})
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the schema definition property should be configured with the right", func() {
				So(schemaDefinitionProperty.PreferredName, ShouldEqual, expectedTerraformName)
			})
			Convey("And the schema definition property should not be computed", func() {
				So(schemaDefinitionProperty.isComputed(), ShouldBeFalse)
			})
		})

		Convey("When createSchemaDefinitionProperty is called with a property schema that has the 'x-terraform-force-new' extension", func() {
			expectedForceNewValue := true
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"string"},
				},
				VendorExtensible: spec.VendorExtensible{
					Extensions: spec.Extensions{
						extTfForceNew: expectedForceNewValue,
					},
				},
			}
			schemaDefinitionProperty, err := r.createSchemaDefinitionProperty("propertyName", propertySchema, []string{})
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the schema definition property should be have force new enabled", func() {
				So(schemaDefinitionProperty.ForceNew, ShouldEqual, expectedForceNewValue)
			})
			Convey("And the schema definition property should not be computed", func() {
				So(schemaDefinitionProperty.isComputed(), ShouldBeFalse)
			})
		})

		Convey("When createSchemaDefinitionProperty is called with a property schema that has the 'x-terraform-sensitive' extension", func() {
			expectedSensitiveValue := true
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"string"},
				},
				VendorExtensible: spec.VendorExtensible{
					Extensions: spec.Extensions{
						extTfSensitive: expectedSensitiveValue,
					},
				},
			}
			schemaDefinitionProperty, err := r.createSchemaDefinitionProperty("propertyName", propertySchema, []string{})
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the schema definition property should be sensitive", func() {
				So(schemaDefinitionProperty.Sensitive, ShouldEqual, expectedSensitiveValue)
			})
			Convey("And the schema definition property should not be computed", func() {
				So(schemaDefinitionProperty.isComputed(), ShouldBeFalse)
			})
		})

		Convey("When createSchemaDefinitionProperty is called with a property schema that has the 'x-terraform-id' extension", func() {
			expectedIsIdentifierValue := true
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"string"},
				},
				VendorExtensible: spec.VendorExtensible{
					Extensions: spec.Extensions{
						extTfID: expectedIsIdentifierValue,
					},
				},
			}
			schemaDefinitionProperty, err := r.createSchemaDefinitionProperty("propertyName", propertySchema, []string{})
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the schema definition property should be marked as identifier", func() {
				So(schemaDefinitionProperty.IsIdentifier, ShouldEqual, expectedIsIdentifierValue)
			})
			Convey("And the schema definition property should not be computed", func() {
				So(schemaDefinitionProperty.isComputed(), ShouldBeFalse)
			})
		})

		Convey("When createSchemaDefinitionProperty is called with a property schema that has the 'x-terraform-immutable' extension", func() {
			expectedIsImmutableValue := true
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"string"},
				},
				VendorExtensible: spec.VendorExtensible{
					Extensions: spec.Extensions{
						extTfImmutable: expectedIsImmutableValue,
					},
				},
			}
			schemaDefinitionProperty, err := r.createSchemaDefinitionProperty("propertyName", propertySchema, []string{})
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the schema definition property should be immutable", func() {
				So(schemaDefinitionProperty.Immutable, ShouldEqual, expectedIsImmutableValue)
			})
			Convey("And the schema definition property should not be computed", func() {
				So(schemaDefinitionProperty.isComputed(), ShouldBeFalse)
			})
		})

		Convey("When createSchemaDefinitionProperty is called with a property schema that has the 'x-terraform-field-status' extension", func() {
			expectedIsStatusFieldValue := true
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"string"},
				},
				VendorExtensible: spec.VendorExtensible{
					Extensions: spec.Extensions{
						extTfFieldStatus: expectedIsStatusFieldValue,
					},
				},
			}
			schemaDefinitionProperty, err := r.createSchemaDefinitionProperty("propertyName", propertySchema, []string{})
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the schema definition property should be marked as the status field", func() {
				So(schemaDefinitionProperty.IsStatusIdentifier, ShouldEqual, expectedIsStatusFieldValue)
			})
			Convey("And the schema definition property should not be computed", func() {
				So(schemaDefinitionProperty.isComputed(), ShouldBeFalse)
			})
		})

		Convey(fmt.Sprintf("When createSchemaDefinitionProperty is called with an optional property schema that has the %s extension (this means the property is optional-computed, and the value computed is not known at runtime)", extTfComputed), func() {
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"string"},
				},
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					ReadOnly: false,
				},
				VendorExtensible: spec.VendorExtensible{
					Extensions: spec.Extensions{
						extTfComputed: true,
					},
				},
			}
			schemaDefinitionProperty, err := r.createSchemaDefinitionProperty("propertyName", propertySchema, []string{})
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the schema definition property should be optional", func() {
				So(schemaDefinitionProperty.isRequired(), ShouldBeFalse)
			})
			Convey("And the schema definition property should not be readOnly", func() {
				So(schemaDefinitionProperty.isReadOnly(), ShouldBeFalse)
			})
			Convey("And the schema definition property should be computed", func() {
				So(schemaDefinitionProperty.isComputed(), ShouldBeTrue)
			})
			Convey("And the schema definition property should have a nil default value", func() {
				So(schemaDefinitionProperty.Default, ShouldBeNil)
			})
		})

		Convey(fmt.Sprintf("When createSchemaDefinitionProperty is called with an optional property schema that violates one optional-computed validation (properties with default attributes cannot have the %s extension too, that does not make any sense)", extTfComputed), func() {
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type:    spec.StringOrArray{"string"},
					Default: "someDefaultValue",
				},
				VendorExtensible: spec.VendorExtensible{
					Extensions: spec.Extensions{
						extTfComputed: true,
					},
				},
			}
			schemaDefinitionProperty, err := r.createSchemaDefinitionProperty("propertyName", propertySchema, []string{})
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey(fmt.Sprintf("Then the error message returned should state that properties with the %s extension can not have a default value attached", extTfComputed), func() {
				So(err.Error(), ShouldEqual, "optional computed property validation failed for property 'propertyName': optional computed properties with default attributes should not have 'x-terraform-computed' extension too")
			})
			Convey("And the schema definition property returned should be nil", func() {
				So(schemaDefinitionProperty, ShouldBeNil)
			})
		})

		Convey(fmt.Sprintf("When createSchemaDefinitionProperty is called with an optional property schema that violates one optional-computed validation (properties with %s extension, should not be readOnly)", extTfComputed), func() {
			propertySchema := spec.Schema{
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					ReadOnly: true,
				},
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"string"},
				},
				VendorExtensible: spec.VendorExtensible{
					Extensions: spec.Extensions{
						extTfComputed: true,
					},
				},
			}
			schemaDefinitionProperty, err := r.createSchemaDefinitionProperty("propertyName", propertySchema, []string{})
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey(fmt.Sprintf("Then the error message returned should state that properties with the %s extension can not be readOnly", extTfComputed), func() {
				So(err.Error(), ShouldEqual, "optional computed property validation failed for property 'propertyName': optional computed properties marked with 'x-terraform-computed' can not be readOnly")
			})
			Convey("And the schema definition property returned should be nil", func() {
				So(schemaDefinitionProperty, ShouldBeNil)
			})
		})
	})
}

func TestIsOptionalComputedProperty(t *testing.T) {
	Convey("Given a SpecV2Resource", t, func() {
		r := &SpecV2Resource{}
		Convey("When isOptionalComputedProperty method is called with a property that is NOT optional", func() {
			property := spec.Schema{
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					ReadOnly: false,
				},
				SchemaProps: spec.SchemaProps{
					Default: "some_defaul_value",
				},
			}
			isOptionalComputedProperty, err := r.isOptionalComputedProperty("some_required_property_name", property, []string{"some_required_property_name"})
			Convey("The the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("AND the result returned should be false since the property is not optional", func() {
				So(isOptionalComputedProperty, ShouldBeFalse)
			})
		})
		Convey("When isOptionalComputedProperty method is called with a property that is optional, and matches the OptionalComputedWithDefault requirements (no computed and has a default value)", func() {
			property := spec.Schema{
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					ReadOnly: false,
				},
				SchemaProps: spec.SchemaProps{
					Default: "some_defaul_value",
				},
			}
			isOptionalComputedProperty, err := r.isOptionalComputedProperty("some_optional_property_name", property, []string{"some_required_property_name"})
			Convey("The the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("AND the result returned should be true since the property is optional computed ", func() {
				So(isOptionalComputedProperty, ShouldBeTrue)
			})
		})
		Convey(fmt.Sprintf("When isOptionalComputedProperty method is called with a property that is optional, and matches the isOptionalComputed requirements (no computed and has the %s extension)", extTfComputed), func() {
			property := spec.Schema{
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					ReadOnly: false,
				},
				VendorExtensible: spec.VendorExtensible{
					Extensions: spec.Extensions{
						extTfComputed: true,
					},
				},
			}
			isOptionalComputedProperty, err := r.isOptionalComputedProperty("some_optional_property_name", property, []string{"some_required_property_name"})
			Convey("The the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("AND the result returned should be true since the property is optional computed ", func() {
				So(isOptionalComputedProperty, ShouldBeTrue)
			})
		})
		Convey(fmt.Sprintf("When isOptionalComputedProperty method is called with a property that is optional, and DOES NOT pass the validation as far as isOptionalComputed requirements is concerned (properties with %s extension cannot be readOnly)", extTfComputed), func() {
			property := spec.Schema{
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					ReadOnly: true,
				},
				VendorExtensible: spec.VendorExtensible{
					Extensions: spec.Extensions{
						extTfComputed: true,
					},
				},
			}
			isOptionalComputedProperty, err := r.isOptionalComputedProperty("some_optional_property_name", property, []string{"some_required_property_name"})
			Convey("The the error returned should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("The the error message returned should not be the expected one", func() {
				So(err.Error(), ShouldEqual, "optional computed property validation failed for property 'some_optional_property_name': optional computed properties marked with 'x-terraform-computed' can not be readOnly")
			})
			Convey("AND the result returned should be false since the property is NOT optional computed ", func() {
				So(isOptionalComputedProperty, ShouldBeFalse)
			})
		})
		Convey("When isOptionalComputedProperty method is called with a property that not optional computed at all (e,g: property is just computed)", func() {
			property := spec.Schema{
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					ReadOnly: true,
				},
			}
			isOptionalComputedProperty, err := r.isOptionalComputedProperty("some_optional_property_name", property, []string{"some_required_property_name"})
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("AND the result returned should be true since the property is NOT optional computed ", func() {
				So(isOptionalComputedProperty, ShouldBeFalse)
			})
		})
	})
}

func TestIsOptionalComputedWithDefault(t *testing.T) {
	Convey("Given a SpecV2Resource", t, func() {
		r := &SpecV2Resource{}
		Convey("When isOptionalComputedWithDefault method is called with a property that is NOT readOnly and has a default attribute", func() {
			property := spec.Schema{
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					ReadOnly: false,
				},
				SchemaProps: spec.SchemaProps{
					Default: "some_defaul_value",
				},
			}
			isOptionalComputedWithDefault, err := r.isOptionalComputedWithDefault("propertyName", property)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("The the result returned should be true since the property matches the requirements to be an optional computed property", func() {
				So(isOptionalComputedWithDefault, ShouldBeTrue)
			})
		})
		Convey("When isOptionalComputedWithDefault method is called with a property that is readOnly and has a default attribute", func() {
			property := spec.Schema{
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					ReadOnly: true,
				},
				SchemaProps: spec.SchemaProps{
					Default: "some_defaul_value",
				},
			}
			isOptionalComputedWithDefault, err := r.isOptionalComputedWithDefault("propertyName", property)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("The the result returned should be false since the property DOES NOT match the requirements to be an optional computed property", func() {
				So(isOptionalComputedWithDefault, ShouldBeFalse)
			})
		})
		Convey("When isOptionalComputedWithDefault method is called with a property that is NOT readOnly and has NO default attribute", func() {
			property := spec.Schema{
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					ReadOnly: false,
				},
				SchemaProps: spec.SchemaProps{
					Default: nil,
				},
			}
			isOptionalComputedWithDefault, err := r.isOptionalComputedWithDefault("propertyName", property)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("The the result returned should be false since the property DOES NOT match the requirements to be an optional computed property", func() {
				So(isOptionalComputedWithDefault, ShouldBeFalse)
			})
		})
		Convey("When isOptionalComputedWithDefault method is called with a property that is just readOnly", func() {
			property := spec.Schema{
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					ReadOnly: true,
				},
				SchemaProps: spec.SchemaProps{
					Default: nil,
				},
			}
			isOptionalComputedWithDefault, err := r.isOptionalComputedWithDefault("propertyName", property)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("The the result returned should be false since the property DOES NOT match the requirements to be an optional computed property", func() {
				So(isOptionalComputedWithDefault, ShouldBeFalse)
			})
		})
		Convey("When isOptionalComputedWithDefault method is called with a property that does not pass the validation phase since it has a default value AND the extension, this is wrong documentation", func() {
			property := spec.Schema{
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					ReadOnly: false,
				},
				SchemaProps: spec.SchemaProps{
					Default: "some_value",
				},
				VendorExtensible: spec.VendorExtensible{
					Extensions: spec.Extensions{
						extTfComputed: true,
					},
				},
			}
			isOptionalComputedWithDefault, err := r.isOptionalComputedWithDefault("propertyName", property)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("Then the error message returned should be the expected one", func() {
				So(err.Error(), ShouldEqual, "optional computed property validation failed for property 'propertyName': optional computed properties with default attributes should not have 'x-terraform-computed' extension too")
			})
			Convey("The the result returned should be false since the property DOES NOT match the requirements to be an optional computed property", func() {
				So(isOptionalComputedWithDefault, ShouldBeFalse)
			})
		})
	})
}

func TestIsOptionalComputed(t *testing.T) {
	Convey("Given a SpecV2Resource", t, func() {
		r := &SpecV2Resource{}
		Convey(fmt.Sprintf("When isOptionalComputed method is called with a property that is NOT computed (readOnly) and has the extension %s with value true", extTfComputed), func() {
			property := spec.Schema{
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					ReadOnly: false,
				},
				VendorExtensible: spec.VendorExtensible{
					Extensions: spec.Extensions{
						extTfComputed: true,
					},
				},
			}
			isOptionalComputed, err := r.isOptionalComputed("propertyName", property)
			Convey("The the result returned should not return an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("The the result returned should be true since the property matches the requirements to be an optional computed property", func() {
				So(isOptionalComputed, ShouldBeTrue)
			})
		})
		Convey(fmt.Sprintf("When isOptionalComputed method is called with a property that is NOT computed (readOnly) and has the extension %s with value false", extTfComputed), func() {
			property := spec.Schema{
				VendorExtensible: spec.VendorExtensible{
					Extensions: spec.Extensions{
						extTfComputed: false,
					},
				},
			}
			isOptionalComputed, err := r.isOptionalComputed("propertyName", property)
			Convey("The the result returned should not return an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("The the result returned should be false since the property DOES NOT match the requirements to be an optional computed property", func() {
				So(isOptionalComputed, ShouldBeFalse)
			})
		})
		Convey(fmt.Sprintf("When isOptionalComputed method is called with a property that is computed (readOnly) and has the extension %s with value true", extTfComputed), func() {
			property := spec.Schema{
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					ReadOnly: true, // this specifies that the property is computed
				},
				VendorExtensible: spec.VendorExtensible{
					Extensions: spec.Extensions{
						extTfComputed: true,
					},
				},
			}
			isOptionalComputed, err := r.isOptionalComputed("propertyName", property)
			Convey(fmt.Sprintf("The the result returned should not be nil since properties with the %s extension cannot be computed,", extTfComputed), func() {
				So(err, ShouldNotBeNil)
			})
			Convey("And the error message should be the expected one", func() {
				So(err.Error(), ShouldEqual, "optional computed property validation failed for property 'propertyName': optional computed properties marked with 'x-terraform-computed' can not be readOnly")
			})
			Convey("The the result returned should be false since the property DOES NOT match the requirements to be an optional computed property", func() {
				So(isOptionalComputed, ShouldBeFalse)
			})
		})
		Convey(fmt.Sprintf("When isOptionalComputed method is called with a property that is optional, and DOES NOT pass the validation as far as isOptionalComputed requirements is concerned (properties with %s extension cannot have default value populated)", extTfComputed), func() {
			property := spec.Schema{
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					ReadOnly: false,
				},
				SchemaProps: spec.SchemaProps{
					Default: "some_defaul_value",
				},
				VendorExtensible: spec.VendorExtensible{
					Extensions: spec.Extensions{
						extTfComputed: true,
					},
				},
			}
			isOptionalComputedProperty, err := r.isOptionalComputed("some_optional_property_name", property)
			Convey("The the error returned should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("The the error message returned should not be the expected one", func() {
				So(err.Error(), ShouldEqual, "optional computed property validation failed for property 'some_optional_property_name': optional computed properties marked with 'x-terraform-computed' can not have the default value as the value is not known at plan time. If the value is known, then this extension should not be used, and rather the 'default' attribute should be populated")
			})
			Convey("AND the result returned should be false since the property is NOT optional computed ", func() {
				So(isOptionalComputedProperty, ShouldBeFalse)
			})
		})
		Convey(fmt.Sprintf("When isOptionalComputed method is called with a property that DOES NOT have the extension %s present", extTfComputed), func() {
			property := spec.Schema{
				SwaggerSchemaProps: spec.SwaggerSchemaProps{},
			}
			isOptionalComputed, err := r.isOptionalComputed("propertyName", property)
			Convey("The the result returned should not return an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("The the result returned should be false", func() {
				So(isOptionalComputed, ShouldBeFalse)
			})
		})
	})
}

func TestIsArrayItemPrimitiveType(t *testing.T) {
	Convey("Given a SpecV2Resource", t, func() {
		r := &SpecV2Resource{}
		Convey("When isArrayItemPrimitiveType method is called with a primitive type typeString", func() {
			isPrimitive := r.isArrayItemPrimitiveType(typeString)
			Convey("The the result returned should be true", func() {
				So(isPrimitive, ShouldBeTrue)
			})
		})
		Convey("When isArrayItemPrimitiveType method is called with a primitive type typeInt", func() {
			isPrimitive := r.isArrayItemPrimitiveType(typeInt)
			Convey("The the result returned should be true", func() {
				So(isPrimitive, ShouldBeTrue)
			})
		})
		Convey("When isArrayItemPrimitiveType method is called with a primitive type typeFloat", func() {
			isPrimitive := r.isArrayItemPrimitiveType(typeFloat)
			Convey("The the result returned should be true", func() {
				So(isPrimitive, ShouldBeTrue)
			})
		})
		Convey("When isArrayItemPrimitiveType method is called with a primitive type typeBool", func() {
			isPrimitive := r.isArrayItemPrimitiveType(typeBool)
			Convey("The the result returned should be true", func() {
				So(isPrimitive, ShouldBeTrue)
			})
		})
		Convey("When isArrayItemPrimitiveType method is called with a NON primitive type typeList", func() {
			isPrimitive := r.isArrayItemPrimitiveType(typeList)
			Convey("The the result returned should be false", func() {
				So(isPrimitive, ShouldBeFalse)
			})
		})
		Convey("When isArrayItemPrimitiveType method is called with a NON primitive type typeObject", func() {
			isPrimitive := r.isArrayItemPrimitiveType(typeObject)
			Convey("The the result returned should be false", func() {
				So(isPrimitive, ShouldBeFalse)
			})
		})
	})
}

func TestValidateArrayItems(t *testing.T) {
	Convey("Given a SpecV2Resource", t, func() {
		r := &SpecV2Resource{}
		Convey("When validateArrayItems method is called with a property that does not have items", func() {
			property := spec.Schema{}
			_, err := r.validateArrayItems(property)
			Convey("The error shoudl not be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("And the error message should be the expected", func() {
				So(err.Error(), ShouldEqual, "array property is missing items schema definition")
			})
		})
		Convey("When validateArrayItems method is called with a property that does have items but they lack the schema", func() {
			property := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Items: &spec.SchemaOrArray{
						// no schema :(
					},
				},
			}
			_, err := r.validateArrayItems(property)
			Convey("The error shoudl not be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("And the error message should be the expected", func() {
				So(err.Error(), ShouldEqual, "array property is missing items schema definition")
			})
		})
		Convey("When validateArrayItems method is called with a property that does have items and a schema BUT the items are of type array (this is not supported at the moment)", func() {
			property := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Items: &spec.SchemaOrArray{
						Schema: &spec.Schema{
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"array"},
							},
						},
					},
				},
			}
			_, err := r.validateArrayItems(property)
			Convey("The error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("And the error message should be the expected", func() {
				So(err.Error(), ShouldEqual, "array property can not have items of type 'array'")
			})
		})
		Convey("When validateArrayItems method is called with an array of unknown type items", func() {
			property := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Items: &spec.SchemaOrArray{
						Schema: &spec.Schema{
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"unknown"},
							},
						},
					},
				},
			}
			_, err := r.validateArrayItems(property)
			Convey("The error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("And the error message should be the expected", func() {
				So(err.Error(), ShouldEqual, "non supported '[unknown]' type")
			})
		})
		Convey("When validateArrayItems method is called with a valid array property that has items of type string", func() {
			property := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Items: &spec.SchemaOrArray{
						Schema: &spec.Schema{
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"string"},
							},
						},
					},
				},
			}
			itemsPropType, err := r.validateArrayItems(property)
			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the type of the items should match the expected string", func() {
				So(itemsPropType, ShouldEqual, typeString)
			})
		})
		Convey("When validateArrayItems method is called with a valid array property that has items of type object", func() {
			property := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Items: &spec.SchemaOrArray{
						Schema: &spec.Schema{
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"object"},
								Properties: map[string]spec.Schema{
									"prop1": spec.Schema{},
								},
							},
						},
					},
				},
			}
			itemsPropType, err := r.validateArrayItems(property)
			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the type of the items should match the expected object", func() {
				So(itemsPropType, ShouldEqual, typeObject)
			})
		})
	})
}

func TestGetPropertyType(t *testing.T) {
	Convey("Given a SpecV2Resource", t, func() {
		r := &SpecV2Resource{}
		Convey("When getPropertyType method is called with a property of type array", func() {
			property := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"array"},
				},
			}
			itemsPropType, err := r.getPropertyType(property)
			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the type of the items should match the expected array", func() {
				So(itemsPropType, ShouldEqual, typeList)
			})
		})

		Convey("When getPropertyType method is called with a property of type object", func() {
			property := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"object"},
					Properties: map[string]spec.Schema{
						"prop1": spec.Schema{},
					},
				},
			}
			itemsPropType, err := r.getPropertyType(property)
			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the type of the items should match the expected object", func() {
				So(itemsPropType, ShouldEqual, typeObject)
			})
		})

		Convey("When getPropertyType method is called with a property of type string", func() {
			property := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"string"},
				},
			}
			itemsPropType, err := r.getPropertyType(property)
			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the type of the items should match the expected string", func() {
				So(itemsPropType, ShouldEqual, typeString)
			})
		})

		Convey("When getPropertyType method is called with a property of type integer", func() {
			property := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"integer"},
				},
			}
			itemsPropType, err := r.getPropertyType(property)
			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the type of the items should match the expected integer", func() {
				So(itemsPropType, ShouldEqual, typeInt)
			})
		})

		Convey("When getPropertyType method is called with a property of type float", func() {
			property := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"number"},
				},
			}
			itemsPropType, err := r.getPropertyType(property)
			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the type of the items should match the expected float", func() {
				So(itemsPropType, ShouldEqual, typeFloat)
			})
		})

		Convey("When getPropertyType method is called with a property of type bool", func() {
			property := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"boolean"},
				},
			}
			itemsPropType, err := r.getPropertyType(property)
			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the type of the items should match the expected bool", func() {
				So(itemsPropType, ShouldEqual, typeBool)
			})
		})

		Convey("When getPropertyType method is called with a property of type non supported", func() {
			property := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"non supported"},
				},
			}
			_, err := r.getPropertyType(property)
			Convey("The error should NOT be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("And the error returned should be ", func() {
				So(err.Error(), ShouldEqual, "non supported '[non supported]' type")
			})
		})
	})
}

func TestResourceIsObjectProperty(t *testing.T) {
	Convey("Given a SpecV2Resource", t, func() {
		r := &SpecV2Resource{}
		Convey("When isObjectProperty method is called with a property of type object that has nested properties", func() {
			property := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"object"},
					Properties: map[string]spec.Schema{
						"prop1": spec.Schema{},
					},
				},
			}
			isObject, objectSchema, err := r.isObjectProperty(property)
			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the result your be true", func() {
				So(isObject, ShouldBeTrue)
			})
			Convey("And the object schema should not be nil", func() {
				So(objectSchema, ShouldNotBeNil)
			})
		})
		Convey("When isObjectProperty method is called with a property of type object that has a ref to an external schema but is missing the type", func() {
			r := SpecV2Resource{
				SchemaDefinitions: map[string]spec.Schema{
					"Listeners": {
						SchemaProps: spec.SchemaProps{
							Type: spec.StringOrArray{"object"},
							Properties: map[string]spec.Schema{
								"protocol": {
									SchemaProps: spec.SchemaProps{
										Type: spec.StringOrArray{"string"},
									},
								},
							},
						},
					},
				},
			}
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					//Type: spec.StringOrArray{"object"}, // Missing type info but still should be considered as object
					Ref: spec.Ref{Ref: jsonreference.MustCreateRef("#/definitions/Listeners")},
				},
			}
			isObject, objectSchema, err := r.isObjectProperty(propertySchema)
			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the result your be true", func() {
				So(isObject, ShouldBeTrue)
			})
			Convey("And the object schema should not be nil", func() {
				So(objectSchema, ShouldNotBeNil)
			})
		})
		Convey("When isObjectProperty method is called with a property of type object that has a ref to an external schema and also has the type", func() {
			r := SpecV2Resource{
				SchemaDefinitions: map[string]spec.Schema{
					"Listeners": {
						SchemaProps: spec.SchemaProps{
							Type: spec.StringOrArray{"object"},
							Properties: map[string]spec.Schema{
								"protocol": {
									SchemaProps: spec.SchemaProps{
										Type: spec.StringOrArray{"string"},
									},
								},
							},
						},
					},
				},
			}
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"object"},
					Ref:  spec.Ref{Ref: jsonreference.MustCreateRef("#/definitions/Listeners")},
				},
			}
			isObject, objectSchema, err := r.isObjectProperty(propertySchema)
			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the result your be true", func() {
				So(isObject, ShouldBeTrue)
			})
			Convey("And the object schema should not be nil", func() {
				So(objectSchema, ShouldNotBeNil)
			})
		})

		Convey("When isObjectProperty method is called with a property of type object that has a ref to a non existing schema", func() {
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"object"},
					Ref:  spec.Ref{Ref: jsonreference.MustCreateRef("#/definitions/nonExisting")},
				},
			}
			isObject, objectSchema, err := r.isObjectProperty(propertySchema)
			Convey("The error should NOT be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("The error message should be the expected", func() {
				So(err.Error(), ShouldEqual, "object ref is poitning to a non existing schema definition: missing schema definition in the swagger file with the supplied ref '#/definitions/nonExisting'")
			})
			Convey("And the result your be true", func() {
				So(isObject, ShouldBeTrue)
			})
			Convey("And the object schema should be nil", func() {
				So(objectSchema, ShouldBeNil)
			})
		})

		Convey("When isObjectProperty method is called with a property that has nested schema with no properties", func() {
			property := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Items: &spec.SchemaOrArray{
						Schema: &spec.Schema{},
					},
				},
			}
			isObject, objectSchema, err := r.isObjectProperty(property)
			Convey("The error should NOT be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the result your be false", func() {
				So(isObject, ShouldBeFalse)
			})
			Convey("And the object schema should be nil", func() {
				So(objectSchema, ShouldBeNil)
			})
		})

		Convey("When isObjectProperty method is called with a property of type string", func() {
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"string"},
				},
			}
			isObject, objectSchema, err := r.isObjectProperty(propertySchema)
			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the result your be false", func() {
				So(isObject, ShouldBeFalse)
			})
			Convey("And the object schema should be nil", func() {
				So(objectSchema, ShouldBeNil)
			})
		})

	})
}

func TestResourceIsArrayProperty(t *testing.T) {
	Convey("Given a SpecV2Resource", t, func() {
		r := &SpecV2Resource{}
		Convey("When isArrayProperty is called with an array type property that has items of type object (nested)", func() {
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"array"},
					Items: &spec.SchemaOrArray{
						Schema: &spec.Schema{
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"object"},
								Properties: map[string]spec.Schema{
									"prop1": spec.Schema{
										SchemaProps: spec.SchemaProps{
											Type: spec.StringOrArray{"string"},
										},
									},
									"prop2": spec.Schema{
										SchemaProps: spec.SchemaProps{
											Type: spec.StringOrArray{"integer"},
										},
									},
								},
							},
						},
					},
				},
			}
			isArray, arrayItemType, objectItemSchema, err := r.isArrayProperty(propertySchema)
			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the result your be true", func() {
				So(isArray, ShouldBeTrue)
			})
			Convey("And the array items should be of type object", func() {
				So(arrayItemType, ShouldEqual, typeObject)
			})
			Convey("And the object schema should not be nil", func() {
				So(objectItemSchema, ShouldNotBeNil)
				exists, _ := assertPropertyExists(objectItemSchema.Properties, "prop1")
				So(exists, ShouldBeTrue)
				exists, _ = assertPropertyExists(objectItemSchema.Properties, "prop2")
				So(exists, ShouldBeTrue)
			})
		})
		Convey("When isArrayProperty is called with an array type property that has items of type primitive (e,g: string)", func() {
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"array"},
					Items: &spec.SchemaOrArray{
						Schema: &spec.Schema{
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"string"},
							},
						},
					},
				},
			}
			isArray, arrayItemType, objectItemSchema, err := r.isArrayProperty(propertySchema)
			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the result your be true", func() {
				So(isArray, ShouldBeTrue)
			})
			Convey("And the array items should be of type string", func() {
				So(arrayItemType, ShouldEqual, typeString)
			})
			Convey("And the object schema should be nil", func() {
				So(objectItemSchema, ShouldBeNil)
			})
		})
		Convey("When isArrayProperty is called with an array type property that has items of type object (ref)", func() {
			r := SpecV2Resource{
				SchemaDefinitions: map[string]spec.Schema{
					"Listeners": {
						SchemaProps: spec.SchemaProps{
							Type: spec.StringOrArray{"object"},
							Properties: map[string]spec.Schema{
								"protocol": {
									SchemaProps: spec.SchemaProps{
										Type: spec.StringOrArray{"string"},
									},
								},
							},
						},
					},
				},
			}
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"array"},
					Items: &spec.SchemaOrArray{
						Schema: &spec.Schema{
							SchemaProps: spec.SchemaProps{
								Ref: spec.Ref{Ref: jsonreference.MustCreateRef("#/definitions/Listeners")},
							},
						},
					},
				},
			}
			isArray, arrayItemType, objectItemSchema, err := r.isArrayProperty(propertySchema)
			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the result your be true", func() {
				So(isArray, ShouldBeTrue)
			})
			Convey("And the array items should be of type object", func() {
				So(arrayItemType, ShouldEqual, typeObject)
			})
			Convey("And the object schema should not be nil", func() {
				So(objectItemSchema, ShouldNotBeNil)
				exists, _ := assertPropertyExists(objectItemSchema.Properties, "protocol")
				So(exists, ShouldBeTrue)
			})
		})
		Convey("When isArrayProperty is called with a NON array property type", func() {
			propertySchema := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"string"},
				},
			}
			isArray, _, objectItemSchema, err := r.isArrayProperty(propertySchema)
			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the result your be false", func() {
				So(isArray, ShouldBeFalse)
			})
			Convey("And the object schema should be nil", func() {
				So(objectItemSchema, ShouldBeNil)
			})
		})
	})
}

func TestIsObjectTypeProperty(t *testing.T) {
	Convey("Given a SpecV2Resource", t, func() {
		r := &SpecV2Resource{}
		Convey("When isObjectTypeProperty method is called a property of type object", func() {
			property := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"object"},
				},
			}
			isArrayType := r.isObjectTypeProperty(property)
			Convey("The the result returned should be true", func() {
				So(isArrayType, ShouldBeTrue)
			})
		})
		Convey("When isObjectTypeProperty method is called a property that IS NOT of type object", func() {
			property := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"array"},
				},
			}
			isArrayType := r.isObjectTypeProperty(property)
			Convey("The the result returned should be false", func() {
				So(isArrayType, ShouldBeFalse)
			})
		})
	})
}

func TestIsArrayTypeProperty(t *testing.T) {
	Convey("Given a SpecV2Resource", t, func() {
		r := &SpecV2Resource{}
		Convey("When isArrayTypeProperty method is called a property of type array", func() {
			property := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"array"},
				},
			}
			isArrayType := r.isArrayTypeProperty(property)
			Convey("The the result returned should be true", func() {
				So(isArrayType, ShouldBeTrue)
			})
		})
		Convey("When isArrayTypeProperty method is called a property that IS NOT of type array", func() {
			property := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"object"},
				},
			}
			isArrayType := r.isArrayTypeProperty(property)
			Convey("The the result returned should be false", func() {
				So(isArrayType, ShouldBeFalse)
			})
		})
	})
}

func TestIsOfType(t *testing.T) {
	Convey("Given a SpecV2Resource", t, func() {
		r := &SpecV2Resource{}
		Convey("When isOfType method is called a property of the expected type", func() {
			property := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"string"},
				},
			}
			isString := r.isOfType(property, "string")
			Convey("The the result returned should be true", func() {
				So(isString, ShouldBeTrue)
			})
		})
		Convey("When isArrayTypeProperty method is called a property that IS NOT of the expected type", func() {
			property := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"string"},
				},
			}
			isInteger := r.isOfType(property, "integer")
			Convey("The the result returned should be false", func() {
				So(isInteger, ShouldBeFalse)
			})
		})
	})
}

func TestSwaggerPropIsRequired(t *testing.T) {
	Convey("Given a SpecV2Resource", t, func() {
		r := &SpecV2Resource{}
		Convey("When isRequired is called with a required prop", func() {
			requiredProp := "requiredProp"
			requiredProps := []string{requiredProp}
			isRequired := r.isRequired(requiredProp, requiredProps)
			Convey("The the result returned should be true", func() {
				So(isRequired, ShouldBeTrue)
			})
		})
		Convey("When isRequired is called with a NON required prop", func() {
			requiredProps := []string{"requiredProp"}
			isRequired := r.isRequired("nonRequired", requiredProps)
			Convey("The the result returned should be true", func() {
				So(isRequired, ShouldBeFalse)
			})
		})
	})
}

func TestGetResourceTerraformName(t *testing.T) {
	Convey("Given a SpecV2Resource with a root path item containing a post operation with the extension 'x-terraform-resource-name'", t, func() {
		extensions := spec.Extensions{}
		expectedResourceName := "example"
		extensions.Add(extTfResourceName, expectedResourceName)
		r := SpecV2Resource{
			RootPathItem: spec.PathItem{
				PathItemProps: spec.PathItemProps{
					Post: &spec.Operation{
						VendorExtensible: spec.VendorExtensible{
							Extensions: extensions,
						},
					},
				},
			},
		}
		Convey("When getResourceTerraformName method is called an existing extension", func() {
			value := r.getResourceTerraformName()
			Convey("Then the value returned should match the value in the extension", func() {
				So(value, ShouldEqual, expectedResourceName)
			})
		})
	})
	Convey("Given a SpecV2Resource with a root path item containing a post operation with the extension 'x-terraform-resource-name'", t, func() {
		r := SpecV2Resource{}
		Convey("When getResourceTerraformName method is called", func() {
			value := r.getResourceTerraformName()
			Convey("Then the value returned should be empty since the resource does not have such extension defined in the post operations", func() {
				So(value, ShouldEqual, "")
			})
		})
	})
}

func TestGetExtensionStringValue(t *testing.T) {
	Convey("Given a SpecV2Resource", t, func() {
		r := SpecV2Resource{}
		Convey("When getExtensionStringValue method is called an existing extension", func() {
			extensions := spec.Extensions{}
			expectedExtensionValue := "example"
			extensions.Add(extTfResourceName, expectedExtensionValue)
			value := r.getExtensionStringValue(extensions, extTfResourceName)
			Convey("Then the value returned should match the value in the extension", func() {
				So(value, ShouldEqual, expectedExtensionValue)
			})
		})
		Convey("When getExtensionStringValue method is called a NON existing extension", func() {
			extensions := spec.Extensions{}
			extensions.Add(extTfResourceName, "example")
			value := r.getExtensionStringValue(extensions, "somethingOtherExtensionName")
			Convey("Then the value returned should match the value in the extension", func() {
				So(value, ShouldEqual, "")
			})
		})
	})
}

func TestCreateResponses(t *testing.T) {
	Convey("Given a SpecV2Resource", t, func() {
		r := SpecV2Resource{}
		Convey("When createResponses method is called with an operation that has the 'x-terraform-resource-poll-enabled' extension set to true", func() {
			expectedTarget := "deployed"
			expectedStatus := "deploy_pending"
			extensions := spec.Extensions{}
			extensions.Add(extTfResourcePollEnabled, true)
			extensions.Add(extTfResourcePollTargetStatuses, expectedTarget)
			extensions.Add(extTfResourcePollPendingStatuses, expectedStatus)
			operation := &spec.Operation{
				OperationProps: spec.OperationProps{
					Responses: &spec.Responses{
						ResponsesProps: spec.ResponsesProps{
							StatusCodeResponses: map[int]spec.Response{
								http.StatusAccepted: {
									VendorExtensible: spec.VendorExtensible{
										Extensions: extensions,
									},
								},
							},
						},
					},
				},
			}
			specResponses := r.createResponses(operation)
			Convey("Then the spec responses map returned should not be empty", func() {
				So(specResponses, ShouldNotBeEmpty)
			})
			Convey("Then there should be an existing key for response code 202", func() {
				So(specResponses, ShouldContainKey, http.StatusAccepted)
			})
			Convey("And the response should meet the configuration", func() {
				So(specResponses[http.StatusAccepted].isPollingEnabled, ShouldBeTrue)
				So(specResponses[http.StatusAccepted].pollTargetStatuses, ShouldContain, expectedTarget)
				So(specResponses[http.StatusAccepted].pollPendingStatuses, ShouldContain, expectedStatus)
			})
		})

		Convey("When createResponses method is called with an operation does not have any status responses", func() {
			operation := &spec.Operation{
				OperationProps: spec.OperationProps{
					Responses: &spec.Responses{
						ResponsesProps: spec.ResponsesProps{
							StatusCodeResponses: map[int]spec.Response{},
						},
					},
				},
			}
			specResponses := r.createResponses(operation)
			Convey("Then the spec responses map returned should not be empty", func() {
				So(specResponses, ShouldBeEmpty)
			})
		})
	})
}

func TestIsResourcePollingEnabled(t *testing.T) {
	Convey("Given a SpecV2Resource", t, func() {
		r := SpecV2Resource{}
		Convey("When isResourcePollingEnabled method is called with a list of responses where one of the reponses matches the response status received and has the 'x-terraform-resource-poll-enabled' extension set to true", func() {
			extensions := spec.Extensions{}
			extensions.Add(extTfResourcePollEnabled, true)
			responses := &spec.Responses{
				ResponsesProps: spec.ResponsesProps{
					StatusCodeResponses: map[int]spec.Response{
						http.StatusAccepted: {
							VendorExtensible: spec.VendorExtensible{
								Extensions: extensions,
							},
						},
					},
				},
			}
			isResourcePollingEnabled := r.isResourcePollingEnabled(responses.StatusCodeResponses[http.StatusAccepted])
			Convey("Then the bool returned should be true", func() {
				So(isResourcePollingEnabled, ShouldBeTrue)
			})
		})
		Convey("When isResourcePollingEnabled method is called with a list of responses where one of the reponses matches the response status received and has the 'x-terraform-resource-poll-enabled' extension set to false", func() {
			extensions := spec.Extensions{}
			extensions.Add(extTfResourcePollEnabled, false)
			responses := &spec.Responses{
				ResponsesProps: spec.ResponsesProps{
					StatusCodeResponses: map[int]spec.Response{
						http.StatusAccepted: {
							VendorExtensible: spec.VendorExtensible{
								Extensions: extensions,
							},
						},
					},
				},
			}
			isResourcePollingEnabled := r.isResourcePollingEnabled(responses.StatusCodeResponses[http.StatusAccepted])
			Convey("Then the bool returned should be false", func() {
				So(isResourcePollingEnabled, ShouldBeFalse)
			})
		})
		Convey("When isResourcePollingEnabled method is called with list of responses where non of the codes match the given response http code", func() {
			responses := &spec.Responses{
				ResponsesProps: spec.ResponsesProps{
					StatusCodeResponses: map[int]spec.Response{
						http.StatusOK: {},
					},
				},
			}
			isResourcePollingEnabled := r.isResourcePollingEnabled(responses.StatusCodeResponses[http.StatusOK])
			Convey("Then bool returned should be false", func() {
				So(isResourcePollingEnabled, ShouldBeFalse)
			})
		})
	})
}

func TestGetResourcePollTargetStatuses(t *testing.T) {
	Convey("Given a SpecV2Resource", t, func() {
		r := SpecV2Resource{}
		Convey("When getResourcePollTargetStatuses method is called with a response that has a given extension 'x-terraform-resource-poll-completed-statuses'", func() {
			expectedTarget := "deployed"
			extensions := spec.Extensions{}
			extensions.Add(extTfResourcePollTargetStatuses, expectedTarget)
			responses := &spec.Responses{
				ResponsesProps: spec.ResponsesProps{
					StatusCodeResponses: map[int]spec.Response{
						http.StatusAccepted: {
							VendorExtensible: spec.VendorExtensible{
								Extensions: extensions,
							},
						},
					},
				},
			}
			statuses := r.getResourcePollTargetStatuses(responses.StatusCodeResponses[http.StatusAccepted])
			Convey("Then the status returned should contain", func() {
				So(statuses, ShouldContain, expectedTarget)
			})
		})
	})
}

func TestGetResourcePollPendingStatuses(t *testing.T) {
	Convey("Given a SpecV2Resource", t, func() {
		r := SpecV2Resource{}
		Convey("When getResourcePollPendingStatuses method is called with a response that has a given extension 'x-terraform-resource-poll-pending-statuses'", func() {
			expectedStatus := "deploy_pending"
			extensions := spec.Extensions{}
			extensions.Add(extTfResourcePollPendingStatuses, expectedStatus)
			responses := spec.Responses{
				ResponsesProps: spec.ResponsesProps{
					StatusCodeResponses: map[int]spec.Response{
						http.StatusAccepted: {
							VendorExtensible: spec.VendorExtensible{
								Extensions: extensions,
							},
						},
					},
				},
			}
			statuses := r.getResourcePollPendingStatuses(responses.StatusCodeResponses[http.StatusAccepted])
			Convey("Then the status returned should contain", func() {
				So(statuses, ShouldContain, expectedStatus)
			})
		})
	})
}

func TestGetPollingStatuses(t *testing.T) {
	Convey("Given a SpecV2Resource", t, func() {
		r := SpecV2Resource{}
		Convey("When getPollingStatuses method is called with a response that has a given extension 'x-terraform-resource-poll-completed-statuses'", func() {
			expectedTarget := "deployed"
			extensions := spec.Extensions{}
			extensions.Add(extTfResourcePollTargetStatuses, expectedTarget)
			responses := spec.Responses{
				ResponsesProps: spec.ResponsesProps{
					StatusCodeResponses: map[int]spec.Response{
						http.StatusAccepted: {
							VendorExtensible: spec.VendorExtensible{
								Extensions: extensions,
							},
						},
					},
				},
			}
			statuses := r.getPollingStatuses(responses.StatusCodeResponses[http.StatusAccepted], extTfResourcePollTargetStatuses)
			Convey("Then the statuses returned should contain", func() {
				So(statuses, ShouldContain, expectedTarget)
			})
		})

		Convey("When getPollingStatuses method is called with a response that has a given extension 'x-terraform-resource-poll-completed-statuses' containing multiple targets (comma separated with spaces)", func() {
			expectedTargets := "deployed, completed, done"
			extensions := spec.Extensions{}
			extensions.Add(extTfResourcePollTargetStatuses, expectedTargets)
			responses := spec.Responses{
				ResponsesProps: spec.ResponsesProps{
					StatusCodeResponses: map[int]spec.Response{
						http.StatusAccepted: {
							VendorExtensible: spec.VendorExtensible{
								Extensions: extensions,
							},
						},
					},
				},
			}
			statuses := r.getPollingStatuses(responses.StatusCodeResponses[http.StatusAccepted], extTfResourcePollTargetStatuses)
			Convey("Then the statuses returned should contain expected targets", func() {
				// the expected Targets are a list of targets with no spaces whatsoever, hence why the removal of spaces
				for _, expectedTarget := range strings.Split(strings.Replace(expectedTargets, " ", "", -1), ",") {
					So(statuses, ShouldContain, expectedTarget)
				}
			})
		})

		Convey("When getPollingStatuses method is called with a response that has a given extension 'x-terraform-resource-poll-completed-statuses' containing multiple targets (comma separated with no spaces)", func() {
			expectedTargets := "deployed,completed,done"
			extensions := spec.Extensions{}
			extensions.Add(extTfResourcePollTargetStatuses, expectedTargets)
			responses := spec.Responses{
				ResponsesProps: spec.ResponsesProps{
					StatusCodeResponses: map[int]spec.Response{
						http.StatusAccepted: {
							VendorExtensible: spec.VendorExtensible{
								Extensions: extensions,
							},
						},
					},
				},
			}
			statuses := r.getPollingStatuses(responses.StatusCodeResponses[http.StatusAccepted], extTfResourcePollTargetStatuses)
			Convey("Then the statuses returned should contain expected targets", func() {
				for _, expectedTarget := range strings.Split(expectedTargets, ",") {
					So(statuses, ShouldContain, expectedTarget)
				}
			})
		})

		Convey("When getPollingStatuses method is called with a response that has does not have a given extension 'x-terraform-resource-poll-completed-statuses'", func() {
			responses := spec.Responses{
				ResponsesProps: spec.ResponsesProps{
					StatusCodeResponses: map[int]spec.Response{
						http.StatusAccepted: {
							VendorExtensible: spec.VendorExtensible{
								Extensions: spec.Extensions{},
							},
						},
					},
				},
			}
			statuses := r.getPollingStatuses(responses.StatusCodeResponses[http.StatusAccepted], extTfResourcePollTargetStatuses)
			Convey("Then the status returned should be empty", func() {
				So(statuses, ShouldBeEmpty)
			})
		})
	})
}

func TestGetTimeouts(t *testing.T) {
	Convey("Given a SpecV2Resource", t, func() {
		expectedTimeout := "30s"
		extensions := spec.Extensions{}
		extensions.Add(extTfResourceTimeout, expectedTimeout)
		op := &spec.Operation{
			VendorExtensible: spec.VendorExtensible{
				Extensions: extensions,
			},
		}
		r := SpecV2Resource{
			RootPathItem: spec.PathItem{
				PathItemProps: spec.PathItemProps{
					Post: op,
				},
			},
			InstancePathItem: spec.PathItem{
				PathItemProps: spec.PathItemProps{
					Put:    op,
					Get:    op,
					Delete: op,
				},
			},
		}
		Convey("When getTimeouts method is called ", func() {
			timeouts, err := r.getTimeouts()
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Then the duration returned should contain the expected duration from the operation", func() {
				So(*timeouts.Post, ShouldEqual, time.Duration(30*time.Second))
				So(*timeouts.Get, ShouldEqual, time.Duration(30*time.Second))
				So(*timeouts.Put, ShouldEqual, time.Duration(30*time.Second))
				So(*timeouts.Delete, ShouldEqual, time.Duration(30*time.Second))
			})
		})
	})
}

func TestGetResourceTimeout(t *testing.T) {
	Convey("Given a SpecV2Resource", t, func() {
		r := SpecV2Resource{}
		Convey(fmt.Sprintf("When getResourceTimeout method is called with an operation that has the extension '%s'", extTfResourceTimeout), func() {
			expectedTimeout := "30s"
			extensions := spec.Extensions{}
			extensions.Add(extTfResourceTimeout, expectedTimeout)
			post := &spec.Operation{
				VendorExtensible: spec.VendorExtensible{
					Extensions: extensions,
				},
			}
			duration, err := r.getResourceTimeout(post)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Then the duration returned should contain", func() {
				So(*duration, ShouldEqual, time.Duration(30*time.Second))
			})
		})
	})
}

func TestGetTimeDuration(t *testing.T) {
	Convey("Given a SpecV2Resource", t, func() {
		r := SpecV2Resource{}
		Convey(fmt.Sprintf("When getTimeDuration method is called with a list of extensions that contains the extension passed in '%s' with value in seconds", extTfResourceTimeout), func() {
			expectedTimeout := "30s"
			extensions := spec.Extensions{}
			extensions.Add(extTfResourceTimeout, expectedTimeout)
			duration, err := r.getTimeDuration(extensions, extTfResourceTimeout)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Then the duration returned should contain", func() {
				So(*duration, ShouldEqual, time.Duration(30*time.Second))
			})
		})
		Convey(fmt.Sprintf("When getTimeDuration method is called with a list of extensions that contains the extension passed in '%s' with value in minutes (using fractions)", extTfResourceTimeout), func() {
			expectedTimeout := "20.5m"
			extensions := spec.Extensions{}
			extensions.Add(extTfResourceTimeout, expectedTimeout)
			duration, err := r.getTimeDuration(extensions, extTfResourceTimeout)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Then the duration returned should contain", func() {
				So(*duration, ShouldEqual, time.Duration((20*time.Minute)+(30*time.Second)))
			})
		})
		Convey(fmt.Sprintf("When getTimeDuration method is called with a list of extensions that contains the extension passed in '%s' with value in hours", extTfResourceTimeout), func() {
			expectedTimeout := "1h"
			extensions := spec.Extensions{}
			extensions.Add(extTfResourceTimeout, expectedTimeout)
			duration, err := r.getTimeDuration(extensions, extTfResourceTimeout)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Then the duration returned should contain", func() {
				So(*duration, ShouldEqual, time.Duration(1*time.Hour))
			})
		})
		Convey(fmt.Sprintf("When getTimeDuration method is called with a list of extensions that DOES NOT contain the extension passed in '%s'", extTfResourceTimeout), func() {
			expectedTimeout := "30s"
			extensions := spec.Extensions{}
			extensions.Add(extTfResourceTimeout, expectedTimeout)
			duration, err := r.getTimeDuration(extensions, "nonExistingExtension")
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Then the duration returned should be nil", func() {
				So(duration, ShouldBeNil)
			})
		})
		Convey(fmt.Sprintf("When getTimeDuration method is called with a list of extensions that DOES contain the extension passed in '%s' BUT the value is an empty string", extTfResourceTimeout), func() {
			extensions := spec.Extensions{}
			extensions.Add(extTfResourceTimeout, "")
			duration, err := r.getTimeDuration(extensions, extTfResourceTimeout)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("Then the duration returned should be nil", func() {
				So(duration, ShouldBeNil)
			})
			Convey("And the error message should be", func() {
				So(err.Error(), ShouldContainSubstring, "invalid duration value: ''. The value must be a sequence of decimal numbers each with optional fraction and a unit suffix (negative durations are not allowed). The value must be formatted either in seconds (s), minutes (m) or hours (h)")
			})
		})
		Convey(fmt.Sprintf("When getTimeDuration method is called with a list of extensions that DOES contain the extension passed in '%s' BUT the value is a negative duration", extTfResourceTimeout), func() {
			extensions := spec.Extensions{}
			extensions.Add(extTfResourceTimeout, "-1.5h")
			duration, err := r.getTimeDuration(extensions, extTfResourceTimeout)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("Then the duration returned should be nil", func() {
				So(duration, ShouldBeNil)
			})
			Convey("And the error message should be", func() {
				So(err.Error(), ShouldContainSubstring, "invalid duration value: '-1.5h'. The value must be a sequence of decimal numbers each with optional fraction and a unit suffix (negative durations are not allowed). The value must be formatted either in seconds (s), minutes (m) or hours (h)")
			})
		})
		Convey(fmt.Sprintf("When getTimeDuration method is called with a list of extensions that DOES contain the extension passed in '%s' BUT the value is NOT supported (distinct than s,m and h)", extTfResourceTimeout), func() {
			extensions := spec.Extensions{}
			extensions.Add(extTfResourceTimeout, "300ms")
			duration, err := r.getTimeDuration(extensions, extTfResourceTimeout)
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("Then the duration returned should be nil", func() {
				So(duration, ShouldBeNil)
			})
			Convey("And the error message should be", func() {
				So(err.Error(), ShouldContainSubstring, "invalid duration value: '300ms'. The value must be a sequence of decimal numbers each with optional fraction and a unit suffix (negative durations are not allowed). The value must be formatted either in seconds (s), minutes (m) or hours (h)")
			})
		})
	})
}

func TestGetDuration(t *testing.T) {
	Convey("Given a SpecV2Resource", t, func() {
		r := SpecV2Resource{}
		Convey("When getDuration method is called a valid formatted time'", func() {
			duration, err := r.getDuration("30s")
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Then the statuses returned should contain", func() {
				fmt.Println(duration)
				So(*duration, ShouldEqual, time.Duration(30*time.Second))
			})
		})
		Convey("When getDuration method is called a invalid formatted time'", func() {
			_, err := r.getDuration("some invalid formatted time")
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestGetResourceOverrideHost(t *testing.T) {
	Convey("Given a terraform compliant resource that has a POST operation containing the x-terraform-resource-host with a non parametrized host containing the host to use", t, func() {
		expectedHost := "some.api.domain.com"
		r := SpecV2Resource{
			RootPathItem: spec.PathItem{
				PathItemProps: spec.PathItemProps{
					Post: &spec.Operation{
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								extTfResourceURL: expectedHost,
							},
						},
					},
				},
			},
		}
		Convey("When getResourceOverrideHost method is called", func() {
			host := getResourceOverrideHost(r.RootPathItem.Post)
			Convey("Then the value returned should be the host value", func() {
				So(host, ShouldEqual, expectedHost)
			})
		})
	})

	Convey("Given a terraform compliant resource that has a POST operation containing the x-terraform-resource-host with a parametrized host containing the host to use", t, func() {
		expectedHost := "some.api.${serviceProviderName}.domain.com"
		r := SpecV2Resource{
			RootPathItem: spec.PathItem{
				PathItemProps: spec.PathItemProps{
					Post: &spec.Operation{
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								extTfResourceURL: expectedHost,
							},
						},
					},
				},
			},
		}
		Convey("When getResourceOverrideHost method is called", func() {
			host := getResourceOverrideHost(r.RootPathItem.Post)
			Convey("Then the value returned should be the host value", func() {
				So(host, ShouldEqual, expectedHost)
			})
		})
	})

	Convey("Given a terraform compliant resource that has a POST operation containing the x-terraform-resource-host with an empty string value", t, func() {
		expectedHost := ""
		r := SpecV2Resource{
			RootPathItem: spec.PathItem{
				PathItemProps: spec.PathItemProps{
					Post: &spec.Operation{
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								extTfResourceURL: expectedHost,
							},
						},
					},
				},
			},
		}
		Convey("When getResourceOverrideHost method is called", func() {
			host := getResourceOverrideHost(r.RootPathItem.Post)
			Convey("Then the value returned should be the host value", func() {
				So(host, ShouldEqual, expectedHost)
			})
		})
	})
}
