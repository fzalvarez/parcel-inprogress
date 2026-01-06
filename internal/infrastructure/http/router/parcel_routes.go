package router

import (
	"github.com/gin-gonic/gin"

	"ms-parcel-core/internal/infrastructure/http/handler"
	parcelclients "ms-parcel-core/internal/parcel/parcel_core/infrastructure/clients"
	parcelrepo "ms-parcel-core/internal/parcel/parcel_core/infrastructure/repository"
	coreport "ms-parcel-core/internal/parcel/parcel_core/port"
	"ms-parcel-core/internal/parcel/parcel_core/usecase"
	docclients "ms-parcel-core/internal/parcel/parcel_documents/infrastructure/clients"
	docrepo "ms-parcel-core/internal/parcel/parcel_documents/infrastructure/repository"
	docusecase "ms-parcel-core/internal/parcel/parcel_documents/usecase"
	itemrepo "ms-parcel-core/internal/parcel/parcel_item/infrastructure/repository"
	itemusecase "ms-parcel-core/internal/parcel/parcel_item/usecase"
	paymentrepo "ms-parcel-core/internal/parcel/parcel_payment/infrastructure/repository"
	paymentusecase "ms-parcel-core/internal/parcel/parcel_payment/usecase"
	pricingrepo "ms-parcel-core/internal/parcel/parcel_pricing/infrastructure/repository"
	pricingusecase "ms-parcel-core/internal/parcel/parcel_pricing/usecase"
	trackingrecorder "ms-parcel-core/internal/parcel/parcel_tracking/infrastructure/recorder"
	trackingrepo "ms-parcel-core/internal/parcel/parcel_tracking/infrastructure/repository"
	trackingusecase "ms-parcel-core/internal/parcel/parcel_tracking/usecase"
)

func RegisterParcelRoutesWithDeps(
	rg *gin.RouterGroup,
	repo *parcelrepo.InMemoryParcelRepository,
	trkRepo *trackingrepo.InMemoryTrackingRepository,
	itemRepo *itemrepo.InMemoryParcelItemRepository,
	payRepo *paymentrepo.InMemoryParcelPaymentRepository,
	tenantConfig coreport.TenantConfigClient,
	tenantOptionsProvider coreport.TenantOptionsProvider,
) {
	trkRecorder := trackingrecorder.NewTrackingRecorderAdapter(trkRepo)

	createUC := usecase.NewCreateParcelUseCase(repo, tenantConfig, trkRecorder, tenantOptionsProvider)
	getUC := usecase.NewGetParcelUseCase(repo)
	listUC := usecase.NewListParcelsUseCase(repo)
	registerUC := usecase.NewRegisterParcelUseCase(repo, trkRecorder)
	boardUC := usecase.NewBoardParcelUseCase(repo, trkRecorder)
	departUC := usecase.NewDepartParcelUseCase(repo, trkRecorder)
	arriveUC := usecase.NewArriveParcelUseCase(repo, trkRecorder)
	deliverUC := usecase.NewDeliverParcelUseCase(repo, trkRecorder)

	parcelsHandler := handler.NewParcelHandler(createUC, getUC, registerUC, boardUC, deliverUC, arriveUC, departUC, listUC)

	priceRuleRepo := pricingrepo.NewInMemoryPriceRuleRepository()
	createRuleUC := pricingusecase.NewCreatePriceRuleUseCase(priceRuleRepo)
	updateRuleUC := pricingusecase.NewUpdatePriceRuleUseCase(priceRuleRepo)
	listRuleUC := pricingusecase.NewListPriceRulesUseCase(priceRuleRepo)
	rulesHandler := handler.NewPriceRuleHandler(createRuleUC, updateRuleUC, listRuleUC)

	addItemUC := itemusecase.NewAddParcelItemUseCase(repo, itemRepo, trkRecorder, tenantOptionsProvider, priceRuleRepo)
	listItemsUC := itemusecase.NewListParcelItemsUseCase(repo, itemRepo)
	deleteItemUC := itemusecase.NewDeleteParcelItemUseCase(repo, itemRepo, trkRecorder)
	itemsHandler := handler.NewParcelItemHandler(addItemUC, listItemsUC, deleteItemUC)

	cashboxClient := parcelclients.NewCashboxStubClient()
	upsertPayUC := paymentusecase.NewUpsertParcelPaymentUseCase(repo, payRepo, tenantOptionsProvider, cashboxClient)
	getPayUC := paymentusecase.NewGetParcelPaymentUseCase(payRepo)
	markPaidUC := paymentusecase.NewMarkPaidParcelPaymentUseCase(repo, payRepo, tenantOptionsProvider)
	paymentHandler := handler.NewParcelPaymentHandler(upsertPayUC, getPayUC, markPaidUC)

	listTrackingUC := trackingusecase.NewListTrackingUseCase(trkRepo)
	trackingHandler := handler.NewParcelTrackingHandler(listTrackingUC)

	// Summary
	summaryUC := usecase.NewGetParcelSummaryUseCase(repo, itemRepo, payRepo, trkRepo)
	summaryHandler := handler.NewParcelSummaryHandler(summaryUC)

	printRepo := docrepo.NewInMemoryPrintRepository()
	qrGen := docclients.NewStubQRGenerator()
	registerPrintUC := docusecase.NewRegisterPrintUseCase(repo, printRepo, tenantOptionsProvider, qrGen)
	docsHandler := handler.NewParcelDocumentsHandler(registerPrintUC, printRepo)

	parcels := rg.Group("/parcels")
	{
		parcels.GET("", parcelsHandler.List)
		parcels.POST("", parcelsHandler.Create)

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

		parcels.PUT("/:id/payment", paymentHandler.Upsert)
		parcels.GET("/:id/payment", paymentHandler.Get)
		parcels.POST("/:id/payment/mark-paid", paymentHandler.MarkPaid)

		parcels.GET("/:id/summary", summaryHandler.Get)

		parcels.POST("/:id/documents/print", docsHandler.RegisterPrint)
		parcels.GET("/:id/documents/prints", docsHandler.ListPrints)
	}

	pricing := rg.Group("/pricing")
	{
		pricing.POST("/rules", rulesHandler.Create)
		pricing.PUT("/rules/:id", rulesHandler.Update)
		pricing.GET("/rules", rulesHandler.List)
	}
}
