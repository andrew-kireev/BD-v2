package models

type Thread struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Forum   string `json:"forum"`
	Message string `json:"message"`
	Votes   int    `json:"votes,omitempty"`
	Slug    string `json:"slug"`
	Created string `json:"created"`
}

type Vote struct {
	ThreadID int    `json:"thread_id"`
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
}
