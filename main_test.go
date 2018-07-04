package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDenyNonDeploymentResources(t *testing.T) {
	podResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	aReq := v1beta1.AdmissionReview{Request: &v1beta1.AdmissionRequest{Resource: podResource}}

	aResp := CheckDeploymentForLimits(aReq)

	assert.False(t, aResp.Allowed)
}
