package webhook

import (
	"reflect"
	"strconv"
	"testing"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

var limitMemory = "1G"
var limitCPU = "0.5"
var requestMemory = "1G"
var requestCPU = "0.1"

func getResourceQuantity(quantity string) resource.Quantity {
	resDef, _ := resource.ParseQuantity(quantity)
	return resDef
}

var singleContainerPodTests = []struct {
	in  []v1.Container
	out []Patch
}{
	{
		[]v1.Container{
			v1.Container{
				Resources: v1.ResourceRequirements{},
			},
		},
		[]Patch{
			Patch{"replace", "/spec/containers/0/resources",
				v1.ResourceRequirements{
					Limits: v1.ResourceList{
						// Memory: "1G",
						// CPU:    "0.5",
						v1.ResourceMemory: getResourceQuantity(limitMemory),
						v1.ResourceCPU:    getResourceQuantity(limitCPU),
					},
					Requests: v1.ResourceList{
						// Memory: "1G",
						// CPU:    "0.1",
						v1.ResourceMemory: getResourceQuantity(requestMemory),
						v1.ResourceCPU:    getResourceQuantity(requestCPU),
					},
				},
			},
		},
	},
	{
		[]v1.Container{
			v1.Container{
				Resources: v1.ResourceRequirements{
					Limits: v1.ResourceList{
						// Memory: "1G",
						v1.ResourceMemory: getResourceQuantity(limitMemory),
					},
				},
			},
		},
		[]Patch{
			Patch{"replace", "/spec/containers/0/resources",
				v1.ResourceRequirements{
					Limits: v1.ResourceList{
						// Memory: "1G",
						// CPU:    "0.5",
						v1.ResourceMemory: getResourceQuantity(limitMemory),
						v1.ResourceCPU:    getResourceQuantity(limitCPU),
					},
					Requests: v1.ResourceList{
						// CPU:    "0.1",
						v1.ResourceCPU: getResourceQuantity(requestCPU),
					},
				},
			},
		},
	},
	{
		[]v1.Container{
			v1.Container{
				Resources: v1.ResourceRequirements{
					Requests: v1.ResourceList{
						// Memory: "1G",
						v1.ResourceMemory: getResourceQuantity(requestMemory),
					},
				},
			},
		},
		[]Patch{
			Patch{"replace", "/spec/containers/0/resources",
				v1.ResourceRequirements{
					Limits: v1.ResourceList{
						// CPU:    "0.5",
						v1.ResourceCPU: getResourceQuantity(limitCPU),
					},
					Requests: v1.ResourceList{
						// Memory: "1G",
						// CPU:    "0.1",
						v1.ResourceMemory: getResourceQuantity(requestMemory),
						v1.ResourceCPU:    getResourceQuantity(requestCPU),
					},
				},
			},
		},
	},
	{
		[]v1.Container{
			v1.Container{
				Resources: v1.ResourceRequirements{
					Limits: v1.ResourceList{
						// CPU: "0.5",
						v1.ResourceCPU: getResourceQuantity(limitCPU),
					},
				},
			},
		},
		[]Patch{
			Patch{"replace", "/spec/containers/0/resources",
				v1.ResourceRequirements{
					Limits: v1.ResourceList{
						// Memory: "1G",
						// CPU:    "0.5",
						v1.ResourceMemory: getResourceQuantity(limitMemory),
						v1.ResourceCPU:    getResourceQuantity(limitCPU),
					},
					Requests: v1.ResourceList{
						// Memory: "1G",
						v1.ResourceMemory: getResourceQuantity(requestMemory),
					},
				},
			},
		},
	},
	{
		[]v1.Container{
			v1.Container{
				Resources: v1.ResourceRequirements{
					Limits: v1.ResourceList{
						// Memory: "1G",
						// CPU:    "0.5",
						v1.ResourceMemory: getResourceQuantity(limitMemory),
						v1.ResourceCPU:    getResourceQuantity(limitCPU),
					},
				},
			},
		},
		[]Patch{},
	},
	{
		[]v1.Container{
			v1.Container{
				Resources: v1.ResourceRequirements{
					Requests: v1.ResourceList{
						// Memory: "1G",
						// CPU:    "0.5",
						v1.ResourceMemory: getResourceQuantity(requestMemory),
						v1.ResourceCPU:    getResourceQuantity(requestCPU),
					},
				},
			},
		},
		[]Patch{},
	},
	{
		[]v1.Container{
			v1.Container{
				Resources: v1.ResourceRequirements{
					Limits: v1.ResourceList{
						// Memory: "1G",
						v1.ResourceMemory: getResourceQuantity(limitMemory),
					},
					Requests: v1.ResourceList{
						// CPU: "0.1",
						v1.ResourceCPU: getResourceQuantity(requestCPU),
					},
				},
			},
		},
		[]Patch{},
	},
	{
		[]v1.Container{
			v1.Container{
				Resources: v1.ResourceRequirements{
					Limits: v1.ResourceList{
						// CPU: "0.5",
						v1.ResourceCPU: getResourceQuantity(limitCPU),
					},
					Requests: v1.ResourceList{
						// Memory: "1G",
						v1.ResourceMemory: getResourceQuantity(requestMemory),
					},
				},
			},
		},
		[]Patch{},
	},
}

var multiContainerPodTests = []struct {
	in  []v1.Container
	out []Patch
}{
	{
		[]v1.Container{
			v1.Container{
				Resources: v1.ResourceRequirements{},
			},
			v1.Container{
				Resources: v1.ResourceRequirements{},
			},
		},
		[]Patch{
			Patch{"replace", "/spec/containers/0/resources",
				v1.ResourceRequirements{
					Limits: v1.ResourceList{
						// Memory: "1G",
						// CPU:    "0.5",
						v1.ResourceMemory: getResourceQuantity(limitMemory),
						v1.ResourceCPU:    getResourceQuantity(limitCPU),
					},
					Requests: v1.ResourceList{
						// Memory: "1G",
						// CPU:    "0.1",
						v1.ResourceMemory: getResourceQuantity(requestMemory),
						v1.ResourceCPU:    getResourceQuantity(requestCPU),
					},
				},
			},
			Patch{"replace", "/spec/containers/1/resources",
				v1.ResourceRequirements{
					Limits: v1.ResourceList{
						// Memory: "1G",
						// CPU:    "0.5",
						v1.ResourceMemory: getResourceQuantity(limitMemory),
						v1.ResourceCPU:    getResourceQuantity(limitCPU),
					},
					Requests: v1.ResourceList{
						// Memory: "1G",
						// CPU:    "0.1",
						v1.ResourceMemory: getResourceQuantity(requestMemory),
						v1.ResourceCPU:    getResourceQuantity(requestCPU),
					},
				},
			},
		},
	},
	{
		[]v1.Container{
			v1.Container{
				Resources: v1.ResourceRequirements{
					Limits: v1.ResourceList{
						// CPU: "0.5",
						v1.ResourceCPU: getResourceQuantity(limitCPU),
					},
					Requests: v1.ResourceList{
						// Memory: "1G",
						v1.ResourceMemory: getResourceQuantity(requestMemory),
					},
				},
			},
			v1.Container{
				Resources: v1.ResourceRequirements{},
			},
		},
		[]Patch{
			Patch{"replace", "/spec/containers/1/resources",
				v1.ResourceRequirements{
					Limits: v1.ResourceList{
						// Memory: "1G",
						// CPU:    "0.5",
						v1.ResourceMemory: getResourceQuantity(limitMemory),
						v1.ResourceCPU:    getResourceQuantity(limitCPU),
					},
					Requests: v1.ResourceList{
						// Memory: "1G",
						// CPU:    "0.1",
						v1.ResourceMemory: getResourceQuantity(requestMemory),
						v1.ResourceCPU:    getResourceQuantity(requestCPU),
					},
				},
			},
		},
	},
}

func TestSingleContainerPodPatches(t *testing.T) {
	for i, tt := range singleContainerPodTests {
		t.Run(""+strconv.Itoa(i), func(t *testing.T) {
			r, _ := podPatches(tt.in, limitMemory, limitCPU, requestMemory, requestCPU)
			// if !reflect.DeepEqual(r, tt.out) {
			// 	t.Errorf("got %q, want %q", r, tt.out)
			// }
			if prettyPrint(r) != prettyPrint(tt.out) {
				t.Errorf("got %q, want %q", prettyPrint(r), prettyPrint(tt.out))
			}
		})
	}
}

func TestMultiContainerPodPatches(t *testing.T) {
	for i, tt := range multiContainerPodTests {
		t.Run(""+strconv.Itoa(i), func(t *testing.T) {
			r, _ := podPatches(tt.in, limitMemory, limitCPU, requestMemory, requestCPU)
			if !reflect.DeepEqual(r, tt.out) {
				t.Errorf("got %q, want %q", r, tt.out)
			}
			// if prettyPrint(r) != prettyPrint(tt.out) {
			// 	t.Errorf("got %q, want %q", prettyPrint(r), prettyPrint(tt.out))
			// }
		})
	}
}
