package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// systemLogs business-level http error codes.
// the systemLogsNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	systemLogsNO       = 55
	systemLogsName     = "systemLogs"
	systemLogsBaseCode = errcode.HCode(systemLogsNO)

	ErrCreateSystemLogs     = errcode.NewError(systemLogsBaseCode+1, "failed to create "+systemLogsName)
	ErrDeleteByIDSystemLogs = errcode.NewError(systemLogsBaseCode+2, "failed to delete "+systemLogsName)
	ErrUpdateByIDSystemLogs = errcode.NewError(systemLogsBaseCode+3, "failed to update "+systemLogsName)
	ErrGetByIDSystemLogs    = errcode.NewError(systemLogsBaseCode+4, "failed to get "+systemLogsName+" details")
	ErrListSystemLogs       = errcode.NewError(systemLogsBaseCode+5, "failed to list of "+systemLogsName)

	// error codes are globally unique, adding 1 to the previous error code
)
