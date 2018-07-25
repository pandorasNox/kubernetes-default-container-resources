/**/
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/pandorasnox/kubernetes-default-container-resources/pkg"
)

func main() {
	tlsDisabled := flag.Bool("tlsDisabled", false, "(optional) disables tls for the server")
	limitMemory := flag.String("limitMemory", "1G", "memory limit (default 1G)")
	limitCPU := flag.String("limitCPU", "0.5", "cpu limit (default 0.5 cores)")
	requestMemory := flag.String("requestMemory", "1G", "memory request (default 1G)")
	requestCPU := flag.String("requestCPU", "0.1", "cpu request (default 0.1 cores)")
	addr := flag.String("addr", ":8083", "address to bind to")
	sslCert := flag.String("sslCert", "/certs/ssl-cert.pem", "address to bind to")
	sslKey := flag.String("sslKey", "/certs/ssl-key.pem", "address to bind to")
	flag.Parse()

	fmt.Println("tlsDisabled: ", *tlsDisabled)
	fmt.Println("")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := webhook.ServeContent(w, r, *limitMemory, *limitCPU, *requestMemory, *requestCPU)
		if err != nil {
			log.Println(err)
		}
	})

	log.Fatal(server(*addr, *sslCert, *sslKey, *tlsDisabled))
}

func server(addr, sslCert, sslKey string, tlsDisabled bool) error {
	if tlsDisabled {
		return listenHTTP(addr)
	}

	return listenHTTPS(addr, sslCert, sslKey)
}

func listenHTTP(addr string) error {
	server := &http.Server{
		Addr: addr,
	}
	return server.ListenAndServe()
}

func listenHTTPS(addr, sslCert, sslKey string) error {
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	server := &http.Server{
		Addr:      addr,
		TLSConfig: cfg,
	}
	return server.ListenAndServeTLS(sslCert, sslKey)
}
