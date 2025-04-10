package model

import "encoding/json"

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func ParseErrorMessage(body []byte) string {
	errMsg := body
	var errRes ErrorResponse
	if err := json.Unmarshal(errMsg, &errRes); err == nil {
		errMsg = []byte(errRes.Error.Message)
	}
	return string(errMsg)
}
