package main

import (
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/certeu/multigeo/multigeo"
	"github.com/certeu/multigeo/router"
)

var buildInfo string

var (
	flagVersion  = flag.Bool("v", false, "Display version and exit")
	flagAddr     = flag.String("addr", ":8001", "Bind address")
	flagCertFile = flag.String("cert", "cert.pem", "Certificate file")
	flagKeyFile  = flag.String("key", "key.pem", "Certificate private key file")
	flagMMDb     = flag.String("mm", "", "MaxMind database file")
	flagIP2LDb   = flag.String("ip2l", "", "IP2Location database file")
)

type GeoResponse struct {
	IP      net.IP                 `json:"ip"`
	GeoData []multigeo.GeoResponse `json:"geodata"`
}

// startHTTPServer starts a HTTP server which will redirects all HTTP requests
// to HTTPS
func startHTTPServer() {
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Connection", "close")
			url := "https://" + req.Host + req.URL.String()
			http.Redirect(w, req, url, http.StatusMovedPermanently)
		}),
	}
	log.Println(srv.ListenAndServe())
}

// getTLSConfig returns the TLS configuration
// props Filippo Valsorda
func getTLSConfig(cert *tls.Certificate) *tls.Config {
	tlsCfg := &tls.Config{
		Certificates:             []tls.Certificate{*cert},
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519, // Go 1.8 only
		},
		MinVersion: tls.VersionTLS11,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			// Remove when upgrading to TLS1.2 only
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	return tlsCfg
}

func geoHandler(w http.ResponseWriter, r *http.Request) {
	i := strings.Split(r.URL.Path, "/")
	ip := net.ParseIP(i[len(i)-1])
	if ip == nil {
		router.NotFound(w, r)
		return
	}

	mm := multigeo.MaxMind{}
	ip2l := multigeo.IP2Location{}

	rr := GeoResponse{IP: ip}

	gr, eer := mm.ToGeo(ip)
	if eer != nil {
		log.Print(eer)
	}
	rr.GeoData = append(rr.GeoData, gr)

	ir, ier := ip2l.ToGeo(ip)
	if ier != nil {
		log.Print(ier)
	}
	rr.GeoData = append(rr.GeoData, ir)

	if i[1] == "xml" {
		w.Header().Set("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(rr)
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rr)
	}
}

func main() {
	flag.Parse()
	if *flagVersion {
		fmt.Printf("Version %s\n", buildInfo)
		return
	}

	err := multigeo.NewMaxMind(*flagMMDb)
	if err != nil {
		log.Fatal(err)
	}

	err = multigeo.NewIP2Location(*flagIP2LDb)
	if err != nil {
		log.Fatal(err)
	}

	go startHTTPServer()

	m := router.NewRouter()
	m.HandleFunc("/json/[0-9A-Fa-f]+", geoHandler)
	m.HandleFunc("/xml/[0-9A-Fa-f]+", geoHandler)

	cert, err := tls.LoadX509KeyPair(*flagCertFile, *flagKeyFile)
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:         *flagAddr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		TLSConfig:    getTLSConfig(&cert),
		Handler:      m,
	}
	log.Fatal(srv.ListenAndServeTLS("", ""))

}
