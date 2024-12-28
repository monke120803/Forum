package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// Try to parse the index.html file and if it fails, log the error
func (data *Forum) Home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/index.html")
	if err != nil {
		http.Error(w, "500 Internal error", http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, ""); err != nil {
		http.Error(w, "500 Internal error", http.StatusInternalServerError)
		return
	}
}

func (data *Forum) CheckCookie(w http.ResponseWriter, r *http.Request) {
	var cookieValue CookieValue

	// Decode the JSON data from the request body into the comment variable
	json.NewDecoder(r.Body).Decode(&cookieValue)

	u := data.GetSession(cookieValue.CookieValue)
	userName := (u.username)

	userInfo := data.GetUserProfile(userName)

	js, err := json.Marshal(userInfo)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(js))
}

func (data *Forum) SendLatestActivity(w http.ResponseWriter, r *http.Request) {
	// Send user information back to client using JSON format

	x, err := r.Cookie("session_token")
	if err != nil {
		log.Fatal(err)
	}
	sessionvalue := x.Value

	sess := data.GetSession(sessionvalue)

	onlineactivity := OnlineActivity{
		Online:        data.OnlineUsers(),
		Offline:       data.OfflineUser(),
		Notifications: data.GetNotifications(sess.username),
	}

	js, err := json.Marshal(onlineactivity)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK) // Checked in authentication.js, alerts user
	w.Write([]byte(js))
}

// Handles receiving the comment data and adding it to the 'comments' table in the database
func (data *Forum) Comment(w http.ResponseWriter, r *http.Request) {
	var comment Comment

	// Decode the JSON data from the request body into the comment variable
	json.NewDecoder(r.Body).Decode(&comment)

	// w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))

	// feches current session value
	x, err := r.Cookie("session_token")
	if err != nil {
		log.Fatal(err)
	}
	sessionvalue := x.Value

	sess := data.GetSession(sessionvalue)
	time := time.Now()

	data.CreateComment(Comment{
		PostID:    comment.PostID,
		Username:  sess.username,
		Content:   comment.Content,
		CreatedAt: time,
	})

	fmt.Println(comment)
}

func (data *Forum) SendComments(w http.ResponseWriter, r *http.Request) {
	var comment Comment

	json.NewDecoder(r.Body).Decode(&comment)

	comments := data.GetComments(comment.PostID)
	js, err := json.Marshal(comments)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(js))
}

// Handles receiving the post data and adding it to the 'posts' table in the database
func (data *Forum) Post(w http.ResponseWriter, r *http.Request) {
	// Decodes posts data into post variable
	var post Post

	// Decode the JSON data from the request body into the post variable
	json.NewDecoder(r.Body).Decode(&post)

	// w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))

	// feches current session value
	x, err := r.Cookie("session_token")
	if err != nil {
		log.Fatal()
	}
	sessionvalue := x.Value

	// Convert data into variables for easier use
	hashtag := post.Hashtag
	time := time.Now()
	content := post.Content

	sess := data.GetSession(sessionvalue)

	// Inserts post into the 'posts' table of the database
	data.CreatePost(Post{
		// username from current session
		Username:  sess.username,
		Content:   content,
		Hashtag:   hashtag,
		CreatedAt: time,
	})
}

func (data *Forum) Chat(w http.ResponseWriter, r *http.Request) {
	var chat Chat

	json.NewDecoder(r.Body).Decode(&chat)

	w.Write([]byte("chat ok"))

	// feches current session value
	x, err := r.Cookie("session_token")
	if err != nil {
		log.Fatal()
	}
	sessionvalue := x.Value

	content := chat.Message
	time := time.Now()
	recipient := chat.MessageRecipient

	sess := data.GetSession(sessionvalue)

	if !data.CheckNotifications(sess.username, recipient) {
		data.SaveNotifications(Notifications{
			Sender:       sess.username,
			Recipient:    recipient,
			Notification: 1,
		})
	}

	chat = data.SaveChat(Chat{
		MessageSender:    sess.username,
		MessageRecipient: recipient,
		Message:          content,
		CreatedAt:        time,
	})

	fmt.Println("SENDER: ", chat.MessageSender)
	fmt.Println("RECEPIENT: ", chat.MessageRecipient)
	fmt.Println("ACTUAL MESSAGE: ", chat.Message)
	fmt.Println("TIME: ", chat.CreatedAt)
}

func (data *Forum) SendLatestPosts(w http.ResponseWriter, r *http.Request) {
	// Send user information back to client using JSON format
	posts := data.getLatestPosts()
	// fmt.Println(userInfo)
	js, err := json.Marshal(posts)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK) // Checked in authentication.js, alerts user
	w.Write([]byte(js))
}

