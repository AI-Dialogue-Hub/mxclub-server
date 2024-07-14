package enum

type MemberStatus string

const (
	Online  MemberStatus = "online"
	Offline MemberStatus = "offline"
	Running MemberStatus = "running" // 游戏中
)
