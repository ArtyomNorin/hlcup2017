package handlers

import (
	"hlcup/services"
	"log"
	"github.com/valyala/fasthttp"
	"github.com/asaskevich/govalidator"
	"strconv"
	"github.com/json-iterator/go"
	"hlcup/entities"
	"strings"
	"github.com/paulbellamy/ratecounter"
)

type VisitApiHandler struct {
	storage     *services.Storage
	errLogger   *log.Logger
	infoLogger  *log.Logger
	rateCounter *ratecounter.RateCounter
}

var oldLocationId, oldUserId uint

func NewVisitApiHandler(storage *services.Storage, errLogger *log.Logger, infoLogger *log.Logger, rateCounter *ratecounter.RateCounter) *VisitApiHandler {

	return &VisitApiHandler{storage: storage, errLogger: errLogger, infoLogger: infoLogger, rateCounter: rateCounter}
}

func (visitApiHandler *VisitApiHandler) GetById(ctx *fasthttp.RequestCtx) {

	visitIdString, ok := ctx.UserValue("visit_id").(string)

	if !govalidator.IsNumeric(visitIdString) || !ok {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.Set("Connection", "keep-alive")
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		visitApiHandler.rateCounter.Incr(1)
		return
	}

	visitId, err := strconv.Atoi(visitIdString)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.Set("Connection", "keep-alive")
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		visitApiHandler.rateCounter.Incr(1)
		return
	}

	visit := visitApiHandler.storage.GetVisitById(uint(visitId))

	if visit == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.Set("Connection", "keep-alive")
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		visitApiHandler.rateCounter.Incr(1)
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.Response.Header.SetContentType("application/json")
		ctx.Response.Header.Set("Connection", "keep-alive")
		ctx.Response.Header.SetContentLength(len(visit))
		ctx.Write(visit)
		visitApiHandler.rateCounter.Incr(1)
	}
}

func (visitApiHandler *VisitApiHandler) CreateOrUpdate(ctx *fasthttp.RequestCtx) {

	if strings.Contains(ctx.URI().String(), "/visits/new") {
		visitApiHandler.Create(ctx)
		return
	}

	visitIdString, ok := ctx.UserValue("visit_id").(string)

	if !govalidator.IsNumeric(visitIdString) || !ok {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.Set("Connection", "keep-alive")
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		visitApiHandler.rateCounter.Incr(1)
		return
	}

	visitId, err := strconv.Atoi(visitIdString)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.Set("Connection", "keep-alive")
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		visitApiHandler.rateCounter.Incr(1)
		return
	}

	visitBytes := visitApiHandler.storage.GetVisitById(uint(visitId))

	if visitBytes == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.Set("Connection", "keep-alive")
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		visitApiHandler.rateCounter.Incr(1)
		return
	}

	newVisitMap := make(map[string]interface{})

	err = jsoniter.ConfigFastest.Unmarshal(ctx.PostBody(), &newVisitMap)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.Set("Connection", "keep-alive")
		ctx.Response.Header.SetContentLength(len("Bad Request"))
		ctx.WriteString("Bad Request")
		visitApiHandler.rateCounter.Incr(1)
		return
	}

	visit := new(entities.Visit)

	err = jsoniter.ConfigFastest.Unmarshal(visitBytes, visit)

	if err != nil {
		visitApiHandler.errLogger.Fatalln(err)
	}

	if value, ok := newVisitMap["location"]; ok {

		locationId, typeOk := value.(float64)

		if value == nil || !typeOk || locationId <= 0 {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.Response.Header.SetContentType("text/plain; charset=utf8")
			ctx.Response.Header.Set("Connection", "keep-alive")
			ctx.Response.Header.SetContentLength(len("Bad Request"))
			ctx.WriteString("Bad Request")
			visitApiHandler.rateCounter.Incr(1)
			return
		}

		locationIdAsUint := uint(locationId)

		oldLocationId = *visit.Location
		visit.Location = &locationIdAsUint
	}

	if value, ok := newVisitMap["user"]; ok {

		userId, typeOk := value.(float64)

		if value == nil || !typeOk || userId <= 0 {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.Response.Header.SetContentType("text/plain; charset=utf8")
			ctx.Response.Header.Set("Connection", "keep-alive")
			ctx.Response.Header.SetContentLength(len("Bad Request"))
			ctx.WriteString("Bad Request")
			visitApiHandler.rateCounter.Incr(1)
			return
		}

		userIdAsUint := uint(userId)

		oldUserId = *visit.User
		visit.User = &userIdAsUint
	}

	if value, ok := newVisitMap["visited_at"]; ok {

		visitedAt, typeOk := value.(float64)

		if value == nil || !typeOk {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.Response.Header.SetContentType("text/plain; charset=utf8")
			ctx.Response.Header.Set("Connection", "keep-alive")
			ctx.Response.Header.SetContentLength(len("Bad Request"))
			ctx.WriteString("Bad Request")
			visitApiHandler.rateCounter.Incr(1)
			return
		}

		visitedAtAsInt := int(visitedAt)

		visit.VisitedAt = &visitedAtAsInt
	}

	if value, ok := newVisitMap["mark"]; ok {

		mark, typeOk := value.(float64)

		if value == nil || !typeOk || (mark != 0 && mark != 1 && mark != 2 && mark != 3 && mark != 4 && mark != 5) {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.Response.Header.SetContentType("text/plain; charset=utf8")
			ctx.Response.Header.Set("Connection", "keep-alive")
			ctx.Response.Header.SetContentLength(len("Bad Request"))
			ctx.WriteString("Bad Request")
			visitApiHandler.rateCounter.Incr(1)
			return
		}

		markAsInt := int(mark)

		visit.Mark = &markAsInt
	}

	if _, ok := newVisitMap["location"]; ok {
		visitApiHandler.storage.DeleteVisitFromLocation(oldLocationId, *visit.Id)
		visitApiHandler.storage.AddVisitByLocationId(visit)
	}

	if _, ok := newVisitMap["user"]; ok {
		visitApiHandler.storage.DeleteVisitFromUser(oldUserId, *visit.Id)
		visitApiHandler.storage.AddVisitByUserId(visit)
	}

	visitApiHandler.storage.AddVisit(visit)

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.Response.Header.SetContentLength(len([]byte("{}")))
	ctx.Write([]byte("{}"))
	visitApiHandler.rateCounter.Incr(1)
}

