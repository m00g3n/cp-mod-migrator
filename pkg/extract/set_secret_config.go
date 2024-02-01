package extract

import (
	"context"
	"fmt"

	v211 "github.tools.sap/framefrog/cp-mod-migrator/pkg/cproxy/api/v211"
	v293 "github.tools.sap/framefrog/cp-mod-migrator/pkg/cproxy/api/v293"
)

var ErrInvalidValue = fmt.Errorf("invalid value")

func SetSecretConfig(_ context.Context, cr *v211.ConnectivityProxy, _ Client) error {
	if cr == nil {
		return fmt.Errorf("%w: %s", ErrInvalidValue, "cr must not be nil")
	}
	// go with defaults if tenant mode is not shared
	if cr.Spec.Config.TenantMode != v211.TenantModeShared {
		return nil
	}

	cr.Spec.SecretConfig.Integration.ConnectivityService.SecretName = v293.CProxyConnectivityServiceSecretName
	return nil
}