func (data *Forum) FetchTime(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in function")
	// Send user information back to client using JSON format
	// orderByTime := data.getMessages()

	// fmt.Println("Fetching time from message struct: ", orderByTime)

	// fmt.Println(userInfo)
	// js, err := json.Marshal(orderByTime)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// w.WriteHeader(http.StatusOK) // Checked in authentication.js, alerts user
	// w.Write([]byte(js))
}

// Updates hashtag count for specific hashtag when called
func (data *Forum) UpdateHashtag(w http.ResponseWriter, r *http.Request) {
	// Decodes posts data into post variable
	var hashtag Hashtag

	// Decode the JSON data from the request body into the post variable
	json.NewDecoder(r.Body).Decode(&hashtag)

	// w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))

	// Convert data into variables for easier use
	hashID := hashtag.ID
	hashName := hashtag.Name
	hashCount := hashtag.Count

	fmt.Println("hashID:", hashID)
	fmt.Println("hashName:", hashName)
	fmt.Println("hashtagCount:", hashCount)

	// Updates hashtag count in the 'hashtags' table of the database
	data.UpdateHashtagCount(Hashtag{
		ID:    hashID,
		Name:  hashName,
		Count: hashCount,
	})
}

func (data *Forum) SendLatestHashtags(w http.ResponseWriter, r *http.Request) {
	// Send user information back to client using JSON format
	hashtags := data.getLatestHashtags()

	js, err := json.Marshal(hashtags)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK) // Checked in authentication.js, alerts user
	w.Write([]byte(js))
}

// Handles the registration of new users - validates the data and adds it to the 'users' table in database
func (data *Forum) RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	// Decodes registration data into user variable
	var user RegisterData
	json.NewDecoder(r.Body).Decode(&user)

	// Used in conjunction with the 'strings.ContainsAny' function to ensure the age entered is strictly numeric
	numChars := "0123456789"

	// Ensures all required fields are filled out, and that the age is strictly numeric
	if len(user.Firstname) == 0 || len(user.Lastname) == 0 || len(user.Email) == 0 || len(user.Username) == 0 || (len(user.Age) == 0 || !strings.ContainsAny(user.Age, numChars)) || user.Gender == "Gender" || len(user.Password) == 0 {
		// This HTTP status code is then checked in authentication.js and the user is alerted to the missing/invalid fields
		w.WriteHeader(http.StatusNotAcceptable)
	} else {
		// Uses web socket to read the information
		w.Header().Set("Content-type", "application/text")

		// These are initially false, and are only set to true if the email/username is not found in the database (registration is available and will not overwrite existing data)
		emailValid := false
		usernameValid := false
		user.LoggedIn = "false"

		/* ---------------------------------------------------------------- */
		/*       CHECKING IF EMAIL/USERNAME ALREADY EXISTS IN DATABASE      */
		/* ---------------------------------------------------------------- */

		/* --- Queries through each table, checks if data already exists -- */

		// EMAIL CHECK
		row := data.DB.QueryRow("select email from users where email= ?", user.Email)
		temp := "" // If email is not found, temp variable will remain empty
		row.Scan(&temp)
		if temp == "" {
			emailValid = true
		}

		// USERNAME CHECK
		row = data.DB.QueryRow("select username from users where username= ?", user.Username)
		temp = "" // If username is not found, temp variable will remain empty
		row.Scan(&temp)
		if temp == "" {
			usernameValid = true
		}

		// If both email and username are valid, we can successfully register the user into the database
		if emailValid && usernameValid {
			// Generates hash from password
			var passwordHash []byte
			passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				fmt.Println("Error hashing password:", err)
				return
			}

			// Inserts registration data into the 'users' table of the database
			query, err := data.DB.Prepare("INSERT INTO users(username, email, password, firstname, lastname, age, gender,loggedin) VALUES(?, ?, ?, ?, ?, ?, ?, ?);")
			if err != nil {
				log.Fatal(err)
			}

			_, err = query.Exec(user.Username, user.Email, string(passwordHash), user.Firstname, user.Lastname, user.Age, user.Gender, user.LoggedIn)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("SUCCESS: User successfully registered into users table.")
			w.WriteHeader(http.StatusOK) // Checked in authentication.js, alerts user
		} else {
			fmt.Println("ERROR: Username or email already exists.")
			w.WriteHeader(http.StatusBadRequest) // Checked in authentication.js, alerts user
		}
	}
}

