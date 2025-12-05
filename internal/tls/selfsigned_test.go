package tls

import (
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"
)

func TestGenerateSelfSignedCert(t *testing.T) {
	certPEM, keyPEM, err := GenerateSelfSignedCert()
	if err != nil {
		t.Fatalf("Failed to generate certificate: %v", err)
	}

	// Verify certificate is valid PEM
	certBlock, _ := pem.Decode(certPEM)
	if certBlock == nil {
		t.Fatal("Failed to decode certificate PEM")
	}
	if certBlock.Type != "CERTIFICATE" {
		t.Errorf("Expected CERTIFICATE block, got %s", certBlock.Type)
	}

	// Verify key is valid PEM
	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil {
		t.Fatal("Failed to decode key PEM")
	}
	if keyBlock.Type != "RSA PRIVATE KEY" {
		t.Errorf("Expected RSA PRIVATE KEY block, got %s", keyBlock.Type)
	}

	// Parse and verify certificate
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		t.Fatalf("Failed to parse certificate: %v", err)
	}

	// Check subject
	if cert.Subject.Organization[0] != "Forge Dev" {
		t.Errorf("Expected Organization 'Forge Dev', got %s", cert.Subject.Organization[0])
	}

	// Check DNS names
	foundLocalhost := false
	for _, name := range cert.DNSNames {
		if name == "localhost" {
			foundLocalhost = true
			break
		}
	}
	if !foundLocalhost {
		t.Error("Certificate should include 'localhost' in DNS names")
	}

	// Check IP addresses
	foundLoopback := false
	for _, ip := range cert.IPAddresses {
		if ip.String() == "127.0.0.1" {
			foundLoopback = true
			break
		}
	}
	if !foundLoopback {
		t.Error("Certificate should include 127.0.0.1 in IP addresses")
	}

	// Check validity period (should be about 1 year)
	validityPeriod := cert.NotAfter.Sub(cert.NotBefore)
	expectedPeriod := 365 * 24 * time.Hour
	if validityPeriod < expectedPeriod-time.Hour || validityPeriod > expectedPeriod+time.Hour {
		t.Errorf("Expected validity period of ~1 year, got %v", validityPeriod)
	}

	// Verify not expired
	if time.Now().After(cert.NotAfter) {
		t.Error("Certificate is already expired")
	}
	if time.Now().Before(cert.NotBefore) {
		t.Error("Certificate is not yet valid")
	}
}

func TestLoadTLSConfig(t *testing.T) {
	certPEM, keyPEM, err := GenerateSelfSignedCert()
	if err != nil {
		t.Fatalf("Failed to generate certificate: %v", err)
	}

	tlsConfig, err := LoadTLSConfig(certPEM, keyPEM)
	if err != nil {
		t.Fatalf("Failed to load TLS config: %v", err)
	}

	if len(tlsConfig.Certificates) != 1 {
		t.Errorf("Expected 1 certificate, got %d", len(tlsConfig.Certificates))
	}

	if tlsConfig.MinVersion != 0x0303 { // TLS 1.2
		t.Errorf("Expected TLS 1.2 minimum, got %x", tlsConfig.MinVersion)
	}
}

func TestLoadTLSConfig_InvalidCert(t *testing.T) {
	_, err := LoadTLSConfig([]byte("invalid"), []byte("invalid"))
	if err == nil {
		t.Error("Expected error for invalid certificate")
	}
}
