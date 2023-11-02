package handlers

import (
	"encoding/json"
	"jobportalapi/internal/auth"
	middlewares "jobportalapi/internal/middlerwares"
	"jobportalapi/internal/models"
	"jobportalapi/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type handler struct {
	S services.Store
	A *auth.Auth
}

// create comapny table in database
func (h *handler) createCompany(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middlewares.TraceIdKey).(string)
	if !ok {
		// If the traceId isn't found in the request, log an error and return
		log.Error().Msg("traceId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}
	var nc models.NewCompany

	// Attempt to decode JSON from the request body into the NewUser variable
	err := json.NewDecoder(c.Request.Body).Decode(&nc)
	if err != nil {
		// If there is an error in decoding, log the error and return
		log.Error().Err(err).Str("Trace Id", traceId)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}

	//for validation
	validate := validator.New()
	err = validate.Struct(&nc)
	if err != nil {
		// If validation fails, log the error and return
		log.Error().Err(err).Str("Trace Id", traceId).Send()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "please provide Name and Location"})
		return
	}

	//store to database
	cmpny, err := h.S.StoreCompany(ctx, nc)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Msg("database creation is not happening")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "company creation in database failed"})
		return

	}
	// If everything goes right, respond with the created user
	c.JSON(http.StatusOK, cmpny)

}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//fetch a particular company from database using comapny id it will get from seeing the database

func (h *handler) getCompanyById(c *gin.Context) {
	///context to trace request
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middlewares.TraceIdKey).(string)
	if !ok {
		log.Error().Msg("traceId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}

	stringCmpnyId := c.Param("cmpny_id")
	cid, err := strconv.ParseUint(stringCmpnyId, 10, 64)
	if err != nil {

		log.Print("conversion string to int error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "error found at conversion.."})
		return

	}
	val, err := h.S.GetCompanyData(cid)
	if err != nil {
		if err != nil {
			log.Error().Err(err).Str("Trace Id", traceId).Msg("not able to hit the database")
			log.Print("company data not found in database %w", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "company record not found"})
			return
		}

	}
	c.JSON(http.StatusOK, val)
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////
// to get all the company details
func (h *handler) getCompanyByAll(c *gin.Context) {
	///context to trace request
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middlewares.TraceIdKey).(string)
	if !ok {
		log.Error().Msg("traceId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}

	val, err := h.S.GetCompanyAllData()
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Msg("not able to hit the database")
		log.Print("company table not present in database")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "Company table not there"})
		return
	}
	c.JSON(http.StatusOK, val)
}
