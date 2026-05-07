package validator

import (
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Field   string `json:"champ"`
	Message string `json:"message"`
}

var (
	reAlphanumUnderscore = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	instance             = func() *validator.Validate {
		v := validator.New()
		v.RegisterValidation("alphanum_underscore", func(fl validator.FieldLevel) bool {
			return reAlphanumUnderscore.MatchString(fl.Field().String())
		})
		return v
	}()
)

/*
 * Validate vérifie la conformité d'un DTO par rapport à ses tags de validation.
 *
 * Attend  : un pointeur vers un struct annoté avec des tags `validate`.
 * Retourne: une liste d'erreurs par champ, ou nil si tout est valide.
 */

func Validate(dto any) []FieldError {
	err := instance.Struct(dto)
	if err == nil {
		return nil
	}

	var fieldErrors []FieldError
	for _, e := range err.(validator.ValidationErrors) {
		fieldErrors = append(fieldErrors, FieldError{
			Field:   e.Field(),
			Message: translateError(e),
		})
	}

	return fieldErrors
}

/*
 * BindAndValidate décode le corps JSON et valide le DTO en une seule opération.
 *
 * Attend  : le contexte Gin et un pointeur vers le DTO cible.
 * Retourne: une liste d'erreurs par champ, ou nil si tout est valide.
 */

func BindAndValidate(c *gin.Context, dto any) []FieldError {
	if err := c.ShouldBindJSON(dto); err != nil {
		return []FieldError{{Field: "body", Message: MsgInvalid}}
	}
	return Validate(dto)
}

func translateError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return MsgRequired
	case "email":
		return MsgInvalidEmail
	case "url":
		return MsgInvalidURL
	case "min":
		return fmt.Sprintf(MsgMinLength, e.Param())
	case "max":
		return fmt.Sprintf(MsgMaxLength, e.Param())
	case "gte":
		return fmt.Sprintf(MsgMinValue, e.Param())
	case "lte":
		return fmt.Sprintf(MsgMaxValue, e.Param())
	case "alphanum", "alphanum_underscore":
		return MsgAlphaNum
	default:
		return MsgInvalid
	}
}
