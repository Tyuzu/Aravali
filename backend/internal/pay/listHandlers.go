package pay

import (
	"net/http"
	"scav/utils"
	"strconv"
)

func (p *PaymentService) ListTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := utils.GetUserIDFromRequest(r)

	skip, _ := strconv.ParseInt(r.URL.Query().Get("skip"), 10, 64)
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)

	if limit <= 0 || limit > 100 {
		limit = 10
	}

	txns, err := p.listUserTransactions(ctx, userID, skip, limit)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, txns)
}

func (p *PaymentService) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := utils.GetUserIDFromRequest(r)

	acc, err := p.getAccountByUserID(ctx, userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "account not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"balance": acc.CachedBalance,
	})
}
