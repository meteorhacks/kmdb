package kmdb

import (
	"github.com/glycerine/go-capnproto"
	"github.com/meteorhacks/bddp"
)

//   Client
// ----------

type Client interface {
	Connect() (err error)
	PutBatch(sz int) (b *PutBatch, err error)
	GetBatch(sz int) (b *GetBatch, err error)
}

type client struct {
	bc bddp.Client
}

func NewClient(address string) (c Client) {
	bc := bddp.NewClient(address)
	return &client{bc}
}

func (c *client) Connect() (err error) {
	return c.bc.Connect()
}

func (c *client) PutBatch(sz int) (b *PutBatch, err error) {
	call, err := c.bc.Method("put")
	if err != nil {
		return nil, err
	}

	seg := call.Segment()
	reqs := NewPutRequestList(seg, sz)

	b = &PutBatch{
		reqs: reqs,
		call: call,
	}

	return b, nil
}

func (c *client) GetBatch(sz int) (b *GetBatch, err error) {
	call, err := c.bc.Method("get")
	if err != nil {
		return nil, err
	}

	seg := call.Segment()
	reqs := NewGetRequestList(seg, sz)

	b = &GetBatch{
		reqs: reqs,
		call: call,
	}

	return b, nil
}

//   PutBatch
// ------------

type PutBatch struct {
	reqs PutRequest_List
	call bddp.MCall
}

func (b *PutBatch) Set(db string, i int, ts int64, vals []string, pld []byte) (err error) {
	seg := b.call.Segment()
	req := NewPutRequest(seg)

	valsCount := len(vals)
	valsList := seg.NewTextList(valsCount)
	for j := 0; j < valsCount; j++ {
		valsList.Set(j, vals[j])
	}

	req.SetDb(db)
	req.SetValues(valsList)
	req.SetPayload(pld)
	req.SetTime(ts)
	b.reqs.Set(i, req)

	return nil
}

func (b *PutBatch) Send() (err error) {
	params := capn.Object(b.reqs)
	_, err = b.call.Call(params)
	return err
}

//   GetBatch
// ------------

type GetBatch struct {
	reqs GetRequest_List
	call bddp.MCall
}

func (b *GetBatch) Set(db string, i int, vals []string, start, end int64) (err error) {
	seg := b.call.Segment()
	req := NewGetRequest(seg)

	valsCount := len(vals)
	valsList := seg.NewTextList(valsCount)
	for i := 0; i < valsCount; i++ {
		valsList.Set(i, vals[i])
	}

	req.SetDb(db)
	req.SetValues(valsList)
	req.SetStart(start)
	req.SetEnd(end)
	b.reqs.Set(i, req)

	return nil
}

func (b *GetBatch) Send() (err error) {
	params := capn.Object(b.reqs)
	_, err = b.call.Call(params)
	return err
}
