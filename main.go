package main

import (
	"context"
	"flag"
	"github.com/febytanzil/dockerapp/data/maps"
	"github.com/febytanzil/dockerapp/data/route"
	"github.com/febytanzil/dockerapp/data/token"
	"github.com/febytanzil/dockerapp/framework/config"
	"github.com/febytanzil/dockerapp/framework/database"
	"github.com/febytanzil/dockerapp/framework/middleware"
	route2 "github.com/febytanzil/dockerapp/services/route"
	"github.com/febytanzil/dockerapp/views/api"
	"github.com/febytanzil/dockerapp/views/async"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()
	log.Println("routeapp starting")

	ok := config.ReadConfig("/etc/dockerapp/", "main") || config.ReadConfig("files/config/", "main")
	if !ok {
		log.Fatal("failed to read config")
	}

	cfg := config.Get()

	cc := make(chan string, cfg.App.Limit)
	inject(cc, cfg)

	r := mux.NewRouter()
	r.HandleFunc("/route", middleware.LimitRate(api.SubmitRoute)).Methods(http.MethodPost)
	r.HandleFunc("/route/{token}", api.GetRoute).Methods(http.MethodGet)
	rt := http.TimeoutHandler(r, time.Second*5, "timeout, please try again some time")

	srv := &http.Server{
		Addr: "0.0.0.0:9000",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      rt, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	go func(ch chan string) {
		for {
			async.CalculateRoute(<-ch)
		}
	}(cc)
	/*go func(ch chan string) {
		for {
			log.Println("channel buffer len", len(ch))
			time.Sleep(time.Second * 3)
		}
	}(cc)*/

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

func inject(ch chan string, c *config.MainConfig) {
	database.InitDB(c.Database.DSN)

	maps.Init(nil)
	route.Init(nil)
	token.Init(nil)

	route2.Init(nil, ch)
}
