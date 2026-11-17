package clientdirectory

import "webtyp.com/crud"

// Update es la regla de negocio de actualizar patient: old es la fila guardada y el receptor trae solo los campos enviados (fields).
// Vuelve a validar el documento solo si cambió alguno de sus campos.
// crud ya validó el esquema, los campos solo del servidor, los límites y el permiso antes de llamarte.
func (next *Client) Update(a crud.CurrentUser, tx crud.Tx, old *Client, fields []string, m *Module) (crud.RuleAnswer, error) {
	if touches(fields, Client_.Phone) {
		if err := m.checkPhone(next); err != nil {
			return crud.RuleAnswer{}, err
		}
	}
	if !touches(fields, Client_.Rut, Client_.DocumentType, Client_.DocumentCountry) {
		return crud.ApplyAndPublish(), nil
	}
	// next trae solo los campos enviados: se mezcla con lo guardado campo por campo.
	merged := *old
	if touches(fields, Client_.Rut) {
		merged.Rut = next.Rut
	}
	if touches(fields, Client_.DocumentType) {
		merged.DocumentType = next.DocumentType
	}
	if touches(fields, Client_.DocumentCountry) {
		merged.DocumentCountry = next.DocumentCountry
	}
	if err := m.checkDocument(&merged); err != nil {
		return crud.RuleAnswer{}, err
	}
	// Solo se escriben las columnas de fields: se devuelve la forma normalizada de las enviadas.
	next.Rut, next.DocumentType, next.DocumentCountry = merged.Rut, merged.DocumentType, merged.DocumentCountry
	if err := m.sameDocument(tx, &merged); err != nil {
		return crud.RuleAnswer{}, err
	}
	return crud.ApplyAndPublish(), nil
}

func touches(fields []string, names ...string) bool {
	for _, f := range fields {
		for _, n := range names {
			if f == n {
				return true
			}
		}
	}
	return false
}
