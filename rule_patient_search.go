package patientdirectory

import (
	"strings"

	"webtyp.com/crud"
	"webtyp.com/input"
	"webtyp.com/storage"
)

// noID es el id que ningún paciente tiene: crud exige un id en cada alta, así que
// buscarlo no devuelve filas.
const noID = ""

// Search traduce el texto que escribió la persona en las condiciones que buscan patient. Solo condiciones unidas con AND (nunca OR).
// crud ya aplicó el tope de página, el cursor y la clínica; tú solo dices qué filas corresponden.
//
// Sin texto no se muestra ningún paciente: el registro arranca vacío y aparece al buscar.
// Un RUT, un documento (pasaporte, cédula extranjera, sin importar el tipo) o un teléfono
// se busca por ese dato; cualquier otro texto busca por nombre.
func (p *Patient) Search(a crud.CurrentUser, filter *Patient, text string, m *Module) ([]storage.Condition, error) {
	if text == "" {
		return []storage.Condition{storage.Eq(Patient_.Id, noID)}, nil
	}
	if norm, err := input.NormalizeRUT(text); err == nil {
		return []storage.Condition{storage.In(Patient_.Rut, []string{norm, strings.ReplaceAll(norm, "-", "")})}, nil
	}
	if isPhone(text) {
		return []storage.Condition{storage.Like(Patient_.Phone, "%"+digitsOnly(text)+"%")}, nil
	}
	if isDocument(text) {
		// Sin distinguir mayúsculas: el recepcionista escribe el pasaporte como puede.
		return []storage.Condition{storage.ILike(Patient_.Rut, text)}, nil
	}
	return []storage.Condition{storage.Like(Patient_.Name, "%"+strings.ToLower(text)+"%")}, nil
}
