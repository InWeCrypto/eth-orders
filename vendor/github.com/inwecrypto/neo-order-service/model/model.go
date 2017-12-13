package model

import (
	"database/sql"
	"fmt"

	"github.com/dynamicgo/config"
	"github.com/dynamicgo/slf4go"
)

// DBModel .
type DBModel struct {
	db  *sql.DB
	cnf *config.Config
	slf4go.Logger
}

// NewDBModel .
func NewDBModel(conf *config.Config, db *sql.DB) *DBModel {
	return &DBModel{
		cnf:    conf,
		db:     db,
		Logger: slf4go.Get("neo-order-service-model"),
	}
}

// GetSQL .
func (model *DBModel) GetSQL(name string) string {
	if !model.cnf.Has(name) {
		panic(fmt.Sprintf("unknown sql %s", name))
	}

	return model.cnf.GetString(name, "xxx")
}

// Tx execute a tx
func (model *DBModel) Tx(proc func(tx *sql.Tx) error) (reterr error) {
	tx, err := model.db.Begin()

	if err != nil {
		return err
	}

	defer func() {
		if err := recover(); err != nil {
			reterr = err.(error)
		}
	}()

	if err := proc(tx); err != nil {
		if rollbackError := tx.Rollback(); rollbackError != nil {
			model.ErrorF("rollback err, %s", rollbackError)
		}

		return err
	}

	return tx.Commit()
}

// WalletModel model model
type WalletModel struct {
	*DBModel
}

// Wallet .
type Wallet struct {
	Address    string `json:"address"`
	UserID     string `json:"userid"`
	CreateTime string `json:"createTime"`
}

// Create create new model
func (model *WalletModel) Create(address string, userid string) error {

	if address == "" {
		return fmt.Errorf("address param can't be empty string")
	}

	if userid == "" {
		return fmt.Errorf("userid param can't be empty string")
	}

	query := model.GetSQL("nos.orm.wallet.create")

	return model.Tx(func(tx *sql.Tx) error {

		model.DebugF("create wallet sql :%s address:%s userid:%s", query, address, userid)

		_, err := tx.Exec(query, address, userid)

		if err != nil {
			return err
		}

		return nil

	})
}

// Delete delete model indicate by userid and address
func (model *WalletModel) Delete(address string, userid string) error {

	if address == "" {
		return fmt.Errorf("address param can't be empty string")
	}

	if userid == "" {
		return fmt.Errorf("userid param can't be empty string")
	}

	query := model.GetSQL("nos.orm.wallet.delete")

	return model.Tx(func(tx *sql.Tx) error {

		model.DebugF("delete model sql :%s address:%s userid:%s", query, address, userid)

		_, err := tx.Exec(query, address, userid)

		if err != nil {
			return err
		}

		return nil

	})
}

// GetByAddress get wallet by address
func (model *WalletModel) GetByAddress(address string) (find *Wallet, err error) {

	if address == "" {
		return nil, fmt.Errorf("address param can't be empty string")
	}

	query := model.GetSQL("nos.orm.wallet.getbyaddress")

	err = model.Tx(func(tx *sql.Tx) error {

		model.DebugF("get wallet sql :%s address:%s ", query, address)

		rows, err := tx.Query(query, address)

		if err != nil {
			return err
		}

		defer rows.Close()

		if rows.Next() {
			var wallet Wallet

			if err := rows.Scan(&wallet.Address, &wallet.UserID, &wallet.CreateTime); err != nil {
				return err
			}

			find = &wallet
		}

		return nil

	})

	return
}

// TxModel .
type TxModel struct {
	*DBModel
}

// Tx .
type Tx struct {
	ID         string
	Blocks     uint64
	TX         string
	Address    string
	Type       string
	Assert     string
	Value      float64
	UpdateTime string
}

