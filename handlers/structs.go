package handlers

import "time"

/* ---------------------------------------------------------------- */
/*             USED FOR CREATING VARIABLES TO STORE DATA            */
/* ---------------------------------------------------------------- */

type RegisterData struct {
	Firstname string `json:"firstName"`
	Lastname  string `json:"lastName"`
	Email     string `json:"email"`
	Username  string `json:"newusername"`
	Age       string `json:"age"`
	Gender    string `json:"gender"`
	Password  string `json:"newpassword"`
	LoggedIn  string
}

type User struct {
	UserID    int    `json:"userID"`
	Firstname string `json:"firstName"`
	Lastname  string `json:"lastName"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	LoggedIn  string
}

type OnlineActivity struct {
	Online        []User
	Offline       []User
	Notifications []Notifications
}

type UserProfile struct {
	User          User
	CreatedPosts  []Post
	Hashtags      []Hashtag
	Notifications []Notifications
	Messages []Message
}


type Hashtags struct {
	Hashtags []Hashtag
}

type Post struct {
	PostID    int
	Username  string `json:"username"`
	Content   string `json:"postBody"`
	Hashtag   string
	CreatedAt time.Time
}

type Message struct {
	MessageID    int       `json:"message_id"`
	Sender       string    `json:"sender"`
	Recipient    string    `json:"recipient"`
	Message      string    `json:"message"`
	CreationDate time.Time `json:"creation_date"`
}

type Hashtag struct {
	ID    int
	Name  string `json:"name"`
	Count string `json:"count"`
}

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	LoggedIn string
}

type UserSession struct {
	username string
	userID   int
	session  string
	max_age  int
}

type Comment struct {
	CommentID int
	PostID    int
	Username  string `json:"username"`
	Content   string `json:"commentBody"`
	CreatedAt time.Time
}

// maybe the fields can be updated?
type Chat struct {
	MessageSender    string    `json:"messagesender"`
	MessageRecipient string    `json:"messagerecipient"`
	SenderID         string    `json:"SenderID"`
	Message          string    `json:"message"`
	ChatID           int       `json:"chatID"`
	MessageID        int       `json:"messageID"`
	CreatedAt        time.Time //`json:"chatDate"`
	UserWithHistroy  []Chat    `json:"userwithhistory"`
	User             []User    `json:"users"`
	Notifications    []Notifications
}

type CookieValue struct {
	CookieValue string
}

type LoadingMessage struct {
	SendersUsername    string `json:"sendersusername"`
	RecipientsUsername string `json:"recipientsusername"`
}

type Notifications struct {
	Sender       string `json:"sendernotification"`
	Recipient    string `json:"recipientnotification"`
	Notification int    `json:"noti"`
}
