package openapi

import (
	"fmt"
	"github.com/go-openapi/spec"
	"github.com/hashicorp/terraform/helper/schema"
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"testing"
)

func TestGetResourceURL(t *testing.T) {

	Convey("Given resource info is configured with https scheme and basePath='/', path='/v1/resource', host='www.host.com'", t, func() {
		expectedBasePath := "/"
		expectedPath := "/v1/resource"
		expectedHost := "www.host.com"
		expectedScheme := "https"

		r := resourceInfo{
			basePath:    expectedBasePath,
			path:        expectedPath,
			host:        expectedHost,
			httpSchemes: []string{expectedScheme},
		}
		Convey("When getResourceURL method is called'", func() {
			resourceURL, err := r.getResourceURL()
			Convey("Then the value returned should be https://www.host.com/v1/resource and the error should be nil", func() {
				So(err, ShouldBeNil)
				So(resourceURL, ShouldEqual, fmt.Sprintf("%s://%s%s", expectedScheme, expectedHost, expectedPath))
			})
		})
	})

	Convey("Given resource info is configured with http scheme and basePath='/', path='/v1/resource', host='www.host.com'", t, func() {
		expectedBasePath := "/"
		expectedPath := "/v1/resource"
		expectedHost := "www.host.com"
		expectedScheme := "http"

		r := resourceInfo{
			basePath:    expectedBasePath,
			path:        expectedPath,
			host:        expectedHost,
			httpSchemes: []string{expectedScheme},
		}
		Convey("When getResourceURL method is called'", func() {
			resourceURL, err := r.getResourceURL()
			Convey("Then the value returned should be http://www.host.com/v1/resource and the error should be nil", func() {
				So(err, ShouldBeNil)
				So(resourceURL, ShouldEqual, fmt.Sprintf("%s://%s%s", expectedScheme, expectedHost, expectedPath))
			})
		})
	})

	Convey("Given resource info is not configured with any scheme and basePath='/', path='/v1/resource', host='www.host.com'", t, func() {
		expectedBasePath := "/"
		expectedPath := "/v1/resource"
		expectedHost := "www.host.com"
		expectedScheme := "http"

		r := resourceInfo{
			basePath:    expectedBasePath,
			path:        expectedPath,
			host:        expectedHost,
			httpSchemes: []string{},
		}
		Convey("When getResourceURL method is called'", func() {
			resourceURL, err := r.getResourceURL()
			Convey("Then the value returned should use the default scheme which is http://www.host.com/v1/resource and the error should be nil", func() {
				So(err, ShouldBeNil)
				So(resourceURL, ShouldEqual, fmt.Sprintf("%s://%s%s", expectedScheme, expectedHost, expectedPath))
			})
		})
	})

	Convey("Given resource info is configured with https scheme, basePath is not empty nor / and path is '/v1/resource''", t, func() {
		expectedBasePath := "/api"
		expectedPath := "/v1/resource"
		expectedHost := "www.host.com"
		expectedScheme := "https"

		r := resourceInfo{
			basePath:    expectedBasePath,
			path:        expectedPath,
			host:        expectedHost,
			httpSchemes: []string{expectedScheme},
		}
		Convey("When getResourceURL method is called'", func() {
			resourceURL, err := r.getResourceURL()
			Convey("Then the value returned should use the default scheme which is https://www.host.com/api/v1/resource and the error should be nil", func() {
				So(err, ShouldBeNil)
				So(resourceURL, ShouldEqual, fmt.Sprintf("%s://%s%s%s", expectedScheme, expectedHost, expectedBasePath, expectedPath))
			})
		})
	})

	Convey("Given resource info is configured with https scheme, basePath is empty and path is '/v1/resource''", t, func() {
		expectedBasePath := ""
		expectedPath := "/v1/resource"
		expectedHost := "www.host.com"
		expectedScheme := "https"

		r := resourceInfo{
			basePath:    expectedBasePath,
			path:        expectedPath,
			host:        expectedHost,
			httpSchemes: []string{expectedScheme},
		}
		Convey("When getResourceURL method is called'", func() {
			resourceURL, err := r.getResourceURL()
			Convey("Then the value returned should use the default scheme which is https://www.host.com/v1/resource and the error should be nil", func() {
				So(err, ShouldBeNil)
				So(resourceURL, ShouldEqual, fmt.Sprintf("%s://%s%s", expectedScheme, expectedHost, expectedPath))
			})
		})
	})

	Convey("Given resource info is configured with https scheme, basePath is / and path is '/v1/resource''", t, func() {
		expectedBasePath := "/"
		expectedPath := "/v1/resource"
		expectedHost := "www.host.com"
		expectedScheme := "https"

		r := resourceInfo{
			basePath:    expectedBasePath,
			path:        expectedPath,
			host:        expectedHost,
			httpSchemes: []string{expectedScheme},
		}
		Convey("When getResourceURL method is called'", func() {
			resourceURL, err := r.getResourceURL()
			Convey("Then the value returned should use the default scheme which is https://www.host.com/v1/resource and the error should be nil", func() {
				So(err, ShouldBeNil)
				So(resourceURL, ShouldEqual, fmt.Sprintf("%s://%s%s", expectedScheme, expectedHost, expectedPath))
			})
		})
	})

	Convey("Given resource info is configured with https scheme, basePath does not start with / and path is '/v1/resource''", t, func() {
		expectedBasePath := "api"
		expectedPath := "/v1/resource"
		expectedHost := "www.host.com"
		expectedScheme := "https"

		r := resourceInfo{
			basePath:    expectedBasePath,
			path:        expectedPath,
			host:        expectedHost,
			httpSchemes: []string{expectedScheme},
		}
		Convey("When getResourceURL method is called'", func() {
			resourceURL, err := r.getResourceURL()
			Convey("Then the value returned should use the default scheme which is https://www.host.com/v1/resource and the error should be nil", func() {
				So(err, ShouldBeNil)
				So(resourceURL, ShouldEqual, fmt.Sprintf("%s://%s/%s%s", expectedScheme, expectedHost, expectedBasePath, expectedPath))
			})
		})
	})

	Convey("Given resource info is configured with a path that does not start with /", t, func() {
		expectedBasePath := "/"
		expectedPath := "v1/resource"
		expectedHost := "www.host.com"
		expectedScheme := "https"

		r := resourceInfo{
			basePath:    expectedBasePath,
			path:        expectedPath,
			host:        expectedHost,
			httpSchemes: []string{expectedScheme},
		}
		Convey("When getResourceURL method is called'", func() {
			resourceURL, err := r.getResourceURL()
			Convey("Then the value returned should use the default scheme which is https://www.host.com/v1/resource and the error should be nil", func() {
				So(err, ShouldBeNil)
				So(resourceURL, ShouldEqual, fmt.Sprintf("%s://%s/%s", expectedScheme, expectedHost, expectedPath))
			})
		})
	})

	Convey("Given resource info is missing the path", t, func() {
		expectedBasePath := ""
		expectedPath := "/v1/resource"
		expectedScheme := "https"

		r := resourceInfo{
			basePath:    expectedBasePath,
			path:        expectedPath,
			host:        "",
			httpSchemes: []string{expectedScheme},
		}
		Convey("When getResourceURL method is called'", func() {
			_, err := r.getResourceURL()
			Convey("Then there should be returned error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Given resource info is missing the host", t, func() {
		expectedBasePath := ""
		expectedHost := "www.host.com"
		expectedScheme := "https"

		r := resourceInfo{
			basePath:    expectedBasePath,
			path:        "",
			host:        expectedHost,
			httpSchemes: []string{expectedScheme},
		}
		Convey("When getResourceURL method is called'", func() {
			_, err := r.getResourceURL()
			Convey("Then there should be returned error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestGetResourceIDURL(t *testing.T) {
	Convey("Given resource info is configured with 'https' scheme and basePath='/', path='/v1/resource', host='www.host.com'", t, func() {
		expectedBasePath := "/"
		expectedPath := "/v1/resource"
		expectedHost := "www.host.com"
		expectedScheme := "https"

		r := resourceInfo{
			basePath:    expectedBasePath,
			path:        expectedPath,
			host:        expectedHost,
			httpSchemes: []string{expectedScheme},
		}
		Convey("When getResourceIDURL method is called with ID=1234", func() {
			id := "1234"
			resourceIDURL, err := r.getResourceIDURL(id)
			Convey("Then the value returned should be https://www.host.com/v1/resource/1234 and the error should be nil", func() {
				So(err, ShouldBeNil)
				So(resourceIDURL, ShouldEqual, fmt.Sprintf("%s://%s%s/%s", expectedScheme, expectedHost, expectedPath, id))
			})
		})
	})

	Convey("Given resource info is missing the host", t, func() {
		expectedBasePath := ""
		expectedHost := "www.host.com"
		expectedScheme := "https"

		r := resourceInfo{
			basePath:    expectedBasePath,
			path:        "",
			host:        expectedHost,
			httpSchemes: []string{expectedScheme},
		}
		Convey("When getResourceIDURL method is called with ID=1234", func() {
			_, err := r.getResourceIDURL("1234")
			Convey("Then there should be returned error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Given resource info is missing the path", t, func() {
		expectedBasePath := ""
		expectedPath := "/v1/resource"
		expectedScheme := "https"

		r := resourceInfo{
			basePath:    expectedBasePath,
			path:        expectedPath,
			host:        "",
			httpSchemes: []string{expectedScheme},
		}
		Convey("When getResourceIDURL method is called with ID=1234", func() {
			_, err := r.getResourceIDURL("1234")
			Convey("Then there should be returned error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestGetImmutableProperties(t *testing.T) {
	Convey("Given resource info is configured with schemaDefinition that contains a property 'immutable_property' that is immutable", t, func() {
		extensions := spec.Extensions{}
		extensions.Add("x-terraform-immutable", true)
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"id": {
							VendorExtensible: spec.VendorExtensible{Extensions: extensions},
						},
						"immutable_property": {
							VendorExtensible: spec.VendorExtensible{Extensions: extensions},
						},
					},
				},
			},
		}
		Convey("When getImmutableProperties method is called", func() {
			immutableProperties := r.getImmutableProperties()
			Convey("Then the array returned should contain 'immutable_property'", func() {
				So(immutableProperties, ShouldContain, "immutable_property")
			})
			Convey("And the 'id' property should be ignored even if it's marked as immutable (should never happen though, edge case)", func() {
				So(immutableProperties, ShouldNotContain, "id")
			})
		})
	})

	Convey("Given resource info is configured with schemaDefinition that DOES NOT contain immutable properties", t, func() {
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"id": {
							VendorExtensible: spec.VendorExtensible{},
						},
						"mutable_property": {
							VendorExtensible: spec.VendorExtensible{Extensions: spec.Extensions{}},
						},
					},
				},
			},
		}
		Convey("When getImmutableProperties method is called", func() {
			immutableProperties := r.getImmutableProperties()
			Convey("Then the array returned should be empty", func() {
				So(immutableProperties, ShouldBeEmpty)
			})
		})
	})

}

func TestCreateTerraformPropertyBasicSchema(t *testing.T) {
	Convey("Given a swagger schema definition that has a property of type 'string'", t, func() {
		propSchema := spec.Schema{
			VendorExtensible: spec.VendorExtensible{},
			SchemaProps: spec.SchemaProps{
				Type: []string{"string"},
			},
		}
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"string_prop": propSchema,
					},
				},
			},
		}
		Convey("When createTerraformPropertyBasicSchema method is called", func() {
			tfPropSchema, err := r.createTerraformPropertyBasicSchema("string_prop", propSchema)
			Convey("Then the resulted terraform property schema should be of type string too", func() {
				So(err, ShouldBeNil)
				So(tfPropSchema.Type, ShouldEqual, schema.TypeString)
			})
		})
	})

	Convey("Given a swagger schema definition that has a property of type 'integer'", t, func() {
		propSchema := spec.Schema{
			VendorExtensible: spec.VendorExtensible{},
			SchemaProps: spec.SchemaProps{
				Type: []string{"integer"},
			},
		}
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"int_prop": propSchema,
					},
				},
			},
		}
		Convey("When createTerraformPropertyBasicSchema method is called", func() {
			tfPropSchema, err := r.createTerraformPropertyBasicSchema("int_prop", propSchema)
			Convey("Then the resulted terraform property schema should be of type int too", func() {
				So(err, ShouldBeNil)
				So(tfPropSchema.Type, ShouldEqual, schema.TypeInt)
			})
		})
	})

	Convey("Given a swagger schema definition that has a property of type 'number'", t, func() {
		propSchema := spec.Schema{
			VendorExtensible: spec.VendorExtensible{},
			SchemaProps: spec.SchemaProps{
				Type: []string{"number"},
			},
		}
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"number_prop": propSchema,
					},
				},
			},
		}
		Convey("When createTerraformPropertyBasicSchema method is called", func() {
			tfPropSchema, err := r.createTerraformPropertyBasicSchema("number_prop", propSchema)
			Convey("Then the resulted terraform property schema should be of type float too", func() {
				So(err, ShouldBeNil)
				So(tfPropSchema.Type, ShouldEqual, schema.TypeFloat)
			})
		})
	})

	Convey("Given a swagger schema definition that has a property of type 'boolean'", t, func() {
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"boolean_prop": {
							VendorExtensible: spec.VendorExtensible{},
							SchemaProps: spec.SchemaProps{
								Type: []string{"boolean"},
							},
						},
					},
				},
			},
		}
		Convey("When createTerraformPropertyBasicSchema method is called", func() {
			tfPropSchema, err := r.createTerraformPropertyBasicSchema("boolean_prop", r.schemaDefinition.Properties["boolean_prop"])
			Convey("Then the resulted terraform property schema should be of type int too", func() {
				So(err, ShouldBeNil)
				So(tfPropSchema.Type, ShouldEqual, schema.TypeBool)
			})
		})
	})

	Convey("Given a swagger schema definition that has a property of type 'array'", t, func() {
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"array_prop": {
							VendorExtensible: spec.VendorExtensible{},
							SchemaProps: spec.SchemaProps{
								Type: []string{"array"},
							},
						},
					},
				},
			},
		}
		Convey("When createTerraformPropertyBasicSchema method is called", func() {
			tfPropSchema, err := r.createTerraformPropertyBasicSchema("array_prop", r.schemaDefinition.Properties["array_prop"])
			Convey("Then the resulted terraform property schema should be of type array too", func() {
				So(err, ShouldBeNil)
				So(tfPropSchema.Type, ShouldEqual, schema.TypeList)
			})
			Convey("And the array elements are of the default type string (only supported type for now)", func() {
				So(reflect.TypeOf(tfPropSchema.Elem).Elem(), ShouldEqual, reflect.TypeOf(schema.Schema{}))
				So(tfPropSchema.Elem.(*schema.Schema).Type, ShouldEqual, schema.TypeString)
			})
		})
	})

	Convey("Given a swagger schema definition that has a property 'string_prop' which is required", t, func() {
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"string_prop": {
							VendorExtensible: spec.VendorExtensible{},
							SchemaProps: spec.SchemaProps{
								Type: []string{"string"},
							},
						},
					},
					Required: []string{"string_prop"}, // This array contains the list of properties that are required
				},
			},
		}
		Convey("When createTerraformPropertyBasicSchema method is called", func() {
			tfPropSchema, err := r.createTerraformPropertyBasicSchema("string_prop", r.schemaDefinition.Properties["string_prop"])
			Convey("Then the returned value should be true", func() {
				So(err, ShouldBeNil)
				So(tfPropSchema.Required, ShouldBeTrue)
			})
		})
	})

	Convey("Given a swagger property schema with 'x-terraform-force-new' metadata", t, func() {
		extensions := spec.Extensions{}
		extensions.Add("x-terraform-force-new", true)
		propSchema := spec.Schema{
			VendorExtensible: spec.VendorExtensible{Extensions: extensions},
			SchemaProps: spec.SchemaProps{
				Type: []string{"boolean"},
			},
		}
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"boolean_prop": propSchema,
					},
				},
			},
		}
		Convey("When createTerraformPropertyBasicSchema method is called", func() {
			tfPropSchema, err := r.createTerraformPropertyBasicSchema("boolean_prop", propSchema)
			Convey("Then the resulted terraform property schema should be of type int too", func() {
				So(err, ShouldBeNil)
				So(tfPropSchema.ForceNew, ShouldBeTrue)
			})
		})
	})

	Convey("Given a swagger property schema configured with readOnly (computed)", t, func() {
		propSchema := spec.Schema{
			VendorExtensible: spec.VendorExtensible{},
			SchemaProps: spec.SchemaProps{
				Type: []string{"boolean"},
			},
			SwaggerSchemaProps: spec.SwaggerSchemaProps{ReadOnly: true},
		}
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"boolean_prop": propSchema,
					},
				},
			},
		}
		Convey("When createTerraformPropertyBasicSchema method is called", func() {
			tfPropSchema, err := r.createTerraformPropertyBasicSchema("boolean_prop", propSchema)
			Convey("Then the resulted terraform property schema should be configured as computed", func() {
				So(err, ShouldBeNil)
				So(tfPropSchema.Computed, ShouldBeTrue)
			})
		})
	})

	Convey("Given a swagger property schema configured with 'x-terraform-force-new' and 'x-terraform-sensitive' metadata", t, func() {
		extensions := spec.Extensions{}
		extensions.Add("x-terraform-force-new", true)
		extensions.Add("x-terraform-sensitive", true)
		propSchema := spec.Schema{
			VendorExtensible: spec.VendorExtensible{Extensions: extensions},
			SchemaProps: spec.SchemaProps{
				Type: []string{"boolean"},
			},
		}
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"boolean_prop": propSchema,
					},
				},
			},
		}
		Convey("When createTerraformPropertyBasicSchema method is called", func() {
			tfPropSchema, err := r.createTerraformPropertyBasicSchema("boolean_prop", propSchema)
			Convey("Then the resulted terraform property schema should be configured as forceNew and sensitive", func() {
				So(err, ShouldBeNil)
				So(tfPropSchema.ForceNew, ShouldBeTrue)
				So(tfPropSchema.Sensitive, ShouldBeTrue)
			})
		})
	})

	Convey("Given a swagger property schema configured with default value", t, func() {
		expectedDefaultValue := "defaultValue"
		propSchema := spec.Schema{
			VendorExtensible: spec.VendorExtensible{},
			SchemaProps: spec.SchemaProps{
				Type:    []string{"boolean"},
				Default: expectedDefaultValue,
			},
		}
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"boolean_prop": propSchema,
					},
				},
			},
		}
		Convey("When createTerraformPropertyBasicSchema method is called", func() {
			tfPropSchema, err := r.createTerraformPropertyBasicSchema("boolean_prop", propSchema)
			Convey("Then the resulted terraform property schema should be configured with the according default value, ", func() {
				So(err, ShouldBeNil)
				So(tfPropSchema.Default, ShouldEqual, expectedDefaultValue)
			})
		})
	})
}

