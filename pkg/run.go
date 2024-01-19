package migration

import (
	"context"
	"log/slog"
	"os"
	"time"

	v294 "github.tools.sap/framefrog/cp-mod-migrator/pkg/cproxy/api/v294"
	"github.tools.sap/framefrog/cp-mod-migrator/pkg/extract"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type getClient = func() (client.Client, error)

func Run(ctx context.Context, getClient getClient, dryRun []string) error {
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

	status, err := GetStatus(ctx, k8sClient)
	if err != nil {
		return err
	}

	if status != StatusMigrationRequired {
		slog.Info(
			"cluster will not be migrated",
			"migrationStatus", status,
		)
		return nil
	}

	var cp v294.ConnectivityProxy
	for _, f := range []extract.Function{
		extract.SetDefaults,
		extract.GetCPConfiguration,
	} {
		if err := f(ctx, &cp, k8sClient); err != nil {
			return err
		}
	}

	data, err := cp.Encode()
	if err != nil {
		return err
	}

	if err := k8sClient.Create(ctx, &cp, &client.CreateOptions{
		DryRun: dryRun,
	}); err != nil {
		return err
	}

	slog.Info("CR created", "data", data)

	return nil
}

func logDuration(msg string, start time.Time) {
	arg := slog.Attr{
		Key:   "duration",
		Value: slog.AnyValue(time.Since(start)),
	}
	slog.Info(msg, arg)
}
