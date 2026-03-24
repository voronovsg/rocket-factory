package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/voronovsg/rocket-factory/inventory/internal/interceptor/logger"
	"github.com/voronovsg/rocket-factory/inventory/internal/interceptor/validate"
	inventoryV1 "github.com/voronovsg/rocket-factory/inventory/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"sync"
)

const (
	grpcAddr = "localhost:50051"
	httpAddr = "localhost:8081"
)

type InventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer
	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

type filterLookup struct {
	uuids      map[string]struct{}
	names      map[string]struct{}
	categories map[inventoryV1.Category]struct{}
	countries  map[string]struct{}
	tags       map[string]struct{}
}

func (s *InventoryService) initParts() {
	parts := generateParts()
	for _, part := range parts {
		s.parts[part.Uuid] = part
	}
}

func generateParts() []*inventoryV1.Part {
	names := []string{
		"Main Engine",
		"Reserve Engine",
		"Thruster",
		"Fuel Tank",
		"Left Wing",
		"Right Wing",
		"Window A",
		"Window B",
		"Control Module",
		"Stabilizer",
	}

	descriptions := []string{
		"Primary propulsion unit",
		"Backup propulsion unit",
		"Thruster for fine adjustments",
		"Main fuel tank",
		"Left aerodynamic wing",
		"Right aerodynamic wing",
		"Front viewing window",
		"Side viewing window",
		"Flight control module",
		"Stabilization fin",
	}

	var parts []*inventoryV1.Part
	for i := 0; i < gofakeit.Number(1, 50); i++ {
		idx := gofakeit.Number(0, len(names)-1)
		parts = append(parts, &inventoryV1.Part{
			Uuid:          uuid.NewString(),
			Name:          names[idx],
			Description:   descriptions[idx],
			Price:         roundTo(gofakeit.Float64Range(100, 10_000)),
			StockQuantity: int64(gofakeit.Number(1, 100)),
			Category:      inventoryV1.Category(gofakeit.Number(1, 4)), //nolint:gosec // safe: gofakeit.Number returns 1..4
			Dimensions:    generateDimensions(),
			Manufacturer:  generateManufacturer(),
			Tags:          generateTags(),
			Metadata:      generateMetadata(),
			CreatedAt:     timestamppb.Now(),
		})
	}

	return parts
}

func generateDimensions() *inventoryV1.Dimensions {
	return &inventoryV1.Dimensions{
		Length: roundTo(gofakeit.Float64Range(1, 1000)),
		Width:  roundTo(gofakeit.Float64Range(1, 1000)),
		Height: roundTo(gofakeit.Float64Range(1, 1000)),
		Weight: roundTo(gofakeit.Float64Range(1, 1000)),
	}
}

func generateManufacturer() *inventoryV1.Manufacturer {
	return &inventoryV1.Manufacturer{
		Name:    gofakeit.Name(),
		Country: gofakeit.Country(),
		Website: gofakeit.URL(),
	}
}

func generateTags() []string {
	var tags []string
	for i := 0; i < gofakeit.Number(1, 10); i++ {
		tags = append(tags, gofakeit.EmojiTag())
	}

	return tags
}

func generateMetadata() map[string]*inventoryV1.Value {
	metadata := make(map[string]*inventoryV1.Value)

	for i := 0; i < gofakeit.Number(1, 10); i++ {
		metadata[gofakeit.Word()] = generateMetadataValue()
	}

	return metadata
}

func generateMetadataValue() *inventoryV1.Value {
	switch gofakeit.Number(0, 3) {
	case 0:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_StringValue{
				StringValue: gofakeit.Word(),
			},
		}

	case 1:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_Int64Value{
				Int64Value: int64(gofakeit.Number(1, 100)),
			},
		}

	case 2:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_DoubleValue{
				DoubleValue: roundTo(gofakeit.Float64Range(1, 100)),
			},
		}

	case 3:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_BoolValue{
				BoolValue: gofakeit.Bool(),
			},
		}

	default:
		return nil
	}
}

func roundTo(x float64) float64 {
	return math.Round(x*100) / 100
}

// preFilterLookup создает lookup таблицы из фильтра для быстрого поиска
func preFilterLookup(filter *inventoryV1.PartsFilter) *filterLookup {
	lookup := &filterLookup{
		uuids:      make(map[string]struct{}, len(filter.GetUuids())),
		names:      make(map[string]struct{}, len(filter.GetNames())),
		categories: make(map[inventoryV1.Category]struct{}, len(filter.GetCategories())),
		countries:  make(map[string]struct{}, len(filter.GetManufacturerCountries())),
		tags:       make(map[string]struct{}, len(filter.GetTags())),
	}

	for _, uuid := range filter.GetUuids() {
		lookup.uuids[uuid] = struct{}{}
	}
	for _, name := range filter.GetNames() {
		lookup.names[strings.ToLower(name)] = struct{}{}
	}
	for _, category := range filter.GetCategories() {
		lookup.categories[category] = struct{}{}
	}
	for _, country := range filter.GetManufacturerCountries() {
		lookup.countries[strings.ToLower(country)] = struct{}{}
	}
	for _, tag := range filter.GetTags() {
		lookup.tags[strings.ToLower(tag)] = struct{}{}
	}

	return lookup
}

