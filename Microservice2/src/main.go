package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type UserTrans struct {
	gorm.Model
	KnownUrl   string
	UnknownUrl string
	Request    string
	Response   string
	RequestAt  time.Time
	ResponseAt time.Time
}

var db *gorm.DB

func UploadBlob(blobname string, containerURL azblob.ContainerURL, file1 multipart.File) string {
	blobURL := containerURL.NewBlockBlobURL(blobname)
	ctx := context.Background()

	fmt.Printf("Uploading the file with blob name: %s\n", blobname)

	o := azblob.UploadToBlockBlobOptions{
		BlobHTTPHeaders: azblob.BlobHTTPHeaders{
			ContentType: "image/jpg",
		},
	}

	buf := bytes.NewBuffer(nil)
	temp, err := io.Copy(buf, file1)

	if err != nil {
		log.Fatal(err.Error(), temp)
	}
	_, err = azblob.UploadBufferToBlockBlob(ctx, buf.Bytes(), blobURL, o)

	if err != nil {
		log.Fatal(err.Error())
	}

	return blobURL.String()
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Reached Upload endpoint...")

	r.ParseMultipartForm(10 << 20)

	file1, _, err := r.FormFile("known")

	if err != nil {
		fmt.Println("Error Retrieving the known File")
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer file1.Close()

	file2, _, err := r.FormFile("unknown")

	if err != nil {
		fmt.Println("Error Retrieving the unknown File")
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer file2.Close()

	// upload the images to Cloud

	accountName, accountKey, containerName := os.Getenv("AZURE_STORAGE_ACCOUNT"), os.Getenv("AZURE_STORAGE_ACCESS_KEY"), os.Getenv("AZURE_STORAGE_CONTAINER")
	if len(accountName) == 0 || len(accountKey) == 0 || len(containerName) == 0 {
		log.Fatal("Either the AZURE_STORAGE_ACCOUNT or AZURE_STORAGE_ACCESS_KEY environment or AZURE_STORAGE_CONTAINER variable is not set")
	}

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		fmt.Println(accountName)
		fmt.Println(accountKey)
		log.Fatal("Invalid credentials with error: " + err.Error())
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

	containerURL := azblob.NewContainerURL(*URL, p)

	blobUrl_known := UploadBlob("known_image.jpg", containerURL, file1)
	blobUrl_unknown := UploadBlob("unknown_image.jpg", containerURL, file2)

	// end upload

	// request to microservice1
	requestAt := time.Now()
	fmt.Println("Calling Microservice1...")

	body := fmt.Sprintf("{\"known_image\":\"%s\",\"unknown_image\":\"%s\"}", blobUrl_known, blobUrl_unknown)
	requestBody := strings.NewReader(body)

	res, err := http.Post(
		"http://py-docker:8000/compare-faces/",
		"application/json; charset=UTF-8",
		requestBody,
	)

	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
	}

	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.StatusCode)
	w.Write([]byte(data))

	// Write to DB
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.Create(&UserTrans{KnownUrl: blobUrl_known, UnknownUrl: blobUrl_unknown, Request: body, Response: string(data), ResponseAt: time.Now(), RequestAt: requestAt})

	fmt.Printf("From Microservice1 - %s\n", []byte(data))
	fmt.Println("Leaving Upload endpoint...")
}

func GetUserTrans(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var usertrans []UserTrans
	db.Find(&usertrans)
	json.NewEncoder(w).Encode(usertrans)
}

func setupRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/upload", uploadFile).Methods("POST")
	r.HandleFunc("/user-trans", GetUserTrans).Methods("GET")
	http.ListenAndServe(":7070", r)
	fmt.Println("Listening at port 7070...")
}

func main() {
	fmt.Println("Starting the Upload server...")

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.AutoMigrate(&UserTrans{})

	setupRoutes()
	fmt.Println("Stopping the Upload server...")
}
