package types

// NonceRequest 获取 nonce 的请求参数
type NonceRequest struct {
	Address string `json:"address" binding:"required"` // MetaMask 钱包地址
}

// NonceReply 获取 nonce 的响应
type NonceReply struct {
	Nonce    string `json:"nonce"`    // 签名用的随机数
	Message  string `json:"message"`  // 需要在 MetaMask 中签名的完整消息
	ExpireIn int    `json:"expireIn"` // nonce 过期时间（秒）
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Address   string `json:"address" binding:"required"`   // MetaMask 钱包地址
	Signature string `json:"signature" binding:"required"` // 签名后的 hex 字符串
}

// LoginReply 登录响应
type LoginReply struct {
	Token         string `json:"token"`         // JWT token
	ExpireIn      int    `json:"expireIn"`      // token 过期时间（秒）
	WalletAddress string `json:"walletAddress"` // 钱包地址
}
