package http

import (
	"errors"
	"strconv"

	"github.com/valyala/fasthttp"
)

type RequestParser struct{}

func NewRequestParser() *RequestParser {
	return &RequestParser{}
}

func (rp *RequestParser) ParseIdFromRequest(ctx *fasthttp.RequestCtx) (int64, error) {
	return rp.ParseInt64FromRequest(ctx, "id")
}

func (rp *RequestParser) ParseInt64FromRequest(ctx *fasthttp.RequestCtx, key string) (int64, error) {
	valueString, ok := ctx.UserValue(key).(string)
	if !ok {
		return 0, errors.New("could not convert value to string")
	}

	valueInt, err := strconv.Atoi(valueString)
	if err != nil {
		return 0, err
	}

	return int64(valueInt), nil
}

func (rp *RequestParser) ParseInt64FromQueryArgs(ctx *fasthttp.RequestCtx, key string) (int64, error) {
	valueString := string(ctx.QueryArgs().Peek(key))

	valueInt, err := strconv.Atoi(valueString)
	if err != nil {
		return 0, err
	}

	return int64(valueInt), nil
}
