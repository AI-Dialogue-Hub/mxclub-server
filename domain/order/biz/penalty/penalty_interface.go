package penalty

// IPenaltyRule 定义PenaltyRule接口
type IPenaltyRule interface {
	ApplyPenalty(req *PenaltyReq) (*PenaltyResp, error)
	MustApplyPenalty(req *PenaltyReq) *PenaltyResp // if happen error, will give default resp
}
