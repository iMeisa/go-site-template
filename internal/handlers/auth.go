package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/iMeisa/errortrace"
	"github.com/iMeisa/go-site-template/internal/config"
	"github.com/iMeisa/go-site-template/internal/models"
	"github.com/iMeisa/go-site-template/internal/tools"
	"io"
	"io/ioutil"
	"net/http"
)

func (m *Repository) Auth(w http.ResponseWriter, r *http.Request) {
	sessionState := tools.RandString(20)
	m.App.Session.Put(r.Context(), "state", sessionState)

	discordConfig := config.AuthDiscordConfig()

	http.Redirect(w, r, discordConfig.AuthCodeURL(sessionState), http.StatusSeeOther)
}

func (m *Repository) Callback(w http.ResponseWriter, r *http.Request) {
	resp := models.JsonResponse{
		Ok: false,
	}

	// Validate state
	sessionState := m.App.Session.Get(r.Context(), "state") // Might update to Session.Pop to destroy key after use
	if authState := r.URL.Query()["state"][0]; authState != sessionState {
		resp.Msg = fmt.Sprintf("states don't match: session: %v, auth: %v", sessionState, authState)
		writeResp(w, resp)
		return
	}

	// Token exchange
	code := r.URL.Query()["code"][0]
	discordConfig := config.AuthDiscordConfig()
	token, err := discordConfig.Exchange(context.Background(), code)
	if err != nil {
		trace := errortrace.NewTrace(err)
		trace.Read()
		resp.Msg = "code exchange failed"
		writeResp(w, resp)
		return
	}

	// Prep API call
	bearer := "Bearer " + token.AccessToken
	request, err := http.NewRequest("GET", "https://discord.com/api/v10/users/@me", nil)
	if err != nil {
		trace := errortrace.NewTrace(err)
		trace.Read()
		resp.Msg = "failed to create request"
		writeResp(w, resp)
	}
	request.Header.Add("Authorization", bearer)

	// Get discord user data
	client := &http.Client{}
	apiResp, err := client.Do(request)
	if err != nil {
		trace := errortrace.NewTrace(err)
		trace.Read()
		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			trace := errortrace.NewTrace(err)
			trace.Read()
			resp.Msg = "could not close body"
			writeResp(w, resp)
			return
		}
	}(apiResp.Body)

	// Read JSON bytes from api
	userBytes, err := ioutil.ReadAll(apiResp.Body)
	if err != nil {
		trace := errortrace.NewTrace(err)
		trace.Read()
		return
	}

	// Update user in db
	var user models.User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		trace := errortrace.NewTrace(err)
		trace.Read()
		return
	}
	trace := m.DB.LoginUpdateUser(&user)
	if trace.HasError() {
		trace.Read()
		return
	}

	// Marshal user to json
	userBytes, err = json.Marshal(user)
	if err != nil {
		trace = errortrace.NewTrace(err)
		trace.Read()
	}

	m.App.Session.Put(r.Context(), "user", string(userBytes))

	// Redirect to last visited page
	lastPage := m.App.Session.Get(r.Context(), "last_page")
	lastPagePath := fmt.Sprintf("%s", lastPage)
	// If last page does not start with /, redirect to /user/leagues
	if lastPagePath[0] != '/' {
		lastPagePath = "/user/leagues"
	}

	http.Redirect(w, r, lastPagePath, http.StatusSeeOther)
}

// Logout removes the user from the session
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {

	m.App.Session.Remove(r.Context(), "user")

	// Redirect to last visited page
	lastPage := m.App.Session.Get(r.Context(), "last_page")
	lastPagePath := fmt.Sprintf("%s", lastPage)
	// If last page does not start with /, redirect to /
	if lastPagePath[0] != '/' {
		lastPagePath = "/"
	}

	http.Redirect(w, r, lastPagePath, http.StatusSeeOther)
}
