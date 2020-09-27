package entities

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"time"
)

// DBConfig holds DB credentials
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// ConnectToDB establishes a connection with the DB
func ConnectToDB(c *DBConfig) (*gorm.DB, error) {
	connectStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", c.Host, c.Port, c.User, c.DBName, c.Password)
	log.Info(connectStr)
	db, err := gorm.Open("postgres", connectStr)
	if err != nil {
		return nil, errors.Wrap(err, "error connecting to the database")
	}
	MigrateModels(db)
	return db, nil
}

// MigrateModels migrates declared models
func MigrateModels(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Branch{})
	db.AutoMigrate(&Plan{})
	db.AutoMigrate(&Subscription{})
	db.AutoMigrate(&Offer{})
	db.AutoMigrate(&Review{})
	db.AutoMigrate(&PartnerProfile{})
	db.AutoMigrate(&PartnerPhoto{})
	db.AutoMigrate(&OTP{})
	db.AutoMigrate(&Category{})
	db.AutoMigrate(&City{})
	db.AutoMigrate(&SupportInfo{})
	db.AutoMigrate(&Exclusive{})
	db.AutoMigrate(&Share{})
	db.AutoMigrate(&PlanCategory{})
	db.AutoMigrate(&CustomerPartnerOffersCount{})
	Seed(db)
}

// Seed seeds the database with basic admin records
func Seed(db *gorm.DB) {
	istanbul := &City{
		EnglishName: "Istanbul",
		TurkishName: "Istanbul",
	}
	db.Where(istanbul).FirstOrCreate(&istanbul)
	hashedPass, _ := EncryptPassword("12345678")
	admin := &User{
		FirstName:       "super",
		LastName:        "admin",
		Email:           "superadmin@tasarruf-admin.com",
		HashedPassword:  hashedPass,
		Verified:        true,
		Mobile:          "1234567899",
		AccountType:     "admin",
		Country:         "turkey",
		CityID:          istanbul.ID,
		DateOfBirth:     time.Now(),
		ProfileImageURL: "https://tasarruf-file-repository.s3.amazonaws.com/profile_images/default.png",
	}
	db.Where(admin).FirstOrCreate(&istanbul)
}
