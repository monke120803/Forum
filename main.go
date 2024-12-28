package main

import (
	"database/sql" // Provides a generic interface for working with SQL databases
	"fmt"
	"log"
	"net/http" // Provides functions for creating and hosting HTTP servers, as well as handling HTTP requests

	"real-time-forum/handlers" // Contains code for handling HTTP requests in the application

	_ "github.com/mattn/go-sqlite3" // Provides an SQLite3 driver for the database/sql package
)

func main() {
	// Open a connection to the database and prints an error if it fails to open
	database, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		fmt.Println(err.Error())
	}
	data := handlers.Connect(database)
	defer database.Close() // This ensures that the database is properly closed when the program exits

	// THIS IS TO DELETE A TABLE WITHOUT DELETING THE WHOLE DATABASE! (USEFUL)
	// rmv := `DROP TABLE IF EXISTS messages;`
	// _, err = data.DB.Exec(rmv)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("should be deleted")

	// Create a file server to serve content from the static directory
	fileServer := http.FileServer(http.Dir("./static"))

	// Redirect requests to the root URL to the static directory
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Handles redirects to the appropriate functions
	http.HandleFunc("/", data.Home)
	http.HandleFunc("/login", data.LoginHandler)
	http.HandleFunc("/register", data.RegistrationHandler)
	http.HandleFunc("/post", data.Post)
	http.HandleFunc("/comment", data.Comment)
	http.HandleFunc("/chat", data.Chat)
	http.HandleFunc("/loadingmessage", data.LoadingMessage)
	http.HandleFunc("/logout", data.LogoutUser)
	http.HandleFunc("/getPosts", data.SendLatestPosts)
	http.HandleFunc("/getHashtags", data.SendLatestHashtags)
	http.HandleFunc("/updateHashtag", data.UpdateHashtag)
	http.HandleFunc("/sendComments", data.SendComments)
	http.HandleFunc("/usersStatus", data.SendLatestActivity)
	http.HandleFunc("/checkCookie", data.CheckCookie)
	http.HandleFunc("/notification", data.SendNotification)
	// http.HandleFunc("/fetchMessages", data.FetchTime)

	// Create the hub that will manage the connections and communication with clients
	hub := handlers.NewHub(data)
	go hub.Run()
	go hub.LogConns()

	// When a request is received at the "/ws" endpoint
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		data.ServeWs(hub, w, r)
	})

	fmt.Println("Server started at http://localhost:8080")

	// Start the HTTP server and listen for incoming requests, printing an error if it fails to start
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
