package e2e

import (
	"context"
	"log"
	"os/exec"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/skybytescode/microservices-proto/golang/order"
)

type CreateOrderTestSuite struct {
	suite.Suite
	project string
}

// compose runs docker compose for this suite's stack, isolated by project name.
func (c *CreateOrderTestSuite) compose(args ...string) error {
	base := []string{"compose", "-f", "resources/docker-compose.yml", "-p", c.project}
	cmd := exec.Command("docker", append(base, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("docker compose %s:\n%s", strings.Join(args, " "), out)
	}
	return err
}

func (c *CreateOrderTestSuite) SetupSuite() {
	c.project = strings.ToLower(uuid.New().String())
	if err := c.compose("up", "-d", "--build", "--wait"); err != nil {
		log.Fatalf("Could not run compose stack: %v", err)
	}
}

func (c *CreateOrderTestSuite) Test_Should_Create_Order() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("Failed to connect order service. Err: %v", err)
	}

	defer conn.Close()

	orderClient := order.NewOrderClient(conn)
	createOrderResponse, errCreate := orderClient.Create(context.Background(), &order.CreateOrderRequest{
		UserId: 23,
		OrderItems: []*order.OrderItem{
			{
				ProductCode: "CAM123",
				Quantity:    3,
				UnitPrice:   1.23,
			},
		},
	})
	c.Require().Nil(errCreate)

	getOrderResponse, errGet := orderClient.Get(context.Background(), &order.GetOrderRequest{OrderId: createOrderResponse.OrderId})
	c.Require().Nil(errGet)
	c.Equal(int64(23), getOrderResponse.UserId)
	orderItem := getOrderResponse.OrderItems[0]
	c.Equal(float32(1.23), orderItem.UnitPrice)
	c.Equal(int32(3), orderItem.Quantity)
	c.Equal("CAM123", orderItem.ProductCode)
}

func (c *CreateOrderTestSuite) TearDownSuite() {
	if err := c.compose("down", "-v"); err != nil {
		log.Fatalf("Could not shutdown compose stack: %v", err)
	}
}

func TestCreateOrderTestSuite(t *testing.T) {
	suite.Run(t, new(CreateOrderTestSuite))
}
