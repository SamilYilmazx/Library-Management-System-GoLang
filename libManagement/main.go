package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Book struct {
	Id     int
	Author string
	Name   string
	Review string
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "tren400"
	dbName := "libraryManagement"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err)
	}
	return db
}

func booksHandler(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDb, err := db.Query("SELECT * FROM books")
	if err != nil {
		panic(err)
	}
	bookInfo := Book{}
	infoToSend := []Book{}
	for selDb.Next() {
		var id int
		var author, name, review string
		err = selDb.Scan(&id, &author, &name, &review)
		if err != nil {
			panic(err.Error())
		}
		bookInfo.Id = id
		bookInfo.Author = author
		bookInfo.Name = name
		bookInfo.Review = review
		infoToSend = append(infoToSend, bookInfo)
	}
	tpl.ExecuteTemplate(w, "books.gohtml", infoToSend)
	defer db.Close()

}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	err := tpl.ExecuteTemplate(w, "main.gohtml", nil)
	if err != nil {
		panic(err)
	}

}

func gotEntry(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	d := Book{}
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == "POST" {
		entryAuth := r.FormValue("author")
		entryName := r.FormValue("nameOfBook")
		entryReview := r.FormValue("reviewOfBook")
		insForm, err := db.Prepare("INSERT INTO books(id,author,name,review) VALUES(?,?,?,?)")
		if err != nil {
			panic(err)
		}
		insForm.Exec(nil, entryAuth, entryName, entryReview)
		log.Println("Entry is handled.")

		d.Author = entryAuth
		d.Name = entryName
		d.Review = entryReview

	}

	tpl.ExecuteTemplate(w, "gotentry.gohtml", d)

	//couldnt handle the redirecting, so we got http:superfluous response.WriteHeader
	//call from main.gotEntry but it doesnt open our index page
	//time.Sleep(3 * time.Second)
	//http.Redirect(w, r, "/", 302)

	defer db.Close()
}

func entryBooks(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "entry.gohtml", nil)
}

func main() {

	http.HandleFunc("/", indexHandler)
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./img"))))
	http.HandleFunc("/books", booksHandler)
	http.HandleFunc("/entryforbooks", entryBooks)
	http.HandleFunc("/gotentry", gotEntry)
	http.ListenAndServe("localhost:8080", nil)
}
