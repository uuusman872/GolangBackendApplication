package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golangcollege/sessions"
	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	errorLog    *log.Logger
	infoLog     *log.Logger
	userRerpo   UserRepository
	templateDir string
	tp          *TemplateRenderer
	publicPath  string
	session     sessions.Session
	postRepo    PostRepository
}

var session *sessions.Session
var secret = []byte("u46IpCV9y5Vlur8YvODJEhgOY8m9JVE4")

var user_schema = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		hashed_password BLOB NOT NULL, -- Storing as BLOB for byte slice
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
`

var profile_schema = `
	CREATE TABLE IF NOT EXISTS profile (
		user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
		avatar TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
`

var posts_schema = `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT NOT NULL,
		title TEXT NOT NULL UNIQUE,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
`

var comments_schema = `
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		body TEXT NOT NULL,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
`

var votes_table = `
	CREATE TABLE IF NOT EXISTS votes (
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (user_id, post_id)
	);
`

func main() {
	db, err := connectToDatabase("users_database.db")
	if err != nil {
		log.Fatal(err)
	}
	session = sessions.New(secret)
	session.Lifetime = 3 * time.Hour
	session.Secure = true
	session.SameSite = http.SameSiteLaxMode

	app := &application{
		errorLog:    log.New(os.Stderr, "ERROR\t", log.LstdFlags|log.Lmicroseconds|log.Lshortfile),
		infoLog:     log.New(os.Stderr, "INFO\t", log.Ltime|log.LstdFlags),
		userRerpo:   NewSQLUserRepository(db),
		templateDir: "./templates",
		publicPath:  "./public",
		session:     *session,
		postRepo:    NewSQLPostRepository(db),
	}

	if err := app.userRerpo.CreateDatabase(posts_schema); err != nil {
		panic("posts Database not created")
	}
	if err := app.userRerpo.CreateDatabase(comments_schema); err != nil {
		panic("comments Database not created")
	}
	if err := app.userRerpo.CreateDatabase(votes_table); err != nil {
		panic("votes Database not created")
	}

	app.tp = NewTemplateRenderer(app.templateDir, true)

	log.Println("Listening on :8080")
	if err := app.serve(); err != nil {
		log.Fatal(err)
	}

	// id, err := app.postRepo.AddComment(1, 1, "lovely comment from us")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(id)

	// comments, err := app.postRepo.GetComments(1)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(comments)

	// if err := app.postRepo.AddVote(1, 1); err != nil {
	// 	panic(err)
	// }

	// posts, _, err := app.postRepo.GetAll(Filter{
	// 	PageSize: 10,
	// 	Page:     1,
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(posts)

	// id, err := app.postRepo.CreatePost("welcome Post", "url", 1)
	// if err != nil {
	// 	panic(err)
	// }

}

func connectToDatabase(name string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
