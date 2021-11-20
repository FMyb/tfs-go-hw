package domain

type Message struct {
	ID       uint `json:"id"`
	UserID   uint `json:"user_id"`
	ToUserID uint `json:"to_user_id"`
	Text     string
}
