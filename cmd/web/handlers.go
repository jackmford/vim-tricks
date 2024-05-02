package main

import (
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	files := []string{
		"./ui/html/pages/error.tmpl",
		"./ui/html/pages/home.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Print(err.Error())
		http.Error(w, "Internal Server error", 500)
		return
	}

	trick, err := app.tricks.Get()
	if err != nil {
		app.errorLog.Print(err.Error())
		err = ts.ExecuteTemplate(w, "error", nil)
		if err != nil {
			app.errorLog.Print(err.Error())
		}
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(w, "home", trick)
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
