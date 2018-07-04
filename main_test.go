package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/api/admission/v1beta1"
	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestDenyNonDeploymentResources(t *testing.T) {
	podResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	aReq := v1beta1.AdmissionReview{Request: &v1beta1.AdmissionRequest{Resource: podResource}}

	aResp := CheckDeploymentForLimits(aReq)

	assert.False(t, aResp.Allowed)
}

func TestDenyDeploymentsWithoutLimits(t *testing.T) {
	deploymentResource := metav1.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}

	raw, _ := json.Marshal(aDeploymentWithOutLimits())

	aReq := v1beta1.AdmissionReview{Request: &v1beta1.AdmissionRequest{Resource: deploymentResource, Object: runtime.RawExtension{Raw: raw}}}
	aResp := CheckDeploymentForLimits(aReq)

	assert.False(t, aResp.Allowed)
}

func TestAllowDeploymentsWithLimits(t *testing.T) {
	deploymentResource := metav1.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}

	raw, _ := json.Marshal(aDeploymentWithLimits())

	aReq := v1beta1.AdmissionReview{Request: &v1beta1.AdmissionRequest{Resource: deploymentResource, Object: runtime.RawExtension{Raw: raw}}}
	aResp := CheckDeploymentForLimits(aReq)

	assert.True(t, aResp.Allowed)
}

func aDeploymentWithOutLimits() *v1.Deployment {
	return &v1.Deployment{
		Spec: v1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{}}}}},
	}
}

func aDeploymentWithLimits() *v1.Deployment {
	return &v1.Deployment{
		Spec: v1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.Quantity{},
							}}}}}}},
	}
}