func TestIsArrayProperty(t *testing.T) {
	Convey("Given a swagger property schema of type 'array'", t, func() {
		propSchema := spec.Schema{
			VendorExtensible: spec.VendorExtensible{},
			SchemaProps: spec.SchemaProps{
				Type: []string{"array"},
			},
		}
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"array_prop": propSchema,
					},
				},
			},
		}
		Convey("When isArrayProperty method is called", func() {
			isArray := r.isArrayProperty(propSchema)
			Convey("Then the returned value should be true", func() {
				So(isArray, ShouldBeTrue)
			})
		})
	})

	Convey("Given a swagger property schema of type different than 'array'", t, func() {
		propSchema := spec.Schema{
			VendorExtensible: spec.VendorExtensible{},
			SchemaProps: spec.SchemaProps{
				Type: []string{"string"},
			},
		}
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"string_prop": propSchema,
					},
				},
			},
		}
		Convey("When isArrayProperty method is called", func() {
			isArray := r.isArrayProperty(propSchema)
			Convey("Then the returned value should be false", func() {
				So(isArray, ShouldBeFalse)
			})
		})
	})
}

func TestIsRequired(t *testing.T) {
	Convey("Given a swagger schema definition that has a property 'string_prop' that is required", t, func() {
		propSchema := spec.Schema{
			VendorExtensible: spec.VendorExtensible{},
			SchemaProps: spec.SchemaProps{
				Type: []string{"string"},
			},
		}
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"string_prop": propSchema,
					},
					Required: []string{"string_prop"}, // This array contains the list of properties that are required
				},
			},
		}
		Convey("When isRequired method is called", func() {
			isRequired := r.isRequired("string_prop", r.schemaDefinition.Required)
			Convey("Then the returned value should be true", func() {
				So(isRequired, ShouldBeTrue)
			})
		})
	})

	Convey("Given a swagger schema definition that has a property 'string_prop' that is not required", t, func() {
		propSchema := spec.Schema{
			VendorExtensible: spec.VendorExtensible{},
			SchemaProps: spec.SchemaProps{
				Type: []string{"string"},
			},
		}
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"string_prop": propSchema,
					},
				},
			},
		}
		Convey("When isRequired method is called", func() {
			isRequired := r.isRequired("string_prop", r.schemaDefinition.Required)
			Convey("Then the returned value should be false", func() {
				So(isRequired, ShouldBeFalse)
			})
		})
	})
}

