package enum

type DasherLevel int

const (
	_ DasherLevel = iota
	DasherLevel_Gold
	DasherLevel_Silver // 银牌
	DasherLevel_Bronze // 铜牌
)
