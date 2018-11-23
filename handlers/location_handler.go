package handlers

import (
	"bufio"
	"github.com/ArtyomNorin/hlcup2017/entities"
	"github.com/ArtyomNorin/hlcup2017/services"
	"github.com/asaskevich/govalidator"
	"github.com/json-iterator/go"
	"github.com/paulbellamy/ratecounter"
	"github.com/valyala/fasthttp"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type LocationApiHandler struct {
	storage            *services.Storage
	errLogger          *log.Logger
	infoLogger         *log.Logger
	timeDataGeneration time.Time
	rateCounter        *ratecounter.RateCounter
}

func NewLocationApiHandler(storage *services.Storage, errLogger *log.Logger, infoLogger *log.Logger, rateCounter *ratecounter.RateCounter, pathToOptions string) *LocationApiHandler {

	file, _ := os.Open(pathToOptions)

	fileScanner := bufio.NewScanner(file)

	fileScanner.Scan()

	timeDataGeneration, err := strconv.Atoi(fileScanner.Text())

	if err != nil {
		errLogger.Fatalln(err)
	}

	return &LocationApiHandler{storage: storage, errLogger: errLogger, infoLogger: infoLogger, timeDataGeneration: time.Unix(int64(timeDataGeneration), 0), rateCounter: rateCounter}
}

func (locationApiHandler *LocationApiHandler) GetById(ctx *fasthttp.RequestCtx) {

	locationIdString, ok := ctx.UserValue("location_id").(string)

	if !govalidator.IsNumeric(locationIdString) || !ok {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.Set("Connection", "keep-alive")
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		locationApiHandler.rateCounter.Incr(1)
		return
	}

	locationId, err := strconv.Atoi(locationIdString)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.Set("Connection", "keep-alive")
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		locationApiHandler.rateCounter.Incr(1)
		return
	}

	location := locationApiHandler.storage.GetLocationById(uint(locationId))

	if location == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.Set("Connection", "keep-alive")
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		locationApiHandler.rateCounter.Incr(1)
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.Response.Header.SetContentType("application/json")
		ctx.Response.Header.Set("Connection", "keep-alive")
		ctx.Response.Header.SetContentLength(len(location))
		ctx.Write(location)
		locationApiHandler.rateCounter.Incr(1)
	}
}

func (locationApiHandler *LocationApiHandler) GetAverageMark(ctx *fasthttp.RequestCtx) {

	locationIdString, ok := ctx.UserValue("location_id").(string)

	if !govalidator.IsNumeric(locationIdString) || !ok {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.Set("Connection", "keep-alive")
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		locationApiHandler.rateCounter.Incr(1)
		return
	}

	locationIdInt, err := strconv.Atoi(locationIdString)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.Set("Connection", "keep-alive")
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		locationApiHandler.rateCounter.Incr(1)
		return
	}

	filter := services.InitVisitFilter(locationApiHandler.timeDataGeneration)

	locationIdUint := uint(locationIdInt)

	filter.LocationId = &locationIdUint

	if ctx.QueryArgs().Has("fromDate") {

		fromDateString := string(ctx.QueryArgs().Peek("fromDate"))

		if !govalidator.IsNumeric(fromDateString) || fromDateString == "" {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.Response.Header.SetContentType("text/plain; charset=utf8")
			ctx.Response.Header.Set("Connection", "keep-alive")
			ctx.Response.Header.SetContentLength(len("Bad Request"))
			ctx.WriteString("Bad Request")
			locationApiHandler.rateCounter.Incr(1)
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
			ctx.Response.Header.Set("Connection", "keep-alive")
			ctx.Response.Header.SetContentLength(len("Bad Request"))
			ctx.WriteString("Bad Request")
			locationApiHandler.rateCounter.Incr(1)
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
			ctx.Response.Header.Set("Connection", "keep-alive")
			ctx.Response.Header.SetContentLength(len("Bad Request"))
			ctx.WriteString("Bad Request")
			locationApiHandler.rateCounter.Incr(1)
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
			ctx.Response.Header.Set("Connection", "keep-alive")
			ctx.Response.Header.SetContentLength(len("Bad Request"))
			ctx.WriteString("Bad Request")
			locationApiHandler.rateCounter.Incr(1)
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
			ctx.Response.Header.Set("Connection", "keep-alive")
			ctx.Response.Header.SetContentLength(len("Bad Request"))
			ctx.WriteString("Bad Request")
			locationApiHandler.rateCounter.Incr(1)
			return
		}

		filter.Gender = &gender
	}

	locationBytes := locationApiHandler.storage.GetLocationById(*filter.LocationId)

	if locationBytes == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.Set("Connection", "keep-alive")
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		locationApiHandler.rateCounter.Incr(1)
		return
	}

	visitsIds := locationApiHandler.storage.GetVisitsByLocationId(*filter.LocationId)

	visitCollection := new(entities.VisitCollection)
	sumOfMarks := 0

	for _, visitId := range visitsIds {

		visitBytes := locationApiHandler.storage.GetVisitById(visitId)

		visit := new(entities.Visit)

		err := jsoniter.ConfigFastest.Unmarshal(visitBytes, visit)

		if err != nil {
			locationApiHandler.errLogger.Fatalln(err)
		}

		userBytes := locationApiHandler.storage.GetUserById(*visit.User)

		user := new(entities.User)

		err = jsoniter.ConfigFastest.Unmarshal(userBytes, user)

		if err != nil {
			locationApiHandler.errLogger.Fatalln(err)
		}

		if !filter.CheckFromAge(*user.BirthDate) ||
			!filter.CheckToAge(*user.BirthDate) ||
			!filter.CheckToDate(*visit.VisitedAt) ||
			!filter.CheckFromDate(*visit.VisitedAt) ||
			!filter.CheckGender(*user.Gender) {
			continue
		}

		sumOfMarks += *visit.Mark
		visitCollection.Visits = append(visitCollection.Visits, visit)
	}

	locationAvgMark := &entities.LocationAvgMark{Avg: 0}

	if len(visitCollection.Visits) != 0 {
		locationAvgMark.Avg = math.Round(float64(sumOfMarks)/float64(len(visitCollection.Visits))*100000) / 100000
	}

	locationAvgMarkBytes, err := jsoniter.ConfigFastest.Marshal(locationAvgMark)

	if err != nil {
		locationApiHandler.errLogger.Fatalln(err)
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.Response.Header.SetContentLength(len(locationAvgMarkBytes))
	ctx.Write(locationAvgMarkBytes)
	locationApiHandler.rateCounter.Incr(1)
}

