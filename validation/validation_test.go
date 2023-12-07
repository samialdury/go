package validation

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestNewValidator(t *testing.T) {
	validate, translator, _ := NewValidator()

	assert.NotNil(t, validate, "Validator should not be nil")
	assert.NotNil(t, translator, "Translator should not be nil")
}

func TestStructValidation(t *testing.T) {
	validate, _, _ := NewValidator()

	type TestStruct struct {
		Field string `validate:"required"`
	}

	testStruct := TestStruct{Field: ""}
	err := validate.Struct(testStruct)
	assert.NotNil(t, err, "Expected validation error for empty field")
}

func TestTranslationFunctionality(t *testing.T) {
	validate, translator, _ := NewValidator()

	type TestStruct struct {
		Name  string `validate:"required"`
		Age   int    `validate:"gte=18"`
		Email string `validate:"required,email"`
	}

	testStruct := TestStruct{Name: "", Age: 17, Email: "not-an-email"}

	err := validate.Struct(testStruct)

	assert.NotNil(t, err, "Expected validation error")

	validationErrors, ok := err.(validator.ValidationErrors)
	assert.True(t, ok, "Errors should be of type ValidationErrors")

	for _, err := range validationErrors {
		translatedErr := err.Translate(translator)

		switch err.Field() {
		case "Name":
			assert.Equal(t, "Name is a required field", translatedErr, "Incorrect translation for Name field")
		case "Age":
			assert.Equal(t, "Age must be 18 or greater", translatedErr, "Incorrect translation for Age field")
		case "Email":
			assert.Equal(t, "Email must be a valid email address", translatedErr, "Incorrect translation for Email field")
		default:
			t.Errorf("Unexpected field: %s", err.Field())
		}
	}
}
