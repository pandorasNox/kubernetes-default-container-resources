package webhook

import (
	"reflect"
	"testing"

	k8s_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

var limitMemory = "1G"
var limitCPU = "0.5"
var requestMemory = "1G"
var requestCPU = "0.1"

var defaults = parseTestResourceRequirements(limitMemory, limitCPU, requestMemory, requestCPU)

func getResourceQuantity(quantity string) resource.Quantity {
	resDef, err := resource.ParseQuantity(quantity)
	if err != nil {
		panic(err)
	}
	return resDef
}

func parseTestResourceRequirements(memoryLimit, CPULimit, memoryRequest, CPURequest string) k8s_v1.ResourceRequirements {
	r := k8s_v1.ResourceRequirements{}

	if (memoryLimit != "" || CPULimit != "") && r.Limits == nil {
		r.Limits = k8s_v1.ResourceList{}
	}
	if (memoryRequest != "" || CPURequest != "") && r.Requests == nil {
		r.Requests = k8s_v1.ResourceList{}
	}

	if memoryLimit != "" {
		r.Limits[k8s_v1.ResourceMemory] = getResourceQuantity(memoryLimit)
	}
	if CPULimit != "" {
		r.Limits[k8s_v1.ResourceCPU] = getResourceQuantity(CPULimit)
	}
	if memoryRequest != "" {
		r.Requests[k8s_v1.ResourceMemory] = getResourceQuantity(memoryRequest)
	}
	if CPURequest != "" {
		r.Requests[k8s_v1.ResourceCPU] = getResourceQuantity(CPURequest)
	}

	return r
}

func Test_addDefaults(t *testing.T) {
	type args struct {
		c k8s_v1.ResourceRequirements
		d k8s_v1.ResourceRequirements
	}
	tests := []struct {
		name    string
		args    args
		want    k8s_v1.ResourceRequirements
		wantErr bool
	}{
		{
			name: "empty ResourceRequirements sets default",
			args: args{
				c: k8s_v1.ResourceRequirements{},
				d: defaults,
			},
			want: defaults,
		},
		{
			name: "empty Limits & Requests sets default",
			args: args{
				c: k8s_v1.ResourceRequirements{
					Limits:   k8s_v1.ResourceList{},
					Requests: k8s_v1.ResourceList{},
				},
				d: defaults,
			},
			want: defaults,
		},
		{
			name: "empty Requests sets default",
			args: args{
				c: k8s_v1.ResourceRequirements{
					Limits: k8s_v1.ResourceList{},
				},
				d: defaults,
			},
			want: defaults,
		},
		{
			name: "empty Limits sets default",
			args: args{
				c: k8s_v1.ResourceRequirements{
					Limits: k8s_v1.ResourceList{},
				},
				d: defaults,
			},
			want: defaults,
		},
		{
			name: "Limits mem equal to default limit mem",
			args: args{
				c: parseTestResourceRequirements("1G", "", "", ""),
				d: defaults,
			},
			want: defaults,
		},
		{
			name: "Limits mem greater than default limit mem",
			args: args{
				c: parseTestResourceRequirements("2G", "", "", ""),
				d: defaults,
			},
			want: parseTestResourceRequirements("2G", limitCPU, requestMemory, requestCPU),
		},
		{
			name: "Limits mem smaller than default limit mem and greate than default requets mem",
			args: args{
				c: parseTestResourceRequirements("512M", "", "", ""),
				d: parseTestResourceRequirements(limitMemory, limitCPU, "256M", requestCPU),
			},
			want: parseTestResourceRequirements("512M", limitCPU, "256M", requestCPU),
		},
		{
			name: "Limits mem smaller than default limit mem",
			args: args{
				c: parseTestResourceRequirements("512M", "", "", ""),
				d: defaults,
			},
			want:    parseTestResourceRequirements("512M", limitCPU, requestMemory, requestCPU),
			wantErr: true,
		},
		{
			name: "Limits cpu equal to default limit cpu",
			args: args{
				c: parseTestResourceRequirements("", "0.5", "", ""),
				d: defaults,
			},
			want: defaults,
		},
		{
			name: "Limits cpu greater than default limit cpu",
			args: args{
				c: parseTestResourceRequirements("", "0.6", "", ""),
				d: defaults,
			},
			want: parseTestResourceRequirements(limitMemory, "0.6", requestMemory, requestCPU),
		},
		{
			name: "Limits cpu smaller than default limit cpu",
			args: args{
				c: parseTestResourceRequirements("", "0.4", "", ""),
				d: defaults,
			},
			want: parseTestResourceRequirements(limitMemory, "0.4", requestMemory, requestCPU),
		},
		{
			name: "Limits cpu smaller than default requests cpu",
			args: args{
				c: parseTestResourceRequirements("", "0.05", "", ""),
				d: defaults,
			},
			want:    parseTestResourceRequirements(limitMemory, "0.05", requestMemory, requestCPU),
			wantErr: true,
		},

		{
			name: "set higher CPU",
			args: args{
				c: parseTestResourceRequirements("", "1000m", "", ""),
				d: defaults,
			},
			want: parseTestResourceRequirements(limitMemory, "1000m", requestMemory, requestCPU),
		},
		{
			name: "get error on to low CPU limit",
			args: args{
				c: parseTestResourceRequirements("", "1m", "", ""),
				d: defaults,
			},
			want:    parseTestResourceRequirements(limitMemory, "1m", requestMemory, requestCPU),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := addDefaults(tt.args.c, tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("addDefaults() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addDefaults() = %v, want %v", got, tt.want)
			}
		})
	}
}
