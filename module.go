package patientdirectory

import (
	"strings"

	"github.com/veltylabs/patient_directory/kinds"
	"webtyp.com/country"
	"webtyp.com/events"
	"webtyp.com/fmt"
	"webtyp.com/model"
	"webtyp.com/orm"
)

type Deps struct {
	Publisher   events.Publisher // opcional — nil deshabilita la publicación silenciosamente
	TenantID    string
	IDs         model.IDGenerator // required: ids for patients created without one
	ValidateRUT func(string) (string, error)
	// Countries is the installation's default country and which countries come first (App.Countries):
	// required. The module holds no country code or prefix of its own.
	Countries country.Catalog
}

type Module struct {
	db          *orm.DB
	store       *orm.Store
	pub         events.Publisher
	validateRUT func(string) (string, error)
	tenantID    string
	countries   country.Catalog
}

func depErr(name string) error {
	return fmt.Err("patient_directory: Deps."+name, "is required")
}

// New conecta el módulo a un *orm.DB ya conectado; se asume que el esquema
// ya existe —consulte el submódulo migrate.
func New(db *orm.DB, deps Deps) (*Module, error) {
	if deps.ValidateRUT == nil {
		return nil, depErr("ValidateRUT")
	}
	if deps.Countries.Default == "" {
		return nil, depErr("Countries.Default")
	}
	if err := deps.Countries.Validate(); err != nil {
		return nil, err
	}
	if deps.TenantID == "" {
		return nil, depErr("TenantID")
	}
	if deps.IDs == nil {
		return nil, depErr("IDs")
	}
	return &Module{
		db:          db,
		store:       orm.TenantStore(db, deps.TenantID, deps.IDs),
		pub:         deps.Publisher,
		validateRUT: deps.ValidateRUT,
		tenantID:    deps.TenantID,
		countries:   deps.Countries,
	}, nil
}

// checkDocument normalizes and validates p's document by its type. An empty type means RUT (what the form sends
// today). Only a RUT is validated; other documents are stored as typed, upper-cased and trimmed.
func (m *Module) checkDocument(p *Patient) error {
	if p.DocumentType == "" {
		p.DocumentType = kinds.DocumentRUT
	}
	p.Rut = strings.TrimSpace(p.Rut)
	p.DocumentCountry = strings.ToUpper(strings.TrimSpace(p.DocumentCountry))
	switch p.DocumentType {
	case kinds.DocumentRUT:
		norm, err := m.validateRUT(p.Rut)
		if err != nil {
			return &model.FieldError{Field: "rut", Err: err}
		}
		p.Rut, p.DocumentCountry = norm, ""
	case kinds.DocumentPassport, kinds.DocumentForeignID:
		if p.Rut == "" {
			return &model.FieldError{Field: "rut", Err: ErrDocumentNumber}
		}
		p.Rut = strings.ToUpper(p.Rut)
		if _, known := country.Find(p.DocumentCountry); !known {
			return &model.FieldError{Field: "document_country", Err: ErrDocumentCountry}
		}
	case kinds.DocumentNone:
		p.Rut, p.DocumentCountry = "", ""
	case kinds.DocumentUnverified:
		// stored as typed: the front desk fixes it
	default:
		return &model.FieldError{Field: "document_type", Err: ErrDocumentType}
	}
	return nil
}

// checkPhone writes the phone in international form (+56911111111). A number
// typed without a prefix takes the installation's default country; the form
// always sends one, so this covers imports and other callers.
func (m *Module) checkPhone(p *Patient) error {
	if p.Phone == "" {
		return nil
	}
	norm, err := country.Normalize(p.Phone, m.countries.Default)
	if err != nil {
		return &model.FieldError{Field: "phone", Err: err}
	}
	p.Phone = norm
	return nil
}
