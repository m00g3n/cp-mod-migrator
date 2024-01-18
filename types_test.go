package main_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	v293 "github.tools.sap/framefrog/cp-mod-migrator/pkg/cproxy/api/v294"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func cp(name, namespace string) v293.ConnectivityProxy {
	return v293.ConnectivityProxy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "connectivityproxy.sap.com/v1",
			Kind:       "ConnectivityProxy",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v293.Spec{
			Config: v293.Config{
				HighAvailabilityMode: v293.HighAvailabilityModeOff,

				Integration: v293.Integration{
					AuditLog: v293.AuditLog{
						Mode: v293.AuditLogModeConsole,
					},
					ConnectivityService: v293.ConnectivityService{
						ServiceCredentialsKey: "test",
					},
				},

				Servers: v293.Servers{
					BusinessDataTunnel: v293.BusinessDataTunnel{
						ExternalHost: "test",
						ExternalPort: 20,
					},
					Proxy: v293.Proxy{
						HTTP: v293.HTTP{
							Port: 123,
						},
					},
				},

				ConnectivityService: v293.ConnectivityService{
					ServiceCredentialsKey: "test",
				},

				SubaccountID: "test-me-plz",
				TenantMode:   v293.TenantModeDedicated,
			},
		},
	}
}

func namespace(name string) corev1.Namespace {
	return corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

var _ = Describe("cproxy CR type", func() {

	It("should be compatible with it's schema", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		ns := namespace("test")
		Expect(k8sClient.Create(ctx, &ns)).ShouldNot(HaveOccurred())

		cr := cp("test", ns.Name)
		Expect(k8sClient.Create(ctx, &cr)).ShouldNot(HaveOccurred())
	})

})
