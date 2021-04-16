package main

import (
   "fmt"
   "log"
   "flag"
   "net/http"
   "crypto/tls"
	"os"
	"os/signal"
	"syscall"
	"context"
)

// --- HTTP Handlers

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

   // --- SSL Certificates

   flag.StringVar(&tlscert, "tlsCertFile", "/etc/certs/cert.pem", "File containing the x509 Certificate for HTTPS.")
   flag.StringVar(&tlskey, "tlsKeyFile", "/etc/certs/key.pem", "File containing the x509 private key to --tlsCertFile.")
   flag.Parse()

   certs, err := tls.LoadX509KeyPair(tlscert, tlskey)
   if err != nil {
      log.Fatal("Filed to load key pair: %v", err)
   }

   // --- TLS HTTP server configuration

   tlsConfig := &tls.Config{
      MinVersion:               tls.VersionTLS12,
      PreferServerCipherSuites: true,
      Certificates:             []tls.Certificate{certs},
      CurvePreferences:         []tls.CurveID{
                                   tls.CurveP521,
                                   tls.CurveP384,
                                   tls.CurveP256,
                                },
      CipherSuites:             []uint16{
                                   tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
                                   tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
                                   tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
                                   tls.TLS_RSA_WITH_AES_256_CBC_SHA,
                                },
   }

   // --- HTTP Mux definition

   validate := lbpValidate{}

   mux := http.NewServeMux()
   mux.HandleFunc("/validate", validate.serve)
   mux.HandleFunc("/headers",  headers)
   mux.HandleFunc("/healthz",  healthz)
   mux.HandleFunc("/readyz",   readyz)
   mux.HandleFunc("/annotate", annotate)

   // --- HTTP Server
   server := &http.Server{
                Addr:         ":8080",
                Handler:      mux,
                TLSConfig:    tlsConfig,
                TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
             }

   go func() {
      log.Fatal(server.ListenAndServeTLS("",""))
      return
   }()

	// listening shutdown singal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("Got shutdown signal, shutting down webhook server gracefully...")
	server.Shutdown(context.Background())

}
