package orders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/dghubble/sling"

	"github.com/inwecrypto/neo-order-service/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateWallet(t *testing.T) {
	resp, err := http.Post("http://localhost:8000/wallet/xxxxx/test", "application/json", strings.NewReader("{}"))

	if assert.NoError(t, err) {
		assert.Equal(t, 200, resp.StatusCode)
	}
}

func TestDeleteWallet(t *testing.T) {

	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8000/wallet/xxxxx/test", nil)

	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)

	if assert.NoError(t, err) {
		assert.Equal(t, 200, resp.StatusCode)
	}
}

// Order .
type Order struct {
	ID          int64   `json:"-" xorm:"pk autoincr"`
	TX          string  `json:"tx" xorm:"notnull"`
	From        string  `json:"from" xorm:"index(from_to)"`
	To          string  `json:"to" xorm:"index(from_to)"`
	Asset       string  `json:"asset" xorm:"notnull"`
	Value       string  `json:"value" xorm:"notnull"`
	Blocks      uint64  `json:"blocks" xorm:""`
	CreateTime  *string `json:"createTime,omitempty" xorm:"TIMESTAMP notnull created"`
	ConfirmTime *string `json:"confirmTime,omitempty" xorm:"TIMESTAMP"`
	Context     *string `json:"context" xorm:"json"`
}

type Comment struct {
	Message string `json:"message"`
}

func TestCreateOrder(t *testing.T) {

	hello := `["test",100]`

	order, err := json.Marshal(&Order{
		TX:      "xxxxxx",
		From:    "xxxxxxx",
		To:      "xxxxxxx",
		Asset:   "xxxxxxxxxxx",
		Value:   "1",
		Context: &hello,
	})

	assert.NoError(t, err)

	resp, err := http.Post("http://localhost:8000/order", "application/json", bytes.NewReader(order))

	if assert.NoError(t, err) {
		assert.Equal(t, 200, resp.StatusCode)
	}
}

func TestGetOrder(t *testing.T) {

	resp, err := http.Get("http://localhost:8000/order/xxxxxx")

	if assert.NoError(t, err) {
		assert.Equal(t, 200, resp.StatusCode)
	}

	data, _ := ioutil.ReadAll(resp.Body)

	println(string(data))
}

func printResult(result interface{}) {

	data, _ := json.MarshalIndent(result, "", "\t")

	fmt.Println(string(data))
}

func TestListOrder(t *testing.T) {
	request, err := sling.New().Get("http://localhost:8000/orders/0x8214b824927a28dc16581cd22e460fe0f7e31994/0x0000000000000000000000000000000000000000/0/20").Request()

	assert.NoError(t, err)

	var orders []*model.Order
	var errmsg interface{}

	_, err = sling.New().Do(request, &orders, &errmsg)
	assert.NoError(t, err)

	assert.NotZero(t, len(orders))

	printResult(orders)

}
