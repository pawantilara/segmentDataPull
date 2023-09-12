package mcx

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"strings"
	"encoding/json"
	"io/ioutil"
	"encoding/csv"

	"github.com/gocolly/colly/v2"
)

func Mcx(date string) error {
	today := time.Now()

    // Parse the given date string
    givenDate, err := time.Parse("02-01-2006", date)
    if err != nil {
        return fmt.Errorf("Error parsing given date: %s", err)
    }
	currentDate := today.Format("2006-01-02")

    // Parse the current date string back into a time.Time value
    today, err = time.Parse("2006-01-02", currentDate)
    if givenDate.Equal(today) || givenDate.After(today) {
        return fmt.Errorf("MCX file is not generated for %s", date)
    } 
	// Create a new collector
	c := colly.NewCollector()

	// Visit the Bhavcopy page
	if err := c.Visit("https://www.mcxindia.com/market-data/bhavcopy"); err != nil {
		return fmt.Errorf("Failed to navigate to the Bhavcopy page: %v", err)
		
	}

	// Replace 'desired_date' with the date you want to select in the format "dd/mm/yyyy"
	desiredDate := "09/09/2023"

	// Execute JavaScript to set the value of the date input field
	c.OnHTML("input#txtDate", func(e *colly.HTMLElement) {
		err := e.DOM.SetAttr("value", desiredDate)
		if err != nil {
			log.Printf("Failed to set date input value: %v", err)
		}
	})

	// Define a callback to submit the form
	submitForm := func(e *colly.HTMLElement) {
		// Extract form fields
		formFields := map[string]string{
			"__EVENTTARGET":        e.ChildAttr("input#__EVENTTARGET", "value"),
			"__EVENTARGUMENT":      e.ChildAttr("input#__EVENTARGUMENT", "value"),
			"__LASTFOCUS":          e.ChildAttr("input#__LASTFOCUS", "value"),
			"__VIEWSTATE":          e.ChildAttr("input#__VIEWSTATE", "value"),
			"__VIEWSTATEGENERATOR": e.ChildAttr("input#__VIEWSTATEGENERATOR", "value"),
			"__EVENTVALIDATION":    e.ChildAttr("input#__EVENTVALIDATION", "value"),
			"ctl00$cph_InnerContainerRight$C001$txtDate": desiredDate,
		}

		// Post the form data to retrieve the Bhavcopy
		err := c.Post("https://www.mcxindia.com/market-data/bhavcopy", formFields)
		if err != nil {
			log.Printf("Failed to submit the form: %v", err)
		}
	}

	// Find and submit the form
	c.OnHTML("form#aspnetForm", submitForm)

	// Wait for the form submission to complete
	c.Wait()

	// Save the Bhavcopy to a local file
	file, err := os.Create("mcx.json")
	if err != nil {
		return fmt.Errorf("Failed to create the local file: %v", err)
	}
	defer file.Close()

	// Download the Bhavcopy from the response body and save it to the local file
	response, err := http.Get("https://www.mcxindia.com/market-data/bhavcopy")
	if err != nil {
		return fmt.Errorf("Failed to download Bhavcopy: %v", err)
	}
	defer response.Body.Close()

	// Create a buffer to store the downloaded CSV data
	var buffer strings.Builder
	inPattern := false

	// Read the CSV data into the buffer
	_, err = io.Copy(&buffer, response.Body)
	if err != nil {
		return fmt.Errorf("Failed to read Bhavcopy data: %v", err)
	}

	// Split the buffer contents into lines
	lines := strings.Split(buffer.String(), "\n")

	// Write the lines to the local file until the pattern occurs again
	for _, line := range lines {
		if strings.Contains(line, "var vBC=[{") {
			inPattern = true
		}
		if inPattern && strings.Contains(line, "</script>") {
			break;
		}
		// Write the line to the local file
		if inPattern{
			_, err = file.WriteString(line + "\n")
		
			if err != nil {
				return fmt.Errorf("Failed to write Bhavcopy data to the file: %v", err)
			}
		}
		
	}

	fmt.Println("Bhavcopy data extracted successfully.")
	err = convert_json_to_csv()
	if err !=nil {
		return fmt.Errorf("Failed to convert to csv: %v", err)
	}
	return nil
}

func convert_json_to_csv() error {
	// Read the JSON-like data from the file
	data, err := ioutil.ReadFile("mcx.json")
	if err != nil {
		return fmt.Errorf("Error reading JSON file: %s", err)
		
	}

	// Convert data to string
	jsonLikeData := string(data)

	// Find the start and end markers
	startIndex := strings.Index(jsonLikeData, "var vBC=")
	endIndex := strings.Index(jsonLikeData, ";var vDate=")

	// Check if both markers are found
	if startIndex == -1 || endIndex == -1 {
		fmt.Println("Start or end marker not found in the file.")
	}

	// Extract the JSON content between the markers
	jsonData := jsonLikeData[startIndex+len("var vBC="):endIndex]

	// Parse the JSON data
	var records []map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &records); err != nil {
		return fmt.Errorf("Error parsing JSON: %s", err)
	}

	// Create and open the CSV file for writing
	csvFile, err := os.Create("output.csv")
	if err != nil {
		return fmt.Errorf("Error creating CSV file: %s", err)
	}
	defer csvFile.Close()

	// Create a CSV writer
	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Write the CSV header
	header := []string{"Date", "Symbol", "ExpiryDate", "Open", "High", "Low", "Close", "PreviousClose", "Volume", "VolumeInThousands", "Value", "OpenInterest", "DateDisplay", "InstrumentName", "StrikePrice", "OptionType"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("Error writing CSV header: %s", err)
		
	}

	// Write the data rows
	for _, record := range records {
		var row []string
		for _, key := range header {
			value := fmt.Sprintf("%v", record[key])
			row = append(row, value)
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("Error writing CSV record: %s", err)
			
		}
	}

	fmt.Println("CSV file 'mcx.csv' created successfully.")
	return nil
}
