package main

import (
   "fmt"
   "log"
   "flag"
   "net/http"
   "crypto/tls"
//   "os"
//   "os/signal"
//   "syscall"
//   "context"

)

func validate(w http.ResponseWriter, req *http.Request) {

   fmt.Fprintf(w, "validate\n")
}

func healthz(w http.ResponseWriter, req *http.Request) {
   fmt.Fprintf(w, "ok\n")
}

func readyz(w http.ResponseWriter, req *http.Request) {
   fmt.Fprintf(w, "ok\n")
}

func annotate(w http.ResponseWriter, req *http.Request) {

   fmt.Fprintf(w, "vmware.com/nsx=true\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
   for name, headers := range req.Header {
      for _, h := range headers {
         fmt.Fprintf(w, "%v: %v\n", name, h)
      }
   }
}


func main() {

   var tlscert, tlskey string

   flag.StringVar(&tlscert, "tlsCertFile", "/etc/certs/cert.pem", "File containing the x509 Certificate for HTTPS.")
   flag.StringVar(&tlskey, "tlsKeyFile", "/etc/certs/key.pem", "File containing the x509 private key to --tlsCertFile.")
   flag.Parse()

   certs, err := tls.LoadX509KeyPair(tlscert, tlskey)
   if err != nil {
      log.Fatal("Filed to load key pair: %v", err)
   }

   server := &http.Server{
     Addr:      ":8080",
     TLSConfig: &tls.Config{Certificates: []tls.Certificate{certs}},
   }

   validate := lbpValidate{}
   mux := http.NewServeMux()
   mux.HandleFunc("/validate", validate.serve)
   mux.HandleFunc("/headers",  headers)
   mux.HandleFunc("/healthz",  healthz)
   mux.HandleFunc("/readyz",   readyz)
   mux.HandleFunc("/annotate", annotate)
   server.Handler = mux

   log.Fatal(server.ListenAndServeTLS("",""))

}