// GetByID get tx object list by tx id
func (model *TxModel) GetByID(id string) ([]*Tx, error) {

	query := model.GetSQL("nos.orm.tx.id")

	model.DebugF("get tx by id: %s with id %s", query, id)

	rows, err := model.db.Query(query, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []*Tx

	for rows.Next() {
		var tx Tx
		if err := rows.Scan(
			&tx.ID,
			&tx.Blocks,
			&tx.Address,
			&tx.Type,
			&tx.Assert,
			&tx.Value,
			&tx.UpdateTime); err != nil {

			return nil, err

		}

		tx.TX = id

		result = append(result, &tx)
	}

	return result, nil
}

// OrderModel neo order model
type OrderModel struct {
	*DBModel
}

// Order neo order object
type Order struct {
	Tx          string `json:"tx" form:"tx" binding:"required"`
	From        string `json:"from" form:"from" binding:"required"`
	To          string `json:"to" form:"to" binding:"required"`
	Asset       string `json:"asset" form:"asset" binding:"required"`
	Value       string `json:"value" form:"value" binding:"required"`
	CreateTime  string `json:"createTime" form:"createTime"`
	ConfirmTime string `json:"confirmTime" form:"confirmTime"`
}

// Create create new ordereapig
func (model *OrderModel) Create(order *Order) error {

	query := model.GetSQL("nos.orm.order.create")

	return model.Tx(func(tx *sql.Tx) error {

		model.DebugF("create order sql :%s order :%s", query, order.Tx)

		_, err := tx.Exec(query, order.Tx, order.From, order.To, order.Asset, order.Value)

		if err != nil {
			return err
		}

		return nil

	})
}

// Confirm confirm order
// func (model *OrderModel) Confirm(txid string) error {

// 	if txid == "" {
// 		return fmt.Errorf("tx param can't be empty string")
// 	}

// 	query := model.GetSQL("nos.orm.order.confirm")

// 	return model.Tx(func(tx *sql.Tx) error {

// 		model.DebugF("confirm sql :%s tx :%s", query, txid)

// 		_, err := tx.Exec(query, txid)

// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// }

// Status .
func (model *OrderModel) Status(txid string) (ok bool, err error) {
	if txid == "" {
		return false, fmt.Errorf("txid param can't be empty string")
	}

	query := model.GetSQL("nos.orm.order.status")

	err = model.Tx(func(tx *sql.Tx) error {

		model.DebugF("status sql :%s tx :%s", query, txid)

		rows, err := tx.Query(query, txid)

		if err != nil {
			return err
		}

		defer rows.Close()

		if rows.Next() {
			ok = true
		}

		return nil

	})

	return
}

// Confirm .
func (model *OrderModel) Confirm(txid string) (err error) {

	model.DebugF("confirm order with tx %s", txid)

	txModel := &TxModel{DBModel: model.DBModel}

	txs, err := txModel.GetByID(txid)

	if err != nil {
		return err
	}

	if len(txs) == 0 {
		model.WarnF("confirm order with tx %s, tx not found", txid)
		return nil
	}

	createQuery := model.GetSQL("nos.orm.order.createWithConfirm")

	confirmQuery := model.GetSQL("nos.orm.order.confirm")

	return model.Tx(func(tx *sql.Tx) error {

		addressed, err := model.getAddresses(tx, txid)

		var selectTx *Tx

		for _, address := range addressed {
			for _, tx := range txs {
				if tx.Address == address {
					selectTx = tx
				}
			}

			if selectTx == nil {
				model.WarnF("tx %s to known address %s not found", txid, address)
				continue
			}

			// model.DebugF("create sql %s with %s %s %s %.8f", createQuery, txid, address, selectTx.Assert, selectTx.Value)

			_, err := tx.Exec(createQuery, txid, "", address, selectTx.Assert, selectTx.Value)

			if err != nil {
				model.ErrorF("create order for tx %s error %s", txid, err)
				return err
			}
		}

		_, err = tx.Exec(confirmQuery, txid)

		return err
	})
}

func (model *OrderModel) getAddresses(tx *sql.Tx, txid string) ([]string, error) {
	query := model.GetSQL("nos.orm.order.wallet")

	rows, err := tx.Query(query, txid)

	if err != nil {
		return nil, err
	}

	model.DebugF("query %s with %s", query, txid)

	defer rows.Close()

	var addresses []string

	for rows.Next() {

		var address string

		err = rows.Scan(&address)

		if err != nil {
			return nil, err
		}

		model.DebugF("tx %s to known address %s", txid, address)

		addresses = append(addresses, address)
	}

	return addresses, nil
}

// Orders get orders
func (model *OrderModel) Orders(address string, asset string, page *Page) (orders []*Order, err error) {
	if address == "" {
		return nil, fmt.Errorf("address param can't be empty string")
	}

	orders = make([]*Order, 0)

	query := model.GetSQL("nos.orm.order.list")

	err = model.Tx(func(tx *sql.Tx) error {

		model.DebugF("list sql :%s address:%s page: %d size: %d", query, address, page.Offset, page.Size)

		rows, err := tx.Query(query, address, asset, page.Offset*page.Size, page.Size)

		if err != nil {
			return err
		}

		defer rows.Close()

		for rows.Next() {

			var order Order

			var confirmTime sql.NullString

			err = rows.Scan(
				&order.Tx,
				&order.From,
				&order.To,
				&order.Asset,
				&order.Value,
				&order.CreateTime,
				&confirmTime)

			if err != nil {
				return err
			}

			order.ConfirmTime = confirmTime.String

			orders = append(orders, &order)
		}

		return nil
	})

	return
}

// Order get order by txid
func (model *OrderModel) Order(txid string) (find *Order, err error) {
	if txid == "" {
		return nil, fmt.Errorf("txid param can't be empty string")
	}

	query := model.GetSQL("nos.orm.order.get")

	err = model.Tx(func(tx *sql.Tx) error {

		model.DebugF("get order sql :%s txid:%s", query, txid)

		rows, err := tx.Query(query, txid)

		if err != nil {
			model.ErrorF("query txid %s err: %s", txid, err)
			return err
		}

		defer rows.Close()

		if rows.Next() {

			var order Order

			var confirmTime sql.NullString

			err = rows.Scan(
				&order.Tx,
				&order.From,
				&order.To,
				&order.Asset,
				&order.Value,
				&order.CreateTime,
				&confirmTime)

			if err != nil {
				model.ErrorF("row next txid %s err: %s", txid, err)
				return err
			}

			order.ConfirmTime = confirmTime.String

			find = &order

		}

		return nil
	})

	return
}

// Page .
type Page struct {
	Offset uint `json:"offset" binding:"required"` // page offset number
	Size   uint `json:"size" binding:"required"`   // page size
}
