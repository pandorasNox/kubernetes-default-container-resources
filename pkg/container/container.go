package container

import (
	k8s_v1 "k8s.io/api/core/v1"
)

// Patch represents a single JSONPatch operation
// @see http://jsonpatch.com/
type Patch struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

type containerPatchStrategy interface {
	patches(containerRR k8s_v1.ResourceRequirements, defaultRR k8s_v1.ResourceRequirements) ([]Patch, error)
}
