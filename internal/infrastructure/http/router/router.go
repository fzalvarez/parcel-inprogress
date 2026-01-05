package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"ms-parcel-core/internal/infrastructure/http/handler"
	parcelclients "ms-parcel-core/internal/parcel/parcel_core/infrastructure/clients"
	parcelrepo "ms-parcel-core/internal/parcel/parcel_core/infrastructure/repository"
	itemrepo "ms-parcel-core/internal/parcel/parcel_item/infrastructure/repository"
	manifestusecase "ms-parcel-core/internal/parcel/parcel_manifest/usecase"
	paymentrepo "ms-parcel-core/internal/parcel/parcel_payment/infrastructure/repository"
	trackingrepo "ms-parcel-core/internal/parcel/parcel_tracking/infrastructure/repository"
)

func RegisterRoutes(engine *gin.Engine) {
	// Health mínimo para verificar server correcto
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := engine.Group("/api/v1")
	{
		// Composition root (deps únicos)
		parcelRepo := parcelrepo.NewInMemoryParcelRepository()
		trkRepo := trackingrepo.NewInMemoryTrackingRepository()
		itemRepo := itemrepo.NewInMemoryParcelItemRepository()
		payRepo := paymentrepo.NewInMemoryParcelPaymentRepository()

		tenantConfig := parcelclients.NewTenantConfigStubClient()
		tenantOptionsProvider := parcelclients.NewCachedTenantOptionsProvider(tenantConfig, 60*time.Second)

		RegisterParcelRoutesWithDeps(v1, parcelRepo, trkRepo, itemRepo, payRepo, tenantConfig, tenantOptionsProvider)

		// Manifests (preview virtual)
		buildUC := manifestusecase.NewBuildManifestPreviewUseCase(parcelRepo)
		h := handler.NewManifestHandler(buildUC)

		manifests := v1.Group("/manifests")
		{
			manifests.POST("/preview", h.PreviewPost)
			manifests.GET("/preview", h.PreviewGet)
		}
	}
}
