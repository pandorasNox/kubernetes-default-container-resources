package webhook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/sirupsen/logrus"
	"k8s.io/api/admission/v1beta1"
	k8s_v1 "k8s.io/api/core/v1"
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
func Mutate(w http.ResponseWriter, r *http.Request, defaults k8s_v1.ResourceRequirements, dryRun bool) error {

	in := &v1beta1.AdmissionReview{}
	err := json.NewDecoder(r.Body).Decode(in)
	if err != nil {
		return fmt.Errorf("failed to json decode body: %s", err)
	}

	pod := k8s_v1.Pod{}
	if err := json.Unmarshal(in.Request.Object.Raw, &pod); err != nil {
		return fmt.Errorf("failed to Unmarshal Pod from incoming AdmissionReview: %s", err)
	}

	resp, err := createResponse(pod.Spec.Containers, defaults)
	if err != nil {
		return fmt.Errorf("failed to create response: %s", err)
	}

	resp.UID = in.Request.UID
	patchType := v1beta1.PatchTypeJSONPatch
	resp.PatchType = &patchType

	if resp.Result != nil && resp.Result.Status == metav1.StatusFailure {
		logrus.WithFields(logrus.Fields{
			"AdmissionResponse": resp,
		}).Warn("failed AdmissionResponse")
	}

	out := v1beta1.AdmissionReview{Response: resp}

	if dryRun {
		logrus.WithFields(logrus.Fields{
			"dryRun":          dryRun,
			"AdmissionReview": out,
		}).Info("DRY-RUN: supposed AdmissionReview")

		out.Response = &v1beta1.AdmissionResponse{Allowed: true}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(out)
	if err != nil {
		return fmt.Errorf("failed to encode and send response: %s", err)
	}

	logrus.WithFields(logrus.Fields{
		"AdmissionReview": out,
	}).Info("Success sended AdmissionReview")

	return nil
}

func createResponse(cc []k8s_v1.Container, defaults k8s_v1.ResourceRequirements) (*v1beta1.AdmissionResponse, error) {

	resp := &v1beta1.AdmissionResponse{}
	patches := []Patch{}
	for i, c := range cc {
		r, err := addDefaults(c.Resources, defaults)
		if err != nil {
			resp.Allowed = false
			resp.Result = &metav1.Status{
				Message: err.Error(),
				Status:  metav1.StatusFailure,
			}
			return resp, nil
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

	if c.Limits == nil {
		c.Limits = k8s_v1.ResourceList{}
	}
	if c.Requests == nil {
		c.Requests = k8s_v1.ResourceList{}
	}

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
