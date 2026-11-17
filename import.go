package clientdirectory

import (
	"webtyp.com/crud"
	"webtyp.com/model"
	"webtyp.com/time"
)

var importUpdateFields = []string{
	Client_.Rut,
	Client_.DocumentType,
	Client_.DocumentCountry,
	Client_.Name,
	Client_.Birthdate,
	Client_.BirthdateEstimated,
	Client_.Sex,
	Client_.Phone,
	Client_.Email,
	Client_.Address,
	Client_.Notes,
	Client_.IsActive,
	Client_.UpdatedAt,
}

// ImportPatient guarda un paciente de un sistema anterior con el ID que tenía allí:
// lo crea, o actualiza todos los campos si ese ID ya existe, de modo que importar
// dos veces deja la misma fila. Valida exactamente como CreatePatient y no publica eventos.
func (m *Module) ImportPatient(p Client) (Client, error) {
	if p.Id == "" {
		return Client{}, &model.FieldError{Field: Client_.Id, Err: ErrIdRequired}
	}
	if p.TenantId != "" && p.TenantId != m.tenantID {
		return Client{}, ErrIdTaken
	}
	p.TenantId = m.tenantID
	if err := m.checkDocument(&p); err != nil {
		return Client{}, err
	}
	if err := m.checkPhone(&p); err != nil {
		return Client{}, err
	}
	if err := model.ValidateFields(model.ActionCreate, &p); err != nil {
		return Client{}, err
	}
	p.UpdatedAt = time.NowNano()

	err := m.store.Atomic(func(tx crud.Tx) error {
		if err := m.sameDocument(tx, &p); err != nil {
			return err
		}
		var existing Client
		err := tx.ReadOne(&existing, p.Id)
		if err == nil {
			return tx.UpdateFields(&p, p.Id, importUpdateFields)
		}
		if model.KindOf(err) == model.ErrorNotFound {
			return tx.Create(&p)
		}
		return err
	})
	if err != nil {
		return Client{}, err
	}
	return p, nil
}
