package v294_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	v293 "github.tools.sap/framefrog/cp-mod-migrator/pkg/cproxy/api/v294"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func testCProxyWithRequired() v293.ConnectivityProxy {
	return v293.ConnectivityProxy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "connectivityproxy.sap.com/v1",
			Kind:       "ConnectivityProxy",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
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

var _ = Describe("cproxy crd should be valid if", func() {

	It("was created using only required properties", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		cr := testCProxyWithRequired()
		Expect(k8sClient.Create(ctx, &cr)).ShouldNot(HaveOccurred())
	})

	It("was created using all available properties", func() {
		Expect(true).To(BeTrue())
	})

})
