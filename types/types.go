// Steve Phillips / elimisteve
// 2013.11.10

package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

var (
	// TODO: Replace with legit DB
	hypeServices = map[string]*HypeService{} // map from URL to *HypeService
)

var (
	ErrServiceDuplicate = errors.New("Service already exists; no duplicates allowed")
	ErrServiceNotFound  = errors.New("Service entry not found")
	ErrServiceInvalid   = errors.New("name, url, and description fields must be populated")
)

type HypeService struct {
	// POSTed by user
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`

	// Filled in by this server
	CreatedBy  string     `json:"created_by"`  // Read IP of original POSTer
	ModifiedBy string     `json:"modified_by"` // Read IP of updater
	CreatedAt  time.Time  `json:"created_at"`
	ModifiedAt time.Time  `json:"modified_at"`
	LastSeen   *time.Time `json:"last_seen"`
}

func (hs *HypeService) Save() error {
	if _, found := hypeServices[hs.URL]; found {
		return ErrDuplicateService
	}

	now := time.Now()
	hs.CreatedAt = now
	hs.ModifiedAt = now

	hypeServices[hs.URL] = hs

	log.Printf("Successfully added `%s` -- %#v\n", hs.URL, hs)
	return nil
}

func (hs *HypeService) Update() error {
	oldHS, found := hypeServices[hs.URL]
	if !found {
		return ErrServiceNotFound
	}

	// Set a priori-known values
	now := time.Now()
	hs.ModifiedAt = now

	// Retain these values from the existing *HypeService
	hs.CreatedBy = oldHS.CreatedBy
	hs.CreatedAt = oldHS.CreatedAt

	// Overwrite the existing one
	hypeServices[hs.URL] = hs

	log.Printf("Successfully updated `%s` -- %#v\n", hs.URL, hs)
	return nil
}

func (hs *HypeService) Validate() error {
	if hs.Name == "" || hs.URL == "" || hs.Description == "" {
		return ErrServiceInvalid
	}
	return nil
}

func GetServiceByURL(url string) (*HypeService, error) {
	service, found := hypeServices[url]
	if !found {
		return nil, ErrServiceNotFound
	}
	return service, nil
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

	jsonData, _ := json.Marshal(services)

	log.Printf("DB dump:\n\n%s\n\n", jsonData)

	// Write JSON to disk

	filename := fmt.Sprintf("hyperactive-%s.json", time.Now().Format(time.RFC3339))
	return ioutil.WriteFile(filename, jsonData, 0644)
}
