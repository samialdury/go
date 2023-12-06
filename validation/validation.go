package validation

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

func NewValidator() (*validator.Validate, ut.Translator) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	english := en.New()
	uni := ut.New(english, english)

	trans, _ := uni.GetTranslator("en")

	err := en_translations.RegisterDefaultTranslations(validate, trans)

	if err != nil {
		panic(err)
	}

	return validate, trans
}
