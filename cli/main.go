package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/hughselway/m-bff/bff"
	"github.com/rs/zerolog/log"
)

func main() {
	grpcAddressHighScore := flag.String("address-m-highscore", "localhost:50051", "The grpc server address for highscore service")
	grpcAddressGameEngine := flag.String("address-m-game-engine", "localhost:60051", "The grpc server address for game-engine service")
	grpcAddressBorderChange := flag.String("address-m-border-change", "localhost:50052", "The grpc server address for border-change service")

	serverAddress := flag.String("address-http", ":8081", "HTTP server address")

	flag.Parse()

	gameClient, err := bff.NewGrpcGameServiceClient(*grpcAddressHighScore)
	if err != nil {
		log.Error().Err(err).Msg("Error in creating a client for m-highscore")
	}

	gameEngineClient, err := bff.NewGrpcGameEngineServiceClient(*grpcAddressGameEngine)
	if err != nil {
		log.Error().Err(err).Msg("Error in creating a client for m-game-engine")
	}

	gameBorderClient, err := bff.NewGrpcGameBorderServiceClient(*grpcAddressBorderChange)
	if err != nil {
		log.Error().Err(err).Msg("Error in creating a client for m-border-change")
	}

	gr := bff.NewGameResource(gameClient, gameEngineClient, gameBorderClient)

	router := gin.Default()
	router.GET("/geths", gr.GetHighScore)
	router.GET("/seths/:hs", gr.SetHighScore)
	router.GET("/getsize", gr.GetSize)
	router.GET("/setscore/:score", gr.SetScore)
	router.GET("/borderchange/st/:st/wd/:wd/scred/:scred/scgreen/:scgreen/scblue/:scblue", gr.GetBorder)

	err = router.Run(*serverAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not start bff")
	}

	log.Info().Msgf("Started http-server at %v", *serverAddress)

}
