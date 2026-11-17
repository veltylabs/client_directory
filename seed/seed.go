package seed

import (
	clientdirectory "github.com/veltylabs/client_directory"
	"github.com/veltylabs/client_directory/kinds"
)

// Data contiene las filas de demostración creadas por Load.
type Data struct {
	Patients []clientdirectory.Client
}

// Load inserta los datos semilla en el módulo usando sus métodos de dominio.
func Load(m *clientdirectory.Module, tenantID string) (Data, error) {
	items := []clientdirectory.Client{
		{
			Id:        "demo-patient-1",
			TenantId:  tenantID,
			Rut:       "11111111-1",
			Name:      "Juan Pérez",
			Birthdate: "1980-05-12",
			Phone:     "+56911111111",
			Email:     "juan.perez@example.cl",
			Address:   "Chillan",
			IsActive:  true,
		},
		{
			Id:        "demo-patient-2",
			TenantId:  tenantID,
			Rut:       "12345678-5",
			Name:      "María González",
			Birthdate: "1992-11-03",
			Phone:     "+56922222222",
			Email:     "maria.gonzalez@example.cl",
			Address:   "Chillan",
			IsActive:  true,
		},
		{
			Id:        "demo-patient-3",
			TenantId:  tenantID,
			Rut:       "15678432-K",
			Name:      "Pedro Soto",
			Birthdate: "1975-02-20",
			Phone:     "+56933333333",
			Email:     "pedro.soto@example.cl",
			Address:   "San Carlos",
			IsActive:  true,
		},
		// Pacientes de otras nacionalidades: un pasaporte, una cédula extranjera y
		// otro pasaporte. Se buscan por su documento, no por RUT.
		{
			Id:              "demo-patient-4",
			TenantId:        tenantID,
			DocumentType:    kinds.DocumentPassport,
			DocumentCountry: "VE",
			Rut:             "AB1234567",
			Name:            "Lucía Mendoza",
			Birthdate:       "1988-07-30",
			Phone:           "+584141234567",
			Email:           "lucia.mendoza@example.com",
			Address:         "Chillan",
			IsActive:        true,
		},
		{
			Id:              "demo-patient-5",
			TenantId:        tenantID,
			DocumentType:    kinds.DocumentForeignID,
			DocumentCountry: "PE",
			Rut:             "CE-004567",
			Name:            "Carlos Ramírez",
			Birthdate:       "1979-09-14",
			Phone:           "+51987654321",
			Email:           "carlos.ramirez@example.com",
			Address:         "San Carlos",
			IsActive:        true,
		},
		{
			Id:              "demo-patient-6",
			TenantId:        tenantID,
			DocumentType:    kinds.DocumentPassport,
			DocumentCountry: "HT",
			Rut:             "HT2024567",
			Name:            "Jean Pierre Louis",
			Birthdate:       "1995-01-08",
			Phone:           "+50937123456",
			Email:           "jean.louis@example.com",
			Address:         "Chillan",
			IsActive:        true,
		},
	}

	var data Data
	for i := range items {
		p, err := m.ImportPatient(items[i])
		if err != nil {
			return Data{}, err
		}
		data.Patients = append(data.Patients, p)
	}

	return data, nil
}
