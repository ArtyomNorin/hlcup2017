package handlers

import (
	"hlcup/services"
	"log"
	"github.com/valyala/fasthttp"
	"github.com/asaskevich/govalidator"
	"strconv"
	"os"
	"bufio"
	"time"
	"hlcup/entities"
	"math"
	"encoding/json"
)

type LocationApiHandler struct {
	storage            *services.Storage
	errLogger          *log.Logger
	infoLogger         *log.Logger
	timeDataGeneration time.Time
}

func NewLocationApiHandler(storage *services.Storage, errLogger *log.Logger, infoLogger *log.Logger, pathToOptions string) *LocationApiHandler {

	file, _ := os.Open(pathToOptions)

	fileScanner := bufio.NewScanner(file)

	fileScanner.Scan()

	timeDataGeneration, err := strconv.Atoi(fileScanner.Text())

	if err != nil {
		errLogger.Fatalln(err)
	}

	return &LocationApiHandler{storage: storage, errLogger: errLogger, infoLogger: infoLogger, timeDataGeneration: time.Unix(int64(timeDataGeneration), 0)}
}

func (locationApiHandler *LocationApiHandler) GetById(ctx *fasthttp.RequestCtx) {

	locationIdString, ok := ctx.UserValue("location_id").(string)

	if !govalidator.IsNumeric(locationIdString) || !ok{
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.SetConnectionClose()
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		return
	}

	locationId, err := strconv.Atoi(locationIdString)

	if err != nil {
		locationApiHandler.errLogger.Fatalln(err)
	}

	location := locationApiHandler.storage.GetLocationById(uint(locationId))

	if location == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.SetConnectionClose()
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.Response.Header.SetContentType("application/json")
		ctx.Response.Header.SetConnectionClose()
		ctx.Response.Header.SetContentLength(len(location))
		ctx.Write(location)
	}
}

func (locationApiHandler *LocationApiHandler) GetAverageMark(ctx *fasthttp.RequestCtx) {

	locationIdString, ok := ctx.UserValue("location_id").(string)

	if !govalidator.IsNumeric(locationIdString) || !ok{
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.SetConnectionClose()
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		return
	}

	locationIdInt, err := strconv.Atoi(locationIdString)

	if err != nil {
		locationApiHandler.errLogger.Fatalln(err)
	}

	filter := services.InitVisitFilter(locationApiHandler.timeDataGeneration)

	locationIdUint := uint(locationIdInt)

	filter.LocationId = &locationIdUint

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
			locationApiHandler.errLogger.Println(err)
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
			locationApiHandler.errLogger.Println(err)
		}

		filter.ToDate = &toDateInt
	}

	if ctx.QueryArgs().Has("fromAge") {

		fromAgeString := string(ctx.QueryArgs().Peek("fromAge"))

		if !govalidator.IsNumeric(fromAgeString) || fromAgeString == "" {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.Response.Header.SetContentType("text/plain; charset=utf8")
			ctx.Response.Header.SetConnectionClose()
			ctx.Response.Header.SetContentLength(len("Bad Request"))
			ctx.WriteString("Bad Request")
			return
		}

		fromAgeInt, err := strconv.Atoi(fromAgeString)

		if err != nil {
			locationApiHandler.errLogger.Println(err)
		}

		filter.FromAge = &fromAgeInt
	}

	if ctx.QueryArgs().Has("toAge") {

		toAgeString := string(ctx.QueryArgs().Peek("toAge"))

		if !govalidator.IsNumeric(toAgeString) || toAgeString == "" {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.Response.Header.SetContentType("text/plain; charset=utf8")
			ctx.Response.Header.SetConnectionClose()
			ctx.Response.Header.SetContentLength(len("Bad Request"))
			ctx.WriteString("Bad Request")
			return
		}

		toAgeInt, err := strconv.Atoi(toAgeString)

		if err != nil {
			locationApiHandler.errLogger.Println(err)
		}

		filter.ToAge = &toAgeInt
	}

	if ctx.QueryArgs().Has("gender") {

		gender := string(ctx.QueryArgs().Peek("gender"))

		if gender == "" || (gender != "m" && gender != "f") {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.Response.Header.SetContentType("text/plain; charset=utf8")
			ctx.Response.Header.SetConnectionClose()
			ctx.Response.Header.SetContentLength(len("Bad Request"))
			ctx.WriteString("Bad Request")
			return
		}

		filter.Gender = &gender
	}

	locationBytes := locationApiHandler.storage.GetLocationById(*filter.LocationId)

	if locationBytes == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.SetConnectionClose()
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		return
	}

	visitsIds := locationApiHandler.storage.GetVisitsByLocationId(*filter.LocationId)

	visitCollection := new(entities.VisitCollection)
	sumOfMarks := 0

	for _, visitId := range visitsIds {

		visitBytes := locationApiHandler.storage.GetVisitId(visitId)

		visit := new(entities.Visit)

		err := json.Unmarshal(visitBytes, visit)

		if err != nil {
			locationApiHandler.errLogger.Fatalln(err)
		}

		userBytes := locationApiHandler.storage.GetUserById(visit.User)

		user := new(entities.User)

		err = json.Unmarshal(userBytes, user)

		if err != nil {
			locationApiHandler.errLogger.Fatalln(err)
		}

		if !filter.CheckFromAge(user.BirthDate) ||
			!filter.CheckToAge(user.BirthDate) ||
			!filter.CheckToDate(visit.VisitedAt) ||
			!filter.CheckFromDate(visit.VisitedAt) ||
			!filter.CheckGender(user.Gender) {
			continue
		}

		sumOfMarks += visit.Mark
		visitCollection.Visits = append(visitCollection.Visits, visit)
	}

	locationAvgMark := &entities.LocationAvgMark{Avg: 0}

	if len(visitCollection.Visits) != 0 {
		locationAvgMark.Avg = math.Round(float64(sumOfMarks)/float64(len(visitCollection.Visits))*100000) / 100000
	}

	locationAvgMarkBytes, err := json.Marshal(locationAvgMark)

	if err != nil {
		locationApiHandler.errLogger.Fatalln(err)
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetConnectionClose()
	ctx.Response.Header.SetContentLength(len(locationAvgMarkBytes))
	ctx.Write(locationAvgMarkBytes)
}