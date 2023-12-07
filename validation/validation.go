package validation

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// NewValidator returns a new validator instance along with
// registered English translator for the validation errors.
func NewValidator() (*validator.Validate, ut.Translator, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	english := en.New()
	uni := ut.New(english, english)

	trans, _ := uni.GetTranslator("en")

	err := en_translations.RegisterDefaultTranslations(validate, trans)

	if err != nil {
		return nil, nil, err
	}

	return validate, trans, nil
}