func TestCreateTerraformPropertySchema(t *testing.T) {
	Convey("Given a swagger schema definition that has a property 'string_prop' of type string, required, sensitive and has a default value 'defaultValue'", t, func() {
		expectedDefaultValue := "defaultValue"
		extensions := spec.Extensions{}
		extensions.Add("x-terraform-sensitive", true)
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"string_prop": {
							VendorExtensible: spec.VendorExtensible{Extensions: extensions},
							SchemaProps: spec.SchemaProps{
								Type:    []string{"string"},
								Default: expectedDefaultValue,
							},
						},
					},
					Required: []string{"string_prop"}, // This array contains the list of properties that are required
				},
			},
		}
		Convey("When createTerraformPropertyBasicSchema method is called", func() {
			tfPropSchema, err := r.createTerraformPropertySchema("string_prop", r.schemaDefinition.Properties["string_prop"])
			Convey("Then the returned tf tfPropSchema should be of type string", func() {
				So(err, ShouldBeNil)
				So(tfPropSchema.Type, ShouldEqual, schema.TypeString)
			})
			Convey("And a validateFunc should be configured", func() {
				So(tfPropSchema.ValidateFunc, ShouldNotBeNil)
			})
			Convey("And be configured as required, sensitive and the default value should match 'defaultValue'", func() {
				So(tfPropSchema.Required, ShouldBeTrue)
			})
			Convey("And be configured as sensitive", func() {
				So(tfPropSchema.Sensitive, ShouldBeTrue)
			})
			Convey("And the default value should match 'defaultValue'", func() {
				So(tfPropSchema.Default, ShouldEqual, expectedDefaultValue)
			})
		})
	})

	Convey("Given a swagger schema definition that has a property 'array_prop' of type array", t, func() {
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"array_prop": {
							VendorExtensible: spec.VendorExtensible{},
							SchemaProps: spec.SchemaProps{
								Type: []string{"array"},
							},
						},
					},
					Required: []string{"array_prop"}, // This array contains the list of properties that are required
				},
			},
		}
		Convey("When createTerraformPropertyBasicSchema method is called", func() {
			tfPropSchema, err := r.createTerraformPropertySchema("array_prop", r.schemaDefinition.Properties["array_prop"])
			Convey("Then the returned tf tfPropSchema should be of type array", func() {
				So(err, ShouldBeNil)
				So(tfPropSchema.Type, ShouldEqual, schema.TypeList)
			})
			Convey("And there should not be any validation function attached to it", func() {
				So(tfPropSchema.ValidateFunc, ShouldBeNil)
			})
		})
	})
}

