package public

import (
	"time"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/models"
	"github.com/gofiber/fiber/v2"
)

type TrendingMusic struct {
	ID             int64 `json:"id"`
	TodayViewCount int64 `json:"today_view_count"`
	TotalViewCount int64 `json:"total_view_count"`
}

// GET /api/v2/public/trending/musics/daily
func GetTrendingMusic(c *fiber.Ctx) error {
	var trending_musics []*TrendingMusic

	config.Database.
		Select("musics.id as id, (trending_musics.total_view_count - musics.view_count) AS today_view_count, musics.view_count AS total_view_count").
		Joins("JOIN trending_musics ON trending_musics.music_id = musics.id AND DATE(trending_musics.release_at) = DATE(?) AND trending_musics.total_view_count > musics.view_count", time.Now().AddDate(0, 0, -1)).
		Limit(10).
		Find(&trending_musics)

	var trending_musics_was_created_today []*models.Music

	config.Database.
		Where("DATE(created_at) = DATE(?)", time.Now()).
		Order("view_count desc").
		Limit(10).
		Find(&trending_musics_was_created_today)

	if len(trending_musics_was_created_today) == 0 {
		return nil
	}

	for _, m := range trending_musics_was_created_today {
		added_to_trending_list := false

		for i, trending_music := range trending_musics {
			if m.ViewCount > trending_music.TodayViewCount {
				trending_musics = append(trending_musics[:i+1], trending_musics[i:]...)
				trending_musics[i] = &TrendingMusic{
					ID:             m.ID,
					TodayViewCount: m.ViewCount,
					TotalViewCount: m.ViewCount,
				}
				added_to_trending_list = true
				break
			}
		}

		if !added_to_trending_list {
			trending_musics[len(trending_musics)+1] = &TrendingMusic{
				ID:             m.ID,
				TodayViewCount: m.ViewCount,
				TotalViewCount: m.ViewCount,
			}
		}
	}

	i := 0
	trendings := make([]*TrendingMusic, 0)
	for _, t_m := range trending_musics {
		if i == 10 {
			break
		}

		trendings = append(trendings, t_m)
		i++
	}

	return c.Status(200).JSON(trendings)
}
