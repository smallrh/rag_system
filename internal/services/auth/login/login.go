package login

import (
	"context"
	"encoding/json"
	"strings"
	"time"
)

func (s *Service) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	loginType := strings.ToLower(strings.TrimSpace(req.Type))
	if loginType != LoginTypePassword && loginType != LoginTypeCode {
		s.record(ctx, req, nil, false, "unsupported_type")
		return LoginResponse{}, ErrUnsupportedType
	}

	provider := strings.ToLower(strings.TrimSpace(req.Provider))
	username := strings.TrimSpace(req.Username)
	if provider == "" || username == "" {
		s.record(ctx, req, nil, false, "empty_input")
		return LoginResponse{}, ErrInvalidCredentials
	}

	// 1. 查身份
	identity, err := s.repo.Identity.GetByProviderIdentifier(ctx, provider, username)
	if err != nil || identity == nil {
		s.record(ctx, req, nil, false, "identity_not_found")
		return LoginResponse{}, ErrInvalidCredentials
	}
	if !isUnset(identity.DeletedAt) || identity.Status != "active" {
		s.record(ctx, req, identity, false, "identity_inactive")
		return LoginResponse{}, ErrIdentityInactive
	}

	// 2. 查用户
	user, err := s.repo.User.GetActiveByID(ctx, identity.UserID)
	if err != nil || user == nil {
		s.record(ctx, req, identity, false, "user_inactive")
		return LoginResponse{}, ErrUserNotActive
	}

	// 3. 按类型校验凭据
	switch loginType {
	case LoginTypePassword:
		if err := s.checkPassword(ctx, identity.UserID, req.Password); err != nil {
			s.record(ctx, req, identity, false, "password_mismatch")
			return LoginResponse{}, ErrInvalidCredentials
		}
	case LoginTypeCode:
		if err := s.consumeAuthCode(ctx, identity, req.AuthCode); err != nil {
			s.record(ctx, req, identity, false, "code_mismatch")
			return LoginResponse{}, ErrInvalidAuthCode
		}
	}

	// 4. 先签 access token（无 DB 副作用）
	now := time.Now()
	expiresAt := now.Add(AccessTokenTTL)
	accessToken, err := s.issueAccessToken(user, identity, now, expiresAt)
	if err != nil {
		s.record(ctx, req, identity, false, "access_token_issue_failed")
		return LoginResponse{}, err
	}

	// 5. 生成 refresh token（随机值）→ 只存哈希
	refreshToken, err := randomToken(refreshTokenN)
	if err != nil {
		return LoginResponse{}, err
	}
	refreshHash := sha256Hex(refreshToken)

	// 6. DeviceInfo → json.RawMessage
	var deviceInfo json.RawMessage
	if len(req.DeviceInfo) > 0 {
		if raw, err := json.Marshal(req.DeviceInfo); err == nil {
			deviceInfo = raw
		}
	}

	// 7. 创建 session
	if _, err := s.repo.Session.Create(
		ctx,
		identity.UserID,
		identity.ID,
		identity.Provider,
		refreshHash,
		deviceInfo,
		time.Now().Add(RefreshTokenTTL),
	); err != nil {
		s.record(ctx, req, identity, false, "session_create_failed")
		return LoginResponse{}, err
	}

	// 8. 审计成功
	s.record(ctx, req, identity, true, "")

	return LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		UserID:       identity.UserID.String(),
	}, nil
}
