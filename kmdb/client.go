package kmdb

import (
	"github.com/golang/protobuf/proto"
	"github.com/meteorhacks/simple-rpc-go"
)

//   Client
// ----------

type Client interface {
	Connect() (err error)
	Put(req *PutReqBatch) (r *PutResBatch, err error)
	Inc(req *IncReqBatch) (r *IncResBatch, err error)
	Get(req *GetReqBatch) (r *GetResBatch, err error)
}

type client struct {
	addr string
	cli  srpc.Client
}

func NewClient(addr string) (c Client) {
	cli := srpc.NewClient(addr)
	return &client{addr, cli}
}

func (c *client) Connect() (err error) {
	return c.cli.Connect()
}

func (c *client) Put(req *PutReqBatch) (r *PutResBatch, err error) {
	pld, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	out, err := c.cli.Call("put", pld)
	if err != nil {
		return nil, err
	}

	r = &PutResBatch{}
	err = proto.Unmarshal(out, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (c *client) Inc(req *IncReqBatch) (r *IncResBatch, err error) {
	pld, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	out, err := c.cli.Call("inc", pld)
	if err != nil {
		return nil, err
	}

	r = &IncResBatch{}
	err = proto.Unmarshal(out, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (c *client) Get(req *GetReqBatch) (r *GetResBatch, err error) {
	pld, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	out, err := c.cli.Call("get", pld)
	if err != nil {
		return nil, err
	}

	r = &GetResBatch{}
	err = proto.Unmarshal(out, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}
