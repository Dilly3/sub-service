package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/dilly3/sub-service/cmd/web"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading env")
	}

}

var sessionManager *scs.SessionManager
var wg = &sync.WaitGroup{}
var infolog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var errorlog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
var PORT = os.Getenv("PORT")

func main() {
	fmt.Println("\n Welcome to sub-service")
	db := loadDb()
	err := db.Ping()
	if err != nil {
		fmt.Println(err)
	}

	app := web.Config{
		DB:       db,
		Session:  sessionManager,
		Wait:     wg,
		ErrorLog: errorlog,
		InfoLog:  infolog,
	}
	app.Serve(PORT)
}

func loadDb() *sql.DB {

	dsn := os.Getenv("DSN")
	var db *sql.DB
	var err error
	count := 0
	for {
		if count == 5 {
			fmt.Println("cant connect to db", err)
			break
		}

		db, err = sql.Open("postgres", dsn)
		if err != nil {
			fmt.Println(err)
		}
		if db != nil {
			break
		}
		count++
	}
	return db
}

func initRedis() *redis.Pool {
	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", os.Getenv("REDIS"))
		},
	}
	return redisPool
}

func InitSession() *scs.SessionManager {
	sessionManager = scs.New()
	sessionManager.Store = redisstore.New(initRedis())
	sessionManager.Lifetime = time.Hour * 24
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = true

	return sessionManager

}
