package login

import "qserver/db"

//1.查询玩家是否已经注册
func getPlayerIdByAccname(accname string) (int, error) {
	var id int
	row := db.DB.QueryRow("SELECT `id` FROM `player` WHERE `accname`=?", accname)
	//if err != nil {
	//	return 0, err
	//}
	row.Scan(&id)

	return id, nil
}

//注册玩家
func registerUser(accname string) (int, error) {
	result, err := db.DB.Exec("INSERT INTO `player`(`accname`) VALUES(?)", accname)
	if err != nil {
		return 0, err
	}
	insId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(insId), nil
}
