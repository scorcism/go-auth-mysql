package secrets

import (
	"database/sql"

	"github.com/scorcism/go-auth/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetSecrets(userId int) ([]types.Secret, error) {
	rows, err := s.db.Query("SELECT * FROM secrets where user_id=?", userId)
	if err != nil {
		return nil, err
	}

	secrets := make([]types.Secret, 0)

	for rows.Next() {
		p, err := scaneRowInSecret(rows)
		if err != nil {
			return nil, err
		}

		secrets = append(secrets, *p)
	}

	return secrets, nil
}

func (s *Store) AddSecret(secret types.Secret) error {

	_, err := s.db.Exec("INSERT INTO secrets (secret_key, label, user_id) VALUES (?,?,?)", secret.SecretKey, secret.Label, secret.UserId)

	if err != nil {
		return err
	}

	return nil
}

func scaneRowInSecret(rows *sql.Rows) (*types.Secret, error) {
	secret := new(types.Secret)

	err := rows.Scan(
		&secret.ID,
		&secret.Label,
		&secret.SecretKey,
		&secret.UserId,
		&secret.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return secret, nil
}
