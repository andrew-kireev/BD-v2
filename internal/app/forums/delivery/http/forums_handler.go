package http

import (
	"BD-v2/internal/app/forums"
	"BD-v2/internal/app/forums/models"
	allModels "BD-v2/internal/models"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"io/ioutil"
	"net/http"
)

type ForumsHandler struct {
	ForumsRep forums.Repository
}

func NewForumsHandler(r *mux.Router, rep forums.Repository) *ForumsHandler {
	forumsHandler := &ForumsHandler{
		ForumsRep: rep,
	}

	r.HandleFunc("/api/forum/create", forumsHandler.CreateForum).Methods(http.MethodPost)
	r.HandleFunc("/api/forum/{slug}/details", forumsHandler.GetForum).Methods(http.MethodGet)

	return forumsHandler
}

func (handler *ForumsHandler) CreateForum(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		fmt.Println("Ошибка чтения body")
		return
	}
	forum := &models.Forum{}
	err = forum.UnmarshalJSON(body)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	err = handler.ForumsRep.CreateForum(forum)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == "23505" {
				f, _ := handler.ForumsRep.GetForumSlug(forum.Slug)
				respBody, _ := f.MarshalJSON()
				w.WriteHeader(409)
				w.Write(respBody)
			} else {
				w.WriteHeader(404)
				resp, _ := allModels.FailedResponse{
					Message: fmt.Sprintf("Не могут юзера найти %s", forum.User),
				}.MarshalJSON()
				w.Write(resp)
			}
		}
		return
	}
	respBody, _ := forum.MarshalJSON()
	w.WriteHeader(201)
	w.Write(respBody)
}

func (handler *ForumsHandler) GetForum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug, ok := vars["slug"]
	if !ok {
		fmt.Println("не шмогли достать slug")
		w.WriteHeader(500)
		return
	}
	forum, err := handler.ForumsRep.GetForumSlug(slug)
	if err != nil {
		w.WriteHeader(404)
		resp, _ := allModels.FailedResponse{
			Message: fmt.Sprintf("Не могут найти форум %s", slug),
		}.MarshalJSON()
		w.Write(resp)
		return
	}
	respBody, _ := forum.MarshalJSON()
	w.WriteHeader(200)
	w.Write(respBody)
}
