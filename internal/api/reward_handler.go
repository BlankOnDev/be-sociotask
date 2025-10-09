package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/harundarat/be-socialtask/internal/store"
	"github.com/harundarat/be-socialtask/internal/utils"
)

var validRewardTypes = map[store.JenisCategory]bool{
	store.CryptoUsdt1: true,
	store.CryptoUsdt2: true,
	store.CryptoUsdt3: true,
}

type RewardHandler struct {
	rewardStore store.TaskRewardStore
	logger      *log.Logger
}

func NewRewardHandler(rewardStore store.TaskRewardStore, logger *log.Logger) *RewardHandler {
	return &RewardHandler{
		rewardStore: rewardStore,
		logger:      logger,
	}
}

func (rh *RewardHandler) HandleCreateReward(w http.ResponseWriter, r *http.Request) {
	var reward store.RewardTask
	err := json.NewDecoder(r.Body).Decode(&reward)
	if err != nil {
		rh.logger.Printf("ERROR: decodingCreateReward: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInvalidRequest, http.StatusBadRequest, nil, nil)
		return
	}

	if !validRewardTypes[reward.RewardType] {
		utils.WriteJSON(w, utils.StatusError, "invalid reward type", http.StatusBadRequest, nil, nil)
		return
	}

	createdReward, err := rh.rewardStore.CreateReward(&reward)
	if err != nil {
		rh.logger.Printf("ERROR: createReward: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, nil)
		return
	}

	utils.WriteJSON(w, utils.StatusSuccess, utils.MessageRewardCreated, http.StatusCreated, utils.Envelope{"reward": createdReward}, nil)
}

func (rh *RewardHandler) HandleGetRewardByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDParam(r)
	if err != nil {
		rh.logger.Printf("ERROR: readIdParam: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInvalidRequest, http.StatusBadRequest, nil, nil)
		return
	}

	reward, err := rh.rewardStore.GetRewardByID(int(id))
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSON(w, utils.StatusError, "reward not found", http.StatusNotFound, nil, nil)
			return
		}
		rh.logger.Printf("ERROR: getRewardByID: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, nil)
		return
	}

	utils.WriteJSON(w, utils.StatusSuccess, utils.MessageRewardRetrieved, http.StatusOK, utils.Envelope{"reward": reward}, nil)
}

func (rh *RewardHandler) HandleGetAllReward(w http.ResponseWriter, r *http.Request) {
	rewards, err := rh.rewardStore.GetReward()
	if err != nil {
		rh.logger.Printf("ERROR: getAllReward: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, nil)
		return
	}

	utils.WriteJSON(w, utils.StatusSuccess, utils.MessageRewardsFetched, http.StatusOK, utils.Envelope{"rewards": rewards}, nil)
}

func (rh *RewardHandler) HandleEditReward(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDParam(r)
	if err != nil {
		rh.logger.Printf("ERROR: readIdParam: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInvalidRequest, http.StatusBadRequest, nil, nil)
		return
	}

	_, err = rh.rewardStore.GetRewardByID(int(id))
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSON(w, utils.StatusError, "reward not found", http.StatusNotFound, nil, nil)
			return
		}
		rh.logger.Printf("ERROR: getRewardByID for edit: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, nil)
		return
	}

	var input store.RewardTask
	input.ID = int(id)
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		rh.logger.Printf("ERROR: decodingEditReward: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInvalidRequest, http.StatusBadRequest, nil, nil)
		return
	}

	if input.RewardType != "" && !validRewardTypes[input.RewardType] {
		utils.WriteJSON(w, utils.StatusError, "invalid reward type", http.StatusBadRequest, nil, nil)
		return
	}

	err = rh.rewardStore.EditReward(&input)
	if err != nil {
		rh.logger.Printf("ERROR: editReward: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, nil)
		return
	}

	utils.WriteJSON(w, utils.StatusSuccess, utils.MessageRewardsUpdated, http.StatusOK, nil, nil)
}

func (rh *RewardHandler) HandleDeleteReward(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDParam(r)
	if err != nil {
		rh.logger.Printf("ERROR: readIdParam: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInvalidRequest, http.StatusBadRequest, nil, nil)
		return
	}

	err = rh.rewardStore.DeleteReward(int(id))
	if err != nil {
		rh.logger.Printf("ERROR: deleteReward: %v", err)
		utils.WriteJSON(w, utils.StatusError, utils.MessageInternalError, http.StatusInternalServerError, nil, nil)
		return
	}

	utils.WriteJSON(w, utils.StatusSuccess, utils.MessageRewardsDelete, http.StatusOK, nil, nil)
}
