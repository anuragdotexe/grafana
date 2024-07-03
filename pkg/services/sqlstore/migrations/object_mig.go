package migrations

import "github.com/grafana/grafana/pkg/services/sqlstore/migrator"

// Add SQL table for a simple unified storage object backend
// NOTE: it would be nice to have this defined in the unified storage package, however that
// introduces a circular dependency.
func addObjectMigrations(mg *migrator.Migrator) {
	mg.AddMigration("create unified storage object table", migrator.NewAddTableMigration(migrator.Table{
		Name: "object",
		Columns: []*migrator.Column{
			// Sequential resource version
			{Name: "rv", Type: migrator.DB_BigInt, Nullable: false, IsPrimaryKey: true, IsAutoIncrement: true},

			// Properties that exist in path/key (and duplicated in the json value)
			{Name: "group", Type: migrator.DB_NVarchar, Length: 190, Nullable: false},
			{Name: "namespace", Type: migrator.DB_NVarchar, Length: 63, Nullable: true}, // namespace is not required (cluster scope)
			{Name: "resource", Type: migrator.DB_NVarchar, Length: 190, Nullable: false},
			{Name: "name", Type: migrator.DB_NVarchar, Length: 190, Nullable: false},

			// The k8s resource JSON text (without the resourceVersion populated)
			{Name: "value", Type: migrator.DB_MediumText, Nullable: false},
		},
		Indices: []*migrator.Index{
			{Cols: []string{"group", "namespace", "resource", "name"}, Type: migrator.UniqueIndex},
		},
	}))
}
