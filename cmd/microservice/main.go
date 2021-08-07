package main

import (
	"net/http"
	"os"

	"github.com/luigizuccarelli/iotpaas-message-producer/pkg/connectors"
	"github.com/luigizuccarelli/iotpaas-message-producer/pkg/handlers"
	"github.com/luigizuccarelli/iotpaas-message-producer/pkg/validator"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gorilla/mux"
	"github.com/microlib/simple"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
)

var (
	logger       *simple.Logger
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "iotpaas_message_producer_http_duration_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path"})
)

// prometheusMiddleware implements mux.MiddlewareFunc.
func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(CONTENTTYPE, APPLICATIONJSON)
		// use this for cors
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept-Language")
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		next.ServeHTTP(w, r)
		timer.ObserveDuration()
	})
}

// startHttpServer a utility function that sets the routes, handlers and starts the http server
func startHttpServer(conn connectors.Clients) *http.Server {

	// set the server props
	srv := &http.Server{Addr: ":" + os.Getenv("SERVER_PORT")}

	// set the router and endpoints
	r := mux.NewRouter()

	r.Use(prometheusMiddleware)
	r.Path("/metrics").Handler(promhttp.Handler())

	r.HandleFunc("/api/v1/streamdata", func(w http.ResponseWriter, req *http.Request) {
		handlers.StreamHandler(w, req, conn)
	}).Methods("POST", "OPTIONS")

	r.HandleFunc("/api/v2/sys/info/isalive", handlers.IsAlive).Methods("GET", "OPTIONS")

	http.Handle("/", r)

	// start our server (concurrent)
	if err := srv.ListenAndServe(); err != nil {
		conn.Error("Httpserver: ListenAndServe() error: %v", err)
		os.Exit(0)
	}

	// return our srv object
	return srv
}

// The main function reads the config file, parsers and validates it and calls our start server function
// using a go channel to intercept sig calls from the os
// A simple curl to test the payload endpoint
// curl -H "Accept: application/json"  -H "Content-Type: application/json" -X PUT -d @sparkpost-webhook-payload.json http://sparkpost-spring-producer-microservice-sparkpost-poc.apps.poc.okd.14west.io/webhook
func main() {

	var logger *simple.Logger

	if os.Getenv("LOG_LEVEL") == "" {
		logger = &simple.Logger{Level: "info"}
	} else {
		logger = &simple.Logger{Level: os.Getenv("LOG_LEVEL")}
	}
	err := validator.ValidateEnvars(logger)
	if err != nil {
		os.Exit(-1)
	}

	// setup our client connectors (message producer)
	conn := connectors.NewClientConnectors(logger)

	// call the start server function
	logger.Info("Starting server on port " + os.Getenv("SERVER_PORT"))
	startHttpServer(conn)
}
