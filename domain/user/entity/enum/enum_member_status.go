package enum

type MemberStatus string

const (
	Online  MemberStatus = "online"
	Offline MemberStatus = "offline"
	Running MemberStatus = "running" // 游戏中
)

var MemberStatusMap = map[MemberStatus]string{
	Online:  "online",
	Offline: "offline",
	Running: "running",
}

func (m MemberStatus) IsValid() bool {
	_, ok := MemberStatusMap[m]
	return ok
}
