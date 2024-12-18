package rest

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sync"
	"testing"
	"time"

	_ "embed"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v2"
)

//go:embed tests/*
var testFiles embed.FS

type TestRequestType struct {
	IntField int
}

func (m TestRequestType) ApplyCustomSchema(s *openapi3.Schema) {
	s.Required = []string{
		"IntField",
	}
}

// TestResponseType description.
type TestResponseType struct {
	// IntField description.
	IntField int
}

func (m TestResponseType) ApplyCustomSchema(s *openapi3.Schema) {
	s.Required = []string{
		"IntField",
	}
}

type AllBasicDataTypes struct {
	Int     int
	Int8    int8
	Int16   int16
	Int32   int32
	Int64   int64
	Uint    uint
	Uint8   uint8
	Uint16  uint16
	Uint32  uint32
	Uint64  uint64
	Uintptr uintptr
	Float32 float32
	Float64 float64
	Byte    byte
	Rune    rune
	String  string
	Bool    bool
}

func (m AllBasicDataTypes) ApplyCustomSchema(s *openapi3.Schema) {
	s.Required = []string{
		"Int",
		"Int8",
		"Int16",
		"Int32",
		"Int64",
		"Uint",
		"Uint8",
		"Uint16",
		"Uint32",
		"Uint64",
		"Uintptr",
		"Float32",
		"Float64",
		"Byte",
		"Rune",
		"String",
		"Bool",
	}
}

type AllBasicDataTypesPointers struct {
	Int     *int
	Int8    *int8
	Int16   *int16
	Int32   *int32
	Int64   *int64
	Uint    *uint
	Uint8   *uint8
	Uint16  *uint16
	Uint32  *uint32
	Uint64  *uint64
	Uintptr *uintptr
	Float32 *float32
	Float64 *float64
	Byte    *byte
	Rune    *rune
	String  *string
	Bool    *bool
}

type OmitEmptyFields struct {
	A string
	B string `validate:",omitempty"`
	C *string
	D *string `validate:",omitempty"`
}

func (m OmitEmptyFields) ApplyCustomSchema(s *openapi3.Schema) {
	s.Required = []string{
		"A",
	}
}

type EmbeddedStructA struct {
	A string
}

func (m EmbeddedStructA) ApplyCustomSchema(s *openapi3.Schema) {
	s.Required = []string{
		"A",
	}
}

type EmbeddedStructB struct {
	B                string
	OptionalB        string `validate:",omitempty"`
	PointerB         *string
	OptionalPointerB *string `validate:",omitempty"`
}

func (m EmbeddedStructB) ApplyCustomSchema(s *openapi3.Schema) {
	s.Required = []string{
		"B",
	}
}

type WithEmbeddedStructs struct {
	EmbeddedStructA
	EmbeddedStructB
	C string
}

func (m WithEmbeddedStructs) ApplyCustomSchema(s *openapi3.Schema) {
	s.Required = []string{
		"A",
		"B",
		"C",
	}
}

type WithNameStructTags struct {
	// FirstName of something.
	FirstName string `json:"firstName"`
	// LastName of something.
	LastName string
	// FullName of something.
	// Deprecated: Use FirstName and LastName
	FullName string
	// MiddleName of something. Deprecated: This deprecation flag is not valid so this field should
	// not be marked as deprecated.
	MiddleName string
}

func (m WithNameStructTags) ApplyCustomSchema(s *openapi3.Schema) {
	s.Required = []string{
		"firstName",
		"LastName",
		"FullName",
		"MiddleName",
	}
}

type KnownTypes struct {
	Time    time.Time  `json:"time"`
	TimePtr *time.Time `json:"timePtr"`
}

