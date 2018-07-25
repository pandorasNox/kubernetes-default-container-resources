package webhook

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

// AdmissionReview is a validation/mutation object readable by the kubernetes api server.
type AdmissionReview struct {
	Response AdmissionReviewResponse `json:"response"`
}

// AdmissionReviewResponse is the response wrapper object for the AdmissionReview.
type AdmissionReviewResponse struct {
	UID       string          `json:"uid"`
	Allowed   bool            `json:"allowed"`
	Status    AdmissionStatus `json:"status,omitempty"`
	Patch     base64String    `json:"patch"`
	PatchType string          `json:"patchType"`
}

type base64String string

// Patch represents a single JSONPatch operation
// @see http://jsonpatch.com/
type Patch struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

// AdmissionStatus JSON/struct wrapper for the status field of the AdmissionReviewResponse
type AdmissionStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
	Code    int    `json:"code"`
}

// AdmissionResponse wrapper for the incomming response from kubernetes
type AdmissionResponse struct {
	Kind    string `json:"kind"`
	Request struct {
		UID    string
		Object struct {
			Spec struct {
				Containers Containers
			}
		}
	}
}

// Containers representation for kubernetesyaml/json container definition
type Containers []struct {
	Name      string
	Env       []Env
	Resources ComputeResources
}

// Env representation for kubernetes yaml/json envirmoent entrie definition
type Env struct {
	Name  string
	Value string
}

// ComputeResources representation for kubernetes yaml/json compute resource definition
type ComputeResources struct {
	Limits   ComputeUnit `json:"limits,omitempty"`
	Requests ComputeUnit `json:"requests,omitempty"`
}

// ComputeUnit representation for kubernetes yaml/json single compute resource definition
type ComputeUnit struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

// ServeContent responds to kubernetes webhooks request to add resource limits.
func ServeContent(w http.ResponseWriter, r *http.Request, limitMemory, limitCPU, requestMemory, requestCPU string) error {

	admissionResponse := &AdmissionResponse{}
	err := json.NewDecoder(r.Body).Decode(admissionResponse)
	if err != nil {
		return fmt.Errorf("failed to decode body: %s", err)
	}

	containers := admissionResponse.Request.Object.Spec.Containers
	requestUID := admissionResponse.Request.UID
	admissionReview, err := getAdmissionReview(containers, requestUID, limitMemory, limitCPU, requestMemory, requestCPU)
	if err != nil {
		return fmt.Errorf("failed to get admissionReview: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(admissionReview)
	if err != nil {
		return fmt.Errorf("failed to send response: %s", err)
	}

	return nil
}

func getAdmissionReview(containers Containers, UID, memoryLimit, CPULimit, memoryRequest, CPURequest string) (AdmissionReview, error) {

	patches := []Patch{}
	for i, c := range containers {
		patches = definePatches(patches, i, c.Resources, memoryLimit, CPULimit, memoryRequest, CPURequest)
	}

	jsonPatch, err := json.Marshal(patches)
	if err != nil {
		return AdmissionReview{}, fmt.Errorf("failed to encode patch: %s", err)
	}

	return AdmissionReview{
		AdmissionReviewResponse{
			UID:       UID,
			Allowed:   true,
			Patch:     base64String(base64.StdEncoding.EncodeToString(jsonPatch)),
			PatchType: "JSONPatch",
		},
	}, nil
}

func definePatches(patches []Patch, i int, cr ComputeResources, memoryLimit, CPULimit, memoryRequest, CPURequest string) []Patch {
	if isMemoryEmpty(cr) {
		patches = append(patches, createPatch(i, "limits/memory", memoryLimit))
		patches = append(patches, createPatch(i, "requests/memory", memoryRequest))
	}
	if isCPUEmpty(cr) {
		patches = append(patches, createPatch(i, "limits/cpu", CPULimit))
		patches = append(patches, createPatch(i, "requests/cpu", CPURequest))
	}
	return patches
}

func isMemoryEmpty(cr ComputeResources) bool {
	return cr.Limits.Memory == "" && cr.Requests.Memory == ""
}

func isCPUEmpty(cr ComputeResources) bool {
	return cr.Limits.CPU == "" && cr.Requests.CPU == ""
}

func createPatch(index int, resource, amount string) Patch {
	return Patch{
		Op:    "add",
		Path:  fmt.Sprintf("/spec/containers/%d/resources/%s", index, resource),
		Value: amount,
	}
}
