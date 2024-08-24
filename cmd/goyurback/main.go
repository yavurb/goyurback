package main

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	go func() {
		host := fmt.Sprintf("0.0.0.0:%s", appCtx.Settings.Port)
		if err := app.Start(host); err != nil && err != http.ErrServerClosed {
			app.Logger.Fatal("shutting down the server", err)
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		app.Logger.Fatal(err)
	}
}
