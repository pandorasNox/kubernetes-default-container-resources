package container

import (
	"fmt"
	"strings"

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
	Patches(index int, containerRR k8s_v1.ResourceRequirements, defaultRR k8s_v1.ResourceRequirements) ([]Patch, error)
}

//ComplementMemOrCPU represents a strategy to patch container.
//It will add Mem if there is no Mem or it will add CPU if there is no CPU
//e.g. request.cpu is set it will patch request.memory and limit.memory
//if nothing is set it will patch the given default
type ComplementMemOrCPU struct{}

//Patches ...
func (c ComplementMemOrCPU) Patches(index int, containerRR k8s_v1.ResourceRequirements, defaultRR k8s_v1.ResourceRequirements) []Patch {
	patches := []Patch{}

	if memoryAndCPUPairExists(containerRR) {
		return patches
	}

	patchValue := k8s_v1.ResourceRequirements{
		Limits:   k8s_v1.ResourceList{},
		Requests: k8s_v1.ResourceList{},
	}

	//keep original demanded compute values
	if isMemoryEmpty(containerRR) {
		patchValue.Limits[k8s_v1.ResourceMemory] = defaultRR.Limits[k8s_v1.ResourceMemory]
		patchValue.Requests[k8s_v1.ResourceMemory] = defaultRR.Requests[k8s_v1.ResourceMemory]
	} else {
		if mapKeyExist(containerRR.Limits, k8s_v1.ResourceMemory) {
			patchValue.Limits[k8s_v1.ResourceMemory] = containerRR.Limits[k8s_v1.ResourceMemory]
		}
		if mapKeyExist(containerRR.Requests, k8s_v1.ResourceMemory) {
			patchValue.Requests[k8s_v1.ResourceMemory] = containerRR.Requests[k8s_v1.ResourceMemory]
		}
	}
	if isCPUEmpty(containerRR) {
		patchValue.Limits[k8s_v1.ResourceCPU] = defaultRR.Limits[k8s_v1.ResourceCPU]
		patchValue.Requests[k8s_v1.ResourceCPU] = defaultRR.Requests[k8s_v1.ResourceCPU]
	} else {
		if mapKeyExist(containerRR.Limits, k8s_v1.ResourceCPU) {
			patchValue.Limits[k8s_v1.ResourceCPU] = containerRR.Limits[k8s_v1.ResourceCPU]
		}
		if mapKeyExist(containerRR.Requests, k8s_v1.ResourceCPU) {
			patchValue.Requests[k8s_v1.ResourceCPU] = containerRR.Requests[k8s_v1.ResourceCPU]
		}
	}

	patches = append(patches, createPatch(
		"replace",
		index,
		"resources",
		patchValue,
	))

	return patches
}

//DefaultOnEmpty ...
type DefaultOnEmpty struct{}

//Patches ...
func (d DefaultOnEmpty) Patches(index int, containerRR k8s_v1.ResourceRequirements, defaultRR k8s_v1.ResourceRequirements) []Patch {
	patches := []Patch{}

	if !isResourcesEmpty(containerRR) {
		return patches
	}

	patches = append(patches, createPatch(
		"replace",
		index,
		"resources",
		defaultRR,
	))

	return patches
}

//ComplementToDefault ...
type ComplementToDefault struct{}

