package tests

import (
	"errors"
	"strings"
	"webtyp.com/country"

	"webtyp.com/events"
	"webtyp.com/orm"
	"webtyp.com/storage/mem"

	clientdirectory "github.com/veltylabs/client_directory"
)

type fakePublisher struct {
	published []events.Event
}

func (f *fakePublisher) Publish(e events.Event) {
	f.published = append(f.published, e)
}

func fakeValidateRUT(rut string) (string, error) {
	if rut == "bad" {
		return "", errors.New("bad check digit")
	}
	cleaned := strings.ReplaceAll(rut, ".", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ToUpper(cleaned)
	if cleaned == "" {
		return "", clientdirectory.ErrDocumentNumber
	}
	return cleaned, nil
}

func setupModule() (*clientdirectory.Module, *fakePublisher, error) {
	return setupModuleWithTenant("test-tenant")
}

func setupModuleWithTenant(tenantID string) (*clientdirectory.Module, *fakePublisher, error) {
	conn := mem.New()
	db := orm.New(conn)
	return setupModuleWithTenantAndDB(tenantID, db)
}

func setupModuleWithTenantAndDB(tenantID string, db *orm.DB) (*clientdirectory.Module, *fakePublisher, error) {
	pub := &fakePublisher{}
	deps := clientdirectory.Deps{
		Countries:   testCountries,
		Publisher:   pub,
		ValidateRUT: fakeValidateRUT,
		TenantID:    tenantID,
		IDs:         &testIDGen{},
	}
	m, err := clientdirectory.New(db, deps)
	if err != nil {
		return nil, nil, err
	}
	return m, pub, nil
}

// testCountries is what an app decides for the tests: Chile by default.
var testCountries = country.Catalog{Default: "CL", First: []string{"VE", "PE", "HT"}}
