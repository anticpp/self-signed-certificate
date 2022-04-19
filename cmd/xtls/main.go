package main

// xtls genca -c ca.yaml
// xtls gencert -c cert.yaml -cacert cacert.pem -key cakey.pem
//
//

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"time"

	"github.com/pkg/errors"
)

func makeca(outCertFile, outKeyFile string) error {
	// Create private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return errors.Wrap(err, "Failed to generate private key")
	}

	// Create certificate template
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return errors.Wrap(err, "Failed to serialNumber")
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: "test-ca",
			/*Organization:  []string{"Company, INC"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Franscisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},*/
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(3 * time.Hour),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// Create certificate data
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return errors.Wrap(err, "Failed to create certificate")
	}

	certPEMBlock := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	if certPEMBlock == nil {
		return errors.New("Failed to encode certificate to PEM")
	}
	if err := ioutil.WriteFile(outCertFile, certPEMBlock, 0644); err != nil {
		return errors.Wrap(err, "Failed to write certificate PEM file")
	}

	// Create private key data
	keyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return errors.Wrap(err, "Unable to marshal private key")
	}
	keyPEMBlock := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyBytes})
	if keyPEMBlock == nil {
		return errors.New("Failed to encode key to PEM")
	}
	if err := ioutil.WriteFile(outKeyFile, keyPEMBlock, 0644); err != nil {
		return errors.Wrap(err, "Failed to write key PEM file")
	}
	return nil
}

func loadca(certFile, keyFile string) (*x509.Certificate, any, error) {
	var cert *x509.Certificate
	var key any
	var err error
	var certBytes, keyBytes []byte
	var block *pem.Block

	// Load certificate
	certBytes, err = ioutil.ReadFile(certFile)
	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("Read cacert file fail: %v", certFile))
	}

	block, _ = pem.Decode(certBytes)
	if block == nil || block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
		return nil, nil, errors.New("cert pem Decode fail")
	}

	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, nil, errors.Wrap(err, "ParseCertificate fail")
	}

	// Load key
	keyBytes, err = ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("Read cacert file fail: %v", certFile))
	}

	block, _ = pem.Decode(keyBytes)
	if block == nil || block.Type != "PRIVATE KEY" || len(block.Headers) != 0 {
		return nil, nil, errors.New("key pem Decode fail")
	}

	key, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, nil, errors.Wrap(err, "ParsePKCS8PrivateKey fail")
	}
	return cert, key, nil
}

func makecert(
	CommonName string,
	outCertFile string,
	outKeyFile string,
	caCert *x509.Certificate,
	caKey any,
) error {
	// Create private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return errors.Wrap(err, "Failed to generate private key")
	}

	// Create certificate template
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return errors.Wrap(err, "Failed to serialNumber")
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: CommonName,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(3 * time.Hour),
	}

	// Create certificate data
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, caCert, &privateKey.PublicKey, caKey)
	if err != nil {
		return errors.Wrap(err, "Failed to create certificate")
	}

	certPEMBlock := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	if certPEMBlock == nil {
		return errors.New("Failed to encode certificate to PEM")
	}
	if err := ioutil.WriteFile(outCertFile, certPEMBlock, 0644); err != nil {
		return errors.Wrap(err, "Failed to write certificate PEM file")
	}

	// Create private key data
	keyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return errors.Wrap(err, "Unable to marshal private key")
	}
	keyPEMBlock := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyBytes})
	if keyPEMBlock == nil {
		return errors.New("Failed to encode key to PEM")
	}
	if err := ioutil.WriteFile(outKeyFile, keyPEMBlock, 0644); err != nil {
		return errors.Wrap(err, "Failed to write key PEM file")
	}
	return nil

}

func main() {
	// Create CA
	log.Println("making CA...")
	err := makeca("./cacert.pem", "./cakey.pem")
	log.Println("DONE")
	log.Println("")

	// Load CA
	log.Println("Loading CA...")
	caCert, caKey, err := loadca("./cacert.pem", "./cakey.pem")
	if err != nil {
		log.Fatalf("loadca fail: %v\n", err)
	}

	fmt.Println("====CERTIFICATE=======")
	fmt.Println("SerialNumber: ", caCert.SerialNumber)
	fmt.Println("Issuer:", caCert.Issuer)
	fmt.Println("Subject:", caCert.Subject)
	fmt.Println("NotBefore:", caCert.NotBefore)
	fmt.Println("NotAfter: ", caCert.NotAfter)

	fmt.Println("====PRIVATE KEY=======")
	switch caKey.(type) {
	case *rsa.PrivateKey:
		fmt.Println("RSA key")
	default:
		fmt.Println("Unknown private key")
	}
	log.Println("")

	// Sign certificate
	log.Println("Signing example.com...")
	err = makecert("example.com", "./example_cert.pem", "./example_key.pem", caCert, caKey)
	if err != nil {
		log.Fatalf("makecert fail: %v\n", err)
	}
	log.Println("DONE")
	log.Println("")
}
