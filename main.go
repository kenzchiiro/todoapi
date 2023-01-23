package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pallat/todoapi/router"
	"github.com/pallat/todoapi/store"
	"github.com/pallat/todoapi/todo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/time/rate"
)

var (
	buildcommit = "dev"
	buildtime   = time.Now().String()
)

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Printf("please consider environment variables: %s\n", err)
	}

	// db, err := gorm.Open(mysql.Open(os.Getenv("DB_CONN")), &gorm.Config{})
	// if err != nil {
	// 	panic("failed to connect database")
	// }

	// if err := db.AutoMigrate(&todo.Todo{}); err != nil {
	// 	log.Println("auto migrate db", err)
	// }

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_CONN")))
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("myapp").Collection("todos")

	r := router.NewMyRouter()

	// gormStore := store.NewGormStore(db)
	mongoStore := store.NewMongoDBStore(collection)

	handler := todo.NewTodoHandler(mongoStore)
	r.POST("/todos", handler.NewTask)
	r.GET("/todos", handler.List)
	r.DELETE("/todos/:id", handler.Remove)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	s := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()
	fmt.Println("shutting down gracefully, press Ctrl+C again to force")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(timeoutCtx); err != nil {
		fmt.Println(err)
	}
}

var limiter = rate.NewLimiter(5, 5)

func limitedHandler(c *gin.Context) {
	if !limiter.Allow() {
		c.AbortWithStatus(http.StatusTooManyRequests)
		return
	}
	c.JSON(200, map[string]interface{}{
		"message": "pong",
	})
}
