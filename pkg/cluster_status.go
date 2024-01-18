package migration

import (
	"context"
	"log/slog"

	v293 "github.tools.sap/framefrog/cp-mod-migrator/pkg/cproxy/api/v294"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Status string

const (
	StatusUnknown           Status = "UNKNOWN"
	StatusMigrationRequired Status = "REQUIRED"
	StatusMigrationSkipped  Status = "SKIPPED"
)

type Client interface {
	client.Client
}

func GetStatus(ctx context.Context, c Client) (Status, error) {
	for _, check := range []Check{
		Not(ModuleInstalled),
	} {
		passed, err := check(ctx, c)
		if err != nil {
			return StatusUnknown, err
		}
		if !passed {
			return StatusMigrationSkipped, nil
		}
	}
	return StatusMigrationRequired, nil
}

type Check func(context.Context, Client) (bool, error)

func Not(check Check) Check {
	return func(ctx context.Context, c Client) (bool, error) {
		passed, err := check(ctx, c)
		return !passed, err
	}
}

func OldConnProxyInstalled(ctx context.Context, c Client) (bool, error) {
	slog.Warn("not implemented yet")
	return true, nil
}

func ModuleInstalled(ctx context.Context, c Client) (bool, error) {
	var cps v293.ConnectivityProxyList
	if err := c.List(ctx, &cps); err != nil {
		return false, err
	}
	if len(cps.Items) != 0 {
		return true, nil
	}
	return false, nil
}
