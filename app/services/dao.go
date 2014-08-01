package services

import (
	"RPGithub/app/db"
	"RPGithub/app/model"
	"github.com/revel/revel"
	"strings"
)

// GetUser gets a user from the database
func GetUser(username string) *model.User {
	var user *model.User

	userData := db.Database.Get(strings.ToLower(username), db.COLLECTION_USER)
	err := userData.One(&user)
	if err != nil {
		return nil
	}

	return user
}

// GetRepository gets a new repository from the database
func GetRepository(id int) *model.Repository {
	var repository *model.Repository

	repositoryData := db.Database.Get(id, db.COLLECTION_REPOSITORY)
	err := repositoryData.One(&repository)
	if err != nil {
		return nil
	}

	return repository
}

// GetUserRepositories gets a list of repositories from the given user
func GetUserRepositories(username string) []*model.Repository {
	var repositories []*model.Repository

	data := db.Database.GetQuery(map[string]string{"owner": strings.ToLower(username)}, db.COLLECTION_REPOSITORY)
	data.All(&repositories)

	return repositories
}

// UpdateUser updates the given user from the database
func UpdateUser(user *model.User) error {
	err := db.Database.Update(user.Id, user, db.COLLECTION_USER)
	if err != nil {
		revel.ERROR.Fatalf("Error while updating data : %s", err.Error())
		return err
	}

	return nil
}

// UpdateRepository updates the repository in the database
func UpdateRepository(repository *model.Repository) error {
	err := db.Database.Update(repository.Id, repository, db.COLLECTION_REPOSITORY)
	if err != nil {
		revel.ERROR.Fatalf("Error while updating data : %s", err.Error())
		return err
	}

	return nil
}

// RegisterRepository register in database a new repository
func RegisterRepository(repository *model.Repository) error {
	err := db.Database.Set(repository, db.COLLECTION_REPOSITORY)
	if err != nil {
		revel.ERROR.Fatalf("Error while saving new repository : %s", err.Error())
		return err
	}

	return nil
}

// RegisterRepository register in database a new repository
func RegisterUser(user *model.User) error {
	err := db.Database.Set(user, db.COLLECTION_USER)
	if err != nil {
		revel.ERROR.Fatalf("Error while saving new user : %s", err.Error())
		return err
	}

	return nil
}
