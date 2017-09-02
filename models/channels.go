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

// AllChannels queries the database for all Channels.
// It returns a list of Channels.
func AllChannels(db *gorm.DB) ([]*Channel, error) {
	channels := []*Channel{}
	if err := db.Find(&channels).Error; err != nil {
		return nil, err
	}
	return channels, nil
}

// FindChannelByUsername queries the database for any Channel with the matching
// username passed in.
// It returns a Channel and any errors that might have been encountered.
func FindChannelByUsername(db *gorm.DB, username string) (*Channel, error) {
	channel := &Channel{}
	if err := db.Where("username = ?", username).First(&channel).Error; err != nil {
		return nil, err
	}
	return channel, nil
}

// FindChannelByStreamKey queries the database for any Channel with the matching
// stream key passed in.
// It returns the Channel and any errors that might have been encountered.
func FindChannelByStreamKey(db *gorm.DB, key string) (*Channel, error) {
	channel := &Channel{}
	if err := db.Where("stream_key = ?", key).First(&channel).Error; err != nil {
		return nil, err
	}
	return channel, nil
}
