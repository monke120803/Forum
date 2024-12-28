package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3" // sqlite3 driver connects go with sql
)

type Forum struct {
	*sql.DB
}

// --------------------------- USER ------------------------//

// Pulls specific user's data and posts data from database and returns it as a User struct
func (data *Forum) GetUserProfile(username string) UserProfile {
	// Used to store the user's profile information
	user := UserProfile{}

	// Get a specific user's information from the 'users' table
	rows, err := data.DB.Query(`SELECT * FROM users where username= ?`, username)
	if err != nil {
		log.Fatal(err)
	}

	// Used to store the user's data so we can add it to struct later on
	var userID int
	var firstname string
	var lastname string
	var email string
	var nickname string
	var password string
	var age int
	var gender string
	var loggedin string

	// Scans through each column in the 'users' row and stores the data in the variables above
	for rows.Next() {
		err := rows.Scan(&userID, &nickname, &email, &password, &firstname, &lastname, &age, &gender, &loggedin)
		if err != nil {
			log.Fatal(err)
		}

		// This contains the specific user's data as well as all of their posts
		user = UserProfile{
			User: User{
				UserID:    userID,
				Username:  nickname,
				Firstname: firstname,
				Lastname:  lastname,
				Email:     email,
				LoggedIn:  loggedin,
			},
			CreatedPosts:  data.GetPosts(username),
			Notifications: data.GetNotifications(username),
			Messages:      data.getMessages(username),
		}
	}
	return user
}

//-------------------------- ACTIVITY STATUS ------------------//

// Updates user status after loginOut
func (data *Forum) UpdateStatus(loggedin string, username string) {
	stmt, err := data.DB.Prepare("UPDATE users SET loggedin = ? WHERE username = ?;")
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec(loggedin, username)
}

func (data *Forum) OnlineUsers() []User {
	var onlineuser User
	var onlineusers []User

	row, err1 := data.DB.Query(`SELECT userID, firstname, lastname, loggedin , username FROM users WHERE loggedin = 'true' ORDER BY username ASC;`)
	if err1 != nil {
		fmt.Println("Error with OnlineUsers func")
		return nil
	}

	// Scans through each column in the 'users' row and stores the data in the variables above
	for row.Next() {
		err := row.Scan(&onlineuser.UserID, &onlineuser.Firstname, &onlineuser.Lastname, &onlineuser.LoggedIn, &onlineuser.Username)
		if err != nil {
			log.Fatal(err)
		}

		onlineusers = append(onlineusers, onlineuser)

	}
	return onlineusers
}

// Offline status function, just setting loggedin to false
func (data *Forum) OfflineUser() []User {
	var offlineuser User
	var offlineusers []User

	row, err1 := data.DB.Query(`SELECT userID, firstname, lastname, loggedin, username FROM users WHERE loggedin = 'false';`)
	if err1 != nil {
		fmt.Println("Error with OfflineUsers func")
		return nil
	}
	// Scans through each column in the 'users' row and stores the data in the variables above
	for row.Next() {
		err := row.Scan(&offlineuser.UserID, &offlineuser.Firstname, &offlineuser.Lastname, &offlineuser.LoggedIn, &offlineuser.Username)
		if err != nil {
			log.Fatal(err)
		}
		offlineusers = append(offlineusers, offlineuser)

	}
	return offlineusers
}

// ------------------ notifications--------------------//

func (data *Forum) GetNotifications(username string) []Notifications {
	var noti Notifications

	var notifi []Notifications

	rows, err := data.DB.Query(`SELECT * FROM notifications WHERE recipient =?`, username)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err := rows.Scan(&noti.Sender, &noti.Recipient, &noti.Notification)
		if err != nil {
			log.Fatal(err)
		}

		notifi = append(notifi, noti)
	}

	return notifi
}

func (data *Forum) DeleteNotification(sender, recipient string) {
	rows, err := data.DB.Prepare("DELETE FROM notifications WHERE sender=? AND recipient=?")
	if err != nil {
		log.Fatal("ERROR Deleting Noti", err)
	}
	defer rows.Close()
	rows.Exec(sender, recipient)
	fmt.Println("DELETING ROW", sender, recipient)
}

func (data *Forum) CheckNotifications(sender, recipient string) bool {
	var noti Notifications

	fmt.Println("inside notification", recipient, sender)

	rows, err := data.DB.Query(`SELECT * FROM notifications WHERE recipient = ? AND sender = ?`, recipient, sender)
	if err != nil {
		log.Fatal(err)
	}
	n := 0
	for rows.Next() {
		err := rows.Scan(&noti.Sender, &noti.Recipient, &noti.Notification)
		if err != nil {
			log.Fatal(err)
		}
		n++
	}
	fmt.Println("check notif", n != 0)
	fmt.Println(noti)
	return n != 0
}

