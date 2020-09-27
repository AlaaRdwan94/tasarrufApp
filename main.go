// Copyright 2019 NOVA Solutions Co. All Rights Reserved.
//

package main

import (
	"flag"
	"fmt"
	"github.com/ahmedaabouzied/tasarruf/entities"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

// Server configuration
type Server struct {
	Port    string        // Port number
	Env     string        // Server environment dev, staging, prod , win
	Timeout time.Duration // Server timeout
	Router  *gin.Engine   // Router
	DB      *gorm.DB      // Gorm DB Connection
}

// ParseENV loads environment variables
func ParseENV() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Error("error parsing ENV variables", err)
	}
	log.Info("=== Loaded environment variables")
}

// CreateDBConfig return a database credentials object
func CreateDBConfig(env string) (*entities.DBConfig, error) {
	port := 5432
	var host, user, password, dbname string
	switch env {
	case "dev":
		host = os.Getenv("DBHOST")
		user = os.Getenv("DBUSER")
		password = os.Getenv("DEVDBPASSWORD")
		dbname = os.Getenv("DBNAME")
	case "staging":
		host = os.Getenv("DBHOST")
		user = os.Getenv("DBUSER")
		password = os.Getenv("DBPASSWORD")
		dbname = os.Getenv("DBNAME")
	case "prod":
		host = os.Getenv("DBHOST")
		user = os.Getenv("DBUSER")
		password = os.Getenv("DBPASSWORD")
		dbname = os.Getenv("DBNAME")
	default:
		host = os.Getenv("DBHOST")
		user = os.Getenv("DBUSER")
		password = os.Getenv("DBPASSWORD")
		dbname = os.Getenv("DBNAME")
	}
	dbC := &entities.DBConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbname,
	}
	log.Info("=== Parsed DB Credentials")
	return dbC, nil
}

// StartServer starts the server or fails with an error
func StartServer(s *Server, debug bool) error {
	srv := &http.Server{
		Handler:      s.Router,
		Addr:         fmt.Sprintf(":%s", s.Port),
		WriteTimeout: s.Timeout,
		ReadTimeout:  s.Timeout,
	}

	InitializeRoutes(s.DB, s.Router)
	if debug {
		log.Info("=== Starting server on port " ,s.Port)
	}

	return srv.ListenAndServe()
}

func main() {
	env := flag.String("env", "dev", "sets the running environment to either dev, staging or prod")
	flag.Parse()
	r := gin.Default() // default gin router
	ParseENV()
	dbConfig, err := CreateDBConfig(*env)
	db, err := entities.ConnectToDB(dbConfig)
	if err != nil {
		log.Fatal(err)
	}
	s := Server{
		Port:    os.Getenv("PORT"),
		Timeout: 15 * time.Second,
		Router:  r,
		Env:     "debug",
		DB:      db,
	}
	debug := s.Env == "debug"
	log.Fatal(StartServer(&s, debug))
}
