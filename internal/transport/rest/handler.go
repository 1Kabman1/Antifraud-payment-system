package rest

import (
	"flag"
	"github.com/1Kabman1/antifraud-payment-system/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func StartHandler() error {
	const (
		registerOperation = "/register"
		CreateRule        = "/aggregation_rule/create"
		GetRule           = "/aggregation_rules/get"
		Port              = ":8080"
	)
	addr := flag.String("addr", Port, "HTTP address")
	flag.Parse()
	h := services.NewApiHandler()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post(CreateRule, h.CreateAggregationRule)
	r.Get(GetRule, h.GetAggregationRules)
	r.Post(registerOperation, h.RegisterOperation)
	log.Println("Star server on port " + *addr)
	return http.ListenAndServe(*addr, r)
}
