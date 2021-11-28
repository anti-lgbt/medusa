package models

import (
	"database/sql"
	"time"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/controllers/entities"
	"github.com/volatiletech/null"
)

type Album struct {
	ID          int64          `json:"id" gorm:"primaryKey"`
	UserID      int64          `json:"user_id" gorm:"type:bigint;not null;index"`
	Name        string         `json:"name" gorm:"type:character varying;not null;index"`
	Description sql.NullString `json:"description" gorm:"type:character varying;index"`
	ViewCount   int64          `json:"view_count" gorm:"type:integer;not null;index;default:0"`
	Image       sql.NullString `json:"-" gorm:"type:character varying"`
	CreatedAt   time.Time      `json:"created_at" gorm:"type:timestamp(0);not null;index"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"type:timestamp(0);not null;index"`
	MusicAlbums []*MusicAlbum  `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	Likes       []*Like        `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	Comments    []*Comment     `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	User        *User
}

func (a *Album) Musics() []*Music {
	var musics []*Music

	config.Database.Find(&musics, "id IN (SELECT music_id FROM \"music_albums\" WHERE \"album_id\" = ?)", a.ID)

	return musics
}

func (a *Album) Comment(user_id int64, content string) *Comment {
	comment := &Comment{
		UserID: user_id,
		AlbumID: sql.NullInt64{
			Int64: a.ID,
			Valid: true,
		},
		Content: content,
	}

	config.Database.Create(&comment)

	return comment
}

func (a *Album) Delete() {
	config.Database.Delete(&a)
}

func (a *Album) Liked(user_id int64) bool {
	var l *Like

	result := config.Database.First(&l, "user_id = ? AND album_id = ?", user_id, a.ID)

	return result.Error == nil
}

func (a *Album) Like(user_id int64) error {
	var l *Like

	if result := config.Database.First(&l, "user_id = ? AND album_id = ?", user_id, a.ID); result.Error == nil {
		return result.Error
	}

	like := Like{
		UserID: user_id,
		AlbumID: sql.NullInt64{
			Int64: a.ID,
			Valid: true,
		},
	}

	config.Database.Create(&like)
	return nil
}

func (a *Album) UnLike(user_id int64) {
	config.Database.Where("user_id = ? AND album_id = ?", user_id, a.ID).Delete(&Like{})
}

func (a *Album) LikeCount() int64 {
	var count int64

	config.Database.Model(&Like{}).Where("album_id = ?", a.ID).Count(&count)

	return count
}

func (a *Album) ToEntity(user_id sql.NullInt64) *entities.Album {
	musics := a.Musics()
	music_entities := make([]*entities.Music, 0)

	for _, music := range musics {
		music_entities = append(music_entities, music.ToEntity(user_id))
	}

	entity := &entities.Album{
		ID:   a.ID,
		Name: a.Name,
		Description: null.String{
			String: a.Description.String,
			Valid:  a.Description.Valid,
		},
		ViewCount: a.ViewCount,
		LikeCount: a.LikeCount(),
		Musics:    music_entities,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}

	if user_id.Valid {
		entity.Liked = a.Liked(user_id.Int64)
	}

	return entity
}
