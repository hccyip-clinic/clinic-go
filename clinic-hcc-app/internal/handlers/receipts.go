package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"clinic-hcc-app/internal/models"

	"github.com/go-chi/chi/v5"
)

func (r *Router) ReceiptList(w http.ResponseWriter, req *http.Request) {
	var (
		receipts []models.Receipt
		err      error
	)
	if req.URL.Query().Get("date") == "today" {
		receipts, err = r.receipts.ListToday(req.Context())
	} else {
		receipts, err = r.receipts.List(req.Context())
	}
	if err != nil {
		http.Error(w, "unable to load receipts", http.StatusInternalServerError)
		return
	}
	r.render(w, req, "receipts", map[string]interface{}{
		"Title":      "Receipts",
		"ActivePage": "receipts",
		"Receipts":   receipts,
	})
}

func (r *Router) ReceiptFormNew(w http.ResponseWriter, req *http.Request) {
	patients, err := r.patients.Search(req.Context(), "")
	if err != nil {
		http.Error(w, "unable to load patients", http.StatusInternalServerError)
		return
	}
	r.render(w, req, "receipt-form", map[string]interface{}{
		"Title":      "New Receipt",
		"ActivePage": "receipts",
		"Receipt": &models.Receipt{
			DiscountType: models.DiscountNone,
			LineItems:    []models.LineItem{{Quantity: 1}},
		},
		"Patients": patients,
	})
}

func (r *Router) ReceiptCreate(w http.ResponseWriter, req *http.Request) {
	receipt, err := receiptFromForm(req)
	if err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	if err := r.receipts.CreateDraft(req.Context(), receipt); err != nil {
		r.renderReceiptFormError(w, req, receipt, err.Error())
		return
	}
	http.Redirect(w, req, "/receipts/"+receipt.ID, http.StatusSeeOther)
}

func (r *Router) ReceiptView(w http.ResponseWriter, req *http.Request) {
	receipt, err := r.receipts.Get(req.Context(), chi.URLParam(req, "id"))
	if err != nil {
		http.Error(w, "receipt not found", http.StatusNotFound)
		return
	}
	r.render(w, req, "receipt-view", map[string]interface{}{
		"Title":      "Receipt",
		"ActivePage": "receipts",
		"Receipt":    receipt,
	})
}

func (r *Router) ReceiptFormEdit(w http.ResponseWriter, req *http.Request) {
	receipt, err := r.receipts.Get(req.Context(), chi.URLParam(req, "id"))
	if err != nil {
		http.Error(w, "receipt not found", http.StatusNotFound)
		return
	}
	if receipt.Status != models.StatusDraft {
		http.Error(w, "finalized receipts are immutable", http.StatusConflict)
		return
	}
	patients, err := r.patients.Search(req.Context(), "")
	if err != nil {
		http.Error(w, "unable to load patients", http.StatusInternalServerError)
		return
	}
	r.render(w, req, "receipt-form", map[string]interface{}{
		"Title":      "Edit Receipt",
		"ActivePage": "receipts",
		"Receipt":    receipt,
		"Patients":   patients,
	})
}

func (r *Router) ReceiptUpdate(w http.ResponseWriter, req *http.Request) {
	receipt, err := receiptFromForm(req)
	if err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	receipt.ID = chi.URLParam(req, "id")
	if err := r.receipts.UpdateDraft(req.Context(), receipt); err != nil {
		r.renderReceiptFormError(w, req, receipt, err.Error())
		return
	}
	http.Redirect(w, req, "/receipts/"+receipt.ID, http.StatusSeeOther)
}

func (r *Router) ReceiptDelete(w http.ResponseWriter, req *http.Request) {
	if err := r.receipts.DeleteDraft(req.Context(), chi.URLParam(req, "id")); err != nil {
		http.Error(w, "receipt cannot be deleted", http.StatusConflict)
		return
	}
	http.Redirect(w, req, "/receipts", http.StatusSeeOther)
}

func (r *Router) ReceiptFinalize(w http.ResponseWriter, req *http.Request) {
	prefix := "RCP"
	_ = r.db.QueryRowContext(req.Context(), `SELECT receipt_prefix FROM settings WHERE id = 1`).Scan(&prefix)
	if err := r.receipts.Finalize(req.Context(), chi.URLParam(req, "id"), prefix); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	http.Redirect(w, req, "/receipts/"+chi.URLParam(req, "id"), http.StatusSeeOther)
}

func receiptFromForm(req *http.Request) (*models.Receipt, error) {
	if err := req.ParseForm(); err != nil {
		return nil, err
	}
	discountValue, err := strconv.Atoi(req.FormValue("discount_value"))
	if err != nil && req.FormValue("discount_value") != "" {
		return nil, err
	}
	descriptions := req.Form["description"]
	quantities := req.Form["quantity"]
	unitPrices := req.Form["unit_price"]
	if len(descriptions) != len(quantities) || len(descriptions) != len(unitPrices) {
		return nil, fmt.Errorf("line item fields are mismatched")
	}
	items := make([]models.LineItem, 0, len(descriptions))
	for index := range descriptions {
		quantity, err := strconv.Atoi(quantities[index])
		if err != nil {
			return nil, err
		}
		unitPrice, err := strconv.Atoi(unitPrices[index])
		if err != nil {
			return nil, err
		}
		items = append(items, models.LineItem{
			Description: descriptions[index],
			Quantity:    quantity,
			UnitPrice:   unitPrice,
		})
	}
	return &models.Receipt{
		PatientID:     req.FormValue("patient_id"),
		VisitDate:     req.FormValue("visit_date"),
		Diagnosis:     req.FormValue("diagnosis"),
		DiscountType:  models.DiscountType(req.FormValue("discount_type")),
		DiscountValue: discountValue,
		LineItems:     items,
	}, nil
}

func (r *Router) renderReceiptFormError(w http.ResponseWriter, req *http.Request, receipt *models.Receipt, message string) {
	patients, _ := r.patients.Search(req.Context(), "")
	r.render(w, req, "receipt-form", map[string]interface{}{
		"Title":      "Receipt",
		"ActivePage": "receipts",
		"Receipt":    receipt,
		"Patients":   patients,
		"Error":      message,
	})
}
