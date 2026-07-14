package p_uniquity_finance_invoices

import (
	"context"
	"strconv"

	finance_accounts "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/getters"
)

type paymentSourceDocType struct{}

func (paymentSourceDocType) GetSourceDocType() string { return PaymentSourceDocType }

func (paymentSourceDocType) GetterDetailUrl(idKey string) getters.Getter[string] {
	return lago.RoutePath("finance_invoices.PaymentDetailRoute", map[string]getters.Getter[any]{
		"id": getters.Any(getters.Key[uint](idKey)),
	})
}

func (paymentSourceDocType) LoadFromID(ctx context.Context, id uint) (finance_accounts.SourceDocInstanceInterface, error) {
	db, err := getters.DBFromContext(ctx)
	if err != nil {
		return nil, err
	}
	var row Payment
	if err := db.WithContext(ctx).Where("id = ?", id).Take(&row).Error; err != nil {
		return nil, err
	}
	return &paymentSourceInst{row: row}, nil
}

type paymentSourceInst struct {
	row Payment
}

func (i *paymentSourceInst) GetSourceDocType() string { return PaymentSourceDocType }
func (i *paymentSourceInst) GetSourceDocID() uint     { return i.row.ID }
func (i *paymentSourceInst) GetDetailUrl() string {
	return AppUrl + "payments/" + strconv.FormatUint(uint64(i.row.ID), 10) + "/"
}

func init() {
	if err := finance_accounts.RegistrySourceDocTypes.Register(PaymentSourceDocType, paymentSourceDocType{}); err != nil {
		panic(err)
	}
}
