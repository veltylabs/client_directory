package tests

import (
	"testing"

	patientdirectory "github.com/veltylabs/patient_directory"
	"github.com/veltylabs/patient_directory/seed"
	"webtyp.com/crud"
	"webtyp.com/events/mock"
	"webtyp.com/input"
	"webtyp.com/model"
	"webtyp.com/orm"
	"webtyp.com/storage/mem"
)

func TestSeed_Load(t *testing.T) {
	db := orm.New(mem.New())
	mod, err := patientdirectory.New(db, patientdirectory.Deps{
		Countries:   testCountries,
		Publisher:   &mock.Broker{},
		TenantID:    "demo",
		IDs:         &testIDGen{},
		ValidateRUT: input.NormalizeRUT,
	})
	if err != nil {
		t.Fatalf("patientdirectory.New: %v", err)
	}

	data, err := seed.Load(mod, "demo")
	if err != nil {
		t.Fatalf("seed.Load: %v", err)
	}

	if len(data.Patients) != 6 {
		t.Fatalf("se esperaban 6 pacientes en seed.Data, se obtuvieron %d", len(data.Patients))
	}

	for _, p := range data.Patients {
		name, _, _, _, found, err := mod.GetPatient("demo", p.Id)
		if err != nil || !found || name != p.Name {
			t.Fatalf("mod.GetPatient(%s): found=%v err=%v name=%s", p.Id, found, err, name)
		}
	}

	// Segundo intento no debería duplicar ni fallar
	data2, err := seed.Load(mod, "demo")
	if err != nil {
		t.Fatalf("seed.Load 2: %v", err)
	}

	if len(data2.Patients) != 6 {
		t.Fatalf("se esperaban 6 pacientes en seed.Data en el segundo intento, se obtuvieron %d", len(data2.Patients))
	}
}

// Los pacientes extranjeros del seed se encuentran por su documento: un
// pasaporte o una cédula extranjera no caen en la búsqueda por nombre.
func TestSeed_ForeignPatientsFoundByDocument(t *testing.T) {
	db := orm.New(mem.New())
	mod, err := patientdirectory.New(db, patientdirectory.Deps{
		Countries:   testCountries,
		Publisher:   &mock.Broker{},
		TenantID:    "demo",
		IDs:         &testIDGen{},
		ValidateRUT: input.NormalizeRUT,
	})
	if err != nil {
		t.Fatalf("patientdirectory.New: %v", err)
	}
	if _, err := seed.Load(mod, "demo"); err != nil {
		t.Fatalf("seed.Load: %v", err)
	}
	caller := crud.LoopbackActingAs(nil, "u1", mod)

	for _, tc := range []struct{ document, wantID string }{
		{"AB1234567", "demo-patient-4"},
		{"CE-004567", "demo-patient-5"},
		{"HT2024567", "demo-patient-6"},
	} {
		var list patientdirectory.PatientList
		var callErr error
		caller.Call(opOf(model.Read),
			&crud.ReadArgs{Search: tc.document}, &list, func(e error) { callErr = e })
		if callErr != nil {
			t.Fatalf("search %s: %v", tc.document, callErr)
		}
		if len(list) != 1 || list[0].Id != tc.wantID {
			t.Errorf("search %s: want %s, got %v", tc.document, tc.wantID, list)
		}
	}
}
