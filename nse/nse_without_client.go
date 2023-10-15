
package nse

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "strconv"
	"time"
    "strings"
)

func NseFo_without_client(date string) error{
	currentDate := time.Now()
	givenDate, err := time.Parse("02-01-2006", date)
    if err != nil {
        return fmt.Errorf("Error parsing given date: %s", err)
    }
	todayDate := currentDate.Format("2006-01-02")

    // Parse the current date string back into a time.Time value
    currDate, err := time.Parse("2006-01-02", todayDate)
    if givenDate.After(currDate) {
        return fmt.Errorf("NseFo file is not generated for %s", date)
    } 

    // Extract the year, month, and day
    year :=  givenDate.Format("2006")
    month := givenDate.Format("Jan")
    day := givenDate.Format("02")

	month = strings.ToUpper(month)

	url := fmt.Sprintf("https://nsearchives.nseindia.com/content/historical/DERIVATIVES/%s/%s/fo%s%s%sbhav.csv.zip", year, month, day, month, year)
    fmt.Println(url)
	// url := "https://nsearchives.nseindia.com/content/historical/DERIVATIVES/2023/SEP/fo08SEP2023bhav.csv.zip"
    
    outputFileName := "nse.csv.zip"

    err = downloadFile_without_client(url, outputFileName, date)
    if err != nil {
        fmt.Println("Error downloading file:", err)
        return err
    }

    fmt.Println("File downloaded successfully:", outputFileName)
	return nil
}

func downloadFile_without_client(url string, outputFileName, date string) error {
    // Create the output file
    outputFile, err := os.Create(outputFileName)
    if err != nil {
        return err
    }
    defer outputFile.Close()

    // Send HTTP GET request to the URL
    // client := &http.Client{}
    // req, err := http.NewRequest("GET", url, nil)
    // if err != nil {
    //     return err
    // }

    // // Set a User-Agent header (mimicking a web browser)
    // req.Header.Set("User-Agent", "Mozilla/5.0")

    // // Perform the HTTP request
    // response, err := client.Do(req)
    // if err != nil {
    //     return fmt.Errorf("HTTP request error: %s", err)
    // }
    // defer response.Body.Close()
    // Create a new HTTP GET request
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error downloading file: %v\n", err)
		
	}
	defer response.Body.Close()

    // Check if the response status code is not 200 OK
    if response.StatusCode != http.StatusOK {
        return fmt.Errorf("Data does not exist for %s", date)
    }

    // Parse the Content-Length header
    contentLengthStr := response.Header.Get("Content-Length")
    contentLength, err := strconv.ParseInt(contentLengthStr, 10, 64)
    if err != nil {
        return err
    }

    // Copy the response body to the output file and validate the size
    n, err := io.Copy(outputFile, response.Body)
    if err != nil {
        return err
    }

    // Check if the downloaded file size matches the expected size
    if n != contentLength {
        return fmt.Errorf("Downloaded file size (%d bytes) does not match expected size (%d bytes)", n, contentLength)
    }

	err  = extract_zip_file()
	if err != nil{
		return err
	}

    return nil
}

// func extract_zip_file() error {
//     zipFileName := "nse.csv.zip"

//     // Unzip the downloaded ZIP file
//     err := unzipFile(zipFileName)
//     if err != nil {
//         fmt.Println("Error unzipping file:", err)
//         return err
//     }
// 	return nil
// }

// func unzipFile(zipFileName string) error {
//     // Open the ZIP file
//     zipFile, err := zip.OpenReader(zipFileName)
//     if err != nil {
//         return err
//     }
//     defer zipFile.Close()

//     // Iterate through the files in the ZIP archive
//     for _, file := range zipFile.File {
//         // Open the ZIP file entry
//         zipEntry, err := file.Open()
//         if err != nil {
//             return err
//         }
//         defer zipEntry.Close()

//         // Create the output file in the root directory with the same name as the ZIP file
//         outFile, err := os.Create("output.csv")
//         if err != nil {
//             return err
//         }
//         defer outFile.Close()

//         // Copy the contents to the output file
//         _, err = io.Copy(outFile, zipEntry)
//         if err != nil {
//             return err
//         }
//     }

//     return nil
// }


