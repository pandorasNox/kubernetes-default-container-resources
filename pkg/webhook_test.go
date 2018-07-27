package webhook

import (
	"reflect"
	"strconv"
	"testing"
)

var singeContainerPodTests = []struct {
	in  []Container
	out []Patch
}{
	{
		[]Container{
			Container{
				Resources: ComputeResources{},
			},
		},
		// []Patch{
		// 	Patch{"add", "/spec/containers/0/resources/limits/memory", "1G"},
		// 	Patch{"add", "/spec/containers/0/resources/requests/memory", "1G"},
		// 	Patch{"add", "/spec/containers/0/resources/limits/cpu", "0.5"},
		// 	Patch{"add", "/spec/containers/0/resources/requests/cpu", "0.1"},
		// },
		[]Patch{
			Patch{"replace", "/spec/containers/0/resources",
				ComputeResources{
					Limits: ComputeUnit{
						Memory: "1G",
						CPU:    "0.5",
					},
					Requests: ComputeUnit{
						Memory: "1G",
						CPU:    "0.1",
					},
				},
			},
		},
	},
	{
		[]Container{
			Container{
				Resources: ComputeResources{
					Limits: ComputeUnit{
						Memory: "1G",
					},
				},
			},
		},
		// []Patch{
		// 	Patch{"add", "/spec/containers/0/resources/limits/cpu", "0.5"},
		// 	Patch{"add", "/spec/containers/0/resources/requests/cpu", "0.1"},
		// },
		[]Patch{
			Patch{"replace", "/spec/containers/0/resources",
				ComputeResources{
					Limits: ComputeUnit{
						Memory: "1G",
						CPU:    "0.5",
					},
					Requests: ComputeUnit{
						CPU: "0.1",
					},
				},
			},
		},
	},
	{
		[]Container{
			Container{
				Resources: ComputeResources{
					Requests: ComputeUnit{
						Memory: "1G",
					},
				},
			},
		},
		// []Patch{
		// 	Patch{"add", "/spec/containers/0/resources/limits/cpu", "0.5"},
		// 	Patch{"add", "/spec/containers/0/resources/requests/cpu", "0.1"},
		// },
		[]Patch{
			Patch{"replace", "/spec/containers/0/resources",
				ComputeResources{
					Limits: ComputeUnit{
						CPU: "0.5",
					},
					Requests: ComputeUnit{
						Memory: "1G",
						CPU:    "0.1",
					},
				},
			},
		},
	},
	{
		[]Container{
			Container{
				Resources: ComputeResources{
					Limits: ComputeUnit{
						CPU: "0.5",
					},
				},
			},
		},
		// []Patch{
		// 	Patch{"add", "/spec/containers/0/resources/limits/memory", "1G"},
		// 	Patch{"add", "/spec/containers/0/resources/requests/memory", "1G"},
		// },
		[]Patch{
			Patch{"replace", "/spec/containers/0/resources",
				ComputeResources{
					Limits: ComputeUnit{
						Memory: "1G",
						CPU:    "0.5",
					},
					Requests: ComputeUnit{
						Memory: "1G",
					},
				},
			},
		},
	},
	{
		[]Container{
			Container{
				Resources: ComputeResources{
					Limits: ComputeUnit{
						Memory: "1G",
						CPU:    "0.5",
					},
				},
			},
		},
		[]Patch{},
	},
	{
		[]Container{
			Container{
				Resources: ComputeResources{
					Requests: ComputeUnit{
						Memory: "1G",
						CPU:    "0.5",
					},
				},
			},
		},
		[]Patch{},
	},
	{
		[]Container{
			Container{
				Resources: ComputeResources{
					Limits: ComputeUnit{
						Memory: "1G",
					},
					Requests: ComputeUnit{
						CPU: "0.5",
					},
				},
			},
		},
		[]Patch{},
	},
	{
		[]Container{
			Container{
				Resources: ComputeResources{
					Limits: ComputeUnit{
						CPU: "0.5",
					},
					Requests: ComputeUnit{
						Memory: "1G",
					},
				},
			},
		},
		[]Patch{},
	},
}

var multiContainerPodTests = []struct {
	in  []Container
	out []Patch
}{
	{
		[]Container{
			Container{
				Resources: ComputeResources{},
			},
			Container{
				Resources: ComputeResources{},
			},
		},
		// []Patch{
		// 	Patch{"add", "/spec/containers/0/resources/limits/memory", "1G"},
		// 	Patch{"add", "/spec/containers/0/resources/requests/memory", "1G"},
		// 	Patch{"add", "/spec/containers/0/resources/limits/cpu", "0.5"},
		// 	Patch{"add", "/spec/containers/0/resources/requests/cpu", "0.1"},
		// 	Patch{"add", "/spec/containers/1/resources/limits/memory", "1G"},
		// 	Patch{"add", "/spec/containers/1/resources/requests/memory", "1G"},
		// 	Patch{"add", "/spec/containers/1/resources/limits/cpu", "0.5"},
		// 	Patch{"add", "/spec/containers/1/resources/requests/cpu", "0.1"},
		// },
		[]Patch{
			Patch{"replace", "/spec/containers/0/resources",
				ComputeResources{
					Limits: ComputeUnit{
						Memory: "1G",
						CPU:    "0.5",
					},
					Requests: ComputeUnit{
						Memory: "1G",
						CPU:    "0.1",
					},
				},
			},
			Patch{"replace", "/spec/containers/1/resources",
				ComputeResources{
					Limits: ComputeUnit{
						Memory: "1G",
						CPU:    "0.5",
					},
					Requests: ComputeUnit{
						Memory: "1G",
						CPU:    "0.1",
					},
				},
			},
		},
	},
	{
		[]Container{
			Container{
				Resources: ComputeResources{
					Limits: ComputeUnit{
						CPU: "0.5",
					},
					Requests: ComputeUnit{
						Memory: "1G",
					},
				},
			},
			Container{
				Resources: ComputeResources{},
			},
		},
		// []Patch{
		// 	Patch{"add", "/spec/containers/1/resources/limits/memory", "1G"},
		// 	Patch{"add", "/spec/containers/1/resources/requests/memory", "1G"},
		// 	Patch{"add", "/spec/containers/1/resources/limits/cpu", "0.5"},
		// 	Patch{"add", "/spec/containers/1/resources/requests/cpu", "0.1"},
		// },
		[]Patch{
			Patch{"replace", "/spec/containers/1/resources",
				ComputeResources{
					Limits: ComputeUnit{
						Memory: "1G",
						CPU:    "0.5",
					},
					Requests: ComputeUnit{
						Memory: "1G",
						CPU:    "0.1",
					},
				},
			},
		},
	},
}

func TestPodPatches(t *testing.T) {
	for i, tt := range singeContainerPodTests {
		t.Run(""+strconv.Itoa(i), func(t *testing.T) {
			r := podPatches(tt.in, "1G", "0.5", "1G", "0.1")
			if !reflect.DeepEqual(r, tt.out) {
				t.Errorf("got %q, want %q", r, tt.out)
			}
		})
	}

	for i, tt := range multiContainerPodTests {
		t.Run(""+strconv.Itoa(i), func(t *testing.T) {
			r := podPatches(tt.in, "1G", "0.5", "1G", "0.1")
			if !reflect.DeepEqual(r, tt.out) {
				t.Errorf("got %q, want %q", r, tt.out)
			}
		})
	}
}
