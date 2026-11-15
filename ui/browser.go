package ui

import (
	patientdirectory "github.com/veltylabs/patient_directory"
	"webtyp.com/country"
	"webtyp.com/lang"
	"webtyp.com/model"
	"webtyp.com/router"
	"webtyp.com/view"
)

const descPatients lang.Text = "Patient directory: personal data, contact information and identification."

// Screens declara las pantallas que este módulo ofrece al menú de la app.
// countries lo decide la app (view.BrowserDeps.Countries): cuáles salen primero y cuál es la de por defecto.
func Screens(caller router.Caller, countries country.Catalog) []view.Screen {
	return []view.Screen{
		{
			ID:         ID,
			Label:      Label,
			Icon:       ID,
			Permission: model.ResourceOf(&patientdirectory.PatientModel),
			Presenter: patientdirectory.NewView(caller,
				view.WithDescription(descPatients),
				view.WithOptions("document_country", countries.Options()...),
				view.WithOptions("phone", countries.PrefixOptions()...),
			),
		},
	}
}
