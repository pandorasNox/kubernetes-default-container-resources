package webhook

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
				Containers Containers
			}
		}
	}
}

// Containers ...
type Containers []struct {
	Name string
	Env  []struct {
		Name  string
		Value string
	}
	Resources ComputeResources
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

func isResourcesEmpty(cr ComputeResources) bool {
	return isMemoryEmpty(cr) && isCPUEmpty(cr)
}

func isMemoryEmpty(cr ComputeResources) bool {
	return cr.Limits.Memory == "" && cr.Requests.Memory == ""
}

func isCPUEmpty(cr ComputeResources) bool {
	return cr.Limits.CPU == "" && cr.Requests.CPU == ""
}

func patchResources(patches []Operation, i, limitMemory, limitCPU, requestMemory, requestCPU string) []Operation {
	// @see http://jsonpatch.com/
	patches = append(
		patches,
		Operation{
			Op:   "add",
			Path: "/spec/containers/" + i + "/resources",
			Value: ComputeResources{
				Limits: ComputeUnit{
					Memory: limitMemory,
					CPU:    limitCPU,
				},
				Requests: ComputeUnit{
					Memory: requestMemory,
					CPU:    requestCPU,
				},
			},
		},
	)

	return patches
}

func patchMemory(patches []Operation, i, limitMemory, requestMemory string) []Operation {
	patches = append(
		patches,
		Operation{
			Op:    "add",
			Path:  "/spec/containers/" + i + "/resources/limits/memory",
			Value: limitMemory,
		},
	)
	patches = append(
		patches,
		Operation{
			Op:    "add",
			Path:  "/spec/containers/" + i + "/resources/requests/memory",
			Value: requestMemory,
		},
	)

	return patches
}

func patchCPU(patches []Operation, i, limitCPU, requestCPU string) []Operation {
	patches = append(
		patches,
		Operation{
			Op:    "add",
			Path:  "/spec/containers/" + i + "/resources/limits/cpu",
			Value: limitCPU,
		},
	)
	patches = append(
		patches,
		Operation{
			Op:    "add",
			Path:  "/spec/containers/" + i + "/resources/requests/cpu",
			Value: requestCPU,
		},
	)

	return patches
}

func getAdmissionReview(c Containers, UID, limitMemory, limitCPU, requestMemory, requestCPU string) (AdmissionReview, error) {

	patches := []Operation{}
	for i, container := range c {

		if isResourcesEmpty(container.Resources) {
			patches = patchResources(patches, strconv.Itoa(i), limitMemory, limitCPU, requestMemory, requestCPU)
			continue
		}

		if isMemoryEmpty(container.Resources) {
			patches = patchMemory(patches, strconv.Itoa(i), limitMemory, requestMemory)
		}

		if isCPUEmpty(container.Resources) {
			patches = patchCPU(patches, strconv.Itoa(i), limitCPU, requestCPU)
		}
	}

	jsonPatch, err := json.Marshal(patches)
	if err != nil {
		return AdmissionReview{}, fmt.Errorf("failed to encode patch: %s", err)
	}

	admissionReview := AdmissionReview{
		AdmissionReviewResponse{
			UID:       UID,
			Allowed:   true,
			Patch:     base64String(base64.StdEncoding.EncodeToString(jsonPatch)),
			PatchType: "JSONPatch",
		},
	}

	return admissionReview, nil
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
