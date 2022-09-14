package http

import (
	"strconv"

	"github.com/arthurshafikov/banner-rotation/internal/core"
)

type RequestParser struct{}

func NewRequestParser() *RequestParser {
	return &RequestParser{}
}

func (rp *RequestParser) ParseInt64FromInterface(value interface{}) (int64, error) {
	valueString, ok := value.(string)
	if !ok {
		return 0, core.ErrConvertToString
	}

	valueInt, err := strconv.Atoi(valueString)
	if err != nil {
		return 0, err
	}

	return int64(valueInt), nil
}

func (rp *RequestParser) ParseInt64FromBytes(bytes []byte) (int64, error) {
	valueString := string(bytes)

	valueInt, err := strconv.Atoi(valueString)
	if err != nil {
		return 0, err
	}

	return int64(valueInt), nil
}
