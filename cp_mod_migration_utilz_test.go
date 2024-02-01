package main_test

import (
	"errors"
	"io"
	"os"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	v211 "github.tools.sap/framefrog/cp-mod-migrator/pkg/cproxy/api/v211"
	"github.tools.sap/framefrog/cp-mod-migrator/pkg/extract"
	"github.tools.sap/framefrog/cp-mod-migrator/pkg/mocks"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func readYaml[T any](name string, out T) {
	file, err := os.Open(name)
	if err != nil {
		Expect(err).ShouldNot(HaveOccurred())
	}

	decoder := yaml.NewYAMLOrJSONDecoder(file, 2048)
	Expect(decoder.Decode(&out)).ShouldNot(HaveOccurred())
}

func namespace(name string) corev1.Namespace {
	return corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

func cp(name, namespace string, migrated bool) v211.ConnectivityProxy {
	annotations := map[string]string{}
	if migrated {
		annotations[v211.CProxyMigratedAnnotation] = "true"
	}

	return v211.ConnectivityProxy{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: annotations,
		},
		Spec: v211.Spec{
			Config: v211.Config{
				HighAvailabilityMode: v211.HighAvailabilityModeOff,

				Integration: v211.Integration{
					AuditLog: v211.AuditLog{
						Mode: v211.AuditLogModeConsole,
					},
					ConnectivityService: &v211.ConnectivityService{
						ServiceCredentialsKey: "test",
					},
				},

				Servers: v211.Servers{
					BusinessDataTunnel: v211.BusinessDataTunnel{
						ExternalHost: "test",
						ExternalPort: 20,
					},
					Proxy: v211.Proxy{
						HTTP: v211.ProxyCfg{
							Port: 123,
						},
						RfcAndLdap: v211.ProxyCfg{
							Port: 123,
						},
						Socks5: v211.ProxyCfg{
							Port: 123,
						},
					},
				},

				ConnectivityService: v211.ConnectivityService{
					ServiceCredentialsKey: "test",
				},

				SubaccountID: "test-me-plz",
				TenantMode:   v211.TenantModeDedicated,
			},

			Ingress: v211.Ingress{
				ClassName: v211.ClassTypeIstio,
			},

			SecretConfig: v211.SecretConfig{
				Integration: v211.SecretConfigIntegration{
					ConnectivityService: v211.ServiceSecretConfig{
						SecretName: "test",
					},
				},
			},
		},
	}
}

func loadCMs(data *[]corev1.ConfigMap) error {
	file, err := os.Open("./hack/testdata/configmaps.yaml")
	if err != nil {
		return err
	}

	decoder := yaml.NewYAMLOrJSONDecoder(file, 2048)
	var cm corev1.ConfigMap
	for {
		err := decoder.Decode(&cm)
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return err
		}

		*data = append(*data, cm)
	}
	return nil
}

func newClient(t *testing.T, cms []corev1.ConfigMap) extract.Client {
	var index int
	runFn := func(args mock.Arguments) {
		cm := args.Get(2).(*corev1.ConfigMap)
		*cm = cms[index]
		if index == len(cms)-1 {
			index = 0
			return
		}
		index++
	}

	client := mocks.NewClient(t)
	client.On("Get", mock.Anything, mock.Anything, mock.Anything).
		Run(runFn).
		Return(nil).
		Times(len(cms))

	return client
}