// --------------------------- POSTS ------------------------//

// Handles creation of new posts
func (data *Forum) CreatePost(post Post) {
	stmt, err := data.DB.Prepare("INSERT INTO posts (username, content, hashtag, creationDate) VALUES (?, ?, ?, ?);")
	if err != nil {
		log.Fatal(err)
	}

	// Uses data from post variable to insert into posts table
	_, err = stmt.Exec(post.Username, post.Content, post.Hashtag, post.CreatedAt)
	if err != nil {
		log.Fatal(err)
	}
}

// Pulls all posts from specific user and returns it as a slice of Post structs
func (data *Forum) GetPosts(username string) []Post {
	// Used to store all of the posts
	var posts []Post
	// Used to store invidiual post data
	var post Post

	rows, err := data.DB.Query(`SELECT * FROM posts WHERE username =?`, username)
	if err != nil {
		log.Fatal(err)
	}

	// Scans through every row where the username matches the username passed in
	for rows.Next() {
		// Populates post var with data from each post found in table
		err := rows.Scan(&post.PostID, &post.Username, &post.Content, &post.Hashtag, &post.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		// Adds each post found from specific user to posts slice
		posts = append(posts, post)
	}
	return posts
}

func (data *Forum) getLatestPosts() []Post {
	// Used to store all of the posts
	var posts []Post
	// Used to store invidiual post data
	var post Post

	rows, err := data.DB.Query(`SELECT * FROM posts`)
	if err != nil {
		log.Fatal(err)
	}

	// Scans through every post
	for rows.Next() {
		// Populates post var with data from each post found in table
		err := rows.Scan(&post.PostID, &post.Username, &post.Content, &post.Hashtag, &post.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		// Adds each post found from specific user to posts slice
		posts = append(posts, post)
	}
	return posts
}

func (data *Forum) getMessages(username string) []Message {
	// Used to store all of the messages
	var messages []Message
	// Used to store invidiual post data
	var message Message

	rows, err := data.DB.Query(`SELECT * FROM messages WHERE recipient=? OR sender=?`, username, username)
	if err != nil {
		log.Fatal(err)
	}

	// Scans through every post
	for rows.Next() {
		// Populates post var with data from each post found in table
		err := rows.Scan(&message.MessageID, &message.Sender, &message.Recipient, &message.Message, &message.CreationDate)
		if err != nil {
			log.Fatal(err)
		}
		// Adds each post found from specific user to messages slice
		messages = append(messages, message)
	}
	return messages
}

// ----------------------- COMMENTS -------------------------//

func (data *Forum) CreateComment(comment Comment) {
	stmt, err := data.DB.Prepare("INSERT INTO comments(postID, username, content, creationDate) VALUES(?, ?, ?, ?);")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(comment.PostID, comment.Username, comment.Content, comment.CreatedAt)
	if err != nil {
		log.Fatal(err)
	}
}

func (data *Forum) GetComments(postID int) []Comment {
	// Used to store all of the comments
	var comments []Comment

	// Used to store individual comment data
	var comment Comment

	rows, err := data.DB.Query(`SELECT * FROM comments WHERE postID =?`, postID)
	if err != nil {
		log.Fatal(err)
	}

	// Scans through every row where the postID matches the postID passed in
	for rows.Next() {
		// Populates post var with data from each post found in table
		err := rows.Scan(&comment.CommentID, &comment.PostID, &comment.Username, &comment.Content, &comment.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		// Adds each comment found from specific post to posts slice
		comments = append(comments, comment)
	}
	return comments
}

// --------------------------- HASHTAG ------------------------//

func (data *Forum) getLatestHashtags() []Hashtag {
	// Used to store all of the posts
	var hashtags []Hashtag
	// Used to store invidiual post data
	var hashtag Hashtag

	rows, err := data.DB.Query(`SELECT * FROM hashtags`)
	if err != nil {
		log.Fatal(err)
	}

	// Scans through every post
	for rows.Next() {
		// Populates post var with data from each post found in table
		err := rows.Scan(&hashtag.ID, &hashtag.Name, &hashtag.Count)
		if err != nil {
			log.Fatal(err)
		}
		// Adds each post found from specific user to posts slice
		hashtags = append(hashtags, hashtag)
	}
	return hashtags
}

// Updates hashtag value
func (data *Forum) UpdateHashtagCount(hashtag Hashtag) {
	// GET COUNT FOR DESIRED HASHTAG
	var hashtagCount int

	var hashtagName string
	var hashtagID int

	rows, err := data.DB.Query(`SELECT * FROM hashtags WHERE hashtagName =?`, hashtag.Name)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err := rows.Scan(&hashtagID, &hashtagName, &hashtagCount)
		if err != nil {
			log.Fatal(err)
		}
	}

	// NOW WE HAVE CURRENT COUNT FOR DESIRED HASHTAG, WE CAN ADD 1 TO IT
	hashtagCount++

	stmt, err := data.DB.Prepare("UPDATE hashtags SET hashtagCount = ? WHERE hashtagName = ?;")
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec(hashtagCount, hashtag.Name)

	fmt.Println("Hashtag count updated to", hashtagCount, "from", hashtag.Count, "for", hashtag.Name)
}

// --------------------------- SESSION ------------------------//

// Inserts session into sessions table
func (data *Forum) InsertSession(sess UserSession) {
	stmnt, err := data.DB.Prepare("INSERT INTO sessions (cookieValue, userID, username) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Println("Error inserting session into table:", err)
	}
	defer stmnt.Close()
	stmnt.Exec(sess.session, sess.userID, sess.username)
}

// TODO: Clarification
// User's cookie expires when browser is closed, delete the cookie from the database.
func (data *Forum) DeleteSession(w http.ResponseWriter, userID int) error {
	cookie := &http.Cookie{
		Name:   "session_token",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)

	stmt, err := data.DB.Prepare("DELETE FROM sessions WHERE userID=?;")
	// defer stmt.Close()
	stmt.Exec(userID)
	if err != nil {
		log.Fatal(err)
		fmt.Println("DeleteSession err: ", err)
		return err
	}
	return nil
}

// Checks all sessions from sessions table and returns latest session
func (data *Forum) GetSession(cookie string) UserSession {
	// Used to store session data
	session := UserSession{}

	// Checks all sessions from sessions table
	rows, err := data.DB.Query(`SELECT * FROM sessions WHERE cookieValue=?;`, cookie)
	if err != nil {
		log.Fatal(err)
	}

	// Used to store individual session data
	var userID int
	var cookieValue string
	var userName string

	// For each session found, populate the variable above
	for rows.Next() {
		err2 := rows.Scan(&userID, &cookieValue, &userName)
		if err2 != nil {
			log.Fatal(err2)
		}
		// Overwrites every session, leaving only data for the latest session
		session = UserSession{
			userID:   userID,
			session:  cookieValue,
			username: userName,
		}
	}

	return session
}

func (data *Forum) SelectingLoadingMessage(username, recipient string) []Chat {
	fmt.Println("SelectingLoadingMessage")
	var loading Chat
	var conversation []Chat

	rows, err := data.DB.Query(`SELECT sender, recipient, message, creationDate FROM messages where (sender = ? AND recipient = ?) OR (sender = ? AND recipient = ?)`, username, recipient, recipient, username)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("after sql:", rows)
	for rows.Next() {
		// fmt.Println("inside query loop")
		err := rows.Scan(&loading.MessageSender, &loading.MessageRecipient, &loading.Message, &loading.CreatedAt)
		if err != nil {
			log.Fatal("conversation error", err)
		}
		// fmt.Println("Messages", username, recipient, ":", &loading.Message)
		conversation = append(conversation, loading)
	}
	// fmt.Println("Con", conversation)
	return conversation
}

func (data *Forum) SaveChat(chat Chat) Chat {
	stmnt, err := data.DB.Prepare("INSERT INTO messages (sender, recipient ,message, creationDate) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmnt.Exec(chat.MessageSender, chat.MessageRecipient, chat.Message, chat.CreatedAt)
	if err != nil {
		log.Fatal(err)
	}
	return chat
}

func (data *Forum) SaveNotifications(noti Notifications) Notifications {
	stmnt, err := data.DB.Prepare("INSERT INTO notifications (sender, recipient, notification) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal("Error inserting to noti table: ", err)
	}

	_, err = stmnt.Exec(noti.Sender, noti.Recipient, noti.Notification)
	if err != nil {
		log.Fatal(err)
	}
	return noti
}

//-------------------------  TABLES -------------------------//

// Used when starting server - Ensures all tables are created to avoid errors
func CheckTablesExist(db *sql.DB, table string) {
	_, table_check := db.Query("select * from " + table + ";")
	if table_check != nil {
		fmt.Println("Error: " + table + " table doesn't exist in database.")

		if table == "users" {
			fmt.Println("Creating users table...")
			users_table := `CREATE TABLE IF NOT EXISTS users (
					"userID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
					"username" TEXT NOT NULL,
					"email" TEXT NOT NULL,
					"password" TEXT NOT NULL,
					"firstname" TEXT,
					"lastname" TEXT,
					"age" INTEGER NOT NULL, 
					"gender" TEXT NOT NULL,
					"loggedin" TEXT
					);`

			users, errUser := db.Prepare(users_table)
			if errUser != nil {
				log.Fatal(errUser)
			}
			users.Exec()
		}

		if table == "posts" {
			fmt.Println("Creating posts table...")
			posts_table := `CREATE TABLE IF NOT EXISTS posts (
					"postID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
					"username" TEXT REFERENCES sesssion(userID),
					"content" TEXT NOT NULL, 
					"hashtag" TEXT NOT NULL,
					"creationDate" TIMESTAMP
					);`

			posts, errTable := db.Prepare(posts_table)
			if errTable != nil {
				log.Fatal(errTable)
			}
			posts.Exec()
		}

		if table == "comments" {
			fmt.Println("Creating comments table...")
			comments_table := `CREATE TABLE IF NOT EXISTS comments (
					"commentID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
					"postID" INTEGER NOT NULL,
					"username" TEXT REFERENCES sesssion(userID),
					"content" TEXT NOT NULL, 
					"creationDate" TIMESTAMP,
					FOREIGN KEY(postID)REFERENCES posts(postID)
					);`

			comments, errCommments := db.Prepare(comments_table)
			if errCommments != nil {
				log.Fatal(errCommments)
			}
			comments.Exec()
		}

		if table == "hashtags" {
			fmt.Println("Creating hashtags table...")
			hashtags_table := `CREATE TABLE IF NOT EXISTS hashtags (
				"hashtagID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
				"hashtagName" TEXT NOT NULL,
				"hashtagCount" INTEGER NOT NULL
				);`

			hashtags, errHashtags := db.Prepare(hashtags_table)
			if errHashtags != nil {
				log.Fatal(errHashtags)
			}
			hashtags.Exec()

			fmt.Println("Inserting hashtags into hashtags table...")
			stmt, err := db.Prepare("INSERT INTO hashtags (hashtagName, hashtagCount) VALUES (?, ?);")
			if err != nil {
				log.Fatal(err)
			}

			// Used to store hashtag names
			hashtagSlice := make([]string, 7)
			hashtagSlice[0] = "#Tech"
			hashtagSlice[1] = "#Food"
			hashtagSlice[2] = "#Art"
			hashtagSlice[3] = "#Sports"
			hashtagSlice[4] = "#Fitness"
			hashtagSlice[5] = "#Travel"
			hashtagSlice[6] = "#Misc"

			// insert all hashtags into hashtags table
			for _, hashtag := range hashtagSlice {
				_, err = stmt.Exec(hashtag, 0)
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		if table == "sessions" {
			fmt.Println("Creating sessions table...")
			sessions_table := `CREATE TABLE IF NOT EXISTS sessions (
				userID INTEGER NOT NULL,
				cookieValue TEXT NOT NULL,
				username TEXT REFERENCES users(username),
				FOREIGN KEY(userID) REFERENCES Users(userID)
					);`

			sessions, errSession := db.Prepare(sessions_table)
			if errSession != nil {
				log.Fatal(errSession)
			}
			sessions.Exec()
		}
		// draft table for message, maybe update fields?
		if table == "messages" {
			fmt.Println("Creating messages table...")
			messages_table := `CREATE TABLE IF NOT EXISTS messages (
					"messageID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
					"sender" TEXT REFERENCES users(username), 
					"recipient" TEXT REFERENCES users(username),
					"message" CHAR(200),
					"creationDate" TIMESTAMP
					);`

			messages, errTable := db.Prepare(messages_table)
			if errTable != nil {
				log.Fatal(errTable)
			}
			messages.Exec()
		}
		// draft table for message, maybe update fields?
		if table == "chat" {
			fmt.Println("Creating chat table...")
			chat_table := `CREATE TABLE IF NOT EXISTS chat (
					"chatID" INTEGER PRIMARY KEY AUTOINCREMENT,
					"username1" TEXT REFERENCES users(username), 
					"username2" TEXT REFERENCES users(username),
					"creationDate" TIMESTAMP
					);`

			chat, errTable := db.Prepare(chat_table)
			if errTable != nil {
				log.Fatal(errTable)
			}
			chat.Exec()
		}
		if table == "notifications" {
			fmt.Println("Creating notifications table...")
			notifications_table := `CREATE TABLE IF NOT EXISTS notifications(
                    "sender" TEXT NOT NULL,
                    "recipient" TEXT,
                    "notification" INTEGER NOT NULL
                    );`

			notify, errNotify := db.Prepare(notifications_table)
			if errNotify != nil {
				log.Fatal(errNotify)
			}
			notify.Exec()
		}

	}
}

// Check all required tables exist in database, and create them if they don't
func Connect(db *sql.DB) *Forum {
	for _, table := range []string{"users", "posts", "comments", "hashtags", "sessions", "messages", "chat", "notifications"} {
		CheckTablesExist(db, table)
	}
	return &Forum{
		DB: db,
	}
}
