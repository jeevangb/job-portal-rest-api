package handlers

import (
	"fmt"
	"jobportalapi/internal/auth"
	middlewares "jobportalapi/internal/middlerwares"
	"jobportalapi/internal/services"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
)

func API(a *auth.Auth, c *services.Conn) *gin.Engine {

	// Create a new Gin engine; Gin is a HTTP web framework written in Go
	r := gin.New()
	m, err := middlewares.NewMid(a)
	store := services.NewStore(c)
	if err != nil {
		log.Panic().Msg("middlewares not set up")

	}
	h := handler{
		A: a,
		S: store,
	}

	r.Use(m.Log(), gin.Recovery())
	//users
	r.GET("/check", check)
	r.POST("/signup", h.Signup)
	r.POST("/userLogin", h.userLogin)

	//company
	r.POST("/createCompany", m.Authenticate(h.createCompany))
	r.GET("/getCompanyById/:cmpny_id", m.Authenticate(h.getCompanyById))
	r.GET("/getCompanyByAll", m.Authenticate(h.getCompanyByAll))

	//Job
	r.POST("/createJob/:cmpny_id", m.Authenticate(h.createJob))
	r.GET("/getJobById/:job_Id", m.Authenticate(h.viewJobById))
	r.GET("/getAllJob", m.Authenticate(h.viewAllJobs))
	r.GET("/viewJobByCid/:cid", m.Authenticate(h.viewJobByCid))

	return r

}

func check(c *gin.Context) {
	//handle panic using recovery function when happening in separate goroutine
	//go func() {
	//	panic("some kind of panic")
	//}()
	time.Sleep(time.Second * 3)
	select {
	case <-c.Request.Context().Done():
		fmt.Println("user not there")
		return
	default:
		c.JSON(http.StatusOK, gin.H{"msg": "statusOk"})

	}

}
