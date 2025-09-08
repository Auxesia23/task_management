-- +goose Up
-- +goose StatementBegin
CREATE TABLE invitations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL,
    user_id UUID NOT NULL,
    inviter_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT fk_project
        FOREIGN KEY(project_id)
        REFERENCES projects(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_inviter
        FOREIGN KEY(inviter_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT unique_invitation UNIQUE (project_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE invitations DROP CONSTRAINT fk_project;
ALTER TABLE invitations DROP CONSTRAINT fk_user;
ALTER TABLE invitations DROP CONSTRAINT fk_inviter;
DROP TABLE invitations;
-- +goose StatementEnd
