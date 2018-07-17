/**/
package main

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"reflect"

	"github.com/golang/glog"
)

// AdmissionReview returns a validation to kubernetes api server
type AdmissionReview struct {
	Response Response `json:"response"`
}

// Response foo
type Response struct {
	UID     string          `json:"uid"`
	Allowed bool            `json:"allowed"`
	Status  AdmissionStatus `json:"status"`
	// Patch     []Operation     `json:"patch"`
	Patch     interface{} `json:"patch"`
	PatchType string      `json:"patchType"`
}

// Operation is foo
type Operation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

// AdmissionStatus is baz
type AdmissionStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
	Code    int    `json:"code"`
}

// AdmissionResponse is foo
type AdmissionResponse struct {
	Kind    string `json:"kind"`
	Request struct {
		UID    string
		Object struct {
			Spec struct {
				Containers []struct {
					Name string
					Env  []struct {
						Name  string
						Value string
					}
					Resources ComputeResources
				}
			}
		}
	}
}

// ComputeResources is foo
type ComputeResources struct {
	Limits struct {
		CPU    string
		Memory string
	}
	Requests struct {
		CPU    string
		Memory string
	}
}

// IsEmpty make bar
func (r ComputeResources) IsEmpty() bool {
	return reflect.DeepEqual(r, ComputeResources{})
}

func serveContent(w http.ResponseWriter, r *http.Request) {
	glog.V(2).Info("mutating")

	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		glog.Errorf("contentType=%s, expect application/json", contentType)
		return
	}

	// admissionStatus := new(AdmissionStatus)
	// admissionReview := &AdmissionReview{
	// 	Response: struct {
	// 		Allowed bool            `json:"allowed"`
	// 		Status  AdmissionStatus `json:"status"`
	// 		Patch   []Operation     `json:"patch"`
	// 	}{
	// 		Allowed: true,
	// 		Status:  *admissionStatus,
	// 		Patch: []Operation{
	// 			// kubectl patch pod valid-pod --type='json' -p='[{"op": "replace", "path": "/spec/containers/0/image", "value":"new image"}]'
	// 			Operation{"add", "/spec/containers", `[{"resources":{"limits":{"memory":"128Mi","cpu":"100m"}}}]`},
	// 		},
	// 	},
	// }

	// anonymousStruct := struct {
	// 	NESCarts      []string
	// 	numberOfCarts int
	// }{
	// 	[]string{"Battletoads", "Mega Man 1", "Clash at Demonhead"},
	// 	len([]string{"Battletoads", "Mega Man 1", "Clash at Demonhead"}),
	// }
	// fmt.Println(" ")
	// fmt.Println("anonymousStruct:")
	// fmt.Println(anonymousStruct)
	// fmt.Println(" ")

	// fmt.Println("body:")
	admissionResponse := &AdmissionResponse{}
	json.NewDecoder(r.Body).Decode(admissionResponse)
	fmt.Println(" ")
	fmt.Println("admissionResponse:")
	fmt.Println(admissionResponse)
	fmt.Println(" ")

	patch := make([]Operation, 0)
	// patch := make([]struct {
	// 	foo string
	// }, 1)
	for containerIndex, container := range admissionResponse.Request.Object.Spec.Containers {
		fmt.Println(" ")
		fmt.Println("container.Resources:")
		fmt.Println(container.Resources)
		fmt.Println(" ")
		if false == container.Resources.IsEmpty() {
			fmt.Println("i:", containerIndex, " ", "container.Resources not empty")
		}
		if container.Resources.IsEmpty() {
			fmt.Println("i:", containerIndex, " ", "container.Resources.IsEmpty")
		}

		// patch = append(patch, Operation{Op: "add", Path: "/data/mutation-stage-1", Value: "yes"})
	}

	fmt.Println(" ")
	fmt.Println("patch")
	fmt.Println(patch)
	fmt.Println(" ")

	// labels := struct {
	// 	Foo string `json:"foo"`
	// }{
	// 	Foo: "bar",
	// }

	// patch :=

	admissionReview := AdmissionReview{
		Response{
			UID:     admissionResponse.Request.UID,
			Allowed: true,
			// Patch: []Operation{
			// 	// Operation{Op: "add", Path: "/data/labels", Value: []byte(`{"foo":"bar"}`)},
			// 	Operation{
			// 		Op:    "add",
			// 		Path:  "/data/mutation-stage-1",
			// 		Value: "yes",
			// 	},
			// },
			// Patch: []struct {
			// 	Op    string      `json:"op"`
			// 	Path  string      `json:"path"`
			// 	Value interface{} `json:"value"`
			// }{
			// 	// Operation{Op: "add", Path: "/data/labels", Value: []byte(`{"foo":"bar"}`)},
			// 	struct {
			// 		Op    string      `json:"op"`
			// 		Path  string      `json:"path"`
			// 		Value interface{} `json:"value"`
			// 	}{
			// 		Op:    "add",
			// 		Path:  "/data/mutation-stage-1",
			// 		Value: "yes",
			// 	},
			// },
			// {"op": "add", "path": "/secrets/", "value": {"name": "whatever" } }
			Patch: base64.StdEncoding.EncodeToString([]byte(`[
				{"op":"add","path":"/spec/initContainers","value":[{"image":"nginx:1.14.0","name":"webhook-added-init-container"}]}
			]`)),
			PatchType: "JSONPatch",
		},
	}

	fmt.Println(" ")
	fmt.Println("admissionReview")
	json.NewEncoder(os.Stdout).Encode(admissionReview)
	fmt.Println(" ")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(admissionReview)

	// addInitContainerPatch := `
	// 	{"response":{"allowed":true,"status":{"status":"","message":"","reason":"","code":0},"patch":	[
	// 		{"op":"add","path":"/spec/initContainers","value":[{"image":"webhook-added-image","name":"webhook-added-init-container","resources":{}}]}
	//    ]}}`
	// w.Write([]byte(addInitContainerPatch))
}

func main() {
	var tlsDisabled *bool
	tlsDisabled = flag.Bool("tlsDisabled", false, "(optional) disables tls for the server")
	flag.Parse()

	// var config Config
	// config.addFlags()
	// flag.Parse()

	fmt.Println("tlsDisabled: ", bool(*tlsDisabled))

	http.HandleFunc("/mutate", serveContent)

	if bool(*tlsDisabled) {
		server := &http.Server{
			Addr: ":8083",
		}
		server.ListenAndServe()
		os.Exit(0)
	}

	// ==============
	// use TLS
	// ==============

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			// tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	// clientset := getClient()
	server := &http.Server{
		// Addr:      ":443",
		Addr: ":8083",
		// TLSConfig: configTLS(config, clientset),
		TLSConfig: cfg,
	}

	err := server.ListenAndServeTLS("/certs/ssl-cert.pem", "/certs/ssl-key.pem")
	log.Fatal(err)
}
