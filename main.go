/**/
package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/pandorasnox/kubernetes-default-container-resources/pkg"
	"github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	// log.SetLevel(log.WarnLevel)
}

func main() {
	tlsDisabled := flag.Bool("tlsDisabled", false, "disabled tls for the server")
	limitMemory := flag.String("limitMemory", "1G", "memory limit (default 1G)")
	limitCPU := flag.String("limitCPU", "0.5", "cpu limit (default 0.5 cores)")
	requestMemory := flag.String("requestMemory", "512M", "memory request (default 1G)")
	requestCPU := flag.String("requestCPU", "0.05", "cpu request (default 0.1 cores)")
	addr := flag.String("addr", ":8083", "address to bind to")
	sslCert := flag.String("sslCert", "/certs/ssl-cert.pem", "address to bind to")
	sslKey := flag.String("sslKey", "/certs/ssl-key.pem", "address to bind to")
	flag.Parse()

	logrus.WithFields(logrus.Fields{
		"tlsDisabled":   *tlsDisabled,
		"addr":          *addr,
		"limitMemory":   *limitMemory,
		"limitCPU":      *limitCPU,
		"requestMemory": *requestMemory,
		"requestCPU":    *requestCPU,
	}).Info("log programm flags")

	defaultResourceRequirements, err := webhook.ParseResourceRequirements(*limitMemory, *limitCPU, *requestMemory, *requestCPU)
	if err != nil {
		log.Fatalf("could not parse resource requirements based on program flags: %s", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := webhook.Mutate(w, r, defaultResourceRequirements)
		if err != nil {
			//todo: use "Fatalf" instead of "Printf"???
			log.Printf("mutation failed: %s", err)
		}
	})

	server := &http.Server{
		Addr: *addr,
	}
	if *tlsDisabled {
		log.Fatalf("server stop because: %s", server.ListenAndServe())
	}

	server.TLSConfig = &tls.Config{
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
	log.Fatalf("tls server stop because: %s", server.ListenAndServeTLS(*sslCert, *sslKey))
}
