package backend

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type Backend struct {
	driver   driver
	db       *sql.DB
	tx       *sql.Tx
	tables   map[reflect.Type]*tableStruct
	idRegexp *regexp.Regexp
}

func (b *Backend) Close() error {
	return b.db.Close()
}

func (b *Backend) Begin() error {
	var err error
	b.tx, err = b.db.Begin()
	return err
}

func (b *Backend) Commit() error {
	err := b.tx.Commit()
	b.tx = nil
	return err
}

func (b *Backend) Rollback() error {
	err := b.tx.Rollback()
	b.tx = nil
	return err
}

func (b *Backend) Get(id string, v interface{}) error {
	// Get table struct
	rv, ts, err := b.parseValue(v, false)
	if err != nil {
		return err
	} else if rv.Kind() != reflect.Ptr {
		return ErrNotAddressable
	}

	// Execute query
	columns := []string{}
	values := []interface{}{}
	for _, c := range ts.columns {
		columns = append(columns, b.driver.quoteName(c.name))
		values = append(values, new(interface{}))
	}

	err = b.queryRow(fmt.Sprintf(
		"SELECT %s FROM %s WHERE id = %s",
		strings.Join(columns, ", "),
		b.driver.quoteName(ts.name),
		b.driver.getBindVar(1),
	), id).Scan(values...)
	if err != nil {
		return err
	}

	// Apply scanned values to struct
	return b.mapValues(rv, ts, values)
}

func (b *Backend) Insert(v interface{}) error {
	// Get table struct
	rv, ts, err := b.parseValue(v, true)
	if err != nil {
		return err
	}

	columns := []string{}
	binds := []string{}
	values := []interface{}{}

	for _, c := range ts.columns {
		f := rv.FieldByName(c.fieldName)

		// Skip field having default value
		if reflect.DeepEqual(f.Interface(), reflect.Zero(f.Type()).Interface()) {
			continue
		}

		// Apply value adaptations
		value, err := b.adaptValue(f)
		if err != nil {
			return err
		}

		columns = append(columns, b.driver.quoteName(c.name))
		binds = append(binds, b.driver.getBindVar(len(columns)))
		values = append(values, value)
	}

	_, err = b.exec(fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		b.driver.quoteName(ts.name),
		strings.Join(columns, ", "),
		strings.Join(binds, ", "),
	), values...)

	return err
}

func (b *Backend) Update(v interface{}) error {
	var id interface{}

	// Get table struct
	rv, ts, err := b.parseValue(v, true)
	if err != nil {
		return err
	}

	sets := []string{}
	values := []interface{}{}

	for _, c := range ts.columns {
		f := rv.FieldByName(c.fieldName)

		if c.name == "id" {
			id = f.Interface()
			continue
		}

		// Skip field having default value
		if reflect.DeepEqual(f.Interface(), reflect.Zero(f.Type()).Interface()) {
			continue
		}

		// Apply value adaptations
		value, err := b.adaptValue(f)
		if err != nil {
			return err
		}

		sets = append(sets, fmt.Sprintf("%s = %s", b.driver.quoteName(c.name), b.driver.getBindVar(len(sets)+1)))
		values = append(values, value)
	}

	if id == nil {
		return ErrInvalidStruct
	}

	values = append(values, id)

	_, err = b.exec(fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = %s",
		b.driver.quoteName(ts.name),
		strings.Join(sets, ", "),
		b.driver.getBindVar(len(sets)+1),
	), values...)

	return err
}

func (b *Backend) Delete(v interface{}) error {
	// Get table struct
	rv, ts, err := b.parseValue(v, true)
	if err != nil {
		return err
	}

	// Get identifier column
	c, err := ts.columnByName("id")
	if err != nil {
		return err
	}

	// Execute query
	_, err = b.exec(fmt.Sprintf(
		"DELETE FROM %s WHERE id = %s",
		b.driver.quoteName(ts.name),
		b.driver.getBindVar(1),
	), rv.FieldByName(c.fieldName).Interface())

	return err
}

func (b *Backend) List(v interface{}, filter string) error {
	// Check for slice element type
	sv := reflect.ValueOf(v)
	if sv.Kind() == reflect.Interface && sv.Elem().Kind() != reflect.Slice {
		return ErrInvalidSlice
	}
	sv = reflect.Indirect(sv).Elem()

	// Create new slice element instance
	rn := reflect.New(sv.Type().Elem())

	// Get table struct
	rv, ts, err := b.parseValue(rn.Interface(), false)
	if err != nil {
		return err
	} else if rv.Kind() != reflect.Ptr {
		return ErrNotAddressable
	}

	// Execute query
	columns := []string{}
	values := []interface{}{}
	for _, c := range ts.columns {
		columns = append(columns, b.driver.quoteName(c.name))
		values = append(values, new(interface{}))
	}

	rows, err := b.query(fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(columns, ", "),
		b.driver.quoteName(ts.name),
	))
	if err != nil {
		return err
	} else if err := rows.Err(); err != nil {
		return err
	}
	defer rows.Close()

	// Scan rows for values
	for rows.Next() {
		rows.Scan(values...)

		// Apply scanned values to struct and append it to result slice
		r := reflect.New(reflect.Indirect(rn).Type())
		if err := b.mapValues(&r, ts, values); err != nil {
			return err
		}
		reflect.Append(sv, reflect.Indirect(r))
	}

	return nil
}

