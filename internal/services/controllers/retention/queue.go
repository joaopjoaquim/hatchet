package retention

import (
	"context"
	"time"

	"github.com/hatchet-dev/hatchet/internal/telemetry"
	"github.com/hatchet-dev/hatchet/pkg/repository/prisma/dbsqlc"
	"github.com/hatchet-dev/hatchet/pkg/repository/prisma/sqlchelpers"
)

func (rc *RetentionControllerImpl) runDeleteQueueItems(ctx context.Context) func() {
	return func() {
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		rc.l.Debug().Msgf("retention controller: deleting queue items")

		err := rc.ForTenants(ctx, rc.runDeleteQueueItemsTenant)

		if err != nil {
			rc.l.Err(err).Msg("could not run delete expired job runs")
		}
	}
}

func (rc *RetentionControllerImpl) runDeleteQueueItemsTenant(ctx context.Context, tenant dbsqlc.Tenant) error {
	ctx, span := telemetry.NewSpan(ctx, "delete-queue-items-tenant")
	defer span.End()

	tenantId := sqlchelpers.UUIDToStr(tenant.ID)
	return rc.repo.StepRun().CleanupQueueItems(ctx, tenantId)
}