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
		return fmt.Errorf("%w: cr must not be nil", ErrInvalidValue)
	}

	if cr.Spec.Config.TenantMode == v211.TenantModeDedicated {
		return nil
	}

	if cr.Spec.Config.TenantMode == v211.TenantModeShared {
		cr.Spec.SecretConfig.Integration.ConnectivityService.SecretName = v293.CProxyConnectivityServiceSecretName
		return nil
	}

	return fmt.Errorf("%w: '%s'", ErrInvalidValue, cr.Spec.Config.TenantMode)
}
