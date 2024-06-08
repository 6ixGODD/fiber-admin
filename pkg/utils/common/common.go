package common

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func GenerateUUID4() (string, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return newUUID.String(), nil
}

func EscapeSpecialChars(input string) string {
	re := regexp.MustCompile(`[\\^$*+?.()|[\]{}]`)
	return re.ReplaceAllStringFunc(
		input, func(s string) string {
			return "\\" + s
		},
	)
}

func FormatValidateError(errs error) error {
	errMessage := fmt.Errorf("failed to validate request")
	for _, err := range errs.(validator.ValidationErrors) {
		errMessage = fmt.Errorf(
			"%s[field: %s] not match the requirement: %s %s, actual: %s",
			errMessage, err.Field(), err.Tag(), err.Param(), err.Value(),
		)
	}
	return errMessage
}
