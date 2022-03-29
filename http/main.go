package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"
)

func runServer(prog string, args []string) error {
	var addr string
	var insecure bool
	var certFile string
	var keyFile string

	fg := flag.NewFlagSet(prog, flag.ExitOnError)
	fg.StringVar(&addr, "addr", ":4433", "Server address")
	fg.BoolVar(&insecure, "insecure", false, "Insecure connection")
	fg.StringVar(&certFile, "cert", "./server.pem", "Certificate file")
	fg.StringVar(&keyFile, "key", "./server-key.pem", "Key file")
	fg.Parse(args)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("%v %v\n", req.Method, req.URL.Path)
		fmt.Fprintf(w, "Prouldly served with Go and HTTPS!\n")
	})

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	var err error
	if insecure == false {
		log.Printf("Using certificate: %v, %v\n", certFile, keyFile)
		log.Printf("Listening on https://%v\n", addr)
		srv.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
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

	fg := flag.NewFlagSet(prog, flag.ExitOnError)
	fg.StringVar(&URL, "url", "https://localhost:4433", "HTTPS Server address")
	fg.StringVar(&caCertFile, "cacert", "cacert.pem", "CA certificate")
	fg.Parse(args)

	var tlsCfg *tls.Config
	if strings.ToLower(URL[:5]) == "https" {
		log.Printf("Using CA: %v\n", caCertFile)
		pemCert, err := ioutil.ReadFile(caCertFile)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Read cacert file fail: %v", caCertFile))
		}
		certPool := x509.NewCertPool()
		if ok := certPool.AppendCertsFromPEM(pemCert); !ok {
			return errors.Wrap(err, fmt.Sprintf("Unable to parse cacert file : %v", caCertFile))
		}
		tlsCfg = &tls.Config{
			RootCAs: certPool,
		}
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsCfg,
		},
	}

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