func (locationApiHandler *LocationApiHandler) CreateOrUpdate(ctx *fasthttp.RequestCtx) {

	if strings.Contains(ctx.URI().String(), "/locations/new") {
		locationApiHandler.Create(ctx)
		return
	}

	locationIdString, ok := ctx.UserValue("location_id").(string)

	if !govalidator.IsNumeric(locationIdString) || !ok {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.Set("Connection", "keep-alive")
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		locationApiHandler.rateCounter.Incr(1)
		return
	}

	locationId, err := strconv.Atoi(locationIdString)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.Set("Connection", "keep-alive")
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		locationApiHandler.rateCounter.Incr(1)
		return
	}

	locationBytes := locationApiHandler.storage.GetLocationById(uint(locationId))

	if locationBytes == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.Set("Connection", "keep-alive")
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		locationApiHandler.rateCounter.Incr(1)
		return
	}

	newLocationMap := make(map[string]interface{})

	err = jsoniter.ConfigFastest.Unmarshal(ctx.PostBody(), &newLocationMap)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.Set("Connection", "keep-alive")
		ctx.Response.Header.SetContentLength(len("Bad Request"))
		ctx.WriteString("Bad Request")
		locationApiHandler.rateCounter.Incr(1)
		return
	}

	location := new(entities.Location)

	err = jsoniter.ConfigFastest.Unmarshal(locationBytes, location)

	if err != nil {
		locationApiHandler.errLogger.Fatalln(err)
	}

	if value, ok := newLocationMap["place"]; ok {

		place, typeOk := value.(string)

		if value == nil || !typeOk {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.Response.Header.SetContentType("text/plain; charset=utf8")
			ctx.Response.Header.Set("Connection", "keep-alive")
			ctx.Response.Header.SetContentLength(len("Bad Request"))
			ctx.WriteString("Bad Request")
			locationApiHandler.rateCounter.Incr(1)
			return
		}

		location.Place = &place
	}

	if value, ok := newLocationMap["country"]; ok {

		country, typeOk := value.(string)

		if value == nil || !typeOk || len(country) > 50 {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.Response.Header.SetContentType("text/plain; charset=utf8")
			ctx.Response.Header.Set("Connection", "keep-alive")
			ctx.Response.Header.SetContentLength(len("Bad Request"))
			ctx.WriteString("Bad Request")
			locationApiHandler.rateCounter.Incr(1)
			return
		}

		location.Country = &country
	}

	if value, ok := newLocationMap["city"]; ok {

		city, typeOk := value.(string)

		if value == nil || !typeOk || len(city) > 50 {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.Response.Header.SetContentType("text/plain; charset=utf8")
			ctx.Response.Header.Set("Connection", "keep-alive")
			ctx.Response.Header.SetContentLength(len("Bad Request"))
			ctx.WriteString("Bad Request")
			locationApiHandler.rateCounter.Incr(1)
			return
		}

		location.City = &city
	}

	if value, ok := newLocationMap["distance"]; ok {

		distance, typeOk := value.(float64)

		if value == nil || !typeOk || distance <= 0 {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.Response.Header.SetContentType("text/plain; charset=utf8")
			ctx.Response.Header.Set("Connection", "keep-alive")
			ctx.Response.Header.SetContentLength(len("Bad Request"))
			ctx.WriteString("Bad Request")
			locationApiHandler.rateCounter.Incr(1)
			return
		}

		distanceAsUint := uint(distance)

		location.Distance = &distanceAsUint
	}

	locationApiHandler.storage.AddLocation(location)

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.Response.Header.SetContentLength(len([]byte("{}")))
	ctx.Write([]byte("{}"))
	locationApiHandler.rateCounter.Incr(1)
}

