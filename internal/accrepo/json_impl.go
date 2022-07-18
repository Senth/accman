package accrepo

import (
	"encoding/json"
	"github.com/Senth/accman/models"
	"log"
	"os"
)

const accountsFile = "accounts.json"

type accRepoJSONImpl struct {
	accounts []models.Account
}

// NewJSONImpl returns a new JSON implementation of the AccRepo interface
func NewJSONImpl() AccRepo {
	repo := &accRepoJSONImpl{}
	repo.load()
	return repo
}

func (a accRepoJSONImpl) Get(accountNumber models.AccountNumber) *models.Account {
	for _, account := range a.accounts {
		if account.Number == accountNumber {
			return &account
		}
	}
	return nil
}

func (a accRepoJSONImpl) GetAll() []models.Account {
	return a.accounts
}

// load the accounts from the JSON file
func (a *accRepoJSONImpl) load() {
	jsonFile, err := os.Open(accountsFile)
	if err != nil {
		log.Panicf("failed to open accounts file: %v", err)
	}
	err = json.NewDecoder(jsonFile).Decode(&a.accounts)
	if err != nil {
		log.Panicf("failed to decode accounts file: %v", err)
	}
}
