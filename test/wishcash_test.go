package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/worker"

	"github.com/ContextLogic/autobots/pkg/clients"
	"github.com/ContextLogic/autobots/pkg/config"
	"github.com/ContextLogic/autobots/pkg/workflows/wishcashpayment"
	"github.com/ContextLogic/autobots/pkg/workflows/wishcashpayment/models"

	cadencepkg "github.com/ContextLogic/cadence/pkg"
	"github.com/sirupsen/logrus"
	temporal "go.temporal.io/sdk/client"
)

var logger = logrus.New()

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	env *testsuite.TestWorkflowEnvironment
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func TestUnitTestSuite(t *testing.T) {
	t.Skip()
	suite.Run(t, new(UnitTestSuite))
}

func (s *UnitTestSuite) Test_Integration_Workflow() {
	config, err := config.Init(config.GetEnvironment())
	if err != nil {
		panic(err)
	}

	clients, err := clients.Init(config)
	if err != nil {
		panic(err)
	}

	cashPaymentWorkflow := wishcashpayment.NewWishCashPaymentWorkflow(
		config.Clients.Temporal.Clients[wishcashpayment.GetNamespace()],
		clients,
	)

	_, filename, _, _ := runtime.Caller(0)
	wp := path.Join(path.Dir(filename), "../pkg/workflows/wish_cash_payment/workflows.json")
	c, err := cadencepkg.New(temporal.Options{HostPort: config.Clients.Temporal.HostPort, Namespace: wish_cash_payment.GetNamespace()})
	if err != nil {
		panic(err)
	}

	activities := wishcashpayment.GetActivityMap(cashPaymentWorkflow)
	workerOptions := map[string]worker.Options{
		// queue: worker_options
		config.Clients.Temporal.TaskQueuePrefix + "_dsl": worker.Options{},
	}

	_, err = c.Register(wp, workerOptions, activities)
	if err != nil {
		panic(err)
	}

	h := make(http.Header)
	h.Add("Accept", "*/*")
	h.Add("Content-Type", "application/x-www-form-urlencoded")
	h.Add("Cookie", "_xsrf=1; sweeper_session=2|1:0|10:1619126921|15:sweeper_session|84:NjdhZWEyMGQtMGVjMS00MTAxLTliOWItMzNmNzU3MDNhYzA5MjAyMS0wNC0yMCAyMToyNzo1Ny40MTYzMzA=|27bbb42d27cb5b8dcf149efd09bf8dce0cee843a47de700ab8e9b4637f157d39")
	body := "_xsrf=1"
	bytes := []byte(body)

	data := &models.WishCashPaymentWorkflowContext{
		Header: h,
		Body:   bytes,
	}
	s.NoError(AddCart(cashPaymentWorkflow, h))

	instance, err := c.ExecuteWorkflow(
		context.Background(),
		temporal.StartWorkflowOptions{
			ID:        strings.Join([]string{wish_cash_payment.GetNamespace(), strconv.Itoa(int(time.Now().Unix()))}, "_"),
			TaskQueue: config.Clients.Temporal.TaskQueuePrefix + "_dsl",
		},
		"WishCashPaymentWorkflow",
		data,
	)

	response := &wishcashpayment_models.WishCashPaymentResponse{}
	err = instance.Get(context.Background(), &response)
	s.NoError(err)

}

func AddCart(wf *wishcashpayment.WishCashPaymentWorkflow, h http.Header) error {
	params := url.Values{}
	params.Add("_xsrf", `1`)
	params.Add("product_id", `5c4a04e0e6a1c633c8876229`)
	params.Add("variation_id", `5c4a04e0e6a1c633c887622a`)
	params.Add("quantity", `1`)
	params.Add("add_to_cart", `true`)
	params.Add("shipping_option_id", `standard`)
	params.Add("product_source", `tabbed_feed_latest`)
	bytes := []byte(params.Encode())

	wf.Clients.Logger.Info("==========calling wish-fe to add cart: started==========")
	wf.Clients.Logger.WithFields(logrus.Fields{"headers": h, "body": string(bytes)}).Info("create order request info")
	bytes, err := wf.Clients.WishFrontend.Post(h, bytes, "api/cart/update")
	if err != nil {
		return err
	}

	response := &models.WishCashPaymentCreateOrderResponse{}
	if err = json.Unmarshal(bytes, response); err != nil {
		return err
	}

	wf.Clients.Logger.Info("==========calling wish-fe to add cart: finished==========")
	wf.Clients.Logger.WithFields(logrus.Fields{"Msg": response.Msg}).Info("add cart response info")
	return nil
}
