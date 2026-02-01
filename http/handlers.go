package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"rest-api/todo"
	"time"

	"github.com/gorilla/mux"
)

type HTTPHandlers struct {
	todolist *todo.List
}

func NewHTTPhandlers(todolist *todo.List) *HTTPHandlers {
	return &HTTPHandlers{
		todolist: todolist,
	}
}

/*
pattern: /tasks
method: POST
info: JSON in HTTP request body
succeed:
  - status code: 201-Created
  - response body: JSON represent created task

faild:
  - status code: 400, 409, 500 ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	var taskDTO TaskDTO
	if err := json.NewDecoder(r.Body).Decode(&taskDTO); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}
	if err := taskDTO.ValidateForCreate(); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}
	todoTask := todo.NewTask(taskDTO.Title, taskDTO.Description)
	if err := h.todolist.AddTask(todoTask); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		if errors.Is(err, todo.ErrTaskAlreadyExists) {
			http.Error(w, errDTO.ToString(), http.StatusConflict)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}
	b, err := json.MarshalIndent(todoTask, "", "    ")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		fmt.Println("fail to write http response")
	}
}

/*
pattern: /tasks/{title}
method: GET
info: pattern
succeed:

	status code: 200-OK
	response body: JSON represented found task

failed:

	status code: 400, 404, 500 ...
	response body: JSON with error + time
*/
func (h *HTTPHandlers) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]
	task, err := h.todolist.GetTask(title)
	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		if errors.Is(todo.ErrTaskNotFound, err) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}

	b, err := json.MarshalIndent(task, "", "    ")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("fail to write http response")
	}
}

/*
pattern: /tasks
method: GET
info: -
succeed:

	status code: 200-OK
	response body: JSON represented found tasks

failed:

	status code: 400, 500 ...
	response body: JSON with error + time
*/
func (h *HTTPHandlers) HandleGetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.todolist.ListTask()
	b, err := json.MarshalIndent(tasks, "", "    ")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("fail to write http response")
	}
}

/*
pattern: /tasks?completed=true
method: GET
info: query params
succeed:

	status code: 200-OK
	response body: JSON represented found task

failed:

	status code: 400, 404, 500 ...
	response body: JSON with error + time
*/
func (h *HTTPHandlers) HandleGetAllUncomletedTasks(w http.ResponseWriter, r *http.Request) {
	uncompletedTasks := h.todolist.ListUnCompletedTasks()
	b, err := json.MarshalIndent(uncompletedTasks, "", "    ")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("fail to write http response")
		return
	}
}

/*
pattern: /tasks/{title}
method: PATCH
info: pattern + JSON in request body
succeed:

	status code: 200-OK
	response body: JSON represented found task

failed:

	status code: 400, 409, 500 ...
	response body: JSON with error + time
*/
func (h *HTTPHandlers) HandleCompleteTask(w http.ResponseWriter, r *http.Request) {
	var completeDTO CompletedTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&completeDTO); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}
	title := mux.Vars(r)["title"]

	var (
		changedTask todo.Task
		err         error
	)

	if completeDTO.Complete {
		changedTask, err = h.todolist.CompleteTask(title)
	} else {
		changedTask, err = h.todolist.UnCompleteTask(title)
	}
	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}
	b, err := json.MarshalIndent(changedTask, "", "    ")
	if err != nil {
		panic(err)
	}
	if _, err := w.Write(b); err != nil {
		fmt.Println("failde to write http response", err)
		return
	}
}

/*
pattern: /tasks/{title}
method: DELETE
info: pattern
succeed:

	status code: 204-No Content
	response body: -

failed:

	status code: 400, 404, 500 ...
	response body: JSON with error + time
*/
func (h *HTTPHandlers) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]
	if err := h.todolist.DeleteTask(title); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
