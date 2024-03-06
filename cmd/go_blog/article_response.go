package go_blog

type ArticleResponse struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Author    string `json:"author"`
	Published bool   `json:"published"`
	Category  string `json:"category"`
}
