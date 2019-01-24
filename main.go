package main

import (
	"github.com/openshift/sre-dashboard/routers"
)

func main() {
	e := routers.Routers
	e.Logger.Info(e.StartTLS(":8443", "/cert/lego/certificates/dashboard.crt", "/cert/lego/certificates/dashboard.key"))
}
