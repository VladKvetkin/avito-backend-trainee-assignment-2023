package server

import (
	"net/http"

	"github.com/VladKvetkin/avito-backend-trainee-assignment-2023/internal/handler"

	"github.com/go-chi/chi"
)

func (s *Server) setupRoutes(handler *handler.Handler) {
	s.mux.Route("/", func(r chi.Router) {
		r.Route("/api", func(r chi.Router) {
			r.Route("/user", func(r chi.Router) {
				r.Post("/change-segments", http.HandlerFunc(handler.ChangeSegments))
				r.Post("/segments", http.HandlerFunc(handler.GetSegments))
				r.Post("/segments-history", http.HandlerFunc(handler.GetSegmentsHistoryReportURL))
			})

			r.Route("/segment", func(r chi.Router) {
				r.Post("/add", http.HandlerFunc(handler.CreateSegment))
				r.Post("/delete", http.HandlerFunc(handler.DeleteSegment))
				r.Get("/report/*", http.HandlerFunc(handler.GetSegmentsHistoryReportFile))
			})
		})
	})
}
