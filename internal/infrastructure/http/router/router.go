package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ms-parcel-core/internal/infrastructure/http/handler"
	parcelrepo "ms-parcel-core/internal/parcel/parcel_core/infrastructure/repository"
	manifestusecase "ms-parcel-core/internal/parcel/parcel_manifest/usecase"
)

func RegisterRoutes(engine *gin.Engine) {
	// Health mínimo para verificar server correcto
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := engine.Group("/api/v1")
	{
		RegisterParcelRoutes(v1)

		// Manifests (preview virtual)
		parcelRepo := parcelrepo.NewInMemoryParcelRepository()
		buildUC := manifestusecase.NewBuildManifestPreviewUseCase(parcelRepo)
		h := handler.NewManifestHandler(buildUC)

		manifests := v1.Group("/manifests")
		{
			manifests.POST("/preview", h.PreviewPost)
			manifests.GET("/preview", h.PreviewGet)
		}
	}
}
