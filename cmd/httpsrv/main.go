package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
)

const autoTLSCertFile = ".cert.pem"
const autoTLSKeyFile = ".key.pem"

// Create temperary certificate for auto-tls.
func makecert() error {
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
			Organization: []string{"No Corp"},
		},
		DNSNames:  []string{"localhost"},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(3 * time.Hour),

		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
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
	if err := ioutil.WriteFile(autoTLSCertFile, certPEMBlock, 0644); err != nil {
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
	if err := ioutil.WriteFile(autoTLSKeyFile, keyPEMBlock, 0644); err != nil {
		return errors.Wrap(err, "Failed to write key PEM file")
	}
	return nil
}

func runServer(prog string, args []string) error {
	var addr string
	var insecure bool
	var verify bool
	var certFile string
	var keyFile string
	var caCertFile string
	var autoTLS bool

	fg := flag.NewFlagSet(prog, flag.ExitOnError)
	fg.StringVar(&addr, "addr", ":4433", "Server address")
	fg.BoolVar(&insecure, "insecure", false, "Insecure connection")
	fg.BoolVar(&verify, "verify", false, "Verify client certificate. Only meaningful when using https")
	fg.StringVar(&certFile, "cert", "", "Certificate file")
	fg.StringVar(&keyFile, "key", "", "Key file")
	fg.StringVar(&caCertFile, "cacert", "", "CA certificate")
	fg.BoolVar(&autoTLS, "auto-tls", false, "Auto TLS")
	fg.Parse(args)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("%v %v\n", req.Method, req.URL.Path)
		fmt.Fprintf(w, "Prouldly served with Go HTTP!\n")
	})
	tlsCfg := &tls.Config{}

	// Load CA certificate
	if len(caCertFile) != 0 {
		log.Printf("Using CA: %v\n", caCertFile)
		pemCert, err := ioutil.ReadFile(caCertFile)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Read cacert file fail: %v", caCertFile))
		}
		certPool := x509.NewCertPool()
		if ok := certPool.AppendCertsFromPEM(pemCert); !ok {
			return errors.Wrap(err, fmt.Sprintf("Unable to parse cacert file : %v", caCertFile))
		}
		tlsCfg.ClientCAs = certPool
	}

	// Verify client certificate
	if verify {
		tlsCfg.ClientAuth = tls.RequireAndVerifyClientCert
	}

	srv := &http.Server{
		Addr:      addr,
		Handler:   mux,
		TLSConfig: tlsCfg,
	}

	var err error
	if autoTLS {
		log.Println("AutoTLS is on")
		log.Printf("Creating temperary certificate %v, %v ...\n", autoTLSCertFile, autoTLSKeyFile)
		err := makecert()
		if err != nil {
			return errors.Wrap(err, "Make certification fail")
		}
		log.Printf("OK..")

		log.Printf("Listening on https://%v\n", addr)
		err = srv.ListenAndServeTLS(autoTLSCertFile, autoTLSKeyFile)
	} else if insecure == false {
		log.Printf("Listening on https://%v\n", addr)
		log.Printf("Client verify: %v\n", verify)

		// Server certificate and key
		if len(certFile) == 0 {
			return errors.New(fmt.Sprintf("-cert must be provided"))
		}
		if len(keyFile) == 0 {
			return errors.New(fmt.Sprintf("-key must be provided"))
		}

		err = srv.ListenAndServeTLS(certFile, keyFile)
	} else {
		log.Printf("Listening on http://%v\n", addr)
		err = srv.ListenAndServe()
	}

	if err != nil {
		return errors.Wrap(err, "Fail to start https server")
	}

	return nil
}

func runClient(prog string, args []string) error {
	var URL string
	var caCertFile string
	var insecure bool
	var certFile string
	var keyFile string

	fg := flag.NewFlagSet(prog, flag.ExitOnError)
	fg.StringVar(&URL, "url", "https://localhost:4433", "URL")
	fg.BoolVar(&insecure, "insecure", false, "Insecure mode, skip verifing server certificate. Only meaningful when using https")
	fg.StringVar(&caCertFile, "cacert", "", "CA certificate")
	fg.StringVar(&certFile, "cert", "", "Certificate file")
	fg.StringVar(&keyFile, "key", "", "Key file")
	fg.Parse(args)

	tlsCfg := &tls.Config{}

	// Load CA certificate
	if len(caCertFile) != 0 {
		log.Printf("Using CA: %v\n", caCertFile)
		pemCert, err := ioutil.ReadFile(caCertFile)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Read cacert file fail: %v", caCertFile))
		}
		certPool := x509.NewCertPool()
		if ok := certPool.AppendCertsFromPEM(pemCert); !ok {
			return errors.Wrap(err, fmt.Sprintf("Unable to parse cacert file : %v", caCertFile))
		}
		tlsCfg.RootCAs = certPool
	}

	// Insecure mode
	if insecure {
		tlsCfg.InsecureSkipVerify = true
	}

	if len(certFile) != 0 && len(keyFile) != 0 {
		log.Printf("Loading %v, %v\n", certFile, keyFile)
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Load certificate and key fail"))
		}
		tlsCfg.Certificates = make([]tls.Certificate, 1)
		tlsCfg.Certificates[0] = cert
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsCfg,
		},
	}

	log.Printf("Requesting %v\n", URL)
	log.Printf("Insecure(skip server certificate verify): %v\n", insecure)
	r, err := client.Get(URL)
	if err != nil {
		return errors.Wrap(err, "HTTP Get fail")
	}
	defer r.Body.Close()

	html, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(err, "Read HTTP request fail")
	}
	log.Printf("%v\n", r.Status)
	log.Println(string(html))
	return nil
}

func main() {
	prog := os.Args[0]
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %v [client|server]\n", prog)
		return
	}
	mode := os.Args[1]
	args := os.Args[2:]
	subProg := fmt.Sprintf("%v %v", prog, mode)

	var err error
	if mode == "server" {
		err = runServer(subProg, args)
	} else {
		err = runClient(subProg, args)
	}
	if err != nil {
		log.Println(err)
	}
}
