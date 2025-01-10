package stat

import (
	"GoAdvanced/pkg/db"
	"gorm.io/datatypes"
	"time"
)

type StatRepository struct {
	Db *db.Db
}

func NewStatRepository(db *db.Db) *StatRepository {
	return &StatRepository{
		Db: db,
	}
}

func (repo *StatRepository) AddClick(linkId uint) {
	var stat Stat
	currentTime := datatypes.Date(time.Now())
	repo.Db.Find(&stat, "link_id = ? and date = ?", linkId, currentTime)
	if stat.ID == 0 {
		repo.Db.Create(&Stat{
			LinkId: linkId,
			Clicks: 1,
			Date:   currentTime,
		})
	} else {
		stat.Clicks++
		repo.Db.Save(&stat)
	}
}

func (repo *StatRepository) GetStat(by string, from, to time.Time) []GetStatResponce {
	var StatResponce []GetStatResponce
	var QuerySelect string
	switch by {
	case GroupByDay:
		QuerySelect = "to_char(date, 'YYYY-MM-DD') as period, sum(clicks)"
	case GroupByMonth:
		QuerySelect = "to_char(date, 'YYYY-MM') as period, sum(clicks)"
	}
	repo.Db.Table("stats").
		Select(QuerySelect).
		Where("date BETWEEN ? AND ?", from, to).
		Group("period").
		Order("period").
		Scan(&StatResponce)
	return StatResponce
}
