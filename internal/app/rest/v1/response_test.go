package v1

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseDto_NewResponseDto(t *testing.T) {

	t.Run("success response with empty data", func(t *testing.T) {
		resp := NewResponseDto(MSG_SUCCESS, nil, "department")
		respByte, err := json.Marshal(resp)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"message":"success","data":{"department":{}},"version":"-"}`, string(respByte))
	})

	t.Run("success response with data", func(t *testing.T) {
		resp := NewResponseDto(MSG_SUCCESS, map[string]int{"id": 1, "age": 2}, "department")
		respByte, err := json.Marshal(resp)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"message":"success","data":{"department":{"id":1,"age":2}},"version":"-"}`, string(respByte))
	})
}

func TestResponseDto_NewResponsesDto(t *testing.T) {

	t.Run("success response with empty data", func(t *testing.T) {
		resp := NewResponsesDto[int](MSG_SUCCESS, nil, "departments")
		respByte, err := json.Marshal(resp)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"message":"success","data":{"departments":[]},"version":"-"}`, string(respByte))
	})

	t.Run("success response with data", func(t *testing.T) {
		resp := NewResponsesDto[int](MSG_SUCCESS, []int{1, 2}, "departments")
		respByte, err := json.Marshal(resp)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"message":"success","data":{"departments": [1,2]},"version":"-"}`, string(respByte))
	})
}
