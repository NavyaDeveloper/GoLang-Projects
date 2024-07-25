package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
)

func newRouter() *httprouter.Router {
	mux := httprouter.New()
	mux.GET("/youtube/channel/stats", getchannelstats())

	return mux
}

func getchannelstats() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Write([]byte("Welcome!"))

	}
}

func main() {
	srv := &http.Server{
		Addr:    ":10101",
		Handler: newRouter(),
	}

	idleconnclosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint

		log.Println("service interrupt received")
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("http server shutdown error:%v", err)
		}

		log.Println("Shutdown Complete")
		close(idleconnclosed)

	}()

	log.Println("Starting server on port 10101")
	if err := srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("fatal http server failed to start:%v", err)
		}
	}

	<-idleconnclosed
	log.Println("Service Stop")
}
