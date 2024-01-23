package extract

import (
	"context"

	v294 "github.tools.sap/framefrog/cp-mod-migrator/pkg/cproxy/api/v294"
)

func SetDefaults(ctx context.Context, cr *v294.ConnectivityProxy, _ Client) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	cr.Name = v294.CProxyDefaultCRName
	cr.Namespace = v294.CProxyDefaultCRNamespace
	cr.Spec.Deployment.RestartWatcher.Enabled = true
	cr.Spec.Ingress.ClassName = v294.ClassTypeIstio
	cr.Spec.Ingress.Tls.SecretName = "cc-certs"
	cr.Spec.Ingress.Timeouts.Proxy.Connect = 20
	cr.Spec.Ingress.Timeouts.Proxy.Read = 120
	cr.Spec.Ingress.Timeouts.Proxy.Send = 120
	cr.Spec.SecretConfig.Integration.ConnectivityService.SecretName = "connectivity-proxy-service-key"
	return nil
}
