package handlers

import (
	"hlcup/services"
	"log"
	"github.com/valyala/fasthttp"
	"github.com/asaskevich/govalidator"
	"strconv"
)

type VisitApiHandler struct {
	storage    *services.Storage
	errLogger  *log.Logger
	infoLogger *log.Logger
}

func NewVisitApiHandler(storage *services.Storage, errLogger *log.Logger, infoLogger *log.Logger) *VisitApiHandler {

	return &VisitApiHandler{storage: storage, errLogger: errLogger, infoLogger: infoLogger}
}

func (visitApiHandler *VisitApiHandler) GetById(ctx *fasthttp.RequestCtx) {

	visitIdString, ok := ctx.UserValue("visit_id").(string)

	if !govalidator.IsNumeric(visitIdString) || !ok {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.SetConnectionClose()
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
		return
	}

	visitId, err := strconv.Atoi(visitIdString)

	if err != nil {
		visitApiHandler.errLogger.Fatalln(err)
	}

	visit := visitApiHandler.storage.GetVisitId(uint(visitId))

	if visit == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.Header.SetContentType("text/plain; charset=utf8")
		ctx.Response.Header.SetConnectionClose()
		ctx.Response.Header.SetContentLength(len("Not Found"))
		ctx.WriteString("Not Found")
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.Response.Header.SetContentType("application/json")
		ctx.Response.Header.SetConnectionClose()
		ctx.Response.Header.SetContentLength(len(visit))
		ctx.Write(visit)
	}
}
