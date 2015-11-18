package backend

import (
	"reflect"
	"testing"
	"time"
)

var (
	date time.Time
)

func init() {
	date = time.Now().Round(time.Second)
}

func Test_Sqlite3_Graph(t *testing.T) {
	execTestGraph("sqlite3", "/tmp/facette-backend.db", t)
}

func Test_Sqlite3_SourceGroup(t *testing.T) {
	execTestSourceGroup("sqlite3", "/tmp/facette-backend.db", t)
}

func Test_Sqlite3_MetricGroup(t *testing.T) {
	execTestMetricGroup("sqlite3", "/tmp/facette-backend.db", t)
}

func Test_Sqlite3_Scale(t *testing.T) {
	execTestScale("sqlite3", "/tmp/facette-backend.db", t)
}

func Test_Sqlite3_Unit(t *testing.T) {
	execTestUnit("sqlite3", "/tmp/facette-backend.db", t)
}

func Test_Postgres_Graph(t *testing.T) {
	execTestGraph("postgres", "dbname=facette user=facette password=facette", t)
}

func Test_Postgres_SourceGroup(t *testing.T) {
	execTestSourceGroup("postgres", "dbname=facette user=facette password=facette", t)
}

func Test_Postgres_MetricGroup(t *testing.T) {
	execTestMetricGroup("postgres", "dbname=facette user=facette password=facette", t)
}

func Test_Postgres_Scale(t *testing.T) {
	execTestScale("postgres", "dbname=facette user=facette password=facette", t)
}

func Test_Postgres_Unit(t *testing.T) {
	execTestUnit("postgres", "dbname=facette user=facette password=facette", t)
}

func Test_MySQL_Graph(t *testing.T) {
	execTestGraph("mysql", "facette:facette@/facette?interpolateParams=true", t)
}

func Test_MySQL_SourceGroup(t *testing.T) {
	execTestSourceGroup("mysql", "facette:facette@/facette?interpolateParams=true", t)
}

func Test_MySQL_MetricGroup(t *testing.T) {
	execTestMetricGroup("mysql", "facette:facette@/facette?interpolateParams=true", t)
}

func Test_MySQL_Scale(t *testing.T) {
	execTestScale("mysql", "facette:facette@/facette?interpolateParams=true", t)
}

func Test_MySQL_Unit(t *testing.T) {
	execTestUnit("mysql", "facette:facette@/facette?interpolateParams=true", t)
}

