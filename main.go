package main

import (
	"github.com/openshift/sre-dashboard/routers"
)

func main() {
	e := routers.Routers
	e.Logger.Info(e.StartTLS(":8443", "/certs/sre-dashboard.openshift.com.crt", "/certs/sre-dashboard.openshift.com.key"))
}
