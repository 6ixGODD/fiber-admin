package validator

import (
	"sync"
	"time"

	"fiber-admin/internal/pkg/config"
	"github.com/go-playground/validator/v10"
)

var (
	validateInstance *validator.Validate
	once             sync.Once
)

func earlierThan(fl validator.FieldLevel) bool {
	endTimePtr := fl.Parent().FieldByName(fl.Param()).Interface().(*string)
	if endTimePtr == nil {
		return false
	}
	startTime, err := time.Parse(time.RFC3339, fl.Field().String())
	if err != nil {
		return false
	}
	endTime, err := time.Parse(time.RFC3339, *endTimePtr)
	if err != nil {
		return false
	}
	return startTime.Before(endTime)
}

func rfc3339(fl validator.FieldLevel) bool {
	_, err := time.Parse(time.RFC3339, fl.Field().String())
	return err == nil
}

func noticeType(fl validator.FieldLevel) bool {
	switch fl.Field().String() {
	case config.NoticeTypeUrgent, config.NoticeTypeNormal:
		return true
	default:
		return false
	}
}

func userRole(fl validator.FieldLevel) bool {
	switch fl.Field().String() {
	case config.UserRoleAdmin, config.UserRoleUser:
		return true
	default:
		return false
	}
}

func operationType(fl validator.FieldLevel) bool {
	switch fl.Field().String() {
	case config.OperationTypeCreate, config.OperationTypeUpdate, config.OperationTypeDelete:
		return true
	default:
		return false
	}
}

func entityType(fl validator.FieldLevel) bool {
	switch fl.Field().String() {
	case config.EntityTypeDocumentation, config.EntityTypeNotice, config.EntityTypeUser:
		return true
	default:
		return false
	}
}

func operationStatus(fl validator.FieldLevel) bool {
	switch fl.Field().String() {
	case config.OperationStatusSuccess, config.OperationStatusFailure:
		return true
	default:
		return false
	}
}

func NewValidator() (*validator.Validate, error) {
	var err error
	once.Do(
		func() {
			validate := validator.New()
			if err = validate.RegisterValidation("earlierThan", earlierThan); err != nil {
				return
			}
			if err = validate.RegisterValidation("rfc3339", rfc3339); err != nil {
				return
			}
			if err = validate.RegisterValidation("noticeType", noticeType); err != nil {
				return
			}
			if err = validate.RegisterValidation("userRole", userRole); err != nil {
				return
			}
			if err = validate.RegisterValidation("operationType", operationType); err != nil {
				return
			}
			if err = validate.RegisterValidation("entityType", entityType); err != nil {
				return
			}
			if err = validate.RegisterValidation("operationStatus", operationStatus); err != nil {
				return
			}
			validateInstance = validate
		},
	)
	return validateInstance, err
}
