package crons

import (
	"fmt"
	"github.com/revel/revel"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Import struct {
}

func (this Import) Run() {
	err := this.download()
	if err != nil {
		revel.ERROR.Fatal(err)
		return
	}
}

// download Downloads the archive file from githubarchive
func (this *Import) download() error {
	date := time.Now().Format("2006-01-02")
	url := fmt.Sprintf("%s/%s-1.json.gz", revel.Config.StringDefault("imports.url", "http://data.githubarchive.org"), date)

	tokens := strings.Split(url, "/")
	file := tokens[len(tokens)-1]
	folder := revel.Config.StringDefault("imports.folder", "imports")

	// Checks if the folder exists
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		if err := os.Mkdir(folder, 0660); err != nil {
			return err
		}
	}

	// Creates a file
	revel.INFO.Printf("Creating file %s", file)
	output, err := os.Create(folder + "/" + file)
	if err != nil {
		return err
	}
	defer output.Close()

	revel.INFO.Println("Downloading the file...")
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	revel.INFO.Println("Copying the file...")
	bytes, err := io.Copy(output, response.Body)
	if err != nil {
		return err
	}

	revel.INFO.Printf("File's downloading done (%s bytes)", bytes)
	return nil
}
