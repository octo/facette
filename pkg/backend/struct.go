package backend

import (
	"reflect"
)

type tableStruct struct {
	name    string
	columns []*columnStruct
}

func (ts *tableStruct) columnByName(name string) (*columnStruct, error) {
	for _, c := range ts.columns {
		if c.name == name {
			return c, nil
		}
	}

	return nil, ErrUnknownColumn
}

func (ts *tableStruct) mapColumns(t reflect.Type) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	n := t.NumField()
	for i := 0; i < n; i++ {
		f := t.Field(i)

		// Check for embedded struct
		if f.Anonymous {
			ts.mapColumns(f.Type)
			continue
		}

		// Skip fields without `db' tag
		name := f.Tag.Get("db")
		if name == "" {
			continue
		}

		ts.columns = append(ts.columns, &columnStruct{
			name:      name,
			fieldName: f.Name,
		})
	}
}

func newTableStruct(v interface{}) (*tableStruct, error) {
	// Check if struct is valid
	t, ok := v.(tabler)
	if !ok {
		return nil, ErrInvalidStruct
	}

	// Map table struct columns
	ts := &tableStruct{name: t.tableName()}
	ts.mapColumns(reflect.TypeOf(v))

	return ts, nil
}

type columnStruct struct {
	name      string
	fieldName string
}

type tabler interface {
	tableName() string
}
