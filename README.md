# Real-Time-Forum

This project is a forum website created using HTML, CSS, Javascript and Go. The forum allows users to register, login, create posts, comment on posts, and send private messages to each other. We also used technologies such as SQLite to store data and the Gorilla websocket package to handle websocket connections and client-server communication.

## How we did it

Firstly, we implemented the registration and login features of the forum. This included creating a form for users to register with, and handling the registration and login requests using Go and websockets. We also added the ability for users to log out from any page on the forum.

To send the user's registration and login credentials through websockets, we used the WebSocket object in JavaScript to establish a websocket connection with the Go server. Then, we used the send method of the WebSocket object to send the user's credentials to the Go server as a stringified JSON object.

On the Go server, we used the Gorilla websocket package to handle the websocket connection and messages. We defined a handler function that listens for websocket messages from the JavaScript client, and uses the json.Unmarshal method to convert the stringified JSON object into a Go data structure. We then processed the user's credentials and performed the appropriate action, such as registering the user or logging them in.

Next, we implemented the ability for users to create posts and comments, and to view posts in a feed display. We also added the ability for users to see comments on a post by clicking on the post in the feed.

To store the user's posts in the posts table, we first used the fetch API in JavaScript to send an HTTP request to the Go server, with the user's post data included in the request body. The Go server then used the sqlite3 package to establish a connection with the SQLite database, and used the INSERT SQL statement to insert the user's post data into the posts table.

Here is an example of how we stored the user's posts into the posts table using JavaScript and Go:

```
// In JavaScript:

// Use the fetch API to send an HTTP request to the Go server, with the user's post data included in the request body
fetch("http://localhost:8080/posts", {
  method: "POST",
  body: JSON.stringify({
    content: "This is the content of my first post",
    hashtag: "#Misc"
  })
});

// In Go:

// Use the sqlite3 package to establish a connection with the SQLite database
db, err := sql.Open("sqlite3", "./forum.db")
if err != nil {
  log.Fatal(err)
}
defer db.Close()

// Use the INSERT SQL statement to insert the user's post data into the posts table
_, err = db.Exec("INSERT INTO posts (content, category, username) VALUES (?, ?, ?)",
  post.Title, post.Content, post.Category, user.ID)
if err != nil {
  log.Fatal(err)
}
```

Finally, we added the private messaging feature to the forum, including the ability to see who is online and able to receive messages, and to send and receive messages in real time using websockets.

## Issues with Go and websockets

One issue we ran into while working on this project was dealing with concurrent access to shared data. Since Go is a concurrent language, it is possible for multiple goroutines to access and modify shared data at the same time. This can lead to race conditions and other problems if not handled properly. To solve this issue, we used Go channels and mutexes to ensure that shared data was accessed and modified in a safe and consistent manner.

Another issue we encountered was with the websocket connections between the Go server and the frontend JavaScript code. It was sometimes difficult to correctly handle websocket events and messages in both the Go server and the JavaScript client, and to keep the websocket connections synchronized. We solved this issue by carefully designing the websocket message protocol and implementing robust error handling in both the Go server and the JavaScript client.

## Running the Go server

To run the Go server, first make sure that you have Go installed on your system. You can download and install Go from the [official Go website](https://golang.org/).

Once Go is installed, you can clone this repository and navigate to the project directory using a command line interface. Then, you can run the following command to build and run the Go server:

```
go build && ./go-forum
```

This will start the Go server on your local machine, listening on port 8080. You can then visit the forum website by opening a web browser and going to the following URL:

```
http://localhost:8080
```

You can then register and login to the forum, and start using the various features of the forum.

## Conclusion

Overall, we found this project to be a challenging but rewarding experience. Working with Go and websockets required us to learn and apply various new concepts.