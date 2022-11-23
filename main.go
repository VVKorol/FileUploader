package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("file_input.txt")
	if err != nil {
		fmt.Println("File file_input.txt not found! Error")
	}
	defer file.Close()

	unuploadedFiles := []string{}

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		url := scanner.Text()
		fmt.Print("Upload ", url)
		count, err := DownloadFile(url)
		if err != nil {
			fmt.Println(" ---> Error")
			unuploadedFiles = append(unuploadedFiles, fmt.Sprint(url, " Err: ", err))
		} else {
			fmt.Println(" ---> Done,", count)
		}
	}

	for _, unuploadedFile := range unuploadedFiles {
		fmt.Println(unuploadedFile)
	}
}

func getFileFormat(url string) string {
	last := strings.LastIndex(url, ".")
	return url[last:]
}

func getFileName(url string) string {
	last := strings.LastIndex(url, "/")
	return url[last:]
}

func DownloadFile(url string) (int64, error) {

	if getFileFormat(url) == "" {
		return 0, errors.New("Error to get format (unknown format)")
	}

	filepath := "." + getFileName(url)

	// Get the data
	resp, err := http.Get(url)

	if err != nil {
		return 0, errors.New("Error in upload: " + err.Error())
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return 0, errors.New("Error in save to file: " + err.Error())
	}
	defer out.Close()

	// Write the body to file
	count, err := io.Copy(out, resp.Body)
	out.Sync()
	return count, err
}
