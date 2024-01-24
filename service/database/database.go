/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {

	// Default methods
	GetName() (string, error)
	SetName(name string) error
	Ping() error

	// User queries
	GetUsernameByToken(Id string) (*string, error)
	GetOwnerUsernameOfComment(CommentID string) (*string, error)
	PostUserID(Username string) (*components.User, error)
	SearchUser(Username string) (*components.UserList, error)
	UpdateUsername(NewUsername string, OldUsername string) error

	// Post queries
	CheckIfOwnerPost(Username string, PostID string) error
	AddLikeToPost(Username string, PostID string) error
	RemoveLikeFromPost(Username string, PostID string) error
	AddCommentToPost(PostID string, Body string, Author string) error
	RemoveCommentFromPost(PostID string, CommentID string) error
	GetUserStream(username string) (*components.Stream, error)
	UploadPost(username string, description string) (*components.Post, error)
	DeletePost(postID string) (*string, error)
	GetPostComments(postID string) (*components.CommentList, error)
	GetPostLikes(postID string) (*components.UserList, error)

	// Profile queries
	GetUserProfile(Username string) (*components.Profile, error)

	// Follow queries
	GetFollowingList(followingUsername string) (*components.UserList, error)
	GetFollowersList(followedUsername string) (*components.UserList, error)
	FollowUser(followerUsername string, followingUsername string) error
	UnfollowUser(followerUsername string, followingUsername string) error

	// Ban queries
	BanUser(bannerUsername, bannedUsername string) error
	UnbanUser(bannerUsername, bannedUsername string) error
	GetBanUserList(bannerUsername string) (*components.UserList, error)
	CheckIfBanned(bannerUsername string, bannedUsername string) error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE example_table (id INTEGER NOT NULL PRIMARY KEY, name TEXT);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS User (
		ID STRING UNIQUE NOT NULL,
		Username STRING PRIMARY KEY NOT NULL,
		ProfilePicPath STRING DEFAULT 'profile_pics/default.png',
		Birthdate STRING,
		Name STRING
	);
	CREATE TABLE IF NOT EXISTS Post (
		PostID INTEGER PRIMARY KEY AUTOINCREMENT,
		Author VARCHAR(16) NOT NULL,
		CreationDatetime STRING NOT NULL,
		Description VARCHAR(128),
		PhotoPath STRING, 
		FOREIGN KEY (Author) REFERENCES User(Username) ON DELETE CASCADE ON UPDATE CASCADE
	);
	CREATE TABLE IF NOT EXISTS Like (
		PostID INTEGER NOT NULL,
		Liker STRING NOT NULL,
		CreationDatetime STRING NOT NULL,
		PRIMARY KEY (PostID, Liker),
		FOREIGN KEY (PostID) REFERENCES Post(PostID), 
		FOREIGN KEY (Liker) REFERENCES User(Username) ON DELETE CASCADE ON UPDATE CASCADE
	);
	CREATE TABLE IF NOT EXISTS Follow (
		Follower STRING NOT NULL,
		Followed STRING NOT NULL,
		CreationDatetime STRING NOT NULL,
		PRIMARY KEY (Follower, Followed),
		FOREIGN KEY (Follower) REFERENCES User(Username) ON DELETE CASCADE ON UPDATE CASCADE,
		FOREIGN KEY (Followed) REFERENCES User(Username) ON DELETE CASCADE ON UPDATE CASCADE
	);
	CREATE TABLE IF NOT EXISTS Comment (
		CommentID INTEGER PRIMARY KEY AUTOINCREMENT,
		PostID INTEGER NOT NULL,
		Author STRING NOT NULL,
		CreationDatetime STRING NOT NULL,
		Comment STRING,
		FOREIGN KEY (Author) REFERENCES User(Username) ON DELETE CASCADE ON UPDATE CASCADE
		FOREIGN KEY (PostID) REFERENCES Post(PostID) ON DELETE CASCADE ON UPDATE CASCADE
	);
	CREATE TABLE IF NOT EXISTS Ban (
		Banner STRING,
		Banned STRING,
		CreationDatetime STRING NOT NULL,
		PRIMARY KEY (Banner, Banned),
		FOREIGN KEY (Banned) REFERENCES User(Username) ON DELETE CASCADE ON UPDATE CASCADE,
		FOREIGN KEY (Banner) REFERENCES User(Username) ON DELETE CASCADE ON UPDATE CASCADE
	);`)
	if err != nil {
		return nil, err
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
