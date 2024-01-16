package extract_test

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"

	. "github.com/onsi/gomega"
	v294 "github.tools.sap/framefrog/cp-mod-migrator/pkg/cproxy/api/v294"
	"github.tools.sap/framefrog/cp-mod-migrator/pkg/extract"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

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

func TestGetCPConfiguration(t *testing.T) {
	g := NewWithT(t)

	var cms []corev1.ConfigMap
	g.Expect(loadCMs(&cms)).ShouldNot(HaveOccurred())

	client := newClient(t, cms)

	var cr v294.ConnectivityProxy
	t.Log("case number:", len(cms))
	for range cms {
		g.Expect(extract.GetCPConfiguration(context.Background(), &cr, client)).
			ShouldNot(HaveOccurred())
	}
}