func TestValidateFunc(t *testing.T) {
	Convey("Given a swagger schema definition that has one property", t, func() {
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"array_prop": {
							VendorExtensible: spec.VendorExtensible{},
							SchemaProps: spec.SchemaProps{
								Type: []string{"array"},
							},
						},
					},
				},
			},
		}
		Convey("When validateFunc method is called", func() {
			validateFunc := r.validateFunc("array_prop", r.schemaDefinition.Properties["array_prop"])
			Convey("Then the returned validateFunc should not be nil", func() {
				So(validateFunc, ShouldNotBeNil)
			})
		})
	})

	Convey("Given a swagger schema definition that has a property which is supposed to be computed but has a default value set", t, func() {
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"array_prop": {
							VendorExtensible: spec.VendorExtensible{},
							SchemaProps: spec.SchemaProps{
								Type:    []string{"array"},
								Default: "defaultValue",
							},
							SwaggerSchemaProps: spec.SwaggerSchemaProps{ReadOnly: true},
						},
					},
				},
			},
		}
		Convey("When validateFunc method is called", func() {
			validateFunc := r.validateFunc("array_prop", r.schemaDefinition.Properties["array_prop"])
			Convey("Then the returned validateFunc should not be nil", func() {
				So(validateFunc, ShouldNotBeNil)
			})
			Convey("And when the function is executed it should return an error as computed properties can not have default values", func() {
				_, errs := validateFunc("", "")
				So(errs, ShouldNotBeEmpty)
			})
		})
	})
}

