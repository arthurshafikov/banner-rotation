package http

import (
	"errors"
	"strconv"
)

type RequestParser struct{}

func NewRequestParser() *RequestParser {
	return &RequestParser{}
}

func (rp *RequestParser) ParseInt64FromInterface(value interface{}) (int64, error) {
	valueString, ok := value.(string)
	if !ok {
		return 0, errors.New("could not convert value to string")
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
