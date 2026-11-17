package ui

import (
	patientdirectory "github.com/veltylabs/patient_directory"
	"github.com/veltylabs/patient_directory/kinds"
	"webtyp.com/lang"
	"webtyp.com/model"
	"webtyp.com/view"
)

const descPatients lang.Text = "Patient directory: personal data, contact information and identification."

// Screens declara las pantallas que este módulo ofrece al menú de la app.
// Los países los decide la app (view.BrowserDeps.Countries): cuáles salen primero y cuál es la de por defecto; de ella
// sale el tipo de documento con que parte un paciente nuevo (kinds.InitialDocument).
func Screens(d view.BrowserDeps) []view.Screen {
	caller, countries := d.Caller, d.Countries
	return []view.Screen{
		{
			ID:         ID,
			Label:      Label,
			Icon:       ID,
			Permission: model.ResourceOf(&patientdirectory.PatientModel),
			Presenter: patientdirectory.NewView(caller,
				view.WithDescription(descPatients),
				view.WithNewRecord(func() view.Record { return &patientdirectory.Patient{DocumentType: kinds.InitialDocument(countries)} }),
				view.WithOptions("document_country", countries.Options()...),
				view.WithOptions("phone", countries.PrefixOptions()...),
			),
		},
	}
}
