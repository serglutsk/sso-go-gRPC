package suite

import (
	"context"
	"net"
	"os"
	"strconv"
	"sso/internal/config"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc"
	"testing"

	ssov1 "github.com/serglutsk/gRPC-sso/gen/go/sso"
)

const (
	grpcHost = "localhost"
)

type Suite struct {
	*testing.T
	Cfg *config.Config
	AuthClient ssov1.AuthClient
}

func New(t *testing.T)(context.Context, *Suite) {
	t.Helper() // Цей метод каже системі тестування Go,
	//  що дана функція є допоміжною (helper), а не самим тестом. Допомагає при дебагу, оскільки при помилці в тесті,
	//  система тестування буде показувати рядок коду, де була викликана ця функція, а не рядок всередині цієї функції. Це робить повідомлення про помилки більш зрозумілими та корисними для розробника.

	t.Parallel() // Цей метод дозволяє запускати цей тест паралельно з іншими тестами. 
	// Це може значно скоротити час виконання тестів, особливо якщо у вас є багато тестів, які не залежать один від одного.
	//  Коли ви викликаєте t.Parallel(), Go запускає цей тест в окремій горутині, дозволяючи іншим тестам виконуватися одночасно.
	
	cfg := config.MustLoadPath("../config/local_test.yaml")
	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)
	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.DialContext(context.Background(),
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials())) // insecure-коннект для тестів
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssov1.NewAuthClient(cc),
	}

}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}

func configPath() string {
	const key = "CONFIG_PATH"

	if v := os.Getenv(key); v != "" {
		return v
	}

	return "../config/local_tests.yaml"
}
