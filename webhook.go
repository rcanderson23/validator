package main

import (
	"encoding/json"
	"io/ioutil"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"net/http"
	//"k8s.io/klog"
	corev1 "k8s.io/api/core/v1"
	"log"
	"regexp"
)

var (
	universalDeserializer = serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer()
)

func validateFuncHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failure to read body: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var adReview admissionv1beta1.AdmissionReview

	if _, _, err := universalDeserializer.Decode(body, nil, &adReview); err != nil {
		log.Printf("unable to decode: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if adReview.Request == nil {
		log.Printf("Request is nil")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	adResponse := validate(adReview.Request)
	bytes, err := json.Marshal(&adResponse)
	_, writeErr := w.Write(bytes)
	if writeErr != nil {
		log.Printf("Could not write response: %v", writeErr)
	}
}

func validate(req *admissionv1beta1.AdmissionRequest) admissionv1beta1.AdmissionReview {
	//setup AdmissionReview. Default to allow and then reject based on conditions
	responseReview := admissionv1beta1.AdmissionReview{
		Response: &admissionv1beta1.AdmissionResponse{
			UID:     req.UID,
			Allowed: true,
		},
	}
	if req.Kind.Kind != "Pod" {

		return responseReview
	}
	pod := corev1.Pod{}

	if _, _, err := universalDeserializer.Decode(req.Object.Raw, nil, &pod); err != nil {
		log.Printf("Couldn't decode object: %v", err)
	}

	containers := pod.Spec.Containers
	for i := 0; i < len(containers); i++ {
		match, _ := regexp.MatchString(`^(\d{12})\.dkr\.ecr\.((\D{2})\-(\D*)\-(\d))\.amazonaws\.com/`, pod.Spec.Containers[i].Image)
		if !match {
			log.Printf("Regex does not match image, admission rejected.")
			responseReview.Response.Allowed = false
			responseReview.Response.Result = &metav1.Status{
				Message: "Regex does not match image, admission rejected",
			}
		}
	}
	return responseReview
}
