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

var defaultResourceRequirements, _ = ParseResourceRequirements(limitMemory, limitCPU, requestMemory, requestCPU)

func getResourceQuantity(quantity string) resource.Quantity {
	resDef, _ := resource.ParseQuantity(quantity)
	return resDef
}

var singleContainerPodTests = []struct {
	name string
	in   []v1.Container
	out  []Patch
}{
	{
		"container without ResourceRequirements",
		[]v1.Container{
			v1.Container{
				Resources: v1.ResourceRequirements{},
			},
		},
		[]Patch{
			Patch{"replace", "/spec/containers/0/resources",
				v1.ResourceRequirements{
					Limits: v1.ResourceList{
						v1.ResourceMemory: defaultResourceRequirements.Limits[v1.ResourceMemory],
						v1.ResourceCPU:    defaultResourceRequirements.Limits[v1.ResourceCPU],
					},
					Requests: v1.ResourceList{
						v1.ResourceMemory: defaultResourceRequirements.Requests[v1.ResourceMemory],
						v1.ResourceCPU:    defaultResourceRequirements.Requests[v1.ResourceCPU],
					},
				},
			},
		},
	},
	{
		"container with limited memory",
		[]v1.Container{
			v1.Container{
				Resources: v1.ResourceRequirements{
					Limits: v1.ResourceList{
						v1.ResourceMemory: defaultResourceRequirements.Limits[v1.ResourceMemory],
					},
				},
			},
		},
		[]Patch{
			Patch{"replace", "/spec/containers/0/resources",
				v1.ResourceRequirements{
					Limits: v1.ResourceList{
						v1.ResourceMemory: defaultResourceRequirements.Limits[v1.ResourceMemory],
						v1.ResourceCPU:    defaultResourceRequirements.Limits[v1.ResourceCPU],
					},
					Requests: v1.ResourceList{
						v1.ResourceCPU: defaultResourceRequirements.Requests[v1.ResourceCPU],
					},
				},
			},
		},
	},
	{
		"container with requested memory",
		[]v1.Container{
			v1.Container{
				Resources: v1.ResourceRequirements{
					Requests: v1.ResourceList{
						v1.ResourceMemory: defaultResourceRequirements.Requests[v1.ResourceMemory],
					},
				},
			},
		},
		[]Patch{
			Patch{"replace", "/spec/containers/0/resources",
				v1.ResourceRequirements{
					Limits: v1.ResourceList{
						v1.ResourceCPU: defaultResourceRequirements.Limits[v1.ResourceCPU],
					},
					Requests: v1.ResourceList{
						v1.ResourceMemory: defaultResourceRequirements.Requests[v1.ResourceMemory],
						v1.ResourceCPU:    defaultResourceRequirements.Requests[v1.ResourceCPU],
					},
				},
			},
		},
	},
	{
		"container with limited cpu",
		[]v1.Container{
			v1.Container{
				Resources: v1.ResourceRequirements{
					Limits: v1.ResourceList{
						v1.ResourceCPU: defaultResourceRequirements.Limits[v1.ResourceCPU],
					},
				},
			},
		},
		[]Patch{
			Patch{"replace", "/spec/containers/0/resources",
				v1.ResourceRequirements{
					Limits: v1.ResourceList{
						v1.ResourceMemory: defaultResourceRequirements.Limits[v1.ResourceMemory],
						v1.ResourceCPU:    defaultResourceRequirements.Limits[v1.ResourceCPU],
					},
					Requests: v1.ResourceList{
						v1.ResourceMemory: defaultResourceRequirements.Requests[v1.ResourceMemory],
					},
				},
			},
		},
	},
	{
		"container with limited memory and cpu",
		[]v1.Container{
			v1.Container{
				Resources: v1.ResourceRequirements{
					Limits: v1.ResourceList{
						v1.ResourceMemory: defaultResourceRequirements.Limits[v1.ResourceMemory],
						v1.ResourceCPU:    defaultResourceRequirements.Limits[v1.ResourceCPU],
					},
				},
			},
		},
		[]Patch{},
	},
	{
		"container with requested memory and cpu",
		[]v1.Container{
			v1.Container{
				Resources: v1.ResourceRequirements{
					Requests: v1.ResourceList{
						v1.ResourceMemory: defaultResourceRequirements.Requests[v1.ResourceMemory],
						v1.ResourceCPU:    defaultResourceRequirements.Requests[v1.ResourceCPU],
					},
				},
			},
		},
		[]Patch{},
	},
	{
		"container with limited memory & requested cpu",
		[]v1.Container{
			v1.Container{
				Resources: v1.ResourceRequirements{
					Limits: v1.ResourceList{
						v1.ResourceMemory: defaultResourceRequirements.Limits[v1.ResourceMemory],
					},
					Requests: v1.ResourceList{
						v1.ResourceCPU: defaultResourceRequirements.Requests[v1.ResourceCPU],
					},
				},
			},
		},
		[]Patch{},
	},
	{
		"container with limited cpu & requested memory",
		[]v1.Container{
			v1.Container{
				Resources: v1.ResourceRequirements{
					Limits: v1.ResourceList{
						v1.ResourceCPU: defaultResourceRequirements.Limits[v1.ResourceCPU],
					},
					Requests: v1.ResourceList{
						v1.ResourceMemory: defaultResourceRequirements.Requests[v1.ResourceMemory],
					},
				},
			},
		},
		[]Patch{},
	},
}

