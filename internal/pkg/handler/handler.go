package handler

import (
	"github.com/FudSy/Diploma/internal/pkg/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Static("/uploads", "./uploads")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	auth := router.Group("/auth")
	{
		auth.POST("/register", h.register)
		auth.POST("/login", h.login)
	}

	authProtected := router.Group("/auth", h.userIdentity)
	{
		authProtected.POST("/logout", h.logout)
		authProtected.GET("/me", h.me)
	}

	admin := router.Group("/auth/admin", h.userIdentity, h.adminIdentity)
	{
		admin.POST("/register", h.registerAdmin)
		admin.GET("/check", h.adminCheck)
	}

	resources := router.Group("/resources", h.userIdentity)
	{
		resources.GET("/", h.getAllResources)
		resources.GET("/:id", h.getResourceByID)
		resources.GET("/:id/availability", h.getResourceAvailability)
	}

	resourceAdmin := router.Group("/resources", h.userIdentity, h.adminIdentity)
	{
		resourceAdmin.POST("/", h.createResource)
		resourceAdmin.PUT("/:id", h.updateResource)
		resourceAdmin.DELETE("/:id", h.deleteResource)
		resourceAdmin.PATCH("/:id/capacity/increase", h.increaseResourceCapacity)
		resourceAdmin.PATCH("/:id/capacity/decrease", h.decreaseResourceCapacity)
		resourceAdmin.POST("/:id/photo", h.uploadResourcePhoto)
	}

	bookings := router.Group("/bookings", h.userIdentity)
	{
		bookings.POST("/", h.createBooking)
		bookings.GET("/", h.getMyBookings)
		bookings.GET("/:id", h.getBookingByID)
		bookings.PUT("/:id", h.updateBooking)
		bookings.DELETE("/:id", h.deleteBooking)
	}

	adminAPI := router.Group("/admin", h.userIdentity, h.adminIdentity)
	{
		adminAPI.GET("/bookings", h.getAllBookingsAdmin)
		adminAPI.GET("/stats", h.getStatsOverview)
	}

	resourceTypes := router.Group("/resource-types", h.userIdentity)
	{
		resourceTypes.GET("/", h.getResourceTypes)
	}

	resourceTypesAdmin := router.Group("/resource-types", h.userIdentity, h.adminIdentity)
	{
		resourceTypesAdmin.POST("/", h.createResourceType)
		resourceTypesAdmin.DELETE("/:id", h.deleteResourceType)
		resourceTypesAdmin.POST("/:id/options", h.addResourceTypeOption)
		resourceTypesAdmin.DELETE("/:id/options/:optionId", h.deleteResourceTypeOption)
	}

	return router
}
