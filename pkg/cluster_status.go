package migration

import (
	"context"

	v293 "github.tools.sap/framefrog/cp-mod-migrator/pkg/cproxy/api/v294"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Status string

const (
	StatusUnknown           Status = "UNKNOWN"
	StatusMigrationRequired Status = "REQUIRED"
	StatusMigrationSkipped  Status = "SKIPPED"
)

//go:generate mockery --name=Client
type Client interface {
	client.Client
}

func GetStatus(ctx context.Context, c Client) (Status, error) {
	for _, check := range []Check{
		Not(ModuleInstalled),
		OldConnProxyInstalled,
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
	var cp appsv1.StatefulSet
	err := c.Get(ctx, client.ObjectKey{
		Namespace: "kyma-system",
		Name:      "connectivity-proxy",
	}, &cp)
	if errors.IsNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

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
