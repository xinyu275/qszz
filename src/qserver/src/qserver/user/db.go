package user

import "qserver/db"

func getPlayerInfoFromDb(playerId int) (*Player, error) {
	var player Player
	row := db.DB.QueryRow("SELECT `id`,`accname`, `nickname`, `lv`, `exp`, "+
		"`reg_time`, `last_login_time`, `last_offline_time`, `last_login_ip`, `gold` "+
		"FROM `player` WHERE `id`=?", playerId)
	err := row.Scan(
		&player.PlayerId,
		&player.Account,
		&player.Nickname,
		&player.Lv,
		&player.Exp,
		&player.RegTime,
		&player.LastLoginTime,
		&player.LastOfflineTime,
		&player.LastLoginIp,
		&player.Gold)
	if err != nil {
		return nil, err
	}
	return &player, nil
}
