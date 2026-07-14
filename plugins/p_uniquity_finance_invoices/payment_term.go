package p_uniquity_finance_invoices

import (
	"context"
	"fmt"
	"time"

	"github.com/lariv-in/lago/getters"
	"github.com/lariv-in/lago/registry"
)

// Discriminator strings stored on [PaymentTerm.Type] for polymorphic backing rows.
const (
	PaymentTermTypeDueDate  = "p_uniquity_finance_invoices.PaymentTermDueDate"
	PaymentTermTypeRelative = "p_uniquity_finance_invoices.PaymentTermRelative"
)

// PaymentTermTypeInterface describes how one payment-term kind participates in polymorphic resolution.
type PaymentTermTypeInterface interface {
	GetPaymentTermType() string
	LoadFromID(ctx context.Context, id uint) (PaymentTermInstanceInterface, error)
}

// PaymentTermInstanceInterface is the loaded backing row for a type/id pair.
type PaymentTermInstanceInterface interface {
	GetPaymentTermType() string
	// GetPaymentTermID returns the backing row primary key.
	GetPaymentTermID() uint
	Summary() string
}

// RegistryPaymentTermTypes maps [PaymentTermTypeInterface.GetPaymentTermType] to loaders.
var RegistryPaymentTermTypes = registry.NewRegistry[PaymentTermTypeInterface]()

// ResolvePaymentTermInstance looks up a registered type and loads its backing row by primary key.
func ResolvePaymentTermInstance(ctx context.Context, typ string, id uint) (PaymentTermInstanceInterface, error) {
	if typ == "" {
		return nil, fmt.Errorf("p_uniquity_finance_invoices: ResolvePaymentTermInstance: empty type")
	}
	loader, ok := RegistryPaymentTermTypes.Get(typ)
	if !ok {
		return nil, fmt.Errorf("p_uniquity_finance_invoices: ResolvePaymentTermInstance: unknown type %q", typ)
	}
	inst, err := loader.LoadFromID(ctx, id)
	if err != nil {
		return nil, err
	}
	if inst.GetPaymentTermType() != typ {
		return nil, fmt.Errorf("p_uniquity_finance_invoices: ResolvePaymentTermInstance: type mismatch: registry key %q, instance %q", typ, inst.GetPaymentTermType())
	}
	return inst, nil
}

// ResolvePaymentTermInstanceFromTerm loads the backing instance for a stored [PaymentTerm] row.
func ResolvePaymentTermInstanceFromTerm(ctx context.Context, pt *PaymentTerm) (PaymentTermInstanceInterface, error) {
	if pt == nil {
		return nil, fmt.Errorf("p_uniquity_finance_invoices: ResolvePaymentTermInstanceFromTerm: nil term")
	}
	return ResolvePaymentTermInstance(ctx, pt.Type, pt.BackingID)
}

type paymentTermDueDateType struct{}

func (paymentTermDueDateType) GetPaymentTermType() string { return PaymentTermTypeDueDate }

func (paymentTermDueDateType) LoadFromID(ctx context.Context, id uint) (PaymentTermInstanceInterface, error) {
	db, err := getters.DBFromContext(ctx)
	if err != nil {
		return nil, err
	}
	var row PaymentTermDueDate
	if err := db.WithContext(ctx).Where("id = ?", id).Take(&row).Error; err != nil {
		return nil, err
	}
	return &paymentTermDueDateInstance{row: row}, nil
}

type paymentTermDueDateInstance struct {
	row PaymentTermDueDate
}

func (i *paymentTermDueDateInstance) GetPaymentTermType() string { return PaymentTermTypeDueDate }
func (i *paymentTermDueDateInstance) GetPaymentTermID() uint     { return i.row.ID }
func (i *paymentTermDueDateInstance) Summary() string {
	if i.row.Datetime.IsZero() {
		return "Due date"
	}
	return i.row.Datetime.Format(time.RFC3339)
}

type paymentTermRelativeType struct{}

func (paymentTermRelativeType) GetPaymentTermType() string { return PaymentTermTypeRelative }

func (paymentTermRelativeType) LoadFromID(ctx context.Context, id uint) (PaymentTermInstanceInterface, error) {
	db, err := getters.DBFromContext(ctx)
	if err != nil {
		return nil, err
	}
	var row PaymentTermRelative
	if err := db.WithContext(ctx).Where("id = ?", id).Take(&row).Error; err != nil {
		return nil, err
	}
	return &paymentTermRelativeInstance{row: row}, nil
}

type paymentTermRelativeInstance struct {
	row PaymentTermRelative
}

func (i *paymentTermRelativeInstance) GetPaymentTermType() string { return PaymentTermTypeRelative }
func (i *paymentTermRelativeInstance) GetPaymentTermID() uint     { return i.row.ID }
func (i *paymentTermRelativeInstance) Summary() string {
	if i.row.Duration <= 0 {
		return fmt.Sprintf("Relative #%d", i.row.ID)
	}
	return i.row.Duration.String()
}

func init() {
	if err := RegistryPaymentTermTypes.Register(PaymentTermTypeDueDate, paymentTermDueDateType{}); err != nil {
		panic(err)
	}
	if err := RegistryPaymentTermTypes.Register(PaymentTermTypeRelative, paymentTermRelativeType{}); err != nil {
		panic(err)
	}
}
