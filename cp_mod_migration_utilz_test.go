package main_test

import (
	"errors"
	"io"
	"os"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	v294 "github.tools.sap/framefrog/cp-mod-migrator/pkg/cproxy/api/v294"
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

func cp(name, namespace string) v294.ConnectivityProxy {
	return v294.ConnectivityProxy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "connectivityproxy.sap.com/v1",
			Kind:       "ConnectivityProxy",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v294.Spec{
			Config: v294.Config{
				HighAvailabilityMode: v294.HighAvailabilityModeOff,

				Integration: v294.Integration{
					AuditLog: v294.AuditLog{
						Mode: v294.AuditLogModeConsole,
					},
					ConnectivityService: v294.ConnectivityService{
						ServiceCredentialsKey: "test",
					},
				},

				Servers: v294.Servers{
					BusinessDataTunnel: v294.BusinessDataTunnel{
						ExternalHost: "test",
						ExternalPort: 20,
					},
					Proxy: v294.Proxy{
						HTTP: v294.HTTP{
							Port: 123,
						},
					},
				},

				ConnectivityService: v294.ConnectivityService{
					ServiceCredentialsKey: "test",
				},

				SubaccountID: "test-me-plz",
				TenantMode:   v294.TenantModeDedicated,
			},
		},
	}
}

func loadCMs(data *[]corev1.ConfigMap) error {
	file, err := os.Open("../../hack/testdata/configmaps.yaml")
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
