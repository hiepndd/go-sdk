package main

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

// Any is a type alias for interface{}
type Any = interface{}

// APIController implements a simple controller.
type APIController struct {
	db     map[string]Any
	dbLock sync.Mutex
}

// Register adds routes for the controller to the app.
func (ac *APIController) Register(app *web.App) {
	app.GET("/api", ac.all, ac.randomFailure)
	app.GET("/api/:key", ac.get, ac.randomFailure)
	app.POST("/api/:key", ac.post, ac.randomFailure)
	app.PUT("/api/:key", ac.put, ac.randomFailure)
	app.DELETE("/api/:key", ac.delete, ac.randomFailure)
}

func (ac *APIController) randomFailure(action web.Action) web.Action {
	return func(r *web.Ctx) web.Result {
		if rand.Int()%2 == 0 {
			return r.JSON().InternalError(fmt.Errorf("random error"))
		}
		return action(r)
	}
}

func (ac *APIController) ensureDB() {
	if ac.db == nil {
		ac.db = map[string]Any{}
	}
}

func (ac *APIController) all(r *web.Ctx) web.Result {
	ac.dbLock.Lock()
	defer ac.dbLock.Unlock()
	ac.ensureDB()

	return r.JSON().Result(ac.db)
}

func (ac *APIController) get(r *web.Ctx) web.Result {
	ac.dbLock.Lock()
	defer ac.dbLock.Unlock()
	ac.ensureDB()

	value, hasValue := ac.db[web.StringValue(r.Param("key"))]
	if !hasValue {
		return r.JSON().NotFound()
	}
	return r.JSON().Result(value)
}

func (ac *APIController) post(r *web.Ctx) web.Result {
	ac.dbLock.Lock()
	defer ac.dbLock.Unlock()
	ac.ensureDB()

	body, err := r.PostBody()
	if err != nil {
		return r.JSON().InternalError(err)
	}
	ac.db[web.StringValue(r.Param("key"))] = string(body)
	return r.JSON().OK()
}

func (ac *APIController) put(r *web.Ctx) web.Result {
	ac.dbLock.Lock()
	defer ac.dbLock.Unlock()
	ac.ensureDB()

	_, hasValue := ac.db[web.StringValue(r.Param("key"))]
	if !hasValue {
		return r.JSON().NotFound()
	}

	body, err := r.PostBody()
	if err != nil {
		return r.JSON().InternalError(err)
	}
	ac.db[web.StringValue(r.Param("key"))] = string(body)

	return r.JSON().OK()
}

func (ac *APIController) delete(r *web.Ctx) web.Result {
	ac.dbLock.Lock()
	defer ac.dbLock.Unlock()
	ac.ensureDB()

	key := web.StringValue(r.Param("key"))
	_, hasValue := ac.db[key]
	if !hasValue {
		return r.JSON().NotFound()
	}
	delete(ac.db, key)
	return r.JSON().OK()
}

func main() {
	log := logger.NewFromEnv()
	app := web.New().WithLogger(log).WithControllers(new(APIController))
	log.SyncFatalExit(app.Start())
}
