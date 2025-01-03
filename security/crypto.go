// security/crypto.go
package security

type Crypto struct {
    cipher    Cipher
    hash      Hash
    signer    Signer
    verifier  Verifier
}

// security/access.go
type AccessControl struct {
    policies  []Policy
    enforcer  *Enforcer
    auditor   *Auditor
    cache     *cache.Cache
}
