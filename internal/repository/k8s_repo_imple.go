package repository

import (
	"context"
	"errors"
	"strings"

	"aurora-adminui/internal/domain/entity"
	domainrepo "aurora-adminui/internal/domain/repository"
	"aurora-adminui/internal/errorx"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const k8sZoneObjectType = "k8s_cluster"

type K8sRepoImple struct {
	db *pgxpool.Pool
}

func NewK8sRepo(db *pgxpool.Pool) domainrepo.K8sRepository {
	return &K8sRepoImple{db: db}
}

func (r *K8sRepoImple) ListClusters(ctx context.Context) ([]entity.K8sCluster, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT
		    c.id,
		    c.name,
		    c.description,
		    c.import_mode,
		    c.api_server_url,
		    c.kubernetes_version,
		    c.validation_status,
		    c.last_validated_at,
		    COALESCE(c.last_validation_error, ''),
		    c.supports_dbaas,
		    c.supports_serverless,
		    c.supports_generic_workloads,
		    c.created_at,
		    z.id,
		    COALESCE(z.name, '')
		  FROM platform.k8s_clusters c
		  LEFT JOIN zone.zone_objects zo
		    ON zo.object_type = $1
		   AND zo.object_id = c.id::text
		  LEFT JOIN zone.zones z ON z.id = zo.zone_id
		  ORDER BY c.created_at DESC, c.name ASC`,
		k8sZoneObjectType,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entity.K8sCluster, 0)
	for rows.Next() {
		var (
			item   entity.K8sCluster
			zoneID *uuid.UUID
		)
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.ImportMode,
			&item.APIServerURL,
			&item.KubernetesVersion,
			&item.ValidationStatus,
			&item.LastValidatedAt,
			&item.LastValidationError,
			&item.SupportsDBAAS,
			&item.SupportsServerless,
			&item.SupportsGenericWorkloads,
			&item.CreatedAt,
			&zoneID,
			&item.ZoneName,
		); err != nil {
			return nil, err
		}
		item.ZoneID = zoneID
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *K8sRepoImple) GetClusterByID(ctx context.Context, id string) (*entity.K8sCluster, error) {
	var (
		item   entity.K8sCluster
		zoneID *uuid.UUID
	)
	err := r.db.QueryRow(
		ctx,
		`SELECT
		    c.id,
		    c.name,
		    c.description,
		    c.import_mode,
		    c.kubeconfig_ciphertext,
		    c.api_server_url,
		    c.current_context,
		    c.kubernetes_version,
		    c.validation_status,
		    c.last_validated_at,
		    COALESCE(c.last_validation_error, ''),
		    c.supports_dbaas,
		    c.supports_serverless,
		    c.supports_generic_workloads,
		    c.created_at,
		    z.id,
		    COALESCE(z.name, '')
		  FROM platform.k8s_clusters c
		  LEFT JOIN zone.zone_objects zo
		    ON zo.object_type = $2
		   AND zo.object_id = c.id::text
		  LEFT JOIN zone.zones z ON z.id = zo.zone_id
		  WHERE c.id = $1`,
		id,
		k8sZoneObjectType,
	).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.ImportMode,
		&item.KubeconfigCiphertext,
		&item.APIServerURL,
		&item.CurrentContext,
		&item.KubernetesVersion,
		&item.ValidationStatus,
		&item.LastValidatedAt,
		&item.LastValidationError,
		&item.SupportsDBAAS,
		&item.SupportsServerless,
		&item.SupportsGenericWorkloads,
		&item.CreatedAt,
		&zoneID,
		&item.ZoneName,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errorx.ErrK8sClusterNotFound
	}
	if err != nil {
		return nil, err
	}
	item.ZoneID = zoneID
	return &item, nil
}

func (r *K8sRepoImple) CreateCluster(ctx context.Context, cluster *entity.K8sCluster) error {
	if cluster == nil {
		return errors.New("k8s cluster is nil")
	}
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	_, err = tx.Exec(
		ctx,
		`INSERT INTO platform.k8s_clusters (
		    id,
		    name,
		    description,
		    import_mode,
		    kubeconfig_ciphertext,
		    api_server_url,
		    current_context,
		    kubernetes_version,
		    validation_status,
		    last_validated_at,
		    last_validation_error,
		    supports_dbaas,
		    supports_serverless,
		    supports_generic_workloads,
		    created_at
		  ) VALUES (
		    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		  )`,
		cluster.ID,
		cluster.Name,
		cluster.Description,
		cluster.ImportMode,
		cluster.KubeconfigCiphertext,
		cluster.APIServerURL,
		cluster.CurrentContext,
		cluster.KubernetesVersion,
		cluster.ValidationStatus,
		cluster.LastValidatedAt,
		cluster.LastValidationError,
		cluster.SupportsDBAAS,
		cluster.SupportsServerless,
		cluster.SupportsGenericWorkloads,
		cluster.CreatedAt,
	)
	if err != nil {
		return mapK8sMutationError(err)
	}

	if cluster.ZoneID != nil {
		_, err = tx.Exec(
			ctx,
			`INSERT INTO zone.zone_objects (object_type, object_id, zone_id, created_at)
			 VALUES ($1, $2, $3, $4)`,
			k8sZoneObjectType,
			cluster.ID.String(),
			*cluster.ZoneID,
			cluster.CreatedAt,
		)
		if err != nil {
			return mapK8sMutationError(err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (r *K8sRepoImple) UpdateClusterValidation(ctx context.Context, cluster *entity.K8sCluster) error {
	if cluster == nil {
		return errors.New("k8s cluster is nil")
	}
	tag, err := r.db.Exec(
		ctx,
		`UPDATE platform.k8s_clusters
		    SET api_server_url = $2,
		        current_context = $3,
		        kubernetes_version = $4,
		        validation_status = $5,
		        last_validated_at = $6,
		        last_validation_error = $7
		  WHERE id = $1`,
		cluster.ID,
		cluster.APIServerURL,
		cluster.CurrentContext,
		cluster.KubernetesVersion,
		cluster.ValidationStatus,
		cluster.LastValidatedAt,
		cluster.LastValidationError,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errorx.ErrK8sClusterNotFound
	}
	return nil
}

func (r *K8sRepoImple) DeleteCluster(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return errorx.ErrInvalidArgument
	}
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	_, err = tx.Exec(
		ctx,
		`DELETE FROM zone.zone_objects
		  WHERE object_type = $1
		    AND object_id = $2`,
		k8sZoneObjectType,
		id,
	)
	if err != nil {
		return err
	}

	tag, err := tx.Exec(ctx, `DELETE FROM platform.k8s_clusters WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errorx.ErrK8sClusterNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func mapK8sMutationError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return errorx.ErrK8sClusterAlreadyExists
		case "23503":
			return errorx.ErrZoneNotFound
		}
	}
	return err
}
