package login

import (
	"context"

	authModel "go_study/internal/model/auth"
	"go_study/internal/services/auth/audit"
)

func (s *Service) record(
	ctx context.Context,
	req LoginRequest,
	identity *authModel.AuthIdentity,
	success bool,
	reason string,
) {
	s.audit.Record(ctx, audit.Event{
		EventType:  "login",
		Success:    success,
		Identity:   identity,
		Identifier: req.Username,
		IPAddress:  req.IPAddress,
		UserAgent:  req.UserAgent,
		Reason:     reason,
	})
}
