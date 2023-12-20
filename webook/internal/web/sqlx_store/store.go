package sqlx_store

import (
	gin_session "github.com/gin-contrib/sessions"
	"github.com/gorilla/sessions"
	"net/http"
)

type Store struct {
}

func (s *Store) Get(r *http.Request, name string) (*sessions.Session, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Store) New(r *http.Request, name string) (*sessions.Session, error) {
	//TODO implement me
	panic("implement me")
}

func (st *Store) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	//TODO implement me
	panic("implement me")
}

func (s *Store) Options(options gin_session.Options) {
	//TODO implement me
	panic("implement me")
}