func TestCreateTerraformResourceSchema(t *testing.T) {
	Convey("Given a swagger schema definition that has multiple properties - 'string_prop', 'int_prop', 'number_prop', 'bool_prop' and 'array_prop'", t, func() {
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"string_prop": {
							VendorExtensible: spec.VendorExtensible{},
							SchemaProps: spec.SchemaProps{
								Type: []string{"string"},
							},
						},
						"intProp": { // This prop does not have a terraform compliant name; however an automatic translation is performed behind the scenes to make it compliant
							VendorExtensible: spec.VendorExtensible{},
							SchemaProps: spec.SchemaProps{
								Type: []string{"integer"},
							},
						},
						"number_prop": {
							VendorExtensible: spec.VendorExtensible{
								Extensions: spec.Extensions{
									"x-terraform-field-name": "numberProp", // making use of specific extension to override field name; but the new field name is not terrafrom name compliant - hence an automatic translation is performed behind the scenes to make it compliant
								},
							},
							SchemaProps: spec.SchemaProps{
								Type: []string{"number"},
							},
						},
						"bool_prop": {
							VendorExtensible: spec.VendorExtensible{},
							SchemaProps: spec.SchemaProps{
								Type: []string{"boolean"},
							},
						},
						"arrayProp": {
							VendorExtensible: spec.VendorExtensible{},
							SchemaProps: spec.SchemaProps{
								Type: []string{"array"},
							},
						},
					},
				},
			},
		}
		Convey("When createTerraformResourceSchema method is called", func() {
			resourceSchema, err := r.createTerraformResourceSchema()
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the tf resource schema returned should match the swagger props - 'string_prop', 'int_prop', 'number_prop' and 'bool_prop' and 'array_prop', ", func() {
				So(resourceSchema, ShouldNotBeNil)
				So(resourceSchema, ShouldContainKey, "string_prop")
				So(resourceSchema, ShouldContainKey, "int_prop")
				So(resourceSchema, ShouldContainKey, "number_prop")
				So(resourceSchema, ShouldContainKey, "bool_prop")
				So(resourceSchema, ShouldContainKey, "array_prop")
			})
		})
	})
}

