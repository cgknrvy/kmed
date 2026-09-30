package cmd

import (
	"context"
	"log/slog"
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

func NewApi(cfg *config.Config) *Api {
	database, err := database.OpenDatabase(cfg.DBPath)
	if err != nil {
		slog.Error("sqlite connection failed", "error", err)
		panic(err)
	}

	if err := database.Migrate(context.Background()); err != nil {
		slog.Error("migrating database failed", "err", err)
		panic(err)
	}

	authMiddleware := auth.NewMiddleware(cfg.AuthSecretKey)

	userHandler := user.NewHandler(database.Client, authMiddleware)
	authHandler := auth.NewHandler(database.Client, authMiddleware, cfg.AuthSecretKey)
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
		slog.Info("request", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
