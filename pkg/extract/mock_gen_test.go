package extract_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.tools.sap/framefrog/cp-mod-migrator/pkg/extract"
	"github.tools.sap/framefrog/cp-mod-migrator/pkg/extract/mocks"
	corev1 "k8s.io/api/core/v1"
)

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