// Handles the login of existing users - validates the data and checks if it exists in the 'users' table in database
func (data *Forum) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Decodes session and login data into variables
	var sess UserSession
	var user LoginData

	json.NewDecoder(r.Body).Decode(&user)
	w.Header().Set("Content-type", "application/text")

	// Only set to true if the email/username IS found in the database
	emailPassCombinationValid := false
	userPassCombinationValid := false

	// Checks if user entered an email or username
	enteredEmail := strings.Contains(user.Username, "@")

	/* ---------------------------------------------------------------- */
	/*               CHECKING EMAIL/USER PASS COMBINATIONS              */
	/* ---------------------------------------------------------------- */

	/* --- Queries through each table, checks if data exists --- */

	// EMAIL CHECK
	if enteredEmail {
		// Checks if email/pass combination exists in database
		var passwordHash string
		row := data.DB.QueryRow("SELECT password FROM users WHERE email = ?", user.Username)
		err := row.Scan(&passwordHash)
		if err != nil {
			fmt.Println("Error with password hash:", err)
		}
		// If the password hash matches the password entered, the email/pass combination is valid
		err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(user.Password))
		if err == nil {
			emailPassCombinationValid = true
		}
	} else {
		// Checks if username/pass combination exists in database
		var passwordHash string
		row := data.DB.QueryRow("SELECT password FROM users WHERE username = ?", user.Username)
		err := row.Scan(&passwordHash)
		if err != nil {
			fmt.Println("Error with password hash:", err)
		}
		// If the password hash matches the password entered, the user/pass combination is valid
		err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(user.Password))
		if err == nil {
			userPassCombinationValid = true
		}
	}
	var usID int = 5

	// If either combination is valid, we can successfully log the user in
	if emailPassCombinationValid || userPassCombinationValid {
		fmt.Println("SUCCESS: User logged in.")

		row := data.DB.QueryRow("SELECT userID FROM users WHERE username = ?;", user.Username)
		err := row.Scan(&usID)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println("usID:", usID)
		// fmt.Println("user.Username:", user.Username)

		// Creates a new session for the user
		sess.username = user.Username
		sess.userID = usID
		sess.max_age = 18000
		sess.session = (uuid.NewV4().String() + "&" + strconv.Itoa(sess.userID))
		user.LoggedIn = "true"

		// Set client cookie for "session_token" as session token we just generated, also set expiry time to 120 minutes
		http.SetCookie(w, &http.Cookie{
			Name:   "session_token",
			Value:  sess.session,
			MaxAge: 900,
		})

		// Insert data into session variable
		data.InsertSession(sess)

		data.UpdateStatus(user.LoggedIn, user.Username)

		// JUST FOR TESTING
		x := data.OfflineUser()
		fmt.Println("offline:", x)

		y := data.OnlineUsers()
		fmt.Println("online:", y)
		z := data.getMessages(user.Username)
		fmt.Println("User details: ", z)
		// Send user information back to client using JSON format
		userInfo := data.GetUserProfile(user.Username)
		// fmt.Println(userInfo)
		fmt.Println("---All Noti", userInfo.Notifications)
		js, err := json.Marshal(userInfo)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusOK) // Checked in authentication.js, alerts user
		w.Write([]byte(js))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Error: Email or password is incorrect.") // Checked in authentication.js, alerts user
	}
}

// // logout handle
func (data *Forum) LogoutUser(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		log.Fatal(err)
	}

	sess := data.GetSession(c.Value)

	fmt.Printf("User %d wants to logout\n", sess.userID)
	loggedin := "false"

	data.DeleteSession(w, sess.userID)
	data.UpdateStatus(loggedin, sess.username)

	// JUST FOR TESTING
	x := data.OfflineUser()
	fmt.Println("offline:", x[0])

	y := data.OnlineUsers()
	fmt.Println("online:", y)

	// Send user information back to client using JSON format
	userInfo := data.GetUserProfile(sess.username)
	// fmt.Println(userInfo)
	js, err := json.Marshal(userInfo)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK) // Checked in authentication.js, alerts user
	w.Write([]byte(js))
}

func (data *Forum) LoadingMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Loading Messages")
	var loading LoadingMessage

	err := json.NewDecoder(r.Body).Decode(&loading)
	if err != nil {
		log.Fatal("Loading handler error: ", err)
	}

	// w.Write([]byte("chat ok"))
	fmt.Println(loading.SendersUsername, loading.RecipientsUsername)
	conv := data.SelectingLoadingMessage(
		loading.SendersUsername,
		loading.RecipientsUsername,
	)

	data.DeleteNotification(loading.RecipientsUsername, loading.SendersUsername)

	js, err := json.Marshal(conv)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("Conversation in auth: ", conv)

	// w.WriteHeader(http.StatusOK) // Checked in authentication.js, alerts user
	w.Write([]byte(js))
}

func (data *Forum) SendNotification(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Sending Notifications")
	var notification Notifications

	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		log.Fatal("SendNotification Handler error: ", err)
	}
	fmt.Println("This should be our notifications: ", notification)
	w.Write([]byte("ok from noti handler"))
}
