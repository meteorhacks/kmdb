package kmdb

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

//   Client
// ----------

type Client interface {
	Connect() (err error)
	Put(req *PutReq) (r *PutRes, err error)
	Inc(req *IncReq) (r *IncRes, err error)
	Get(req *GetReq) (r *GetRes, err error)
	PutBatch(batch *PutReqBatch) (r *PutResBatch, err error)
	IncBatch(batch *IncReqBatch) (r *IncResBatch, err error)
	GetBatch(batch *GetReqBatch) (r *GetResBatch, err error)
}

type client struct {
	addr string
	con  *grpc.ClientConn
	svc  DatabaseServiceClient
}

func NewClient(addr string) (c Client) {
	return &client{addr, nil, nil}
}

func (c *client) Connect() (err error) {
	if c.con != nil {
		err = c.con.Close()
		if err != nil {
			return err
		}
	}

	c.con, err = grpc.Dial(c.addr)
	if err != nil {
		return err
	}

	c.svc = NewDatabaseServiceClient(c.con)

	return nil
}

func (c *client) Put(req *PutReq) (r *PutRes, err error) {
	ctx := context.TODO()
	return c.svc.Put(ctx, req)
}

func (c *client) Inc(req *IncReq) (r *IncRes, err error) {
	ctx := context.TODO()
	return c.svc.Inc(ctx, req)
}

func (c *client) Get(req *GetReq) (r *GetRes, err error) {
	ctx := context.TODO()
	return c.svc.Get(ctx, req)
}

func (c *client) PutBatch(batch *PutReqBatch) (r *PutResBatch, err error) {
	ctx := context.TODO()
	return c.svc.PutBatch(ctx, batch)
}

func (c *client) IncBatch(batch *IncReqBatch) (r *IncResBatch, err error) {
	ctx := context.TODO()
	return c.svc.IncBatch(ctx, batch)
}

func (c *client) GetBatch(batch *GetReqBatch) (r *GetResBatch, err error) {
	ctx := context.TODO()
	return c.svc.GetBatch(ctx, batch)
}
