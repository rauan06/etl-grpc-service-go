package main

func main() {
	// ctx := context.Background()
	// cfg := config.LoadConfig()
	// logger := logger.SetupPrettySlog(os.Stdout)

	// // Start gRPC server
	// go func() {
	// 	lis, err := net.Listen("tcp", ":5059")
	// 	if err != nil {
	// 		log.Fatalf("failed to listen: %v", err)
	// 	}

	// 	grpcServer := grpc.NewServer()
	// 	grpcCategoryClient, err := client.NewCategoryClient(ctx, "0.0.0.0:5059")
	// 	if err != nil {
	// 		log.Fatalf("Failed to create category client: %v", err)
	// 	}
	// 	grpcCityClient, err := client.NewCityClient(ctx, "0.0.0.0:5059")
	// 	if err != nil {
	// 		log.Fatalf("Failed to create city client: %v", err)
	// 	}
	// 	grpcPriceClient, err := client.NewPriceClient(ctx, "0.0.0.0:5059")
	// 	if err != nil {
	// 		log.Fatalf("Failed to create price client: %v", err)
	// 	}
	// 	grpcStockClient, err := client.NewStockClient(ctx, "0.0.0.0:5059")
	// 	if err != nil {
	// 		log.Fatalf("Failed to create stock client: %v", err)
	// 	}
	// 	grpcProductClient, err := client.NewProductClient(ctx, "0.0.0.0:5059")
	// 	if err != nil {
	// 		log.Fatalf("Failed to create product client: %v", err)
	// 	}

	// 	URL, _ := url.Parse("http://0.0.0.0:8080")

	// 	categoryHttpClient := clientHttp.NewCategoryClient(URL)
	// 	cityHttpClient := clientHttp.NewCityClient(URL)
	// 	priceHttpClient := clientHttp.NewPriceClient(URL)
	// 	stockHttpClient := clientHttp.NewStockClient(URL)
	// 	productHttpClient := clientHttp.NewProductClient(URL)

	// 	cache, err := redis.New(ctx, cfg)
	// 	if err != nil {
	// 		log.Fatalf("Failed to initialize Redis: %v", err)
	// 	}

	// 	// Construct service layers from clients
	// 	categorySvc := service.NewCategoryService(grpcCategoryClient, categoryHttpClient, cache, logger)
	// 	citySvc := service.NewCityService(grpcCityClient, cityHttpClient, cache, logger)
	// 	priceSvc := service.NewPriceService(grpcPriceClient, priceHttpClient, cache, logger)
	// 	stockSvc := service.NewStockService(grpcStockClient, stockHttpClient, cache, logger)
	// 	productSvc := service.NewProductService(grpcProductClient, productHttpClient, cache, logger)

	// 	// Create ETL handler
	// 	h := handler.NewEtlHandler(categorySvc, citySvc, priceSvc, stockSvc, productSvc)

	// 	// Register ETL handler to the gRPC server (only ONCE, and pass the handler itself)
	// 	mux := runtime.NewServeMux()
	// 	etlPb.RegisterCategoryServiceServer(grpcServer, h)
	// 	etlPb.RegisterCategoryServiceHandlerFromEndpoint(ctx, mux)

	// 	log.Println("gRPC server listening on :5059")
	// 	if err := grpcServer.Serve(lis); err != nil {
	// 		log.Fatalf("failed to serve gRPC: %v", err)
	// 	}
	// }()

	// // Start HTTP gateway
	// mux := runtime.NewServeMux()
	// opts := []grpc.DialOption{grpc.WithInsecure()}
	// if err := etlPb.RegisterCategoryServiceHandlerFromEndpoint(ctx, mux, "localhost:5059", opts); err != nil {
	// 	log.Fatalf("failed to start HTTP gateway: %v", err)
	// }

	// log.Println("HTTP Gateway listening on :8099")
	// if err := http.ListenAndServe(":8099", mux); err != nil {
	// 	log.Fatalf("failed to start HTTP server: %v", err)
	// }
}
