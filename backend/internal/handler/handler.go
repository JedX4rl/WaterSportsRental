package handler

import (
	"WaterSportsRental/internal/filter"
	accessToken "WaterSportsRental/internal/jwt"
	"WaterSportsRental/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.URLFormat)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*", "*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Route("/auth", func(r chi.Router) { //Аутентификация
		r.Post("/sign-up", h.signUp)
		r.Post("/sign-in", h.signIn)
		//r.Post("/sign-in", h.UpdateUser)
	})

	router.Route("/profile", func(r chi.Router) {
		r.Use(accessToken.JwtAuthMiddleware())
		r.Post("/update", h.UpdateUser)
		r.Get("/user", h.UserInfo)
	})

	router.Get("/allLocations", h.GetAllLocations)

	router.Route("/products", func(r chi.Router) {
		r.Route("/rent", func(r chi.Router) {
			r.Use(accessToken.JwtAuthMiddleware())
			r.Post("/", h.rent)
		})
		r.Route("/info", func(r chi.Router) { //информация об истории заказов
			r.Use(accessToken.JwtAuthMiddleware())
			r.Get("/", h.getRentedItems)
		})

		r.Route("/review", func(r chi.Router) {
			r.Use(accessToken.JwtAuthMiddleware())
			r.Post("/", h.postReview)        //оставть отзыв о товаре, учитывается его айди]
			r.Get("/", h.getReviewsByUserId) //отзывы в профиле ???
		})

		r.Get("/locations", h.getAvailableProducts)           //фильрация по локациям
		r.Get("/", filter.Middleware(h.getAllEquipments, 10)) //динамическая фильтрация
		r.Get("/reviews", h.getReviewsByItemId)               //отзывы по конкретному товару (учитывается его айди)
		r.Get("/dates", h.getAvailableDates)                  //период, когда товар доступен

	})

	router.Route("/admin", func(r chi.Router) {
		r.Use(accessToken.JwtAuthMiddleware())
		r.Use(AdminMiddleware(h))
		r.Get("/permission", h.hasPermission)
		r.Route("/locations", func(r chi.Router) {
			r.Get("/", h.GetAllLocations)
			r.Post("/", h.CreateLocation)
			r.Post("/update", h.UpdateLocationInfo)
			r.Delete("/", h.DeleteLocation)
			r.Delete("/deleteAll", h.DeleteAllLocations)
		})
		r.Route("/users", func(r chi.Router) {
			r.Get("/", h.GetUserById)
			r.Get("/all", h.GetAllUsers)
			r.Post("/", h.CreateUser)
			r.Post("/update", h.UpdateWithRole)
			r.Post("/role", h.GiveRole)
			r.Delete("/", h.DeleteUserById)
			r.Get("/agreements", h.GetAllRentalAgreementsByID)
		})
		r.Route("/items", func(r chi.Router) {
			r.Post("/", h.Create)
			r.Get("/", filter.Middleware(h.getAllEquipments, 10))
			r.Post("/update", h.UpdateItem)
			r.Delete("/", h.DeleteItem)
			r.Post("/add", h.addProduct)
		})
		r.Route("/reviews", func(r chi.Router) {
			r.Get("/", h.getAllReviews)
			r.Delete("/", h.DeleteReviewById)
		})
		r.Route("/database", func(r chi.Router) {
			r.Get("/dump", h.createDump)
			r.Post("/restore", h.restoreDump)
			r.Get("/all", h.getAllDumps)
		})
		r.Get("/logs", h.getLogs)
	})

	return router
}
