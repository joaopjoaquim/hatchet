// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: webhook_workers.sql

package dbsqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const deleteWebhookWorker = `-- name: DeleteWebhookWorker :exec
UPDATE "WebhookWorker"
SET "deleted" = true
WHERE
  "id" = $1::uuid
  and "tenantId" = $2::uuid
`

type DeleteWebhookWorkerParams struct {
	ID       pgtype.UUID `json:"id"`
	Tenantid pgtype.UUID `json:"tenantid"`
}

func (q *Queries) DeleteWebhookWorker(ctx context.Context, db DBTX, arg DeleteWebhookWorkerParams) error {
	_, err := db.Exec(ctx, deleteWebhookWorker, arg.ID, arg.Tenantid)
	return err
}

const listActiveWebhookWorkers = `-- name: ListActiveWebhookWorkers :many
SELECT id, "createdAt", "updatedAt", name, secret, url, "tokenValue", deleted, "tokenId", "tenantId"
FROM "WebhookWorker"
WHERE "tenantId" = $1::uuid AND "deleted" = false
`

func (q *Queries) ListActiveWebhookWorkers(ctx context.Context, db DBTX, tenantid pgtype.UUID) ([]*WebhookWorker, error) {
	rows, err := db.Query(ctx, listActiveWebhookWorkers, tenantid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*WebhookWorker
	for rows.Next() {
		var i WebhookWorker
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Secret,
			&i.Url,
			&i.TokenValue,
			&i.Deleted,
			&i.TokenId,
			&i.TenantId,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listWebhookWorkersByPartitionId = `-- name: ListWebhookWorkersByPartitionId :many
WITH tenants AS (
    SELECT
        "id"
    FROM
        "Tenant"
    WHERE
        "workerPartitionId" = $1::text OR
        "workerPartitionId" IS NULL
), update_partition AS (
    UPDATE
        "TenantWorkerPartition"
    SET
        "lastHeartbeat" = NOW()
    WHERE
        "id" = $1::text
)
SELECT id, "createdAt", "updatedAt", name, secret, url, "tokenValue", deleted, "tokenId", "tenantId"
FROM "WebhookWorker"
WHERE "tenantId" IN (SELECT "id" FROM tenants)
`

func (q *Queries) ListWebhookWorkersByPartitionId(ctx context.Context, db DBTX, workerpartitionid string) ([]*WebhookWorker, error) {
	rows, err := db.Query(ctx, listWebhookWorkersByPartitionId, workerpartitionid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*WebhookWorker
	for rows.Next() {
		var i WebhookWorker
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Secret,
			&i.Url,
			&i.TokenValue,
			&i.Deleted,
			&i.TokenId,
			&i.TenantId,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const upsertWebhookWorker = `-- name: UpsertWebhookWorker :one
INSERT INTO "WebhookWorker" (
    "id",
    "createdAt",
    "updatedAt",
    "name",
    "secret",
    "url",
    "tenantId",
    "tokenId",
    "tokenValue",
    "deleted"
)
VALUES (
    gen_random_uuid(),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    $1::text,
    $2::text,
    $3::text,
    $4::uuid,
    $5::uuid,
    $6::text,
    coalesce($7::boolean, false)
)
ON CONFLICT ("url") DO
UPDATE
SET
    "tokenId" = coalesce($5::uuid, excluded."tokenId"),
    "tokenValue" = coalesce($6::text, excluded."tokenValue"),
    "name" = coalesce($1::text, excluded."name"),
    "secret" = coalesce($2::text, excluded."secret"),
    "url" = coalesce($3::text, excluded."url"),
    "deleted" = coalesce($7::boolean, excluded."deleted")
RETURNING id, "createdAt", "updatedAt", name, secret, url, "tokenValue", deleted, "tokenId", "tenantId"
`

type UpsertWebhookWorkerParams struct {
	Name       pgtype.Text `json:"name"`
	Secret     pgtype.Text `json:"secret"`
	Url        pgtype.Text `json:"url"`
	Tenantid   pgtype.UUID `json:"tenantid"`
	TokenId    pgtype.UUID `json:"tokenId"`
	TokenValue pgtype.Text `json:"tokenValue"`
	Deleted    pgtype.Bool `json:"deleted"`
}

func (q *Queries) UpsertWebhookWorker(ctx context.Context, db DBTX, arg UpsertWebhookWorkerParams) (*WebhookWorker, error) {
	row := db.QueryRow(ctx, upsertWebhookWorker,
		arg.Name,
		arg.Secret,
		arg.Url,
		arg.Tenantid,
		arg.TokenId,
		arg.TokenValue,
		arg.Deleted,
	)
	var i WebhookWorker
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Secret,
		&i.Url,
		&i.TokenValue,
		&i.Deleted,
		&i.TokenId,
		&i.TenantId,
	)
	return &i, err
}