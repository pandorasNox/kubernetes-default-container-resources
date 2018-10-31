package webhook

import (
	"encoding/json"
	"reflect"
	"strconv"
	"testing"

	"github.com/pandorasnox/kubernetes-default-container-resources/pkg/container"
	k8s_v1 "k8s.io/api/core/v1"
)

var limitMemory = "1G"
var limitCPU = "0.5"
var requestMemory = "512M"
var requestCPU = "0.05"

var defaultResourceRequirements, _ = ParseResourceRequirements(limitMemory, limitCPU, requestMemory, requestCPU)

func prettyPrint(i interface{}) string {
	// s, _ := json.MarshalIndent(i, "", "\t")
	s, _ := json.Marshal(i)
	return string(s)
}

var singleContainerPodTests = []struct {
	name string
	in   []k8s_v1.Container
	out  []container.Patch
}{
	{
		"container without ResourceRequirements",
		[]k8s_v1.Container{
			k8s_v1.Container{
				Resources: k8s_v1.ResourceRequirements{},
			},
		},
		[]container.Patch{
			container.Patch{"replace", "/spec/containers/0/resources",
				k8s_v1.ResourceRequirements{
					Limits: k8s_v1.ResourceList{
						k8s_v1.ResourceMemory: defaultResourceRequirements.Limits[k8s_v1.ResourceMemory],
						k8s_v1.ResourceCPU:    defaultResourceRequirements.Limits[k8s_v1.ResourceCPU],
					},
					Requests: k8s_v1.ResourceList{
						k8s_v1.ResourceMemory: defaultResourceRequirements.Requests[k8s_v1.ResourceMemory],
						k8s_v1.ResourceCPU:    defaultResourceRequirements.Requests[k8s_v1.ResourceCPU],
					},
				},
			},
		},
	},
	{
		"container with limited memory",
		[]k8s_v1.Container{
			k8s_v1.Container{
				Resources: k8s_v1.ResourceRequirements{
					Limits: k8s_v1.ResourceList{
						k8s_v1.ResourceMemory: defaultResourceRequirements.Limits[k8s_v1.ResourceMemory],
					},
				},
			},
		},
		[]container.Patch{
			container.Patch{"replace", "/spec/containers/0/resources",
				k8s_v1.ResourceRequirements{
					Limits: k8s_v1.ResourceList{
						k8s_v1.ResourceMemory: defaultResourceRequirements.Limits[k8s_v1.ResourceMemory],
						k8s_v1.ResourceCPU:    defaultResourceRequirements.Limits[k8s_v1.ResourceCPU],
					},
					Requests: k8s_v1.ResourceList{
						k8s_v1.ResourceCPU: defaultResourceRequirements.Requests[k8s_v1.ResourceCPU],
					},
				},
			},
		},
	},
	{
		"container with requested memory",
		[]k8s_v1.Container{
			k8s_v1.Container{
				Resources: k8s_v1.ResourceRequirements{
					Requests: k8s_v1.ResourceList{
						k8s_v1.ResourceMemory: defaultResourceRequirements.Requests[k8s_v1.ResourceMemory],
					},
				},
			},
		},
		[]container.Patch{
			container.Patch{"replace", "/spec/containers/0/resources",
				k8s_v1.ResourceRequirements{
					Limits: k8s_v1.ResourceList{
						k8s_v1.ResourceCPU: defaultResourceRequirements.Limits[k8s_v1.ResourceCPU],
					},
					Requests: k8s_v1.ResourceList{
						k8s_v1.ResourceMemory: defaultResourceRequirements.Requests[k8s_v1.ResourceMemory],
						k8s_v1.ResourceCPU:    defaultResourceRequirements.Requests[k8s_v1.ResourceCPU],
					},
				},
			},
		},
	},
	{
		"container with limited cpu",
		[]k8s_v1.Container{
			k8s_v1.Container{
				Resources: k8s_v1.ResourceRequirements{
					Limits: k8s_v1.ResourceList{
						k8s_v1.ResourceCPU: defaultResourceRequirements.Limits[k8s_v1.ResourceCPU],
					},
				},
			},
		},
		[]container.Patch{
			container.Patch{"replace", "/spec/containers/0/resources",
				k8s_v1.ResourceRequirements{
					Limits: k8s_v1.ResourceList{
						k8s_v1.ResourceMemory: defaultResourceRequirements.Limits[k8s_v1.ResourceMemory],
						k8s_v1.ResourceCPU:    defaultResourceRequirements.Limits[k8s_v1.ResourceCPU],
					},
					Requests: k8s_v1.ResourceList{
						k8s_v1.ResourceMemory: defaultResourceRequirements.Requests[k8s_v1.ResourceMemory],
					},
				},
			},
		},
	},
	{
		"container with limited memory and cpu",
		[]k8s_v1.Container{
			k8s_v1.Container{
				Resources: k8s_v1.ResourceRequirements{
					Limits: k8s_v1.ResourceList{
						k8s_v1.ResourceMemory: defaultResourceRequirements.Limits[k8s_v1.ResourceMemory],
						k8s_v1.ResourceCPU:    defaultResourceRequirements.Limits[k8s_v1.ResourceCPU],
					},
				},
			},
		},
		[]container.Patch{},
	},
	{
		"container with requested memory and cpu",
		[]k8s_v1.Container{
			k8s_v1.Container{
				Resources: k8s_v1.ResourceRequirements{
					Requests: k8s_v1.ResourceList{
						k8s_v1.ResourceMemory: defaultResourceRequirements.Requests[k8s_v1.ResourceMemory],
						k8s_v1.ResourceCPU:    defaultResourceRequirements.Requests[k8s_v1.ResourceCPU],
					},
				},
			},
		},
		[]container.Patch{},
	},
	{
		"container with limited memory & requested cpu",
		[]k8s_v1.Container{
			k8s_v1.Container{
				Resources: k8s_v1.ResourceRequirements{
					Limits: k8s_v1.ResourceList{
						k8s_v1.ResourceMemory: defaultResourceRequirements.Limits[k8s_v1.ResourceMemory],
					},
					Requests: k8s_v1.ResourceList{
						k8s_v1.ResourceCPU: defaultResourceRequirements.Requests[k8s_v1.ResourceCPU],
					},
				},
			},
		},
		[]container.Patch{},
	},
	{
		"container with limited cpu & requested memory",
		[]k8s_v1.Container{
			k8s_v1.Container{
				Resources: k8s_v1.ResourceRequirements{
					Limits: k8s_v1.ResourceList{
						k8s_v1.ResourceCPU: defaultResourceRequirements.Limits[k8s_v1.ResourceCPU],
					},
					Requests: k8s_v1.ResourceList{
						k8s_v1.ResourceMemory: defaultResourceRequirements.Requests[k8s_v1.ResourceMemory],
					},
				},
			},
		},
		[]container.Patch{},
	},
}

