package backend

import (
	"database/sql"
	"reflect"
)

type driver interface {
	getBindVar(int) string
	quoteName(string) string
	sqlSchema() []string
	transformValue(reflect.Value, reflect.Value) (reflect.Value, error)
}

func driverSupported(driver string) bool {
	for _, d := range sql.Drivers() {
		if d == driver {
			return true
		}
	}

	return false
}
