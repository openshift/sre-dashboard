package auth

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/openshift/sre-dashboard/databases"
	"github.com/openshift/sre-dashboard/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
)

var (
	oauthStateString  = "random" // TODO randomize
	googleOauthConfig = &oauth2.Config{
		ClientID:     databases.OAuthID,
		ClientSecret: databases.OAuthKey,
		RedirectURL:  "https://sre-dashboard.openshift.com/oauth/callback", //"http://127.0.0.1:8080/oauth/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
)

// GET /oauth/callback
func HandleGoogleCallback(c echo.Context) error {
	state := c.QueryParam("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	code := c.QueryParam("code")
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("Code exchange failed with '%s'\n", err)
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Println("error getting response")
		fmt.Println(err)
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error reading response")
		fmt.Println(err)
	}

	var gUser models.GoogleUser
	err = json.Unmarshal(contents, &gUser)
	if err != nil {
		fmt.Println("Error Unmarshaling google user json: ", err)
	}

	if gUser.VerifiedEmail == true && gUser.HD == "redhat.com" {
		sess, _ := session.Get("session", c)
		sess.Values["authenticated"] = "true"
		sess.Values["google_logged_in"] = gUser.Email
		sess.Save(c.Request(), c.Response())

		return c.Render(http.StatusOK, "main.html", nil)
	}
	return c.String(200, string(contents)+`https://www.googleapis.com/oauth2/v1/tokeninfo?access_token=`+token.AccessToken)
}

// GET /login/google
func HandleGoogleLogin(c echo.Context) error {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}
