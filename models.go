package main

import (
	"time"
)

// User represent a customer in the database with a backing table of users
type User struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	Profile        Profile   `json:"profile"`
}

// Profile belongs to a user
type Profile struct {
	UserID  int       `json:"user_id"`
	Avatar  string    `json:"avatar"`
	Created time.Time `json:"created"`
}

// Post represents a submitted link shown on the home page
// type Post struct {
// 	ID        int
// 	Title     string
// 	URL       string
// 	Host      string
// 	UserName  string
// 	Votes     int
// 	Comments  int
// 	CreatedAt time.Time
// }

// GetVoteCountsHuman renders the vote count for display
// func (p Post) GetVoteCountsHuman() string {
// 	return fmt.Sprintf("%d points", p.Votes)
// }

// // GetCommentCountsHuman renders the comment count for display
// func (p Post) GetCommentCountsHuman() string {
// 	return fmt.Sprintf("%d comments", p.Comments)
// }

// // CreatedAtHuman renders the post's age for display
// func (p Post) CreatedAtHuman() string {
// 	return p.CreatedAt.Format("Jan 2, 2006")
// }

// Metadata describes pagination state for a listing page
// type Metadata struct {
// 	TotalRecords int
// 	PageSize     int
// 	CurrentPage  int
// 	PrevPage     int
// 	NextPage     int
// }
