package webhook

import (
	"reflect"
	"strconv"
	"testing"
)

// func TestAverage(t *testing.T) {
// 	fmt.Println("", Patch{})

// 	containers := Containers{
// 		struct {
// 			Name      string
// 			Env       []Envs
// 			Resources ComputeResources
// 		}{
// 			Name:      "myName",
// 			Env:       []Envs{Envs{}},
// 			Resources: ComputeResources{},
// 		},
// 	}

// 	admissionReview, err := getAdmissionReview(containers, "ABC123", "64Mi", "0.5", "32Mi", "0.1")
// 	if err != nil {
// 		t.Error("couldn't got getAdmissionReview")
// 	}

// 	fmt.Println(admissionReview)

// 	expectedAdmissionReview := AdmissionReview{
// 		Response: AdmissionReviewResponse{
// 			UID:       "ABC123",
// 			Allowed:   true,
// 			Status:    AdmissionStatus{},
// 			Patch:     "",
// 			PatchType: "JSONPatch",
// 		},
// 	}

// 	fmt.Println(expectedAdmissionReview)

// 	var v = 1.5
// 	if v != 1.5 {
// 		t.Error("Expected 1.5, got ", v)
// 	}
// }

type defChanInput struct {
	patches          []Patch
	computeResources ComputeResources
}

var definePatchesTests = []struct {
	in  defChanInput
	out []Patch
}{
	{
		defChanInput{
			[]Patch{},
			ComputeResources{},
		},
		[]Patch{
			Patch{"add", "/spec/containers/0/resources/limits/memory", "1G"},
			Patch{"add", "/spec/containers/0/resources/requests/memory", "1G"},
			Patch{"add", "/spec/containers/0/resources/limits/cpu", "0.5"},
			Patch{"add", "/spec/containers/0/resources/requests/cpu", "0.1"},
		},
	},
	{
		defChanInput{
			[]Patch{},
			ComputeResources{
				Limits: ComputeUnit{
					Memory: "1G",
				},
			},
		},
		[]Patch{
			Patch{"add", "/spec/containers/1/resources/limits/cpu", "0.5"},
			Patch{"add", "/spec/containers/1/resources/requests/cpu", "0.1"},
		},
	},
	{
		defChanInput{
			[]Patch{},
			ComputeResources{
				Requests: ComputeUnit{
					Memory: "1G",
				},
			},
		},
		[]Patch{
			Patch{"add", "/spec/containers/2/resources/limits/cpu", "0.5"},
			Patch{"add", "/spec/containers/2/resources/requests/cpu", "0.1"},
		},
	},
	{
		defChanInput{
			[]Patch{},
			ComputeResources{
				Limits: ComputeUnit{
					CPU: "0.5",
				},
			},
		},
		[]Patch{
			Patch{"add", "/spec/containers/3/resources/limits/memory", "1G"},
			Patch{"add", "/spec/containers/3/resources/requests/memory", "1G"},
		},
	},
	{
		defChanInput{
			[]Patch{},
			ComputeResources{
				Limits: ComputeUnit{
					Memory: "1G",
					CPU:    "0.5",
				},
			},
		},
		[]Patch{},
	},
	{
		defChanInput{
			[]Patch{},
			ComputeResources{
				Requests: ComputeUnit{
					Memory: "1G",
					CPU:    "0.5",
				},
			},
		},
		[]Patch{},
	},
	{
		defChanInput{
			[]Patch{},
			ComputeResources{
				Limits: ComputeUnit{
					Memory: "1G",
				},
				Requests: ComputeUnit{
					CPU: "0.5",
				},
			},
		},
		[]Patch{},
	},
	{
		defChanInput{
			[]Patch{},
			ComputeResources{
				Limits: ComputeUnit{
					CPU: "0.5",
				},
				Requests: ComputeUnit{
					Memory: "1G",
				},
			},
		},
		[]Patch{},
	},
}

func TestDefineChanges(t *testing.T) {
	for i, tt := range definePatchesTests {
		t.Run(""+strconv.Itoa(i), func(t *testing.T) {
			r := definePatches(tt.in.patches, i, tt.in.computeResources, "1G", "0.5", "1G", "0.1")
			if !reflect.DeepEqual(r, tt.out) {
				t.Errorf("got %q, want %q", r, tt.out)
			}
		})
	}
}
