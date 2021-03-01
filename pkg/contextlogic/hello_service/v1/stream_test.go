package v1

import (
	"context"
	"io"
	"testing"

	"github.com/ContextLogic/hello-service/api/proto_gen/contextlogic/hello_service/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func testStreamService() (*Stream, error) {
	return &Stream{}, nil
}

type testClient struct {
	grpc.ServerStream
	c chan *v1.StreamDataResponse
}

func (tc *testClient) Send(m *v1.StreamDataResponse) error {
	tc.c <- m
	return nil
}

func (tc *testClient) SendAndClose(m *v1.StreamDataResponse) error {
	tc.c <- m
	tc.c <- nil
	return nil
}

func (tc *testClient) Recv() (*v1.StreamDataRequest, error) {

	r := <-tc.c
	if r == nil || r.StreamData == nil {
		return nil, io.EOF
	}
	return &v1.StreamDataRequest{StreamData: r.StreamData}, nil

}

type testServerStream struct {
	ctx context.Context
}

func (ts *testServerStream) SendMsg(m interface{}) error  { return nil }
func (ts *testServerStream) RecvMsg(m interface{}) error  { return nil }
func (ts *testServerStream) SetHeader(metadata.MD) error  { return nil }
func (ts *testServerStream) SendHeader(metadata.MD) error { return nil }
func (ts *testServerStream) SetTrailer(metadata.MD)       {}
func (ts *testServerStream) Context() context.Context     { return ts.ctx }

func TestBiDirectStream(t *testing.T) {

	s, err := testStreamService()
	assert.Emptyf(t, err, "Failed creating HelloStreamService: %v", err)

	ts := testServerStream{context.Background()}
	tc := testClient{&ts, make(chan *v1.StreamDataResponse)}

	// go routine the BiDirectStream
	go s.BiDirectStream(&tc)

	// stream the testing data
	sd1 := v1.StreamDataResponse{
		StreamData: &v1.StreamData{
			Id:    1,
			Value: 1,
		},
	}
	sd2 := v1.StreamDataResponse{
		StreamData: &v1.StreamData{
			Id:    2,
			Value: 2,
		},
	}
	tc.Send(&sd1)
	tc.Send(&sd2)
	tc.Send(nil)

	// verify the returned stream
	in, err := tc.Recv()
	assert.Equal(t, in.StreamData.Id, sd1.StreamData.Id, "BiDirectStream returned incorrect stream")
	assert.Equal(t, in.StreamData.Value, sd1.StreamData.Value, "BiDirectStream returned incorrect stream")
	assert.Emptyf(t, err, "Recv returned non nil error: %v", err)

	in, err = tc.Recv()
	assert.Equal(t, in.StreamData.Id, sd2.StreamData.Id, "BiDirectStream returned incorrect stream")
	assert.Equal(t, in.StreamData.Value, sd2.StreamData.Value, "BiDirectStream returned incorrect stream")
	assert.Emptyf(t, err, "Recv returned non nil error: %v", err)

	in, err = tc.Recv()
	assert.Emptyf(t, in, "returned stream did not close properly")
	assert.Equal(t, io.EOF, err, "Recv returned non EOF error: %v", err)

}

func TestClientStream(t *testing.T) {
	s, err := testStreamService()
	assert.Emptyf(t, err, "Failed creating HelloStreamService: %v", err)

	ts := testServerStream{context.Background()}
	tc := testClient{&ts, make(chan *v1.StreamDataResponse)}

	go s.ClientStream(&tc)

	for i := 1; i < 11; i++ {
		tc.Send(&v1.StreamDataResponse{
			StreamData: &v1.StreamData{
				Id:    int64(i),
				Value: float64(i),
			},
		})
	}
	tc.Send(nil)

	in, err := tc.Recv()
	assert.Equal(t, in.StreamData.Value, float64(55), "ClientStream returned incorrect value")
	assert.Emptyf(t, err, "Recv returned non nil error: %v", err)

	in, err = tc.Recv()
	assert.Emptyf(t, in, "returned stream did not close properly")
	assert.Equal(t, io.EOF, err, "Recv returned non EOF error: %v", err)

}

func TestServerStream(t *testing.T) {
	s, err := testStreamService()
	assert.Emptyf(t, err, "Failed creating HelloStreamService: %v", err)

	ts := testServerStream{context.Background()}
	tc := testClient{&ts, make(chan *v1.StreamDataResponse)}

	go s.ServerStream(&v1.StreamDataRequest{
		StreamData: &v1.StreamData{
			Id:    1,
			Value: 1,
		},
	}, &tc)

	for i := 0; i < 10; i++ {
		in, err := tc.Recv()
		assert.Emptyf(t, err, "Recv returned non nil error: %v", err)
		assert.Equal(t, in.StreamData.Value, float64(i), "SeverStream returned incorrect value")
	}
	in, err := tc.Recv()
	assert.Emptyf(t, in, "returned stream did not close properly")
	assert.Equal(t, io.EOF, err, "Recv returned non EOF error: %v", err)
}
