package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	oauth2api "github.com/cfthoo/todo-app/api/oauth"
	"github.com/cfthoo/todo-app/pkg/db/model"
	"github.com/gorilla/mux"
)

// DAO is the todoList data access object
type DAO interface {
	Create(ctx context.Context, task *model.Task) (*model.Task, error)
	FetchAll(ctx context.Context, userId string) ([]model.Task, error)
	FetchByID(ctx context.Context, id int) (*model.Task, error)
	Update(ctx context.Context, task *model.Task) (*model.Task, error)
	Delete(ctx context.Context, id int) error
}

// Handler provides all of the task handlers
type Handler struct {
	TodoListDAO DAO
}

type errorMessage struct {
	ID      string `json:"id,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

// Create will create task for todolist
func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task := &model.Task{}
		if err := json.NewDecoder(r.Body).Decode(task); err != nil {
			msg := &errorMessage{
				Error:   err.Error(),
				Message: "json decode error",
			}
			StdResponse(w, http.StatusBadRequest, msg)
			return
		}

		// set currently logged in userId to createdBy
		task.CreatedBy = oauth2api.UserId
		resp, err := h.TodoListDAO.Create(r.Context(), task)
		if err != nil {
			msg := &errorMessage{
				Error:   err.Error(),
				Message: "database error",
			}
			StdResponse(w, http.StatusInternalServerError, msg)
			return
		}
		StdResponse(w, http.StatusCreated, resp)
	}
}

// List will return all of the tasks
func (h *Handler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// set currently logged in userId to userID
		// here we only want to retrieve the task for a currently logged in user
		userID := oauth2api.UserId
		tasks, err := h.TodoListDAO.FetchAll(r.Context(), userID)
		if err != nil {
			msg := &errorMessage{
				Error:   err.Error(),
				Message: "json decode error",
			}
			StdResponse(w, http.StatusNotFound, msg)
		}

		StdResponse(w, http.StatusOK, tasks)

	}
}

// FetchByID will return task by id
func (h *Handler) FetchByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])

		task, err := h.TodoListDAO.FetchByID(r.Context(), id)

		if err != nil {
			msg := &errorMessage{
				Error:   err.Error(),
				Message: "Please check your datastore",
			}
			StdResponse(w, http.StatusInternalServerError, msg)
			return

		}

		StdResponse(w, http.StatusOK, task)

	}
}

// update will return the updated task
func (h *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task := &model.Task{}
		if err := json.NewDecoder(r.Body).Decode(task); err != nil {
			msg := &errorMessage{
				Error:   err.Error(),
				Message: "json decode error",
			}
			StdResponse(w, http.StatusBadRequest, msg)
			return
		}
		if task.ID <= 0 {
			msg := &errorMessage{
				Message: "Invalid Task Id",
			}
			StdResponse(w, http.StatusBadRequest, msg)
			return
		}
		if len(task.Name) == 0 {
			msg := &errorMessage{
				Message: "Task must have fields to update",
			}
			StdResponse(w, http.StatusBadRequest, msg)
			return
		}

		resp, err := h.TodoListDAO.Update(r.Context(), task)

		if err != nil {
			msg := &errorMessage{
				Error:   err.Error(),
				Message: "database error",
			}
			StdResponse(w, http.StatusInternalServerError, msg)
			return
		}
		StdResponse(w, http.StatusOK, resp)
	}
}

// delete will remove the task
func (h *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])

		err := h.TodoListDAO.Delete(r.Context(), id)

		if err != nil {
			msg := &errorMessage{
				Error:   err.Error(),
				Message: "database error",
			}
			StdResponse(w, http.StatusInternalServerError, msg)
			return

		}

		StdResponse(w, http.StatusNoContent, nil)
	}

}

// StdResponse will send a standard response with a json body
func StdResponse(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if body == nil {
		return
	}

	enc, err := json.Marshal(body)
	if err != nil {
		return
	}
	_, _ = w.Write(enc)
}
