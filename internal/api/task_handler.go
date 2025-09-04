package api

import (
	"encoding/json"
	"log"
	"net/http"

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
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	createdTask, err := th.taskStore.CreateTask(&task)
	if err != nil {
		th.logger.Printf("ERROR: createTask: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create task"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"task": createdTask})

}
