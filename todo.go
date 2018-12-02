package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bmizerany/pat"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"
)

type Todo struct {
	ID        string    `json:"id"`
	Todo      string    `json:"todo"`
	CreatedAt time.Time `json:"created_at"`
}

type Todos []Todo

var mainDB *sql.DB

func main() {

	dbPtr := flag.String("db", "sqlite3://:memory:", "Database URL")
	portPtr := flag.String("addr", "0.0.0.0:8080", "HTTP Address")

	flag.Parse()

	u, err := url.Parse(*dbPtr)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(u.Scheme, strings.Replace(*dbPtr, "sqlite3://", "", -1))
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS todos (id CHAR(36) PRIMARY KEY, todo VARCHAR(255) NULL, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`)
	if err != nil {
		panic(err)
	}

	if u.Scheme == "sqlite3" {
		_, err = db.Exec(`INSERT INTO todos(id, todo) VALUES
('0a946288-f258-11e8-8470-f7f86366161e', 'feed the cats'),
('0a92b9d8-f258-11e8-9aa4-afea51c69f05', 'prepare workshop demo'),
('0a9587c6-f258-11e8-9963-274a0346e269', 'foo'),
('0a962046-f258-11e8-903b-3bc836797b33', 'bar');`)
		if err != nil {
			panic(err)
		}
	}

	mainDB = db

	r := pat.New()
	r.Get("/ping", http.HandlerFunc(ping))
	r.Get("/echo", http.HandlerFunc(echo))

	r.Del("/:id", http.HandlerFunc(deleteByID))
	r.Get("/:id", http.HandlerFunc(getByID))
	r.Put("/:id", http.HandlerFunc(updateByID))

	r.Get("/", http.HandlerFunc(getAll))
	r.Post("/", http.HandlerFunc(insert))

	http.Handle("/", r)

	port := strings.Split(*portPtr, ":")[1]

	fmt.Printf("Todo app started. Listening on port %s.\n", port)
	err = http.ListenAndServe(*portPtr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "🏓 PONG !")
}

func echo(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Host          string
		Method        string
		URL           string
		Proto         string
		RemoteAddress string
		FormValues    map[string][]string
		PostValues    map[string][]string
		Header        map[string][]string
	}{r.Host, r.Method, r.URL.String(), r.Proto, r.RemoteAddr, r.Form, r.PostForm, r.Header}

	b, _ := json.Marshal(req)

	fmt.Fprintf(w, string(b))
}

func getAll(w http.ResponseWriter, r *http.Request) {
	rows, err := mainDB.Query("SELECT * FROM todos")
	checkErr(err)
	var todos Todos
	for rows.Next() {
		var todo Todo
		err = rows.Scan(&todo.ID, &todo.Todo, &todo.CreatedAt)
		checkErr(err)
		todos = append(todos, todo)
	}

	t := template.New("fieldname example")
	t, _ = t.Parse(`<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>TODO</title>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap-theme.min.css" integrity="sha384-rHyoN1iRsVXV4nD0JutlnGaslCJuC7uwjduW9SVrLvRYooPp2bWYgmgJQIXwl/Sp" crossorigin="anonymous">
  </head>
  <body>
    <br>
    <div class="container">
      <div class="row">
        <div class="col-md-4 offset-md-4">
          <div class="panel panel-default">
            <div class="panel-heading">todo demo app</div>
            <div class="panel-body">
              <form class="form-inline" method="POST">
                <div class="form-group">
                  <div class="input-group">
                    <input type="text" class="form-control" name="todo">
                  </div>
                </div>
                <button type="submit" class="btn btn-primary glyphicon glyphicon-floppy-save"></button>
              </form>
              <br>
              <ul>
              {{range .}}
                <li data-toggle="tooltip" data-placement="top" title="ID: {{.ID}}"><a href="/{{.ID}}">{{.Todo}}</a></li>
              {{end}}
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
    <script src="https://code.jquery.com/jquery-1.12.4.min.js" integrity="sha256-ZosEbRLbNQzLpnKIkEdrPv7lOy9C27hHQ+Xp8a4MxAQ=" crossorigin="anonymous"></script>
    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js" integrity="sha384-Tc5IQib027qvyjSMfHjOMaLkfuWVxZxUPnCJA7l2mCWNIpG9mGCD8wGNIcPD7Txa" crossorigin="anonymous"></script>
  </body>
</html>`)
	t.Execute(w, todos)
}

func getByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")
	rows, err := mainDB.Query("SELECT * FROM todos where id=$1", id)
	checkErr(err)
	var todo Todo
	for rows.Next() {
		err = rows.Scan(&todo.ID, &todo.Todo, &todo.CreatedAt)
		checkErr(err)
	}
	jsonB, errMarshal := json.Marshal(todo)
	checkErr(errMarshal)
	fmt.Fprintf(w, "%s", string(jsonB))
}

//foo') DROP TABLE todos; '
func insert(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	todo.Todo = r.FormValue("todo")
	todo.ID = uuid.NewV4().String()
	todo.CreatedAt = time.Now().UTC()

	_, err := mainDB.Exec("INSERT INTO todos(id, todo, created_at) values($1, $2, $3)", todo.ID, todo.Todo, todo.CreatedAt)
	checkErr(err)
	jsonB, errMarshal := json.Marshal(todo)
	checkErr(errMarshal)

	w.Header().Set("Location", "/")
	w.WriteHeader(301)
	w.Write(jsonB)
}

func updateByID(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	todo.Todo = r.FormValue("todo")
	todo.ID = r.URL.Query().Get(":id")
	result, err := mainDB.Exec("UPDATE todos SET todo=$1 WHERE id=$2", todo.Todo, todo.ID)
	checkErr(err)
	rowAffected, errLast := result.RowsAffected()
	checkErr(errLast)
	if rowAffected > 0 {
		jsonB, errMarshal := json.Marshal(todo)
		checkErr(errMarshal)
		fmt.Fprintf(w, "%s", string(jsonB))
	} else {
		fmt.Fprintf(w, "{row_affected=%d}", rowAffected)
	}

}

func deleteByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")
	result, err := mainDB.Exec("DELETE FROM todos WHERE id=$1", id)
	checkErr(err)
	rowAffected, errRow := result.RowsAffected()
	checkErr(errRow)
	fmt.Fprintf(w, "{row_affected=%d}", rowAffected)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
