package errx

import (
	"errors"
	"fmt"
)

type Error struct {
	error
	LogData      any
	ResponseData any
}

func New(text string) Error {
	return Error{
		error: errors.New(text),
	}
}

func Wrap(err error, text string) Error {
	return Error{
		error: fmt.Errorf("%s: %w", text, err),
	}
}

func (e Error) Unwrap() error {
	return errors.Unwrap(e.error)
}

func (e Error) Is(target error) bool {
	var targetErr Error
	if !errors.As(target, &targetErr) {
		return false
	}

	return e.error.Error() == targetErr.Error()
}

func (e Error) WithLogData(logData any) Error {
	e.LogData = logData
	return e
}

func (e Error) AddLogStr(key, value string) Error {
	if e.LogData == nil {
		e.LogData = map[string]string{key: value}
		return e
	}

	logDataMap, ok := e.LogData.(map[string]string)
	if !ok {
		return e
	}

	logDataMap[key] = value

	e.LogData = logDataMap

	return e
}

func (e Error) WithResponseData(responseData any) Error {
	e.ResponseData = responseData
	return e
}
