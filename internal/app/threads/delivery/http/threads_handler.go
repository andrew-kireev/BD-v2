package http

import (
	"BD-v2/internal/app/forums"
	models2 "BD-v2/internal/app/posts/models"
	"BD-v2/internal/app/threads"
	"BD-v2/internal/app/threads/models"
	allModels "BD-v2/internal/models"
	tools "BD-v2/pkg"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgconn"
	"io/ioutil"
	"net/http"
	"strconv"
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

	r.HandleFunc("/api/thread/{slug}/create", theadHandler.CreatePostsByThreadSlug).Methods(http.MethodPost)
	r.HandleFunc("/api/thread/{slug}/vote", theadHandler.AddVoice).Methods(http.MethodPost)
	r.HandleFunc("/api/forum/{slug}/create", theadHandler.CrateThread).Methods(http.MethodPost)
	r.HandleFunc("/api/forum/{slug}/threads", theadHandler.GetTreads).Methods(http.MethodGet)

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
	defer r.Body.Close()
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

func (handler *ThreadsHandler) GetTreads(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 1
	}
	desc := tools.ConvertToBool(r.URL.Query().Get("desc"))
	since := r.URL.Query().Get("since")
	slug, _ := (mux.Vars(r))["slug"]

	forum, err := handler.forumRep.GetForumSlug(slug)
	if err != nil {
		w.WriteHeader(404)
		resp, _ := allModels.FailedResponse{
			Message: fmt.Sprintf("Не могут найти форум %s", forum.Slug),
		}.MarshalJSON()
		w.Write(resp)
		return
	}

	threads, err := handler.forumRep.GetTreads(limit, slug, since, desc)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(404)
		resp, _ := allModels.FailedResponse{
			Message: fmt.Sprintf("Не могут юзера найти %s", "fds"),
		}.MarshalJSON()
		w.Write(resp)
		return
	}
	respBody, _ := json.Marshal(threads)
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func (handler *ThreadsHandler) CreatePostsByThreadID(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("thread_id")
	posts := make([]*models2.Post, 0)
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &posts)
	if len(posts) == 0 {
		respBody, _ := json.Marshal(posts)
		w.WriteHeader(http.StatusCreated)
		w.Write(respBody)
		return
	}
	var err error
	ctx := context.Background()
	threadID, _ := strconv.Atoi(slug)
	thread := &models.Thread{}
	thread, err = handler.threadRep.FindThreadID(ctx, threadID)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(404)
		resp, _ := allModels.FailedResponse{
			Message: fmt.Sprintf("Не могут юзера найти %s", "fds"),
		}.MarshalJSON()
		w.Write(resp)
		return
	}

	for idx, post := range posts {
		post.Thread = thread.ID
		post.Forum = thread.Forum
		posts[idx], err = handler.threadRep.CreatePost(ctx, post)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(404)
			resp, _ := allModels.FailedResponse{
				Message: fmt.Sprintf("Не могут юзера найти %s", "fds"),
			}.MarshalJSON()
			w.Write(resp)
			return
		}
	}
	respBody, _ := json.Marshal(posts)
	w.WriteHeader(http.StatusCreated)
	w.Write(respBody)
}

func (handler *ThreadsHandler) CreatePostsByThreadSlug(w http.ResponseWriter, r *http.Request) {
	posts := make([]*models2.Post, 0)
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &posts)
	if len(posts) == 0 {
		respBody, _ := json.Marshal(posts)
		w.WriteHeader(http.StatusCreated)
		w.Write(respBody)
		return
	}
	vars := mux.Vars(r)
	slug, _ := vars["slug"]
	threadID, err := strconv.Atoi(slug)
	thread := &models.Thread{}
	ctx := context.Background()
	if err != nil {
		thread, err = handler.threadRep.FindThreadSlug(ctx, slug)
	} else {
		thread, err = handler.threadRep.FindThreadID(ctx, threadID)
	}

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(404)
		resp, _ := allModels.FailedResponse{
			Message: fmt.Sprintf("Не могут юзера найти %s", "fds"),
		}.MarshalJSON()
		w.Write(resp)
		return
	}

	for idx, post := range posts {
		post.Thread = thread.ID
		post.Forum = thread.Forum
		posts[idx], err = handler.threadRep.CreatePost(ctx, post)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(404)
			resp, _ := allModels.FailedResponse{
				Message: fmt.Sprintf("Не могут юзера найти %s", "fds"),
			}.MarshalJSON()
			w.Write(resp)
			return
		}
	}
	respBody, _ := json.Marshal(posts)
	w.WriteHeader(http.StatusCreated)
	w.Write(respBody)
}

func (handler *ThreadsHandler) AddVoice(w http.ResponseWriter, r *http.Request) {
	voice := &models.Vote{}
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &voice)

	vars := mux.Vars(r)
	slug, _ := vars["slug"]
	threadID, err := strconv.Atoi(slug)
	thread := &models.Thread{}
	ctx := context.Background()
	if err != nil {
		thread, err = handler.threadRep.FindThreadSlug(ctx, slug)
	} else {
		thread, err = handler.threadRep.FindThreadID(ctx, threadID)
	}
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(404)
		resp, _ := allModels.FailedResponse{
			Message: fmt.Sprintf("Не могут юзера найти %s", "fds"),
		}.MarshalJSON()
		w.Write(resp)
		return
	}
	voice.ThreadID = thread.ID
	err = handler.threadRep.AddVoice(ctx, voice)
	if err != nil {
		fmt.Print(err)
		oldVoice := voice.Voice
		voice, err := handler.threadRep.GetVoice(ctx, voice)
		if err != nil {
			fmt.Println(err)
		}
		if voice.Voice != oldVoice {
			err = handler.threadRep.UpdateVoice(ctx, voice)
		}
	}
	thread.Votes += voice.Voice
	respBody, _ := json.Marshal(thread)
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}
