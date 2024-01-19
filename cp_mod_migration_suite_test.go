package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	v293 "github.tools.sap/framefrog/cp-mod-migrator/pkg/cproxy/api/v294"
	"k8s.io/client-go/rest"
	"k8s.io/kubectl/pkg/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

func TestCproxy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cproxy Suite")
}

var (
	externalDependencyDataPath = "./hack/cproxy/crd.yaml"

	testEnv   *envtest.Environment
	k8sClient client.Client
	config    *rest.Config
)

func getK8sClient() (client.Client, error) {
	return k8sClient, nil
}

var _ = BeforeSuite(func() {
	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{
			externalDependencyDataPath,
		}, ErrorIfCRDPathMissing: true,
	}

	var err error
	config, err = testEnv.Start()

	Expect(err).NotTo(HaveOccurred())
	Expect(config).NotTo(BeNil())

	err = v293.AddToScheme(scheme.Scheme)
	Expect(err).ShouldNot(HaveOccurred())

	k8sClient, err = client.New(config, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
