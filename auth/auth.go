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
	"strings"
)

var (
	g_id, g_key = databases.OAuthID, databases.OAuthKey // getOauth("/secrets/google_auth_creds")
	//g_id, g_key       = getOauth("/home/remote/dedgar/ansible/google_auth_creds")
	oauthStateString  = "random"
	googleOauthConfig = &oauth2.Config{
		ClientID:     g_id,
		ClientSecret: g_key,
		RedirectURL:  "http://127.0.0.1:8080/oauth/callback", //"https://sre-dashboard.openshift.com/oauth/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
)

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

	fmt.Println("accessToken", token.AccessToken)

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
		fmt.Println(err)
	}

	if gUser.VerifiedEmail == true && strings.HasSuffix(gUser.Email, "@redhat.com") {
		sess, _ := session.Get("session", c)
		sess.Values["authenticated"] = "true"
		sess.Values["google_logged_in"] = gUser.Email
		sess.Save(c.Request(), c.Response())

		return c.Render(http.StatusOK, "main.html", nil)
	}
	return c.String(200, string(contents)+`https://www.googleapis.com/oauth2/v1/tokeninfo?access_token=`+token.AccessToken)
}

func HandleGoogleLogin(c echo.Context) error {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

/*
func getOauth(filepath string) (id, key string) {
	var appSecrets models.OauthCreds

	filePath := "/secrets/sre_dashboard_secrets.json"
	//filePath := "/home/remote/dedgar/ansible/sre_dashboard_secrets.json"
	fileBytes, err := ioutil.ReadFile(filePath)

	if err != nil {
		fmt.Println("Error loading secrets json: ", err)
	}

	err = json.Unmarshal(fileBytes, &appSecrets)
	if err != nil {
		fmt.Println("Error Unmarshaling secrets json: ", err)
	}

	id = appSecrets.GoogleAuthID
	key = appSecrets.GoogleAuthKey
  	filebytes, err := ioutil.ReadFile(filepath)
		if err != nil {
			fmt.Println("Error getting OAuth Credentials: ", err)
		}
		file_str := string(filebytes)

		id, key = strings.Split(file_str, "\n")[0], strings.Split(file_str, "\n")[1]

	return id, key
}*/
