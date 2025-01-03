//security/auth.go
package security

import (
    "context"
    "time"
    
    "github.com/Corphon/daoframe/tools/cache"
)

// AuthService 认证服务
type AuthService struct {
    // 认证提供者
    providers  map[string]AuthProvider
    // 会话管理
    sessions   *SessionManager
    // 令牌管理
    tokens     *TokenManager
    // 密码管理
    passwords  *PasswordManager
    // 监控指标
    metrics    *AuthMetrics
}

// AuthProvider 认证提供者接口
type AuthProvider interface {
    Authenticate(ctx context.Context, credentials interface{}) (*Principal, error)
    Validate(ctx context.Context, token string) (*Principal, error)
    Revoke(ctx context.Context, token string) error
}

// TokenManager 令牌管理器
type TokenManager struct {
    store      TokenStore
    generator  TokenGenerator
    validator  TokenValidator
    cache      *cache.Cache
}

// Token 认证令牌
type Token struct {
    ID        string
    Type      string
    Principal *Principal
    Claims    map[string]interface{}
    IssuedAt  time.Time
    ExpireAt  time.Time
    Metadata  map[string]string
}

// Authenticate 认证用户
func (as *AuthService) Authenticate(ctx context.Context, providerID string, credentials interface{}) (*AuthResult, error) {
    // 获取认证提供者
    provider, exists := as.providers[providerID]
    if !exists {
        return nil, ErrProviderNotFound
    }

    // 执行认证
    start := time.Now()
    principal, err := provider.Authenticate(ctx, credentials)
    if err != nil {
        as.metrics.AuthenticationErrors.Inc()
        return nil, err
    }
    as.metrics.AuthenticationDuration.Observe(time.Since(start).Seconds())

    // 生成令牌
    token, err := as.tokens.Generate(ctx, principal)
    if err != nil {
        return nil, err
    }

    // 创建会话
    session, err := as.sessions.Create(ctx, principal, token)
    if err != nil {
        return nil, err
    }

    return &AuthResult{
        Principal: principal,
        Token:     token,
        Session:   session,
    }, nil
}
