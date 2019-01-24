package routers

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"

	"github.com/gorilla/sessions"

	"github.com/openshift/sre-dashboard/auth"
	"github.com/openshift/sre-dashboard/controllers"
	"github.com/openshift/sre-dashboard/databases"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	Routers *echo.Echo
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func init() {
	t := &Template{
		templates: func() *template.Template {
			tmpl := template.New("")
			if err := filepath.Walk("/usr/local/bin/tmpl", func(path string, info os.FileInfo, err error) error {
				if strings.HasSuffix(path, ".html") {
					_, err = tmpl.ParseFiles(path)
					if err != nil {
						log.Println(err)
					}
				}
				return err
			}); err != nil {
				panic(err)
			}
			return tmpl
		}(),
	}

	Routers = echo.New()
	Routers.Static("/", "/usr/local/bin/static")
	Routers.Renderer = t

	Routers.Use(middleware.Logger())
	Routers.Use(middleware.Recover())

	Routers.Use(session.Middleware(sessions.NewCookieStore([]byte(databases.CookieSecret))))

	// AuthMiddleware Requires users be logged in with an @redhat.com email
	Routers.GET("/", controllers.GetMain, controllers.AuthMiddleware())
	Routers.POST("/", controllers.GetMain)
	Routers.GET("/takedowns", controllers.GetGraph, controllers.AuthMiddleware())
	Routers.GET("/api/takedowns", controllers.GetApiGraph)
	Routers.GET("/login/google", auth.HandleGoogleLogin)
	Routers.GET("/oauth/callback", auth.HandleGoogleCallback)
}
