package repository

import (
	"context"
	"errors"

	"aurora-adminui/internal/domain/entity"
	domainrepo "aurora-adminui/internal/domain/repository"
	"aurora-adminui/internal/errorx"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MarketplaceRepoImple struct {
	db *pgxpool.Pool
}

func NewMarketplaceRepo(db *pgxpool.Pool) domainrepo.MarketplaceRepository {
	return &MarketplaceRepoImple{db: db}
}

func (r *MarketplaceRepoImple) ListMarketplaceModelOptions(ctx context.Context) ([]entity.MarketplaceApp, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT
		    (array_agg(id ORDER BY resource_version DESC, id ASC))[1]::text AS id,
		    (array_agg(id ORDER BY resource_version DESC, id ASC))[1]::text AS resource_definition_id,
		    resource_type,
		    COALESCE(NULLIF(model_name, ''), resource_model) AS resource_model,
		    COALESCE(array_remove(array_agg(resource_version ORDER BY resource_version DESC), NULL), '{}') AS versions
		  FROM resource_platform.resource_definitions
		  WHERE resource_type = 'application'
		  GROUP BY resource_type, COALESCE(NULLIF(model_name, ''), resource_model)
		  ORDER BY resource_type ASC, COALESCE(NULLIF(model_name, ''), resource_model) ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entity.MarketplaceApp, 0)
	for rows.Next() {
		var item entity.MarketplaceApp
		if err := rows.Scan(
			&item.ID,
			&item.ResourceDefinitionID,
			&item.ResourceType,
			&item.ResourceModel,
			&item.Versions,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *MarketplaceRepoImple) ListMarketplaceTemplateOptions(ctx context.Context) ([]entity.MarketplaceTemplate, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT
		    template.id::text,
		    template.name,
		    definition.resource_type,
		    COALESCE(NULLIF(definition.model_name, ''), definition.resource_model) AS resource_model,
		    definition.resource_version
		  FROM resource_platform.resource_templates template
		  JOIN resource_platform.resource_definitions definition
		    ON definition.id = template.resource_definition_id
		  WHERE definition.resource_type = 'application'
		  ORDER BY COALESCE(NULLIF(definition.model_name, ''), definition.resource_model) ASC, definition.resource_version DESC, template.name ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entity.MarketplaceTemplate, 0)
	for rows.Next() {
		var item entity.MarketplaceTemplate
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.ResourceType,
			&item.ResourceModel,
			&item.Version,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *MarketplaceRepoImple) ListMarketplaceApps(ctx context.Context) ([]entity.MarketplaceApp, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT
		    app.id,
		    app.name,
		    app.slug,
		    app.summary,
		    app.resource_definition_id::text,
		    app.template_id::text,
		    template.name,
		    linked_def.resource_type,
		    COALESCE(NULLIF(linked_def.model_name, ''), linked_def.resource_model) AS resource_model,
		    COALESCE(def_versions.versions, '{}') AS versions
		  FROM workspace.marketplace_apps app
		  JOIN resource_platform.resource_definitions linked_def
		    ON linked_def.id = app.resource_definition_id
		  JOIN resource_platform.resource_templates template
		    ON template.id = app.template_id
		  LEFT JOIN LATERAL (
		    SELECT array_agg(version ORDER BY version DESC) AS versions
		    FROM (
		      SELECT DISTINCT version_def.resource_version AS version
		      FROM resource_platform.resource_definitions version_def
		      WHERE version_def.resource_type = linked_def.resource_type
		        AND COALESCE(NULLIF(version_def.model_name, ''), version_def.resource_model) = COALESCE(NULLIF(linked_def.model_name, ''), linked_def.resource_model)
		    ) versions
		  ) def_versions ON TRUE
		  ORDER BY app.created_at DESC, app.name ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entity.MarketplaceApp, 0)
	for rows.Next() {
		var item entity.MarketplaceApp
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Slug,
			&item.Summary,
			&item.ResourceDefinitionID,
			&item.TemplateID,
			&item.TemplateName,
			&item.ResourceType,
			&item.ResourceModel,
			&item.Versions,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *MarketplaceRepoImple) GetMarketplaceApp(ctx context.Context, id string) (*entity.MarketplaceApp, error) {
	row := r.db.QueryRow(
		ctx,
		`SELECT
		    app.id,
		    app.name,
		    app.slug,
		    app.summary,
		    app.description,
		    app.resource_definition_id::text,
		    app.template_id::text,
		    template.name,
		    linked_def.resource_type,
		    COALESCE(NULLIF(linked_def.model_name, ''), linked_def.resource_model) AS resource_model,
		    COALESCE(def_versions.versions, '{}') AS versions
		  FROM workspace.marketplace_apps app
		  JOIN resource_platform.resource_definitions linked_def
		    ON linked_def.id = app.resource_definition_id
		  JOIN resource_platform.resource_templates template
		    ON template.id = app.template_id
		  LEFT JOIN LATERAL (
		    SELECT array_agg(version ORDER BY version DESC) AS versions
		    FROM (
		      SELECT DISTINCT version_def.resource_version AS version
		      FROM resource_platform.resource_definitions version_def
		      WHERE version_def.resource_type = linked_def.resource_type
		        AND COALESCE(NULLIF(version_def.model_name, ''), version_def.resource_model) = COALESCE(NULLIF(linked_def.model_name, ''), linked_def.resource_model)
		    ) versions
		  ) def_versions ON TRUE
		  WHERE app.id = $1`,
		id,
	)

	var item entity.MarketplaceApp
	err := row.Scan(
		&item.ID,
		&item.Name,
		&item.Slug,
		&item.Summary,
		&item.Description,
		&item.ResourceDefinitionID,
		&item.TemplateID,
		&item.TemplateName,
		&item.ResourceType,
		&item.ResourceModel,
		&item.Versions,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *MarketplaceRepoImple) CreateMarketplaceApp(ctx context.Context, item *entity.MarketplaceApp) error {
	if item == nil {
		return errors.New("marketplace app is nil")
	}
	resourceDefinitionID, err := uuid.Parse(item.ResourceDefinitionID)
	if err != nil {
		return errorx.ErrInvalidArgument
	}
	templateID, err := uuid.Parse(item.TemplateID)
	if err != nil {
		return errorx.ErrInvalidArgument
	}

	_, err = r.db.Exec(
		ctx,
		`INSERT INTO workspace.marketplace_apps (
		    id, name, slug, summary, description, resource_definition_id, template_id, created_at, updated_at
		  ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		item.ID,
		item.Name,
		item.Slug,
		item.Summary,
		item.Description,
		resourceDefinitionID,
		templateID,
		item.CreatedAt,
		item.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return errorx.ErrMarketplaceAppAlreadyExists
		}
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return errorx.ErrInvalidArgument
		}
		return err
	}
	return nil
}
