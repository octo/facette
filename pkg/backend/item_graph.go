package backend

type Graph struct {
	Item
	Series     []SeriesGroup          `db:"series" json:"series,omitempty"`
	Link       string                 `db:"link" json:"link,omitempty"`
	Attributes map[string]interface{} `db:"attributes" json:"attributes,omitempty"`
	Options    map[string]interface{} `db:"options" json:"options,omitempty"`
	Template   bool                   `db:"template" json:"template"`
}

func (g Graph) tableName() string {
	return "graphs"
}

type SeriesGroup struct {
	Series   []Series `json:"series"`
	Operator uint8    `json:"operator"`
}

type Series struct {
	Name   string `json:"name"`
	Origin string `json:"origin"`
	Source string `json:"source"`
	Metric string `json:"metric"`
}
