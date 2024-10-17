package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var db *sql.DB
var tracer trace.Tracer

func initDB() {
	var err error
	connStr := "postgres://admin:password@postgres/mydb?sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func initTracer() func() {
	// Create OTLP trace exporter (sending traces to Tempo)
	exporter, err := otlptracehttp.New(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Create trace provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)

	// Return cleanup function to ensure all data is flushed
	return func() {
		err := tp.Shutdown(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}
}

func createTransaction(w http.ResponseWriter, r *http.Request) {
	var txn Transaction
	err := json.NewDecoder(r.Body).Decode(&txn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, span := tracer.Start(r.Context(), "CreateTransaction")
	defer span.End()

	span.SetAttributes(
		attribute.String("event_id", txn.EventID),
		attribute.String("user", txn.User),
	)

	query := `INSERT INTO transactions (event_id, timestamp, user, item_id) VALUES ($1, $2, $3, $4)`
	_, err = db.ExecContext(ctx, query, txn.EventID, txn.Timestamp, txn.User, txn.ItemID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func main() {
	// Initialize the database
	initDB()
	defer db.Close()

	// Initialize OpenTelemetry tracer
	shutdown := initTracer()
	defer shutdown()

	tracer = otel.Tracer("golang-api")

	http.HandleFunc("/transaction", createTransaction)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
