package main

import (
	"data/absensi"
	"data/config"
	"data/config/api"
	"data/migrations"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	config.ConnectDB()
	config.InitMidtrans()

	router := mux.NewRouter()

	// --- Jalankan migrasi jika ada argumen 'migrate' ---
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		config.ConnectDB()
		migrations.Migrate()
		return
	}

	// //setting payment gateway use midtrans
	config.InitMidtrans()

	// --- Setup Gin routes ---
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	//file aplikasi router gin
	api.ApiAllRoleWithToken(r)
	router.PathPrefix("/api/murid/").Handler(r)
	router.PathPrefix("/api/auth/").Handler(r)
	router.PathPrefix("/api/pembimbing").Handler(r)
	router.PathPrefix("/api/jadwal").Handler(r)
	router.PathPrefix("/api/payment").Handler(r)
	router.PathPrefix("/api/evaluasi").Handler(r)

	//gorilla/mux
	router.HandleFunc("/api/v1/absensi", absensi.GetRekapAbsensi).Methods("GET")
	router.HandleFunc("/api/v1/createAbsensi/{id_jadwal}", absensi.CreateAbsensi).Methods("POST")

	// Serve static files
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})
	router.PathPrefix("/templates/").Handler(http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))

	// === TAMBAHKAN INI: CORS + HEADER JSON ===
	corsObj := handlers.CORS(
		handlers.AllowedOrigins([]string{"*", "http://localhost:3000", "https://bimbelsomagede.loca.lt", "https://sailing-vampire-pleased-statewide.trycloudflare.com", "http://192.168.1.28:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Wrapper agar semua response punya header JSON
	loggedRouter := handlers.LoggingHandler(log.Writer(), corsObj(router))

	// Server dengan timeout
	srv := &http.Server{
		Handler:      loggedRouter,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Server berjalan di http://localhost:8080")
	log.Fatal(srv.ListenAndServe())
}
