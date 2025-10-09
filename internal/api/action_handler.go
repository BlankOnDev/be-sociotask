package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/harundarat/be-socialtask/internal/store"
	"github.com/harundarat/be-socialtask/internal/utils"
)

var validTaskTypes = map[store.TypeAction]bool{
	store.Type1: true,
	store.Type2: true,
	store.Type3: true,
}

type ActionHandler struct {
	actionStore store.TaskActionStore
	logger      *log.Logger
}

func NewActionHandler(actionStore store.TaskActionStore, logger *log.Logger) *ActionHandler {
	return &ActionHandler{
		actionStore: actionStore,
		logger:      logger,
	}
}

func (th *ActionHandler) HandleCreateAction(w http.ResponseWriter, r *http.Request) {
	var action store.ActionTask
	err := json.NewDecoder(r.Body).Decode(&action)
	if err != nil {
		th.logger.Printf("ERROR: decodingCreateAction: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInvalidRequest, http.StatusBadRequest, nil, nil)
		return
	}

	if !validTaskTypes[action.Type] {
		utils.WriteJSON(w, utils.StatusError, utils.MessageActionInvalidType, http.StatusBadRequest, nil, nil)
		return
	}

	createdAction, err := th.actionStore.CreateAction(&action)
	if err != nil {
		th.logger.Printf("ERROR: createAction: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, nil)
		return
	}

	utils.WriteJSON(w, utils.StatusSuccess, utils.MessageActionCreated, http.StatusCreated, utils.Envelope{"action": createdAction}, nil)
}

func (th *ActionHandler) HandleGetActionByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDParam(r)
	if err != nil {
		th.logger.Printf("ERROR: readIdParam: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInvalidRequest, http.StatusBadRequest, nil, nil)
		return
	}

	action, err := th.actionStore.GetActionByID(int(id))
	if err != nil {
		th.logger.Printf("ERROR: getActionByID: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, nil)
		return
	}

	utils.WriteJSON(w, utils.StatusSuccess, utils.MessageActionRetrieved, http.StatusOK, utils.Envelope{"action": action}, nil)
}

func (th *ActionHandler) HandleGetAllAction(w http.ResponseWriter, r *http.Request) {
	action, err := th.actionStore.GetAction()
	if err != nil {
		th.logger.Printf("ERROR: getAllTask: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, nil)
		return
	}

	utils.WriteJSON(w, utils.StatusSuccess, utils.MessageTasksFetched, http.StatusOK, utils.Envelope{"action": action}, nil)
}

func (th *ActionHandler) HandleEditAction(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDParam(r)
	if err != nil {
		th.logger.Printf("ERROR: readIdParam: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInvalidRequest, http.StatusBadRequest, nil, nil)
		return
	}

	_, err = th.actionStore.GetActionByID(int(id))
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSON(w, utils.StatusError, "action not found", http.StatusNotFound, nil, nil)
			return
		}
		th.logger.Printf("ERROR: getActionByID for edit: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, nil)
		return
	}

	var action store.ActionTask
	action.ID = int(id)
	err = json.NewDecoder(r.Body).Decode(&action)
	if err != nil {
		th.logger.Printf("ERROR: decodingEditAction: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInvalidRequest, http.StatusBadRequest, nil, nil)
		return
	}
	if action.Type != "" && !validTaskTypes[action.Type] {
		utils.WriteJSON(w, utils.StatusError, "invalid task type", http.StatusBadRequest, nil, nil)
		return
	}

	err = th.actionStore.EditAction(&action)
	if err != nil {
		th.logger.Printf("ERROR: getActionByID: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, nil)
		return
	}

	utils.WriteJSON(w, utils.StatusSuccess, utils.MessageActionsUpdated, http.StatusOK, nil, nil)
}

func (th *ActionHandler) HandleDeleteAction(w http.ResponseWriter, r *http.Request) {

	id, err := utils.ReadIDParam(r)
	if err != nil {
		th.logger.Printf("ERROR: readIdParam: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInvalidRequest, http.StatusBadRequest, nil, nil)
		return
	}

	err = th.actionStore.DeleteAction(int(id))
	if err != nil {
		th.logger.Printf("ERROR: deleteAction: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, nil)
		return
	}

	utils.WriteJSON(w, utils.StatusSuccess, utils.MessageActionsDelete, http.StatusOK, nil, nil)
}