func TestConvertToTerraformCompliantFieldName(t *testing.T) {
	Convey("Given a property with a name that is terraform field name compliant", t, func() {
		propertyName := "some_prop_name_that_is_terraform_field_name_compliant"
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						propertyName: {
							VendorExtensible: spec.VendorExtensible{},
							SchemaProps: spec.SchemaProps{
								Type: []string{"string"},
							},
						},
					},
				},
			},
		}
		Convey("When convertToTerraformCompliantFieldName method is called", func() {
			fieldName := r.convertToTerraformCompliantFieldName(propertyName, r.schemaDefinition.Properties[propertyName])
			Convey("And string return is terraform field name compliant, ", func() {
				So(fieldName, ShouldEqual, propertyName)
			})
		})
	})

	Convey("Given a property with a name that is NOT terraform field name compliant", t, func() {
		propertyName := "thisPropIsNotTerraformField_Compliant"
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						propertyName: {
							VendorExtensible: spec.VendorExtensible{},
							SchemaProps: spec.SchemaProps{
								Type: []string{"string"},
							},
						},
					},
				},
			},
		}
		Convey("When convertToTerraformCompliantFieldName method is called", func() {
			fieldName := r.convertToTerraformCompliantFieldName(propertyName, r.schemaDefinition.Properties[propertyName])
			Convey("And string return is terraform field name compliant, ", func() {
				So(fieldName, ShouldEqual, "this_prop_is_not_terraform_field_compliant")
			})
		})
	})

	Convey("Given a property with a name that is NOT terraform field name compliant but has an extension that overrides it", t, func() {
		propertyName := "thisPropIsNotTerraformField_Compliant"
		expectedPropertyName := "this_property_is_now_terraform_field_compliant"
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						propertyName: {
							VendorExtensible: spec.VendorExtensible{
								Extensions: spec.Extensions{
									extTfFieldName: expectedPropertyName,
								},
							},
							SchemaProps: spec.SchemaProps{
								Type: []string{"string"},
							},
						},
					},
				},
			},
		}
		Convey("When convertToTerraformCompliantFieldName method is called", func() {
			fieldName := r.convertToTerraformCompliantFieldName(propertyName, r.schemaDefinition.Properties[propertyName])
			Convey("And string return is terraform field name compliant, ", func() {
				So(fieldName, ShouldEqual, expectedPropertyName)
			})
		})
	})

	Convey("Given a property with a name that is NOT terraform field name compliant but has an extension that overrides it which in turn is also not terraform name compliant", t, func() {
		propertyName := "thisPropIsNotTerraformField_Compliant"
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						propertyName: {
							VendorExtensible: spec.VendorExtensible{
								Extensions: spec.Extensions{
									extTfFieldName: "thisPropIsAlsoNotTerraformField_Compliant",
								},
							},
							SchemaProps: spec.SchemaProps{
								Type: []string{"string"},
							},
						},
					},
				},
			},
		}
		Convey("When convertToTerraformCompliantFieldName method is called", func() {
			fieldName := r.convertToTerraformCompliantFieldName(propertyName, r.schemaDefinition.Properties[propertyName])
			Convey("And string return is terraform field name compliant, ", func() {
				So(fieldName, ShouldEqual, "this_prop_is_also_not_terraform_field_compliant")
			})
		})
	})
}

