package contract

import (
	"context"
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"be/internal/bindcode"
)

// ErrTransactorNotInitialized 无私钥配置时交易方法返回的错误
var ErrTransactorNotInitialized = errors.New("vault transactor not initialized: strategistPrivateKey not configured")

// VaultContractClient 封装 PolyVault 合约的读写调用
// 只读方法始终可用，交易方法（WithdrawToStrategy / DepositFromStrategy）仅在配置私钥后可用
type VaultContractClient struct {
	rpcClient    *ethclient.Client
	vaultAddress common.Address
	vaultCaller  *bindcode.PolyVaultCaller
	transactor   *bindcode.PolyVaultTransactor // 有私钥时初始化，用于 write 操作
	auth         *bind.TransactOpts            // 交易签名者
}

// NewVaultContractClient 创建 VaultContractClient 实例
// rpcURL: 链上RPC节点地址, contractAddr: 已部署的PolyVault合约地址
// privateKey: 策略角色私钥（可选，为空时交易方法返回 ErrTransactorNotInitialized）
func NewVaultContractClient(rpcURL, contractAddr, privateKey string) (*VaultContractClient, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	addr := common.HexToAddress(contractAddr)
	caller, err := bindcode.NewPolyVaultCaller(addr, client)
	if err != nil {
		client.Close()
		return nil, err
	}

	c := &VaultContractClient{
		rpcClient:    client,
		vaultAddress: addr,
		vaultCaller:  caller,
	}

	// 可选：有私钥时初始化 transactor
	if privateKey != "" {
		// 去除 0x 前缀，HexToECDSA 不支持带前缀的 hex 字符串
		trimmedKey := strings.TrimPrefix(privateKey, "0x")
		trimmedKey = strings.TrimPrefix(trimmedKey, "0X")
		key, err := crypto.HexToECDSA(trimmedKey)
		if err != nil {
			client.Close()
			return nil, err
		}
		chainID, err := client.NetworkID(context.Background())
		if err != nil {
			client.Close()
			return nil, err
		}
		auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
		if err != nil {
			client.Close()
			return nil, err
		}
		// 设置 gas 限制和价格使用默认值（合约调用时自动估算）
		c.auth = auth

		transactor, err := bindcode.NewPolyVaultTransactor(addr, client)
		if err != nil {
			client.Close()
			return nil, err
		}
		c.transactor = transactor
	}

	return c, nil
}

// ========== 只读方法 ==========

// TotalAssets 调用 totalAssets()，返回金库总资产（USDC余额 + 策略债务）
func (c *VaultContractClient) TotalAssets(ctx context.Context) (*big.Int, error) {
	return c.vaultCaller.TotalAssets(&bind.CallOpts{Context: ctx})
}

// TotalSupply 调用 totalSupply()，返回总份额发行量
func (c *VaultContractClient) TotalSupply(ctx context.Context) (*big.Int, error) {
	return c.vaultCaller.TotalSupply(&bind.CallOpts{Context: ctx})
}

// AvailableBalance 调用 availableBalance()，返回金库中可用USDC余额
func (c *VaultContractClient) AvailableBalance(ctx context.Context) (*big.Int, error) {
	return c.vaultCaller.AvailableBalance(&bind.CallOpts{Context: ctx})
}

// StrategyDebt 调用 strategyDebt()，返回已部署到策略的资金量
func (c *VaultContractClient) StrategyDebt(ctx context.Context) (*big.Int, error) {
	return c.vaultCaller.StrategyDebt(&bind.CallOpts{Context: ctx})
}

// SharePrice 调用 convertToAssets(1e18)，返回 1 份额对应的 USDC 数量（6位精度）
func (c *VaultContractClient) SharePrice(ctx context.Context) (*big.Int, error) {
	return c.vaultCaller.ConvertToAssets(&bind.CallOpts{Context: ctx}, new(big.Int).SetUint64(1e18))
}

// ========== 交易方法（需私钥初始化） ==========

// WithdrawToStrategy 调用 withdrawToStrategy(amount)，从金库提取 USDC 到策略师地址
// 需要 STRATEGIST_ROLE，需要私钥已配置
func (c *VaultContractClient) WithdrawToStrategy(ctx context.Context, amount *big.Int) error {
	if c.transactor == nil || c.auth == nil {
		return ErrTransactorNotInitialized
	}
	tx, err := c.transactor.WithdrawToStrategy(c.auth, amount)
	if err != nil {
		return err
	}
	// 等待交易确认
	_, err = bind.WaitMined(ctx, c.rpcClient, tx)
	return err
}

// DepositFromStrategy 调用 depositFromStrategy(amount)，将 USDC 从策略师地址归还到金库
// 需要 STRATEGIST_ROLE，需要私钥已配置
func (c *VaultContractClient) DepositFromStrategy(ctx context.Context, amount *big.Int) error {
	if c.transactor == nil || c.auth == nil {
		return ErrTransactorNotInitialized
	}
	tx, err := c.transactor.DepositFromStrategy(c.auth, amount)
	if err != nil {
		return err
	}
	// 等待交易确认
	_, err = bind.WaitMined(ctx, c.rpcClient, tx)
	return err
}

// Close 关闭 RPC 连接
func (c *VaultContractClient) Close() {
	if c.rpcClient != nil {
		c.rpcClient.Close()
	}
}
