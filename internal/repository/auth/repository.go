package auth

import "github.com/jmoiron/sqlx"

type Repository struct {
	db *sqlx.DB

	User *UserRepository
	//Identity   *IdentityRepository
	//Credential *CredentialRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db:   db,
		User: NewUserRepository(db),
		//Identity:   NewIdentityRepository(db),
		//Credential: NewCredentialRepository(db),
	}
}
