package env

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func intPtr(i int) *int {
	return &i
}

type testConfig struct {
	PtrField                       *string                     `env:"PTR_FIELD"`
	StrField                       string                      `env:"STR_FIELD"`
	IntField                       int                         `env:"INT_FIELD"`
	IntFieldNotSet                 int                         `env:"INT_FIELD_NOT_SET"`
	BoolField                      bool                        `env:"BOOL_FIELD"`
	UintField                      uint                        `env:"UIT_FIELD"`
	FloatField                     float64                     `env:"FLOAT_FIELD"`
	ArrayField                     *arrayWithTextUnmarshaller  `env:"ARRAY_FIELD"`
	TypedField                     *enumWithTextUnmarshaller   `env:"TYPED_FIELD"`
	TypedFieldNoPtr                enumWithTextUnmarshaller    `env:"TYPED_FIELD"`
	StructField                    *structWithTextUnmarshaller `env:"STRUCT_FIELD"`
	StructFieldNoPtr               structWithTextUnmarshaller  `env:"STRUCT_FIELD"`
	StructNestedField              *structWithNestedField
	StructNestedFieldWithDive      *structWithNestedField `env:",dive"`
	StructNestedFieldWithDiveNoPtr structWithNestedField  `env:",dive"`
	unexportedField                string                 //lint:ignore U1000 used by unit tests
}

type enumWithTextUnmarshaller int

func (e *enumWithTextUnmarshaller) UnmarshalText(text []byte) error {
	if e == nil {
		return fmt.Errorf("unmarshal text: nil pointer")
	}

	const (
		a = iota
		b
		c
	)

	switch string(text) {
	case "a":
		*e = a
	case "b":
		*e = b
	default:
		*e = c
	}
	return nil
}

type arrayWithTextUnmarshaller []string

func (a *arrayWithTextUnmarshaller) UnmarshalText(text []byte) error {
	if a == nil {
		return fmt.Errorf("unmarshal text: nil pointer")
	}
	*a = strings.Split(string(text), ",")
	return nil
}

type structWithNestedField struct {
	StrField string `env:"STRUCT_NESTED_STR_FIELD"`
}

type structWithTextUnmarshaller struct {
	str string
}

func (s *structWithTextUnmarshaller) UnmarshalText(text []byte) error {
	if s == nil {
		return fmt.Errorf("unmarshal text: nil pointer")
	}
	s.str = string(text)
	return nil
}

type structErrUnmarshaller struct{}

func (s *structErrUnmarshaller) UnmarshalText(text []byte) error {
	return errors.New("test_err")
}

type testConfigInvalidArray struct {
	Arr []string `env:"ARR"`
}

type testConfigUnsupportedKind struct {
	Num complex64 `env:"COMPLEX"`
}

type testConfigInvalidDive struct {
	Info string `env:",dive"`
}

type testConfigUnmarshalErr struct {
	StructErrUnmarshaler *structErrUnmarshaller `env:"STRUCT_ERR_UNMARSHALER"`
}

func TestApplyWithPrefix(t *testing.T) {
	type args struct {
		target any
		prefix string
	}
	tests := []struct {
		name    string
		args    args
		testFn  func(*testing.T, args)
		envVars map[string]string
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				target: &testConfig{},
				prefix: "TEST",
			},
			testFn: func(t *testing.T, a args) {
				tc := a.target.(*testConfig)
				em := enumWithTextUnmarshaller(2)

				assert.Equal(t, "test", tc.StrField)
				assert.Equal(t, int(1), tc.IntField)
				assert.Equal(t, true, tc.BoolField)
				assert.Equal(t, uint(1), tc.UintField)
				assert.Equal(t, float64(1), tc.FloatField)
				assert.Equal(t, &arrayWithTextUnmarshaller{"a", "b", "c"}, tc.ArrayField)
				assert.Equal(t, em, *tc.TypedField)
				assert.Equal(t, em, tc.TypedFieldNoPtr)
				assert.Equal(t, structWithTextUnmarshaller{"test"}, *tc.StructField)
				assert.Equal(t, structWithTextUnmarshaller{"test"}, tc.StructFieldNoPtr)
				assert.Equal(t, "test", tc.StructNestedFieldWithDive.StrField)
				assert.Equal(t, "test", tc.StructNestedFieldWithDiveNoPtr.StrField)
				// left untouched
				assert.Nil(t, tc.StructNestedField)
			},
			envVars: map[string]string{
				"TEST_PTR_FIELD":               "ptr",
				"TEST_STR_FIELD":               "test",
				"TEST_INT_FIELD":               "1",
				"TEST_BOOL_FIELD":              "true",
				"TEST_UIT_FIELD":               "1",
				"TEST_FLOAT_FIELD":             "1.0",
				"TEST_ARRAY_FIELD":             "a,b,c",
				"TEST_TYPED_FIELD":             "c",
				"TEST_STRUCT_FIELD":            "test",
				"TEST_STRUCT_NESTED_STR_FIELD": "test",
			},
		},
		{
			name: "invalid_array",
			args: args{
				target: &testConfigInvalidArray{},
				prefix: "TEST",
			},
			testFn: func(t *testing.T, a args) {
			},
			envVars: map[string]string{
				"TEST_ARR": `[test]`,
			},
			wantErr: true,
		},
		{
			name: "invalid_kind",
			args: args{
				target: &testConfigUnsupportedKind{},
				prefix: "TEST",
			},
			testFn: func(t *testing.T, a args) {
			},
			envVars: map[string]string{
				"TEST_COMPLEX": `42`,
			},
			wantErr: true,
		},
		{
			name: "invalid_dive",
			args: args{
				target: &testConfigInvalidDive{},
				prefix: "TEST",
			},
			testFn: func(t *testing.T, a args) {
			},
			envVars: map[string]string{},
			wantErr: true,
		},
		{
			name: "invalid_config_type",
			args: args{
				target: testConfigInvalidDive{},
				prefix: "TEST",
			},
			testFn: func(t *testing.T, a args) {
			},
			envVars: map[string]string{},
			wantErr: true,
		},
		{
			name: "invalid_not_config_type",
			args: args{
				target: intPtr(1),
				prefix: "TEST",
			},
			testFn: func(t *testing.T, a args) {
			},
			envVars: map[string]string{},
			wantErr: true,
		},
		{
			name: "invalid_unmarshal",
			args: args{
				target: &testConfigUnmarshalErr{},
				prefix: "TEST",
			},
			testFn: func(t *testing.T, a args) {
			},
			envVars: map[string]string{
				"TEST_STRUCT_ERR_UNMARSHALER": "err",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		// Prepare environment for the test.
		for k, v := range tt.envVars {
			require.NoError(t, os.Setenv(k, v))
		}

		t.Run(tt.name, func(t *testing.T) {
			if err := ApplyWithPrefix(tt.args.target, tt.args.prefix); (err != nil) != tt.wantErr {
				t.Errorf("ApplyWithPrefix() error = %v, wantErr %v", err, tt.wantErr)
				t.FailNow()
			}
			tt.testFn(t, tt.args)
		})
	}
}
