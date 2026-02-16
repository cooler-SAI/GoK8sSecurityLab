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

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"GoK8s Security Lab"},
			CommonName:   "sac-webhook.default.svc", // Important for K8s
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// 3. Create the certificate itself (self-signed)
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		fmt.Printf("Failed to create certificate: %v\n", err)
		return
	}

	// 4. Write cert.pem
	certOut, _ := os.Create("cert.pem")
	err2 := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err2 != nil {
		fmt.Printf("Failed to write data to cert.pem: %v\n", err)
		err := certOut.Close()
		if err != nil {
			return
		}
		return
	}
	err3 := certOut.Close()
	if err3 != nil {
		fmt.Printf("Error closing cert.pem: %v\n", err3)
		return
	}
	fmt.Println("Created cert.pem")

	// 5. Write key.pem
	keyOut, _ := os.OpenFile("key.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	privBytes, _ := x509.MarshalECPrivateKey(priv)
	err4 := pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})
	if err4 != nil {
		fmt.Printf("Failed to write data to key.pem: %v\n", err4)
		return
	}
	err5 := keyOut.Close()
	if err5 != nil {
		fmt.Printf("Error closing key.pem: %v\n", err5)
		return
	}
	fmt.Println("Created key.pem")
}
