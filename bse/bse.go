package bse

import (
	"time"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"archive/zip"
	"encoding/csv"
)

func Bse(date string) error {
    // Parse the input date
	fmt.Println(date)
    inputDateFormat := "02-01-2006"
    inputDate, err := time.Parse(inputDateFormat, date)
    if err != nil {
        return fmt.Errorf("Error parsing input date: %s", date)
    }

    // Format the parsed date as "110923"
    outputDateFormat := "020106"
    dateString := inputDate.Format(outputDateFormat)

	zipURL := fmt.Sprintf("https://www.bseindia.com/download/BhavCopy/Equity/EQ%s_CSV.ZIP", dateString)

	// Create an HTTP client
	client := &http.Client{}

	// Create a request with headers
	req, err := http.NewRequest("GET", zipURL, nil)
	if err != nil {
		return fmt.Errorf("Failed to create a request: %v", err)
		
	}
	req.Header.Set("Referer", "https://www.bseindia.com") // Set the Referer header
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36") // Set a user-agent header

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to download ZIP file. Status code: %d", resp.StatusCode)
	}

	// Create a file to save the ZIP content
	file, err := os.Create("bse_bhavcopy.zip")
	if err != nil {
		return fmt.Errorf("Failed to create ZIP file: %v", err)
	}
	defer file.Close()

	// Copy the ZIP content to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to save ZIP file: %v", err)
	}

	fmt.Println("ZIP file downloaded successfully.")
	return nil
}


func ConvertZipToCsv(fileName string) error{
	zipFileName := fileName 

	// Open the ZIP file for reading
	zipFile, err := zip.OpenReader(zipFileName)
	if err != nil {
		return fmt.Errorf("Failed to open ZIP file: %v", err)
	}
	defer zipFile.Close()

	// Find and extract the CSV file from the ZIP archive
	var csvFile *zip.File
	for _, file := range zipFile.File {
		if strings.HasSuffix(file.Name, ".CSV") {
			csvFile = file
			break
		}
	}

	if csvFile == nil {
		return fmt.Errorf("CSV file not found in the ZIP archive")
	}

	// Open the CSV file for reading
	csvFileReader, err := csvFile.Open()
	if err != nil {
		return fmt.Errorf("Failed to open CSV file in ZIP archive: %v", err)
	}
	defer csvFileReader.Close()

	// Create a CSV reader for the CSV file
	csvReader := csv.NewReader(csvFileReader)

	// Create a CSV writer for the output CSV file
	outputFileName := "output.csv"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return fmt.Errorf("Failed to create output CSV file: %v", err)
	}
	defer outputFile.Close()

	csvWriter := csv.NewWriter(outputFile)
	defer csvWriter.Flush()

	// Read and write each line from the CSV file
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("Error reading CSV record: %v", err)
		}

		if err := csvWriter.Write(record); err != nil {
			return fmt.Errorf("Error writing CSV record: %v", err)
		}
	}

	fmt.Printf("CSV conversion completed. Output file: %s\n", outputFileName)
	return nil

}