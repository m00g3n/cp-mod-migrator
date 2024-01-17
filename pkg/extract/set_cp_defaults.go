package extract

import (
	"context"

	v294 "github.tools.sap/framefrog/cp-mod-migrator/pkg/cproxy/api/v294"
)

func SetDefaults(ctx context.Context, cr *v294.ConnectivityProxy, _ Client) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	cr.Name = "default"
	cr.Namespace = "default"
	return nil
}
