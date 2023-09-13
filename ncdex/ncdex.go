package ncdex

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"archive/zip"
	"strings"
	"time"
)

func Ncdex(date string) error  {
	currentDate := time.Now()
    givenDate, err := time.Parse("02-01-2006", date)
    if err != nil {
        return fmt.Errorf("Error parsing given date: %s", err)
    }
	todayDate := currentDate.Format("2006-01-02")

    // Parse the current date string back into a time.Time value
    currDate, err := time.Parse("2006-01-02", todayDate)
    if givenDate.After(currDate) {
        return fmt.Errorf("Ncdex file is not generated for %s", date)
    } 

    // Extract the year, month, and day
    year :=  givenDate.Format("2006")
    month := givenDate.Format("Jan")
    day := givenDate.Format("02")

	month = strings.ToUpper(month)
	date = fmt.Sprintf("%s-%s-%s",day,month,year )

	url := fmt.Sprintf("https://ncdex.com/markets/bhavcopy?file_type=final&type=bhavcopy&filedate=%s&format=csv_file", date)

	// Define the output file name
	outputFileName := "bhavcopy.zip"

	// Create a new HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error downloading file: %v\n", err)
		
	}
	defer resp.Body.Close()

	// Check if the response status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to get the data for this date : %s\n", date)
		
	}

	// Create or open the output file for writing
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return fmt.Errorf("Error creating output file: %v\n", err)
		
	}
	defer outputFile.Close()

	// Copy the response body (file content) to the output file
	_, err = io.Copy(outputFile, resp.Body)
	if err != nil {
		return fmt.Errorf("Error copying file content: %v\n", err)
	}

	fmt.Printf("File downloaded successfully as %s\n", outputFileName)
	return nil
}
func UnzipAndSaveAsOutputCSV(zipFileName string) error {
    // Open the ZIP file
    zipFile, err := zip.OpenReader(zipFileName)
    if err != nil {
        return err
    }
    defer zipFile.Close()

    // Create the output file
    outputFileName := "output.csv"
    outFile, err := os.Create(outputFileName)
    if err != nil {
        return err
    }
    defer outFile.Close()

    // Search for a CSV file in the ZIP archive and extract its contents
    for _, file := range zipFile.File {
        if strings.HasSuffix(file.Name, ".csv") {
            // Open the ZIP file entry
            zipEntry, err := file.Open()
            if err != nil {
                return err
            }
            defer zipEntry.Close()

            // Copy the contents to the output file
            _, err = io.Copy(outFile, zipEntry)
            if err != nil {
                return err
            }
			break;

            // You can choose to break here if you want to extract the first CSV file found
            // break
        }
    }

    return nil
}