package thermpro_exporter

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	mcsv "thermpro_exporter/internal/csv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

/*
Writing this helps me clarify my thoughts. I think I found the problem.
Numbers keeps asking me to format the data and I think this is why.
There is a format difference between Bodner generated data, and ThermPro generated data.
The data fields are {date, temperature, humidity}.  Bodner formatted data:
Saturday, June 15, 2024 00:00:00
76.10
58

	ThermPro formatted data: 01-13-2024 0:00 AM

65.3
39
I am trying to find a way to change the format using numbers.
*/
const dateformat = `01-02-2006 03:04 PM`

func GenerateCSV(startDate time.Time, endDate time.Time, interval time.Duration) error {

	// find the file
	fileName, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	containerDir := fileName + "/Library/Containers"
	dirContents, err := os.ReadDir(containerDir)
	if err != nil {
		return err
	}
	for _, v := range dirContents {
		if v.IsDir() {
			maybeFilename := containerDir + string(os.PathSeparator) + v.Name() + "/Data/Documents/LocalData.db"
			_, err := os.Stat(maybeFilename)
			if !errors.Is(err, os.ErrNotExist) {
				fileName = maybeFilename
				break
			}
		}
	}
	fmt.Println(fileName)
	database, err := sql.Open("sqlite3", fileName)
	if err != nil {
		return err
	}
	defer database.Close()

	// find name
	rows, err := database.Query(`SELECT name FROM sqlite_schema 
WHERE type IN ('table','view') 
AND name NOT LIKE 'sqlite_%'
ORDER BY 1;
`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var correctName string

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return err
		}
		if strings.HasPrefix(name, "NewLocalData") {
			correctName = name
		}
	}
	err = rows.Err()
	if err != nil {
		return err
	}

	type Data struct {
		DateTime    string `csv:"date_time"`
		Temperature string `csv:"temperature"`
		Humidity    string `csv:"humidity"`
	}

	// get data
	rows, err = database.Query(fmt.Sprintf(`select Tem, Hum, timeInterval 
from %s 
where isValid = 1 
order by timeInterval asc;
`, correctName))
	if err != nil {
		return err
	}
	defer rows.Close()

	var results []Data
	nextTime := startDate
	for rows.Next() {
		var curData Data
		var timestamp int64
		err = rows.Scan(&curData.Temperature, &curData.Humidity, &timestamp)
		if err != nil {
			return err
		}
		// convert from C to F
		tempC, _ := strconv.ParseFloat(curData.Temperature, 64)
		tempF := (tempC * 1.8) + 32
		curData.Temperature = fmt.Sprintf("%.2f", tempF)
		foundTime := time.UnixMilli(timestamp * 1000)
		if !foundTime.Before(nextTime) {
			curData.DateTime = foundTime.Format(dateformat)
			results = append(results, curData)
			nextTime = foundTime.Add(interval)
		}
		if nextTime.After(endDate) {
			break
		}
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	contents, err := mcsv.Marshal(results)
	if err != nil {
		return err
	}
	var out bytes.Buffer
	cout := csv.NewWriter(&out)
	err = cout.WriteAll(contents)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("%d-%02d-%02d to %d-%02d-%02d.csv",
		startDate.Year(), startDate.Month(), startDate.Day(),
		endDate.Year(), endDate.Month(), endDate.Day())
	err = os.WriteFile(filename, out.Bytes(), 0700)
	if err != nil {
		return err
	}
	// open the file in its default app
	cmd := exec.Command("open", filename)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
