package main

import (
	"log"
	"time"

	"github.com/FrankOHara43/go-grpc-microservice/account"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/tools/go/analysis/passes/defers"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `env:"DATABASE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("",&cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r account.Respository
	retry.ForeverSleep(2*time.Second, func(_ int)(err error){
		r, err =account.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return 
	})
	defer r.close()
	log.Println("Listening on port 8080...")
	s := account.NewService(r)
	log.Fatal(account.ListenGRPC(s, 8080))
}