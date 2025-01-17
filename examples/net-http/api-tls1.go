// Go in action
// @jeffotoni
// 2019-01-01

// generating .key and .csr
// openssl req -nodes -newkey rsa:2048 -keyout server.key -out server.csr -subj "/C=BR/ST=Minas/L=Belo Horizonte/O=s3wf Ltd./OU=IT/CN=localhost"

// generating server .crt or .pem
// openssl x509 -req -sha256 -in server.csr -signkey server.key -out server.crt -days 365

package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
)

var (
	addr = ":443"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/v1/api/ping",
		func(w http.ResponseWriter, req *http.Request) {
			w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
			io.WriteString(w, "DevopsBH, Golang for Devops TLS MUX!\n")
		})

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	// show
	log.Printf("Server Run %s TLS / https://localhost%s", addr, addr)

	// conf listen TLS
	err := srv.ListenAndServeTLS("server.crt", "server.key")
	log.Fatal(err)
}

// curl --insecure -i -XGET https://localhost:8443/v1/api/ping
// curl -k -i -XGET https://localhost:8443/v1/api/ping
