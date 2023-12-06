package config

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/go-playground/validator/v10"
	"github.com/samialdury/go/validation"
)

// Parse parses a struct containing `env` tags and loads its values from
// environment variables.
func Parse[T any](s *T) error {
	err := env.Parse(s)

	return err
}

// Validate validates a struct containing `env` tags.
func Validate[T any](s *T) error {
	v, t := validation.NewValidator()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("env"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	if err := v.Struct(s); err != nil {
		errs, ok := err.(validator.ValidationErrors)

		if !ok {
			return err
		}

		translated := errs.Translate(t)

		// Extract the env var name from the error message.
		re := regexp.MustCompile(`^\S+`)

		// Join all error messages into a single string.
		var sb strings.Builder
		for _, value := range translated {
			envVarName := re.FindString(value)
			sb.WriteString(fmt.Sprintf("env.%s. Got `%s`.\n", value, os.Getenv(envVarName)))
		}

		return fmt.Errorf(sb.String())
	}

	return nil
}