var multiContainerPodTests = []struct {
	name string
	in   []k8s_v1.Container
	out  []container.Patch
}{
	{
		"two container without ResourceRequirements",
		[]k8s_v1.Container{
			k8s_v1.Container{
				Resources: k8s_v1.ResourceRequirements{},
			},
			k8s_v1.Container{
				Resources: k8s_v1.ResourceRequirements{},
			},
		},
		[]container.Patch{
			container.Patch{"replace", "/spec/containers/0/resources",
				k8s_v1.ResourceRequirements{
					Limits: k8s_v1.ResourceList{
						k8s_v1.ResourceMemory: defaultResourceRequirements.Limits[k8s_v1.ResourceMemory],
						k8s_v1.ResourceCPU:    defaultResourceRequirements.Limits[k8s_v1.ResourceCPU],
					},
					Requests: k8s_v1.ResourceList{
						k8s_v1.ResourceMemory: defaultResourceRequirements.Requests[k8s_v1.ResourceMemory],
						k8s_v1.ResourceCPU:    defaultResourceRequirements.Requests[k8s_v1.ResourceCPU],
					},
				},
			},
			container.Patch{"replace", "/spec/containers/1/resources",
				k8s_v1.ResourceRequirements{
					Limits: k8s_v1.ResourceList{
						k8s_v1.ResourceMemory: defaultResourceRequirements.Limits[k8s_v1.ResourceMemory],
						k8s_v1.ResourceCPU:    defaultResourceRequirements.Limits[k8s_v1.ResourceCPU],
					},
					Requests: k8s_v1.ResourceList{
						k8s_v1.ResourceMemory: defaultResourceRequirements.Requests[k8s_v1.ResourceMemory],
						k8s_v1.ResourceCPU:    defaultResourceRequirements.Requests[k8s_v1.ResourceCPU],
					},
				},
			},
		},
	},
	{
		"fst container with limited cpu and requested memory snd container without ResourceRequirements",
		[]k8s_v1.Container{
			k8s_v1.Container{
				Resources: k8s_v1.ResourceRequirements{
					Limits: k8s_v1.ResourceList{
						k8s_v1.ResourceCPU: defaultResourceRequirements.Limits[k8s_v1.ResourceCPU],
					},
					Requests: k8s_v1.ResourceList{
						k8s_v1.ResourceMemory: defaultResourceRequirements.Requests[k8s_v1.ResourceMemory],
					},
				},
			},
			k8s_v1.Container{
				Resources: k8s_v1.ResourceRequirements{},
			},
		},
		[]container.Patch{
			container.Patch{"replace", "/spec/containers/1/resources",
				k8s_v1.ResourceRequirements{
					Limits: k8s_v1.ResourceList{
						k8s_v1.ResourceMemory: defaultResourceRequirements.Limits[k8s_v1.ResourceMemory],
						k8s_v1.ResourceCPU:    defaultResourceRequirements.Limits[k8s_v1.ResourceCPU],
					},
					Requests: k8s_v1.ResourceList{
						k8s_v1.ResourceMemory: defaultResourceRequirements.Requests[k8s_v1.ResourceMemory],
						k8s_v1.ResourceCPU:    defaultResourceRequirements.Requests[k8s_v1.ResourceCPU],
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
