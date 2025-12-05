# Issue #037: Add HTTPS/TLS Support

**Priority:** 🟡 HIGH  
**Estimated Tokens:** ~1,200 (Low complexity)  
**Agent Role:** Implementation

---

## 1. 🎫 Related Issue Context

**Gap Reference:** GAP-004 from v1-analysis.md

Current server only runs HTTP:
```go
if err := http.ListenAndServe(":8080", mux); err != nil {
```

Per Project Charter: "HTTPS/WSS used locally."

---

## 2. 📋 Acceptance Criteria

### Backend (Go)
- [ ] Add `FORGE_TLS_CERT` environment variable for certificate path
- [ ] Add `FORGE_TLS_KEY` environment variable for private key path
- [ ] If both are set, use `http.ListenAndServeTLS()`
- [ ] If neither is set, fall back to HTTP with warning log
- [ ] Add `--dev-tls` flag to auto-generate self-signed cert for development

### Self-Signed Certificate Generation
- [ ] Generate cert/key pair in memory or `.forge/certs/` directory
- [ ] Cert valid for: localhost, 127.0.0.1
- [ ] 1 year expiration
- [ ] Log warning that self-signed cert is for development only

### Frontend (React)
- [ ] Update WebSocket hook to detect if running on HTTPS
- [ ] Use `wss://` protocol when on HTTPS, `ws://` when on HTTP

### Documentation
- [ ] Add TLS configuration section to README
- [ ] Document how to generate production certificates

---

## 3. 📊 Token Efficiency Strategy

- Changes primarily in main.go
- Use standard library crypto/tls for self-signed generation
- Minimal frontend changes (protocol detection)

---

## 4. 🏗️ Technical Specification

### Main.go TLS Logic
```go
func main() {
    // ... existing setup ...
    
    tlsCert := os.Getenv("FORGE_TLS_CERT")
    tlsKey := os.Getenv("FORGE_TLS_KEY")
    devTLS := flag.Bool("dev-tls", false, "Generate self-signed cert for development")
    flag.Parse()
    
    addr := ":8080"
    
    if tlsCert != "" && tlsKey != "" {
        log.Printf("Starting HTTPS server on %s", addr)
        if err := http.ListenAndServeTLS(addr, tlsCert, tlsKey, mux); err != nil {
            log.Fatal(err)
        }
    } else if *devTLS {
        cert, key := generateSelfSignedCert()
        log.Println("⚠️  Using self-signed certificate for development")
        // ... use cert/key ...
    } else {
        log.Println("⚠️  Running HTTP (no TLS) - not recommended for production")
        if err := http.ListenAndServe(addr, mux); err != nil {
            log.Fatal(err)
        }
    }
}
```

### Self-Signed Generation
```go
func generateSelfSignedCert() (certPEM, keyPEM []byte) {
    priv, _ := rsa.GenerateKey(rand.Reader, 2048)
    
    template := x509.Certificate{
        SerialNumber: big.NewInt(1),
        Subject:      pkix.Name{Organization: []string{"Forge Dev"}},
        NotBefore:    time.Now(),
        NotAfter:     time.Now().Add(365 * 24 * time.Hour),
        KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
        ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
        IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
        DNSNames:     []string{"localhost"},
    }
    
    certDER, _ := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
    // ... encode to PEM ...
}
```

### Frontend Protocol Detection
```typescript
const getWebSocketURL = () => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    return `${protocol}//${window.location.host}/ws`;
};
```

---

## 5. 📁 Files to Create/Modify

| Action | File |
|--------|------|
| MODIFY | `main.go` (add TLS support) |
| CREATE | `internal/tls/selfsigned.go` (cert generation) |
| MODIFY | `frontend/src/hooks/useWebSocket.ts` (protocol detection) |
| MODIFY | `README.md` (TLS documentation) |

---

## 6. ✅ Definition of Done

1. `FORGE_TLS_CERT` + `FORGE_TLS_KEY` enables HTTPS
2. `--dev-tls` flag generates and uses self-signed certificate
3. Server logs which mode it's running in (HTTP vs HTTPS)
4. Frontend WebSocket automatically uses correct protocol
5. README documents TLS configuration
