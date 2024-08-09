package repository

import (
	"database/sql"
	"feather/config"
	"feather/types/schema"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

type Repository struct {
	config *config.Config
	db     *sql.DB
}

const (
	user             = "feather.user"
	githubUser       = "feather.github_user"
	jenkinsUser      = "feather.jenkins_user"
	githubRepository = "feather.github_repository"
)

func NewRepository(config *config.Config) (*Repository, error) {
	repository := &Repository{config: config}
	var err error

	if repository.db, err = sql.Open(config.DB.Database, config.DB.URL); err != nil {
		return nil, err
	}
	return repository, nil
}

func (repository *Repository) CreateUser(email string, password string) error {
	if _, err := repository.db.Exec("INSERT INTO feather.user(email, password) VALUES(?, ?)", email, password); err != nil {
		return err
	}
	log.Println("CreateUser Query run successfully!")
	return nil
}

func (repository *Repository) CreateGithubUser(userId int64, nickname string, email string, token string) error {
	if _, err := repository.db.Exec("INSERT INTO feather.github_user(user_id, nickname, email, token) VALUES(?, ?, ?, ?)", userId, nickname, email, token); err != nil {
		return err
	}
	log.Println("CreateGithubUser Query run successfully!")
	return nil
}

func (repository *Repository) GithubUser(githubUserId int64) (*schema.GithubUser, error) {
	u := new(schema.GithubUser)
	qs := query([]string{"SELECT * FROM", githubUser, "WHERE github_user_id = ?"})
	if err := repository.db.QueryRow(qs, githubUserId).Scan(&u.ID, &u.UserId, &u.Nickname, &u.Email, &u.Token); err != nil {
		if err := noResult(err); err != nil {
			return nil, err
		}
	}

	log.Println("GithubUser Query run successfully!")
	return u, nil
}

func (repository *Repository) CreateGithubRepository(githubUserId int64, name string, description string, isPrivate bool) error {
	if _, err := repository.db.Exec("INSERT INTO feather.github_repository(github_user_id, name, description, is_private) VALUES(?, ?, ?, ?)", githubUserId, name, description, isPrivate); err != nil {
		return err
	}

	log.Println("CreateGithubRepository Query run successfully!")
	return nil
}

func (repository *Repository) CreateJenkinsUser(userId int64, nickname string, token string) error {
	if _, err := repository.db.Exec("INSERT INTO feather.jenkins_user(user_id, nickname, token) VALUES(?, ?, ?)", userId, nickname, token); err != nil {
		return err
	}

	log.Println("CreateJenkinsUser Query run successfully!")
	return nil
}

func (repository *Repository) JenkinsUserByUserId(userId int64) (*schema.JenkinsUser, error) {
	u := new(schema.JenkinsUser)
	qs := query([]string{"SELECT * FROM", jenkinsUser, "WHERE user_id = ?"})
	if err := repository.db.QueryRow(qs, userId).Scan(&u.ID, &u.UserId, &u.Nickname, &u.Token); err != nil {
		if err := noResult(err); err != nil {
			return nil, err
		}
	}

	log.Println("JenkinsUserByUserId Query run successfully!")
	return u, nil
}

func (repository *Repository) JenkinsUserByJenkinsUserId(jenkinsUserId int64) (*schema.JenkinsUser, error) {
	u := new(schema.JenkinsUser)
	qs := query([]string{"SELECT * FROM", jenkinsUser, "WHERE jenkins_user_id = ?"})
	if err := repository.db.QueryRow(qs, jenkinsUserId).Scan(&u.ID, &u.UserId, &u.Nickname, &u.Token); err != nil {
		if err := noResult(err); err != nil {
			return nil, err
		}
	}

	log.Println("JenkinsUserByJenkinsUserId Query run successfully!")
	return u, nil
}

func query(qs []string) string {
	return strings.Join(qs, " ") + ";"
}

func noResult(err error) error {
	if strings.Contains(err.Error(), "sql: no rows in result set") {
		return nil
	} else {
		return err
	}
}
