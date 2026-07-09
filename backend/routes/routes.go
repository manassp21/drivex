package routes

import(
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes() http.Handler{
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("DriveX api is running"))
	})

	r.Route("/auth", func(r chi.Router){

	})

	r.Route("/files", func(r chi.Router){
		
	})

	r.Route("/folders", func(r chi.Router){
		
	})

	return r

}