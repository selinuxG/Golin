package web

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"golin/global"
	"math/big"
	"os"
	"time"
)

func CreateCert() {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	template := x509.Certificate{
		// 设置证书版本为 v3
		Version:      3,
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"高业尚"},
			CommonName:   "Golin安全加密",
		},
		Issuer: pkix.Name{
			CommonName: "高业尚",
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		panic(err)
	}

	if !global.PathExists("cert") {
		os.Mkdir("cert", os.FileMode(global.FilePer))
	}
	certOut, err := os.Create("cert/cert.pem")
	if err != nil {
		panic(err)
	}
	defer certOut.Close()

	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	keyOut, err := os.Create("cert/key.pem")
	if err != nil {
		panic(err)
	}
	defer keyOut.Close()

	pkcs1PrivateKey := x509.MarshalPKCS1PrivateKey(priv)
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: pkcs1PrivateKey})
}
