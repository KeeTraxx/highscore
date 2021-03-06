package main

import (
	"fmt"
	"net/http"

	"encoding/gob"

	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/facebook"
	"github.com/dghubble/gologin/google"
	scrypt "github.com/elithrar/simple-scrypt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"golang.org/x/oauth2"
	facebookOAuth2 "golang.org/x/oauth2/facebook"
	googleOAuth2 "golang.org/x/oauth2/google"
)

func setupAuth(e *echo.Echo) {
	gob.Register(&Profile{})
	setupGoogle(e)
	setupFacebook(e)

	e.POST("/logout", func(c echo.Context) error {
		//sessionStore.Destroy(c.Request(), sessionName)

		sess, _ := sessionStore.Get(c.Request(), sessionName)
		sess.Options.MaxAge = -1

		sess.Save(c.Request(), c.Response())

		return c.Redirect(http.StatusFound, "/")
	})

	e.GET("/api/profile", getProfile)
	e.POST("/api/register", register)
	e.POST("/local/login", localLogin)

}

func localLogin(c echo.Context) error {
	var loginData map[string]string
	c.Bind(&loginData)

	name := loginData["name"]
	password := loginData["password"]

	var user User

	DB.LogMode(true)

	if DB.
		Where("name = ? AND type = ?", name, "local").
		Or("email = ? AND type = ?", name, "local").
		First(&user).
		RecordNotFound() {
		return c.NoContent(http.StatusForbidden)
	}

	Info.Printf("%+v", user)

	err := scrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password))

	if err != nil {
		Error.Println(err)
		return c.NoContent(http.StatusForbidden)
	}

	sess, err := sessionStore.New(c.Request(), sessionName)

	if err != nil {
		Error.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	sess.Values["profile"] = &Profile{
		Email: user.Email,
		Name:  user.Name,
		ID:    user.ID,
	}

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		Error.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func getProfile(c echo.Context) error {
	sess, _ := sessionStore.Get(c.Request(), sessionName)
	//sess, _ := session.Get(sessionName, c)
	Info.Printf("%+v", sess)

	_, exists := sess.Values["profile"]

	if !exists {
		return c.JSON(http.StatusOK, nil)
	}

	profile := sess.Values["profile"].(*Profile)

	var user User

	DB.First(&user, profile.ID)

	return c.JSON(http.StatusOK, &user)
}

func register(c echo.Context) error {
	var registration map[string]string
	c.Bind(&registration)

	var user User

	if !DB.First(&user, &User{
		Email: registration["email"],
		Type:  "local",
	}).RecordNotFound() {
		Error.Printf("Email already. %+v", user.Email)
		return c.String(http.StatusNotAcceptable, "Email already exists...")
	}

	if !DB.First(&user, &User{
		Name: registration["name"],
		Type: "local",
	}).RecordNotFound() {
		Error.Println("User already.")
		return c.String(http.StatusNotAcceptable, "Username already exists...")
	}

	pw, _ := scrypt.GenerateFromPassword([]byte(registration["password"]), scrypt.DefaultParams)
	Info.Printf("%+v %v", registration, string(pw))
	pwstr := string(pw)

	DB.Save(&User{
		Name:       registration["name"],
		Password:   &pwstr,
		ProviderID: "local",
		Type:       "local",
		Email:      registration["email"],
	})

	return c.NoContent(http.StatusOK)
}

func setupGoogle(e *echo.Echo) {
	oauth2Config := &oauth2.Config{
		ClientID:     "736599901494-vet694v3bdbum6n2bbfibdanuam6b02v.apps.googleusercontent.com",
		ClientSecret: "OVpLOVrGlQNW28S0UuApuqO4",
		RedirectURL:  getenv("BASE_URL", "http://localhost:1323") + "/google/callback",
		Endpoint:     googleOAuth2.Endpoint,
		Scopes:       []string{"profile", "email"},
	}
	stateConfig := gologin.DebugOnlyCookieConfig

	e.Any("/google/login", echo.WrapHandler(google.StateHandler(stateConfig, google.LoginHandler(oauth2Config, nil))))
	e.Any("/google/callback", echo.WrapHandler(google.StateHandler(stateConfig, google.CallbackHandler(oauth2Config, issueSessionGoogle(), nil))))
}

func setupFacebook(e *echo.Echo) {
	oauth2Config := &oauth2.Config{
		ClientID:     "796487397192190",
		ClientSecret: "7970506ceaef3e7c1c3ad25f3533ea7e",
		RedirectURL:  getenv("BASE_URL", "http://localhost:1323") + "/facebook/callback",
		Endpoint:     facebookOAuth2.Endpoint,
		Scopes:       []string{"email"},
	}
	stateConfig := gologin.DebugOnlyCookieConfig

	e.Any("/facebook/login", echo.WrapHandler(facebook.StateHandler(stateConfig, facebook.LoginHandler(oauth2Config, nil))))
	e.Any("/facebook/callback", echo.WrapHandler(facebook.StateHandler(stateConfig, facebook.CallbackHandler(oauth2Config, issueSessionFacebook(), nil))))
}

const (
	sessionName   = "highscore-session"
	sessionSecret = "highscore-secret-cookie-salt"
)

// sessionStore encodes and decodes session data stored in signed cookies
var sessionStore = sessions.NewCookieStore([]byte(sessionSecret), nil)

// Profile represents basic info about an user
type Profile struct {
	ID         uint   `json:"id"`
	ProviderID string `json:"provider_id"`
	Locale     string `json:"locale"`
	Link       string `json:"link"`
	Email      string `json:"email"`
	FamilyName string `json:"family-name"`
	GivenName  string `json:"given-name"`
	Name       string `json:"name"`
	Picture    string `json:"picture"`
	Gender     string `json:"gender"`
}

func issueSessionGoogle() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		googleUser, err := google.UserFromContext(ctx)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		Info.Printf("Got userino: %+v", googleUser)

		// 2. Implement a success handler to issue some form of session
		session, err := sessionStore.New(req, sessionName)

		profile := new(Profile)

		profile.Email = googleUser.Email
		profile.FamilyName = googleUser.FamilyName
		profile.Gender = googleUser.Email
		profile.GivenName = googleUser.GivenName
		profile.ProviderID = googleUser.Id
		profile.Link = googleUser.Link
		profile.Locale = googleUser.Locale
		profile.Name = googleUser.Name
		profile.Picture = googleUser.Picture

		session.Values["profile"] = profile

		var user User

		DB.FirstOrInit(&user, &User{
			ProviderID: profile.ProviderID,
			Type:       "google",
		})

		if DB.NewRecord(&user) {
			user.Email = profile.Email
			user.Picture = &profile.Picture
			user.Name = profile.Name
			DB.Save(&user)
		}

		profile.ID = user.ID

		err = session.Save(req, w)

		if err != nil {
			Error.Println(err)
		}

		http.Redirect(w, req, "/", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}

func authRequired() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, err := sessionStore.Get(c.Request(), sessionName)

			if err != nil {
				Error.Println(err)
				return c.NoContent(http.StatusInternalServerError)
			}

			if _, authenticated := sess.Values["profile"]; authenticated {
				return next(c)
			}

			return echo.ErrUnauthorized

		}
	}
}

func issueSessionFacebook() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		facebookUser, err := facebook.UserFromContext(ctx)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		Info.Printf("Got userino: %+v", facebookUser)

		// 2. Implement a success handler to issue some form of session
		session, err := sessionStore.New(req, sessionName)

		profile := new(Profile)

		profile.Email = facebookUser.Email
		profile.Name = facebookUser.Name
		profile.ProviderID = facebookUser.ID
		profile.Picture = fmt.Sprintf("https://graph.facebook.com/%v/picture?type=large", profile.ProviderID)

		session.Values["profile"] = profile

		var user User

		DB.FirstOrInit(&user, &User{
			ProviderID: profile.ProviderID,
			Type:       "facebook",
		})

		if DB.NewRecord(&user) {
			user.Email = profile.Email
			user.Picture = &profile.Picture
			user.Name = profile.Name
			DB.Save(&user)
		}

		profile.ID = user.ID

		err = session.Save(req, w)

		if err != nil {
			Error.Println(err)
		}

		http.Redirect(w, req, "/", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}
