package main

import (
	"bookcabin-test/backend/config"
	"bookcabin-test/backend/internal/domain"
	"bookcabin-test/backend/internal/handler"
	"bookcabin-test/backend/internal/repository"
	"bookcabin-test/backend/internal/usecase"
	"log"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func main() {
	// init DB conennection
	db, err := config.NewDB()
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}
	db.AutoMigrate(&domain.Voucher{})

	// init repositories
	voucherRepo := repository.NewVoucherRepository(db)

	//init usecases
	voucherUC := usecase.NewVoucherUsecase(voucherRepo)
	aircraftUC := usecase.NewAircraftUsecase()

	// init handlers
	voucherHandler := handler.NewVoucherHandler(voucherUC)
	aircraftHandler := handler.NewAircraftHandler(aircraftUC)

	// init gin router
	r := gin.Default()
	r.Use(cors.Default())
	voucherHandler.RegisterRoutes(r)
	aircraftHandler.RegisterRoutes(r)
	r.Run()
}