func (visitApiHandler *VisitApiHandler) Create(ctx *fasthttp.RequestCtx) {

	newVisitMap := make(map[string]interface{})

	err := jsoniter.ConfigFastest.Unmarshal(ctx.PostBody(), &newVisitMap)

	if err != nil {
		visitApiHandler.returnBadRequest(ctx)
		visitApiHandler.rateCounter.Incr(1)
		return
	}

	visitIdInterface, ok := newVisitMap["id"]

	if !ok {
		visitApiHandler.returnBadRequest(ctx)
		visitApiHandler.rateCounter.Incr(1)
		return
	}

	visitIdFloat, typeOk := visitIdInterface.(float64)

	if !typeOk {
		visitApiHandler.returnBadRequest(ctx)
		visitApiHandler.rateCounter.Incr(1)
		return
	}

	visitIdUint := uint(visitIdFloat)

	visitBytes := visitApiHandler.storage.GetVisitById(visitIdUint)

	if visitBytes != nil {
		visitApiHandler.returnBadRequest(ctx)
		visitApiHandler.rateCounter.Incr(1)
		return
	}

	locationIdInterface, ok := newVisitMap["location"]

	if !ok {
		visitApiHandler.returnBadRequest(ctx)
		visitApiHandler.rateCounter.Incr(1)
		return
	}

	userIdInterface, ok := newVisitMap["user"]

	if !ok {
		visitApiHandler.returnBadRequest(ctx)
		visitApiHandler.rateCounter.Incr(1)
		return
	}

	VisitedAtInterface, ok := newVisitMap["visited_at"]

	if !ok {
		visitApiHandler.returnBadRequest(ctx)
		visitApiHandler.rateCounter.Incr(1)
		return
	}

	markInterface, ok := newVisitMap["mark"]

	if !ok {
		visitApiHandler.returnBadRequest(ctx)
		visitApiHandler.rateCounter.Incr(1)
		return
	}

	visit := new(entities.Visit)

	visit.Id = &visitIdUint

	locationIdFloat, typeOk := locationIdInterface.(float64)

	if !typeOk || locationIdFloat <= 0 {
		visitApiHandler.returnBadRequest(ctx)
		visitApiHandler.rateCounter.Incr(1)
		return
	}

	locationIdUint := uint(locationIdFloat)

	location := visitApiHandler.storage.GetLocationById(locationIdUint)

	if location == nil {
		visitApiHandler.returnBadRequest(ctx)
		visitApiHandler.rateCounter.Incr(1)
		return
	}

	visit.Location = &locationIdUint

	userIdFloat, typeOk := userIdInterface.(float64)

	if !typeOk || userIdFloat <= 0 {
		visitApiHandler.returnBadRequest(ctx)
		visitApiHandler.rateCounter.Incr(1)
		return
	}

	userIdUint := uint(userIdFloat)

	user := visitApiHandler.storage.GetUserById(userIdUint)

	if user == nil {
		visitApiHandler.returnBadRequest(ctx)
		visitApiHandler.rateCounter.Incr(1)
		return
	}

	visit.User = &userIdUint

	visitedAtFloat, typeOk := VisitedAtInterface.(float64)

	if !typeOk {
		visitApiHandler.returnBadRequest(ctx)
		visitApiHandler.rateCounter.Incr(1)
		return
	}

	visitedAtInt := int(visitedAtFloat)

	visit.VisitedAt = &visitedAtInt

	markFloat, typeOk := markInterface.(float64)

	if !typeOk || (markFloat != 0 && markFloat != 1 && markFloat != 2 && markFloat != 3 && markFloat != 4 && markFloat != 5) {
		visitApiHandler.returnBadRequest(ctx)
		visitApiHandler.rateCounter.Incr(1)
		return
	}

	markInt := int(markFloat)

	visit.Mark = &markInt

	visitApiHandler.storage.AddVisit(visit)
	visitApiHandler.storage.AddVisitByUserId(visit)
	visitApiHandler.storage.AddVisitByLocationId(visit)

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.Response.Header.SetContentLength(len([]byte("{}")))
	ctx.Write([]byte("{}"))
	visitApiHandler.rateCounter.Incr(1)
}

func (visitApiHandler *VisitApiHandler) returnBadRequest(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusBadRequest)
	ctx.Response.Header.SetContentType("text/plain; charset=utf8")
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.Response.Header.SetContentLength(len("Bad Request"))
	ctx.WriteString("Bad Request")
}
