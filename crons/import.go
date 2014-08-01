package crons

import (
	"RPGithub/app/model"
	"RPGithub/app/services"
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
	Payload    Payload         `json:"payload"`
}

type ActorAttributes struct {
	Login string `json:"login"`
	Type  string `json:"type"`
}

type Payload struct {
	Action string `json:"action"`
}

type Repository struct {
	Language     string `json:"language"`
	Organization string `json:"organization"`
	Stars        int    `json:"stargazers_count"`
	Size         int    `json:"size"`
	Id           int    `json:"id"`
	Url          string `json:"url"`
	Description  string `json:"description"`
	Owner        string `json:"owner"`
	Name         string `json:"name"`
	Wiki         bool   `json:"has_wiki"`
	Downloads    bool   `json:"has_downloads"`
	Forks        int    `json:"forks_count"`
	Issues       int    `json:"open_issues_count"`
	IsFork       bool   `json:"fork"`
}

var steps [12]int = [12]int{5, 10, 30, 50, 100, 300, 500, 1000, 3000, 5000, 10000, 100000000}

// Structure that implements the Job interface
type Import struct{}

// Run is the method called by the cronjob
// It downloads the archive file and update the database
func (this Import) Run() {
	// fullPath, err := this.Download("")
	// if err != nil {
	// 	revel.ERROR.Fatal(err)
	// 	return
	// }

	fullPath := "imports/2014-07-27-1.json.gz"
	data, err := this.Ungzip(fullPath)
	if err != nil {
		revel.ERROR.Fatal(err)
	}

	this.Parse(data, true)
}

// Download Downloads the archive file from githubarchive
func (this *Import) Download(date string) (string, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	url := fmt.Sprintf("%s/%s.json.gz", revel.Config.StringDefault("imports.url", "http://data.githubarchive.org"), date)

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

// Ungzip ungzip the given gzipped file and returns its content
func (this *Import) Ungzip(file string) (string, error) {
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

// Parse the given json string and updates the database with it
func (this *Import) Parse(data string, ranking bool) error {
	array := strings.Split(data, "\n")
	var total int = len(array)
	var err error

	if ranking == true {
		// Clear all events
		services.ClearEventDay()
	}

	for key, event := range array {
		revel.INFO.Printf("-> Event %d/%d", key, total)

		var jsonmap ImportedData
		_ = json.Unmarshal([]byte(event), &jsonmap)

		// Only user type for moment
		if jsonmap.User.Type != "User" {
			continue
		}

		// ------------------------------------- GET USER
		user := services.GetUser(strings.ToLower(jsonmap.User.Login))
		if user == nil {
			revel.INFO.Printf("Get user %s : %s", jsonmap.User.Login, err)

			// New user
			user = model.NewUser(jsonmap.User.Login)

			// Register the user
			services.RegisterUser(user)
		}

		// ------------------------------------- GET REPOSITORY
		repository := services.GetRepository(jsonmap.Repository.Id)
		if repository == nil {
			revel.INFO.Printf("Get repository %s : %s", jsonmap.Repository.Id, err)

			// New repository
			repository = model.NewRepository(
				jsonmap.Repository.Id,
				jsonmap.Repository.Name,
			)

			// Register the repository
			services.RegisterRepository(repository)
		}

		repository.Size = jsonmap.Repository.Size
		repository.Url = jsonmap.Repository.Url
		repository.Language = jsonmap.Repository.Language
		repository.Owner = strings.ToLower(jsonmap.Repository.Owner)
		repository.Organization = strings.ToLower(jsonmap.Repository.Organization)
		repository.Wiki = jsonmap.Repository.Wiki
		repository.Downloads = jsonmap.Repository.Downloads
		repository.Forks = jsonmap.Repository.Forks
		repository.Stars = jsonmap.Repository.Stars
		repository.Issues = jsonmap.Repository.Issues
		repository.IsFork = jsonmap.Repository.IsFork
		repository.Description = jsonmap.Repository.Description

		language := user.GetLanguage(jsonmap.Repository.Language)

		// --------------------------------- UPDATES
		var xp int
		switch strings.ToLower(jsonmap.Type) {
		case "pushevent":
			language.Events.Pushes += 1
			for key, value := range steps {
				if jsonmap.Repository.Stars < value {
					xp = 50 * (key + (key + 1))
					break
				}
			}

		case "createevent":
			xp = 1
			language.Events.Creates += 1

		case "deleteevent":
			xp = 1
			language.Events.Deletes += 1

		case "issuesevent":
			language.Events.Issues += 1
			for key, value := range steps {
				if jsonmap.Repository.Stars < value {
					xp = 5 * (key + (key + 1))
					break
				}
			}

		case "issuecommentevent":
			language.Events.Comments += 1
			for key, value := range steps {
				if jsonmap.Repository.Stars < value {
					xp = 1 * (key + (key + 1))
					break
				}
			}

		case "watchevent":
			language.Events.Stars += 1
			xp = 1

		case "forkevent":
			language.Events.Forks += 1
			xp = 5

		case "pullrequestevent":
			language.Events.Pullrequests += 1
			for key, value := range steps {
				if jsonmap.Repository.Stars < value {
					xp = 300 * (key + (key + 1))
					break
				}
			}

		case "pullrequestreviewcommentevent":
			language.Events.Comments += 1
			for key, value := range steps {
				if jsonmap.Repository.Stars < value {
					xp = 1 * (key + (key + 1))
					break
				}
			}
		}

		// Register a daily event
		if ranking == true {
			services.RegisterEventDay(model.NewEventDay(
				strings.ToLower(jsonmap.Type),
				language.Name,
				user.Id,
				xp,
			))
		}

		// Updates level & experience
		language.AddExperience(xp)
		user.AddExperience(xp)

		// Updates database data
		services.UpdateUser(user)
		services.UpdateRepository(repository)
	}

	return nil
}
