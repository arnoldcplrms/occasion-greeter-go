package greeter

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var sheetColumnMap = map[string]string{
	"Male's First Name":                              "maleName",
	"Male's Last Name":                               "maleLastName",
	"Male's Nickname (If Applicable)":                "maleNickname",
	"Male's Birthday":                                "maleBirthday",
	"Male's Email Address":                           "maleEmail",
	"Male's Profile Picture":                         "maleProfilePicture",
	"Female's First Name":                            "femaleName",
	"Female's Last Name":                             "femaleLastName",
	"Female's Nickname (If Applicable)":              "femaleNickname",
	"Female's Birthday":                              "femaleBirthday",
	"Female's Email Address":                         "femaleEmail",
	"Female's Profile Picture":                       "femaleProfilePicture",
	"Couple's Profile Picture":                       "weddingProfilePicture",
	"Wedding Anniversary (If engaged, set the future date/ If Not Applicable Skip)": "weddingAnniversary",
}

func loadCSVFromReader(r io.Reader) ([]OccasionData, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read csv: %w", err)
	}

	if len(rows) == 0 {
		return []OccasionData{}, nil
	}

	headers := rows[0]
	results := make([]OccasionData, 0, len(rows)-1)

	for idx, row := range rows[1:] {
		rowNumber := idx + 2
		record := make(map[string]string, len(headers))
		for i, h := range headers {
			if i < len(row) {
				key, ok := sheetColumnMap[h]
				if !ok {
					key = h
				}
				record[key] = strings.TrimSpace(row[i])
			}
		}

		maleName, ok := record["maleName"]
		if !ok || maleName == "" {
			return nil, fmt.Errorf("CSV row %d: missing required column 'maleName'", rowNumber)
		}
		femaleName, ok := record["femaleName"]
		if !ok || femaleName == "" {
			return nil, fmt.Errorf("CSV row %d: missing required column 'femaleName'", rowNumber)
		}

		maleLastName := record["maleLastName"]
		if maleLastName == "" {
			maleLastName = record["coupleLastName"]
		}

		results = append(results, OccasionData{
			MaleName:              maleName,
			MaleLastName:          maleLastName,
			FemaleName:            femaleName,
			FemaleLastName:        record["femaleLastName"],
			MaleNickname:          record["maleNickname"],
			FemaleNickname:        record["femaleNickname"],
			MaleBirthday:          record["maleBirthday"],
			FemaleBirthday:        record["femaleBirthday"],
			MaleEmail:             record["maleEmail"],
			FemaleEmail:           record["femaleEmail"],
			MaleProfilePicture:    record["maleProfilePicture"],
			FemaleProfilePicture:  record["femaleProfilePicture"],
			CoupleLastName:        maleLastName,
			WeddingAnniversary:    record["weddingAnniversary"],
			WeddingProfilePicture: record["weddingProfilePicture"],
		})
	}

	return results, nil
}

func LoadOccasionsDataFromFile(path string) ([]OccasionData, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open csv: %w", err)
	}
	defer f.Close()
	return loadCSVFromReader(f)
}

func LoadOccasionsData() ([]OccasionData, error) {
	source := os.Getenv("CSV_URL")
	if source == "" {
		return nil, fmt.Errorf("CSV_URL environment variable is not set")
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(source)
	if err != nil {
		return nil, fmt.Errorf("download CSV_URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed to download CSV_URL (status %d)", resp.StatusCode)
	}

	return loadCSVFromReader(resp.Body)
}
