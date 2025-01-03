// security/crypto.go
package security

import (
    "crypto/aes"
    "crypto/rsa"
    "crypto/rand"
    "encoding/base64"
)

// Crypto 加密服务
type Crypto struct {
    // 配置
    config     *CryptoConfig
    // 密钥管理
    keyManager *KeyManager
    // 算法提供者
    providers  map[string]CryptoProvider
    // 性能指标
    metrics    *CryptoMetrics
}

// CryptoConfig 加密配置
type CryptoConfig struct {
    // 默认算法
    DefaultAlgorithm string
    // 密钥轮换周期
    KeyRotationPeriod time.Duration
    // 密钥大小
    KeySize int
    // 加密选项
    Options CryptoOptions
}

// CryptoProvider 加密提供者接口
type CryptoProvider interface {
    // 加密
    Encrypt(plaintext []byte, key []byte) ([]byte, error)
    // 解密
    Decrypt(ciphertext []byte, key []byte) ([]byte, error)
    // 生成密钥
    GenerateKey(size int) ([]byte, error)
}

// Encrypt 加密数据
func (c *Crypto) Encrypt(ctx context.Context, data []byte, opts ...CryptoOption) ([]byte, error) {
    options := c.config.Options
    for _, opt := range opts {
        opt(&options)
    }

    // 获取加密提供者
    provider, err := c.getProvider(options.Algorithm)
    if err != nil {
        return nil, err
    }

    // 获取加密密钥
    key, err := c.keyManager.GetKey(ctx, options.KeyID)
    if err != nil {
        return nil, err
    }

    // 执行加密
    start := time.Now()
    ciphertext, err := provider.Encrypt(data, key)
    if err != nil {
        c.metrics.EncryptionErrors.Inc()
        return nil, err
    }
    c.metrics.EncryptionDuration.Observe(time.Since(start).Seconds())

    return ciphertext, nil
}
