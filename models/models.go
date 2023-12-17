package models

import "database/sql"

type RiddleBase struct {
	ID       int     `json:"id"`
	Riddle   string  `json:"riddle"`
	Solution string  `json:"solution"`
	Synonyms *string `json:"synonyms,omitempty"`
}

type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

type Riddle struct {
	RiddleBase
	Username  sql.NullString `json:"username,omitempty"`
	UserEmail sql.NullString `json:"user_email,omitempty"`
}

type RiddleResponse struct {
	RiddleBase
	Links []Link `json:"links,omitempty"`
}
