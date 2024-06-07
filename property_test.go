package parseme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsBoolean(t *testing.T) {
	testcases := []struct {
		name     string
		property *Property
		expected bool
	}{
		{
			"boolean property",
			NewProperty(Boolean, "random", "true"),
			true,
		},
		{
			"non-boolean property",
			NewProperty(Value, "lang", "en-US"),
			false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			result := tc.property.IsBoolean()
			assert.Equal(result, tc.expected)
		})
	}
}

func Test_Name(t *testing.T) {
	testcases := []struct {
		name     string
		property *Property
		expected string
	}{
		{
			"boolean property",
			NewProperty(Boolean, "random", "true"),
			"random",
		},
		{
			"non-boolean property",
			NewProperty(Value, "lang", "en-US"),
			"lang",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			result := tc.property.Name()
			assert.Equal(result, tc.expected)
		})
	}
}

func Test_Value(t *testing.T) {
	testcases := []struct {
		name     string
		property *Property
		expected string
	}{
		{
			"boolean property",
			NewProperty(Boolean, "random", "true"),
			"true",
		},
		{
			"non-boolean property",
			NewProperty(Value, "lang", "en-US"),
			"en-US",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			result := tc.property.Value()
			assert.Equal(result, tc.expected)
		})
	}
}

func Test_BooleanValue(t *testing.T) {
	testcases := []struct {
		name        string
		property    *Property
		expected    bool
		expectedErr error
	}{
		{
			"boolean true property",
			NewProperty(Boolean, "random", "true"),
			true,
			nil,
		},
		{
			"boolean false property",
			NewProperty(Boolean, "random", "false"),
			false,
			nil,
		},

		{
			"non-boolean property",
			NewProperty(Value, "lang", "en-US"),
			false,
			&PropertyBooleanValueError{},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			result, err := tc.property.BooleanValue()

			if tc.expectedErr != nil {
				assert.Equal(err, tc.expectedErr)
				assert.EqualError(err, tc.expectedErr.Error())
			} else if err != nil && tc.expectedErr == nil {
				t.Log("BooleanValue() returned an error, but expected error is nil.")
				t.FailNow()
			} else {
				assert.Equal(result, tc.expected)
			}
		})
	}
}

func Test_SetType(t *testing.T) {
	testcases := []struct {
		name     string
		property *Property
		expected *Property
	}{
		{
			"set value property to boolean",
			NewProperty(Value, "random", "true"),
			NewProperty(Boolean, "random", "true"),
		},
		{
			"set boolean property to value",
			NewProperty(Boolean, "random", "false"),
			NewProperty(Value, "random", "false"),
		},
		{
			"set to boolean with value different from true",
			NewProperty(Value, "lang", "en-US"),
			NewProperty(Boolean, "lang", "true"),
		},
		{
			"set to boolean with value equal to true",
			NewProperty(Value, "lang", "true"),
			NewProperty(Boolean, "lang", "true"),
		},
		{
			"set to boolean with value equal to false",
			NewProperty(Value, "lang", "false"),
			NewProperty(Boolean, "lang", "false"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			property := tc.property

			if tc.property.propertyType == Value {
				property.SetType(Boolean)
			} else {
				property.SetType(Value)
			}
			assert.Equal(property, tc.expected)
		})
	}
}

func Test_SetName(t *testing.T) {
	testcases := []struct {
		name         string
		propertyName string
		property     *Property
		expected     *Property
		expectedErr  error
	}{
		{
			"value property valid name",
			"country",
			NewProperty(Value, "lang", "en-US"),
			NewProperty(Value, "country", "en-US"),
			nil,
		},
		{
			"boolean property valid name",
			"country",
			NewProperty(Boolean, "activated", "true"),
			NewProperty(Boolean, "country", "true"),
			nil,
		},
		{
			"name with leading and trailing spaces",
			"     country      ",
			NewProperty(Boolean, "activated", "true"),
			NewProperty(Boolean, "country", "true"),
			nil,
		},
		{
			"empty value",
			"",
			NewProperty(Boolean, "activated", "true"),
			nil,
			&PropertyEmptyNameError{},
		},
		{
			"value property invalid name",
			"1nvalid",
			NewProperty(Value, "lang", "en-US"),
			nil,
			&PropertyInvalidNameError{name: "1nvalid"},
		},
		{
			"boolean property invalid name",
			"1nvalid",
			NewProperty(Boolean, "random", "false"),
			nil,
			&PropertyInvalidNameError{name: "1nvalid"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			property := tc.property
			err := property.SetName(tc.propertyName)

			if tc.expectedErr != nil {
				assert.Equal(err, tc.expectedErr)
				assert.EqualError(err, tc.expectedErr.Error())
			} else if err != nil && tc.expectedErr == nil {
				t.Log("SetName() returned an error, but expected error is nil")
				t.FailNow()
			} else {
				assert.Equal(property, tc.expected)
			}
		})
	}
}

func Test_SetValue(t *testing.T) {
	testcases := []struct {
		name          string
		propertyValue string
		property      *Property
		expected      *Property
		expectedErr   error
	}{
		{
			"property value",
			"pt-BR",
			NewProperty(Value, "lang", "en-US"),
			NewProperty(Value, "lang", "pt-BR"),
			nil,
		},
		{
			"assign true to boolean property",
			"true",
			NewProperty(Boolean, "activated", "false"),
			NewProperty(Boolean, "activated", "true"),
			nil,
		},
		{
			"assign false to boolean property",
			"false",
			NewProperty(Boolean, "activated", "true"),
			NewProperty(Boolean, "activated", "false"),
			nil,
		},
		{
			"assign unsupported value to boolean property",
			"invalid",
			NewProperty(Boolean, "activated", "true"),
			nil,
			&PropertyInvalidBooleanError{value: "invalid"},
		},
		{
			"name with leading and trailing spaces",
			"     country      ",
			NewProperty(Value, "type", "city"),
			NewProperty(Value, "type", "     country      "),
			nil,
		},
		{
			"empty value",
			"",
			NewProperty(Value, "lang", "en-US"),
			NewProperty(Value, "lang", ""),
			nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			property := tc.property
			err := property.SetValue(tc.propertyValue)

			if tc.expectedErr != nil {
				assert.Equal(err, tc.expectedErr)
				assert.EqualError(err, tc.expectedErr.Error())
			} else if err != nil && tc.expectedErr == nil {
				t.Log("SetName() returned an error, but expected error is nil")
				t.FailNow()
			} else {
				assert.Equal(property, tc.expected)
			}
		})
	}
}
