package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/riskibarqy/go-template/utils"
)

// ErrNotEnough declare specific error for Not Enough
// ErrExisted declare specific error for data already exist
var (
	ErrNotFound     = fmt.Errorf("data is not found")
	ErrAlreadyExist = fmt.Errorf("data already exists")
)

// GenericStorage represents the generic Storage
// for the domain models that matches with its database models
type GenericStorage interface {
	Single(ctx context.Context, elem interface{}, where string, arg map[string]interface{}) error
	Where(ctx context.Context, elems interface{}, where string, arg map[string]interface{}) error
	SelectWithQuery(ctx context.Context, elem interface{}, query string, args map[string]interface{}) error
	FindByID(ctx context.Context, elem interface{}, id interface{}) error
	FindAll(ctx context.Context, elems interface{}, page int, limit int) error
	Insert(ctx context.Context, elem interface{}) error
	// InsertMany(ctx context.Context, elem interface{}) error
	Update(ctx context.Context, elem interface{}) error
	Delete(ctx context.Context, id interface{}) error
	DeleteHard(ctx context.Context, id interface{}) error
}

// PostgresStorage is the postgres implementation of generic Storage
type PostgresStorage struct {
	db              Queryer
	tableName       string
	elemType        reflect.Type
	selectFields    string
	insertFields    string
	insertParams    string
	updateSetFields string
}

// Single queries an element according to the query & argument provided
func (r *PostgresStorage) Single(ctx context.Context, elem interface{}, where string, arg map[string]interface{}) error {
	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}

	statement, err := db.PrepareNamed(fmt.Sprintf(`SELECT %s FROM "%s" WHERE %s`,
		r.selectFields, r.tableName, where))
	if err != nil {
		return err
	}
	defer statement.Close()

	err = statement.Get(elem, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}

	return nil
}

// Where queries the elements according to the query & argument provided
func (r *PostgresStorage) Where(ctx context.Context, elems interface{}, where string, arg map[string]interface{}) error {
	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}

	query := fmt.Sprintf(`SELECT %s FROM "%s" WHERE %s`, r.selectFields, r.tableName, where)
	query, args, err := sqlx.Named(query, arg)
	if err != nil {
		return err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}

	query = db.Rebind(query)

	err = db.Select(elems, query, args...)
	if err != nil {
		return err
	}

	return nil
}

// SelectWithQuery Customizable Query for Select
func (r *PostgresStorage) SelectWithQuery(ctx context.Context, elems interface{}, query string, arg map[string]interface{}) error {
	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}

	query, args, err := sqlx.Named(query, arg)
	if err != nil {
		return err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}

	query = db.Rebind(query)

	err = db.Select(elems, query, args...)
	if err != nil {
		return err
	}

	return nil
}

// FindByID finds an element by its id
// it's defined in this project context that
// the element id column in the db should be "id"
func (r *PostgresStorage) FindByID(ctx context.Context, elem interface{}, id interface{}) error {
	where := `"id" = :id`
	err := r.Single(ctx, elem, where, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return err
	}

	return nil
}

// FindAll finds all elements from the database.
func (r *PostgresStorage) FindAll(ctx context.Context, elems interface{}, page int, limit int) error {
	where := `true`
	where = fmt.Sprintf(`%s ORDER BY "id" DESC LIMIT :limit OFFSET :offset`, where)

	err := r.Where(ctx, elems, where, map[string]interface{}{
		"limit":  limit,
		"offset": (page - 1) * limit,
	})

	if err != nil {
		return err
	}

	return nil
}

// Insert inserts a new element into the database.
// It assumes the primary key of the table is "id" with serial type.
// It will set the "owner" field of the element with the current account in the context if exists.
// It will set the "createdAt" and "updatedAt" fields with current time.
// If immutable set true, it won't insert the updatedAt
func (r *PostgresStorage) Insert(ctx context.Context, elem interface{}) error {
	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}

	query := fmt.Sprintf(`
	INSERT INTO "%s"(%s)
	VALUES (%s)
	RETURNING %s`, r.tableName, r.insertFields, r.insertParams, r.selectFields)
	statement, err := db.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer statement.Close()

	dbArgs := r.insertArgs(elem, 0)
	err = statement.Get(elem, dbArgs)
	if err != nil {
		return err
	}

	return nil
}

// // InsertMany is function for creating many datas into specific table in database.
// func (r *PostgresStorage) InsertMany(ctx context.Context, elem interface{}) error {
// 	currentAccount := appcontext.CurrentAccount(ctx)
// 	db := r.db
// 	tx, ok := TxFromContext(ctx)
// 	if ok {
// 		db = tx
// 	}

