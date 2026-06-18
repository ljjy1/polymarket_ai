package ecode

import "github.com/go-dev-frame/sponge/pkg/errcode"

// auth business-level http error codes
// the auth code range is 20100~20199
var (
	ErrNonceExpired     = errcode.NewError(20100, "nonce已过期或不存在，请重新获取") // nonce expired or not found
	ErrSignatureInvalid = errcode.NewError(20101, "签名验证失败，请重试")         // signature verification failed
	ErrAddressMismatch  = errcode.NewError(20102, "钱包地址与签名不匹配")         // address does not match signature
	ErrAuthFailed       = errcode.NewError(20103, "认证失败，请重新登录")         // authentication failed
)
