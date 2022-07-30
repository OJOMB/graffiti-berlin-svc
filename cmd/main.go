package main

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/OJOMB/graffiti-berlin-svc/internal/app"
	"github.com/OJOMB/graffiti-berlin-svc/internal/pkg/auth"
	"github.com/OJOMB/graffiti-berlin-svc/internal/pkg/domain"
	"github.com/OJOMB/graffiti-berlin-svc/internal/pkg/passwords"
	"github.com/OJOMB/graffiti-berlin-svc/internal/pkg/repo"
	"github.com/OJOMB/graffiti-berlin-svc/internal/pkg/uuidv4"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
)

const (
	versionEnv     = "SVC_VERSION"
	portEnv        = "SVC_PORT"
	environmentEnv = "SVC_ENVIRONMENT"
	dbHostEnv      = "DB_HOST"
	dbPortEnv      = "DB_PORT"
	dbUserEnv      = "DB_USER"
	dbPasswordEnv  = "DB_PASSWORD"

	defaultVersion     = "v0.0.0"
	defaultPort        = 8080
	defaultHost        = "0.0.0.0"
	defaultEnvironment = "dev"
	defaultDBPort      = 3306
	defaultDBUser      = "root"
	defaultDBPassword  = "pass"

	dbName  = "graffiti"
	appName = "graffiti-berlin-svc"

	passwordGeneratorCost = 15
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	version := os.Getenv(versionEnv)
	if version == "" {
		logger.Info("failed to retrieve app version number from env...using default")
		version = defaultVersion
	}

	logger.Infof("app version number %s", version)

	var port int
	var err error
	portStr := os.Getenv(portEnv)
	if portStr == "" {
		logger.Infof("failed to retrieve app port number from env...using default %d", defaultPort)
		port = defaultPort
	} else {
		port, err = strconv.Atoi(portStr)
		if err != nil {
			logger.WithError(err).Fatalf("retrieved invalid service port number from env: %s", portStr)
		}
	}

	logger.Infof("app port number %d", port)

	environment := os.Getenv(environmentEnv)
	if environment == "" {
		logger.Info("failed to retrieve app environment from env...using default")
		environment = defaultEnvironment
	}

	var docsEnabled bool
	if environment != "prod" {
		logger.Info("app docs endpoint enabled")
		// swaggerUI should not be available in prod environment
		docsEnabled = true
	}

	dbHost := os.Getenv(dbHostEnv)
	if version == "" {
		logger.Info("failed to retrieve DB host from env...using default")
		dbHost = defaultHost
	}

	var dbPort int
	dbPortStr := os.Getenv(dbPortEnv)
	if dbPortStr == "" {
		logger.Info("failed to retrieve DB host from env...using default")
		dbPort = defaultDBPort
	} else {
		dbPort, err = strconv.Atoi(dbPortStr)
		if err != nil {
			logger.WithError(err).Fatalf("retrieved invalid DB port number from env: %s", portStr)
		}
	}

	dbUser := os.Getenv(dbUserEnv)
	if dbUser == "" {
		logger.Info("failed to retrieve DB host from env...using default")
		dbUser = defaultDBUser
	}

	dbPassword := os.Getenv(dbPasswordEnv)
	if dbUser == "" {
		logger.Info("failed to retrieve DB host from env...using default")
		dbPassword = defaultDBPassword
	}

	logger.Infof("connecting to DB @ %s:%d as %s", dbHost, dbPort, dbUser)

	// default loc - so UTC
	dbCnxnStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dbCnxnStr)
	if err != nil {
		logger.WithError(err).Fatalf("failed to establish connection to DB @ %s:%d", dbHost, dbPort)
	}

	defer db.Close()

	logger.Info("successfully connected to DB")

	server := app.New(
		mux.NewRouter(),
		logger, &net.TCPAddr{IP: net.ParseIP(defaultHost), Port: port},
		version,
		environment,
		auth.NewJWTTool("secretKey", 12*time.Hour, appName),
		domain.NewService(
			logger,
			repo.NewSQLRepo(db, logger),
			uuidv4.NewGenerator(),
			passwords.NewGenerator(passwordGeneratorCost),
		),
	)

	server.Run()
}
