package backend

type SourceGroup struct {
	Item
	Entries []GroupEntry `db:"entries" json:"entries"`
}

func (s SourceGroup) tableName() string {
	return "sourcegroups"
}

type MetricGroup struct {
	Item
	Entries []GroupEntry `db:"entries" json:"entries"`
}

func (s MetricGroup) tableName() string {
	return "metricgroups"
}

type GroupEntry struct {
	Pattern string `json:"pattern"`
	Origin  string `json:"origin"`
}
