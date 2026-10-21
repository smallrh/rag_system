
-- ==========================================
-- 1. 用户主体表
-- ==========================================
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       display_name VARCHAR(100),
                       status VARCHAR(20) NOT NULL DEFAULT 'active',
                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       deleted_at TIMESTAMPTZ,

                       CONSTRAINT chk_users_status
                           CHECK (status IN ('active', 'disabled', 'locked'))
);

-- ==========================================
-- 2. 认证身份表
-- ==========================================
CREATE TABLE auth_identities (
                                 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                 user_id UUID NOT NULL REFERENCES users(id),
                                 provider VARCHAR(30) NOT NULL,
                                 identifier VARCHAR(320) NOT NULL,
                                 verified_at TIMESTAMPTZ,
                                 status VARCHAR(20) NOT NULL DEFAULT 'active',
                                 metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
                                 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                 updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                 deleted_at TIMESTAMPTZ,

                                 CONSTRAINT chk_identity_provider
                                     CHECK (provider IN ('email', 'phone', 'google', 'github', 'apple')),

                                 CONSTRAINT chk_identity_status
                                     CHECK (status IN ('active', 'disabled', 'revoked'))
);

-- 同一有效身份不能重复绑定
CREATE UNIQUE INDEX uq_auth_identity_identifier
    ON auth_identities(provider, LOWER(identifier))
    WHERE deleted_at IS NULL;

CREATE INDEX idx_auth_identities_user_id
    ON auth_identities(user_id);

-- ==========================================
-- 3. 认证凭证表
-- ==========================================
CREATE TABLE auth_credentials (
                                  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                  user_id UUID NOT NULL REFERENCES users(id),
                                  identity_id UUID REFERENCES auth_identities(id),
                                  credential_type VARCHAR(30) NOT NULL,
                                  secret_hash TEXT,
                                  metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
                                  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                  deleted_at TIMESTAMPTZ,

                                  CONSTRAINT chk_credential_type
                                      CHECK (credential_type IN ('password', 'passkey')),

                                  CONSTRAINT chk_credential_secret
                                      CHECK (
                                          (credential_type = 'password' AND secret_hash IS NOT NULL)
                                              OR
                                          (credential_type = 'passkey')
                                          )
);

-- 一个用户只能有一条有效密码凭证
CREATE UNIQUE INDEX uq_auth_credentials_password
    ON auth_credentials(user_id, credential_type)
    WHERE credential_type = 'password'
      AND deleted_at IS NULL;

CREATE INDEX idx_auth_credentials_user_id
    ON auth_credentials(user_id);

-- ==========================================
-- 4. 登录会话表
-- ==========================================
CREATE TABLE auth_sessions (
                               id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                               user_id UUID NOT NULL REFERENCES users(id),
                               identity_id UUID REFERENCES auth_identities(id),
                               provider VARCHAR(30),
                               refresh_token_hash VARCHAR(128) NOT NULL UNIQUE,
                               device_info JSONB NOT NULL DEFAULT '{}'::jsonb,
                               expires_at TIMESTAMPTZ NOT NULL,
                               last_seen_at TIMESTAMPTZ,
                               rotated_at TIMESTAMPTZ,
                               replaced_by_session_id UUID REFERENCES auth_sessions(id),
                               revoked_at TIMESTAMPTZ,
                               created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                               updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                               deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_auth_sessions_user_id
    ON auth_sessions(user_id);

CREATE INDEX idx_auth_sessions_active
    ON auth_sessions(user_id, expires_at)
    WHERE revoked_at IS NULL AND deleted_at IS NULL;

-- ==========================================
-- 5. 验证与重置令牌表
-- ==========================================
CREATE TABLE auth_tokens (
                             id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                             user_id UUID REFERENCES users(id),
                             identity_id UUID REFERENCES auth_identities(id),
                             purpose VARCHAR(30) NOT NULL,
                             target VARCHAR(320) NOT NULL,
                             token_hash VARCHAR(128) NOT NULL,
                             expires_at TIMESTAMPTZ NOT NULL,
                             consumed_at TIMESTAMPTZ,
                             created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                             updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                             deleted_at TIMESTAMPTZ,

                             CONSTRAINT chk_auth_token_purpose
                                 CHECK (
                                     purpose IN (
                                                 'email_verify',
                                                 'phone_verify',
                                                 'password_reset',
                                                 'email_login',
                                                 'phone_login'
                                         )
                                     ),

                             CONSTRAINT uq_auth_tokens_token_hash
                                 UNIQUE (token_hash)
);

CREATE INDEX idx_auth_tokens_target_purpose
    ON auth_tokens(target, purpose, expires_at);

-- ==========================================
-- 6. 认证审计日志表
-- ==========================================
CREATE TABLE auth_audit_logs (
                                 id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                                 user_id UUID REFERENCES users(id),
                                 identity_id UUID REFERENCES auth_identities(id),
                                 event_type VARCHAR(50) NOT NULL,
                                 success BOOLEAN NOT NULL,
                                 identifier VARCHAR(320),
                                 ip_address INET,
                                 user_agent TEXT,
                                 failure_reason VARCHAR(100),
                                 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_auth_audit_user_created
    ON auth_audit_logs(user_id, created_at DESC);

CREATE INDEX idx_auth_audit_identifier_created
    ON auth_audit_logs(identifier, created_at DESC);

-- ==========================================
-- 7. 自动更新 updated_at
-- ==========================================
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_auth_identities_updated_at
    BEFORE UPDATE ON auth_identities
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_auth_credentials_updated_at
    BEFORE UPDATE ON auth_credentials
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_auth_sessions_updated_at
    BEFORE UPDATE ON auth_sessions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_auth_tokens_updated_at
    BEFORE UPDATE ON auth_tokens
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
