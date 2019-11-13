package main

import (
	"context"
	"html/template"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4"
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
		rows, err := db.Query(context.Background(), "SELECT item_id, quantity FROM items WHERE guild_id=$1", 1)
		if err != nil {
			log.Error("Unable to complete query", "error", err)
		}
		defer rows.Close()

		var items []Item
		for rows.Next() {
			var itemID int
			var amount int
			err = rows.Scan(&itemID, &amount)
			if err != nil {
				log.Error("unable to scan row", "error", err)
			}
			items = append(items, Item{Id: itemID, Amt: amount})
		}

		testdata := indexData{
			Guild: "Its ok to be Whitemane",
			Items: items,
		}
		tpl.ExecuteTemplate(w, "index.html", testdata)
	}
}

func updateItems(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		rows := [][]interface{}{}

		rows = append(rows, []interface{}{1, 2, 4})
		copyCount, err := db.CopyFrom(context.Background(), pgx.Identifier{"items"}, []string{"item_id, guild_id, quantity"}, pgx.CopyFromRows(rows))
		if err != nil {
			log.Error("Unexpected error for CopyFrom: %v", err)
		}
		if int(copyCount) != len(rows) {
			log.Error("Expected CopyFrom to return %d copied rows, but got %d", len(rows), copyCount)
		}

	} else {
		log.Warn("Wrong request method at /update", req.Method)
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
	http.HandleFunc("/update", updateItems)

	log.Info("Starting App on localhost:8080")
	err = http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Crit("Unable to start web server", "error", err)
		os.Exit(1)
	}
}
