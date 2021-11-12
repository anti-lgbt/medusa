package models

import (
	"database/sql"
	"time"

	"github.com/anti-lgbt/medusa/config"
)

type Like struct {
	ID        int64         `gorm:"primaryKey"`
	UserID    int64         `gorm:"type:bigint;not null;uniqueIndex:idx_user_id_and_reply_id_and_comment_id_and_music_id_and_album_id"`
	ReplyID   sql.NullInt64 `gorm:"type:bigint;uniqueIndex:idx_user_id_and_reply_id_and_comment_id_and_music_id_and_album_id"`
	CommentID sql.NullInt64 `gorm:"type:bigint;uniqueIndex:idx_user_id_and_reply_id_and_comment_id_and_music_id_and_album_id"`
	MusicID   sql.NullInt64 `gorm:"type:bigint;uniqueIndex:idx_user_id_and_reply_id_and_comment_id_and_music_id_and_album_id"`
	AlbumID   sql.NullInt64 `gorm:"type:bigint;uniqueIndex:idx_user_id_and_reply_id_and_comment_id_and_music_id_and_album_id"`
	CreatedAt time.Time     `gorm:"type:timestamp(0);not null;index"`
	UpdatedAt time.Time     `gorm:"type:timestamp(0);not null;index"`
	User      *User
	Reply     *Reply
	Comment   *Comment
	Music     *Music
	Album     *Album
}

func (l *Like) UnLike() {
	config.Database.Delete(&l)
}
