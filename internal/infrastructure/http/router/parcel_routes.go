package router

import (
	"github.com/gin-gonic/gin"

	"ms-parcel-core/internal/infrastructure/http/handler"
	parcelclients "ms-parcel-core/internal/parcel/parcel_core/infrastructure/clients"
	parcelrepo "ms-parcel-core/internal/parcel/parcel_core/infrastructure/repository"
	"ms-parcel-core/internal/parcel/parcel_core/usecase"
)

func RegisterParcelRoutes(rg *gin.RouterGroup) {
	repo := parcelrepo.NewInMemoryParcelRepository()
	tenantConfig := parcelclients.NewTenantConfigStubClient()
	createUC := usecase.NewCreateParcelUseCase(repo, tenantConfig)
	h := handler.NewParcelHandler(createUC)

	parcels := rg.Group("/parcels")
	{
		parcels.POST("", h.Create)
		parcels.POST("/", h.Create)
	}
}
