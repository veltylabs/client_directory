package clientdirectory

import "webtyp.com/crud"

// StoreClient is the tenant-sealed store of the patients.
func (m *Module) StoreClient() crud.Store { return m.store }
