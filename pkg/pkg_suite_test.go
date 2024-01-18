package migration_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	migration "github.tools.sap/framefrog/cp-mod-migrator/pkg"
	v294 "github.tools.sap/framefrog/cp-mod-migrator/pkg/cproxy/api/v294"
	"github.tools.sap/framefrog/cp-mod-migrator/pkg/mocks"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	clientCpCrNotFound migration.Client
	clientCpCrFound    migration.Client
	clientCpCrErr      migration.Client
	clientSsNotFound   migration.Client
	clientSsFound      migration.Client
	clientSsErr        migration.Client

	ErrNotFoundTest = errors.NewNotFound(schema.GroupResource{}, "test not found")
	ErrTest         = errors.NewUnauthorized("test error")
	mockArgsAny     = []interface{}{mock.Anything, mock.Anything, mock.Anything}
)

func setupClientCpCR(t *testing.T, items ...v294.ConnectivityProxy) (migration.Client, error) {
	runFn := func(args mock.Arguments) {
		cps := args.Get(1).(*v294.ConnectivityProxyList)
		*cps = v294.ConnectivityProxyList{
			Items: append([]v294.ConnectivityProxy{}, items...),
		}
	}

	mockClient := mocks.NewClient(t)
	mockClient.On("List", mockArgsAny...).Return(nil).Run(runFn)

	return mockClient, nil
}

func setupClientCpCRErr(t *testing.T) error {
	mockClient := mocks.NewClient(t)
	mockClient.On("List", mockArgsAny...).Return(ErrNotFoundTest)

	clientCpCrErr = mockClient
	return nil
}

func setupClientCpCRFound(t *testing.T) (err error) {
	clientCpCrFound, err = setupClientCpCR(t, v294.ConnectivityProxy{})
	return err
}

func setupClientCpCRNotFound(t *testing.T) (err error) {
	clientCpCrNotFound, err = setupClientCpCR(t)
	return err
}

func setupClientSsNotFound(t *testing.T) (err error) {
	mockedClient := mocks.NewClient(t)
	mockedClient.On("Get", mockArgsAny...).Return(ErrNotFoundTest)

	clientSsNotFound = mockedClient
	return
}

func setupClientSsFound(t *testing.T) (err error) {
	runFn := func(args mock.Arguments) {
		cps := args.Get(2).(*appsv1.StatefulSet)
		*cps = appsv1.StatefulSet{}
	}

	mockedClient := mocks.NewClient(t)
	mockedClient.On("Get", mockArgsAny...).Return(nil).Run(runFn)

	clientSsFound = mockedClient
	return
}

func setupClientSsErr(t *testing.T) error {
	mockClient := mocks.NewClient(t)
	mockClient.On("Get", mockArgsAny...).Return(ErrTest)

	clientSsErr = mockClient
	return nil
}

type setupMockClient func(*testing.T) error

func TestPkg(t *testing.T) {
	RegisterFailHandler(Fail)

	BeforeSuite(func() {
		for _, setupClient := range []setupMockClient{
			setupClientCpCRFound,
			setupClientCpCRNotFound,
			setupClientCpCRErr,
			setupClientSsNotFound,
			setupClientSsFound,
			setupClientSsErr,
		} {
			err := setupClient(t)
			Expect(err).ShouldNot(HaveOccurred(), "unable to setup mock client")
		}
	})

	RunSpecs(t, "Pkg Suite")
}
