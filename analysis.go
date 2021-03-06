package lookout

import (
	"github.com/src-d/lookout/pb"

	"google.golang.org/grpc"
)

type EventResponse = pb.EventResponse
type Comment = pb.Comment

type AnalyzerClient = pb.AnalyzerClient
type AnalyzerServer = pb.AnalyzerServer

func RegisterAnalyzerServer(s *grpc.Server, srv AnalyzerServer) {
	pb.RegisterAnalyzerServer(s, srv)
}

func NewAnalyzerClient(conn *grpc.ClientConn) AnalyzerClient {
	return pb.NewAnalyzerClient(conn)
}
