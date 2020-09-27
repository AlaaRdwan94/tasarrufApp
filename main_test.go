// Copyright 2019 NOVA Solutions Co. All Rights Reserved.
//

package main

import (
	"github.com/ahmedaabouzied/tasarruf/entities"
	"testing"
)

func TestConnectToDB(t *testing.T) {
	conf := entities.DBConfig{
		Port:     5432,
		Host:     "localhost",
		User:     "tasarruf",
		Password: "password",
		DBName:   "tasarruftestdb",
	}
	db, err := entities.ConnectToDB(&conf)
	if err != nil {
		t.Fail()
	}
	db.Close()
}

func TestParseENV(t *testing.T) {
	ParseENV()
}

func TestCreateDBConfig(t *testing.T) {
	ParseENV()
	_, err := CreateDBConfig()
	if err != nil {
		t.Error(err)
	}
}