// 	sqlStr := fmt.Sprintf(`
// 	INSERT INTO "%s"(%s)
// 	VALUES `, r.tableName, r.insertFields)

// 	var dbArgs map[string]interface{}

// 	datas := reflect.ValueOf(elem)

// 	insertFields := strings.Split(r.insertFields, ",")
// 	limit := 60000 / len(insertFields)
// 	indexData := 0
// 	if datas.Kind() == reflect.Slice {
// 		for i := 0; i < datas.Len(); i++ {
// 			sqlStr += fmt.Sprintf("(%s),", insertParams(r.elemType, r.isImmutable, i+1))
// 			arg := r.insertArgs(currentAccount, 0, datas.Index(i), i+1)
// 			if indexData == 0 {
// 				dbArgs = arg
// 			} else {
// 				for k, v := range arg {
// 					dbArgs[k] = v
// 				}
// 			}
// 			indexData++
// 			if indexData == limit {
// 				err := r.insertData(ctx, sqlStr, dbArgs)
// 				if err != nil {
// 					return err
// 				}

// 				indexData = 0
// 				sqlStr = fmt.Sprintf(`
// 				INSERT INTO "%s"(%s)
// 				VALUES `, r.tableName, r.insertFields)
// 				dbArgs = map[string]interface{}{}
// 			}
// 		}
// 	}

// 	if datas.Kind() == reflect.Map {
// 		for key, element := range datas.MapKeys() {
// 			sqlStr += fmt.Sprintf("(%s),", insertParams(r.elemType, r.isImmutable, key+1))
// 			arg := r.insertArgs(currentAccount, currentUserID, datas.MapIndex(element), key+1)
// 			if indexData == 0 {
// 				dbArgs = arg
// 			} else {
// 				for k, v := range arg {
// 					dbArgs[k] = v
// 				}
// 			}
// 			indexData++
// 			if indexData == limit {
// 				err := r.insertData(ctx, sqlStr, dbArgs)
// 				if err != nil {
// 					return err
// 				}

// 				indexData = 0
// 				sqlStr = fmt.Sprintf(`
// 				INSERT INTO "%s"(%s)
// 				VALUES `, r.tableName, r.insertFields)
// 				dbArgs = map[string]interface{}{}
// 			}
// 		}
// 	}

// 	sqlStr = strings.TrimSuffix(sqlStr, ",")

// 	statement, err := db.PrepareNamed(sqlStr)
// 	if err != nil {
// 		return err
// 	}
// 	defer statement.Close()

// 	_, err = statement.Exec(dbArgs)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (r *PostgresStorage) insertArgs(elem interface{}, index int) map[string]interface{} {
	res := map[string]interface{}{}

	var v reflect.Value

	if reflect.TypeOf(elem) == reflect.TypeOf(reflect.Value{}) {
		data := elem.(reflect.Value)
		v = reflect.Indirect(data)
	} else {
		v = reflect.ValueOf(elem).Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		dbTag := r.elemType.Field(i).Tag.Get("db")
		if !idTag(dbTag) && !emptyTag(dbTag) {
			var typeMapString map[string]interface{}
			var val interface{}
			if v.Field(i).Type() == reflect.TypeOf(typeMapString) {
				metadataBytes, err := json.Marshal(v.Field(i).Interface())
				if err != nil {
					val = "{}"
				} else {
					val = string(metadataBytes)
				}
			} else {
				val = v.Field(i).Interface()
			}
			res[dbTag] = val
		}
	}

	if index != 0 {
		s := strconv.Itoa(index)
		res = renamingKey(res, s)
	}

	return res
}

// RenamingKey is function for renaming key for map
func renamingKey(m map[string]interface{}, add string) map[string]interface{} {
	newMap := map[string]interface{}{}
	for k, v := range m {
		newKey := fmt.Sprint(k, add)
		newMap[newKey] = v
	}
	return newMap
}

// Update updates the element in the database.
// It will update the "updatedAt" field.
func (r *PostgresStorage) Update(ctx context.Context, elem interface{}) error {
	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}
	id := r.findID(elem)
	existingElem := reflect.New(r.elemType).Interface()
	err := r.FindByID(ctx, existingElem, id)

	if err != nil {
		return err
	}

	statement, err := db.PrepareNamed(fmt.Sprintf(`
		UPDATE "%s" SET %s WHERE "id" = :id RETURNING %s`,
		r.tableName,
		r.updateSetFields,
		r.selectFields))
	if err != nil {
		return err
	}
	defer statement.Close()

	updateArgs := r.updateArgs(existingElem, elem)
	updateArgs["id"] = id
	err = statement.Get(elem, updateArgs)
	if err != nil {
		return err
	}

	return nil
}

