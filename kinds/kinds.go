package kinds

import (
	"webtyp.com/country"
	"webtyp.com/fmt"
	"webtyp.com/input"
)

const (
	DocumentRUT        = "rut"        // Chilean RUT: check digit validated; unique per tenant
	DocumentPassport   = "passport"   // any country's passport, stored as typed
	DocumentForeignID  = "foreign_id" // another country's id card, stored as typed
	DocumentNone       = "none"       // the person has no document at hand
	DocumentUnverified = "unverified" // typed but invalid; the front desk must fix it
)

// RUTCountry is the country whose people have a RUT: an installation whose default country is this one starts a new
// record with a RUT; any other starts with a passport.
const RUTCountry = "CL"

// InitialDocument is the document type a new record starts with: the RUT when the installation's default country is
// RUTCountry, the passport otherwise. Example: p := &Patient{DocumentType: kinds.InitialDocument(countries)}
func InitialDocument(countries country.Catalog) string {
	if countries.Default == RUTCountry {
		return DocumentRUT
	}
	return DocumentPassport
}

// rutName is the national id name, the same in every language: data, not UI text.
const rutName = "RUT"

// DocumentType: closed options for patient.document_type. Example: in := kinds.DocumentType()
func DocumentType() input.Input {
	return input.Radio(
		fmt.KeyValue{Key: DocumentRUT, Value: rutName},
		fmt.KeyValue{Key: DocumentPassport, Value: "Passport"},
		fmt.KeyValue{Key: DocumentForeignID, Value: "Foreign id card"},
		fmt.KeyValue{Key: DocumentNone, Value: "No document"},
		fmt.KeyValue{Key: DocumentUnverified, Value: "To verify"},
	)
}
