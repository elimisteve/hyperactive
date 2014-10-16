// Steve Phillips / elimisteve
// 2013.11.10

package types

import (
	"fmt"
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
	PostedBy   string    `json:"posted_by"`   // Read IP of original POSTer
	ModifiedBy string    `json:"modified_by"` // Read IP of updater
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	LastSeen   time.Time `json:"last_seen"`
}

func (hs *HypeService) Save() (err error) {
	defer func() {
		if err != nil {
			fmt.Printf("%s failed to be updated or saved\n", hs.URL)
			return
		}
		fmt.Printf("%s successfully updated or saved\n", hs.URL)
	}()

	if hs.Name == "" || hs.URL == "" || hs.Description == "" {
		err = fmt.Errorf("name, url, and description fields must be populated\n")
		return
	}

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

func (hs *HypeService) populateFields() error {
	now := time.Now()
	// TODO: Only populate if blank
	hs.PostedBy = ""   // TODO: List IP of POSTer
	hs.ModifiedBy = "" // TODO: List IP of POSTer
	hs.ModifiedAt = now
	hs.CreatedAt = now
	return nil
}

func (hs *HypeService) updateFromOld(oldHS *HypeService) {
	// New version wins, except `PostedBy` and `CreatedAt` fields
	hs.PostedBy = oldHS.PostedBy
	hs.CreatedAt = oldHS.CreatedAt
}

func ServicesList() ([]*HypeService, error) {
	list := make([]*HypeService, len(hypeServices))
	var ndx int
	for _, hs := range hypeServices {
		list[ndx] = hs
		ndx++
	}
	return list, nil
}
