package services

import (
	"bytes"
	"encoding/json"
	"github.com/1Kabman1/antifraud-payment-system/internal/hashStorage"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	amount = "amount"
)

type ApiHandler struct {
	s        hashStorage.Storage
	errorLog *log.Logger
	infoLog  *log.Logger
}

// NewApiHandler - set storage and logs
func NewApiHandler() ApiHandler {
	return ApiHandler{
		s:        hashStorage.NewStorage(),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
	}
}

// GetAggregationRules - Get aggregation data
func (h *ApiHandler) GetAggregationRules(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]hashStorage.Rule, h.s.RulesLen())
	for _, aRule := range h.s.Rules() {
		resp[strconv.Itoa(aRule.AggregationRuleId)] = *aRule
	}
	ruleJson, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.errorLog.Println(err)
		return
	}
	if json.Valid(ruleJson) {
		if _, err = w.Write([]byte("Status " + "success \n")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			h.errorLog.Println(err)
			return
		}
		if _, err = w.Write(ruleJson); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			h.errorLog.Println(err)
			return
		}
	}
}

// CreateAggregationRule - create aggregation Rule
func (h *ApiHandler) CreateAggregationRule(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	aRule := hashStorage.NewRule()
	if r.Body != nil {
		defer func() {
			if err := r.Body.Close(); err != nil {
				h.errorLog.Println(err)
			}
		}()
	}
	if _, err := buf.ReadFrom(r.Body); err != nil {
		h.errorLog.Println(err)
		return
	}
	if err := json.Unmarshal(buf.Bytes(), &aRule); err != nil {
		h.errorLog.Println(err)
		return
	}
	h.s.SetRule(aRule.Name, &aRule)
	if _, err := w.Write([]byte("Message " + "Rule " + strconv.Itoa(aRule.AggregationRuleId) + " created")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.errorLog.Println(err)
		return
	}
}

// RegisterOperation - counts aggregated based on the rules
func (h *ApiHandler) RegisterOperation(w http.ResponseWriter, r *http.Request) {
	payment := map[string]interface{}{}
	if r.Body != nil {
		defer func() {
			if err := r.Body.Close(); err != nil {
				h.errorLog.Println(err)
			}
		}()
	}
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.errorLog.Println(err.Error())
		return
	}
	aggregatesBy, aErr := prepareTheDataForHashing(h.s.Rules(), payment)
	if aErr != nil {
		http.Error(w, aErr.Error(), http.StatusInternalServerError)
		h.errorLog.Println(aErr)
		return
	}
	for nameRule, keyCounter := range calculateHash(aggregatesBy) {
		hteErr, aRule := h.s.Rule(nameRule)
		if hteErr != nil {
			http.Error(w, hteErr.Error(), http.StatusInternalServerError)
			h.errorLog.Println(hteErr)
			return
		}
		h.s.SetCounter(keyCounter, aRule.AggregationRuleId)
		h.s.FixingThePayment(keyCounter, aRule.AggregateValue, payment[amount].(float64))
	}
}
