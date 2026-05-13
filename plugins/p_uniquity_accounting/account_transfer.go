package p_uniquity_accounting

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"strconv"

	"github.com/UniquityVentures/lamu/fields"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/views"
	"gorm.io/gorm"
)

// AccountTransferForm is the POST body for account transfer (not a DB model).
type AccountTransferForm struct {
	JournalID   uint
	ToAccountID uint
	Amount      fields.DecimalSix
}

// JournalAccountTransferForm is the POST body for journal-scoped account-to-account transfer (not a DB model).
type JournalAccountTransferForm struct {
	FromAccountID uint
	ToAccountID   uint
	Amount        fields.DecimalSix
}

// journalAccountTransferExcludeByQueryParam drops an account from picker lists using a GET query parameter (e.g. other FK from hx-include="closest form").
type journalAccountTransferExcludeByQueryParam struct {
	Param string
}

func (p journalAccountTransferExcludeByQueryParam) Patch(_ views.View, r *http.Request, query gorm.ChainInterface[Account]) gorm.ChainInterface[Account] {
	if p.Param == "" {
		return query
	}
	s := r.URL.Query().Get(p.Param)
	if s == "" {
		return query
	}
	u64, err := strconv.ParseUint(s, 10, 32)
	if err != nil || u64 == 0 {
		return query
	}
	return query.Where("id <> ?", uint(u64))
}

// accountTransferExcludeSourceQueryPatcher drops the source account row from the "to account" picker list (path param id).
type accountTransferExcludeSourceQueryPatcher struct{}

func (accountTransferExcludeSourceQueryPatcher) Patch(_ views.View, r *http.Request, query gorm.ChainInterface[Account]) gorm.ChainInterface[Account] {
	idStr := r.PathValue("id")
	if idStr == "" {
		return query
	}
	u64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return query
	}
	return query.Where("id <> ?", uint(u64))
}

func decimalSixNeg(a fields.DecimalSix) fields.DecimalSix {
	norm := a.NormalizeDecimals()
	if norm.R == nil {
		return fields.DecimalSix{R: big.NewRat(0, 1)}
	}
	neg := new(big.Rat).Neg(norm.R)
	return fields.DecimalSix{R: neg}.NormalizeDecimals()
}

// AccountToAccountTransfer creates one journal entry and two lines in a single transaction:
// the first line debits "from" (negative amount), the second credits "to" (positive amount).
func AccountToAccountTransfer(ctx context.Context, db *gorm.DB, journal Journal, amount fields.DecimalSix, from, to Account) (JournalEntry, JournalEntryItem, JournalEntryItem, error) {
	if from.ID == 0 || to.ID == 0 || from.ID == to.ID {
		return JournalEntry{}, JournalEntryItem{}, JournalEntryItem{}, fmt.Errorf("invalid source or destination account")
	}
	if journal.ID == 0 {
		return JournalEntry{}, JournalEntryItem{}, JournalEntryItem{}, fmt.Errorf("journal is required")
	}
	amount = amount.NormalizeDecimals()
	if amount.R == nil || amount.R.Sign() <= 0 {
		return JournalEntry{}, JournalEntryItem{}, JournalEntryItem{}, fmt.Errorf("amount must be positive")
	}
	negAmt := decimalSixNeg(amount)
	var je JournalEntry
	var fromLine, toLine JournalEntryItem
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		je = JournalEntry{JournalID: journal.ID}
		if err := gorm.G[JournalEntry](tx).Create(ctx, &je); err != nil {
			return err
		}
		fromLine = JournalEntryItem{
			JournalEntryID: je.ID,
			AccountID:      from.ID,
			Amount:         negAmt,
		}
		if err := gorm.G[JournalEntryItem](tx).Create(ctx, &fromLine); err != nil {
			return err
		}
		toLine = JournalEntryItem{
			JournalEntryID: je.ID,
			AccountID:      to.ID,
			Amount:         amount,
		}
		return gorm.G[JournalEntryItem](tx).Create(ctx, &toLine)
	})
	return je, fromLine, toLine, err
}

