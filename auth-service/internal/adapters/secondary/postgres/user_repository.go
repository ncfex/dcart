package postgres

import (
	"database/sql"

	"github.com/ncfex/dcart/auth-service/internal/core/ports"
	"github.com/ncfex/dcart/auth-service/internal/domain"

	_ "github.com/lib/pq"
)

type repository struct {
	db *sql.DB
}

func NewUserRepository(dsn string) ports.UserRepository {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	return &repository{db: db}
}

func (r *repository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow("SELECT id, username, password FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) Create(user *domain.User) error {
	_, err := r.db.Exec("INSERT INTO users (id, username, password) VALUES (get_random_uuid(), $1, $2)", user.Username, user.Password)
	return err
}
