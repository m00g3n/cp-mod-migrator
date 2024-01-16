package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	migration "github.tools.sap/framefrog/cp-mod-migrator/pkg"
	v294 "github.tools.sap/framefrog/cp-mod-migrator/pkg/cproxy/api/v294"
	"github.tools.sap/framefrog/cp-mod-migrator/pkg/extract"
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

type getClient = func() (client.Client, error)

func logDuration(msg string, start time.Time) {
	arg := slog.Attr{
		Key:   "duration",
		Value: slog.AnyValue(time.Since(start)),
	}
	slog.Info(msg, arg)
}

func run(ctx context.Context, getClient getClient) error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	slog.SetDefault(logger)

	start := time.Now()
	slog.Info("started")
	defer logDuration("finished in", start)

	k8sClient, err := getClient()
	if err != nil {
		return err
	}

	status, err := migration.GetStatus(ctx, k8sClient)
	if err != nil {
		return err
	}

	if status != migration.StatusMigrationRequired {
		slog.Info(
			"cluster will not be migrated",
			"migrationStatus", status,
		)
	}

	var cp v294.ConnectivityProxy
	for _, f := range []extract.Function{
		extract.GetCPConfiguration,
	} {
		if err := f(ctx, &cp, k8sClient); err != nil {
			return err
		}
	}

	if err := json.NewEncoder(os.Stdout).Encode(&cp); err != nil {
		return err
	}

	return nil
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

	if err := run(ctx, cfg.client); err != nil {
		exit1(err)
	}
}
