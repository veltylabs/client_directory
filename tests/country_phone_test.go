package tests

import (
	"strings"
	"testing"

	"github.com/veltylabs/patient_directory/kinds"
	"webtyp.com/country"
	"webtyp.com/crud"
	"webtyp.com/model"
	"webtyp.com/orm"
	"webtyp.com/storage/mem"

	patientdirectory "github.com/veltylabs/patient_directory"
)

func createPatient(t *testing.T, p *patientdirectory.Patient) error {
	t.Helper()
	mod, _, err := setupModule()
	if err != nil {
		t.Fatal(err)
	}
	caller := crud.LoopbackActingAs(nil, "u1", mod)
	var saved patientdirectory.Patient
	var callErr error
	caller.Call(opOf(model.Create), &crud.CreateArgs{Records: []model.Model{p}}, &saved, func(e error) { callErr = e })
	return callErr
}

// The module holds no country of its own: without the installation's default it does not start.
func TestNew_RequiresDefaultCountry(t *testing.T) {
	_, err := patientdirectory.New(orm.New(mem.New()), patientdirectory.Deps{
		Publisher: &fakePublisher{}, ValidateRUT: fakeValidateRUT, TenantID: "t1", IDs: &testIDGen{},
	})
	if err == nil || !strings.Contains(err.Error(), "Countries") {
		t.Fatalf("New without Deps.Countries must fail naming it, got %v", err)
	}
}

// A passport or a foreign id card is not valid without the country that issued it.
func TestCreate_ForeignDocumentNeedsIssuingCountry(t *testing.T) {
	err := createPatient(t, &patientdirectory.Patient{Id: "p1", DocumentType: kinds.DocumentPassport, Rut: "ab123456", Name: "Lucía"})
	if model.FieldOf(err) != "document_country" || !strings.Contains(err.Error(), string(patientdirectory.ErrDocumentCountry)) {
		t.Fatalf("want a field error on document_country, got %v", err)
	}
	err = createPatient(t, &patientdirectory.Patient{Id: "p2", DocumentType: kinds.DocumentPassport, DocumentCountry: "ZZ", Rut: "ab123456", Name: "Lucía"})
	if model.FieldOf(err) != "document_country" {
		t.Fatalf("an unknown country must be refused too, got %v", err)
	}
	if err := createPatient(t, &patientdirectory.Patient{Id: "p3", DocumentType: kinds.DocumentPassport, DocumentCountry: "VE", Rut: "ab123456", Name: "Lucía"}); err != nil {
		t.Fatalf("a passport with its country must be accepted: %v", err)
	}
}

// The phone is always international: a screen or a caller that goes through crud must send the prefix;
// an import or another module's call (plain strings) takes the default country.
func TestPhone_InternationalThroughCrudDefaultCountryThroughImport(t *testing.T) {
	if err := createPatient(t, &patientdirectory.Patient{Id: "p1", DocumentType: kinds.DocumentRUT, Rut: "11111111-1", Name: "Juan", Phone: "9 1111 1111"}); model.FieldOf(err) != "phone" {
		t.Fatalf("a phone with no prefix must be refused on the phone field, got %v", err)
	}
	if err := createPatient(t, &patientdirectory.Patient{Id: "p2", DocumentType: kinds.DocumentRUT, Rut: "11111111-1", Name: "Juan", Phone: "+56911111111"}); err != nil {
		t.Fatalf("an international phone must be accepted: %v", err)
	}
	mod, _, err := setupModule()
	if err != nil {
		t.Fatal(err)
	}
	got, err := mod.ImportPatient(patientdirectory.Patient{Id: "p3", DocumentType: kinds.DocumentRUT, Rut: "11111111-1", Name: "Juan", Phone: "9 1111 1111"})
	if err != nil || got.Phone != "+56911111111" {
		t.Fatalf("import phone = %q, %v; want +56911111111 (default country CL)", got.Phone, err)
	}
}

// A new record starts with the RUT in the country that has one and with a passport anywhere else.
func TestInitialDocument_FollowsTheDefaultCountry(t *testing.T) {
	if got := kinds.InitialDocument(country.Catalog{Default: "CL"}); got != kinds.DocumentRUT {
		t.Errorf("CL: got %q, want %q", got, kinds.DocumentRUT)
	}
	if got := kinds.InitialDocument(country.Catalog{Default: "AR"}); got != kinds.DocumentPassport {
		t.Errorf("AR: got %q, want %q", got, kinds.DocumentPassport)
	}
}
