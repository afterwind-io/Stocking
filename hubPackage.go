package stocking

import (
	"bytes"
)

// HubPackge TODO
type HubPackge struct {
	client  *Client
	raw     []byte
	mtype   string
	ccode   string
	content string
}

func (p *HubPackge) decode() error {
	segs := bytes.SplitN(p.raw, []byte(","), 3)

	if len(segs) != 3 {
		return sError{"Invalid message format."}
	}

	mtype := segs[tmisType]
	if len(mtype) == 0 {
		return sError{"Invalid message type."}
	}

	p.mtype = string(mtype)
	p.ccode = string(segs[tmisCCode])
	p.content = string(segs[tmisContent])

	return nil
}

func (p *HubPackge) encode() []byte {
	return []byte(p.mtype + "," + p.ccode + "," + p.content)
}

func (p *HubPackge) hasAck() bool {
	return p.ccode != "0"
}

func (p *HubPackge) error(ccode, content string) {
	p.mtype = tmitError
	p.ccode = ccode
	p.content = content
}

func newHubPackage(c *Client, m []byte) (*HubPackge, error) {
	p := HubPackge{
		client: c,
		raw:    m,
	}

	if err := p.decode(); err != nil {
		return nil, err
	}

	return &p, nil
}
