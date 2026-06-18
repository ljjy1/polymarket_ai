package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// predictions business-level http error codes.
// the predictionsNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	predictionsNO       = 7
	predictionsName     = "predictions"
	predictionsBaseCode = errcode.HCode(predictionsNO)

	ErrCreatePredictions     = errcode.NewError(predictionsBaseCode+1, "failed to create "+predictionsName)
	ErrDeleteByIDPredictions = errcode.NewError(predictionsBaseCode+2, "failed to delete "+predictionsName)
	ErrUpdateByIDPredictions = errcode.NewError(predictionsBaseCode+3, "failed to update "+predictionsName)
	ErrGetByIDPredictions    = errcode.NewError(predictionsBaseCode+4, "failed to get "+predictionsName+" details")
	ErrListPredictions       = errcode.NewError(predictionsBaseCode+5, "failed to list of "+predictionsName)

	// error codes are globally unique, adding 1 to the previous error code
)
