package p_uniquity_finance_accounts

import (
	"context"
	"fmt"

	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/registry"
	"gorm.io/gorm"
)

// SourceDoc points at an arbitrary backing row identified by polymorphic Type + SourceDocID.
// No database foreign key constraint is enforced on SourceDocID.
type SourceDoc struct {
	gorm.Model

	Type        string `gorm:"column:source_doc_type;not null"`
	SourceDocID uint   `gorm:"column:source_doc_id;not null"`
}

// SourceDocTypeInterface describes how one document kind participates in linking and URLs.
type SourceDocTypeInterface interface {
	GetSourceDocType() string
	GetterDetailUrl(idKey string) getters.Getter[string]
	LoadFromID(ctx context.Context, id uint) (SourceDocInstanceInterface, error)
}

// SourceDocInstanceInterface is the loaded document for a resolved type/id pair.
type SourceDocInstanceInterface interface {
	GetSourceDocType() string
	GetSourceDocID() uint
	GetDetailUrl() string
}

// RegistrySourceDocTypes maps [SourceDocTypeInterface.GetSourceDocType] to loader/getter implementations.
var RegistrySourceDocTypes = registry.NewRegistry[SourceDocTypeInterface]()

// ResolveSourceDocInstance looks up a registered type and loads its instance by primary key.
func ResolveSourceDocInstance(ctx context.Context, typ string, id uint) (SourceDocInstanceInterface, error) {
	if typ == "" {
		return nil, fmt.Errorf("p_uniquity_finance_accounts: ResolveSourceDocInstance: empty type")
	}
	loader, ok := RegistrySourceDocTypes.Get(typ)
	if !ok {
		return nil, fmt.Errorf("p_uniquity_finance_accounts: ResolveSourceDocInstance: unknown type %q", typ)
	}
	inst, err := loader.LoadFromID(ctx, id)
	if err != nil {
		return nil, err
	}
	if inst.GetSourceDocType() != typ {
		return nil, fmt.Errorf("p_uniquity_finance_accounts: ResolveSourceDocInstance: type mismatch: registry key %q, instance %q", typ, inst.GetSourceDocType())
	}
	return inst, nil
}
