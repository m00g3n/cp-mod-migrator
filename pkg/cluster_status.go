package migration

import (
	"context"

	v294 "github.tools.sap/framefrog/cp-mod-migrator/pkg/cproxy/api/v294"
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

func GetStatus(ctx context.Context, c Client) (Status, v294.ConnectivityProxy, error) {
	var cp v294.ConnectivityProxy
	// get connectivity proxy
	key := client.ObjectKey{
		Namespace: v294.CProxyDefaultCRNamespace,
		Name:      v294.CProxyDefaultCRName,
	}
	if err := c.Get(ctx, key, &cp); err != nil {
		return StatusUnknown, cp, err
	}
	// check if migration was not already performed
	if cp.Migrated() {
		return StatusMigrationSkipped, cp, nil
	}
	// check if old connectivity proxy is installed
	installed, err := OldConnProxyInstalled(ctx, c)
	// status unknown
	if err != nil {
		return StatusUnknown, cp, err
	}
	// migration is not required
	if !installed {
		return StatusMigrationSkipped, cp, nil
	}
	// migration is required
	return StatusMigrationRequired, cp, nil
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
