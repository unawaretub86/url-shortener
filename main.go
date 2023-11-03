package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	h "github.com/unawaretub86/url-shortener/api"
	mr "github.com/unawaretub86/url-shortener/repository/mongodb"
	rr "github.com/unawaretub86/url-shortener/repository/redis"
	"github.com/unawaretub86/url-shortener/shortener"
)

// This architecture works like this
// repo <-- service --> serializer --> http
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	repo := chooseRepo()
	service := shortener.NewRedirectService(repo)
	handler := h.NewHandler(service)

	r := gin.Default()

	r.GET("/:code", handler.Get)
	r.POST("/", handler.Post)

	errs := make(chan error, 2) // Create a channel for errors with a capacity of 2.

	// First goroutine: Start a web server on port 8080
	go func() {
		fmt.Println("Listening port :8080")
		errs <- http.ListenAndServe(httpPort(), r) // Attempt to start the web server and send any error to the "errs" channel.
	}()

	// Second goroutine: Capture the SIGINT signal (Ctrl+C) for controlled termination
	go func() {
		c := make(chan os.Signal, 1)     // Create a channel for signals (in this case, SIGINT).
		signal.Notify(c, syscall.SIGINT) // Register the SIGINT signal to the "c" channel.
		errs <- fmt.Errorf("%s", <-c)    // When the SIGINT signal is received, send an error to the "errs" channel.
	}()

	// Block until one of the previously created goroutines sends an error (either due to a server error or receiving SIGINT).
	fmt.Printf("Terminated %s", <-errs)

}

// This gets the port and set it if it is empty
func httpPort() string {
	port := "8080"

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	return fmt.Sprintf(":%s", port)
}

// This choose repo depending the database to use
func chooseRepo() shortener.RedirectRepository {
	switch os.Getenv("URL_DB") {
	case "redis":
		redisURL := os.Getenv("REDIS_URL")
		repo, err := rr.NewRedisRepository(redisURL)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "mongo":
		mongoURL := os.Getenv("MONGO_URL")
		mongodb := os.Getenv("MONGO_DB")
		mongoTimeOut, _ := strconv.Atoi(os.Getenv("MONGO_TIME_OUT"))
		repo, err := mr.NewMongoRepository(mongoURL, mongodb, mongoTimeOut)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	}
	return nil
}
