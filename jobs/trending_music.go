package jobs

import (
	"time"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/models"
	"github.com/jasonlvhit/gocron"
)

type TrendingMusicJob struct{}

func NewTrendingMusicJob() *TrendingMusicJob {
	return &TrendingMusicJob{}
}

func (job *TrendingMusicJob) Process() {
	s := gocron.NewScheduler()
	s.Every(1).Day().At("00:00:00").Do(func() {
		var musics []*models.Music

		config.Database.Find(&musics, "DATE(updated_at) >= DATE(?)", time.Now().AddDate(0, 0, -1))

		for _, music := range musics {
			var yesterday_trending_music *models.TrendingMusic

			if result := config.Database.First(&yesterday_trending_music, "release_at = DATE(?)", time.Now().AddDate(0, 0, -2)); result.Error != nil {
				today_trending_music := &models.TrendingMusic{
					MusicID:        music.ID,
					TotalViewCount: music.ViewCount,
					DayViewCount:   music.ViewCount,
					ReleaseAt:      time.Now(),
				}

				config.Database.Create(&today_trending_music)
			} else {
				if yesterday_trending_music.TotalViewCount == music.ViewCount {
					return
				}

				today_trending_music := &models.TrendingMusic{
					MusicID:        music.ID,
					TotalViewCount: music.ViewCount,
					DayViewCount:   music.ViewCount - yesterday_trending_music.TotalViewCount,
					ReleaseAt:      time.Now(),
				}

				config.Database.Create(&today_trending_music)
			}
		}
	})
	<-s.Start()

	for {
	}
}
