package factories

import (
	"testing"
	"time"

	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type RawContent struct {
	ArtistName    string
	Album         RawAlbum
	MusicBrainzID string
}

type RawAlbum struct {
	Title         string
	Tracks        []RawTracks
	Type          models.AlbumType
	MusicBrainzID string
}

type RawTracks struct {
	Title         string
	MusicBrainzID string
}

func GetStatsSeedData(t *testing.T, db *gorm.DB) (models.User, models.User) {
	_ = SeedUser(t, db, "admin", "123", "ADMIN", false)
	privateUser := SeedUser(t, db, "private", "123", "USER", false)
	regularUser := SeedUser(t, db, "regular", "123", "USER", true)

	tracks := CreateScrobbleContent(t, db)
	for _, track := range tracks {
		CreateScrobbles(t, db, track.ID, privateUser.ID, 1)
		CreateScrobbles(t, db, track.ID, regularUser.ID, 1)
	}

	return privateUser, regularUser
}

// CreateScrobbles creates a given number of scrobbles for a given user and track.
func CreateScrobbles(t *testing.T, db *gorm.DB, trackID int64, userID int64, count int) {
	t.Helper()

	for i := 0; i < count; i++ {
		scrobble := &models.Scrobble{
			UserID:      userID,
			TrackID:     trackID,
			Origin:      "SUBSONIC",
			ScrobbledAt: time.Now(),
		}

		err := db.Table("scrobbles").Create(&scrobble).Error
		require.NoError(t, err)
	}
}

// CreateScrobbleContent creates the whole chain (track, album artists) used to test stats routes
func CreateScrobbleContent(t *testing.T, db *gorm.DB) []models.Track {
	rawContents := []RawContent{
		{
			ArtistName:    "nothing,nowhere.",
			MusicBrainzID: "90e8bc57-78b0-4992-8784-898d5122d3b3",
			Album: RawAlbum{
				Title:         "miserymaker",
				Type:          "EP",
				MusicBrainzID: "84727998-d10c-4be1-b4f7-7d58e422b02a",
				Tracks: []RawTracks{
					{Title: "DEAD2ME", MusicBrainzID: "58c6d768-a1d4-4de9-98d5-6aa47b050564"},
					{Title: "AURA", MusicBrainzID: "518f8e50-8f16-407f-985e-f62ffc14c81f"},
					{Title: "THE HEART MECHANIC", MusicBrainzID: "45df8a99-4c58-40f8-83da-7f6b0e8f74a7"},
				},
			},
		},
		{
			ArtistName:    "Underoath",
			MusicBrainzID: "674a7e8c-9682-419a-8e05-2358e28b5359",
			Album: RawAlbum{
				Title:         "They're Only Chasing Safety",
				Type:          "ALBUM",
				MusicBrainzID: "4f313e9f-21e9-4e7e-a0f5-9f2c0248bedd",
				Tracks: []RawTracks{
					{Title: "Young And Aspiring", MusicBrainzID: "f8a0158f-8236-47b6-8258-131be47e811f"},
					{Title: "The Impact Of Reason", MusicBrainzID: "1a15c8fb-e248-426d-bb19-5e9f48476dfe"},
					{Title: "Some Will Seek Forgiveness, Others Escape", MusicBrainzID: "6be12749-02be-4ac3-97f0-afa0f14d9bf3"},
				},
			},
		},
		{
			ArtistName:    "Ptite Soeur",
			MusicBrainzID: "dac6f818-9b72-42ad-a2ff-1d6ee332ab4c",
			Album: RawAlbum{
				Title:         "KAYFABE CHIMERA",
				Type:          "EP",
				MusicBrainzID: "49aa6697-651f-491b-be4a-a274975607a0",
				Tracks: []RawTracks{
					{Title: "GEM KARSON", MusicBrainzID: "0caed678-3365-4cd9-bc91-f3b18ae7906d"},
					{Title: "KAYFABE", MusicBrainzID: "3f057c69-fdc2-4a91-b975-af38d7451367"},
				},
			},
		},
	}

	var tracks []models.Track

	for _, rawContent := range rawContents {
		// First create artist
		artist := &models.Artist{
			Name:          rawContent.ArtistName,
			MusicBrainzID: rawContent.MusicBrainzID,
		}

		err := db.Table("artists").Create(artist).Error
		require.NoError(t, err)

		// Then album
		album := &models.Album{
			Title:         rawContent.Album.Title,
			Artists:       []models.Artist{*artist},
			Type:          rawContent.Album.Type,
			MusicBrainzID: rawContent.Album.MusicBrainzID,
		}

		err = db.Table("albums").Create(album).Error
		require.NoError(t, err)

		// And finally the tracks
		for _, rawTrack := range rawContent.Album.Tracks {
			track := &models.Track{
				Artists:       []models.Artist{*artist},
				AlbumID:       album.ID,
				Title:         rawTrack.Title,
				MusicBrainzID: rawTrack.MusicBrainzID,
			}

			err = db.Table("tracks").Create(track).Error
			require.NoError(t, err)

			tracks = append(tracks, *track)
		}
	}
	return tracks
}
