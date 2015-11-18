package backend

type Unit struct {
	Item
	Label string `db:"label" json:"label"`
}

func (u Unit) tableName() string {
	return "units"
}
