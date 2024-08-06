package logics

import "forum/internal/tag/repositories"

func UpdateTagUserCountLogic(tagID uint) (string, error) {
	fansCount, err := repositories.UpdateTagUserCountReq(tagID)
	return fansCount, err
}