//Patches ...
func (d ComplementToDefault) Patches(index int, containerRR k8s_v1.ResourceRequirements, defaultRR k8s_v1.ResourceRequirements) []Patch {
	patches := []Patch{}

	patchValue := defaultRR

	//keep original demanded compute resource values
	if mapKeyExist(containerRR.Limits, k8s_v1.ResourceMemory) {
		patchValue.Limits[k8s_v1.ResourceMemory] = defaultRR.Limits[k8s_v1.ResourceMemory]
	}
	if mapKeyExist(containerRR.Limits, k8s_v1.ResourceCPU) {
		patchValue.Limits[k8s_v1.ResourceCPU] = defaultRR.Limits[k8s_v1.ResourceCPU]
	}
	if mapKeyExist(containerRR.Requests, k8s_v1.ResourceMemory) {
		patchValue.Requests[k8s_v1.ResourceMemory] = defaultRR.Requests[k8s_v1.ResourceMemory]
	}
	if mapKeyExist(containerRR.Requests, k8s_v1.ResourceCPU) {
		patchValue.Requests[k8s_v1.ResourceCPU] = defaultRR.Requests[k8s_v1.ResourceCPU]
	}

	requestMem := patchValue.Requests[k8s_v1.ResourceMemory]
	limitMem := patchValue.Limits[k8s_v1.ResourceMemory]
	if requestMem.Cmp(limitMem) == 1 && mapKeyExist(containerRR.Requests, k8s_v1.ResourceMemory) &&
		mapKeyExist(containerRR.Limits, k8s_v1.ResourceMemory) {
		return []Patch{}
	}
	if requestMem.Cmp(limitMem) == 1 && mapKeyExist(containerRR.Limits, k8s_v1.ResourceMemory) {
		patchValue.Requests[k8s_v1.ResourceMemory] = patchValue.Limits[k8s_v1.ResourceMemory]
	} else if requestMem.Cmp(limitMem) == 1 {
		patchValue.Limits[k8s_v1.ResourceMemory] = patchValue.Requests[k8s_v1.ResourceMemory]
	}

	requestCPU := patchValue.Requests[k8s_v1.ResourceCPU]
	limitCPU := patchValue.Limits[k8s_v1.ResourceCPU]
	if requestCPU.Cmp(limitCPU) == 1 && mapKeyExist(containerRR.Requests, k8s_v1.ResourceCPU) &&
		mapKeyExist(containerRR.Limits, k8s_v1.ResourceCPU) {
		return []Patch{}
	}
	if requestCPU.Cmp(limitCPU) == 1 && mapKeyExist(containerRR.Limits, k8s_v1.ResourceCPU) {
		patchValue.Requests[k8s_v1.ResourceCPU] = patchValue.Limits[k8s_v1.ResourceCPU]
	} else if requestCPU.Cmp(limitCPU) == 1 {
		patchValue.Limits[k8s_v1.ResourceCPU] = patchValue.Requests[k8s_v1.ResourceCPU]
	}

	patches = append(patches, createPatch(
		"replace",
		index,
		"resources",
		patchValue,
	))

	return patches
}

func isResourcesEmpty(rr k8s_v1.ResourceRequirements) bool {
	return !mapKeyExist(rr.Limits, k8s_v1.ResourceMemory) && !mapKeyExist(rr.Limits, k8s_v1.ResourceCPU) &&
		!mapKeyExist(rr.Requests, k8s_v1.ResourceMemory) && !mapKeyExist(rr.Requests, k8s_v1.ResourceCPU)
}

func memoryAndCPUPairExists(rr k8s_v1.ResourceRequirements) bool {
	return (mapKeyExist(rr.Limits, k8s_v1.ResourceMemory) && mapKeyExist(rr.Limits, k8s_v1.ResourceCPU)) ||
		(mapKeyExist(rr.Requests, k8s_v1.ResourceMemory) && mapKeyExist(rr.Requests, k8s_v1.ResourceCPU)) ||
		(mapKeyExist(rr.Limits, k8s_v1.ResourceMemory) && mapKeyExist(rr.Requests, k8s_v1.ResourceCPU)) ||
		(mapKeyExist(rr.Requests, k8s_v1.ResourceMemory) && mapKeyExist(rr.Limits, k8s_v1.ResourceCPU))
}

func isMemoryEmpty(rr k8s_v1.ResourceRequirements) bool {
	return !mapKeyExist(rr.Limits, k8s_v1.ResourceMemory) && !mapKeyExist(rr.Requests, k8s_v1.ResourceMemory)
}

func isCPUEmpty(rr k8s_v1.ResourceRequirements) bool {
	return !mapKeyExist(rr.Limits, k8s_v1.ResourceCPU) && !mapKeyExist(rr.Requests, k8s_v1.ResourceCPU)
}

func mapKeyExist(rl k8s_v1.ResourceList, key k8s_v1.ResourceName) bool {
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
