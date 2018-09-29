package stocking

import (
	"bytes"
	"fmt"
)

// HubPackge TODO
type HubPackge struct {
	client  *Client
	raw     []byte
	mtype   TextMessageInboundType
	ack     string
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
	p.ack = string(segs[tmisAck])
	p.content = string(segs[tmisContent])

	return nil
}

func (p *HubPackge) encode() []byte {
	return []byte(fmt.Sprintf(`%v,%v,%v`, p.mtype, p.ack, p.content))
}

func (p *HubPackge) hasAck() bool {
	return p.ack != ""
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
