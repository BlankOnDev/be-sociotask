package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/harundarat/be-socialtask/internal/store"
	"github.com/harundarat/be-socialtask/internal/utils"
)

type createRewardRequest struct {
	UserID int64 `json:"user_id"`
	TaskID int64 `json:"task_id"`
}

type RewardsHandler struct {
	rewardsStore store.RewardsStore
	logger       *log.Logger
}

func NewRewardsHandler(rewardsStore store.RewardsStore, logger *log.Logger) *RewardsHandler {
	return &RewardsHandler{
		rewardsStore: rewardsStore,
		logger:       logger,
	}
}

func (h *RewardsHandler) validateCreateRewardRequest(req *createRewardRequest) error {
	if req.UserID == 0 {
		return errors.New("user_id is required and cannot be zero")
	}
	if req.TaskID == 0 {
		return errors.New("task_id is required and cannot be zero")
	}

	return nil
}

func (rh *RewardsHandler) HandleCreateReward(w http.ResponseWriter, r *http.Request) {
	var req createRewardRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		rh.logger.Printf("error decoding request body: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInvalidRequest, http.StatusBadRequest, nil, nil)
		return
	}

	err = rh.validateCreateRewardRequest(&req)
	if err != nil {
		utils.WriteJSON(w, utils.StatusError, utils.MessageValidationFailed, http.StatusBadRequest, nil, []string{err.Error()})
		return
	}

	reward := &store.Reward{
		UserID: req.UserID,
		TaskID: req.TaskID,
	}

	reward, err = rh.rewardsStore.Create(reward)
	if err != nil {
		rh.logger.Printf("ERROR: creating reward: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, nil)
		return
	}

	utils.WriteJSON(w, utils.StatusSuccess, "Reward created successfully", http.StatusCreated, utils.Envelope{"reward": reward}, nil)
}
