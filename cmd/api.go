package cmd

import (
	"context"
	"log"
	"net/http"

	"kmed/api/internal/auth"
	"kmed/api/internal/config"
	"kmed/api/internal/consultation"
	"kmed/api/internal/database"
	"kmed/api/internal/icd"
	"kmed/api/internal/patient"
	"kmed/api/internal/user"

	_ "github.com/mattn/go-sqlite3"
)

type Api struct {
	database *database.Database
	Router   *http.ServeMux
}

func NewApi() *Api {
	configStore := config.NewConfigStore("kmed")
	config, err := configStore.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	database, err := database.OpenDatabase(config.DBPath)
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}

	if err := database.Migrate(context.Background()); err != nil {
		log.Fatalf("failed migrating the database: %v", err)
	}

	authMiddleware := auth.NewMiddleware(config.AuthSecretKey)

	userHandler := user.NewHandler(database.Client, authMiddleware)
	authHandler := auth.NewHandler(database.Client, authMiddleware, config.AuthSecretKey)
	patientHandler := setupPatientHandler(database, authMiddleware)
	consultationHandler := consultation.NewHandler(database.Client, authMiddleware)
	icd10Handler := setupICD10Handler(database, authMiddleware)

	router := http.NewServeMux()
	router.Handle("/users/", logRequest(http.StripPrefix("/users", userHandler.Router())))
	router.Handle("/auth/", logRequest(http.StripPrefix("/auth", authHandler.Router())))
	router.Handle("/patients/", logRequest(http.StripPrefix("/patients", patientHandler.Router())))
	router.Handle(
		"/consultations/",
		logRequest(http.StripPrefix("/consultations", consultationHandler.Router())),
	)
	router.Handle("/icd/", logRequest(http.StripPrefix("/icd", icd10Handler.Router())))

	return &Api{database: database, Router: router}
}

func setupPatientHandler(
	database *database.Database,
	authMiddleware auth.Middleware,
) *patient.Handler {
	patientService := patient.NewService(database)
	patientHandler := patient.NewHandler(patientService, authMiddleware)
	return patientHandler
}

func setupICD10Handler(
	db *database.Database,
	authMiddleware auth.Middleware,
) *icd.Handler {
	icdService := icd.NewService(db)
	icd10Handler := icd.NewHandler(icdService, authMiddleware)
	return icd10Handler
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
