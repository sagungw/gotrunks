package errors

import (
	"errors"
	"fmt"
)

type Code string

type SubCode string

type CodedError struct {
	Code     Code
	SubCode  SubCode
	original error
}

func NewCodedError(code Code, subCode SubCode, original error) *CodedError {
	return &CodedError{
		Code:     code,
		SubCode:  subCode,
		original: original,
	}
}

func NewCodedErrorMessage(code Code, subCode SubCode, originalMessage string) *CodedError {
	return &CodedError{
		Code:     code,
		SubCode:  subCode,
		original: errors.New(originalMessage),
	}
}

func (e *CodedError) Error() string {
	return fmt.Sprintf("[%s][%s]-%s", e.Code, e.SubCode, e.original.Error())
}

func (e *CodedError) Unwrap() error {
	return e.original
}
