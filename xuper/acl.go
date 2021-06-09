package xuper

type Acl struct {
	Pm        *Pm
	AksWeight map[string]int
}

type Pm struct {
	Rule        int
	AcceptValue int
}

func (a *Acl) AddAK(ak string, weight int) {
	a.AksWeight[ak] = weight
}

func getDefaultACL(address string) string {
	return `
	{
		"pm": {
			"rule": 1,
			"acceptValue": 1.0
		},
		"aksWeight": {
			"` + address + `": 1.0
		}
	}
	`
}
