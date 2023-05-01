package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cfthoo/todo-app/api"
	oauth2api "github.com/cfthoo/todo-app/api/oauth"
	"github.com/cfthoo/todo-app/pkg/controller"
	conn "github.com/cfthoo/todo-app/pkg/db"
	"github.com/cfthoo/todo-app/pkg/db/repo"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load(".env")
	db, err := conn.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// db migration
	conn.RunDBMigration(os.Getenv("MIGRATION_URL"), os.Getenv("DB_SOURCE"))

	u := &api.Handler{
		TodoListDAO: &repo.TodoList{
			DB: db,
		},
	}

	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/google/login", oauth2api.LoginHandler)
	r.HandleFunc("/google/callback", oauth2api.CallbackHandler)
	r.HandleFunc("/fb/login", oauth2api.LoginHandler)
	r.HandleFunc("/fb/callback", oauth2api.CallbackHandler)
	r.HandleFunc("/github/login", oauth2api.LoginHandler)
	r.HandleFunc("/github/callback", oauth2api.CallbackHandler)
	r.Methods(http.MethodPost).Path("/todolist").Handler(controller.ValidateJWT(u.Create()))
	r.Methods(http.MethodGet).Path("/todolist").Handler(controller.ValidateJWT(u.List()))
	r.Methods(http.MethodGet).Path(fmt.Sprintf("/todolist/{%s}", "id")).Handler(controller.ValidateJWT(u.FetchByID()))
	r.Methods(http.MethodPut).Path("/todolist").Handler(controller.ValidateJWT(u.Update()))
	r.Methods(http.MethodDelete).Path(fmt.Sprintf("/todolist/{%s}", "id")).Handler(controller.ValidateJWT(u.Delete()))
	// Start the HTTP server
	addr := ":8080"
	log.Printf("Server listening on %s", addr)
	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	var html = `
        <html>
            <body>
                <a href="/google/login">Google Log In</a><br>
                <a href="/fb/login">Facebook Log In</a><br>
                <a href="/github/login">Github Log In</a><br>
            </body>
        </html>`
	fmt.Fprint(w, html)
}