// it assumes the id column named "id"
func (r *PostgresStorage) findID(elem interface{}) interface{} {
	v := reflect.ValueOf(elem).Elem()
	for i := 0; i < v.NumField(); i++ {
		dbTag := r.elemType.Field(i).Tag.Get("db")
		if idTag(dbTag) {
			return v.Field(i).Interface()
		}
	}
	return nil
}

func (r *PostgresStorage) updateArgs(existingElem interface{}, elem interface{}) map[string]interface{} {
	res := map[string]interface{}{}

	v := reflect.ValueOf(elem).Elem()
	ev := reflect.ValueOf(existingElem).Elem()
	for i := 0; i < ev.NumField(); i++ {
		dbTag := r.elemType.Field(i).Tag.Get("db")
		if !idTag(dbTag) && !emptyTag(dbTag) {
			var typeMapString map[string]interface{}
			var val interface{}

			if v.Field(i).Type() == reflect.TypeOf(typeMapString) {
				metadataBytes, err := json.Marshal(v.Field(i).Interface())
				if err != nil {
					val = "{}"
				} else {
					val = string(metadataBytes)
				}
			} else {
				val = v.Field(i).Interface()
			}
			res[dbTag] = val
		}
	}
	return res
}

// Delete deletes the elem from database.
// Delete not really deletes the elem from the db, but it will set the
// "deletedAt" column to current time.
func (r *PostgresStorage) Delete(ctx context.Context, id interface{}) error {
	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}
	statement, err := db.PrepareNamed(fmt.Sprintf(`UPDATE "%s" SET "deleted_at" = :deletedAt WHERE "id" = :id RETURNING %s
	`, r.tableName, r.selectFields))
	if err != nil {
		return err
	}
	defer statement.Close()

	deleteArgs := map[string]interface{}{
		"id":        id,
		"deletedAt": utils.Now(),
	}

	_, err = statement.Exec(deleteArgs)
	if err != nil {
		return err
	}
	return nil
}

// DeleteHard hard delete the elem from database.
func (r *PostgresStorage) DeleteHard(ctx context.Context, id interface{}) error {
	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}

	statement, err := db.PrepareNamed(fmt.Sprintf(`
		DELETE FROM "%s" WHERE "id" = :id
	`, r.tableName))
	if err != nil {
		return err
	}
	defer statement.Close()

	deleteArgs := map[string]interface{}{
		"id": id,
	}

	_, err = statement.Exec(deleteArgs)
	if err != nil {
		return err
	}
	return nil
}

// NewPostgresStorage creates a new generic postgres Storage
func NewPostgresStorage(db *sqlx.DB, tableName string, elem interface{}) *PostgresStorage {
	elemType := reflect.TypeOf(elem)
	return &PostgresStorage{
		db:              db,
		tableName:       tableName,
		elemType:        elemType,
		selectFields:    selectFields(elemType),
		insertFields:    insertFields(elemType),
		insertParams:    insertParams(elemType, 0),
		updateSetFields: updateSetFields(elemType),
	}
}

func selectFields(elemType reflect.Type) string {
	dbFields := []string{}
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag != "" && dbTag != "-" {
			dbFields = append(dbFields, fmt.Sprintf("\"%s\"", dbTag))
		}
	}
	return strings.Join(dbFields, ",")
}

func insertFields(elemType reflect.Type) string {
	dbFields := []string{}
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		dbTag := field.Tag.Get("db")
		if !idTag(dbTag) && !emptyTag(dbTag) {
			dbFields = append(dbFields, fmt.Sprintf("\"%s\"", dbTag))
		}
	}
	return strings.Join(dbFields, ",")
}

func insertParams(elemType reflect.Type, index int) string {
	dbParams := []string{}
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		dbTag := field.Tag.Get("db")
		if !idTag(dbTag) && !emptyTag(dbTag) {
			dbParams = append(dbParams, fmt.Sprintf(":%s", dbTag))
		}
	}

	if index != 0 {
		s := strconv.Itoa(index)
		for i, v := range dbParams {
			dbParams[i] = fmt.Sprint(v, s)
		}
	}

	return strings.Join(dbParams, ",")
}

func updateSetFields(elemType reflect.Type) string {
	setFields := []string{}
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		dbTag := field.Tag.Get("db")
		if !idTag(dbTag) && !emptyTag(dbTag) {
			setFields = append(setFields, fmt.Sprintf("\"%s\" = :%s", dbTag, dbTag))
		}
	}
	return strings.Join(setFields, ",")
}

func idTag(dbTag string) bool {
	return dbTag == "id"
}

func emptyTag(dbTag string) bool {
	emptyTags := []string{"", "-"}
	for _, t := range emptyTags {
		if dbTag == t {
			return true
		}
	}
	return false
}
