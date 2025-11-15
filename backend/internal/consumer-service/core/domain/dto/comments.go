package dto

type Comment struct {
	CommentID    string `db:"comment_id" json:"comment_id"`
	UserID       string `db:"user_id" json:"user_id"`
	Content      string `db:"content" json:"content"`
	Rating       int    `db:"rating" json:"rating"`
	HelpfulVotes int    `db:"helpful_votes" json:"helpful_votes"`
	IsToxic      bool   `db:"is_toxic" json:"is_toxic"`
}

type AddCommentRequest struct {
	UserID  string `json:"user_id"`
	Content string `json:"content"`
	Rating  int    `json:"rating"`
}

type StoreComment struct {
	CommentID    string `db:"comment_id" json:"comment_id"`
	UserID       string `db:"user_id" json:"user_id"`
	Content      string `db:"content" json:"content"`
	Rating       int    `db:"rating" json:"rating"`
	HelpfulVotes int    `db:"helpful_votes" json:"helpful_votes"`
	IsToxic      bool   `db:"is_toxic" json:"is_toxic"`
}
