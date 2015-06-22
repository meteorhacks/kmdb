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

func (b *PutBatch) Set(i int, db string, ts int64, fields []string, val float64, num int64) (err error) {
	seg := b.call.Segment()
	req := NewPutRequest(seg)
	req.SetDatabase(db)
	req.SetTimestamp(ts)
	req.SetValue(val)
	req.SetCount(num)
	req.SetFields(toTextList(seg, fields))
	b.reqs.Set(i, req)

	return nil
}

func (b *PutBatch) Send() (res capn.Object, err error) {
	params := capn.Object(b.reqs)
	return b.call.Call(params)
}

//   GetBatch
// ------------

type GetBatch struct {
	reqs GetRequest_List
	call bddp.MCall
}

func (b *GetBatch) Set(i int, db string, start, end int64, fields []string, groupBy []bool) (err error) {
	seg := b.call.Segment()
	req := NewGetRequest(seg)
	req.SetDatabase(db)
	req.SetStartTime(start)
	req.SetEndTime(end)
	req.SetFields(toTextList(seg, fields))
	req.SetGroupBy(toBitList(seg, groupBy))
	b.reqs.Set(i, req)

	return nil
}

func (b *GetBatch) Send() (res capn.Object, err error) {
	params := capn.Object(b.reqs)
	return b.call.Call(params)
}

//   Cap'n Proto
// ---------------

func toTextList(seg *capn.Segment, vals []string) (list capn.TextList) {
	count := len(vals)
	list = seg.NewTextList(count)

	for i := 0; i < count; i++ {
		list.Set(i, vals[i])
	}

	return list
}

func toBitList(seg *capn.Segment, vals []bool) (list capn.BitList) {
	count := len(vals)
	list = seg.NewBitList(count)

	for i := 0; i < count; i++ {
		list.Set(i, vals[i])
	}

	return list
}
