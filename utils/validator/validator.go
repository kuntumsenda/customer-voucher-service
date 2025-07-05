package validator

import (
	"errors"
	"reflect"
	"strconv"

	"customer-voucher-service/constants/message"

	"github.com/go-playground/validator/v10"
)

func ValidateReqField(req interface{}) error {
	v := validator.New()
	val := reflect.ValueOf(req)
	typ := reflect.TypeOf(req)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}
	if err := v.Struct(req); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			field, _ := typ.FieldByName(fieldErr.Field())
			label := field.Tag.Get("label")
			if label == "" {
				label = fieldErr.Field()
			}
			switch fieldErr.Tag() {
			case "required":
				return errors.New(message.RequiredMessage(label))
			case "email":
				return errors.New(message.EmailMessage(label))
			case "max":
				param := fieldErr.Param()
				return errors.New(message.MaxLengthMessage(label, toInt(param)))
				// Tambahkan case lain sesuai kebutuhan
			}
		}
		return errors.New("invalid request")
	}
	return nil
}

func toInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
