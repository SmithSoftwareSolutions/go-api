package services

import (
	"database/sql"
	"errors"
	"log/slog"
	"strings"

	"smithsolutions/go-api/internal/util"
)

type ServiceStatus int

const (
	ServiceStatusRunning ServiceStatus = iota
	ServiceStatusFailed
)

type GetOneOverrider[modelT, includeT any] interface {
	GetOneById(id int, include *includeT) (*modelT, error)
}
type GetManyOverrider[modelT, whereT, includeT any] interface {
	GetMany(where whereT, include *includeT) (*[]modelT, error)
}
type AttachRelationsOverrider[modelT, includeT any] interface {
	AttachRelations(model *modelT, include includeT) error
}

type Creater interface {
	SQL() ([]string, []any, error)
}

type Updater interface {
	SQL() (string, []any, error)
}

type Wherer interface {
	SQL() (string, []any, error)
}

type ResourceService[modelT any, createT Creater, updateT Updater, whereT Wherer, includeT any] struct {
	tableName string
	db        *sql.DB
	columns   []string

	status ServiceStatus

	getOneOverrider          GetOneOverrider[modelT, includeT]
	getManyOverrider         GetManyOverrider[modelT, whereT, includeT]
	attachRelationsOverrider AttachRelationsOverrider[modelT, includeT]
}

func SetupResourceService[modelT any, createT Creater, updateT Updater, whereT Wherer, includeT any](db *sql.DB, tableName string, model any) ResourceService[modelT, createT, updateT, whereT, includeT] {
	columns, err := util.GetColumnsFromModel(model)

	status := ServiceStatusRunning
	if err != nil {
		slog.Error("failed to setup new "+tableName+" service, could not retrieve columns from row struct", "err", err)
		status = ServiceStatusFailed
	}

	return ResourceService[modelT, createT, updateT, whereT, includeT]{
		tableName: tableName,
		db:        db,
		columns:   columns,
		status:    status,
	}
}

func (s *ResourceService[modelT, createT, updateT, whereT, includeT]) Create(data createT) (int, error) {
	// TODO: add overrider switch logic

	if s.status == ServiceStatusFailed {
		return 0, errors.New("service failed to setup or is currently in failed state")
	}

	columns, params, err := data.SQL()

	if err != nil {
		return 0, err
	}

	if len(params) <= 0 {
		return 0, errors.New("no values provided for insert statement")
	}

	paramPlaceholders := strings.Repeat("?, ", len(params))
	paramPlaceholders = paramPlaceholders[:len(paramPlaceholders)-2]

	sql := "INSERT INTO " + s.tableName + " (" + strings.Join(columns, ",") + ") VALUES (" + paramPlaceholders + ")"
	result, err := s.db.Exec(sql, params...)

	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()

	return int(rowsAffected), err
}

func (s *ResourceService[modelT, createT, updateT, whereT, includeT]) GetOneById(id int, include *includeT) (*modelT, error) {
	if s.getOneOverrider != nil {
		return s.getOneOverrider.GetOneById(id, include)
	}

	if s.status == ServiceStatusFailed {
		return nil, errors.New("service failed to setup or is currently in failed state")
	}

	sql := "SELECT " + strings.Join(s.columns, ", ") + " FROM " + s.tableName + " WHERE id=? LIMIT 1"
	params := []any{id}

	var row modelT
	err := util.ScanRow(s.db, &row, sql, params...)

	if err != nil {

		return nil, err
	}

	if include != nil {
		err := s.AttachRelations(&row, *include)
		if err != nil {
			return nil, err
		}
	}

	return &row, nil
}

func (s *ResourceService[modelT, createT, updateT, whereT, includeT]) GetMany(where whereT, include *includeT) (*[]modelT, error) {
	if s.getManyOverrider != nil {
		return s.getManyOverrider.GetMany(where, include)
	}

	whereString, params, err := where.SQL()

	if err != nil {
		return nil, err
	}

	if whereString != "" {
		whereString = " WHERE " + whereString
	}

	sql := "SELECT " + strings.Join(s.columns, ", ") + " FROM " + s.tableName + whereString

	var rows []modelT
	err = util.ScanRows(s.db, &rows, sql, params...)

	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		if include != nil {
			err = s.AttachRelations(&row, *include)
			if err != nil {
				return nil, err
			}
		}
	}

	if rows == nil {
		return &([]modelT{}), nil
	}

	return &rows, nil
}

func (s *ResourceService[modelT, createT, updateT, whereT, includeT]) UpdateOne(id int, data updateT) (int, error) {
	// TODO: add overrider switch logic

	if s.status == ServiceStatusFailed {
		return 0, errors.New("service failed to setup or is currently in failed state")
	}

	setString, params, err := data.SQL()

	if err != nil {
		return 0, err
	}

	params = append(params, id)

	if len(params) <= 0 {
		return 0, errors.New("no values provided for update statement")
	}

	sql := "UPDATE " + s.tableName + " SET " + setString + " WHERE id=? LIMIT 1"
	result, err := s.db.Exec(sql, params...)

	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()

	return int(rowsAffected), err
}

func (s *ResourceService[modelT, createT, updateT, whereT, includeT]) DeleteOneById(id int) (int, error) {
	// TODO: add overrider switch logic

	if s.status == ServiceStatusFailed {
		return 0, errors.New("service failed to setup or is currently in failed state")
	}

	sql := "DELETE FROM " + s.tableName + " WHERE id =? LIMIT 1"
	params := []any{id}

	result, err := s.db.Exec(sql, params...)

	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()

	return int(rowsAffected), err
}

func (s *ResourceService[modelT, createT, updateT, whereT, includeT]) AttachRelations(user *modelT, include includeT) error {
	if s.attachRelationsOverrider != nil {

		return s.attachRelationsOverrider.AttachRelations(user, include)
	}

	return nil
}
