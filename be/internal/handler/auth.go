package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/go-dev-frame/sponge/pkg/gin/middleware/auth"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/sgorm"

	"be/internal/cache"
	"be/internal/config"
	"be/internal/database"
	"be/internal/ecode"
	"be/internal/model"
	"be/internal/types"
)

const (
	nonceTTL        = 5 * time.Minute // nonce 过期时间
	signMessageTmpl = "Sign in to Polymarket AI\nNonce: %s"
)

var (
	errSigInvalid   = errors.New("signature invalid")
	errAddrMismatch = errors.New("address mismatch")
)

// AuthHandler 定义鉴权相关的 handler 接口
type AuthHandler interface {
	Nonce(c *gin.Context)
	Login(c *gin.Context)
}

type authHandler struct {
	db         *sgorm.DB
	nonceCache cache.NonceCache
}

// NewAuthHandler 创建 AuthHandler
func NewAuthHandler() AuthHandler {
	return &authHandler{
		db:         database.GetDB(),
		nonceCache: cache.NewNonceCache(database.GetCacheType()),
	}
}

// Nonce 生成一个随机 nonce，用于 MetaMask 签名
// @Summary 获取签名 nonce
// @Description 生成一个随机 nonce，前端需在 MetaMask 中对返回的 message 进行签名
// @Tags auth
// @Accept json
// @Produce json
// @Param data body types.NonceRequest true "钱包地址"
// @Success 200 {object} types.Result{data=types.NonceReply}
// @Router /api/v1/auth/nonce [post]
func (h *authHandler) Nonce(c *gin.Context) {
	form := &types.NonceRequest{}
	if err := c.ShouldBindJSON(form); err != nil {
		logger.Warn("ShouldBindJSON error", logger.Err(err))
		response.Error(c, ecode.InvalidParams)
		return
	}

	// 校验地址格式
	address := strings.TrimSpace(form.Address)
	if !common.IsHexAddress(address) {
		response.Error(c, ecode.InvalidParams)
		return
	}

	// 生成 16 字节随机 nonce
	nonceBytes := make([]byte, 16)
	if _, err := rand.Read(nonceBytes); err != nil {
		logger.Error("rand.Read error", logger.Err(err))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}
	nonce := hex.EncodeToString(nonceBytes)

	// 存入缓存，TTL 5 分钟
	ctx := c.Request.Context()
	if err := h.nonceCache.Set(ctx, address, nonce, nonceTTL); err != nil {
		logger.Error("nonceCache.Set error", logger.Err(err))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	// 构造签名消息 - 前端需要在 MetaMask 中签名此消息
	message := buildSignMessage(nonce)

	response.Success(c, &types.NonceReply{
		Nonce:    nonce,
		Message:  message,
		ExpireIn: int(nonceTTL.Seconds()),
	})
}

// Login 验证 MetaMask 签名并返回 JWT token
// @Summary MetaMask 签名登录
// @Description 验证 MetaMask 签名，验证通过后返回 JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param data body types.LoginRequest true "地址和签名"
// @Success 200 {object} types.Result{data=types.LoginReply}
// @Router /api/v1/auth/login [post]
func (h *authHandler) Login(c *gin.Context) {
	form := &types.LoginRequest{}
	if err := c.ShouldBindJSON(form); err != nil {
		logger.Warn("ShouldBindJSON error", logger.Err(err))
		response.Error(c, ecode.InvalidParams)
		return
	}

	address := strings.TrimSpace(form.Address)

	// 校验地址格式
	if !common.IsHexAddress(address) {
		response.Error(c, ecode.InvalidParams)
		return
	}

	// 从缓存获取 nonce
	ctx := c.Request.Context()
	storedNonce, err := h.nonceCache.Get(ctx, address)
	if err != nil {
		if errors.Is(err, database.ErrCacheNotFound) {
			response.Error(c, ecode.ErrNonceExpired)
		} else {
			logger.Error("nonceCache.Get error", logger.Err(err))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	// 验证签名
	signMessage := buildSignMessage(storedNonce)
	if err := verifySignature(address, signMessage, form.Signature); err != nil {
		logger.Warn("verifySignature failed", logger.String("address", address), logger.Err(err))
		switch {
		case errors.Is(err, errSigInvalid):
			response.Error(c, ecode.ErrSignatureInvalid)
		case errors.Is(err, errAddrMismatch):
			response.Error(c, ecode.ErrAddressMismatch)
		default:
			response.Error(c, ecode.ErrAuthFailed)
		}
		return
	}

	// 验证通过，删除已使用的 nonce（防止重放攻击）
	go func() {
		_ = h.nonceCache.Del(ctx, address)
	}()

	// 查找或创建用户
	user, err := h.findOrCreateUser(ctx, address)
	if err != nil {
		logger.Error("findOrCreateUser error", logger.Err(err))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	// 生成 JWT token
	jwtCfg := config.Get().JWT
	token, err := auth.GenerateToken(
		address,
		auth.WithGenerateTokenFields(map[string]any{
			"user_id": user.ID,
			"address": address,
		}),
	)
	if err != nil {
		logger.Error("auth.GenerateToken error", logger.Err(err))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, &types.LoginReply{
		Token:         token,
		ExpireIn:      jwtCfg.ExpireTime * 3600,
		WalletAddress: address,
	})
}

// findOrCreateUser 查找或创建用户
func (h *authHandler) findOrCreateUser(ctx context.Context, address string) (*model.User, error) {
	user := &model.User{}
	err := h.db.WithContext(ctx).Where("wallet_address = ?", strings.ToLower(address)).First(user).Error
	if err == nil {
		// 更新最后登录时间
		now := time.Now()
		_ = h.db.WithContext(ctx).Model(user).Update("last_login_at", now).Error
		return user, nil
	}

	if !errors.Is(err, database.ErrRecordNotFound) {
		return nil, err
	}

	// 创建新用户
	newUser := &model.User{
		WalletAddress: strings.ToLower(address),
		LastLoginAt:   timePtr(time.Now()),
	}
	if err := h.db.WithContext(ctx).Create(newUser).Error; err != nil {
		return nil, err
	}
	return newUser, nil
}

// verifySignature 验证以太坊签名
func verifySignature(address, message, signatureHex string) error {
	// 构造以太坊签名消息 (加 \x19Ethereum Signed Message:\n 前缀)
	hash := accounts.TextHash([]byte(message))

	// 解码签名（处理 0x 前缀）
	sig := common.FromHex(signatureHex)
	if len(sig) != 65 {
		return errSigInvalid
	}

	// MetaMask 返回的 v 值是 27/28，需要转为 0/1
	if sig[64] == 27 || sig[64] == 28 {
		sig[64] -= 27
	}

	// 从签名中恢复公钥
	pubKey, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return errSigInvalid
	}

	// 从公钥获取以太坊地址
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	// 比对地址（大小写不敏感）
	if !strings.EqualFold(recoveredAddr.Hex(), address) {
		return errAddrMismatch
	}

	return nil
}

// buildSignMessage 构造签名消息
func buildSignMessage(nonce string) string {
	return fmt.Sprintf(signMessageTmpl, nonce)
}

// timePtr 返回 time.Time 的指针
func timePtr(t time.Time) *time.Time {
	return &t
}
