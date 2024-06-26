package app

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/doktorupnos/crow/backend/internal/respond"
	"github.com/doktorupnos/crow/backend/internal/user"
	"github.com/go-chi/chi/v5"
)

var errDecodeRequestBody = errors.New("failed to decode request body")

func (app *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	body := RequestBody{}

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, errDecodeRequestBody)
		return
	}

	userID, err := app.userService.Create(body.Name, body.Password)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			err = user.ErrUserNameTaken
		}
		respond.Error(w, http.StatusBadRequest, err)
		return
	}
	respond.JWT(w, http.StatusCreated, app.Env.JWT.Secret, userID.String(), app.Env.JWT.Lifetime)
}

func (app *App) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := app.userService.GetAll()
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}
	respond.JSON(w, http.StatusOK, users)
}

func (app *App) GetUserByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		respond.Error(w, http.StatusBadRequest, errMissingURLParameter)
		return
	}

	u, err := app.userService.GetByName(name)
	if err != nil {
		respond.Error(w, http.StatusNotFound, err)
		return
	}
	respond.JSON(w, http.StatusOK, u)
}

func (app *App) DeleteUser(w http.ResponseWriter, r *http.Request, u user.User) {
	if err := app.userService.Delete(u); err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (app *App) UpdateUser(w http.ResponseWriter, r *http.Request, u user.User) {
	type RequestBody struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	body := RequestBody{}

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, errDecodeRequestBody)
		return
	}

	if err := app.userService.Update(u, body.Name, body.Password); err != nil {
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	u, err = app.userService.GetByID(u.ID)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
	}
	respond.JSON(w, http.StatusOK, u)
}

func (app *App) extractTargetUser(r *http.Request, u user.User) (user.User, error) {
	target := u
	if uq := r.URL.Query().Get("u"); uq != "" {
		var err error
		target, err = app.userService.GetByName(uq)
		if err != nil {
			return user.User{}, err
		}
	}
	return target, nil
}
