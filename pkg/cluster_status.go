package migration

import (
	"context"

	v293 "github.tools.sap/framefrog/cp-mod-migrator/pkg/cproxy/api/v294"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Status string

const (
	StatusUnknown           Status = "UNKNOWN"
	StatusMigrationRequired Status = "REQUIRED"
	StatusMigrationSkipped  Status = "SKIPPED"
)

func GetStatus(ctx context.Context, c client.Client) (Status, error) {
	var cps v293.ConnectivityProxyList
	if err := c.List(ctx, &cps); err != nil {
		return StatusUnknown, err
	}

	if len(cps.Items) != 0 {
		return StatusMigrationSkipped, nil
	}

	return StatusMigrationRequired, nil
}