func TestGetResourceIdentifier(t *testing.T) {
	Convey("Given a swagger schema definition that has an id property", t, func() {
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"id": {
							VendorExtensible: spec.VendorExtensible{},
							SchemaProps: spec.SchemaProps{
								Type: []string{"string"},
							},
						},
					},
				},
			},
		}
		Convey("When getResourceIdentifier method is called", func() {
			id, err := r.getResourceIdentifier()
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Then the value returned should be 'id'", func() {
				So(id, ShouldEqual, "id")
			})
		})
	})

	Convey("Given a swagger schema definition that DOES NOT have an 'id' property but has a property configured with x-terraform-id set to TRUE", t, func() {
		extensions := spec.Extensions{}
		extensions.Add("x-terraform-id", true)
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"some-other-id": {
							VendorExtensible: spec.VendorExtensible{Extensions: extensions},
							SchemaProps: spec.SchemaProps{
								Type: []string{"string"},
							},
						},
					},
				},
			},
		}
		Convey("When getResourceIdentifier method is called", func() {
			id, err := r.getResourceIdentifier()
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Then the value returned should be 'some-other-id'", func() {
				So(id, ShouldEqual, "some-other-id")
			})
		})
	})

	Convey("Given a swagger schema definition that HAS BOTH an 'id' property AND ALSO a property configured with x-terraform-id set to true", t, func() {
		extensions := spec.Extensions{}
		extensions.Add("x-terraform-id", true)
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"id": {
							VendorExtensible: spec.VendorExtensible{},
							SchemaProps: spec.SchemaProps{
								Type: []string{"string"},
							},
						},
						"some-other-id": {
							VendorExtensible: spec.VendorExtensible{Extensions: extensions},
							SchemaProps: spec.SchemaProps{
								Type: []string{"string"},
							},
						},
					},
				},
			},
		}
		Convey("When getResourceIdentifier method is called", func() {
			id, err := r.getResourceIdentifier()
			Convey("Then the error returned should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Then the value returned should be 'some-other-id' as it takes preference over the default 'id' property", func() {
				So(id, ShouldEqual, "some-other-id")
			})
		})
	})

	Convey("Given a swagger schema definition that DOES NOT have an 'id' property but has a property configured with x-terraform-id set to FALSE", t, func() {
		extensions := spec.Extensions{}
		extensions.Add("x-terraform-id", false)
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"some-other-id": {
							VendorExtensible: spec.VendorExtensible{Extensions: extensions},
							SchemaProps: spec.SchemaProps{
								Type: []string{"string"},
							},
						},
					},
				},
			},
		}
		Convey("When getResourceIdentifier method is called", func() {
			_, err := r.getResourceIdentifier()
			Convey("Then the error returned should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Given a swagger schema definition that NEITHER HAS an 'id' property NOR a property configured with x-terraform-id set to true", t, func() {
		r := resourceInfo{
			schemaDefinition: spec.Schema{
				SchemaProps: spec.SchemaProps{
					Properties: map[string]spec.Schema{
						"prop-that-is-not-id": {
							VendorExtensible: spec.VendorExtensible{},
							SchemaProps: spec.SchemaProps{
								Type: []string{"string"},
							},
						},
						"prop-that-is-not-id-and-does-not-have-id-metadata": {
							VendorExtensible: spec.VendorExtensible{},
							SchemaProps: spec.SchemaProps{
								Type: []string{"string"},
							},
						},
					},
				},
			},
		}
		Convey("When getResourceIdentifier method is called", func() {
			_, err := r.getResourceIdentifier()
			Convey("Then the error returned should NOT be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestShouldIgnoreResource(t *testing.T) {
	Convey("Given a terraform compliant resource that has a POST operation containing the x-terraform-exclude-resource with value true", t, func() {
		r := resourceInfo{
			createPathInfo: spec.PathItem{
				PathItemProps: spec.PathItemProps{
					Post: &spec.Operation{
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-terraform-exclude-resource": true,
							},
						},
					},
				},
			},
		}
		Convey("When shouldIgnoreResource method is called", func() {
			shouldIgnoreResource := r.shouldIgnoreResource()
			Convey("Then the value returned should be true", func() {
				So(shouldIgnoreResource, ShouldBeTrue)
			})
		})
	})
	Convey("Given a terraform compliant resource that has a POST operation containing the x-terraform-exclude-resource with value false", t, func() {
		r := resourceInfo{
			createPathInfo: spec.PathItem{
				PathItemProps: spec.PathItemProps{
					Post: &spec.Operation{
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-terraform-exclude-resource": false,
							},
						},
					},
				},
			},
		}
		Convey("When shouldIgnoreResource method is called", func() {
			shouldIgnoreResource := r.shouldIgnoreResource()
			Convey("Then the value returned should be true", func() {
				So(shouldIgnoreResource, ShouldBeFalse)
			})
		})
	})
	Convey("Given a terraform compliant resource that has a POST operation that DOES NOT contain the x-terraform-exclude-resource extension", t, func() {
		r := resourceInfo{
			createPathInfo: spec.PathItem{
				PathItemProps: spec.PathItemProps{
					Post: &spec.Operation{
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{},
						},
					},
				},
			},
		}
		Convey("When shouldIgnoreResource method is called", func() {
			shouldIgnoreResource := r.shouldIgnoreResource()
			Convey("Then the value returned should be true", func() {
				So(shouldIgnoreResource, ShouldBeFalse)
			})
		})
	})
}
