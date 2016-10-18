package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	dmp "github.com/sergi/go-diff/diffmatchpatch"

	"app/utils"
)

type database interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Get(interface{}, string, ...interface{}) error
	Select(interface{}, string, ...interface{}) error
}

type Env struct {
	Db  database
	Dmp *dmp.DiffMatchPatch
}

func New(db *sqlx.DB) *Env {
	return &Env{db, dmp.New()}
}

// http://stackoverflow.com/a/23502629/242684
func (e *Env) Atomic(op func(e *Env) error) error {
	// TODO: check the performance hit of this casting
	db, ok := e.Db.(*sqlx.DB)
	if !ok {
		// TODO: should this be a panic, maybe? or is the sole need for all this an
		//       indicator of bad design? decisions decisions...
		return utils.ErrTODO
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	etx := *e
	etx.Db = tx
	err = op(&etx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func BetterGetterErrors(err error) error {
	if err == sql.ErrNoRows {
		return utils.ErrNotFound
	}

	// 22P02 stands for invalid_text_representation, for when something other than
	// a valid integer is supplied
	if err, ok := err.(*pq.Error); ok && err.Code == "22P02" {
		return utils.ErrNotFound
	}

	return err
}

type JSONArray []string

func (a JSONArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}

	return json.Marshal(a)
}

func (a JSONArray) Scan(src interface{}) error {
	if v := reflect.ValueOf(src); !v.IsValid() || v.IsNil() {
		return nil
	} else if data, ok := src.([]byte); ok {
		return json.Unmarshal(data, &a)
	}

	return errors.New("JSONArray: scan source was not []byte")
}

type JSONMap map[string]string

func (m JSONMap) Value() (driver.Value, error) {
	if len(m) == 0 {
		return nil, nil
	}

	return json.Marshal(m)
}

func (m JSONMap) Scan(src interface{}) error {
	if v := reflect.ValueOf(src); !v.IsValid() || v.IsNil() {
		return nil
	} else if data, ok := src.([]byte); ok {
		return json.Unmarshal(data, &m)
	}

	return errors.New("JSONMap: scan source was not []byte")
}

func (e *Env) Diff(
	tableName, rowId, columnName, by, before, after string,
) error {
	patch := e.Dmp.PatchToText(e.Dmp.PatchMake(after, before))
	if patch == "" {
		return nil
	}

	_, err := e.Db.Exec(`
    INSERT INTO diff (
      table_name,
      row_id,
      column_name,
      reverse_patch,
      created_by
    )
    VALUES ($1, $2, $3, $4, $5)`,
		tableName,
		rowId,
		columnName,
		patch,
		by,
	)
	return err
}
