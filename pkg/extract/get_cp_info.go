package extract

import (
	"context"
	"fmt"
	"strconv"

	v211 "github.tools.sap/framefrog/cp-mod-migrator/pkg/cproxy/api/v211"
	v293 "github.tools.sap/framefrog/cp-mod-migrator/pkg/cproxy/api/v293"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	cmKeyInfo = client.ObjectKey{
		Namespace: v293.CProxyCMInfoNamespace,
		Name:      v293.CProxyCMInfoName,
	}
)

func GetCPInfo(ctx context.Context, cr *v211.ConnectivityProxy, c Client) error {
	var cm corev1.ConfigMap
	if err := c.Get(ctx, cmKeyInfo, &cm); err != nil {
		return err
	}

	type data struct {
		key string
		val *int
	}

	for _, data := range []data{
		{
			key: v293.CProxyOnpremiseProxyHttpPort,
			val: &cr.Spec.Config.Servers.Proxy.HTTP.Port,
		},
		{
			key: v293.CProxyOnpremiseProxyLdapPort,
			val: &cr.Spec.Config.Servers.Proxy.RfcAndLdap.Port,
		},
		{
			key: v293.CProxyOnpremiseSocks5ProxyPort,
			val: &cr.Spec.Config.Servers.Proxy.Socks5.Port,
		},
	} {
		err := valueOf(cm.Data, data.key, data.val)
		if err != nil {
			return err
		}
	}
	return nil
}

func valueOf(data map[string]string, key string, val *int) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			return
		}
	}()

	strValue, found := data[key]
	if !found {
		return fmt.Errorf("%w: cm '%s:%s' missing '%s' key",
			ErrNotFound,
			v293.CProxyCMInfoNamespace,
			v293.CProxyCMInfoName,
			key,
		)
	}
	value, err := strconv.Atoi(strValue)
	if err != nil {
		return fmt.Errorf("%w: invalid value %s.%s[%s]=%s",
			err,
			v293.CProxyCMInfoNamespace,
			v293.CProxyCMInfoName,
			key,
			strValue,
		)
	}
	*val = value
	return
}
