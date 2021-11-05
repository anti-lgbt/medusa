package models

import (
	"time"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/models/datatypes"
)

type Like struct {
	ID        int64               `gorm:"primaryKey"`
	UserID    int64               `gorm:"type:bigint;not null;uniqueIndex:idx_user_id_and_reply_id_and_comment_id_and_music_id_and_album_id"`
	ReplyID   datatypes.NullInt64 `gorm:"type:bigint;uniqueIndex:idx_user_id_and_reply_id_and_comment_id_and_music_id_and_album_id"`
	CommentID datatypes.NullInt64 `gorm:"type:bigint;uniqueIndex:idx_user_id_and_reply_id_and_comment_id_and_music_id_and_album_id"`
	MusicID   datatypes.NullInt64 `gorm:"type:bigint;uniqueIndex:idx_user_id_and_reply_id_and_comment_id_and_music_id_and_album_id"`
	AlbumID   datatypes.NullInt64 `gorm:"type:bigint;uniqueIndex:idx_user_id_and_reply_id_and_comment_id_and_music_id_and_album_id"`
	CreatedAt time.Time           `gorm:"type:timestamp(0);not null;index"`
	UpdatedAt time.Time           `gorm:"type:timestamp(0);not null;index"`
	User      *User
	Reply     *Reply
	Comment   *Comment
	Music     *Music
	Album     *Album
}

func (l *Like) UnLike() {
	config.Database.Delete(&l)
}
