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
	"strings"
	"github.com/json-iterator/go"
	"hlcup/entities"
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
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.SetConnectionClose()
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		return
	}

	if userId <= 0 {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.SetConnectionClose()
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		return
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
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.SetConnectionClose()
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		return
	}

	if userId <= 0 {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.SetConnectionClose()
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		return
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

func (userApiHandler *UserApiHandler) CreateOrUpdate(ctx *fasthttp.RequestCtx) {

	if strings.Contains(ctx.URI().String(), "/users/new") {
		userApiHandler.Create(ctx)
		return
	}

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
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.SetConnectionClose()
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		return
	}

	userBytes := userApiHandler.storage.GetUserById(uint(userId))

	if userBytes == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.SetConnectionClose()
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		return
	}

	newUserMap := make(map[string]interface{})

	err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(ctx.PostBody(), &newUserMap)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.SetConnectionClose()
		ctx.Response.Header.SetContentLength(len("Bad Request"))
		ctx.WriteString("Bad Request")
		return
	}

	user := new(entities.User)

	err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(userBytes, user)

	if err != nil {
		userApiHandler.errLogger.Fatalln(err)
	}

	if value, ok := newUserMap["email"]; ok {

		email, typeOk := value.(string)

		if value == nil || !typeOk || len(email) > 100 {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.Response.Header.SetContentType("text/plain; charset=utf8")
			ctx.Response.Header.SetConnectionClose()
			ctx.Response.Header.SetContentLength(len("Bad Request"))
			ctx.WriteString("Bad Request")
			return
		}

		user.Email = &email
	}

	if value, ok := newUserMap["first_name"]; ok {

		firstName, typeOk := value.(string)

		if value == nil || !typeOk || len(firstName) > 50 {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.Response.Header.SetContentType("text/plain; charset=utf8")
			ctx.Response.Header.SetConnectionClose()
			ctx.Response.Header.SetContentLength(len("Bad Request"))
			ctx.WriteString("Bad Request")
			return
		}

		user.FirstName = &firstName
	}

	if value, ok := newUserMap["last_name"]; ok {

		lastName, typeOk := value.(string)

		if value == nil || !typeOk || len(lastName) > 50 {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.Response.Header.SetContentType("text/plain; charset=utf8")
			ctx.Response.Header.SetConnectionClose()
			ctx.Response.Header.SetContentLength(len("Bad Request"))
			ctx.WriteString("Bad Request")
			return
		}

		user.LastName = &lastName
	}

	if value, ok := newUserMap["gender"]; ok {

		gender, typeOk := value.(string)

		if value == nil || !typeOk || (gender != "m" && gender != "f") {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.Response.Header.SetContentType("text/plain; charset=utf8")
			ctx.Response.Header.SetConnectionClose()
			ctx.Response.Header.SetContentLength(len("Bad Request"))
			ctx.WriteString("Bad Request")
			return
		}

		user.Gender = &gender
	}

	if value, ok := newUserMap["birth_date"]; ok {

		birthDateFloat, typeOk := value.(float64)

		if value == nil || !typeOk {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.Response.Header.SetContentType("text/plain; charset=utf8")
			ctx.Response.Header.SetConnectionClose()
			ctx.Response.Header.SetContentLength(len("Bad Request"))
			ctx.WriteString("Bad Request")
			return
		}

		birthDate := int(birthDateFloat)

		user.BirthDate = &birthDate
	}

	userApiHandler.storage.AddUser(user)

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetConnectionClose()
	ctx.Response.Header.SetContentLength(len([]byte("{}")))
	ctx.Write([]byte("{}"))
}

func (userApiHandler *UserApiHandler) Create(ctx *fasthttp.RequestCtx) {

	newUserMap := make(map[string]interface{})

	err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(ctx.PostBody(), &newUserMap)

	if err != nil {
		userApiHandler.returnBadRequest(ctx)
		return
	}

	userIdInterface, ok := newUserMap["id"]

	if !ok {
		userApiHandler.returnBadRequest(ctx)
		return
	}

	userIdFloat, typeOk := userIdInterface.(float64)

	if !typeOk {
		userApiHandler.returnBadRequest(ctx)
		return
	}

	userIdUint := uint(userIdFloat)

	userBytes := userApiHandler.storage.GetLocationById(userIdUint)

	if userBytes != nil {
		userApiHandler.returnBadRequest(ctx)
		return
	}

	emailInterface, ok := newUserMap["email"]

	if !ok {
		userApiHandler.returnBadRequest(ctx)
		return
	}

	email, typeOk := emailInterface.(string)

	if !typeOk || len(email) > 100 {
		userApiHandler.returnBadRequest(ctx)
		return
	}

	isEmailExist := userApiHandler.storage.IsEmailExist(email)

	if isEmailExist {
		userApiHandler.returnBadRequest(ctx)
		return
	}

	firstNameInterface, ok := newUserMap["first_name"]

	if !ok {
		userApiHandler.returnBadRequest(ctx)
		return
	}

	lastNameInterface, ok := newUserMap["last_name"]

	if !ok {
		userApiHandler.returnBadRequest(ctx)
		return
	}

	genderInterface, ok := newUserMap["gender"]

	if !ok {
		userApiHandler.returnBadRequest(ctx)
		return
	}

	birthDateInterface, ok := newUserMap["birth_date"]

	if !ok {
		userApiHandler.returnBadRequest(ctx)
		return
	}

	user := new(entities.User)

	user.Id = &userIdUint

	user.Email = &email

	firstName, typeOk := firstNameInterface.(string)

	if !typeOk || len(firstName) > 50 {
		userApiHandler.returnBadRequest(ctx)
		return
	}

	user.FirstName = &firstName

	lastName, typeOk := lastNameInterface.(string)

	if !typeOk || len(lastName) > 50 {
		userApiHandler.returnBadRequest(ctx)
		return
	}

	user.LastName = &lastName

	gender, typeOk := genderInterface.(string)

	if !typeOk || (gender != "m" && gender != "f") {
		userApiHandler.returnBadRequest(ctx)
		return
	}

	user.Gender = &gender

	birthDateFloat, typeOk := birthDateInterface.(float64)

	if !typeOk {
		userApiHandler.returnBadRequest(ctx)
		return
	}

	birthDate := int(birthDateFloat)

	user.BirthDate = &birthDate

	userApiHandler.storage.AddUser(user)

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetConnectionClose()
	ctx.Response.Header.SetContentLength(len([]byte("{}")))
	ctx.Write([]byte("{}"))
}

func (userApiHandler *UserApiHandler) returnBadRequest(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusBadRequest)
	ctx.Response.Header.SetContentType("text/plain; charset=utf8")
	ctx.Response.Header.SetConnectionClose()
	ctx.Response.Header.SetContentLength(len("Bad Request"))
	ctx.WriteString("Bad Request")
}
