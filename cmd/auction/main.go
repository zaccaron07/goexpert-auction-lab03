package main

import (
	"context"
	"log"

	"github.com/zaccaron07/goexpert-auction-lab03/configuration/database/mongodb"
	"github.com/zaccaron07/goexpert-auction-lab03/internal/infra/api/web/controller/auction_controller"
	"github.com/zaccaron07/goexpert-auction-lab03/internal/infra/api/web/controller/bid_controller"
	"github.com/zaccaron07/goexpert-auction-lab03/internal/infra/api/web/controller/user_controller"
	"github.com/zaccaron07/goexpert-auction-lab03/internal/infra/database/auction"
	"github.com/zaccaron07/goexpert-auction-lab03/internal/infra/database/bid"
	"github.com/zaccaron07/goexpert-auction-lab03/internal/infra/database/user"
	"github.com/zaccaron07/goexpert-auction-lab03/internal/usecase/auction_usecase"
	"github.com/zaccaron07/goexpert-auction-lab03/internal/usecase/bid_usecase"
	"github.com/zaccaron07/goexpert-auction-lab03/internal/usecase/user_usecase"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
99c5ab02-dcfc-470b-a03e-fb1fc0078857
5421869b-dfda-4028-b04c-a695667c77d8
*/
func main() {
	ctx := context.Background()

	if err := godotenv.Load("cmd/auction/.env"); err != nil {
		log.Fatal("Error trying to load env variables")
		return
	} else {
		log.Println("Env variables loaded successfully")
	}

	databaseConnection, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return
	} else {
		log.Println("MongoDB connection established successfully")
	}

	router := gin.Default()

	userController, bidController, auctionsController := initDependencies(databaseConnection)

	router.GET("/auction", auctionsController.FindAuctions)
	router.GET("/auction/:auctionId", auctionsController.FindAuctionById)
	router.POST("/auction", auctionsController.CreateAuction)
	router.GET("/auction/winner/:auctionId", auctionsController.FindWinningBidByAuctionId)
	router.POST("/bid", bidController.CreateBid)
	router.GET("/bid/:auctionId", bidController.FindBidByAuctionId)
	router.GET("/user/:userId", userController.FindUserById)

	router.Run("localhost:8080")
}

func initDependencies(database *mongo.Database) (
	userController *user_controller.UserController,
	bidController *bid_controller.BidController,
	auctionController *auction_controller.AuctionController) {

	auctionRepository := auction.NewAuctionRepository(database)
	bidRepository := bid.NewBidRepository(database, auctionRepository)
	userRepository := user.NewUserRepository(database)

	userController = user_controller.NewUserController(
		user_usecase.NewUserUseCase(userRepository))
	auctionController = auction_controller.NewAuctionController(
		auction_usecase.NewAuctionUseCase(auctionRepository, bidRepository))
	bidController = bid_controller.NewBidController(bid_usecase.NewBidUseCase(bidRepository))

	return
}
