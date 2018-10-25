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

func TestDebug(t *testing.T) {
	// fmt.Println("")
	// incommingAdmissionReviewStr := `{"kind":"AdmissionReview","apiVersion":"admission.k8s.io/v1beta1","request":{"uid":"9e9322a0-85f0-11e8-b78d-080027d242b4","kind":{"group":"","version":"v1","kind":"Pod"},"resource":{"group":"","version":"v1","resource":"pods"},"namespace":"foo","operation":"CREATE","userInfo":{"username":"system:serviceaccount:kube-system:replicaset-controller","uid":"1df83e81-85e6-11e8-b78d-080027d242b4","groups":["system:serviceaccounts","system:serviceaccounts:kube-system","system:authenticated"]},"object":{"metadata":{"name":"nginx-deployment-75675f5897-gsbdl","generateName":"nginx-deployment-75675f5897-","namespace":"foo","uid":"9e931740-85f0-11e8-b78d-080027d242b4","creationTimestamp":"2018-07-12T16:28:33Z","labels":{"app":"nginx","pod-template-hash":"3123191453"},"ownerReferences":[{"apiVersion":"extensions/v1beta1","kind":"ReplicaSet","name":"nginx-deployment-75675f5897","uid":"9e8c0c5e-85f0-11e8-b78d-080027d242b4","controller":true,"blockOwnerDeletion":true}]},"spec":{"volumes":[{"name":"default-token-cj7hv","secret":{"secretName":"default-token-cj7hv"}}],"containers":[{"name":"nginx","image":"nginx:1.7.9","env":[{"name":"foo","value":"bar"}],"ports":[{"containerPort":80,"protocol":"TCP"}],"volumeMounts":[{"name":"default-token-cj7hv","readOnly":true,"mountPath":"/var/run/secrets/kubernetes.io/serviceaccount"}],"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent"}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","serviceAccountName":"default","serviceAccount":"default","securityContext":{},"schedulerName":"default-scheduler","tolerations":[{"key":"node.kubernetes.io/not-ready","operator":"Exists","effect":"NoExecute","tolerationSeconds":300},{"key":"node.kubernetes.io/unreachable","operator":"Exists","effect":"NoExecute","tolerationSeconds":300}]},"status":{"phase":"Pending","qosClass":"BestEffort"}},"oldObject":null}}`
	// incommingAdmissionReviewReader := strings.NewReader(incommingAdmissionReviewStr)
	// // fmt.Println(incommingAdmissionReviewStr)
	// fmt.Println("")

	// incomingAdmissionReview := &v1beta1.AdmissionReview{}
	// json.NewDecoder(incommingAdmissionReviewReader).Decode(incomingAdmissionReview)

	// // fmt.Println("")
	// // fmt.Printf("show k8s incomingAdmissionReview: %s", prettyPrint(incomingAdmissionReview))
	// // fmt.Println("")

	// raw := incomingAdmissionReview.Request.Object.Raw
	// pod := v1.Pod{}
	// json.Unmarshal(raw, &pod)

	// // fmt.Println("")
	// // fmt.Printf("show k8s pod: %s", prettyPrint(pod))
	// // fmt.Println("")

	// containers := pod.Spec.Containers

	// // fmt.Println("")
	// // fmt.Printf("show k8s containers: %s", prettyPrint(containers))
	// // fmt.Println("")

	// fstCon := containers[0]

	// // fmt.Println("")
	// fmt.Printf("show k8s fstCon: %s", prettyPrint(fstCon))
	// // fmt.Println("")

	// RR := v1.ResourceRequirements{}
	// RR2 := v1.ResourceRequirements{}
	// fmt.Println("")
	// fmt.Println("RR.Limits == nil: ", RR.Limits == nil)
	// fmt.Println("")
	// if mapKeyExist(RR.Limits, v1.ResourceCPU) {
	// 	fmt.Println("RR.Limits[v1.ResourceCPU] - keyExist")
	// } else {
	// 	fmt.Println("RR.Limits[v1.ResourceCPU] - key NOT Exist")
	// }
	// fmt.Println("")
	// fmt.Println("deep equal: ", reflect.DeepEqual(RR, RR2))
	// fmt.Println("")
}
