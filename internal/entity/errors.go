package entity

import "errors"

type ResponseError struct {
	ErrMsg string `json:"error"`
}

const (
	//=======================================
	MsgErrorQuery = "invalid query parametr's"
	MsgErrorBody  = "invalid body request"
	MsgErrorPath  = "invalid path parametr"
)

var (
	ErrorsNotBody  = errors.New(MsgErrorBody)
	ErrorsGetPath  = errors.New("GetValueFromUrl: invalid get path")
	ErrorsNotFound = errors.New("Not found id's")
)
