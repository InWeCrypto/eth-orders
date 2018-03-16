package orders

import (
	"sync"

	"github.com/dynamicgo/config"
	"github.com/dynamicgo/slf4go"
	"github.com/go-xorm/xorm"
	"github.com/inwecrypto/ethdb"
	"github.com/inwecrypto/gomq"
	kafka "github.com/inwecrypto/gomq-kafka"
)

type txWatcher struct {
	mq gomq.Consumer
	slf4go.Logger
	db       *xorm.Engine
	marked   int64
	handlers int64
	sync.Mutex
}

func newTxWatcher(conf *config.Config, db *xorm.Engine) (*txWatcher, error) {

	mq, err := kafka.NewAliyunConsumer(conf)

	if err != nil {
		return nil, err
	}

	return &txWatcher{
		mq:       mq,
		Logger:   slf4go.Get("txwatcher"),
		db:       db,
		handlers: config.GetInt64("orders.handlers", 10),
	}, nil
}

func (watcher *txWatcher) Run() {

	go func() {
		for err := range watcher.mq.Errors() {
			watcher.ErrorF("mq error %s", err)
		}
	}()

	for i := int64(0); i < watcher.handlers; i++ {
		watcher.doRun()
	}
}

func (watcher *txWatcher) doRun() {
	for message := range watcher.mq.Messages() {
		if err := watcher.handleTx(string(message.Key())); err != nil {
			watcher.ErrorF("handle tx %s err, %s", string(message.Key()), err)
		}

		watcher.commitMessage(message)
	}
}
func (watcher *txWatcher) commitMessage(message gomq.Message) {
	watcher.Lock()
	defer watcher.Unlock()

	if watcher.marked < message.Offset() {
		watcher.marked = message.Offset()

		watcher.mq.Commit(message)
	}
}

func (watcher *txWatcher) handleTx(tx string) error {

	// watcher.DebugF("handle tx %s", tx)

	ethTx := new(ethdb.TableTx)

	ok, err := watcher.db.Where("t_x = ?", tx).Get(ethTx)

	if err != nil {
		return err
	}

	if !ok {
		watcher.WarnF("handle tx %s -- not found", tx)
		return nil
	}

	order := new(ethdb.TableOrder)

	order.ConfirmTime = &ethTx.CreateTime
	order.Blocks = int64(ethTx.Blocks)

	updated, err := watcher.db.Where("t_x = ?", tx).Cols("confirm_time", "blocks").Update(order)

	if err != nil {
		return err
	}

	if updated != 0 {
		watcher.DebugF("updated orders(%d) for tx %s", updated, tx)
		return nil
	}

	wallet := new(ethdb.TableWallet)

	count, err := watcher.db.Where(`"address" = ? or "address" = ?`, ethTx.From, ethTx.To).Count(wallet)

	if err != nil {
		return err
	}

	if count > 0 {
		order.Asset = ethTx.Asset
		order.From = ethTx.From
		order.To = ethTx.To
		order.TX = ethTx.TX
		order.Value = ethTx.Value
		order.CreateTime = ethTx.CreateTime

		_, err = watcher.db.Insert(order)

		return err
	}

	return nil
}
