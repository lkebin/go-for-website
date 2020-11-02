package model

type Post struct {
	ID      int32  `db:"id,unsafe"`
	Title   string `db:"title"`
	Content string `db:"content"`
}
