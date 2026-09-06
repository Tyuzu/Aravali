package pay

import (
	"context"
	"errors"
	"net/http"
	"scav/utils"
)

// CreateWallet explicitly provisions an account for an onboarded user
func (p *PaymentService) CreateWallet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := utils.GetUserIDFromRequest(r)

	if _, err := p.getAccountByUserID(ctx, userID); err == nil {
		utils.RespondWithError(w, http.StatusConflict, "wallet already exists")
		return
	}

	if !p.userExists(ctx, userID) {
		utils.RespondWithError(w, http.StatusBadRequest, "user does not exist")
		return
	}

	newAcc, err := p.createWalletAccount(ctx, userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to create wallet")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]any{"success": true, "account_id": newAcc.ID})
}

// GetAccountStrict retrieves an account or returns an error if not found (No lazy generation!)
func (p *PaymentService) GetAccountStrict(ctx context.Context, userID string) (string, error) {
	acc, err := p.getAccountByUserID(ctx, userID)
	if err != nil {
		return "", errors.New("account_not_found")
	}
	return acc.ID, nil
}
