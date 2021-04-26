package http

import (
	"BD-v2/internal/app/users"
	"BD-v2/internal/app/users/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type UsersHandler struct {
	UsersRep users.Repository
}

func NewUsersHandler(r *mux.Router, rep users.Repository) *UsersHandler {
	usersHandler := &UsersHandler{
		UsersRep: rep,
	}

	r.HandleFunc("/api/user/{nickname}/create", usersHandler.CreateUser).Methods(http.MethodPost)

	return usersHandler
}

func (userHandler *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		fmt.Println("Ошибка чтения body")
		return
	}
	vars := mux.Vars(r)
	nickname, ok := vars["nickname"]
	if !ok {
		fmt.Println("не шмогли достать nickname")
		w.WriteHeader(500)
		return
	}
	user := &models.User{}
	user.Nickname = nickname
	err = user.UnmarshalJSON(body)
	if err != nil {
		w.WriteHeader(500)
		fmt.Println("Ошибка распаковки json")
		return
	}

	err = userHandler.UsersRep.CreateUser(context.Background(), user)
	if err != nil {
		existedUsers, err := userHandler.UsersRep.CheckIfUserExist(context.Background(), user)
		if err != nil {
			w.WriteHeader(500)
			fmt.Println("какая-то проблемка")
			return
		}
		respBody, _ := json.Marshal(existedUsers)
		w.WriteHeader(http.StatusConflict)
		w.Write(respBody)
		return
	}
	respBody, _ := user.MarshalJSON()
	w.WriteHeader(http.StatusCreated)
	w.Write(respBody)
}
