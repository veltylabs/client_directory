package migrate_test

import (
	"testing"
	"webtyp.com/storage"
	"webtyp.com/storage/mem"

	"webtyp.com/ddl"
	"webtyp.com/model"

	"github.com/veltylabs/patient_directory/migrate"
)

// dummyExecer es la base que ve la migración: anota el DDL que recibe y, como ddl reconstruye
// text_search leyendo filas, deja pasar las lecturas y escrituras a una base en memoria.
type dummyExecer struct {
	storage.Conn
	queries []string
}

func (e *dummyExecer) Exec(query string, args ...any) (storage.Result, error) {
	if query == "mem" { // lo que compila la base en memoria para las lecturas y escrituras
		return e.Conn.Exec(query, args...)
	}
	e.queries = append(e.queries, query)
	return storage.Affected(0), nil
}

type dummyCompiler struct{}

func (c *dummyCompiler) CompileDDL(s ddl.Stmt, m model.Model) (string, []any, error) {
	return "CREATE TABLE IF NOT EXISTS patient (id TEXT PRIMARY KEY)", nil, nil
}

func TestMigrate(t *testing.T) {
	execer := &dummyExecer{Conn: mem.New()}
	compiler := &dummyCompiler{}

	err := migrate.Migrate(execer, compiler)
	if err != nil {
		t.Fatalf("Migrate falló: %v", err)
	}

	if len(execer.queries) == 0 {
		t.Fatalf("se esperaban consultas a ejecutar, se obtuvieron 0")
	}
}
