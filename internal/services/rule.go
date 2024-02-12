package services

type Rule struct {
	AggregationRuleId int      `json:"id"`
	Name              string   `json:"Name"`
	AggregateBy       []string `json:"AggregateBy"`
	AggregateValue    string   `json:"AggregateValue"`
}