func accountTransferPostHandler(v *views.View) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values, fieldErrors, err := v.ParseForm(w, r)
		if err != nil {
			return
		}
		if fieldErrors == nil {
			fieldErrors = make(map[string]error)
		}

		fromAcct, ok := r.Context().Value("account").(Account)
		if !ok {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		journalID, journalOK := values["JournalID"].(uint)
		if !journalOK || journalID == 0 {
			fieldErrors["JournalID"] = fmt.Errorf("choose a journal")
		}
		toAccountID, toOK := values["ToAccountID"].(uint)
		if !toOK || toAccountID == 0 {
			fieldErrors["ToAccountID"] = fmt.Errorf("choose an account to transfer to")
		} else if toAccountID == fromAcct.ID {
			fieldErrors["ToAccountID"] = fmt.Errorf("cannot transfer to the same account")
		}
		amount, amtOK := values["Amount"].(fields.DecimalSix)
		if !amtOK {
			fieldErrors["Amount"] = fmt.Errorf("invalid amount")
		} else {
			norm := amount.NormalizeDecimals()
			if norm.R == nil || norm.R.Sign() <= 0 {
				fieldErrors["Amount"] = fmt.Errorf("amount must be positive")
			}
		}

		if len(fieldErrors) > 0 {
			ctx := views.ContextWithErrorsAndValues(r.Context(), values, fieldErrors)
			v.RenderPage(w, r.WithContext(ctx))
			return
		}

		db, dberr := getters.DBFromContext(r.Context())
		if dberr != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		_, _, _, txErr := AccountToAccountTransfer(ctx, db, Journal{Model: gorm.Model{ID: journalID}}, amount, fromAcct, Account{Model: gorm.Model{ID: toAccountID}})
		if txErr != nil {
			fieldErrors["_form"] = fmt.Errorf("%v", txErr)
			ctx := views.ContextWithErrorsAndValues(r.Context(), values, fieldErrors)
			v.RenderPage(w, r.WithContext(ctx))
			return
		}

		successURL, err := lamu.RoutePath("accounting.AccountDetailRoute", map[string]getters.Getter[any]{
			"id": getters.Any(getters.Static(fromAcct.ID)),
		})(ctx)
		if err != nil {
			fieldErrors["_form"] = err
			ctx := views.ContextWithErrorsAndValues(r.Context(), values, fieldErrors)
			v.RenderPage(w, r.WithContext(ctx))
			return
		}
		views.HtmxRedirect(w, r, successURL, http.StatusSeeOther)
	})
}

func journalAccountTransferPostHandler(v *views.View) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values, fieldErrors, err := v.ParseForm(w, r)
		if err != nil {
			return
		}
		if fieldErrors == nil {
			fieldErrors = make(map[string]error)
		}

		journalRow, ok := r.Context().Value("journal").(Journal)
		if !ok {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		fromID, fromOK := values["FromAccountID"].(uint)
		if !fromOK || fromID == 0 {
			fieldErrors["FromAccountID"] = fmt.Errorf("choose a from account")
		}
		toID, toOK := values["ToAccountID"].(uint)
		if !toOK || toID == 0 {
			fieldErrors["ToAccountID"] = fmt.Errorf("choose a to account")
		} else if fromOK && fromID != 0 && toID == fromID {
			fieldErrors["ToAccountID"] = fmt.Errorf("cannot transfer to the same account")
		}
		amount, amtOK := values["Amount"].(fields.DecimalSix)
		if !amtOK {
			fieldErrors["Amount"] = fmt.Errorf("invalid amount")
		} else {
			norm := amount.NormalizeDecimals()
			if norm.R == nil || norm.R.Sign() <= 0 {
				fieldErrors["Amount"] = fmt.Errorf("amount must be positive")
			}
		}

		if len(fieldErrors) > 0 {
			ctx := views.ContextWithErrorsAndValues(r.Context(), values, fieldErrors)
			v.RenderPage(w, r.WithContext(ctx))
			return
		}

		db, dberr := getters.DBFromContext(r.Context())
		if dberr != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		_, _, _, txErr := AccountToAccountTransfer(ctx, db, journalRow, amount,
			Account{Model: gorm.Model{ID: fromID}},
			Account{Model: gorm.Model{ID: toID}},
		)
		if txErr != nil {
			fieldErrors["_form"] = fmt.Errorf("%v", txErr)
			ctx := views.ContextWithErrorsAndValues(r.Context(), values, fieldErrors)
			v.RenderPage(w, r.WithContext(ctx))
			return
		}

		successURL, err := lamu.RoutePath("accounting.JournalDetailRoute", map[string]getters.Getter[any]{
			"id": getters.Any(getters.Static(journalRow.ID)),
		})(ctx)
		if err != nil {
			fieldErrors["_form"] = err
			ctx := views.ContextWithErrorsAndValues(r.Context(), values, fieldErrors)
			v.RenderPage(w, r.WithContext(ctx))
			return
		}
		views.HtmxRedirect(w, r, successURL, http.StatusSeeOther)
	})
}
