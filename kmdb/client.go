package kmdb

import "git.apache.org/thrift.git/lib/go/thrift"

//   Client
// ----------

type Client interface {
	ThriftService
	Connect() (err error)
}

type client struct {
	addr string
	svc  *ThriftServiceClient
}

func NewClient(addr string) (c Client) {
	return &client{addr, nil}
}

func (c *client) Connect() (err error) {
	// pfac := thrift.NewTBinaryProtocolFactoryDefault()
	// pfac := thrift.NewTCompactProtocolFactory()
	pfac := thrift.NewTJSONProtocolFactory()

	trans, err := thrift.NewTSocket(c.addr)
	if err != nil {
		return err
	}

	if err := trans.Open(); err != nil {
		return err
	}

	c.svc = NewThriftServiceClientFactory(trans, pfac)

	return nil
}

func (c *client) Put(req *PutReq) (r *PutRes, err error) {
	return c.svc.Put(req)
}

func (c *client) Inc(req *IncReq) (r *IncRes, err error) {
	return c.svc.Inc(req)
}

func (c *client) Get(req *GetReq) (r *GetRes, err error) {
	return c.svc.Get(req)
}

func (c *client) PutBatch(batch []*PutReq) (r []*PutRes, err error) {
	return c.svc.PutBatch(batch)
}

func (c *client) IncBatch(batch []*IncReq) (r []*IncRes, err error) {
	return c.svc.IncBatch(batch)
}

func (c *client) GetBatch(batch []*GetReq) (r []*GetRes, err error) {
	return c.svc.GetBatch(batch)
}
