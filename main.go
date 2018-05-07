package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/julienschmidt/httprouter"
)

type ServerConfiguration struct {
	Address         string
	Port            string
	StorageLocation string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getServerConfiguration() *ServerConfiguration {
	return &ServerConfiguration{
		Address:         getEnv("DIVVY_UP_SERVER_ADDRESS", ""),
		Port:            getEnv("DIVVY_UP_SERVER_PORT", "8080"),
		StorageLocation: getEnv("DIVVY_UP_SERVER_STORAGE_LOCATION", "/tmp/files"),
	}
}

var config = getServerConfiguration()

func Downloader(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fileName := ps.ByName("id")
	filePath := path.Join(config.StorageLocation, fileName)

	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Read %d bytes. Replying to client...", len(contents))

	w.Write(contents)
	// :fmt.Fprint(w)
}

func Uploader(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fileName := ps.ByName("id")
	filePath := path.Join(config.StorageLocation, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	n, err := io.Copy(file, r.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%d bytes have been recieved and writen to %s.\n", n, fileName)
	fmt.Fprintf(w, "\nFile %s written with success!\n", fileName)
}

func main() {

	config := getServerConfiguration()
	log.Printf("The files will be stored under '%s'", config.StorageLocation)

	router := httprouter.New()
	router.GET("/:id", Downloader)
	router.PUT("/:id", Uploader)

	log.Printf("Serving at address '%s' on port '%s'", config.Address, config.Port)
	log.Fatal(http.ListenAndServe(config.Address+":"+config.Port, router))
}
