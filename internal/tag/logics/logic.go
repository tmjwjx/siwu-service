package logics

import "forum/internal/tag/repositorys"

func UpdateTagUserCountLogic(tagID uint) (string, error) {
	fansCount, err := repositorys.UpdateTagUserCountReq(tagID)
	return fansCount, err
}
