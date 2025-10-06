package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/harundarat/be-socialtask/internal/middleware"
	"github.com/harundarat/be-socialtask/internal/store"
	"github.com/harundarat/be-socialtask/internal/utils"
)

type TaskHandler struct {
	taskStore store.TaskStore
	logger    *log.Logger
}

func NewTaskHandler(taskStore store.TaskStore, logger *log.Logger) *TaskHandler {
	return &TaskHandler{
		taskStore: taskStore,
		logger:    logger,
	}
}

func (th *TaskHandler) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	var task store.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		th.logger.Printf("ERROR: decodingCreateTask: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInvalidRequest, http.StatusBadRequest, nil, nil)
		return
	}

	users := middleware.GetUser(r)
	task.UserID = users.ID

	createdTask, err := th.taskStore.CreateTask(&task)
	if err != nil {
		th.logger.Printf("ERROR: createTask: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, nil)
		return
	}

	utils.WriteJSON(w, utils.StatusSuccess, utils.MessageTaskCreated, http.StatusCreated, utils.Envelope{"task": createdTask}, nil)

}

func (th *TaskHandler) HandleGetTaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDParam(r)
	if err != nil {
		th.logger.Printf("ERROR: readIdParam: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInvalidRequest, http.StatusBadRequest, nil, nil)
		return
	}

	task, err := th.taskStore.GetTaskByID(id)
	if err != nil {
		th.logger.Printf("ERROR: getTaskByID: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, nil)
		return
	}

	utils.WriteJSON(w, utils.StatusSuccess, utils.MessageTaskRetrieved, http.StatusOK, utils.Envelope{"task": task}, nil)
}

func (th *TaskHandler) HandleGetAllTask(w http.ResponseWriter, r *http.Request) {
	var (
		limit int64 = 10
	)

	page := utils.ReadPageParam(r)
	offset := (page - 1) * limit
	tasks, totalPages, err := th.taskStore.GetAllTask(limit, offset)
	if err != nil {
		th.logger.Printf("ERROR: getAllTask: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, nil)
		return
	}

	utils.WriteJSON(w, utils.StatusSuccess, utils.MessageTasksFetched, http.StatusOK, utils.Envelope{
		"tasks": tasks,
		"meta": map[string]int64{
			"page":  page,
			"limit": limit,
			"total": totalPages,
		},
	}, nil)
}

func (th *TaskHandler) HandleEditTask(w http.ResponseWriter, r *http.Request) {
	_ = middleware.GetUser(r)

	id, err := utils.ReadIDParam(r)
	if err != nil {
		th.logger.Printf("ERROR: readIdParam: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInvalidRequest, http.StatusBadRequest, nil, nil)
		return
	}

	var task store.Task
	// task.UserID = users.ID // bisa, if userID != task_userID

	task.ID = int(id)

	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		th.logger.Printf("ERROR: decodingCreateTask: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInvalidRequest, http.StatusBadRequest, nil, nil)
		return
	}

	err = th.taskStore.EditTask(&task)
	if err != nil {
		th.logger.Printf("ERROR: getTaskByID: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, nil)
		return
	}

	utils.WriteJSON(w, utils.StatusSuccess, utils.MessageTaskRetrieved, http.StatusOK, nil, nil)
}

func (th *TaskHandler) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	_ = middleware.GetUser(r)
	id, err := utils.ReadIDParam(r)
	if err != nil {
		th.logger.Printf("ERROR: readIdParam: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInvalidRequest, http.StatusBadRequest, nil, nil)
		return
	}

	err = th.taskStore.DeleteTask(id)
	if err != nil {
		th.logger.Printf("ERROR: deleteTask: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, nil)
		return
	}

	utils.WriteJSON(w, utils.StatusSuccess, utils.MessageTaskRetrieved, http.StatusOK, nil, nil)
}
