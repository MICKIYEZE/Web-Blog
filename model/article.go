package model

type Article struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Published string `json:"published"`
	Author    string `json:"author"`
}