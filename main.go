package main

import (
	"github.com/openshift/sre-dashboard/routers"
)

func main() {
	e := routers.Routers
	e.Logger.Info(e.Start(":8443"))
}
