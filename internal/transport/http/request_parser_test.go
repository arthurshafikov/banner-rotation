package http

import (
	"testing"

	"github.com/arthurshafikov/banner-rotation/internal/core"
	"github.com/stretchr/testify/require"
)

var expectedInt64 int64 = 64

func TestParseInt64FromInterface(t *testing.T) {
	rp := NewRequestParser()
	var value interface{} = "64"

	result, err := rp.ParseInt64FromInterface(value)

	require.NoError(t, err)
	require.Equal(t, expectedInt64, result)
}

func TestParseInt64FromInterfaceCouldntConvertToString(t *testing.T) {
	rp := NewRequestParser()
	var value interface{} = 312

	result, err := rp.ParseInt64FromInterface(value)

	require.ErrorIs(t, err, core.ErrConvertToString)
	require.Zero(t, result)
}

func TestParseInt64FromInterfaceStringIsNotNumeric(t *testing.T) {
	rp := NewRequestParser()
	var value interface{} = "Some string"

	result, err := rp.ParseInt64FromInterface(value)

	require.Contains(t, err.Error(), "strconv.Atoi: parsing \"Some string\": invalid syntax")
	require.Zero(t, result)
}

func TestParseInt64FromBytes(t *testing.T) {
	rp := NewRequestParser()
	value := []byte("64")

	result, err := rp.ParseInt64FromBytes(value)

	require.NoError(t, err)
	require.Equal(t, expectedInt64, result)
}

func TestParseInt64FromBytesStringIsNotNumeric(t *testing.T) {
	rp := NewRequestParser()
	value := []byte("Some string")

	result, err := rp.ParseInt64FromBytes(value)

	require.Contains(t, err.Error(), "strconv.Atoi: parsing \"Some string\": invalid syntax")
	require.Zero(t, result)
}
