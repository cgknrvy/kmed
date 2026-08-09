package cmd

import (
	"context"
	"log"
	"net/http"

	"kmed/api/internal/auth"
	"kmed/api/internal/consultation"
	"kmed/api/internal/migration"
	"kmed/api/internal/patient"
	"kmed/api/internal/user"

	_ "github.com/mattn/go-sqlite3"
)

type Api struct {
	database *migration.Database
	Router   *http.ServeMux
}

func NewApi() *Api {
	database, err := migration.OpenDatabase("./kmed.db?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}

	if err := database.Migrate(context.Background()); err != nil {
		log.Fatalf("failed migrating the database: %v", err)
	}

	authMiddleware := auth.NewMiddleware()

	userHandler := user.NewHandler(database.Client, authMiddleware)
	authHandler := auth.NewHandler(database.Client)
	patientHandler := setupPatientHandler(database, authMiddleware)
	consultationHandler := consultation.NewHandler(database.Client, authMiddleware)

	router := http.NewServeMux()
	router.Handle("/users/", logRequest(http.StripPrefix("/users", userHandler.Router())))
	router.Handle("/auth/", logRequest(http.StripPrefix("/auth", authHandler.Router())))
	router.Handle("/patients/", logRequest(http.StripPrefix("/patients", patientHandler.Router())))
	router.Handle(
		"/consultations/",
		logRequest(http.StripPrefix("/consultations", consultationHandler.Router())),
	)

	return &Api{database: database, Router: router}
}

func setupPatientHandler(
	database *migration.Database,
	authMiddleware auth.Middleware,
) *patient.Handler {
	patientService := patient.NewService(database)
	patientHandler := patient.NewHandler(patientService, authMiddleware)
	return patientHandler
}

func (api *Api) Close() error {
	// database.DB is closed by the ent.Client.Close() too
	return api.database.Client.Close()
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
