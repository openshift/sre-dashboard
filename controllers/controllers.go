package controllers

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"

	"github.com/gorilla/sessions"

	"fmt"
	"github.com/openshift/sre-dashboard/databases"
	"github.com/openshift/sre-dashboard/models"
	"net/http"
)

// GET /
func GetMain(c echo.Context) error {
	return c.Render(http.StatusOK, "main.html", nil)
}

// GET /login
func GetLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}

// GET /graph
func GetGraph(c echo.Context) error {
	// Default chart shows 1 week of takedowns.
	catCount := databases.QueryTakedowns(7)
	catMap := map[string]interface{}{"catCount": catCount}

	return c.Render(http.StatusOK, "graph_j_pie.html", catMap)
}

// GET /api/graph
func GetApiGraph(c echo.Context) error {
	/* Example of the returned data structure:
	{
	  "response": {
	    "bots": 491,
	    "disposable_email_address": 240,
	    "duplicate_accounts": 1369,
	    "other": 6
	  }
	}
	*/
	var content = models.Content
	var catCount map[string]int

	dateParam := c.QueryParam("dateparam")
	callBack := c.QueryParam("callback")

	switch {
	case dateParam == "day":
		catCount = databases.QueryTakedowns(1)
	case dateParam == "week":
		catCount = databases.QueryTakedowns(7)
	case dateParam == "month":
		catCount = databases.QueryTakedowns(30)
	case dateParam == "quarter":
		catCount = databases.QueryTakedowns(90)
	default:
		catCount = databases.QueryTakedowns(7)
	}

	content.Response = catCount

	return c.JSONP(http.StatusOK, callBack, &content)
}

// GET /trial
func GetTrial(c echo.Context) error {
	sess, _ := session.Get("session", c)
	logged_in_user := sess.Values["current_user"].(string)
	return c.String(http.StatusOK, logged_in_user)
}

// handle any error by attempting to render a custom page for it
func Custom404Handler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	errorPage := fmt.Sprintf("%d.html", code)
	if err := c.Render(code, errorPage, code); err != nil {
		c.Logger().Error(err)
	}
	c.Logger().Error(err)
}

func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			MainSession(c)
			sess, _ := session.Get("session", c)

			if sess.Values["authenticated"] == "true" {
				fmt.Println("in if block")
				fmt.Println(sess.Values)
				return next(c)
			}
			return c.Redirect(http.StatusTemporaryRedirect, "/login/google")
		}
	}
}

func MainSession(c echo.Context) { //error {
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 86400,
	}
	sess.Save(c.Request(), c.Response())
}
