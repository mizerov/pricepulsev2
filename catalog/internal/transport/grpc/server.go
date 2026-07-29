package grpc

import (
	"context"

	"github.com/google/uuid"
	catalogv1 "github.com/mizerov/pricepulse/api/catalog/v1"
	"github.com/mizerov/pricepulse/catalog/pkg/domains"
	"github.com/samber/do/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	catalogv1.UnimplementedCatalogServiceServer
	productRepo domains.ProductRepository
}

func NewServer(i do.Injector) (*Server, error) {
	productRepo, err := do.InvokeAs[domains.ProductRepository](i)
	if err != nil {
		return nil, err
	}
	return &Server{productRepo: productRepo}, nil
}

func (s *Server) CheckProductExists(ctx context.Context, req *catalogv1.CheckProductExistsRequest) (*catalogv1.CheckProductExistsResponse, error) {
	productID, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product_id: %v", err)
	}

	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		if err == domains.ErrProductNotFound {
			return &catalogv1.CheckProductExistsResponse{Exists: false}, nil
		}
		return nil, status.Errorf(codes.Internal, "failed to find product: %v", err)
	}

	if product.Status == domains.ProductStatusDraft {
		if req.AccountId != product.OwnerID.String() {
			return &catalogv1.CheckProductExistsResponse{Exists: false}, nil
		}
	}

	return &catalogv1.CheckProductExistsResponse{Exists: true}, nil
}
