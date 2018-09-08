package handlers

import (
	"hlcup/services"
	"log"
	"github.com/valyala/fasthttp"
	"github.com/asaskevich/govalidator"
	"strconv"
	"time"
	"os"
	"bufio"
	"encoding/json"
)

type UserApiHandler struct {
	storage            *services.Storage
	errLogger          *log.Logger
	infoLogger         *log.Logger
	timeDataGeneration time.Time
}

func NewUserApiHandler(storage *services.Storage, errLogger *log.Logger, infoLogger *log.Logger, pathToOptions string) *UserApiHandler {

	file, _ := os.Open(pathToOptions)

	fileScanner := bufio.NewScanner(file)

	fileScanner.Scan()

	timeDataGeneration, err := strconv.Atoi(fileScanner.Text())

	if err != nil {
		errLogger.Fatalln(err)
	}

	return &UserApiHandler{storage: storage, errLogger: errLogger, infoLogger: infoLogger, timeDataGeneration: time.Unix(int64(timeDataGeneration), 0)}
}

func (userApiHandler *UserApiHandler) GetById(ctx *fasthttp.RequestCtx) {

	userIdString, ok := ctx.UserValue("user_id").(string)

	if !govalidator.IsNumeric(userIdString) || !ok {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.SetConnectionClose()
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		return
	}

	userId, err := strconv.Atoi(userIdString)

	if err != nil {
		userApiHandler.errLogger.Fatalln(err)
	}

	user := userApiHandler.storage.GetUserById(uint(userId))

	if user == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.SetConnectionClose()
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.Response.Header.SetContentType("application/json")
		ctx.Response.Header.SetConnectionClose()
		ctx.Response.Header.SetContentLength(len(user))
		ctx.Write(user)
	}
}

func (userApiHandler *UserApiHandler) GetVisitedPlaces(ctx *fasthttp.RequestCtx) {

	userIdString, ok := ctx.UserValue("user_id").(string)

	if !govalidator.IsNumeric(userIdString) || !ok {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.SetConnectionClose()
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		return
	}

	userId, err := strconv.Atoi(userIdString)

	if err != nil {
		userApiHandler.errLogger.Fatalln(err)
	}

	filter := services.InitVisitFilter(userApiHandler.timeDataGeneration)

	userIdUint := uint(userId)

	filter.UserId = &userIdUint

	if ctx.QueryArgs().Has("fromDate") {

		fromDateString := string(ctx.QueryArgs().Peek("fromDate"))

		if !govalidator.IsNumeric(fromDateString) || fromDateString == "" {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.Response.Header.SetContentType("text/plain; charset=utf8")
			ctx.Response.Header.SetConnectionClose()
			ctx.Response.Header.SetContentLength(len("Bad Request"))
			ctx.WriteString("Bad Request")
			return
		}

		fromDateInt, err := strconv.Atoi(fromDateString)

		if err != nil {
			userApiHandler.errLogger.Println(err)
		}

		filter.FromDate = &fromDateInt
	}

	if ctx.QueryArgs().Has("toDate") {

		toDateString := string(ctx.QueryArgs().Peek("toDate"))

		if !govalidator.IsNumeric(toDateString) || toDateString == "" {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.Response.Header.SetContentType("text/plain; charset=utf8")
			ctx.Response.Header.SetConnectionClose()
			ctx.Response.Header.SetContentLength(len("Bad Request"))
			ctx.WriteString("Bad Request")
			return
		}

		toDateInt, err := strconv.Atoi(toDateString)

		if err != nil {
			userApiHandler.errLogger.Println(err)
		}

		filter.ToDate = &toDateInt
	}

	if ctx.QueryArgs().Has("toDistance") {

		toDistanceString := string(ctx.QueryArgs().Peek("toDistance"))

		if !govalidator.IsNumeric(toDistanceString) || toDistanceString == "" {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.Response.Header.SetContentType("text/plain; charset=utf8")
			ctx.Response.Header.SetConnectionClose()
			ctx.Response.Header.SetContentLength(len("Bad Request"))
			ctx.WriteString("Bad Request")
			return
		}

		toDistanceInt, err := strconv.Atoi(toDistanceString)

		if err != nil {
			userApiHandler.errLogger.Println(err)
		}

		toDistanceUint := uint(toDistanceInt)

		filter.ToDistance = &toDistanceUint
	}

	if ctx.QueryArgs().Has("country") {

		country := string(ctx.QueryArgs().Peek("country"))

		if country == "" || len(country) > 50 {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.Response.Header.SetContentType("text/plain; charset=utf8")
			ctx.Response.Header.SetConnectionClose()
			ctx.Response.Header.SetContentLength(len("Bad Request"))
			ctx.WriteString("Bad Request")
			return
		}

		filter.Country = &country
	}

	visitedPlaceCollection := userApiHandler.storage.GetVisitedPlacesByUser(filter)

	if visitedPlaceCollection == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.SetConnectionClose()
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		return
	}

	visitedPlaceCollectionBytes, err := json.Marshal(visitedPlaceCollection)

	if err != nil {
		userApiHandler.errLogger.Fatalln(err)
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetConnectionClose()
	ctx.Response.Header.SetContentLength(len(visitedPlaceCollectionBytes))
	ctx.Write(visitedPlaceCollectionBytes)
}
