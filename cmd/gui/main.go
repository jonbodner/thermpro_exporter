package main

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"thermpro_exporter"
	"time"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Data Exporter")

	label1 := widget.NewLabel("Start Date:")
	start := widget.NewEntry()
	start.SetPlaceHolder("YYYY-MM-DD")
	fmt.Println(start.Size())
	label2 := widget.NewLabel("End Date:")
	end := widget.NewEntry()
	end.SetPlaceHolder("YYYY-MM-DD")
	label3 := widget.NewLabel("Interval")
	interval := widget.NewEntry()
	interval.SetPlaceHolder("1m")
	cancel := widget.NewButton("Cancel", func() {
		myApp.Quit()
	})
	run := widget.NewButton("Export", func() {
		err := process(start.Text, end.Text, interval.Text)
		if err == nil {
			myApp.Quit()
		}
		// todo dialog for error
		fmt.Println(err)
	})
	form := container.New(layout.NewFormLayout(), label1, start, label2, end, label3, interval)

	buttons := container.New(layout.NewGridLayoutWithColumns(2), cancel, run)
	main := container.NewBorder(nil, buttons, nil, nil, form)
	myWindow.SetContent(main)
	myWindow.Resize(fyne.Size{
		Width:  300,
		Height: myWindow.Content().Size().Height,
	})
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun()
}

func process(start, end, interval string) error {
	var outErr []error
	var err error
	var startDate time.Time
	if len(start) == 0 {
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -1)
	} else {
		startDate, err = time.Parse(time.RFC3339, start+"T00:00:00-04:00")
		if err != nil {
			outErr = append(outErr, err)
		}
	}
	fmt.Println(startDate)
	var endDate time.Time
	if len(end) == 0 {
		now := time.Now()
		endDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	} else {
		endDate, err = time.Parse(time.RFC3339, end+"T00:00:00-04:00")
		if err != nil {
			outErr = append(outErr, err)
		}
	}
	fmt.Println(endDate)
	var toSkip time.Duration
	if len(interval) == 0 {
		toSkip = time.Minute
	} else {
		toSkip, err = time.ParseDuration(interval)
		if err != nil {
			outErr = append(outErr, err)
		}
	}
	fmt.Println(toSkip)
	if len(outErr) != 0 {
		return errors.Join(outErr...)
	}
	err = thermpro_exporter.GenerateCSV(startDate, endDate, toSkip)
	if err != nil {
		return err
	}
	return nil
}
