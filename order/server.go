package order

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/FrankOHara43/go-grpc-microservice/account"
	"github.com/FrankOHara43/go-grpc-microservice/catalog"
	"github.com/FrankOHara43/go-grpc-microservice/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedOrderServiceServer
	service       Service
	accountClient *account.Client
	catalogClient *catalog.Client
}

func ListenGRPC(s Service, accountURL, catalogURL string, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		return err
	}

	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		accountClient.Close()
		return err
	}

	serv := grpc.NewServer()
	pb.RegisterOrderServiceServer(serv, &grpcServer{
		UnimplementedOrderServiceServer: pb.UnimplementedOrderServiceServer{},
		service:                         s,
		accountClient:                   accountClient,
		catalogClient:                   catalogClient,
	})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) PostOrder(ctx context.Context, r *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	if _, err := s.accountClient.GetAccount(ctx, r.AccountId); err != nil {
		log.Println("error getting account:", err)
		return nil, errors.New("account not found")
	}

	productIDs := make([]string, 0, len(r.Products))
	requestedQty := make(map[string]uint32, len(r.Products))
	for _, p := range r.Products {
		productIDs = append(productIDs, p.ProductId)
		requestedQty[p.ProductId] = p.Quantity
	}

	catalogProducts, err := s.catalogClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		log.Println("error getting products:", err)
		return nil, errors.New("products not found")
	}

	products := make([]OrderedProduct, 0, len(catalogProducts))
	for _, p := range catalogProducts {
		qty := requestedQty[p.ID]
		if qty == 0 {
			continue
		}
		products = append(products, OrderedProduct{
			ID:          p.ID,
			Quantity:    qty,
			Price:       p.Price,
			Name:        p.Name,
			Description: p.Description,
		})
	}

	newOrder, err := s.service.PostOrder(ctx, r.AccountId, products)
	if err != nil {
		log.Println("error posting order:", err)
		return nil, errors.New("error posting order")
	}

	orderProto := &pb.Order{
		Id:         newOrder.ID,
		AccountId:  newOrder.AccountID,
		TotalPrice: newOrder.TotalPrice,
		Products:   []*pb.Order_OrderProduct{},
	}
	orderProto.CreatedAt, _ = newOrder.CreatedAt.MarshalBinary()

	for _, p := range newOrder.Products {
		orderProto.Products = append(orderProto.Products, &pb.Order_OrderProduct{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
		})
	}

	return &pb.PostOrderResponse{Order: orderProto}, nil
}

func (s *grpcServer) GetOrdersForAccount(ctx context.Context, r *pb.GetOrdersForAccountRequest) (*pb.GetOrdersForAccountResponse, error) {
	accountOrders, err := s.service.GetOrdersForAccount(ctx, r.AccountId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	productIDMap := map[string]bool{}
	for _, o := range accountOrders {
		for _, p := range o.Products {
			productIDMap[p.ID] = true
		}
	}

	productIDs := []string{}
	for id := range productIDMap {
		productIDs = append(productIDs, id)
	}

	catalogProducts, err := s.catalogClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		log.Println("error getting products:", err)
		return nil, errors.New("products not found")
	}

	catalogByID := make(map[string]catalog.Product, len(catalogProducts))
	for _, p := range catalogProducts {
		catalogByID[p.ID] = p
	}

	orders := []*pb.Order{}
	for _, o := range accountOrders {
		op := &pb.Order{
			AccountId:  o.AccountID,
			Id:         o.ID,
			TotalPrice: o.TotalPrice,
			Products:   []*pb.Order_OrderProduct{},
		}
		op.CreatedAt, _ = o.CreatedAt.MarshalBinary()

		for _, product := range o.Products {
			if p, ok := catalogByID[product.ID]; ok {
				product.Name = p.Name
				product.Description = p.Description
				product.Price = p.Price
			}
			op.Products = append(op.Products, &pb.Order_OrderProduct{
				Id:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				Quantity:    product.Quantity,
			})
		}

		orders = append(orders, op)
	}

	return &pb.GetOrdersForAccountResponse{Orders: orders}, nil
}
