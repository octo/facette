package backend

import (
	"fmt"
	"reflect"
	"strconv"

	_ "github.com/lib/pq"
)

type postgresDriver struct{}

func (d postgresDriver) getBindVar(i int) string {
	return fmt.Sprintf("$%d", i)
}

func (d postgresDriver) quoteName(s string) string {
	return fmt.Sprintf("%q", s)
}

func (d postgresDriver) sqlSchema() []string {
	return []string{
		`CREATE TABLE IF NOT EXISTS scales (
			id UUID NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
			modified TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
			value NUMERIC NOT NULL,
			CONSTRAINT "pk_scales" PRIMARY KEY (id),
			CONSTRAINT "un_scales_name" UNIQUE (name),
			CONSTRAINT "un_scales_value" UNIQUE (value)
		)`,
		`CREATE TABLE IF NOT EXISTS units (
			id UUID NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
			modified TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
			label VARCHAR(32) NOT NULL,
			CONSTRAINT "pk_units" PRIMARY KEY (id),
			CONSTRAINT "un_units_name" UNIQUE (name),
			CONSTRAINT "un_units_value" UNIQUE (label)
		)`,
		`CREATE TABLE IF NOT EXISTS sourcegroups (
			id UUID NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
			modified TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
			entries TEXT NOT NULL,
			CONSTRAINT "pk_sourcegroups" PRIMARY KEY (id),
			CONSTRAINT "un_sourcegroups_name" UNIQUE (name)
		)`,
		`CREATE TABLE IF NOT EXISTS metricgroups (
			id UUID NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
			modified TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
			entries TEXT NOT NULL,
			CONSTRAINT "pk_metricgroups" PRIMARY KEY (id),
			CONSTRAINT "un_metricgroups_name" UNIQUE (name)
		)`,
		`CREATE TABLE IF NOT EXISTS graphs (
			id UUID NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
			modified TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
			series TEXT,
			link UUID,
			attributes TEXT,
			options TEXT,
			template BOOLEAN NOT NULL DEFAULT false,
			CONSTRAINT "pk_graphs" PRIMARY KEY (id),
			CONSTRAINT "fk_graphs_link" FOREIGN KEY (link) REFERENCES graphs (id)
				ON DELETE CASCADE ON UPDATE CASCADE,
			CONSTRAINT "un_graphs_name" UNIQUE (name),
			CONSTRAINT "ck_graphs_entry" CHECK ((series IS NOT NULL AND link IS NULL AND attributes IS NULL) OR
				(series IS NULL AND template = false AND link IS NOT NULL AND attributes IS NOT NULL))
		)`,
		`CREATE TABLE IF NOT EXISTS collections (
			id UUID NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
			modified TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
			CONSTRAINT "pk_collections" PRIMARY KEY (id),
			CONSTRAINT "un_collections_name" UNIQUE (name)
		)`,
		`CREATE TABLE IF NOT EXISTS collections_graphs (
			collection_id UUID NOT NULL,
			graph_id UUID NOT NULL,
			options TEXT,
			CONSTRAINT "pk_collections_graphs" PRIMARY KEY (collection_id, graph_id),
			CONSTRAINT "fk_collections_graphs_collection_id" FOREIGN KEY (collection_id) REFERENCES collections (id)
				ON DELETE CASCADE ON UPDATE CASCADE,
			CONSTRAINT "fk_collections_graphs_graph_id" FOREIGN KEY (graph_id) REFERENCES graphs (id)
				ON DELETE CASCADE ON UPDATE CASCADE
		)`,
	}
}

func (d postgresDriver) transformValue(rv, value reflect.Value) (reflect.Value, error) {
	switch rv.Kind() {
	case reflect.Float64:
		float, err := strconv.ParseFloat(string(value.Interface().([]byte)), 64)
		if err != nil {
			return reflect.ValueOf(nil), err
		}

		return reflect.ValueOf(float), nil
	}

	return value, ErrNotTransformable
}