func isEmptyFilter(filter *inventoryV1.PartsFilter) bool {
	if filter == nil {
		return true
	}
	return len(filter.GetUuids()) == 0 &&
		len(filter.GetNames()) == 0 &&
		len(filter.GetCategories()) == 0 &&
		len(filter.GetManufacturerCountries()) == 0 &&
		len(filter.GetTags()) == 0
}

// matchesFilter проверяет, соответствует ли деталь всем критериям фильтра
func matchesFilter(part *inventoryV1.Part, lookup *filterLookup) bool {
	if !matchesUUIDFilter(part, lookup) {
		return false
	}
	if !matchesNameFilter(part, lookup) {
		return false
	}
	if !matchesCategoryFilter(part, lookup) {
		return false
	}
	if !matchesCountryFilter(part, lookup) {
		return false
	}
	if !matchesTagsFilter(part, lookup) {
		return false
	}
	return true
}

// matchesUUIDFilter проверяет соответствие по UUID
func matchesUUIDFilter(part *inventoryV1.Part, lookup *filterLookup) bool {
	if len(lookup.uuids) == 0 {
		return true
	}
	_, exists := lookup.uuids[part.GetUuid()]
	return exists
}

// matchesNameFilter проверяет соответствие по названию
func matchesNameFilter(part *inventoryV1.Part, lookup *filterLookup) bool {
	if len(lookup.names) == 0 {
		return true
	}
	_, exists := lookup.names[strings.ToLower(part.GetName())]
	return exists
}

// matchesCategoryFilter проверяет соответствие по категории
func matchesCategoryFilter(part *inventoryV1.Part, lookup *filterLookup) bool {
	if len(lookup.categories) == 0 {
		return true
	}
	_, exists := lookup.categories[part.GetCategory()]
	return exists
}

// matchesCountryFilter проверяет соответствие по стране
func matchesCountryFilter(part *inventoryV1.Part, lookup *filterLookup) bool {
	if len(lookup.countries) == 0 {
		return true
	}
	manufacturer := part.GetManufacturer()
	if manufacturer == nil {
		return false
	}
	_, exists := lookup.countries[strings.ToLower(manufacturer.GetCountry())]
	return exists
}

// matchesTagsFilter проверяет соответствие по тегам
func matchesTagsFilter(part *inventoryV1.Part, lookup *filterLookup) bool {
	if len(lookup.tags) == 0 {
		return true
	}
	tags := part.GetTags()
	for _, tag := range tags {
		if _, exists := lookup.tags[strings.ToLower(tag)]; exists {
			return true
		}
	}
	return false
}

// GetPart Возвращает информацию о детали по её UUID
func (s *InventoryService) GetPart(_ context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	parts, ok := s.parts[req.Uuid]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.Uuid)
	}

	return &inventoryV1.GetPartResponse{Part: parts}, nil
}

// ListParts возвращает список деталей с возможностью фильтрации
func (s *InventoryService) ListParts(_ context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filter := req.GetFilter()

	// Если фильтр пустой, возвращаем все элементы
	if isEmptyFilter(filter) {
		result := make([]*inventoryV1.Part, 0, len(s.parts))
		for _, part := range s.parts {
			result = append(result, part)
		}
		return &inventoryV1.ListPartsResponse{Parts: result}, nil
	}

	// Подготавливаем структуры для быстрого поиска
	lookup := preFilterLookup(filter)
	result := make([]*inventoryV1.Part, 0, len(s.parts))
	for _, part := range s.parts {
		if matchesFilter(part, lookup) {
			result = append(result, part)
		}
	}

	return &inventoryV1.ListPartsResponse{Parts: result}, nil
}

func main() {
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			validate.UnaryValidateInterceptor(),
			logger.UnaryLoggerInterceptor(),
		),
	)
	reflection.Register(s)

	service := &InventoryService{
		parts: make(map[string]*inventoryV1.Part),
	}
	service.initParts()
	log.Printf("📦 Generated %d parts for inventory", len(service.parts))

	inventoryV1.RegisterInventoryServiceServer(s, service)

	go func() {
		log.Printf("🚀 gRPC server listening on %v\n", grpcAddr)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	var gwServer *http.Server
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		err = inventoryV1.RegisterInventoryServiceHandlerFromEndpoint(
			ctx,
			mux,
			fmt.Sprintf(grpcAddr),
			opts,
		)
		if err != nil {
			log.Printf("Failed to register gateway: %v\n", err)
			return
		}

		fileServer := http.FileServer(http.Dir("api"))
		httpMux := http.NewServeMux()
		httpMux.Handle("/api/", mux)
		httpMux.Handle("/swagger-ui.html", fileServer)
		httpMux.Handle("/inventory.swagger.json", fileServer)
		httpMux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				http.Redirect(w, r, "/swagger-ui.html", http.StatusMovedPermanently)
				return
			}
			fileServer.ServeHTTP(w, r)
		}))

		gwServer = &http.Server{
			Addr:              fmt.Sprintf(httpAddr),
			Handler:           httpMux,
			ReadHeaderTimeout: 10 * time.Second,
		}
		log.Printf("🌐 HTTP server with gRPC-Gateway and Swagger UI listening on %v\n", httpAddr)
		err = gwServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Failed to serve HTTP: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down servers...")

	if gwServer != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := gwServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
		log.Println("✅ HTTP server stopped")
	}

	s.GracefulStop()
	log.Println("✅ Server stopped")
}
