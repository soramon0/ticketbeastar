package utils

import (
	"reflect"
	"strings"
	"ticketbeastar/pkg/models"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type ValidatorTransaltor struct {
	Validator    *validator.Validate
	Translations ut.Translator
}

func NewValidator() (*ValidatorTransaltor, error) {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	validate := validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)

	if err := registerOverrides(trans, validate); err != nil {
		return nil, err
	}

	return &ValidatorTransaltor{
		Validator:    validate,
		Translations: trans,
	}, nil
}

func (vt *ValidatorTransaltor) ValidationErrors(ve validator.ValidationErrors) *models.APIValidaitonErrors {
	out := make([]models.APIFieldError, len(ve))
	for i, fe := range ve {
		t := ve.Translate(vt.Translations)
		out[i] = models.APIFieldError{Field: fe.Field(), Message: t[fe.Namespace()]}
	}

	return &models.APIValidaitonErrors{Errors: out}
}

func registerOverrides(trans ut.Translator, v *validator.Validate) error {
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is required", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())

		return t
	})
}
