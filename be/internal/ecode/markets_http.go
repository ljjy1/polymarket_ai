package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// markets business-level http error codes.
// the marketsNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	marketsNO       = 8
	marketsName     = "markets"
	marketsBaseCode = errcode.HCode(marketsNO)

	ErrCreateMarkets     = errcode.NewError(marketsBaseCode+1, "failed to create "+marketsName)
	ErrDeleteByIDMarkets = errcode.NewError(marketsBaseCode+2, "failed to delete "+marketsName)
	ErrUpdateByIDMarkets = errcode.NewError(marketsBaseCode+3, "failed to update "+marketsName)
	ErrGetByIDMarkets    = errcode.NewError(marketsBaseCode+4, "failed to get "+marketsName+" details")
	ErrListMarkets       = errcode.NewError(marketsBaseCode+5, "failed to list of "+marketsName)

	// error codes are globally unique, adding 1 to the previous error code
)
