package webhook_test

import (
	"fmt"
	"testing"

	webhook "github.com/pandorasnox/kubernetes-default-container-resources/pkg"
)

func TestAverage(t *testing.T) {
	fmt.Println("", webhook.Operation{})
	var v float64 = 1.7
	if v != 1.5 {
		t.Error("Expected 1.5, got ", v)
	}
}
