package main

import (
	_ "embed"
	"fmt"
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

	host := fmt.Sprintf("0.0.0.0:%s", appCtx.Settings.Port)
	app.Logger.Fatal(app.Start(host))
}