func (m KnownTypes) ApplyCustomSchema(s *openapi3.Schema) {
	s.Required = []string{
		"time",
	}
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (m User) ApplyCustomSchema(s *openapi3.Schema) {
	s.Required = []string{
		"id",
		"name",
	}
}

type OK struct {
	OK bool `json:"ok"`
}

func (m OK) ApplyCustomSchema(s *openapi3.Schema) {
	s.Required = []string{
		"ok",
	}
}

type StringEnum string

const (
	StringEnumA StringEnum = "A"
	StringEnumB StringEnum = "B"
	StringEnumC StringEnum = "B"
)

type IntEnum int64

const (
	IntEnum1 IntEnum = 1
	IntEnum2 IntEnum = 2
	IntEnum3 IntEnum = 3
)

type WithEnums struct {
	S  StringEnum   `json:"s"`
	SS []StringEnum `json:"ss"`
	I  IntEnum      `json:"i"`
	V  string       `json:"v"`
}

func (m WithEnums) ApplyCustomSchema(s *openapi3.Schema) {
	s.Required = []string{
		"s",
		"ss",
		"i",
		"v",
	}
}

type Pence int64

type WithMaps struct {
	Amounts map[string]Pence `json:"amounts"`
}

func (m WithMaps) ApplyCustomSchema(s *openapi3.Schema) {
	s.Required = []string{
		"amounts",
	}
}

type MultipleDateFieldsWithComments struct {
	// DateField is a field containing a date
	DateField time.Time `json:"dateField"`
	// DateFieldA is another field containing a date
	DateFieldA time.Time `json:"dateFieldA"`
}

func (m MultipleDateFieldsWithComments) ApplyCustomSchema(s *openapi3.Schema) {
	s.Required = []string{
		"dateField",
		"dateFieldA",
	}
}

type StructWithCustomisation struct {
	A string                  `json:"a"`
	B FieldWithCustomisation  `json:"b"`
	C *FieldWithCustomisation `json:"c"`
}

func (*StructWithCustomisation) ApplyCustomSchema(s *openapi3.Schema) {
	s.Properties["a"].Value.Description = "A string"
	s.Properties["a"].Value.Example = "test"
	s.Properties["b"].Value.Description = "A custom field"
	s.Required = []string{
		"a",
		"b",
	}
}

type FieldWithCustomisation string

func (*FieldWithCustomisation) ApplyCustomSchema(s *openapi3.Schema) {
	s.Format = "custom"
	s.Example = "model_field_customisation"
}

type StructWithTags struct {
	A string `json:"a" rest:"A is a string."`
}

func (m StructWithTags) ApplyCustomSchema(s *openapi3.Schema) {
	s.Required = []string{
		"a",
	}
}

type RecursiveModelModel struct {
	Model *RecursiveModel `json:"model" validate:"omitempty"`
	Bar   string          `json:"bar" validate:"omitempty"`
}

type RecursiveModel struct {
	Recursive *RecursiveModelModel `json:"recursive" validate:"omitempty"`
	Foo       string               `json:"foo" validate:"omitempty"`
}

type WithSwaggerType struct {
	Foo []uint8 `json:"foo" swaggertype:"string" validate:"omitempty"`
}

type WithWithSwaggerType struct {
	*WithSwaggerType
	Bar string `json:"bar" validate:"omitempty"`
}

type WithExamplesMinMaxEnumSet struct {
	Foo   string  `json:"foo" minLength:"2" maxLength:"5"`
	Bar   int     `json:"bar" minimum:"0" maximum:"255"`
	Baz   string  `json:"baz" enums:"foo,bar,baz"`
	Qux   int     `json:"qux" enums:"1,2,3"`
	Fred  *string `json:"fred" enums:"foo,bar,baz"`
	Thud  float64 `json:"thud" minimum:"0" maximum:"9.9999"`
	Waldo float64 `json:"waldo" minimum:"0" maximum:"99999.9999999"`
	Set   string  `json:"set" set:"foo,bar"`
}

func (m WithExamplesMinMaxEnumSet) ApplyCustomSchema(s *openapi3.Schema) {
	s.Required = []string{
		"foo",
		"bar",
		"baz",
		"qux",
		"thud",
		"waldo",
		"set",
	}
}

func TestSchema(t *testing.T) {
	tests := []struct {
		name  string
		opts  []APIOpts
		setup func(api *API) error
	}{
		{
			name:  "test000.yaml",
			setup: func(api *API) error { return nil },
		},
		{
			name: "test001.yaml",
			setup: func(api *API) error {
				api.Post("/test").
					HasRequestModel(ModelOf[TestRequestType]()).
					HasResponseModel(http.StatusOK, ModelOf[TestResponseType]()).
					HasDescription("Test request type description").
					HasTags([]string{"TestRequest"})
				return nil
			},
		},
		{
			name: "basic-data-types.yaml",
			setup: func(api *API) error {
				api.Post("/test").
					HasRequestModel(ModelOf[AllBasicDataTypes]()).
					HasResponseModel(http.StatusOK, ModelOf[AllBasicDataTypes]()).
					HasOperationID("postAllBasicDataTypes").
					HasTags([]string{"BasicData"}).
					HasDescription("Post all basic data types description")
				return nil
			},
		},
		{
			name: "basic-data-types-pointers.yaml",
			setup: func(api *API) error {
				api.Post("/test").
					HasRequestModel(ModelOf[AllBasicDataTypesPointers]()).
					HasResponseModel(http.StatusOK, ModelOf[AllBasicDataTypesPointers]())
				return nil
			},
		},
		{
			name: "omit-empty-fields.yaml",
			setup: func(api *API) error {
				api.Post("/test").
					HasRequestModel(ModelOf[OmitEmptyFields]()).
					HasResponseModel(http.StatusOK, ModelOf[OmitEmptyFields]())
				return nil
			},
		},
		{
			name: "anonymous-type.yaml",
			setup: func(api *API) error {
				api.Post("/test").
					HasRequestModel(ModelOf[struct{ A string }]()).
					HasResponseModel(http.StatusOK, ModelOf[struct{ B string }]())
				return nil
			},
		},
		{
			name: "embedded-structs.yaml",
			setup: func(api *API) error {
				api.Get("/embedded").
					HasResponseModel(http.StatusOK, ModelOf[EmbeddedStructA]())
				api.Post("/test").
					HasRequestModel(ModelOf[WithEmbeddedStructs]()).
					HasResponseModel(http.StatusOK, ModelOf[WithEmbeddedStructs]())
				return nil
			},
		},
		{
			name: "with-name-struct-tags.yaml",
			setup: func(api *API) error {
				api.Post("/test").
					HasRequestModel(ModelOf[WithNameStructTags]()).
					HasResponseModel(http.StatusOK, ModelOf[WithNameStructTags]())
				return nil
			},
		},
		{
			name: "known-types.yaml",
			setup: func(api *API) error {
				api.Route(http.MethodGet, "/test").
					HasResponseModel(http.StatusOK, ModelOf[KnownTypes]())
				return nil
			},
		},
		{
			name: "recursive-models.yaml",
			setup: func(api *API) error {
				api.Get("/recursive-models").
					HasResponseModel(http.StatusOK, ModelOf[RecursiveModel]())
				return nil
			},
		},
		{
			name: "all-methods.yaml",
			setup: func(api *API) (err error) {
				api.Get("/get").HasResponseModel(http.StatusOK, ModelOf[OK]())
				api.Head("/head").HasResponseModel(http.StatusOK, ModelOf[OK]())
				api.Post("/post").HasResponseModel(http.StatusOK, ModelOf[OK]())
				api.Put("/put").HasResponseModel(http.StatusOK, ModelOf[OK]())
				api.Patch("/patch").HasResponseModel(http.StatusOK, ModelOf[OK]())
				api.Delete("/delete").HasResponseModel(http.StatusOK, ModelOf[OK]())
				api.Connect("/connect").HasResponseModel(http.StatusOK, ModelOf[OK]())
				api.Options("/options").HasResponseModel(http.StatusOK, ModelOf[OK]())
				api.Trace("/trace").HasResponseModel(http.StatusOK, ModelOf[OK]())
				return
			},
		},
		{
			name: "enums.yaml",
			setup: func(api *API) (err error) {
				// Register the enums and values.
				api.RegisterModel(ModelOf[StringEnum](), WithEnumValues(StringEnumA, StringEnumB, StringEnumC))
				api.RegisterModel(ModelOf[IntEnum](), WithEnumValues(IntEnum1, IntEnum2, IntEnum3))

				api.Get("/get").HasResponseModel(http.StatusOK, ModelOf[WithEnums]())
				return
			},
		},
		{
			name: "enum-constants.yaml",
			setup: func(api *API) (err error) {
				// Register the enums and values.
				api.RegisterModel(ModelOf[StringEnum](), WithEnumConstants[StringEnum]())
				api.RegisterModel(ModelOf[IntEnum](), WithEnumConstants[IntEnum]())

				api.Get("/get").HasResponseModel(http.StatusOK, ModelOf[WithEnums]())
				return
			},
		},
		{
			name: "with-maps.yaml",
			setup: func(api *API) (err error) {
				api.Get("/get").HasResponseModel(http.StatusOK, ModelOf[WithMaps]())
				return
			},
		},
		{
			name: "route-params.yaml",
			setup: func(api *API) (err error) {
				api.Get(`/organisation/{orgId:\d+}/user/{userId}`).
					HasPathParameter("orgId", PathParam{
						Description: "Organisation ID",
						Regexp:      `\d+`,
					}).
					HasPathParameter("userId", PathParam{
						Description: "User ID",
					}).
					HasResponseModel(http.StatusOK, ModelOf[User]())
				return
			},
		},
		{
			name: "route-params.yaml",
			setup: func(api *API) (err error) {
				api.Get(`/organisation/{orgId:\d+}/user/{userId}`).
					HasPathParameter("orgId", PathParam{
						Regexp: `\d+`,
						ApplyCustomSchema: func(s *openapi3.Parameter) {
							s.Description = "Organisation ID"
						},
					}).
					HasPathParameter("userId", PathParam{
						Description: "User ID",
					}).
					HasResponseModel(http.StatusOK, ModelOf[User]())
				return
			},
		},
		{
			name: "query-params.yaml",
			setup: func(api *API) (err error) {
				api.Get(`/users?orgId=123&orderBy=field`).
					HasQueryParameter("orgId", QueryParam{
						Description: "ID of the organisation",
						Required:    true,
						Type:        PrimitiveTypeInteger,
					}).
					HasQueryParameter("orderBy", QueryParam{
						Description: "The field to order the results by",
						Required:    false,
						Type:        PrimitiveTypeString,
						Regexp:      `field|otherField`,
					}).
					HasResponseModel(http.StatusOK, ModelOf[User]())
				return
			},
		},
		{
			name: "query-params.yaml",
			setup: func(api *API) (err error) {
				api.Get(`/users?orgId=123&orderBy=field`).
					HasQueryParameter("orgId", QueryParam{
						Required: true,
						Type:     PrimitiveTypeInteger,
						ApplyCustomSchema: func(s *openapi3.Parameter) {
							s.Description = "ID of the organisation"
						},
					}).
					HasQueryParameter("orderBy", QueryParam{
						Required: false,
						Type:     PrimitiveTypeString,
						Regexp:   `field|otherField`,
						ApplyCustomSchema: func(s *openapi3.Parameter) {
							s.Description = "The field to order the results by"
						},
					}).
					HasResponseModel(http.StatusOK, ModelOf[User]())
				return
			},
		},
		{
			name: "multiple-dates-with-comments.yaml",
			setup: func(api *API) (err error) {
				api.Get("/dates").
					HasResponseModel(http.StatusOK, ModelOf[MultipleDateFieldsWithComments]())
				return
			},
		},
		{
			name: "custom-models.yaml",
			setup: func(api *API) (err error) {
				api.Get("/struct-with-customisation").
					HasResponseModel(http.StatusOK, ModelOf[StructWithCustomisation]())
				api.Get("/struct-ptr-with-customisation").
					HasResponseModel(http.StatusOK, ModelOf[*StructWithCustomisation]())
				return
			},
		},
		{
			name: "with-swaggertype.yaml",
			setup: func(api *API) (err error) {
				api.Get("/with-swaggertype").
					HasResponseModel(http.StatusOK, ModelOf[WithWithSwaggerType]())
				return
			},
		},
		{
			name: "with-model-examples.yaml",
			setup: func(api *API) (err error) {
				foo := map[string]any{"a": "foo"}
				api.RegisterModel(ModelOf[StructWithTags](), WithExample(foo))
				api.Get("/with-model-examples").
					HasResponseModel(http.StatusOK, ModelOf[StructWithTags]())
				return
			},
		},
		{
			name: "with-field-examples-min-max.yaml",
			setup: func(api *API) (err error) {
				api.RegisterModel(ModelOf[WithExamplesMinMaxEnumSet]())
				api.Get("/with-field-examples-min-max").
					HasResponseModel(http.StatusOK, ModelOf[WithExamplesMinMaxEnumSet]())
				return
			},
		},
		{
			name: "global-customisation.yaml",
			opts: []APIOpts{
				WithApplyCustomSchemaToType(func(t reflect.Type, s *openapi3.Schema) {
					if t != reflect.TypeOf(StructWithTags{}) {
						return
					}
					for fi := 0; fi < t.NumField(); fi++ {
						// Get the field name.
						var name string
						name = t.Field(fi).Tag.Get("json")
						if name == "" {
							name = t.Field(fi).Name
						}

						// Get the custom description from the struct tag.
						desc := t.Field(fi).Tag.Get("rest")
						if desc == "" {
							continue
						}
						if s.Properties == nil {
							s.Properties = make(map[string]*openapi3.SchemaRef)
						}
						if s.Properties[name] == nil {
							s.Properties[name] = &openapi3.SchemaRef{
								Value: &openapi3.Schema{},
							}
						}
						s.Properties[name].Value.Description = desc
					}
				}),
			},
			setup: func(api *API) error {
				api.Get("/").
					HasResponseModel(http.StatusOK, ModelOf[StructWithTags]())
				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var expected, actual []byte

			var wg sync.WaitGroup
			wg.Add(2)
			errs := make([]error, 2)

			// Validate the test file.
			go func() {
				defer wg.Done()
				// Load test file.
				expectedYAML, err := testFiles.ReadFile("tests/" + test.name)
				if err != nil {
					errs[0] = fmt.Errorf("could not read file %q: %v", test.name, err)
					return
				}
				expectedSpec, err := openapi3.NewLoader().LoadFromData(expectedYAML)
				if err != nil {
					errs[0] = fmt.Errorf("error in expected YAML: %w", err)
					return
				}
				expected, errs[0] = specToYAML(expectedSpec)
			}()

			go func() {
				defer wg.Done()
				// Create the API.
				api := NewAPI(test.name, test.opts...)
				api.StripPkgPaths = []string{"github.com/heimspiel/rest"}
				// Configure it.
				test.setup(api)
				// Create the actual spec.
				spec, err := api.Spec()
				if err != nil {
					t.Errorf("failed to generate spec: %v", err)
				}
				actual, errs[1] = specToYAML(spec)
			}()

			wg.Wait()
			var setupFailed bool
			for _, err := range errs {
				if err != nil {
					setupFailed = true
					t.Error(err)
				}
			}
			if setupFailed {
				t.Fatal("test setup failed")
			}

			// Compare the JSON marshalled output to ignore unexported fields and internal state.
			if diff := cmp.Diff(expected, actual); diff != "" {
				t.Error(diff)
				t.Error("\n\n" + string(actual))
			}
		})
	}
}

func specToYAML(spec *openapi3.T) (out []byte, err error) {
	// Use JSON, because kin-openapi doesn't customise the YAML output.
	// For example, AdditionalProperties only has a MarshalJSON capability.
	out, err = json.Marshal(spec)
	if err != nil {
		err = fmt.Errorf("could not marshal spec to JSON: %w", err)
		return
	}
	var m map[string]interface{}
	err = json.Unmarshal(out, &m)
	if err != nil {
		return
	}
	return yaml.Marshal(m)
}
