package specification

import (
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	sq "github.com/Masterminds/squirrel"
)

type ApiKeySpecification struct {
	WorkspaceIDs []int64
	Environments []string
}

func NewApiKeySpecification(spec *ApiKeySpecification) (string, []interface{}, error) {
	builder := sq.
		Select("id", "name", "workspace_id", "environment", "created_at", "updated_at", "revoked", "expires_at", "raw_prefix").
		From("api_keys")

	if spec != nil {
		if len(spec.WorkspaceIDs) > 0 {
			builder = builder.Where(sq.Eq{"workspace_id": spec.WorkspaceIDs})
		}
		if len(spec.Environments) > 0 {
			builder = builder.Where(sq.Eq{"environment": spec.Environments})
		}
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return "", nil, err
	}

	return query, args, nil

}
func ToApiKeySpecification(filter *request.GetApiKeyRequestFilter) *ApiKeySpecification {
	if filter == nil {
		return nil
	}
	return &ApiKeySpecification{
		WorkspaceIDs: filter.WorkspaceIDs,
		Environments: filter.Environments,
	}
}
