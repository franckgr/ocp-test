package main

import (
   "fmt"
   "log"
   "io/ioutil"
   "net/http"
   "encoding/json"

   "github.com/golang/glog"
   "k8s.io/api/admission/v1beta1"
   "k8s.io/api/core/v1"
   metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// lbpValidate listen to admission requests and serve responses
type lbpValidate struct {
}

func (validate *lbpValidate) serve(w http.ResponseWriter, r *http.Request) {
   var body []byte
   if r.Body != nil {
      if data, err := ioutil.ReadAll(r.Body); err == nil {
      body = data
      }
   }

   if len(body) == 0 {
      log.Fatal("empty body")
      http.Error(w, "empty body", http.StatusBadRequest)
      return
   }
   log.Println("Received request")

   if r.URL.Path != "/validate" {
      log.Fatal("no validate")
      http.Error(w, "no validate", http.StatusBadRequest)
      return
   }

   arRequest := v1beta1.AdmissionReview{}
   if err := json.Unmarshal(body, &arRequest); err != nil {
      log.Fatal("incorrect body")
      http.Error(w, "incorrect body", http.StatusBadRequest)
   }

   raw := arRequest.Request.Object.Raw
   pod := v1.Pod{}
   if err := json.Unmarshal(raw, &pod); err != nil {
      log.Fatal("error deserializing pod")
      return
   }
   if pod.Name == "smooth-app" {
      return
   }

   arResponse := v1beta1.AdmissionReview{
      Response: &v1beta1.AdmissionResponse{
         Allowed: false,
         Result: &metav1.Status{
            Message: "Keep calm and not add more crap in the cluster!",
         },
      },
   }
   resp, err := json.Marshal(arResponse)
   if err != nil {
      log.Fatal("Can't encode response: %v", err)
      http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
   }
   glog.Infof("Ready to write reponse ...")
   if _, err := w.Write(resp); err != nil {
      log.Fatal("Can't write response: %v", err)
      http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
   }
}
