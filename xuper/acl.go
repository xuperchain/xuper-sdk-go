package xuper

// ACL acl.
type ACL struct {
	PM        PermissionModel    `json:"pm"`
	AksWeight map[string]float64 `json:"aksWeight"`
}

// PermissionModel acl permission model.
type PermissionModel struct {
	Rule        int32   `json:"rule"`
	AcceptValue float64 `json:"acceptValue"`
}

// NewACL new ACl instance.
func NewACL(rule int32, acceptValue float64) *ACL {
	return &ACL{
		PM: PermissionModel{
			Rule:        rule,
			AcceptValue: acceptValue,
		},
		AksWeight: map[string]float64{},
	}
}

// AddAK add ak and weight pair.
func (a *ACL) AddAK(ak string, weight float64) {
	if a.AksWeight == nil {
		a.AksWeight = make(map[string]float64, 1)
	}
	a.AksWeight[ak] = weight
}

func getDefaultACL(address string) *ACL {
	return &ACL{
		PM: PermissionModel{
			Rule:        1,
			AcceptValue: 1.0,
		},
		AksWeight: map[string]float64{
			address: 1.0,
		},
	}
}
