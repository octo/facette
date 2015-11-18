package backend

type Scale struct {
	Item
	Value float64 `db:"value" json:"value"`
}

func (s Scale) tableName() string {
	return "scales"
}
