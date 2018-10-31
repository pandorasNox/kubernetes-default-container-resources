package webhook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pandorasnox/kubernetes-default-container-resources/pkg/container"
	"k8s.io/api/admission/v1beta1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
)

//ParseResourceRequirements parses string resource representations to ResourceRequirements
func ParseResourceRequirements(memoryLimit, CPULimit, memoryRequest, CPURequest string) (v1.ResourceRequirements, error) {
	defaultMemoryLimit, err := resource.ParseQuantity(memoryLimit)
	if err != nil {
		return v1.ResourceRequirements{}, fmt.Errorf("failed to parse memoryLimit quanttiy: %s", err)
	}
	defaultCPULimit, err := resource.ParseQuantity(CPULimit)
	if err != nil {
		return v1.ResourceRequirements{}, fmt.Errorf("failed to parse CPULimit quanttiy: %s", err)
	}
	defaultMemoryRequest, err := resource.ParseQuantity(memoryRequest)
	if err != nil {
		return v1.ResourceRequirements{}, fmt.Errorf("failed to parse memoryRequest quanttiy: %s", err)
	}
	defaultCPURequest, err := resource.ParseQuantity(CPURequest)
	if err != nil {
		return v1.ResourceRequirements{}, fmt.Errorf("failed to parse CPURequest quanttiy: %s", err)
	}

	resourceRequirements := v1.ResourceRequirements{
		Limits: v1.ResourceList{
			v1.ResourceMemory: defaultMemoryLimit,
			v1.ResourceCPU:    defaultCPULimit,
		},
		Requests: v1.ResourceList{
			v1.ResourceMemory: defaultMemoryRequest,
			v1.ResourceCPU:    defaultCPURequest,
		},
	}

	return resourceRequirements, nil
}

// Mutate responds to kubernetes webhooks request to add resource limits.
func Mutate(w http.ResponseWriter, r *http.Request, defaultResourceRequirements v1.ResourceRequirements) error {

	incomingAdmissionReview := &v1beta1.AdmissionReview{}
	err := json.NewDecoder(r.Body).Decode(incomingAdmissionReview)
	if err != nil {
		return fmt.Errorf("failed to decode body: %s", err)
	}

	raw := incomingAdmissionReview.Request.Object.Raw
	pod := v1.Pod{}
	if err := json.Unmarshal(raw, &pod); err != nil {
		return fmt.Errorf("failed to Unmarshal Pod from incoming AdmissionReview: %s", err)

	}

	containers := pod.Spec.Containers
	requestUID := incomingAdmissionReview.Request.UID
	outgoingAdmissionReview, err := admissionReview(containers, requestUID, defaultResourceRequirements)
	if err != nil {
		return fmt.Errorf("failed to get outgoingAdmissionReview: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(outgoingAdmissionReview)
	if err != nil {
		return fmt.Errorf("failed to encode outgoingAdmissionReview and send response: %s", err)
	}

	return nil
}

func admissionReview(containers []v1.Container, UID types.UID, defaultRR v1.ResourceRequirements) (v1beta1.AdmissionReview, error) {

	patches, err := podPatches(containers, defaultRR)
	if err != nil {
		return v1beta1.AdmissionReview{}, fmt.Errorf("failed to get patches for pod: %s", err)
	}

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

func podPatches(containers []v1.Container, defaultRR v1.ResourceRequirements) ([]container.Patch, error) {
	patches := []container.Patch{}
	for i, c := range containers {
		containerPatches, err := containerPatches(i, c.Resources, defaultRR)
		if err != nil {
			return nil, fmt.Errorf("failed to get containerPatches: %s", err)
		}

		for _, p := range containerPatches {
			patches = append(patches, p)
		}
	}
	return patches, nil
}

func containerPatches(index int, containerRR v1.ResourceRequirements, defaultRR v1.ResourceRequirements) ([]container.Patch, error) {
	patches := []container.Patch{}

	if memoryAndCPUPairExists(containerRR) {
		return patches, nil
	}

	patchValue := v1.ResourceRequirements{
		Limits:   v1.ResourceList{},
		Requests: v1.ResourceList{},
	}

	//keep original demanded compute values
	if isMemoryEmpty(containerRR) {
		patchValue.Limits[v1.ResourceMemory] = defaultRR.Limits[v1.ResourceMemory]
		patchValue.Requests[v1.ResourceMemory] = defaultRR.Requests[v1.ResourceMemory]
	} else {
		v, ok := containerRR.Limits[v1.ResourceMemory]
		if ok {
			patchValue.Limits[v1.ResourceMemory] = v
		}
		v, ok = containerRR.Requests[v1.ResourceMemory]
		if ok {
			patchValue.Requests[v1.ResourceMemory] = v
		}
	}
	if isCPUEmpty(containerRR) {
		patchValue.Limits[v1.ResourceCPU] = defaultRR.Limits[v1.ResourceCPU]
		patchValue.Requests[v1.ResourceCPU] = defaultRR.Requests[v1.ResourceCPU]
	} else {
		v, ok := containerRR.Limits[v1.ResourceCPU]
		if ok {
			patchValue.Limits[v1.ResourceCPU] = v
		}
		v, ok = containerRR.Requests[v1.ResourceCPU]
		if ok {
			patchValue.Requests[v1.ResourceCPU] = v
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

func memoryAndCPUPairExists(rr v1.ResourceRequirements) bool {
	return (mapKeyExist(rr.Limits, v1.ResourceMemory) && mapKeyExist(rr.Limits, v1.ResourceCPU)) ||
		(mapKeyExist(rr.Requests, v1.ResourceMemory) && mapKeyExist(rr.Requests, v1.ResourceCPU)) ||
		(mapKeyExist(rr.Limits, v1.ResourceMemory) && mapKeyExist(rr.Requests, v1.ResourceCPU)) ||
		(mapKeyExist(rr.Requests, v1.ResourceMemory) && mapKeyExist(rr.Limits, v1.ResourceCPU))
}

func isMemoryEmpty(rr v1.ResourceRequirements) bool {
	return !mapKeyExist(rr.Limits, v1.ResourceMemory) && !mapKeyExist(rr.Requests, v1.ResourceMemory)
}

func isCPUEmpty(rr v1.ResourceRequirements) bool {
	return !mapKeyExist(rr.Limits, v1.ResourceCPU) && !mapKeyExist(rr.Requests, v1.ResourceCPU)
}

func mapKeyExist(rl v1.ResourceList, key v1.ResourceName) bool {
	if _, keyExist := rl[key]; keyExist {
		return true
	}

	return false
}

func createPatch(op string, index int, containerSubPath, value interface{}) container.Patch {
	return container.Patch{
		Op:    op,
		Path:  strings.TrimRight(fmt.Sprintf("/spec/containers/%d/%s", index, containerSubPath), "/"),
		Value: value,
	}
}
