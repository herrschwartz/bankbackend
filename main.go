package main

import (
	"context"
	"html/template"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4/log/log15adapter"
	"github.com/jackc/pgx/v4/pgxpool"
	log "gopkg.in/inconshreveable/log15.v2"
)

var db *pgxpool.Pool
var tpl *template.Template

//Item for things
type Item struct {
	Id  int
	Amt int
}

type indexData struct {
	Guild string
	Items []Item
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		testdata := indexData{
			Guild: "Its ok to be Whitemane",
			Items: []Item{
				{Id: 5113, Amt: 1},
				{Id: 16204, Amt: 20},
				{Id: 16393, Amt: 1},
				{Id: 19019, Amt: 1},
				{Id: 14551, Amt: 1},
				{Id: 14552, Amt: 1},
				{Id: 16223, Amt: 1},
			},
		}
		tpl.ExecuteTemplate(w, "index.html", testdata)
	}
}

func main() {
	logger := log15adapter.NewLogger(log.New("module", "pgx"))

	poolConfig, err := pgxpool.ParseConfig("postgres://postgres:admin@127.0.0.1:5432/postgres?pool_max_conns=10")
	if err != nil {
		log.Crit("Unable to parse DATABASE_URL", "error", err)
		os.Exit(1)
	}

	poolConfig.ConnConfig.Logger = logger

	db, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Crit("Unable to create connection pool", "error", err)
		os.Exit(1)
	}

	tpl = template.Must(template.ParseFiles("templates/index.html"))

	http.HandleFunc("/", indexHandler)

	log.Info("Starting App on localhost:8080")
	err = http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Crit("Unable to start web server", "error", err)
		os.Exit(1)
	}
}
