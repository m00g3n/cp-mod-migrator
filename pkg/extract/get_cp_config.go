package extract

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"

	v293 "github.tools.sap/framefrog/cp-mod-migrator/pkg/cproxy/api/v293"
	v294 "github.tools.sap/framefrog/cp-mod-migrator/pkg/cproxy/api/v294"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	cmKey = client.ObjectKey{
		Namespace: v293.CProxyCMNamespace,
		Name:      v293.CProxyCMName,
	}

	ErrNotFound = fmt.Errorf("not foud")
)

//go:generate mockery --name=Client
type Client interface {
	client.Client
}

type Function func(context.Context, *v294.ConnectivityProxy, Client) error

func GetCPConfiguration(ctx context.Context, cr *v294.ConnectivityProxy, c Client) error {
	var cm corev1.ConfigMap
	if err := c.Get(ctx, cmKey, &cm); err != nil {
		return err
	}

	cfgFile, found := cm.Data[v293.CProxyConfigFilename]
	if !found {
		return fmt.Errorf("%w: cm '%s:%s' missing '%s' key",
			ErrNotFound,
			v293.CProxyCMNamespace,
			v293.CProxyCMName,
			v293.CProxyConfigFilename,
		)
	}

	if err := yaml.UnmarshalStrict([]byte(cfgFile), &cr.Spec.Config); err != nil {
		data := []byte(cfgFile)
		encode := base64.StdEncoding.EncodeToString(data)

		slog.Warn("unable to unmarshal config file",
			v293.CProxyConfigFilename, encode,
		)

		return err
	}

	return nil
}
