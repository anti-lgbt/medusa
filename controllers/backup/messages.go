package controllers

// UID: UID1564133 | GID: GID156615616
// func SendMessage(uid, content string) {
// 	var CurrentUserID int64 = 1
// 	var group *models.Group

// 	r := regexp.MustCompile("UID([^;]*).*")
// 	if result := config.Database.Joins("UserGroup").Where("user_groups.").First(&group, "uid = ?", uid); result.Error != nil && r.MatchString(uid) {
// 		var target_user *models.User

// 		if result := config.Database.First(&target_user, "uid = ?", uid); result.Error != nil {
// 			// TODO: ERROR
// 		}

// 		group := models.Group{
// 			UID: ,
// 		}

// 		config.Database.Create()

// 		user_groups = []*models.UserGroup{
// 			&models.UserGroup{
// 				UserID: CurrentUserID,
// 				GroupID: ,
// 			},
// 		}

// 		config.Database.Create()
// 	}
// }
