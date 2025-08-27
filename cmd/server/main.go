package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jelisavac-l/GBattleships/internal/routes"
)

// it did not work on my machine ???
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	fmt.Println("Battleships server started!")
	routes.RegisterServerRoutes()

	if err := http.ListenAndServe(":8080", corsMiddleware(http.DefaultServeMux)); err != nil {
		log.Fatal(err)
	}

}
