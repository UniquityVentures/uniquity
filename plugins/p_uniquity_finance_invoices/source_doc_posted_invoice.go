package p_uniquity_finance_invoices

import (
	"context"
	"strconv"

	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	finance_accounts "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
)

// PostedInvoiceSourceDocType is stored on [finance_accounts.SourceDoc.Type] for posted invoices.
const PostedInvoiceSourceDocType = "p_uniquity_finance_invoices.PostedInvoice"

type postedInvoiceSourceDocType struct{}

func (postedInvoiceSourceDocType) GetSourceDocType() string { return PostedInvoiceSourceDocType }

func (postedInvoiceSourceDocType) GetterDetailUrl(idKey string) getters.Getter[string] {
	return lamu.RoutePath("finance_invoices.PostedInvoiceDetailRoute", map[string]getters.Getter[any]{
		"id": getters.Any(getters.Key[uint](idKey)),
	})
}

func (postedInvoiceSourceDocType) LoadFromID(ctx context.Context, id uint) (finance_accounts.SourceDocInstanceInterface, error) {
	db, err := getters.DBFromContext(ctx)
	if err != nil {
		return nil, err
	}
	var row PostedInvoice
	if err := db.WithContext(ctx).Where("id = ?", id).Take(&row).Error; err != nil {
		return nil, err
	}
	return &postedInvoiceSourceInst{row: row}, nil
}

type postedInvoiceSourceInst struct {
	row PostedInvoice
}

func (i *postedInvoiceSourceInst) GetSourceDocType() string { return PostedInvoiceSourceDocType }
func (i *postedInvoiceSourceInst) GetSourceDocID() uint     { return i.row.ID }
func (i *postedInvoiceSourceInst) GetDetailUrl() string {
	return AppUrl + "posted/" + strconv.FormatUint(uint64(i.row.ID), 10) + "/"
}

func init() {
	if err := finance_accounts.RegistrySourceDocTypes.Register(PostedInvoiceSourceDocType, postedInvoiceSourceDocType{}); err != nil {
		panic(err)
	}
}
