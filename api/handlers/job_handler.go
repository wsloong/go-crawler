package handlers

import (
	"net/http"
	"sync"
)

var JobHandler = newJobHandler()

type jobHandler struct {
	once *sync.Once
}

func newJobHandler() *jobHandler {
	return &jobHandler{
		once: new(sync.Once),
	}
}

func (j *jobHandler) RegisterRoutes() {
	j.once.Do(func() {
		http.HandleFunc("/", j.Index)
	})
}

func (j *jobHandler) Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("success"))
	return
}
