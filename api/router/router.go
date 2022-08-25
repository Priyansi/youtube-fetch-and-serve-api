package router

import (
	"github.com/julienschmidt/httprouter"
	"github.com/priyansi/fampay-backend-assignment/api/handlers/getvideos"
)

func Get() *httprouter.Router {
	mux := httprouter.New()
	mux.GET("/get_videos", getvideos.Do())
	return mux
}
