// Steve Phillips / elimisteve
// 2013.11.10

package types

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

var (
	// TODO: Replace with legit DB
	hypeServices = map[string]*HypeService{} // map from URL to *HypeService
)

type HypeService struct {
	// POSTed by user
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`

	// Filled in by this server
	CreatedBy  string    `json:"created_by"`  // Read IP of original POSTer
	ModifiedBy string    `json:"modified_by"` // Read IP of updater
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	LastSeen   time.Time `json:"last_seen"`
}

func (hs *HypeService) Save() (err error) {
	defer func() {
		if err != nil {
			log.Printf("Failed to add `%s` -- %#v\n", hs.URL, hs)
			return
		}
		log.Printf("Successfully added `%s` -- %#v\n", hs.URL, hs)
	}()

	oldHS, found := hypeServices[hs.URL]
	if !found {
		if err = hs.populateFields(); err != nil {
			return
		}
		hypeServices[hs.URL] = hs
		return
	}
	if err = hs.populateFields(); err != nil {
		return
	}
	hs.updateFromOld(oldHS)
	hypeServices[hs.URL] = hs
	return
}

func (hs *HypeService) Validate() error {
	if hs.Name == "" || hs.URL == "" || hs.Description == "" {
		return fmt.Errorf("name, url, and description fields must be populated")
	}
	return nil
}

func (hs *HypeService) populateFields() error {
	now := time.Now()
	// TODO: Only populate if blank
	hs.CreatedBy = ""  // TODO: List IP of POSTer
	hs.ModifiedBy = "" // TODO: List IP of POSTer
	hs.ModifiedAt = now
	hs.CreatedAt = now
	return nil
}

func (hs *HypeService) updateFromOld(oldHS *HypeService) {
	// New version wins, except `CreatedBy` and `CreatedAt` fields
	hs.CreatedBy = oldHS.CreatedBy
	hs.CreatedAt = oldHS.CreatedAt
}

func ServicesList() ([]*HypeService, error) {
	list := make([]*HypeService, 0, len(hypeServices))
	for _, hs := range hypeServices {
		list = append(list, hs)
	}
	return list, nil
}

func DumpDB() error {
	services, err := ServicesList()
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(services)
	if err != nil {
		return err
	}

	log.Printf("DB dump:\n\n%s\n\n", jsonData)

	// Write JSON to disk

	filename := fmt.Sprintf("hyperactive-%s.json", time.Now().Format(time.RFC3339))
	return ioutil.WriteFile(filename, jsonData, 0644)
}
