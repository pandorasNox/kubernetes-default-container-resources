package webhook

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Code    int    `json:"code,omitempty"`
}

// AdmissionResponse wrapper for the incomming response from kubernetes
type AdmissionResponse struct {
	Kind    string `json:"kind"`
	Request struct {
		UID    string
		Object struct {
			Spec struct {
				Containers []Container
			}
		}
	}
}

// Container representation for kubernetesyaml/json container definition
type Container struct {
	Resources ComputeResources
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

// Mutate responds to kubernetes webhooks request to add resource limits.
func Mutate(w http.ResponseWriter, r *http.Request, limitMemory, limitCPU, requestMemory, requestCPU string) error {

	admissionResponse := &AdmissionResponse{}
	err := json.NewDecoder(r.Body).Decode(admissionResponse)
	if err != nil {
		return fmt.Errorf("failed to decode body: %s", err)
	}

	containers := admissionResponse.Request.Object.Spec.Containers
	requestUID := admissionResponse.Request.UID
	admissionReview, err := admissionReview(containers, requestUID, limitMemory, limitCPU, requestMemory, requestCPU)
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

func admissionReview(containers []Container, UID, memoryLimit, CPULimit, memoryRequest, CPURequest string) (AdmissionReview, error) {

	patches := podPatches(containers, memoryLimit, CPULimit, memoryRequest, CPURequest)

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

func podPatches(containers []Container, memoryLimit, CPULimit, memoryRequest, CPURequest string) []Patch {
	patches := []Patch{}
	for i, c := range containers {
		containerPatches := containerPatches(i, c.Resources, memoryLimit, CPULimit, memoryRequest, CPURequest)
		for _, p := range containerPatches {
			patches = append(patches, p)
		}
	}
	return patches
}

func containerPatches(index int, cr ComputeResources, memoryLimit, CPULimit, memoryRequest, CPURequest string) []Patch {
	patches := []Patch{}

	if memoryAndCPUPairExists(cr) {
		return patches
	}

	patchValue := ComputeResources{
		Limits: ComputeUnit{
			Memory: memoryLimit,
			CPU:    CPULimit,
		},
		Requests: ComputeUnit{
			Memory: memoryRequest,
			CPU:    CPURequest,
		},
	}

	if !isMemoryEmpty(cr) {
		patchValue.Limits.Memory = cr.Limits.Memory
		patchValue.Requests.Memory = cr.Requests.Memory
	}
	if !isCPUEmpty(cr) {
		patchValue.Limits.CPU = cr.Limits.CPU
		patchValue.Requests.CPU = cr.Requests.CPU
	}

	patches = append(patches, createPatch(
		"replace",
		index,
		"resources",
		patchValue,
	))

	return patches
}

func memoryAndCPUPairExists(cr ComputeResources) bool {
	return (cr.Limits.Memory != "" && cr.Limits.CPU != "") ||
		(cr.Requests.Memory != "" && cr.Requests.CPU != "") ||
		(cr.Limits.Memory != "" && cr.Requests.CPU != "") ||
		(cr.Requests.Memory != "" && cr.Limits.CPU != "")
}

func isMemoryEmpty(cr ComputeResources) bool {
	return cr.Limits.Memory == "" && cr.Requests.Memory == ""
}

func isCPUEmpty(cr ComputeResources) bool {
	return cr.Limits.CPU == "" && cr.Requests.CPU == ""
}

func createPatch(op string, index int, containerSubPath, value interface{}) Patch {
	return Patch{
		Op:    op,
		Path:  strings.TrimRight(fmt.Sprintf("/spec/containers/%d/%s", index, containerSubPath), "/"),
		Value: value,
	}
}
