/**/
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"

	webhook "github.com/pandorasnox/kubernetes-default-container-resources/pkg"
)

func main() {
	tlsDisabled := flag.Bool("tlsDisabled", false, "(optional) disables tls for the server")
	limitMemory := flag.String("limitMemory", "1G", "memory limit (default 1G)")
	limitCPU := flag.String("limitCPU", "0.5", "cpu limit (default 0.5 cores)")
	requestMemory := flag.String("requestMemory", "1G", "memory request (default 1G)")
	requestCPU := flag.String("requestCPU", "0.1", "cpu request (default 0.1 cores)")
	flag.Parse()

	fmt.Println("tlsDisabled: ", *tlsDisabled)
	fmt.Println("")

	http.HandleFunc("/mutate", func(w http.ResponseWriter, r *http.Request) {
		err := webhook.ServeContent(w, r, *limitMemory, *limitCPU, *requestMemory, *requestCPU)
		if err != nil {
			log.Println(err)
		}
	})

	server := listenHTTPS
	if *tlsDisabled {
		server = listenHTTP
	}
	log.Fatal(server())
}

func listenHTTP() error {
	server := &http.Server{
		Addr: ":8083",
	}
	return server.ListenAndServe()
}

func listenHTTPS() error {
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
		Addr:      ":8083",
		TLSConfig: cfg,
	}
	return server.ListenAndServeTLS("/certs/ssl-cert.pem", "/certs/ssl-key.pem")
}
