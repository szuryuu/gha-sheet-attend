package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	spreadsheetID = "1-pqF4HCNXjvON226ZTlkifF0VKiKPvJkOdOxm9N0kl4"
	sheetName     = "Shafwan Ilham Dzaky"
	dateFormat    = "1/2/2006"
)

func sheetService(ctx context.Context) (*sheets.Service, error) {
	creds_base64 := os.Getenv("GCP_SA_KEY")
	if creds_base64 == "" {
		return nil, fmt.Errorf("GCP_SA_KEY environment variable not set")
	}

	creds_json_str, err := base64.StdEncoding.DecodeString(creds_base64)
	if err != nil {
		return nil, fmt.Errorf("Error decoding data: %v", err)
	}

	creds_google, err := google.CredentialsFromJSON(ctx, creds_json_str, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, fmt.Errorf("Error creating credentials: %v", err)
	}

	srv, err := sheets.NewService(ctx, option.WithTokenSource(creds_google.TokenSource))
	if err != nil {
		return nil, fmt.Errorf("Error creating service: %v", err)
	}

	return srv, nil
}

func getNextRowNumber(service *sheets.Service) (int, error) {
	readRange := fmt.Sprintf("%s!A1:A", sheetName)
	resp, err := service.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		return 0, err
	}

	const headerRows = 4
	if len(resp.Values) <= headerRows {
		return 1, nil
	}

	lastRow := resp.Values[len(resp.Values)-1]
	if len(lastRow) == 0 {
		return len(resp.Values) + 1, nil
	}

	lastNumberStr, ok := lastRow[0].(string)
	if !ok {
		return 0, fmt.Errorf("Invalid data type for last number")
	}

	lastNumber, err := strconv.Atoi(lastNumberStr)
	if err != nil {
		return 0, fmt.Errorf("Error converting last number to int: %v", err)
	}

	return lastNumber + 1, nil
}

func main() {
	ctx := context.Background()
	sheetService, err := sheetService(ctx)
	if err != nil {
		fmt.Println("Error creating service:", err)
		return
	}

	nextNumberRow, err := getNextRowNumber(sheetService)
	if err != nil {
		fmt.Println("Error getting next row number:", err)
		return
	}
	log.Printf("Next number row: %d", nextNumberRow)

	attendRecord := os.Getenv("INPUT_ATTEND_RECORD")
	startTime := os.Getenv("INPUT_START_TIME")
	endTime := os.Getenv("INPUT_END_TIME")
	additionalInfo := os.Getenv("INPUT_ADDITIONAL_INFO")

	if attendRecord == "Libur" {
		startTime = ""
		endTime = ""
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}

	todayDate := time.Now().In(loc).Format(dateFormat)
	newRow := sheets.ValueRange{
		Values: [][]any{
			{
				nextNumberRow,
				todayDate,
				attendRecord,
				startTime,
				endTime,
				additionalInfo,
			},
		},
	}

	appendCall := sheetService.Spreadsheets.Values.Append(spreadsheetID, sheetName, &newRow).
		ValueInputOption("USER_ENTERED")
	_, err = appendCall.Do()
	if err != nil {
		fmt.Println("Error appending row:", err)
		return
	}

	log.Printf("Row appended successfully")
}
