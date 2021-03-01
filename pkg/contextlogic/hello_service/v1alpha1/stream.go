package v1alpha1

import (
	"io"

	"github.com/ContextLogic/hello-service/api/proto_gen/contextlogic/hello_service/v1alpha1"
)

//Stream is Stream
type Stream struct {
}

//BiDirectStream is the code implementation of the function defined in the proto files
func (s *Stream) BiDirectStream(stream v1alpha1.Stream_BiDirectStreamServer) (err error) {

	var streamDataList = []*v1alpha1.StreamData{}
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		streamDataList = append(streamDataList, in.StreamData)
	}
	for _, sd := range streamDataList {
		res := &v1alpha1.StreamDataResponse{
			StreamData: &v1alpha1.StreamData{
				Id:    sd.Id,
				Value: float64(sd.Id),
			},
		}
		if err := stream.Send(res); err != nil {
			return err
		}
	}
	return stream.Send(&v1alpha1.StreamDataResponse{})
}

//ClientStream is the code implementation of the function defined in the proto files
func (s *Stream) ClientStream(stream v1alpha1.Stream_ClientStreamServer) (err error) {

	var total float64

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		total += in.StreamData.Value
	}

	return stream.SendAndClose(&v1alpha1.StreamDataResponse{
		StreamData: &v1alpha1.StreamData{
			Value: total,
		},
	})
}

//ServerStream is the code implementation of the function defined in the proto files
func (s *Stream) ServerStream(req *v1alpha1.StreamDataRequest, stream v1alpha1.Stream_ServerStreamServer) (err error) {

	for i := 0; i < 10; i++ {
		res := &v1alpha1.StreamDataResponse{
			StreamData: &v1alpha1.StreamData{
				Id:    req.StreamData.Id,
				Value: float64(i),
			},
		}
		if err := stream.Send(res); err != nil {
			return err
		}
	}
	return stream.Send(&v1alpha1.StreamDataResponse{})
}
