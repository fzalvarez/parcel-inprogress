package router

import (
	"github.com/gin-gonic/gin"

	"ms-parcel-core/internal/infrastructure/http/handler"
	parcelclients "ms-parcel-core/internal/parcel/parcel_core/infrastructure/clients"
	parcelrepo "ms-parcel-core/internal/parcel/parcel_core/infrastructure/repository"
	"ms-parcel-core/internal/parcel/parcel_core/usecase"
	itemrepo "ms-parcel-core/internal/parcel/parcel_item/infrastructure/repository"
	itemusecase "ms-parcel-core/internal/parcel/parcel_item/usecase"
	trackingrecorder "ms-parcel-core/internal/parcel/parcel_tracking/infrastructure/recorder"
	trackingrepo "ms-parcel-core/internal/parcel/parcel_tracking/infrastructure/repository"
	trackingusecase "ms-parcel-core/internal/parcel/parcel_tracking/usecase"
)

func RegisterParcelRoutes(rg *gin.RouterGroup) {
	repo := parcelrepo.NewInMemoryParcelRepository()
	tenantConfig := parcelclients.NewTenantConfigStubClient()

	trkRepo := trackingrepo.NewInMemoryTrackingRepository()
	trkRecorder := trackingrecorder.NewTrackingRecorderAdapter(trkRepo)

	createUC := usecase.NewCreateParcelUseCase(repo, tenantConfig, trkRecorder)
	getUC := usecase.NewGetParcelUseCase(repo)
	listUC := usecase.NewListParcelsUseCase(repo)
	registerUC := usecase.NewRegisterParcelUseCase(repo, trkRecorder)
	boardUC := usecase.NewBoardParcelUseCase(repo, trkRecorder)
	departUC := usecase.NewDepartParcelUseCase(repo, trkRecorder)
	arriveUC := usecase.NewArriveParcelUseCase(repo, trkRecorder)
	deliverUC := usecase.NewDeliverParcelUseCase(repo, trkRecorder)

	parcelsHandler := handler.NewParcelHandler(createUC, getUC, registerUC, boardUC, deliverUC, arriveUC, departUC, listUC)

	itemRepo := itemrepo.NewInMemoryParcelItemRepository()
	addItemUC := itemusecase.NewAddParcelItemUseCase(repo, itemRepo, trkRecorder)
	listItemsUC := itemusecase.NewListParcelItemsUseCase(repo, itemRepo)
	deleteItemUC := itemusecase.NewDeleteParcelItemUseCase(repo, itemRepo, trkRecorder)
	itemsHandler := handler.NewParcelItemHandler(addItemUC, listItemsUC, deleteItemUC)

	listTrackingUC := trackingusecase.NewListTrackingUseCase(trkRepo)
	trackingHandler := handler.NewParcelTrackingHandler(listTrackingUC)

	parcels := rg.Group("/parcels")
	{
		parcels.GET("", parcelsHandler.List)
		parcels.GET("/", parcelsHandler.List)

		parcels.POST("", parcelsHandler.Create)
		parcels.POST("/", parcelsHandler.Create)

		parcels.GET("/:id", parcelsHandler.GetByID)

		parcels.POST("/:id/register", parcelsHandler.Register)
		parcels.POST("/:id/board", parcelsHandler.Board)
		parcels.POST("/:id/depart", parcelsHandler.Depart)
		parcels.POST("/:id/arrive", parcelsHandler.Arrive)
		parcels.POST("/:id/deliver", parcelsHandler.Deliver)

		parcels.GET("/:id/tracking", trackingHandler.ListByParcelID)

		parcels.POST("/:id/items", itemsHandler.Add)
		parcels.GET("/:id/items", itemsHandler.List)
		parcels.DELETE("/:id/items/:item_id", itemsHandler.Delete)
	}
}
