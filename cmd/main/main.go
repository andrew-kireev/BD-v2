package main

import (
	forumHandler "BD-v2/internal/app/forums/delivery/http"
	forumRep "BD-v2/internal/app/forums/repository"
	userHandler "BD-v2/internal/app/users/delivery/http"
	userRep "BD-v2/internal/app/users/repository"
	"BD-v2/internal/middlware"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jmoiron/sqlx"
	"net/http"
)

func main() {
	pool, err := pgxpool.Connect(context.Background(),
		"host=localhost port=5432 dbname=lera_bd sslmode=disable",
	)
	db, err := sqlx.Connect("postgres",
		"host=localhost port=5432 dbname=lera_bd sslmode=disable",
	)
	if err != nil {
		fmt.Println("Не смогли подключиться к бд")
	}

	forumRep := forumRep.NewForumRepository(db)
	userRep := userRep.NewUsersRepository(pool)
	router := mux.NewRouter()
	_ = userHandler.NewUsersHandler(router, userRep)
	_ = forumHandler.NewForumsHandler(router, forumRep)

	router.Use(middlware.ContentType)
	http.ListenAndServe(":5000", router)
}
