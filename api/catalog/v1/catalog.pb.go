package catalogv1

import (
	"context"
	"google.golang.org/grpc"
)

type CheckProductExistsRequest struct {
	ProductId string
	AccountId string
}

type CheckProductExistsResponse struct {
	Exists bool
}

type CatalogServiceClient interface {
	CheckProductExists(ctx context.Context, in *CheckProductExistsRequest, opts ...grpc.CallOption) (*CheckProductExistsResponse, error)
}

type catalogServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCatalogServiceClient(cc grpc.ClientConnInterface) CatalogServiceClient {
	return &catalogServiceClient{cc}
}

func (c *catalogServiceClient) CheckProductExists(ctx context.Context, in *CheckProductExistsRequest, opts ...grpc.CallOption) (*CheckProductExistsResponse, error) {
	out := new(CheckProductExistsResponse)
	err := c.cc.Invoke(ctx, "/catalog.v1.CatalogService/CheckProductExists", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type CatalogServiceServer interface {
	CheckProductExists(context.Context, *CheckProductExistsRequest) (*CheckProductExistsResponse, error)
}

type UnimplementedCatalogServiceServer struct {
}

func (UnimplementedCatalogServiceServer) CheckProductExists(context.Context, *CheckProductExistsRequest) (*CheckProductExistsResponse, error) {
	return nil, nil
}

func RegisterCatalogServiceServer(s grpc.ServiceRegistrar, srv CatalogServiceServer) {
	s.RegisterService(&CatalogService_ServiceDesc, srv)
}

var CatalogService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "catalog.v1.CatalogService",
	HandlerType: (*CatalogServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckProductExists",
			Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
				in := new(CheckProductExistsRequest)
				if err := dec(in); err != nil {
					return nil, err
				}
				if interceptor == nil {
					return srv.(CatalogServiceServer).CheckProductExists(ctx, in)
				}
				info := &grpc.UnaryServerInfo{
					Server:     srv,
					FullMethod: "/catalog.v1.CatalogService/CheckProductExists",
				}
				handler := func(ctx context.Context, req interface{}) (interface{}, error) {
					return srv.(CatalogServiceServer).CheckProductExists(ctx, req.(*CheckProductExistsRequest))
				}
				return interceptor(ctx, in, info, handler)
			},
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "catalog.proto",
}
