package client

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	ID        int    `json:"id"`
	LoginTime string `json:"login_time"`
	Token     string `json:"token"`
}
