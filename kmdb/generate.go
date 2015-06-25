//go:generate thrift --gen go -out ../ protocol.thrift
//go:generate rm -rf thrift_service-remote

package kmdb
