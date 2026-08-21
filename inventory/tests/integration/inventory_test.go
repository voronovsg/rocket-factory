//go:build integration
// +build integration

package integration

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	grpcMiddleware "github.com/voronovsg/rocket-factory/platform/pkg/middleware/grpc"
	inventoryV1 "github.com/voronovsg/rocket-factory/shared/pkg/proto/inventory/v1"
)

var _ = Describe("InventoryService", func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		inventoryClient inventoryV1.InventoryServiceClient
		partUUID        string
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)
		ctx = metadata.AppendToOutgoingContext(ctx, grpcMiddleware.SessionUUIDMetadataKey, testSessionUUID)

		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

		inventoryClient = inventoryV1.NewInventoryServiceClient(conn)

		partUUID, err = env.InsertTestPart(ctx)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестовой детали в MongoDB")
	})

	AfterEach(func() {
		err := env.ClearPartsCollection(ctx)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешную очистку коллекции parts")

		cancel()
	})

	It("должен успешно возвращать деталь по UUID", func() {
		resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
			Uuid: partUUID,
		})

		Expect(err).ToNot(HaveOccurred())
		Expect(resp.GetPart()).ToNot(BeNil())
		Expect(resp.GetPart().Uuid).To(Equal(partUUID))
		Expect(resp.GetPart().Name).ToNot(BeEmpty())
		Expect(resp.GetPart().GetManufacturer().Name).ToNot(BeEmpty())
	})
})
