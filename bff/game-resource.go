package bff

import (
	"context"
	"github.com/gin-gonic/gin"
	pbborderchange "github.com/hughselway/m-apis/m-borderchange"
	pbgameengine "github.com/hughselway/m-apis/m-game-engine/v1"
	pbhighscore "github.com/hughselway/m-apis/m-highscore/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"strconv"
)

type gameResource struct {
	gameClient       pbhighscore.GameClient
	gameEngineClient pbgameengine.GameEngineClient
	gameBorderClient pbborderchange.GameBorderClient
}

func NewGameResource(gameClient pbhighscore.GameClient, gameEngineClient pbgameengine.GameEngineClient, gameBorderClient pbborderchange.GameBorderClient) *gameResource {
	return &gameResource{
		gameClient:       gameClient,
		gameEngineClient: gameEngineClient,
		gameBorderClient: gameBorderClient,
	}
}

func NewGrpcGameServiceClient(serverAddr string) (pbhighscore.GameClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())

	if err != nil {
		log.Fatal().Msgf("Failed to dial: %v", err)
		return nil, err
	} else {
		log.Info().Msgf("Successfully connected to [%s]", serverAddr)
	}

	if conn == nil {
		log.Info().Msg("m-highscore connection is nil in m-bff")
	}

	client := pbhighscore.NewGameClient(conn)

	return client, nil
}

func NewGrpcGameEngineServiceClient(serverAddr string) (pbgameengine.GameEngineClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())

	if err != nil {
		log.Fatal().Msgf("Failed to dial: %v", err)
		return nil, err
	} else {
		log.Info().Msgf("Successfully connected to [%s]", serverAddr)
	}

	if conn == nil {
		log.Info().Msg("m-game-engine connection is nil in m-bff")
	}

	client := pbgameengine.NewGameEngineClient(conn)

	return client, nil
}

func NewGrpcGameBorderServiceClient(serverAddr string) (pbborderchange.GameBorderClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())

	if err != nil {
		log.Fatal().Msgf("Failed to dial: %v", err)
		return nil, err
	} else {
		log.Info().Msgf("Successfully connected to [%s]", serverAddr)
	}

	if conn == nil {
		log.Info().Msg("m-borderchange connection is nil in m-bff")
	}

	client := pbborderchange.NewGameBorderClient(conn)

	return client, nil
}

func (gr *gameResource) SetHighScore(c *gin.Context) {
	highScoreString := c.Param("hs")
	highScoreFloat64, err := strconv.ParseFloat(highScoreString, 64)
	if err != nil {
		log.Error().Err(err).Msg("Failed to convert highscore to float")
	}
	_, err = gr.gameClient.SetHighScore(context.Background(), &pbhighscore.SetHighScoreRequest{
		HighScore: highScoreFloat64,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to call SetHighScore")
	}
}

func (gr *gameResource) GetHighScore(c *gin.Context) {

	highScoreResponse, err := gr.gameClient.GetHighScore(context.Background(), &pbhighscore.GetHighScoreRequest{})
	if err != nil {
		log.Error().Err(err).Msg("Error while getting highscore")
		return
	}
	hsString := strconv.FormatFloat(highScoreResponse.HighScore, 'e', -1, 64)

	c.JSONP(200, gin.H{
		"hs": hsString,
	})

}

func (gr *gameResource) GetSize(c *gin.Context) {
	sizeResponse, err := gr.gameEngineClient.GetSize(context.Background(), &pbgameengine.GetSizeRequest{})
	if err != nil {
		log.Error().Err(err).Msg("Error while getting size")
	}
	c.JSON(200, gin.H{
		"size": sizeResponse.GetSize(),
	})
}

func (gr *gameResource) SetScore(c *gin.Context) {
	scoreString := c.Param("score")
	scoreFloat64, _ := strconv.ParseFloat(scoreString, 64)

	_, err := gr.gameEngineClient.SetScore(context.Background(), &pbgameengine.SetScoreRequest{
		Score: scoreFloat64,
	})
	if err != nil {
		log.Error().Err(err).Msg("Error while setting score in m-game-engine")
	}
}

func (gr *gameResource) GetBorder(c *gin.Context) {
	style := c.Param("st")
	width := c.Param("wd")
	shapered, err := strconv.ParseFloat(c.Param("scred"), 64)
	if err != nil {
		log.Error().Err(err).Msg("failed to convert red colour to float64")
	}
	shapegreen, err := strconv.ParseFloat(c.Param("scgreen"), 64)
	if err != nil {
		log.Error().Err(err).Msg("failed to convert green colour to float64")
	}
	shapeblue, err := strconv.ParseFloat(c.Param("scblue"), 64)
	if err != nil {
		log.Error().Err(err).Msg("failed to convert blue colour to float64")
	}

	getborderresponse, err := gr.gameBorderClient.GetBorder(context.Background(), &pbborderchange.GetBorderRequest{
		//get stuff from the c context as in other function
		Style:           style,
		Width:           width,
		ShapeColorRed:   shapered,
		ShapeColorGreen: shapegreen,
		ShapeColorBlue:  shapeblue,
	})
	if err != nil {
		log.Error().Err(err).Msg("error in changing border from m-borderchange")
	}

	crString := strconv.FormatFloat(getborderresponse.BorderColorRed, 'e', -1, 64)
	cgString := strconv.FormatFloat(getborderresponse.BorderColorGreen, 'e', -1, 64)
	cbString := strconv.FormatFloat(getborderresponse.BorderColorBlue, 'e', -1, 64)

	c.JSONP(200, gin.H{
		"style":   getborderresponse.Style,
		"width":   getborderresponse.Width,
		"scred":   crString,
		"scgreen": cgString,
		"scblue":  cbString,
	})
}
