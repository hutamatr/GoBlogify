package go_blog

import "time"

type ArticleResponse struct {
	Id         int       `json:"id"`
	Title      string    `json:"title"`
	Body       string    `json:"body"`
	Author     string    `json:"author"`
	Published  bool      `json:"published"`
	Deleted    bool      `json:"deleted"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
	Deleted_At time.Time `json:"deleted_at"`
	Category   Category  `json:"category"`
}
