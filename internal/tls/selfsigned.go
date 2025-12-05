package tls

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

// GenerateSelfSignedCert creates a self-signed certificate for development use.
// Returns the certificate and key as PEM-encoded bytes.
func GenerateSelfSignedCert() (certPEM, keyPEM []byte, err error) {
	// Generate RSA private key
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	// Create certificate template
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate serial number: %w", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Forge Dev"},
			CommonName:   "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour), // 1 year validity
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
		DNSNames:              []string{"localhost"},
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	// Encode certificate to PEM
	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	// Encode private key to PEM
	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return certPEM, keyPEM, nil
}

// GenerateAndSaveCert generates a self-signed certificate and saves it to the .forge/certs directory.
// Returns the paths to the certificate and key files.
func GenerateAndSaveCert() (certPath, keyPath string, err error) {
	// Create .forge/certs directory
	certsDir := ".forge/certs"
	if err := os.MkdirAll(certsDir, 0700); err != nil {
		return "", "", fmt.Errorf("failed to create certs directory: %w", err)
	}

	certPath = filepath.Join(certsDir, "localhost.crt")
	keyPath = filepath.Join(certsDir, "localhost.key")

	// Check if certs already exist
	if _, err := os.Stat(certPath); err == nil {
		if _, err := os.Stat(keyPath); err == nil {
			log.Println("Using existing self-signed certificate from", certsDir)
			return certPath, keyPath, nil
		}
	}

	// Generate new certificate
	certPEM, keyPEM, err := GenerateSelfSignedCert()
	if err != nil {
		return "", "", err
	}

	// Save certificate
	if err := os.WriteFile(certPath, certPEM, 0644); err != nil {
		return "", "", fmt.Errorf("failed to write certificate: %w", err)
	}

	// Save private key (more restrictive permissions)
	if err := os.WriteFile(keyPath, keyPEM, 0600); err != nil {
		return "", "", fmt.Errorf("failed to write private key: %w", err)
	}

	log.Printf("Generated self-signed certificate in %s", certsDir)
	return certPath, keyPath, nil
}

// LoadTLSConfig creates a tls.Config from certificate and key PEM data.
func LoadTLSConfig(certPEM, keyPEM []byte) (*tls.Config, error) {
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to load key pair: %w", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}, nil
}
