package auth

import (
	"context"

	authModel "go_study/internal/model/auth"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AuditLogRepository struct {
	db *sqlx.DB
}

func NewAuditLogRepository(db *sqlx.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Create(
	ctx context.Context,
	userID *uuid.UUID,
	identityID *uuid.UUID,
	eventType string,
	success bool,
	identifier *string,
	ipAddress *string,
	userAgent *string,
	failureReason *string,
) (*authModel.AuthAuditLog, error) {
	var log authModel.AuthAuditLog

	err := r.db.GetContext(ctx, &log, `
        INSERT INTO auth_audit_logs (
            user_id, identity_id, event_type, success,
            identifier, ip_address, user_agent, failure_reason
        )
        VALUES ($1, $2, $3, $4, $5, $6::inet, $7, $8)
        RETURNING
            id, user_id, identity_id, event_type, success,
            identifier, ip_address::text AS ip_address, user_agent, failure_reason, created_at
    `, userID, identityID, eventType, success, identifier, ipAddress, userAgent, failureReason)
	if err != nil {
		return nil, err
	}

	return &log, nil
}

func (r *AuditLogRepository) ListByUserID(
	ctx context.Context,
	userID uuid.UUID,
	limit int,
) ([]authModel.AuthAuditLog, error) {
	if limit <= 0 {
		limit = 50
	}

	var logs []authModel.AuthAuditLog

	err := r.db.SelectContext(ctx, &logs, `
        SELECT
            id, user_id, identity_id, event_type, success,
            identifier, ip_address::text AS ip_address, user_agent, failure_reason, created_at
        FROM auth_audit_logs
        WHERE user_id = $1
        ORDER BY created_at DESC
        LIMIT $2
    `, userID, limit)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (r *AuditLogRepository) ListByIdentifier(
	ctx context.Context,
	identifier string,
	limit int,
) ([]authModel.AuthAuditLog, error) {
	if limit <= 0 {
		limit = 50
	}

	var logs []authModel.AuthAuditLog

	err := r.db.SelectContext(ctx, &logs, `
        SELECT
            id, user_id, identity_id, event_type, success,
            identifier, ip_address::text AS ip_address, user_agent, failure_reason, created_at
        FROM auth_audit_logs
        WHERE identifier = $1
        ORDER BY created_at DESC
        LIMIT $2
    `, identifier, limit)
	if err != nil {
		return nil, err
	}

	return logs, nil
}
