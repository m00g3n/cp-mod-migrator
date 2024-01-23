package main_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	migration "github.tools.sap/framefrog/cp-mod-migrator/pkg"
	v294 "github.tools.sap/framefrog/cp-mod-migrator/pkg/cproxy/api/v294"
	"github.tools.sap/framefrog/cp-mod-migrator/pkg/extract"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func deleteObjs(ctx context.Context, obj client.Object, objs ...client.Object) {
	for _, obj := range append(objs, obj) {
		err := k8sClient.Delete(ctx, obj)
		if !errors.IsNotFound(err) {
			Expect(err).ShouldNot(HaveOccurred())
		}
	}
}

var _ = Describe("cp-mod-migrator", Ordered, func() {

	var (
		ns     corev1.Namespace
		cm     corev1.ConfigMap
		cmInfo corev1.ConfigMap
		sSet   appsv1.StatefulSet
	)

	BeforeAll(func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		// create kyma-system namespace
		ns = namespace("kyma-system")
		Expect(k8sClient.Create(ctx, &ns)).ShouldNot(HaveOccurred())
		// read data
		readYaml("hack/testdata/cp_cm.yaml", &cm)
		readYaml("hack/testdata/cp_cm_info.yaml", &cmInfo)
		readYaml("hack/testdata/cp_stateful_set.yaml", &sSet)
	})

	It("should have types compatible with connectivity-porxy schema", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		cr := cp("connectivity-proxy", ns.Name, true)
		Expect(k8sClient.Create(ctx, &cr)).ShouldNot(HaveOccurred())
		Expect(k8sClient.Delete(ctx, &cr)).ShouldNot(HaveOccurred())
	})

	It("should cover all existing cases", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		for range cms {
			var cr v294.ConnectivityProxy
			err := extract.GetCPConfiguration(ctx, &cr, mockedClient)
			Expect(err).ShouldNot(HaveOccurred())

			err = extract.SetDefaults(ctx, &cr, mockedClient)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(k8sClient.Create(ctx, &cr)).ShouldNot(HaveOccurred())
			deleteObjs(ctx, &cr)
		}
	})

	It("should migrate data", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		// create config-map with configuration
		cmCopy := cm.DeepCopy()
		Expect(k8sClient.Create(ctx, cmCopy)).ShouldNot(HaveOccurred())
		// create config-map with proxy information
		cmInfoCopy := cmInfo.DeepCopy()
		Expect(k8sClient.Create(ctx, cmInfoCopy)).ShouldNot(HaveOccurred())
		// create statefu-set
		sSetCopy := sSet.DeepCopy()
		Expect(k8sClient.Create(ctx, sSetCopy)).ShouldNot(HaveOccurred())
		// create CR to be migrated
		cr := cp("connectivity-proxy", ns.Name, false)
		Expect(k8sClient.Create(ctx, cr.DeepCopy())).ShouldNot(HaveOccurred())
		// start migration
		Expect(migration.Run(ctx, getK8sClient, []string{})).ShouldNot(HaveOccurred())
		// fetch created CR
		key := client.ObjectKey{Name: "connectivity-proxy", Namespace: "kyma-system"}
		Expect(k8sClient.Get(ctx, key, &cr)).ShouldNot(HaveOccurred())
		Expect(cr.Annotations).Should(HaveKeyWithValue(v294.CProxyMigratedAnnotation, ""))
		// clean up
		deleteObjs(ctx, cmCopy, cmInfoCopy, sSetCopy, &cr)
	})

	It("should not migrate data when CP module is installed on a cluster", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		// create CR
		cr := cp(v294.CProxyDefaultCRName, ns.Name, true)
		Expect(k8sClient.Create(ctx, &cr)).ShouldNot(HaveOccurred())
		// create config-map with configuration
		cmCopy := cm.DeepCopy()
		Expect(k8sClient.Create(ctx, cmCopy)).ShouldNot(HaveOccurred())
		// create statefu-set
		sSetCopy := sSet.DeepCopy()
		Expect(k8sClient.Create(ctx, sSetCopy)).ShouldNot(HaveOccurred())
		// start migration
		Expect(migration.Run(ctx, getK8sClient, []string{})).ShouldNot(HaveOccurred())
		deleteObjs(ctx, cmCopy, sSetCopy, &cr)
	})

	It("should not migrate data when non modular CP component is not installed on a cluster", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		// create config-map with configuration
		cmCopy := cm.DeepCopy()
		Expect(k8sClient.Create(ctx, cmCopy)).ShouldNot(HaveOccurred())
		// create CR
		cr := cp(v294.CProxyDefaultCRName, ns.Name, true)
		Expect(k8sClient.Create(ctx, &cr)).ShouldNot(HaveOccurred())
		// start migration
		Expect(migration.Run(ctx, getK8sClient, []string{})).ShouldNot(HaveOccurred())
		deleteObjs(ctx, cmCopy, &cr)
	})

	It("should not migrate when non modular CP component is corrupted (missing configuration)", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		// create CR
		cr := cp(v294.CProxyDefaultCRName, ns.Name, false)
		Expect(k8sClient.Create(ctx, &cr)).ShouldNot(HaveOccurred())
		// create statefu-set
		sSetCopy := sSet.DeepCopy()
		Expect(k8sClient.Create(ctx, sSetCopy)).ShouldNot(HaveOccurred())
		// start migration
		Expect(migration.Run(ctx, getK8sClient, []string{})).Should(MatchError(`configmaps "connectivity-proxy" not found`))
		deleteObjs(ctx, sSetCopy)
	})

	It("should not migrate empty cluster", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		// start migration
		Expect(migration.Run(ctx, getK8sClient, []string{})).ShouldNot(HaveOccurred())
	})
})
