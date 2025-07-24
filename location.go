package location

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/garden-raccoon/location-pkg/models"
	proto "github.com/garden-raccoon/location-pkg/protocols/location-pkg"
)

type LocationPkgAPI interface {
	CreateOrUpdateLocation(s *models.Location) error
	DeleteLocation(locationUuid uuid.UUID) error
	LocationByUuid(locationUuid uuid.UUID) (*models.Location, error)
	GetAllLocations() ([]*models.Location, error)
	UpdateLocation(s *models.Location) error
	HealthCheck() error
	// Close GRPC Api connection
	Close() error
}

// Api is profile-service GRPC Api
// structure with client Connection
type Api struct {
	addr    string
	timeout time.Duration
	*grpc.ClientConn
	mu sync.Mutex
	proto.LocationServiceClient
	grpc_health_v1.HealthClient
}

// New create new Battles Api instance
func New(addr string, timeOut time.Duration) (LocationPkgAPI, error) {
	api := &Api{timeout: timeOut}

	if err := api.initConn(addr); err != nil {
		return nil, fmt.Errorf("create LocationApi:  %w", err)
	}
	api.HealthClient = grpc_health_v1.NewHealthClient(api.ClientConn)

	api.LocationServiceClient = proto.NewLocationServiceClient(api.ClientConn)
	return api, nil
}
func (api *Api) UpdateLocation(s *models.Location) error {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()
	_, err := api.LocationServiceClient.UpdateLocation(ctx, s.Proto())
	if err != nil {
		return fmt.Errorf("call MealsByLocation: %w", err)
	}
	return nil
}
func (api *Api) GetAllLocations() ([]*models.Location, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()
	resp, err := api.LocationServiceClient.GetAllLocations(ctx, &proto.EmptyLocation{})
	if err != nil {
		return nil, fmt.Errorf("GetMeals api request: %w", err)
	}
	return models.LocationsFromProto(resp), nil
}
func (api *Api) DeleteLocation(locationUuid uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()
	req := &proto.LocationDeleteReq{
		LocationUuid: locationUuid.Bytes(),
	}
	_, err := api.LocationServiceClient.DeleteLocation(ctx, req)
	if err != nil {
		return fmt.Errorf("DeleteLocation api request: %w", err)
	}
	return nil
}

func (api *Api) CreateOrUpdateLocation(s *models.Location) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()
	_, err = api.LocationServiceClient.CreateOrUpdateLocation(ctx, s.Proto())
	if err != nil {
		return fmt.Errorf("create Location api request: %w", err)
	}
	return nil
}

// initConn initialize connection to Grpc servers
func (api *Api) initConn(addr string) (err error) {
	var kacp = keepalive.ClientParameters{
		Time:                5 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,     // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,            // send pings even without active streams
	}
	connParams := grpc.WithConnectParams(grpc.ConnectParams{
		Backoff: backoff.Config{
			BaseDelay:  100 * time.Millisecond,
			Multiplier: 1.2,
			MaxDelay:   1 * time.Second,
		},
		MinConnectTimeout: 5 * time.Second,
	})
	api.ClientConn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithKeepaliveParams(kacp), connParams)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}
	return
}

func (api *Api) LocationByUuid(LocationUuid uuid.UUID) (*models.Location, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()
	getReq := &proto.LocationGetReq{LocationUuid: LocationUuid.Bytes()}
	resp, err := api.LocationServiceClient.LocationByUUID(ctx, getReq)
	if err != nil {
		return nil, fmt.Errorf("LocationAPI LocationById request failed: %w", err)
	}
	return models.LocationFromProto(resp), nil
}

func (api *Api) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()

	api.mu.Lock()
	defer api.mu.Unlock()

	resp, err := api.HealthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: "locationapi"})
	if err != nil {
		return fmt.Errorf("healthcheck error: %w", err)
	}

	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		return fmt.Errorf("node is %s", errors.New("service is unhealthy"))
	}
	//api.status = service.StatusHealthy
	return nil
}
