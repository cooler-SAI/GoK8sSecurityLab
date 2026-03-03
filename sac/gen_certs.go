package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

func main() {
	// 1. Generate private key (ECDSA)
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Printf("Failed to generate private key: %v\n", err)
		return
	}

	// 2. Configure certificate parameters
	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour) // Validity period 1 year

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, serialNumberLimit)

	// 3. Create template with SAN (Subject Alternative Name)
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"GoK8s Security Lab"},
			CommonName:   "sac-webhook-service.default.svc", // CN is deprecated, but included for compatibility
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true, // Self-signed cert is its own CA

		// 👇 IMPORTANT: Add Subject Alternative Names
		DNSNames: []string{
			"sac-webhook-service.default.svc", // K8s service FQDN
			"sac-webhook-service",             // Short service name
			"localhost",                       // For local testing
		},
		IPAddresses: []net.IP{
			net.ParseIP("127.0.0.1"), // For local testing
			net.ParseIP("0.0.0.0"),
		},
	}

	// 4. Create the certificate (self-signed)
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		fmt.Printf("Failed to create certificate: %v\n", err)
		return
	}

	// 5. Write cert.pem
	certOut, err := os.Create("cert.pem")
	if err != nil {
		fmt.Printf("Failed to create cert.pem: %v\n", err)
		return
	}
	defer func(certOut *os.File) {
		err := certOut.Close()
		if err != nil {
			fmt.Printf("Failed to close cert.pem: %v\n", err)
		}
	}(certOut)

	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		fmt.Printf("Failed to write cert.pem: %v\n", err)
		return
	}
	fmt.Println("✅ Created cert.pem with SAN extensions")

	// 6. Write key.pem
	keyOut, err := os.OpenFile("key.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Printf("Failed to create key.pem: %v\n", err)
		return
	}
	defer func(keyOut *os.File) {
		err := keyOut.Close()
		if err != nil {
			fmt.Printf("Failed to close key.pem: %v\n", err)
		}
	}(keyOut)

	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		fmt.Printf("Failed to marshal private key: %v\n", err)
		return
	}

	err = pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})
	if err != nil {
		fmt.Printf("Failed to write key.pem: %v\n", err)
		return
	}
	fmt.Println("✅ Created key.pem")

	// 7. Verify the certificate
	certPEM, err := os.ReadFile("cert.pem")
	if err == nil {
		block, _ := pem.Decode(certPEM)
		if block != nil {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err == nil {
				fmt.Println("\n📋 Certificate Information:")
				fmt.Printf("   Subject: %s\n", cert.Subject.CommonName)
				fmt.Printf("   DNS Names: %v\n", cert.DNSNames)
				fmt.Printf("   IP Addresses: %v\n", cert.IPAddresses)
				fmt.Printf("   Expires: %v\n", cert.NotAfter)
			}
		}
	}
}
