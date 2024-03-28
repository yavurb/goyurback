package main

import (
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/yavurb/goyurback/internal/app"
)

var (
	//go:embed .version
	_version string

	Version string = strings.TrimSpace(_version)
)

func main() {
	appCtx := app.NewAppContext()
	defer appCtx.Connpool.Close()

	app := appCtx.NewRouter()

	fmt.Printf(`
 ██████╗  ██████╗ ██╗   ██╗██╗   ██╗██████╗ ██████╗  █████╗  ██████╗██╗  ██╗
██╔════╝ ██╔═══██╗╚██╗ ██╔╝██║   ██║██╔══██╗██╔══██╗██╔══██╗██╔════╝██║ ██╔╝
██║  ███╗██║   ██║ ╚████╔╝ ██║   ██║██████╔╝██████╔╝███████║██║     █████╔╝
██║   ██║██║   ██║  ╚██╔╝  ██║   ██║██╔══██╗██╔══██╗██╔══██║██║     ██╔═██╗
╚██████╔╝╚██████╔╝   ██║   ╚██████╔╝██║  ██║██████╔╝██║  ██║╚██████╗██║  ██╗
 ╚═════╝  ╚═════╝    ╚═╝    ╚═════╝ ╚═╝  ╚═╝╚═════╝ ╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝.%s
	`, Version)
	fmt.Println()

	caCertFile, err := os.ReadFile("certs/cert.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertFile)

	mtlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
		MinVersion: tls.VersionTLS13,
	}

	host := fmt.Sprintf("0.0.0.0:%s", appCtx.Settings.Port)
	server := http.Server{
		Addr:      host,
		Handler:   app,
		TLSConfig: mtlsConfig,
	}

	app.Logger.Fatal(server.ListenAndServeTLS("certs/cert.pem", "certs/key.pem"))
}
