package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"goelster/s3"
	"goelster/bse"
	"goelster/mcx"
	"goelster/nse"
	"goelster/ncdex"
)

var bucketName = "bhavcopy"
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/download", downloadCSV).Methods("GET")
	http.Handle("/", router)
	fmt.Println("server is running on :8080")
	http.ListenAndServe(":8080", nil)
	}

func upload_download_file_from_s3(w http.ResponseWriter, r *http.Request, bucketName, key string) {
	err := s3.UploadFileToS3(bucketName, key, "output.csv")
	if err != nil {
		http.Error(w, "Failed to upload file to S3", http.StatusInternalServerError)
		return
	}
	key = key +"output.csv"
	err = s3.DownloadFileHandler(w, r, bucketName, key)
	if err != nil{
		http.Error(w, "Failed to download the file from S3", http.StatusInternalServerError)
	}
}
func downloadCSV(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	segment := r.URL.Query().Get("segment")
	date := r.URL.Query().Get("date")

	if segment == "" || date == "" {
		http.Error(w, "Segment and date parameters are required", http.StatusBadRequest)
		return
	}
	
	switch segment {
	case "bse":
		key := fmt.Sprintf("bsebhavcopy/dt=%s/", date)
		implement_bse(w, r, date, key)
	case "mcx":
		key := fmt.Sprintf("mcxbhavcopy/dt=%s/", date)
		implement_mcx(w, r, date, key)
	case "nsefo":
		key := fmt.Sprintf("nsefobhavcopy/dt=%s/", date)
		implement_nsefo(w, r, date, key)
	case "ncdex":
		key := fmt.Sprintf("ncdexbhavcopy/dt=%s/", date)
		implement_ncdex(w, r, date, key)
	case "nsefo_":
		key := fmt.Sprintf("nsefobhavcopy/dt=%s/", date)
		implement_nsefo_without_client(w, r, date, key)
	}
}

func check_file_exist_then_download(w http.ResponseWriter, r *http.Request, key string) bool {
	fmt.Println(key)
	isExist, err := s3.CheckIfS3FileExists(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	if isExist {
		key = key+"output.csv"
		err := s3.DownloadFileHandler(w, r, bucketName, key)
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return true
	}
	return false
}

func implement_mcx(w http.ResponseWriter, r *http.Request, date, key string) {
	isExist := check_file_exist_then_download(w, r, key)
	if !isExist{
		err := mcx.Mcx(date)
		if err !=nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		upload_download_file_from_s3(w,r, bucketName, key)
	}
}

func implement_nsefo(w http.ResponseWriter, r *http.Request, date, key string) {
	isExist := check_file_exist_then_download(w, r, key)
	if !isExist{
		err := nse.NseFo(date)
		if err !=nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		upload_download_file_from_s3(w,r, bucketName, key)
	}
}
func implement_nsefo_without_client(w http.ResponseWriter, r *http.Request, date, key string){
	isExist := check_file_exist_then_download(w, r, key)
	if !isExist{
		err := nse.NseFo_without_client(date)
		if err !=nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		upload_download_file_from_s3(w,r, bucketName, key)
	}
}

func implement_ncdex(w http.ResponseWriter, r *http.Request, date, key string){
	isExist := check_file_exist_then_download(w, r, key)
	if !isExist{
		err := ncdex.Ncdex(date)
		if err !=nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err  = ncdex.UnzipAndSaveAsOutputCSV("bhavcopy.zip")
		if err !=nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		upload_download_file_from_s3(w,r, bucketName, key)
	}
}

func implement_bse(w http.ResponseWriter, r *http.Request, date, key string){
		isExist := check_file_exist_then_download(w, r, key)
		if !isExist{
			err := bse.Bse(date)
			if err !=nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err  = bse.ConvertZipToCsv("bse_bhavcopy.zip")
			if err !=nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			upload_download_file_from_s3(w,r, bucketName, key)
		}
}