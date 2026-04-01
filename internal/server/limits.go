package server

type Tier string

const (
	TierFree Tier = "free"
	TierPro  Tier = "pro"
)

type Limits struct {
	Tier        Tier
	Description string
}

func LimitsFor(tier string) Limits {
	if tier == "pro" {
		return Limits{Tier: TierPro, Description: "Unlimited links and clicks"}
	}
	return Limits{Tier: TierFree, Description: "10 links, 1k clicks tracked"}
}

func (l Limits) IsPro() bool {
	return l.Tier == TierPro
}
