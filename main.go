/**/
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/pandorasnox/kubernetes-default-container-resources/pkg"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func main() {
	tlsDisabled := flag.Bool("tlsDisabled", false, "disabled tls for the server")
	limitMemory := flag.String("limitMemory", "1G", "memory limit (default 1G)")
	limitCPU := flag.String("limitCPU", "0.5", "cpu limit (default 0.5 cores)")
	requestMemory := flag.String("requestMemory", "512M", "memory request (default 1G)")
	requestCPU := flag.String("requestCPU", "0.05", "cpu request (default 0.1 cores)")
	addr := flag.String("addr", ":8083", "address to bind to")
	sslCert := flag.String("sslCert", "/certs/ssl-cert.pem", "address to bind to")
	sslKey := flag.String("sslKey", "/certs/ssl-key.pem", "address to bind to")
	dryRun := flag.Bool("dry-run", false, "enables dry-run mode, always returning success AdmissionReview")
	flag.Parse()

	logrus.SetFormatter(&logrus.JSONFormatter{
		DataKey: "data",
		// PrettyPrint: true,
	})

	logrus.WithFields(logrus.Fields{
		"dry-run": *dryRun,
	}).Info("dry-run status")

	logrus.WithFields(logrus.Fields{
		"tlsDisabled":   *tlsDisabled,
		"addr":          *addr,
		"limitMemory":   *limitMemory,
		"limitCPU":      *limitCPU,
		"requestMemory": *requestMemory,
		"requestCPU":    *requestCPU,
	}).Info("programm flags")

	defaultResourceRequirements, err := parseResourceRequirements(*limitMemory, *limitCPU, *requestMemory, *requestCPU)
	if err != nil {
		log.Fatalf("could not parse resource requirements based on program flags: %s", err)
	}

	logrus.WithFields(logrus.Fields{
		"defaultResourceRequirements": defaultResourceRequirements,
	}).Info("parsed defaultResourceRequirements")

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := webhook.Mutate(w, r, defaultResourceRequirements, *dryRun)
		if err != nil {
			//todo: use "Fatalf" instead of "Printf"???
			log.Printf("mutation failed: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte("400 - Bad request!"))
			if err != nil {
				log.Fatalf("could not write to response: %s", err)
			}
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

func parseResourceRequirements(memoryLimit, CPULimit, memoryRequest, CPURequest string) (v1.ResourceRequirements, error) {
	defaultMemoryLimit, err := resource.ParseQuantity(memoryLimit)
	if err != nil {
		return v1.ResourceRequirements{}, fmt.Errorf("failed to parse memoryLimit quantity: %s", err)
	}
	defaultCPULimit, err := resource.ParseQuantity(CPULimit)
	if err != nil {
		return v1.ResourceRequirements{}, fmt.Errorf("failed to parse CPULimit quantity: %s", err)
	}
	defaultMemoryRequest, err := resource.ParseQuantity(memoryRequest)
	if err != nil {
		return v1.ResourceRequirements{}, fmt.Errorf("failed to parse memoryRequest quantity: %s", err)
	}
	defaultCPURequest, err := resource.ParseQuantity(CPURequest)
	if err != nil {
		return v1.ResourceRequirements{}, fmt.Errorf("failed to parse CPURequest quantity: %s", err)
	}

	return v1.ResourceRequirements{
		Limits: v1.ResourceList{
			v1.ResourceMemory: defaultMemoryLimit,
			v1.ResourceCPU:    defaultCPULimit,
		},
		Requests: v1.ResourceList{
			v1.ResourceMemory: defaultMemoryRequest,
			v1.ResourceCPU:    defaultCPURequest,
		},
	}, nil
}
