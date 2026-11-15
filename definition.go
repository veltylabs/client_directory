package patientdirectory

import (
	"github.com/veltylabs/patient_directory/kinds"
	"webtyp.com/components/agedate"
	"webtyp.com/input"
	"webtyp.com/model"
)

// PatientModel es el registro del establecimiento para las personas que atiende:
// IDENTIDAD y CONTACTO únicamente. Los datos clínicos —diagnósticos, atenciones, recetas—
// pertenecen a clinical_encounter y nunca deben agregarse aquí.
//
// rut NO lleva la marca Unique a propósito: la propiedad Unique de model.FieldDB es
// de columna única, y un RUT es único por TENANT, no globalmente. La restricción
// se aplica en CreatePatient contra ErrRutAlreadyExists —la misma regla de capa de
// aplicación que device_manager aplica a una IP de dispositivo.
var Definition = model.ModuleDefinition{Name: "patient_directory", Entities: []model.Definition{PatientModel}}

var PatientModel = model.Definition{
	Name:  "patient",
	Allow: model.Create | model.Read | model.Update,
	Rules: model.Rules{Create: model.AddMyBusinessRule, Read: model.AddMyBusinessRule, Update: model.AddMyBusinessRule},
	Fields: model.Fields{
		{Name: "id", Type: model.Text(), DB: &model.FieldDB{PK: true}, OmitEmpty: true},
		{Name: "tenant_id", Type: model.Text(), NotNull: true, ServerOnly: true},
		{Name: "document_type", Type: kinds.DocumentType(), Label: "Document type", NotNull: true, Item: model.ItemDescription},
		{Name: "rut", Type: input.Text(), Label: "Document number", Help: "Document number, as written on the document.", OmitEmpty: true, Permitted: model.Permitted{Letters: true, Numbers: true, Extra: []rune{'-', '.'}, Maximum: 30}, Item: model.ItemDescription},
		{Name: "document_country", Type: input.Select(), Label: "Issuing country", Help: "The country that issued the passport or id card.", OmitEmpty: true, Permitted: model.Permitted{Letters: true, Minimum: 2, Maximum: 2},
			ShownWhen: model.When{Field: "document_type", In: []string{kinds.DocumentPassport, kinds.DocumentForeignID}}},
		{Name: "name", Type: input.Text(), Label: "Name", NotNull: true, Permitted: model.Permitted{Minimum: 1, Maximum: 255}, Item: model.ItemLabel},
		{Name: "birthdate", Type: agedate.New(), Label: "Birthdate", OmitEmpty: true, AlsoSetsField: "birthdate_estimated"},
		{Name: "birthdate_estimated", Type: input.Checkbox(), Label: "Estimated birth date", Help: "Marked by itself when the date comes from the age the patient said, not from a document."},
		{Name: "sex", Type: input.Gender(), Label: "Sex", OmitEmpty: true},
		{Name: "phone", Type: input.Phone(), Label: "Phone", OmitEmpty: true},
		{Name: "email", Type: input.Email(), Label: "Email", OmitEmpty: true},
		{Name: "address", Type: input.Address(), Label: "Address", OmitEmpty: true},
		{Name: "notes", Type: input.Textarea(), Label: "Notes", Help: "Anything the front desk must know about this patient, e.g. the old paper file number.", OmitEmpty: true},
		{Name: "is_active", Type: input.Checkbox(), Label: "Active patient", Help: "Uncheck to retire the patient; the record is never deleted."},
		{Name: "updated_at", Type: model.NanoTime(), OmitEmpty: true, ServerOnly: true},
	},
}

const (
	ErrNotFound         model.NotFound = "patient not found"
	ErrRutAlreadyExists model.Conflict = "patient rut already exists in this tenant"
	ErrIdRequired       model.Invalid  = "patient id is required"
	ErrIdTaken          model.Conflict = "patient belongs to another tenant"

	ErrDocumentType   model.Invalid = "the document type must be RUT, passport, foreign id card, no document or to verify"
	ErrDocumentCountry model.Invalid = "the issuing country is required for a foreign document"
	ErrDocumentNumber model.Invalid = "a passport or foreign id card needs its number"
	ErrSex            model.Invalid = "sex must be female, male or not recorded"
)
