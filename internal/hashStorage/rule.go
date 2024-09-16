package hashStorage

type Rule struct {
	AggregationRuleId int
	Name              string   `json:"Name"`
	AggregateBy       []string `json:"AggregateBy"`
	AggregateValue    string   `json:"AggregateValue"`
	ExpirationTime    int      `json:"ExpirationTime"`
	TimePeriod        int      `json:"TimePeriod"`
}

func NewRule() Rule {
	return Rule{}
}
