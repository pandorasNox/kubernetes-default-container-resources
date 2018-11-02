package webhook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"k8s.io/api/admission/v1beta1"
	"k8s.io/api/core/v1"
	k8s_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Patch represents a single JSONPatch operation
// @see http://jsonpatch.com/
type Patch struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

// Mutate responds to kubernetes webhooks request to add resource limits.
func Mutate(w http.ResponseWriter, r *http.Request, defaults v1.ResourceRequirements) error {

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

	resp, err := podPatches(pod.Spec.Containers, defaults)

	resp.UID = incomingAdmissionReview.Request.UID
	patchType := v1beta1.PatchTypeJSONPatch
	resp.PatchType = &patchType

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(v1beta1.AdmissionReview{Response: resp})
	if err != nil {
		return fmt.Errorf("failed to encode outgoingAdmissionReview and send response: %s", err)
	}

	return nil
}

func podPatches(cc []v1.Container, defaults v1.ResourceRequirements) (*v1beta1.AdmissionResponse, error) {

	resp := &v1beta1.AdmissionResponse{}
	patches := []Patch{}
	for i, c := range cc {
		r, err := addDefaults(c.Resources, defaults)
		if err != nil {
			resp.Allowed = false
			resp.Result = &metav1.Status{
				Message: err.Error(),
			}
			return nil, fmt.Errorf("can't patch container with name: %s, reason: %s", c.Name, err)
		}
		patches = append(patches, Patch{
			Op:    "replace",
			Path:  filepath.Join("/spec/containers", strconv.Itoa(i), "resources"),
			Value: r,
		})
	}

	json, err := json.Marshal(patches)
	if err != nil {
		return nil, fmt.Errorf("failed to encode patch: %s", err)
	}

	resp.Allowed = true
	resp.Patch = []byte(json)

	return resp, nil
}

func addDefaults(c k8s_v1.ResourceRequirements, d k8s_v1.ResourceRequirements) (k8s_v1.ResourceRequirements, error) {

	if _, found := c.Limits[k8s_v1.ResourceMemory]; !found {
		c.Limits[k8s_v1.ResourceMemory] = d.Limits[k8s_v1.ResourceMemory]
	}
	if _, found := c.Limits[k8s_v1.ResourceCPU]; !found {
		c.Limits[k8s_v1.ResourceCPU] = d.Limits[k8s_v1.ResourceCPU]
	}
	if _, found := c.Requests[k8s_v1.ResourceMemory]; !found {
		c.Requests[k8s_v1.ResourceMemory] = d.Requests[k8s_v1.ResourceMemory]
	}
	if _, found := c.Requests[k8s_v1.ResourceCPU]; !found {
		c.Requests[k8s_v1.ResourceCPU] = d.Requests[k8s_v1.ResourceCPU]
	}

	requestMem := c.Requests[k8s_v1.ResourceMemory]
	limitMem := c.Limits[k8s_v1.ResourceMemory]
	if requestMem.Cmp(limitMem) == 1 {
		return c, fmt.Errorf("requested memory is greater than memory limit")
	}

	requestCPU := c.Requests[k8s_v1.ResourceCPU]
	limitCPU := c.Limits[k8s_v1.ResourceCPU]
	if requestCPU.Cmp(limitCPU) == 1 {
		return c, fmt.Errorf("requested cpu is greater than cpu limit")
	}

	return c, nil
}

func ParseResourceRequirements(memoryLimit, CPULimit, memoryRequest, CPURequest string) (v1.ResourceRequirements, error) {
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
