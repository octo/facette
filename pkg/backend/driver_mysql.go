package backend

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type mysqlDriver struct{}

func (d mysqlDriver) getBindVar(i int) string {
	return "?"
}

func (d mysqlDriver) quoteName(s string) string {
	return fmt.Sprintf("`%s`", s)
}

func (d mysqlDriver) sqlSchema() []string {
	return []string{
		`CREATE TABLE IF NOT EXISTS scales (
			id VARCHAR(36) NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			value NUMERIC NOT NULL,
			CONSTRAINT ` + "`pk_scales`" + ` PRIMARY KEY (id),
			CONSTRAINT ` + "`un_scales_name`" + `UNIQUE (name),
			CONSTRAINT ` + "`un_scales_value`" + `UNIQUE (value)
		)`,
		`CREATE TABLE IF NOT EXISTS units (
			id VARCHAR(36) NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			label VARCHAR(32) NOT NULL,
			CONSTRAINT ` + "`pk_units`" + `PRIMARY KEY (id),
			CONSTRAINT ` + "`un_units_name`" + `UNIQUE (name),
			CONSTRAINT ` + "`un_units_value`" + `UNIQUE (label)
		)`,
		`CREATE TABLE IF NOT EXISTS sourcegroups (
			id VARCHAR(36) NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			entries TEXT NOT NULL,
			CONSTRAINT ` + "`pk_sourcegroups`" + `PRIMARY KEY (id),
			CONSTRAINT ` + "`un_sourcegroups_name`" + `UNIQUE (name)
		)`,
		`CREATE TABLE IF NOT EXISTS metricgroups (
			id VARCHAR(36) NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			entries TEXT NOT NULL,
			CONSTRAINT ` + "`pk_metricgroups`" + `PRIMARY KEY (id),
			CONSTRAINT ` + "`un_metricgroups_name`" + `UNIQUE (name)
		)`,
		`CREATE TABLE IF NOT EXISTS graphs (
			id VARCHAR(36) NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			series TEXT,
			link VARCHAR(36),
			attributes TEXT,
			options TEXT,
			template BOOLEAN NOT NULL DEFAULT false,
			CONSTRAINT ` + "`pk_graphs`" + `PRIMARY KEY (id),
			CONSTRAINT ` + "`fk_graphs_link`" + `FOREIGN KEY (link) REFERENCES graphs (id)
				ON DELETE CASCADE ON UPDATE CASCADE,
			CONSTRAINT ` + "`un_graphs_name`" + `UNIQUE (name),
			CONSTRAINT ` + "`ck_graphs_entry`" + `CHECK ((series IS NOT NULL AND link IS NULL AND attributes IS NULL) OR
				(series IS NULL AND template = false AND link IS NOT NULL AND attributes IS NOT NULL))
		)`,
		`CREATE TABLE IF NOT EXISTS collections (
			id VARCHAR(36) NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT ` + "`pk_collections`" + `PRIMARY KEY (id),
			CONSTRAINT ` + "`un_collections_name`" + `UNIQUE (name)
		)`,
		`CREATE TABLE IF NOT EXISTS collections_graphs (
			collection_id VARCHAR(36) NOT NULL,
			graph_id VARCHAR(36) NOT NULL,
			options TEXT,
			CONSTRAINT ` + "`pk_collections_graphs`" + `PRIMARY KEY (collection_id, graph_id),
			CONSTRAINT ` + "`fk_collections_graphs_collection_id`" + `FOREIGN KEY (collection_id) REFERENCES collections (id)
				ON DELETE CASCADE ON UPDATE CASCADE,
			CONSTRAINT ` + "`fk_collections_graphs_graph_id`" + `FOREIGN KEY (graph_id) REFERENCES graphs (id)
				ON DELETE CASCADE ON UPDATE CASCADE
		)`,
	}
}

func (d mysqlDriver) transformValue(rv, value reflect.Value) (reflect.Value, error) {
	switch rv.Kind() {
	case reflect.Bool:
		v, err := strconv.ParseBool(string(value.Interface().([]byte)))
		if err != nil {
			return reflect.ValueOf(nil), err
		}

		return reflect.ValueOf(v), nil

	case reflect.Float64:
		float, err := strconv.ParseFloat(string(value.Interface().([]byte)), 64)
		if err != nil {
			return reflect.ValueOf(nil), err
		}

		return reflect.ValueOf(float), nil

	case reflect.Struct:
		if _, ok := rv.Interface().(time.Time); ok {
			t, err := time.Parse("2006-01-02 15:04:05", string(value.Interface().([]byte)))
			if err != nil {
				return reflect.ValueOf(nil), err
			}

			return reflect.ValueOf(t), nil
		}
	}

	return value, ErrNotTransformable
}
