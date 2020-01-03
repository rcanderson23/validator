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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"log"
	"regexp"
)

var (
	universalDeserializer = serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer()
)

func (vs *ValidatorSpec) validateHandler(w http.ResponseWriter, r *http.Request) {
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
	adResponse := validate(adReview.Request, vs)
	bytes, err := json.Marshal(&adResponse)
	_, writeErr := w.Write(bytes)
	if writeErr != nil {
		log.Printf("Could not write response: %v", writeErr)
	}
}

func validate(req *admissionv1beta1.AdmissionRequest, vs *ValidatorSpec) admissionv1beta1.AdmissionReview {
	//setup AdmissionReview. Default to allow and then reject based on conditions
	responseReview := admissionv1beta1.AdmissionReview{
		Response: &admissionv1beta1.AdmissionResponse{
			UID:     req.UID,
			Allowed: true,
			Result: &metav1.Status{
				Message: "",
			},
		},
	}

	var allowed bool
	var message string
	switch req.Kind.Kind {
	case "Pod":
		if vs.Pod == nil {
			break
		}
		object := &corev1.Pod{}
		if _, _, err := universalDeserializer.Decode(req.Object.Raw, nil, object); err != nil {
			log.Printf("Couldn't decode object: %v", err)
		}
		log.Println("checking labels")
		log.Println(vs.Pod.Labels)
		if vs.Pod.Labels != nil {
			allowed, message = checkLabels(object.Labels, vs.Pod.Labels)
			if !allowed {
				responseReview.Response.Allowed = allowed
				responseReview.Response.Result.Message = message
				break
			}
		}

		if vs.Pod.Image != "" {
			allowed, message = checkContainers(&object.Spec, vs.Pod.Image)
			responseReview.Response.Allowed = allowed
			responseReview.Response.Result.Message = message
		}
	case "Deployment":
		object := &appsv1.Deployment{}
		if _, _, err := universalDeserializer.Decode(req.Object.Raw, nil, object); err != nil {
			log.Printf("Couldn't decode object: %v", err)
		}
		allowed, message = checkLabels(object.Labels, vs.Deployment.Labels)
		if !allowed {
			responseReview.Response.Allowed = allowed
			responseReview.Response.Result.Message = message
			break
		}
		allowed, message = checkContainers(&object.Spec.Template.Spec, vs.Deployment.Image)
		responseReview.Response.Allowed = allowed
		responseReview.Response.Result.Message = message
	case "ReplicaSet":
		object := &appsv1.ReplicaSet{}
		if _, _, err := universalDeserializer.Decode(req.Object.Raw, nil, object); err != nil {
			log.Printf("Couldn't decode object: %v", err)
		}
		allowed, message = checkLabels(object.Labels, vs.ReplicaSet.Labels)
		if !allowed {
			responseReview.Response.Allowed = allowed
			responseReview.Response.Result.Message = message
			break
		}
		allowed, message = checkContainers(&object.Spec.Template.Spec, vs.ReplicaSet.Image)
		responseReview.Response.Allowed = allowed
		responseReview.Response.Result.Message = message
	default:
		return responseReview
	}

	return responseReview
}

func checkLabels(objectLabels map[string]string, requiredLabels map[string]string) (bool, string) {
	for key := range requiredLabels {
		log.Printf("checking label %s with value %s", key, objectLabels[key])
		if _, ok := objectLabels[key]; !ok {
			return false, "Required label not present"
		}
	}
	return true, ""
}

func checkContainers(podSpec *corev1.PodSpec, imageRegex string) (bool, string) {
	for i := 0; i < len(podSpec.Containers); i++ {
		match, _ := regexp.MatchString(imageRegex, podSpec.Containers[i].Image)
		if !match {
			log.Println("Regex does not match image, admission rejected.")
			return false, "Regex does not match image, admission rejected."
		}
	}
	return true, ""
}
