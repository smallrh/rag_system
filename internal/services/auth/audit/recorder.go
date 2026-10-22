package audit

import (
	"context"
	"strings"

	"github.com/google/uuid"

	authModel "go_study/internal/model/auth"
	repoAuth "go_study/internal/repository/auth"
)

type Event struct {
	EventType  string
	Success    bool
	Identity   *authModel.AuthIdentity
	Identifier string
	IPAddress  string
	UserAgent  string
	Reason     string
}

type Recorder struct {
	repo *repoAuth.Repository
}

func NewRecorder(repo *repoAuth.Repository) *Recorder {
	return &Recorder{repo: repo}
}

func (r *Recorder) Record(ctx context.Context, e Event) {
	var uid, iid *uuid.UUID
	if e.Identity != nil {
		id := e.Identity.ID
		u := e.Identity.UserID
		iid = &id
		uid = &u
	}

	identifier := e.Identifier
	identifierPtr := &identifier

	var ipPtr *string
	if ip := strings.TrimSpace(e.IPAddress); ip != "" {
		ipPtr = &ip
	}

	var uaPtr *string
	if ua := strings.TrimSpace(e.UserAgent); ua != "" {
		uaPtr = &ua
	}

	var reasonPtr *string
	if e.Reason != "" {
		reason := e.Reason
		reasonPtr = &reason
	}

	_, _ = r.repo.AuditLog.Create(
		ctx, uid, iid, e.EventType, e.Success, identifierPtr, ipPtr, uaPtr, reasonPtr,
	)
}
