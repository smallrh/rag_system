package login

import (
	repoAuth "go_study/internal/repository/auth"
	"go_study/internal/services/auth/audit"
)

type Service struct {
	repo      *repoAuth.Repository
	jwtSecret []byte
	audit     *audit.Recorder
}

func NewService(repo *repoAuth.Repository, jwtSecret []byte) *Service {
	return &Service{
		repo:      repo,
		jwtSecret: jwtSecret,
		audit:     audit.NewRecorder(repo),
	}
}
