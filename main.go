package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"vimtricks.jackmitchellfordyce.com/internal/models"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	tricks   *models.TrickModel
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	trick, err := app.tricks.Get()
	if err != nil {
		app.errorLog.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	ts, err := template.ParseFiles("./ui/html/pages/home.tmpl")
	if err != nil {
		app.errorLog.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

	err = ts.Execute(w, trick)
	if err != nil {
		app.errorLog.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func (app *application) trickCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		//app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "test4"
	content := "testing4 vimtricks"

	id, err := app.tricks.Insert(title, content)
	if err != nil {
		app.errorLog.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	app.infoLog.Print(id)
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")

	dsn := flag.String("dsn", "web:pass@/vimtricks?parseTime=true", "MySQL data source name")

	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		tricks:   &models.TrickModel{DB: db},
	}

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/create", app.trickCreate)

	infoLog.Print("Starting server on :4000")
	err = http.ListenAndServe(*addr, mux)
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
