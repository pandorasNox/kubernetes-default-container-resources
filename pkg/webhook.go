package webhook

import (
	"encoding/json"
	"fmt"
	"net/http"

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
func Mutate(w http.ResponseWriter, r *http.Request, patchStrategy container.ComplementMemOrCPU, defaultResourceRequirements v1.ResourceRequirements) error {

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
	outgoingAdmissionReview, err := admissionReview(patchStrategy, containers, requestUID, defaultResourceRequirements)
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

func admissionReview(patchStrategy container.ComplementMemOrCPU, containers []v1.Container, UID types.UID, defaultRR v1.ResourceRequirements) (v1beta1.AdmissionReview, error) {

	patches, err := podPatches(patchStrategy, containers, defaultRR)
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

func podPatches(patchStrategy container.ComplementMemOrCPU, containers []v1.Container, defaultRR v1.ResourceRequirements) ([]container.Patch, error) {
	patches := []container.Patch{}
	for i, c := range containers {
		containerPatches := patchStrategy.Patches(i, c.Resources, defaultRR)

		for _, p := range containerPatches {
			patches = append(patches, p)
		}
	}
	return patches, nil
}
