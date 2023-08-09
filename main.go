package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"instantmsg/models"
	"instantmsg/server"
	"instantmsg/server/routes"
	"os"
)

func main() {

	db, err := gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/chat?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		fmt.Printf("Erreur opening database%v\n", err)
		os.Exit(1)
	}

	ctx := &models.AppContext{
		DB: db,
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		fmt.Printf("Erreur lors de l'automigration %v\n", err)
	}
	s := server.NewServer()
	ur := routes.NewUserRoute(ctx)
	s.Router.AddRoutes(ur.GetRoutes())
	s.Router.RegisterRoutes()
	s.Run()

}