func (locationApiHandler *LocationApiHandler) Create(ctx *fasthttp.RequestCtx) {

	newLocationMap := make(map[string]interface{})

	err := jsoniter.ConfigFastest.Unmarshal(ctx.PostBody(), &newLocationMap)

	if err != nil {
		locationApiHandler.returnBadRequest(ctx)
		locationApiHandler.rateCounter.Incr(1)
		return
	}

	locationIdInterface, ok := newLocationMap["id"]

	if !ok {
		locationApiHandler.returnBadRequest(ctx)
		locationApiHandler.rateCounter.Incr(1)
		return
	}

	locationIdFloat, typeOk := locationIdInterface.(float64)

	if !typeOk {
		locationApiHandler.returnBadRequest(ctx)
		locationApiHandler.rateCounter.Incr(1)
		return
	}

	locationIdUint := uint(locationIdFloat)

	locationBytes := locationApiHandler.storage.GetLocationById(locationIdUint)

	if locationBytes != nil {
		locationApiHandler.returnBadRequest(ctx)
		locationApiHandler.rateCounter.Incr(1)
		return
	}

	placeInterface, ok := newLocationMap["place"]

	if !ok {
		locationApiHandler.returnBadRequest(ctx)
		locationApiHandler.rateCounter.Incr(1)
		return
	}

	countryInterface, ok := newLocationMap["country"]

	if !ok {
		locationApiHandler.returnBadRequest(ctx)
		locationApiHandler.rateCounter.Incr(1)
		return
	}

	cityInterface, ok := newLocationMap["city"]

	if !ok {
		locationApiHandler.returnBadRequest(ctx)
		locationApiHandler.rateCounter.Incr(1)
		return
	}

	distanceInterface, ok := newLocationMap["distance"]

	if !ok {
		locationApiHandler.returnBadRequest(ctx)
		locationApiHandler.rateCounter.Incr(1)
		return
	}

	location := new(entities.Location)

	location.Id = &locationIdUint

	place, typeOk := placeInterface.(string)

	if !typeOk {
		locationApiHandler.returnBadRequest(ctx)
		locationApiHandler.rateCounter.Incr(1)
		return
	}

	location.Place = &place

	country, typeOk := countryInterface.(string)

	if !typeOk || len(country) > 50 {
		locationApiHandler.returnBadRequest(ctx)
		locationApiHandler.rateCounter.Incr(1)
		return
	}

	location.Country = &country

	city, typeOk := cityInterface.(string)

	if !typeOk || len(city) > 50 {
		locationApiHandler.returnBadRequest(ctx)
		locationApiHandler.rateCounter.Incr(1)
		return
	}

	location.City = &city

	distanceFloat, typeOk := distanceInterface.(float64)

	if !typeOk || distanceFloat <= 0 {
		locationApiHandler.returnBadRequest(ctx)
		locationApiHandler.rateCounter.Incr(1)
		return
	}

	distance := uint(distanceFloat)

	location.Distance = &distance

	locationApiHandler.storage.AddLocation(location)

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.Response.Header.SetContentLength(len([]byte("{}")))
	ctx.Write([]byte("{}"))
	locationApiHandler.rateCounter.Incr(1)
}

func (locationApiHandler *LocationApiHandler) returnBadRequest(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusBadRequest)
	ctx.Response.Header.SetContentType("text/plain; charset=utf8")
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.Response.Header.SetContentLength(len("Bad Request"))
	ctx.WriteString("Bad Request")
}