func execTestGraph(driver, dsn string, t *testing.T) {
	b, err := NewBackend(driver, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()

	execTest(b, Graph{
		Item: Item{
			ID:          "00000000-0000-0000-0000-000000000000",
			Name:        "graph1",
			Description: "A great graph description",
			Created:     date,
			Modified:    date,
		},
		Series: []SeriesGroup{
			SeriesGroup{
				Series: []Series{
					Series{
						Name:   "series1",
						Origin: "origin1",
						Source: "source1",
						Metric: "metric1",
					},
				},
			},
		},
		Options: map[string]interface{}{
			"title": "A great graph title",
		},
		Template: false,
	}, Graph{
		Item: Item{
			ID:          "00000000-0000-0000-0000-000000000000",
			Name:        "graph1",
			Description: "A great graph description (updated)",
			Created:     date,
			Modified:    date.Add(time.Hour),
		},
		Series: []SeriesGroup{
			SeriesGroup{
				Series: []Series{
					Series{
						Name:   "series1",
						Origin: "origin1",
						Source: "source1",
						Metric: "metric1",
					},
					Series{
						Name:   "series2",
						Origin: "origin1",
						Source: "source1",
						Metric: "metric2",
					},
				},
			},
		},
		Options: map[string]interface{}{
			"title": "A great graph title (updated)",
		},
		Template: false,
	}, t)
}

func execTestSourceGroup(driver, dsn string, t *testing.T) {
	b, err := NewBackend(driver, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()

	execTest(b, SourceGroup{
		Item: Item{
			ID:          "00000000-0000-0000-0000-000000000000",
			Name:        "sourcegroup1",
			Description: "A great sourcegroup description",
			Created:     date,
			Modified:    date,
		},
		Entries: []GroupEntry{
			GroupEntry{
				Pattern: "glob:host*.example.net",
				Origin:  "origin1",
			},
		},
	}, SourceGroup{
		Item: Item{
			ID:          "00000000-0000-0000-0000-000000000000",
			Name:        "sourcegroup1",
			Description: "A great sourcegroup description (updated)",
			Created:     date,
			Modified:    date.Add(time.Hour),
		},
		Entries: []GroupEntry{
			GroupEntry{
				Pattern: "glob:host*.example.net",
				Origin:  "",
			},
		},
	}, t)
}

func execTestMetricGroup(driver, dsn string, t *testing.T) {
	b, err := NewBackend(driver, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()

	execTest(b, MetricGroup{
		Item: Item{
			ID:          "00000000-0000-0000-0000-000000000000",
			Name:        "metricgroup1",
			Description: "A great metricgroup description",
			Created:     date,
			Modified:    date,
		},
		Entries: []GroupEntry{
			GroupEntry{
				Pattern: "glob:load.*",
				Origin:  "origin1",
			},
		},
	}, MetricGroup{
		Item: Item{
			ID:          "00000000-0000-0000-0000-000000000000",
			Name:        "metricgroup1",
			Description: "A great metricgroup description (updated)",
			Created:     date,
			Modified:    date.Add(time.Hour),
		},
		Entries: []GroupEntry{
			GroupEntry{
				Pattern: "glob:load.*",
				Origin:  "",
			},
		},
	}, t)
}

func execTestScale(driver, dsn string, t *testing.T) {
	b, err := NewBackend(driver, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()

	execTest(b, Scale{
		Item: Item{
			ID:          "00000000-0000-0000-0000-000000000000",
			Name:        "scale1",
			Description: "A great scale description",
			Created:     date,
			Modified:    date,
		},
		Value: 0.123,
	}, Scale{
		Item: Item{
			ID:          "00000000-0000-0000-0000-000000000000",
			Name:        "scale1",
			Description: "A great scale description (updated)",
			Created:     date,
			Modified:    date.Add(time.Hour),
		},
		Value: 0.456,
	}, t)
}

func execTestUnit(driver, dsn string, t *testing.T) {
	b, err := NewBackend(driver, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()

	execTest(b, Unit{
		Item: Item{
			ID:          "00000000-0000-0000-0000-000000000000",
			Name:        "unit1",
			Description: "A great unit description",
			Created:     date,
			Modified:    date,
		},
		Label: "ms",
	}, Unit{
		Item: Item{
			ID:          "00000000-0000-0000-0000-000000000000",
			Name:        "unit1",
			Description: "A great unit description (updated)",
			Created:     date,
			Modified:    date.Add(time.Hour),
		},
		Label: "s",
	}, t)
}

func execTest(b *Backend, item, update interface{}, t *testing.T) {
	var result interface{}

	// Check item insertion
	if err := b.Insert(item); err != nil {
		t.Fatal(err)
	}

	result = reflect.New(reflect.TypeOf(item)).Interface()
	if err := b.Get("00000000-0000-0000-0000-000000000000", result); err != nil {
		t.Fatal(err)
	}

	result = reflect.Indirect(reflect.ValueOf(result)).Interface()
	if !deepEqual(item, result) {
		t.Logf("\nExpected %#v\nbut got  %#v", item, result)
		t.Fail()
	}

	// Insert second element
	// TODO

	// Check items list
	v := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(item)), 0, 0).Interface()
	if err := b.List(&v, ""); err != nil {
		t.Fatal(err)
	}

	// TODO

	// Check item update
	if err := b.Update(update); err != nil {
		t.Fatal(err)
	}

	result = reflect.New(reflect.TypeOf(item)).Interface()
	if err := b.Get("00000000-0000-0000-0000-000000000000", result); err != nil {
		t.Fatal(err)
	}

	result = reflect.Indirect(reflect.ValueOf(result)).Interface()
	if !deepEqual(update, result) {
		t.Logf("\nExpected %#v\nbut got  %#v", update, result)
		t.Fail()
	}

	// Check item deletion
	if err := b.Delete(item); err != nil {
		t.Fatal(err)
	}
}

func deepEqual(a, b interface{}) bool {
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	if va.Kind() == reflect.Struct {
		if va.NumField() != vb.NumField() {
			return false
		}

		for i, n := 0, va.NumField(); i < n; i++ {
			ia := va.Field(i).Interface()
			ib := vb.Field(i).Interface()

			if ta, ok := ia.(time.Time); ok {
				if tb, ok := ib.(time.Time); !ok {
					return false
				} else if !ta.Equal(tb) {
					return false
				}
			}

			return deepEqual(ia, ib)
		}
	}

	return reflect.DeepEqual(a, b)
}
