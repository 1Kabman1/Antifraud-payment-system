package services

import (
	"crypto/md5"
	"errors"
	"github.com/1Kabman1/antifraud-payment-system/internal/hashStorage"
	"strconv"
	"strings"
)

// calculateHash - hash function
func calculateHash(data map[string]string) map[string][16]byte {
	result := make(map[string][16]byte)

	for key, val := range data {
		h := md5.Sum([]byte(val))
		result[key] = h
	}

	return result

}

// prepareTheDataForHashing - prepares data for hashing
func prepareTheDataForHashing(rules, operationProperties map[string]interface{}) (map[string]string, error) {
	aggregatesBy := make(map[string]string, len(rules))
	var aBuilder strings.Builder

	for _, tempRule := range rules {

		aRule := tempRule.(rule)
		aggregate := ""
		var flag bool

		for _, agg := range aRule.AggregateBy {
			v, ok := operationProperties[agg]
			if !ok {
				flag = true
				aBuilder.Reset()
				break
			}
			switch aInterface := v.(type) {
			case float64:
				intInterface := int(aInterface)
				_, err := aBuilder.WriteString(strconv.Itoa(intInterface))
				if err != nil {
					return nil, errors.New("the type is not float")
				}
				aBuilder.WriteString(agg) // добавляю имя агрегирующего чтобы убрать совпадение
			case string:
				_, err := aBuilder.WriteString(aInterface)
				if err != nil {
					return nil, errors.New("the type is not string")
				}
				aBuilder.WriteString(agg) // добавляю имя агрегирующего чтобы убрать совпадение
			}

		}

		if !flag {
			aggregate = aBuilder.String()
			aggregatesBy[aRule.Name] = aggregate
			aBuilder.Reset()
		}

	}
	return aggregatesBy, nil
}

func increaseTheCounter(keyCounter [16]byte, h *hashStorage.Storage, c counter,
	AggregateValue string, mapping map[string]interface{}) {

	if AggregateValue == count {
		c.Value += 1
		h.SetCounter(keyCounter, c)
	} else {
		c.Value += int(mapping[amount].(float64))
		h.SetCounter(keyCounter, c)
	}
}

func addTheCounter(keyCounter [16]byte, AggregationRuleId int, h *hashStorage.Storage, aNewCounter counter,
	AggregateValue string, mapping map[string]interface{}) {

	if AggregateValue == amount {
		aNewCounter.Value = int(mapping[amount].(float64))
	} else {
		aNewCounter.Value++
	}
	idCounter := h.CounterLen()
	aNewCounter.id = idCounter
	h.SetCounter(keyCounter, aNewCounter)
	h.AddToArchivist(AggregationRuleId, idCounter)

}
