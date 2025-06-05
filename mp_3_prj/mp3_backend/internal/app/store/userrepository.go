package store

import "backend/internal/app/model"

type UserRepository struct {
	store *Store
}

// Create...
func (r *UserRepository) Create(u *model.User) (*model.User, error) {
	query := `SELECT create_user($1, $2, $3)` //Внутрення функция в бд для добавления новогго юзера
	err := r.store.db.QueryRow(query, u.Username, u.Email, u.PasswordHash).Scan(&u.ID)

	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	query := `SELECT id, username, email, password_hash FROM users WHERE email = $1`

	err := r.store.db.QueryRow(query, email).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.PasswordHash,
	)

	if err != nil {
		return nil, err
	}

	return u, nil
}