var multiContainerPodTests = []struct {
	name string
	in   []v1.Container
	out  []Patch
}{
	{
		"two container without ResourceRequirements",
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
						v1.ResourceMemory: defaultResourceRequirements.Limits[v1.ResourceMemory],
						v1.ResourceCPU:    defaultResourceRequirements.Limits[v1.ResourceCPU],
					},
					Requests: v1.ResourceList{
						v1.ResourceMemory: defaultResourceRequirements.Requests[v1.ResourceMemory],
						v1.ResourceCPU:    defaultResourceRequirements.Requests[v1.ResourceCPU],
					},
				},
			},
			Patch{"replace", "/spec/containers/1/resources",
				v1.ResourceRequirements{
					Limits: v1.ResourceList{
						v1.ResourceMemory: defaultResourceRequirements.Limits[v1.ResourceMemory],
						v1.ResourceCPU:    defaultResourceRequirements.Limits[v1.ResourceCPU],
					},
					Requests: v1.ResourceList{
						v1.ResourceMemory: defaultResourceRequirements.Requests[v1.ResourceMemory],
						v1.ResourceCPU:    defaultResourceRequirements.Requests[v1.ResourceCPU],
					},
				},
			},
		},
	},
	{
		"fst container with limited cpu and requested memory, snd container without ResourceRequirements",
		[]v1.Container{
			v1.Container{
				Resources: v1.ResourceRequirements{
					Limits: v1.ResourceList{
						v1.ResourceCPU: defaultResourceRequirements.Limits[v1.ResourceCPU],
					},
					Requests: v1.ResourceList{
						v1.ResourceMemory: defaultResourceRequirements.Requests[v1.ResourceMemory],
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
						v1.ResourceMemory: defaultResourceRequirements.Limits[v1.ResourceMemory],
						v1.ResourceCPU:    defaultResourceRequirements.Limits[v1.ResourceCPU],
					},
					Requests: v1.ResourceList{
						v1.ResourceMemory: defaultResourceRequirements.Requests[v1.ResourceMemory],
						v1.ResourceCPU:    defaultResourceRequirements.Requests[v1.ResourceCPU],
					},
				},
			},
		},
	},
}

func TestSingleContainerPodPatches(t *testing.T) {
	for i, tt := range singleContainerPodTests {
		t.Run(strconv.Itoa(i)+"_"+tt.name, func(t *testing.T) {
			r, _ := podPatches(tt.in, defaultResourceRequirements)
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
		t.Run(strconv.Itoa(i)+"_"+tt.name, func(t *testing.T) {
			r, _ := podPatches(tt.in, defaultResourceRequirements)
			if !reflect.DeepEqual(r, tt.out) {
				t.Errorf("got %q, want %q", r, tt.out)
			}
			// if prettyPrint(r) != prettyPrint(tt.out) {
			// 	t.Errorf("got %q, want %q", prettyPrint(r), prettyPrint(tt.out))
			// }
		})
	}
}
