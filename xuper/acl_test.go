package xuper

import "testing"

func TestAcl(t *testing.T) {
	acl := NewACL(1, 1.0)
	if acl.PM.AcceptValue != 1 || acl.PM.Rule != 1.0 {
		t.Error("Acl assert failed")
	}

	acl.AddAK("a", 1.0)
	if len(acl.AksWeight) != 1 {
		t.Error("Acl AddAK assert failed")
	}

	if v, ok := acl.AksWeight["a"]; !ok {
		t.Error("Acl AddAK assert failed")
	} else if v != 1.0 {
		t.Error("Acl AddAK assert failed")
	}

	defaultACL := getDefaultACL("bob")
	if v, ok := defaultACL.AksWeight["bob"]; !ok {
		t.Error("Acl AddAK assert failed")
	} else if v != 1.0 {
		t.Error("Acl AddAK assert failed")
	}
}
