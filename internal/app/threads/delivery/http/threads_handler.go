package http

import (
	"BD-v2/internal/app/forums"
	"BD-v2/internal/app/threads"
	"BD-v2/internal/app/threads/models"
	allModels "BD-v2/internal/models"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgconn"
	"io/ioutil"
	"net/http"
)

type ThreadsHandler struct {
	threadRep threads.Repository
	forumRep  forums.Repository
}

func NewThreadsHandler(r *mux.Router, threadRep threads.Repository,
	forumRep forums.Repository) *ThreadsHandler {
	theadHandler := &ThreadsHandler{
		threadRep: threadRep,
		forumRep:  forumRep,
	}

	r.HandleFunc("/api/forum/{slug}/create", theadHandler.CrateThread).Methods(http.MethodPost)

	return theadHandler
}

func (handler *ThreadsHandler) CrateThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	forum, ok := vars["slug"]
	if !ok {
		fmt.Println("не шмогли достать slug")
		w.WriteHeader(500)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		fmt.Println("Ошибка чтения body")
		return
	}
	thread := &models.Thread{}
	thread.UnmarshalJSON(body)
	forumModel, err := handler.forumRep.GetForumSlug(forum)
	if err != nil {
		w.WriteHeader(404)
		resp, _ := allModels.FailedResponse{
			Message: fmt.Sprintf("Не могут юзера найти %s", thread.Author),
		}.MarshalJSON()
		w.Write(resp)
		return
	}
	thread.Forum = forumModel.Slug
	threadSlug := thread.Slug
	thread, err = handler.threadRep.CreateThread(context.Background(), thread)
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok {
			if err.Code == "23505" {
				thr, _ := handler.threadRep.FindThreadSlug(context.Background(), threadSlug)
				respBody, _ := thr.MarshalJSON()
				w.WriteHeader(409)
				w.Write(respBody)
			} else {
				w.WriteHeader(404)
				resp, _ := allModels.FailedResponse{
					Message: fmt.Sprintf("Не могут юзера найти %s", thread.Author),
				}.MarshalJSON()
				w.Write(resp)
			}
		}
		return
	}
	respBody, _ := thread.MarshalJSON()
	w.WriteHeader(http.StatusCreated)
	w.Write(respBody)
}
