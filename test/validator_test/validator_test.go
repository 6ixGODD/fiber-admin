package validator_test

import (
	"testing"
	"time"

	"fiber-admin/test/mock"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FakeStruct struct {
	FakeID     *string `query:"fake_id" validate:"required,mongodb"`
	StartDate  *string `query:"start_date" validate:"omitnil,rfc3339,earlierThan=EndDate"`
	EndDate    *string `query:"end_date" validate:"omitnil,rfc3339"`
	FakeString *string `query:"fake_string" validate:"omitnil,max=1000,min=1"`
	FakeNumber *int64  `query:"fake_number" validate:"omitnil,numeric,min=1,max=100"`
	FakeBool   *bool   `query:"fake_bool" validate:"required"`
	FakeQuery  *string `query:"fake_query" validate:"searchQuery"`
	FakeEmail  *string `query:"fake_email" validate:"email"`
}

func TestValidator(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation(
		"searchQuery", func(fl validator.FieldLevel) bool {
			return len(fl.Field().String()) >= 3
		},
	)
	assert.Nil(t, err)

	err = validate.RegisterValidation(
		"earlierThan", func(fl validator.FieldLevel) bool {
			startStr := fl.Field().String()
			params := fl.Param()
			endStrPtr := fl.Parent().FieldByName(params).Interface().(*string)
			if endStrPtr == nil {
				return false
			}
			endStr := *endStrPtr
			start, err := time.Parse(time.RFC3339, startStr)
			if err != nil {
				return false
			}
			end, err := time.Parse(time.RFC3339, endStr)
			if err != nil {
				return false
			}
			return start.Before(end)
		},
	)
	assert.Nil(t, err)

	err = validate.RegisterValidation(
		"rfc3339", func(fl validator.FieldLevel) bool {
			_, err := time.Parse(time.RFC3339, fl.Field().String())
			return err == nil
		},
	)

	foo := &FakeStruct{
		FakeID:     nil,
		StartDate:  nil,
		EndDate:    nil,
		FakeString: nil,
		FakeNumber: nil,
		FakeBool:   nil,
		FakeQuery:  nil,
		FakeEmail:  nil,
	}
	err = validate.Struct(foo)
	assert.NotNil(t, err)
	assert.Equal(t, 4, len(err.(validator.ValidationErrors)))
	t.Log(err)

	var (
		wrongObjectId  = mock.RandomString(24)
		wrongStartDate = time.Now().AddDate(0, 0, 1).Format(time.RFC3339)
		wrongEndDate   = time.Now().AddDate(0, 0, -1).Format(time.RFC3339)
		wrongString    = mock.RandomString(1001)
		wrongNumber    = int64(0)
		fakeBool       = true
		wrongQuery     = "a"
		wrongEmail     = "a"
	)
	bar := &FakeStruct{
		FakeID:     &wrongObjectId,
		StartDate:  &wrongStartDate,
		EndDate:    &wrongEndDate,
		FakeString: &wrongString,
		FakeNumber: &wrongNumber,
		FakeBool:   &fakeBool,
		FakeQuery:  &wrongQuery,
		FakeEmail:  &wrongEmail,
	}
	err = validate.Struct(bar)
	assert.NotNil(t, err)
	assert.Equal(t, 6, len(err.(validator.ValidationErrors)))
	t.Log(err)

	var (
		rightObjectId  = primitive.NewObjectID().Hex()
		rightStartDate = time.Now().AddDate(0, 0, -1).Format(time.RFC3339)
		rightEndDate   = time.Now().AddDate(0, 0, 1).Format(time.RFC3339)
		rightString    = mock.RandomString(1000)
		rightNumber    = int64(1)
		rightBool      = true
		rightQuery     = "abc"
		rightEmail     = "123456789@qq.com"
	)
	baz := &FakeStruct{
		FakeID:     &rightObjectId,
		StartDate:  &rightStartDate,
		EndDate:    &rightEndDate,
		FakeString: &rightString,
		FakeNumber: &rightNumber,
		FakeBool:   &rightBool,
		FakeQuery:  &rightQuery,
		FakeEmail:  &rightEmail,
	}
	err = validate.Struct(baz)
	assert.Nil(t, err)
	t.Log(err)

}
