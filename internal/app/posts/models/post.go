package models

type Post struct {
	ID       int    `json:"id"`
	Parent   *int   `json:"-"`
	Author   string `json:"author"`
	Message  string `json:"message"`
	ISEdited bool   `json:"is_edited,omitempty"`
	Forum    string `json:"forum"`
	Thread   int    `json:"thread"`
	Created  string `json:"created"`
}
