package go_blog

type Article struct {
	Id         int
	Title      string
	Body       string
	Author     string
	Created_At string
	Updated_At string
	Deleted_At string
	Deleted    bool
	Published  bool
	Category   string
}

type Category struct {
	Id   int
	Name string
}
