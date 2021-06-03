package xuper

type Proposal struct {
	xclient *XClient
	request *Request
}

func (p *Proposal) Build() (*Transaction, error) {

	return nil, nil
}
