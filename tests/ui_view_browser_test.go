package tests

import (
	"github.com/veltylabs/client_directory/kinds"
	"testing"
	"webtyp.com/model"
	"webtyp.com/view"

	clientdirectory "github.com/veltylabs/client_directory"
	"github.com/veltylabs/client_directory/ui"
)

func TestUI_Screens(t *testing.T) {
	caller := view.NewRecords().Caller(&mockCaller{})
	screens := ui.Screens(view.BrowserDeps{Caller: caller, Countries: testCountries})
	if len(screens) != 1 {
		t.Fatalf("expected 1 screen, got %d", len(screens))
	}
	s := screens[0]
	if err := s.Validate(); err != nil {
		t.Fatalf("screen validation failed: %v", err)
	}
	if s.ID != ui.ID {
		t.Errorf("ID = %q, want %q", s.ID, ui.ID)
	}
	if string(s.Permission) != string(model.ResourceOf(&clientdirectory.ClientModel)) {
		t.Errorf("Permission = %q, want patient", s.Permission)
	}
	if s.Presenter.Description() == "" {
		t.Errorf("presenter description is empty")
	}
}

// The "+" button's draft starts with the installation's default document type.
func TestUI_NewRecordStartsWithTheDocumentOfTheDefaultCountry(t *testing.T) {
	foreign := testCountries
	foreign.Default = "AR"
	screens := ui.Screens(view.BrowserDeps{Caller: view.NewRecords().Caller(&mockCaller{}), Countries: foreign})
	p, ok := screens[0].Presenter.NewRecord().(*clientdirectory.Client)
	if !ok || p.DocumentType != kinds.DocumentPassport {
		t.Fatalf("new record must start with passport, got %#v", screens[0].Presenter.NewRecord())
	}
}
