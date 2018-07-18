package webhook

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"reflect"
	"strconv"

	"github.com/kubernetes/kubernetes/pkg/kubelet/kubeletconfig/util/log"
)

// AdmissionReview is a validation/mutation object readable by the kubernetes api server
type AdmissionReview struct {
	Response AdmissionReviewResponse `json:"response"`
}

// AdmissionReviewResponse is the response wrapper object for the AdmissionReview
type AdmissionReviewResponse struct {
	UID       string          `json:"uid"`
	Allowed   bool            `json:"allowed"`
	Status    AdmissionStatus `json:"status"`
	Patch     Base64String    `json:"patch"`
	PatchType string          `json:"patchType"`
}

//Base64String should be a base64 encoded string
type Base64String string

// Operation ...
// @see http://jsonpatch.com/
type Operation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

// AdmissionStatus ...
type AdmissionStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
	Code    int    `json:"code"`
}

// AdmissionResponse ...
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

// ComputeResources ...
type ComputeResources struct {
	Limits   ComputeUnit `json:"limits,omitempty"`
	Requests ComputeUnit `json:"requests,omitempty"`
}

// ComputeUnit ...
type ComputeUnit struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

// IsEmpty ...
func (r ComputeResources) IsEmpty() bool {
	return reflect.DeepEqual(r, ComputeResources{})
}

// ServeContent ...
func ServeContent(w http.ResponseWriter, r *http.Request) {
	fmt.Println("requestDump:")
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))
	fmt.Println("")

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		log.Errorf("contentType=%s, expect application/json", contentType)
		return
	}

	admissionResponse := &AdmissionResponse{}
	json.NewDecoder(r.Body).Decode(admissionResponse)
	fmt.Println("admissionResponse:")
	fmt.Println(admissionResponse)
	fmt.Println("")

	patch := []Operation{}
	for i, container := range admissionResponse.Request.Object.Spec.Containers {

		if false == container.Resources.IsEmpty() {
			continue
		}

		// @see http://jsonpatch.com/
		patch = append(
			patch,
			Operation{
				Op:   "add",
				Path: "/spec/containers/" + strconv.Itoa(i) + "/resources",
				Value: ComputeResources{
					Limits: ComputeUnit{
						Memory: "512Mi",
					},
					Requests: ComputeUnit{
						Memory: "512Mi",
					},
				},
			},
		)
	}

	jsonPatch, err := json.Marshal(patch)

	fmt.Println("patch:")
	fmt.Println(string(jsonPatch))
	fmt.Println("")

	admissionReview := AdmissionReview{
		AdmissionReviewResponse{
			UID:       admissionResponse.Request.UID,
			Allowed:   true,
			Patch:     Base64String(base64.StdEncoding.EncodeToString(jsonPatch)),
			PatchType: "JSONPatch",
		},
	}

	fmt.Println("admissionReview:")
	json.NewEncoder(os.Stdout).Encode(admissionReview)
	fmt.Println("")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(admissionReview)
}
