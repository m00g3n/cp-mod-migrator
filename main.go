package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	migration "github.tools.sap/framefrog/cp-mod-migrator/pkg"
	v294 "github.tools.sap/framefrog/cp-mod-migrator/pkg/cproxy/api/v294"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type config struct {
	kubeconfigPath string
	isDryRun       bool
}

// dryRun - returns dry run val accepted by client patch method
func (c *config) dryRun() []string {
	if c.isDryRun {
		return []string{"All"}
	}
	return []string{}
}

func addToScheme(s *runtime.Scheme) error {
	for _, add := range []func(s *runtime.Scheme) error{
		v294.AddToScheme,
		corev1.AddToScheme,
	} {
		if err := add(s); err != nil {
			return fmt.Errorf("unable to add scheme: %s", err)
		}

	}
	return nil
}

func (cfg *config) client() (client.Client, error) {
	restCfg, err := clientcmd.BuildConfigFromFlags("", cfg.kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch rest config: %w", err)
	}

	scheme := runtime.NewScheme()
	if err := addToScheme(scheme); err != nil {
		return nil, err
	}

	return client.New(restCfg, client.Options{
		Scheme: scheme,
		Cache: &client.CacheOptions{
			DisableFor: []client.Object{
				&corev1.Secret{},
			},
		},
	})
}

// kubeconfigPathVar - registers kubeconfigPath flag
func kubeconfigPathVar(kubeconfigFlag *string) {
	if home := homedir.HomeDir(); home != "" {
		flag.StringVar(
			kubeconfigFlag,
			"kubeconfig",
			filepath.Join(home, ".kube", "config"),
			"(optional) absolute path to the kubeconfig file")
		return
	}

	flag.StringVar(kubeconfigFlag, "kubeconfig", "", "absolute path to the kubeconfig file")
}

// newConfig - creates new application configuration base on passed flags
func newConfig() config {
	result := config{}
	kubeconfigPathVar(&result.kubeconfigPath)
	flag.BoolVar(
		&result.isDryRun,
		"dry-run",
		false,
		"(optional) indicates that modifications should not be persisted",
	)

	flag.Parse()
	return result
}

func exit1(err error) {
	slog.Error(err.Error())
	os.Exit(1)
}

func main() {
	cfg := newConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dryRun := cfg.dryRun()
	if err := migration.Run(ctx, cfg.client, dryRun); err != nil {
		exit1(err)
	}
}
