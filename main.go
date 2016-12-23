package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func init() {
	// Assign a user store
	InitUserStore()
	// Assign session store
	InitSessionStore()
	// Assign mongoDB
	InitMongoDB()
}

func main() {

	addr := ":3000"
	router := NewRouter()

	router.Handle("GET", "/", HandleHome)
	router.Handle("GET", "/register", HandleUserNew)
	router.Handle("POST", "/register", HandleUserCreate)
	router.Handle("GET", "/login", HandleSessionNew)
	router.Handle("POST", "/login", HandleSessionCreate)
	router.Handle("GET", "/image/:imageID", HandleImageShow)
	router.Handle("GET", "/user/:userID", HandleUserShow)

	router.ServeFiles("/assets/*filepath", http.Dir("assets/"))
	router.ServeFiles("/im/*filepath", http.Dir("data/images/"))

	secureRouter := NewRouter()
	secureRouter.Handle("GET", "/sign-out", HandleSessionDestroy)
	secureRouter.Handle("GET", "/account", HandleUserEdit)
	secureRouter.Handle("POST", "/account", HandleUserUpdate)
	secureRouter.Handle("GET", "/images/new", HandleImageNew)
	secureRouter.Handle("POST", "/images/new", HandleImageCreate)

	notFoundRouter := NewRouterCustom()

	middleware := Middleware{}
	middleware.Add(router)
	middleware.Add(http.HandlerFunc(RequireLogin))
	middleware.Add(secureRouter)
	middleware.Add(notFoundRouter)

	log.Fatal(http.ListenAndServe(addr, middleware))

	defer CloseMongoDBSession(mongoSession)
}

type NotFound struct{}

func (n *NotFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//RenderTemplate(w, r, "index/404", nil)
}
func NewRouter() *httprouter.Router {
	router := httprouter.New()
	notFound := new(NotFound)
	router.NotFound = notFound
	return router
}

type NotFoundCustom struct{}

func (n *NotFoundCustom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "index/404", nil)
}
func NewRouterCustom() *httprouter.Router {
	router := httprouter.New()
	notFound := new(NotFoundCustom)
	router.NotFound = notFound
	return router
}