func (b *Backend) parseValue(v interface{}, check bool) (*reflect.Value, *tableStruct, error) {
	// Check if value is addressable
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return nil, nil, ErrInvalidStruct
	}

	rt := rv.Type()

	// Check for registered table struct or map a new one
	if _, ok := b.tables[rt]; !ok {
		var err error

		b.tables[rt], err = newTableStruct(v)
		if err != nil {
			return nil, nil, err
		}
	}

	// Check identifier validity
	if check {
		c, err := b.tables[rt].columnByName("id")
		if err != nil {
			return nil, nil, err
		}

		f := reflect.Indirect(rv).FieldByName(c.fieldName).Interface()
		if id, ok := f.(string); !ok || !b.idRegexp.MatchString(id) {
			return nil, nil, ErrInvalidIdentifier
		}
	}

	return &rv, b.tables[rt], nil
}

func (b *Backend) mapValues(rv *reflect.Value, ts *tableStruct, values []interface{}) error {
	for idx, c := range ts.columns {
		value := reflect.ValueOf(values[idx]).Elem()
		if value.Kind() == reflect.Interface {
			value = value.Elem()
		}

		if !value.IsValid() {
			continue
		}

		f := rv.Elem().FieldByName(c.fieldName)

		// Handle reflect conversion for complex types
		if !value.Type().ConvertibleTo(f.Type()) {
			var err error

			// Transform value according to driver rules if any
			value, err = b.driver.transformValue(f, value)
			if err != nil && err != ErrNotTransformable {
				return fmt.Errorf("failed to transform value: %s", err)
			} else if err == nil {
				goto set
			}

			switch f.Kind() {
			case reflect.Map, reflect.Slice, reflect.Struct:
				r := reflect.New(f.Type())

				// Unmarshal non-stringable value JSON encoded data
				err := json.Unmarshal(value.Interface().([]byte), r.Interface())
				if err != nil {
					return fmt.Errorf("failed to unmarshal data: %s", err)
				}

				value = r.Elem()
			}
		}

	set:
		f.Set(value.Convert(f.Type()))
	}

	return nil
}

func (b *Backend) adaptValue(rv reflect.Value) (interface{}, error) {
	// Marshal non-stringable value to JSON encoded data
	_, stringable := rv.Type().MethodByName("String")
	if !stringable {
		switch rv.Kind() {
		case reflect.Map, reflect.Slice, reflect.Struct:
			data, err := json.Marshal(rv.Interface())
			if err != nil {
				return nil, err
			}

			rv = reflect.ValueOf(data)
		}
	}

	return rv.Interface(), nil
}

func (b *Backend) exec(query string, values ...interface{}) (sql.Result, error) {
	if b.tx != nil {
		return b.tx.Exec(query, values...)
	}

	return b.db.Exec(query, values...)
}

func (b *Backend) query(query string, values ...interface{}) (*sql.Rows, error) {
	if b.tx != nil {
		return b.tx.Query(query, values...)
	}

	return b.db.Query(query, values...)
}

func (b *Backend) queryRow(query string, values ...interface{}) *sql.Row {
	if b.tx != nil {
		return b.tx.QueryRow(query, values...)
	}

	return b.db.QueryRow(query, values...)
}

func NewBackend(driver, dsn string) (*Backend, error) {
	// Check if database driver is supported
	if !driverSupported(driver) {
		return nil, fmt.Errorf("unsupported database driver `%s'", driver)
	}

	// Open database connection
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	// Create new backend
	b := &Backend{
		db:       db,
		tables:   make(map[reflect.Type]*tableStruct),
		idRegexp: regexp.MustCompile("^\\d{8}-(?:\\d{4}-){3}\\d{12}$"),
	}

	switch driver {
	case "mysql":
		b.driver = mysqlDriver{}
	case "postgres":
		b.driver = postgresDriver{}
	case "sqlite3":
		b.driver = sqlite3Driver{}
	default:
		return nil, fmt.Errorf("unsupported database backend `%s'", driver)
	}

	// Initialize database schema
	for _, q := range b.driver.sqlSchema() {
		if _, err := b.db.Exec(q); err != nil {
			return nil, err
		}
	}

	return b, nil
}
