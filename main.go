package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	"github.com/golang/glog"
	"k8s.io/api/admission/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var codecs = serializer.NewCodecFactory(runtime.NewScheme())

func main() {
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServeTLS(":443", "/app/server.crt", "/app/server.key", nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		glog.Errorf("contentType=%s, expect application/json", contentType)
		return
	}

	var reviewResponse *v1beta1.AdmissionResponse
	ar := v1beta1.AdmissionReview{}

	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		glog.Error(err)
		reviewResponse = toAdmissionResponse(err)
	} else {
		reviewResponse = CheckDeploymentForLimits(ar)
	}

	response := v1beta1.AdmissionReview{}
	if reviewResponse != nil {
		response.Response = reviewResponse
		response.Response.UID = ar.Request.UID
	}
	// reset the Object and OldObject, they are not needed in a response.
	ar.Request.Object = runtime.RawExtension{}
	ar.Request.OldObject = runtime.RawExtension{}

	resp, err := json.Marshal(response)
	if err != nil {
		glog.Error(err)
	}
	if _, err := w.Write(resp); err != nil {
		glog.Error(err)
	}

}

func toAdmissionResponse(err error) *v1beta1.AdmissionResponse {
	return &v1beta1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}

// CheckDeploymentForLimits ensures that a reviewed deployment has limits set
func CheckDeploymentForLimits(aReq v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	deploymentResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "deployments"}
	if aReq.Request.Resource != deploymentResource {
		err := fmt.Errorf("expect resource to be %s", deploymentResource)
		glog.Error(err)
		return toAdmissionResponse(err)
	}

	raw := aReq.Request.Object.Raw
	deployment := appsv1.Deployment{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &deployment); err != nil {
		glog.Error(err)
		return toAdmissionResponse(err)
	}

	if hasNoLimits(deployment) {
		return toAdmissionResponse(fmt.Errorf("no resource limits set"))
	}
	return &v1beta1.AdmissionResponse{Allowed: true}
}

func hasNoLimits(d appsv1.Deployment) bool {
	for _, c := range d.Spec.Template.Spec.Containers {
		if reflect.DeepEqual(c.Resources, corev1.ResourceRequirements{}) {
			return true
		}
	}
	return false
}
