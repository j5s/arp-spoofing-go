package models

//Host 主机信息
type Host struct {
	ID      int    `db:"id"`
	IP      string `db:"ip"`
	MAC     string `db:"mac"`
	MACInfo string `db:"mac_info"`
}
