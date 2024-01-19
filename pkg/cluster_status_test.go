package migration_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	migration "github.tools.sap/framefrog/cp-mod-migrator/pkg"
)

var _ = Describe("ModuleInstalled function", func() {
	ctx := context.Background()

	It("should return false if module is not installed", func() {
		result, err := migration.ModuleInstalled(ctx, clientCpCrNotFound)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(result).Should(BeFalse())
	})

	It("should return true if module is installed", func() {
		result, err := migration.ModuleInstalled(ctx, clientCpCrFound)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(result).Should(BeTrue())
	})

	It("should return client error", func() {
		_, err := migration.ModuleInstalled(ctx, clientCpCrErr)
		Expect(err).Should(MatchError(ErrNotFoundTest))
	})
})

var _ = Describe("OldConnProxyInstalled function", func() {
	ctx := context.Background()

	It("should return false if connectivity-porxy stateful-set was not found", func() {
		result, err := migration.OldConnProxyInstalled(ctx, clientSsNotFound)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(result).Should(BeFalse())
	})

	It("should return true if connectivity-porxy stateful-set was found", func() {
		result, err := migration.OldConnProxyInstalled(ctx, clientSsFound)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(result).Should(BeTrue())
	})

	It("should return client error", func() {
		_, err := migration.OldConnProxyInstalled(ctx, clientSsErr)
		Expect(err).Should(MatchError(ErrTest))
	})
})
