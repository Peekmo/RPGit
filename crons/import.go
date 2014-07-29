package crons

import (
	"RPGithub/app/db"
	"RPGithub/app/model"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/revel/revel"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type ImportedData struct {
	Type       string          `json:"type"`
	User       ActorAttributes `json:"actor_attributes"`
	Repository Repository      `json:"repository"`
}

type ActorAttributes struct {
	Login string `json:"login"`
	Type  string `json:"type"`
}

type Repository struct {
	Language     string `json:"language"`
	Organization string `json:"organization"`
	Stars        int    `json:"stargazers"`
	Size         int    `json:"size"`
	Id           int    `json:"id"`
	Url          string `json:"url"`
	Owner        string `json:"owner"`
	Name         string `json:"name"`
}

// Structure that implements the Job interface
type Import struct{}

// Run is the method called by the cronjob
// It downloads the archive file and update the database
func (this Import) Run() {
	// fullPath, err := this.download()
	// if err != nil {
	// 	revel.ERROR.Fatal(err)
	// 	return
	// }

	fullPath := "imports/2014-07-27-1.json.gz"
	data, err := this.ungzip(fullPath)
	if err != nil {
		revel.ERROR.Fatal(err)
	}

	this.parse(data)
}

// download Downloads the archive file from githubarchive
func (this *Import) download() (string, error) {
	date := time.Now().Format("2006-01-02")
	url := fmt.Sprintf("%s/%s-1.json.gz", revel.Config.StringDefault("imports.url", "http://data.githubarchive.org"), date)

	tokens := strings.Split(url, "/")
	file := tokens[len(tokens)-1]
	folder := revel.Config.StringDefault("imports.folder", "imports")
	fullPath := folder + "/" + file

	// Checks if the folder exists
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		if err := os.Mkdir(folder, 0660); err != nil {
			return fullPath, err
		}
	}

	// Creates a file
	revel.INFO.Printf("Creating file %s", file)
	output, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer output.Close()

	// Get the archive from githubarchive
	revel.INFO.Println("Downloading the file...")
	response, err := http.Get(url)
	if err != nil {
		return fullPath, err
	}
	defer response.Body.Close()

	// Write the file into the created one
	revel.INFO.Println("Copying the file...")
	bytes, err := io.Copy(output, response.Body)
	if err != nil {
		return fullPath, err
	}

	revel.INFO.Printf("File's downloading done (%s bytes)", bytes)
	return fullPath, nil
}

// ungzip ungzip the given gzipped file and returns its content
func (this *Import) ungzip(file string) (string, error) {
	// Read the file
	fileReader, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer fileReader.Close()

	reader, err := gzip.NewReader(fileReader)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	barray, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	str := string(barray[:])
	return str, nil
}

// parse the given json string and updates the database with it
func (this *Import) parse(data string) error {
	array := strings.Split(data, "\n")

	for _, event := range array {

		var jsonmap ImportedData
		_ = json.Unmarshal([]byte(event), &jsonmap)

		if jsonmap.User.Type != "User" {
			continue
		}

		var user *model.User

		userData := db.Database.Get(strings.ToLower(jsonmap.User.Login), db.COLLECTION_USER)
		err := userData.One(&user)
		if err != nil {
			revel.INFO.Printf("Get user %s : %s", jsonmap.User.Login, err)

			// New user
			user = model.NewUser(jsonmap.User.Login)

			// Register the user
			err = db.Database.Set(user, db.COLLECTION_USER)
			if err != nil {
				revel.ERROR.Fatalf("Error while saving new user : %s", err)
			}
		}

		_ = user.GetLanguage(jsonmap.Repository.Language)
		err = db.Database.Update(user.Id, user, db.COLLECTION_USER)
		if err != nil {
			revel.ERROR.Fatalf("Error while updating user : %s", err)
		}

		// fmt.Println(event)
		break
	}

	return nil
}
