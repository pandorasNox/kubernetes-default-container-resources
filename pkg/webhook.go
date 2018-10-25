package webhook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"k8s.io/api/admission/v1beta1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
)

// Patch represents a single JSONPatch operation
// @see http://jsonpatch.com/
type Patch struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

func prettyPrint(i interface{}) string {
	// s, _ := json.MarshalIndent(i, "", "\t")
	s, _ := json.Marshal(i)
	return string(s)
}

// Mutate responds to kubernetes webhooks request to add resource limits.
func Mutate(w http.ResponseWriter, r *http.Request, limitMemory, limitCPU,
	requestMemory, requestCPU string) error {

	incomingAdmissionReview := &v1beta1.AdmissionReview{}
	err := json.NewDecoder(r.Body).Decode(incomingAdmissionReview)
	if err != nil {
		return fmt.Errorf("failed to decode body: %s", err)
	}
	// fmt.Printf("show k8s incomingAdmissionReview: %s", prettyPrint(incomingAdmissionReview))

	raw := incomingAdmissionReview.Request.Object.Raw
	pod := v1.Pod{}
	if err := json.Unmarshal(raw, &pod); err != nil {
		return fmt.Errorf("failed to Unmarshal Pod from incoming AdmissionReview: %s", err)

	}

	containers := pod.Spec.Containers
	requestUID := incomingAdmissionReview.Request.UID
	outgoingAdmissionReview, err := admissionReview(containers, requestUID, limitMemory, limitCPU, requestMemory, requestCPU)
	if err != nil {
		return fmt.Errorf("failed to get outgoingAdmissionReview: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(outgoingAdmissionReview)
	if err != nil {
		// failed to encode outgoingAdmissionReview into header???
		return fmt.Errorf("failed to send response: %s", err)
	}

	return nil
}

func admissionReview(containers []v1.Container, UID types.UID, memoryLimit, CPULimit, memoryRequest,
	CPURequest string) (v1beta1.AdmissionReview, error) {

	patches, err := podPatches(containers, memoryLimit, CPULimit, memoryRequest, CPURequest)
	//todo: handle err

	jsonPatch, err := json.Marshal(patches)
	if err != nil {
		return v1beta1.AdmissionReview{}, fmt.Errorf("failed to encode patch: %s", err)
	}

	patchType := v1beta1.PatchTypeJSONPatch
	return v1beta1.AdmissionReview{
		Response: &v1beta1.AdmissionResponse{
			UID:       UID,
			Allowed:   true,
			Patch:     []byte(jsonPatch),
			PatchType: &patchType,
		},
	}, nil
}

func podPatches(containers []v1.Container, memoryLimit, CPULimit, memoryRequest,
	CPURequest string) ([]Patch, error) {
	patches := []Patch{}
	for i, c := range containers {
		containerPatches, err := containerPatches(i, c.Resources, memoryLimit, CPULimit, memoryRequest, CPURequest)
		if err != nil {
			return nil, fmt.Errorf("failed to get containerPatches: %s", err)
		}

		for _, p := range containerPatches {
			patches = append(patches, p)
		}
	}
	return patches, nil
}

func containerPatches(index int, cr v1.ResourceRequirements, memoryLimit, CPULimit, memoryRequest,
	CPURequest string) ([]Patch, error) {
	patches := []Patch{}

	if memoryAndCPUPairExists(cr) {
		return patches, nil
	}

	defaultMemoryLimit, err := resource.ParseQuantity(memoryLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to parse memoryLimit quanttiy: %s", err)
	}
	defaultCPULimit, err := resource.ParseQuantity(CPULimit)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CPULimit quanttiy: %s", err)
	}
	defaultMemoryRequest, err := resource.ParseQuantity(memoryRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to parse memoryRequest quanttiy: %s", err)
	}
	defaultCPURequest, err := resource.ParseQuantity(CPURequest)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CPURequest quanttiy: %s", err)
	}

	patchValue := v1.ResourceRequirements{
		Limits: v1.ResourceList{
			v1.ResourceMemory: defaultMemoryLimit,
			v1.ResourceCPU:    defaultCPULimit,
		},
		Requests: v1.ResourceList{
			v1.ResourceMemory: defaultMemoryRequest,
			v1.ResourceCPU:    defaultCPURequest,
		},
	}

	//keep original demanded compute values
	if !isMemoryEmpty(cr) {
		patchValue.Limits[v1.ResourceMemory] = cr.Limits[v1.ResourceMemory]
		if !mapKeyExist(cr.Limits, v1.ResourceMemory) {
			delete(patchValue.Limits, v1.ResourceMemory)
		}

		patchValue.Requests[v1.ResourceMemory] = cr.Requests[v1.ResourceMemory]
		if !mapKeyExist(cr.Requests, v1.ResourceMemory) {
			delete(patchValue.Requests, v1.ResourceMemory)
		}
	}
	if !isCPUEmpty(cr) {
		patchValue.Limits[v1.ResourceCPU] = cr.Limits[v1.ResourceCPU]
		if !mapKeyExist(cr.Limits, v1.ResourceCPU) {
			delete(patchValue.Limits, v1.ResourceCPU)
		}

		patchValue.Requests[v1.ResourceCPU] = cr.Requests[v1.ResourceCPU]
		if !mapKeyExist(cr.Requests, v1.ResourceCPU) {
			delete(patchValue.Requests, v1.ResourceCPU)
		}
	}

	patches = append(patches, createPatch(
		"replace",
		index,
		"resources",
		patchValue,
	))

	return patches, nil
}

func memoryAndCPUPairExists(cr v1.ResourceRequirements) bool {
	return (mapKeyExist(cr.Limits, v1.ResourceMemory) && mapKeyExist(cr.Limits, v1.ResourceCPU)) ||
		(mapKeyExist(cr.Requests, v1.ResourceMemory) && mapKeyExist(cr.Requests, v1.ResourceCPU)) ||
		(mapKeyExist(cr.Limits, v1.ResourceMemory) && mapKeyExist(cr.Requests, v1.ResourceCPU)) ||
		(mapKeyExist(cr.Requests, v1.ResourceMemory) && mapKeyExist(cr.Limits, v1.ResourceCPU))
}

func isMemoryEmpty(cr v1.ResourceRequirements) bool {
	return !mapKeyExist(cr.Limits, v1.ResourceMemory) && !mapKeyExist(cr.Requests, v1.ResourceMemory)
}

func isCPUEmpty(cr v1.ResourceRequirements) bool {
	return !mapKeyExist(cr.Limits, v1.ResourceCPU) && !mapKeyExist(cr.Requests, v1.ResourceCPU)
}

func mapKeyExist(rl v1.ResourceList, key v1.ResourceName) bool {
	if _, keyExist := rl[key]; keyExist {
		return true
	}

	return false
}

func createPatch(op string, index int, containerSubPath, value interface{}) Patch {
	return Patch{
		Op:    op,
		Path:  strings.TrimRight(fmt.Sprintf("/spec/containers/%d/%s", index, containerSubPath), "/"),
		Value: value,
	}
}
