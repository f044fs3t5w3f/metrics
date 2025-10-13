package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/f044fs3t5w3f/metrics/internal/handler"
	"github.com/f044fs3t5w3f/metrics/internal/repository"
)

func main() {
	fmt.Println(os.Args)
	parseFlags()
	storage := repository.NewMemStorage()
	r := handler.GetRouter(storage)
	fmt.Println(flagRunAddr)
	err := http.ListenAndServe(flagRunAddr, r)
	if err != nil {
		panic(err.Error())
	}
}
