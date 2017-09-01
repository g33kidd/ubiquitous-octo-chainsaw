package models

import "github.com/jinzhu/gorm"

// Channel stores data for an account/channel
type Channel struct {
	gorm.Model

	Username  string
	StreamKey string
	Title     string
	Status    string
}

// AllChannels : Finds all channels in the database
func AllChannels(db *gorm.DB) ([]*Channel, error) {
	return nil, nil
}

// FindChannelByUsername : Finds a channel by it's username
func FindChannelByUsername(db *gorm.DB, username string) (*Channel, error) {
	return nil, nil
}

// FindChannelByStreamKey : Finds a channel by it's stream_key
func FindChannelByStreamKey(db *gorm.DB, key string) (*Channel, error) {
	return nil, nil
}
