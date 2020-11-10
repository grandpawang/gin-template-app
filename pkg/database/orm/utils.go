package orm

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/jinzhu/gorm"
)

// DataMap data key map
type DataMap = map[string]interface{}

// AssociationID table association id
type AssociationID = map[string][]uint

// Detail KV
type Detail struct {
	ID     uint   `json:"id" form:"id"`
	Common string `json:"common" form:"common"`
}

// Base Basic database data
type Base struct {
	ID        uint       `gorm:"type:serial;primary_key;auto_increment;"`
	CreatedAt *time.Time `json:"created_at" time_format:"0000-00-00 00:00:00"`
	UpdatedAt *time.Time `json:"updated_at" time_format:"0000-00-00 00:00:00"`
}

// BaseView basic database data view
type BaseView struct {
	ID        uint       `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// DBOrder Database query command
type DBOrder struct {
	Type  string
	Order interface{}
	Value []interface{}
}

// NumberQuery Fuzzy query number type
type NumberQuery struct {
	Max int `form:"max" json:"max"`
	Min int `form:"min" json:"min"`
}

// PageQuery page query type
type PageQuery struct {
	Page  *uint64 `form:"page" json:"page"`
	Limit *uint64 `form:"limit" json:"limit"`
	Sort  *string `form:"sort" json:"sort"`
}

type timeQuery struct {
	Start *time.Time `form:"start" json:"start"`
	End   *time.Time `form:"end" json:"end"`
}

// TimeQuery time query
type TimeQuery struct {
	CreatedAt *timeQuery `form:"created_at" json:"created_at"`
	UpdateAt  *timeQuery `form:"update_at" json:"update_at"`
}

// GetByKey interface{} get item by it's Field Name
func getByKey(form interface{}, fieldName string) reflect.Value {
	val := reflect.ValueOf(form)
	if val.Kind() != reflect.Ptr {
		val = val.FieldByName(fieldName)
	} else {
		val = val.Elem().FieldByName(fieldName)
	}
	return val
}

// Form2PageWhereOrderList whereForm to PageWhereOrder List
// whereForm is a fuzzy query and form type is ptr and form item also is ptr
// example: dbOrder := Form2PageWhereOrderList(dbOrder, whereForm, model.TableName())
func Form2PageWhereOrderList(res []DBOrder, whereForm interface{}, modelName string) []DBOrder {
	formVal := reflect.ValueOf(whereForm)
	formType := formVal.Type()

	for i := 0; i < formVal.NumField(); i++ {
		fieldVal := formVal.Field(i).Elem()
		fieldType := formType.Field(i)
		if fieldKind := fieldVal.Kind(); fieldKind != reflect.Invalid {
			fieldVal := fieldVal.Interface()
			// Custom type does not convert
			notCustoType := fieldType.Type.Elem().String() == fieldKind.String()
			// form name
			formName := fieldType.Tag.Get("form")
			if fieldKind == reflect.String && notCustoType { // reflect.String
				// string type
				res = append(res, DBOrder{
					Type:  "Where",
					Order: modelName + ".\"" + formName + "\" like ? ",
					Value: []interface{}{"%" + fieldVal.(string) + "%"},
				})
			} else if uint(fieldKind) >= 2 && uint(fieldKind) <= 14 && notCustoType { // reflect.Int = 2 reflect.Float64 = 14
				// number type
				res = append(res, DBOrder{
					Type:  "Where",
					Order: modelName + ".\"" + formName + "\" > ? and \"" + formName + "\" < ?",
					Value: []interface{}{
						getByKey(fieldVal, "Min").Interface(),
						getByKey(fieldVal, "Max").Interface(),
					},
				})
			} else if fieldKind == reflect.Array || fieldKind == reflect.Slice { // many 2 many

			} else {
				res = append(res, DBOrder{
					// else type
					Type:  "Where",
					Order: modelName + ".\"" + formName + "\" = ?",
					Value: []interface{}{fieldVal},
				})
			}
		}
	}
	return res
}

// ValidPageForm valid page query form
func ValidPageForm(pageForm *PageQuery) {
	// PageQuery to PageWhereOrder
	if pageForm.Page == nil || *pageForm.Page < 1 {
		*pageForm.Page = 1
	}
	if pageForm.Limit == nil {
		*pageForm.Limit = 10
	}
	if pageForm.Sort != nil {
		sort := *pageForm.Sort
		if len(sort) > 2 {
			orderType := sort[0:1]
			order := sort[1:]
			if orderType == "+" {
				order += " ASC"
			} else {
				order += " DESC"
			}
			*pageForm.Sort = order
		}
	} else {
		pageForm.Sort = new(string)
		*pageForm.Sort = "id ASC"

	}
}

// ValidTimeForm valid time query form
func ValidTimeForm(timeForm *TimeQuery) {
	fmt.Println(timeForm)
	fmt.Println("asdasdasda")
}

// MakeDBOrder Generate database query commands using pageForm and whereForm
func MakeDBOrder(pageForm PageQuery, timeForm TimeQuery, whereForm interface{}, TableName string) (dbOrder []DBOrder) {
	// valid page form
	ValidPageForm(&pageForm)
	dbOrder = []DBOrder{
		{Type: "Page", Order: *pageForm.Page},
		{Type: "PageSize", Order: *pageForm.Limit},
		{Type: "Order", Order: TableName + "." + *pageForm.Sort},
	}
	// add time filter
	if timeForm.CreatedAt.Start != nil && timeForm.CreatedAt.End != nil {
		dbOrder = append(dbOrder, DBOrder{
			Type:  "Where",
			Order: TableName + `."created_at" between ? and ?`,
			Value: []interface{}{timeForm.CreatedAt.Start, timeForm.CreatedAt.End},
		})
	}

	if timeForm.UpdateAt.Start != nil && timeForm.UpdateAt.End != nil {
		dbOrder = append(dbOrder, DBOrder{
			Type:  "Where",
			Order: TableName + `."update_at" between ? and ?`,
			Value: []interface{}{timeForm.UpdateAt.Start, timeForm.UpdateAt.End},
		})
	}
	dbOrder = Form2PageWhereOrderList(dbOrder, whereForm, TableName)
	return
}

// WhereArgs where args
type WhereArgs []interface{}

// WhereForm where form
type WhereForm struct {
	Join  string
	Where interface{}
	Value []interface{} // Nested subquery: DB.(...).QueryExpr()
}

// JoinWhere join where string autofill ids
type JoinWhere struct {
	Join  string
	Where string
}

// ParseWhereForm parse where form
func (orm *ORM) ParseWhereForm(whereForm []WhereForm) *gorm.DB {
	db := orm.DB
	for _, where := range whereForm {
		if where.Join != "" {
			db = db.Joins(where.Join)
		}
		db = db.Where(where.Where, where.Value...)
	}
	return db
}

// MakeWhereForm Generate a Where form for a single WhereForm
func (orm *ORM) MakeWhereForm(where interface{}, Value ...interface{}) WhereForm {
	return WhereForm{Where: where, Value: Value}
}

// SingleWhereForm Generate a Where form list for a single WhereForm
func (orm *ORM) SingleWhereForm(where interface{}, Value ...interface{}) []WhereForm {
	return []WhereForm{{Where: where, Value: Value}}
}

func (orm *ORM) appendIDsWhereForm(whereForm []WhereForm, ids *[]uint) []WhereForm {
	if ids != nil && len(*ids) > 0 {
		whereForm = append(whereForm, orm.MakeWhereForm("id in (?)", *ids))
	}
	return whereForm
}

// MakeJoinsIDsWhere use Joins method make association where.
// return common.MakeJoinsIDsWhere(ids, where, user.TableName(), map[string]common.JoinWhere{
// 	"role": common.JoinWhere{
// 		Join:  fmt.Sprintf(`join %s ur on ur.user_id = %s.id`, user.tableAssociationRole(), user.TableName()),
// 		Where: `ur.role_id in (?)`,
// 	},
// })
func (orm *ORM) MakeJoinsIDsWhere(ids *[]uint, where map[string]interface{}, tableName string, associationData map[string]JoinWhere) []WhereForm {
	whereForm := []WhereForm{}
	if where != nil && len(where) > 0 {
		for k, v := range associationData {
			if where[k] != nil {
				if ids := where[k].([]uint); len(ids) > 0 {
					whereForm = append(whereForm, WhereForm{
						Join:  v.Join,
						Where: v.Where,
						Value: WhereArgs{ids},
					})
				}
				delete(where, k)
			}
		}
		whereForm = append(whereForm, orm.MakeWhereForm(where))
	}
	whereForm = orm.appendIDsWhereForm(whereForm, ids)
	return whereForm
}

// IDWhereWhereForm Generate a Where form list for id or where
func (orm *ORM) IDWhereWhereForm(ids *[]uint, where map[string]interface{}) []WhereForm {
	whereForm := []WhereForm{}
	if len(where) > 0 {
		whereForm = orm.SingleWhereForm(where)
	}
	whereForm = orm.appendIDsWhereForm(whereForm, ids)
	return whereForm
}

// WhereForm2IDs where form generate ids
func (orm *ORM) WhereForm2IDs(ctx context.Context, modelName string, whereForm []WhereForm, ids *[]uint) (err error) {
	// exec query
	err = orm.PluckList(ctx, modelName, "ID", ids, whereForm)
	return
}

// Create a piece of data in the database
func (orm *ORM) Create(ctx context.Context, model interface{}) error {
	return orm.DB.Create(model).Error
}

// Exists or not in db by table name
func (orm *ORM) Exists(ctx context.Context, model string, whereForm []WhereForm, args ...interface{}) bool {
	var count int64

	db := orm.ParseWhereForm(whereForm).Table(model).Count(&count)
	if db.Error == nil && count >= 1 {
		return true
	}
	fmt.Println("common.Exists: ", db.Error, count)
	return false
}

// Single table operation

func (orm *ORM) getCountAndErr(db *gorm.DB) (count int64, err error) {
	err = db.Error
	if err != nil {
		return
	}
	count = db.RowsAffected
	return
}

// Update multiple pieces of data in the database by table name
func (orm *ORM) Update(ctx context.Context, model string, value interface{}, whereForm []WhereForm) (count int64, err error) {
	db := orm.ParseWhereForm(whereForm).Table(model).Updates(value)
	return orm.getCountAndErr(db)
}

// UpdateColums multiple pieces of data in the database without hooks
func (orm *ORM) UpdateColums(ctx context.Context, model string, value interface{}, whereForm []WhereForm) (count int64, err error) {
	db := orm.ParseWhereForm(whereForm).Table(model).UpdateColumns(value) // this method without hooks
	return orm.getCountAndErr(db)
}

// PhysicalDeleted multiple pieces of data in the database
func (orm *ORM) PhysicalDeleted(ctx context.Context, model interface{}, whereForm []WhereForm) (count int64, err error) {
	db := orm.ParseWhereForm(whereForm).Unscoped().Delete(model)
	return orm.getCountAndErr(db)

}

// Select operate

// GetPage from the database by table name
// like common.GetPage(tag.TableName(), &viewList, whereOrder...)
func (orm *ORM) GetPage(ctx context.Context, model string, out interface{}, whereOrder ...DBOrder) (totalCount int64, err error) {
	db := orm.DB.Table(model)
	var page uint64 = 1
	var pageSize uint64 = 10

	if len(whereOrder) > 0 {
		for _, wo := range whereOrder {
			switch wo.Type {
			case "Select":
				db = db.Select(wo.Order)
			case "Where":
				db = db.Where(wo.Order, wo.Value...)
			case "Order":
				db = db.Order(wo.Order)
			case "Joins":
				db = db.Joins(wo.Order.(string))
			case "Group":
				db = db.Group(wo.Order.(string))
			case "Page":
				page = wo.Order.(uint64)
			case "PageSize":
				pageSize = wo.Order.(uint64)
			}
		}
	}

	err = db.Count(&totalCount).Error
	if err != nil || totalCount == 0 {
		return
	}
	err = db.Offset((page - 1) * pageSize).Limit(pageSize).Scan(out).Error
	return
}

// FirstByID Get a piece of data from the database by id
func (orm *ORM) FirstByID(ctx context.Context, out interface{}, id int) (notFound bool, err error) {
	err = orm.DB.First(out, id).Error
	if err != nil {
		notFound = gorm.IsRecordNotFoundError(err)
	}
	return
}

// First Get a piece of data from the database by where form
func (orm *ORM) First(ctx context.Context, out interface{}, whereForm []WhereForm) (notFound bool, err error) {
	err = orm.ParseWhereForm(whereForm).First(out).Error
	if err != nil || err == gorm.ErrRecordNotFound {
		notFound = gorm.IsRecordNotFoundError(err)
		err = nil
	}
	return
}

// Find Get multiple pieces of data from the database and Structured data as `out` data struct
func (orm *ORM) Find(ctx context.Context, whereForm []WhereForm, out interface{}, orders ...string) error {
	db := orm.ParseWhereForm(whereForm)
	if len(orders) > 0 {
		for _, order := range orders {
			db = db.Order(order)
		}
	}
	return db.Find(out).Error
}

// Scan Get multiple pieces of data from the database and data struct is `out` data struct
func (orm *ORM) Scan(ctx context.Context, model, out interface{}, whereForm []WhereForm) (notFound bool, err error) {
	err = orm.ParseWhereForm(whereForm).Model(model).Scan(out).Error
	if err != nil {
		notFound = gorm.IsRecordNotFoundError(err)
	}
	return
}

// ScanTable Get multiple pieces of data from the database table name and data struct is `out` data struct
func (orm *ORM) ScanTable(ctx context.Context, model string, out interface{}, whereForm []WhereForm) (err error) {
	return orm.ParseWhereForm(whereForm).Table(model).Scan(out).Error
}

// PluckList pluck multiple pieces of data from the database table name.
func (orm *ORM) PluckList(ctx context.Context, model, fieldName string, out interface{}, whereForm []WhereForm) error {
	return orm.ParseWhereForm(whereForm).Table(model).Pluck(fieldName, out).Error
}

// CheckIDs Verify that the ids exists and return the ids that exists
func (orm *ORM) CheckIDs(ctx context.Context, tableName string, ids *[]uint) (err error) {
	if err = orm.PluckList(ctx, tableName, "ID", ids, orm.SingleWhereForm("id in (?)", *ids)); err != nil {
		err = fmt.Errorf("validate role ids error")
		return
	}
	return
}

// HandleAssociationFunc Handling with database association data
type HandleAssociationFunc func(ctx context.Context, model interface{}, ids *[]uint, column string, associationData interface{}) error

// HandleSelfFunc Handling with Table itself data
type HandleSelfFunc func(ctx context.Context, modelName string, model interface{}, ids *[]uint, form map[string]interface{}) (int64, error)

// HandleFunc assciation data handle methods
type HandleFunc struct {
	HandleAssociationFunc
	HandleSelfFunc
}

// HandleAssociation simple handle association data methods.
func (orm *ORM) HandleAssociation(
	ctx context.Context,
	modelName string, model interface{},
	ids *[]uint, form, association map[string]interface{}, handle HandleFunc) (count int64, err error) {

	count, err = 0, nil
	tx := orm.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err = tx.Error; err != nil {
		tx.Rollback()
		return
	}

	if len(*ids) <= 0 {
		tx.Rollback()
		return
	}

	for column, associationData := range association {
		// Handling associations
		if err = handle.HandleAssociationFunc(ctx, model, ids, column, associationData); err != nil {
			tx.Rollback()
			return
		}
	}

	// Handle Self
	count, err = handle.HandleSelfFunc(ctx, modelName, model, ids, form)
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit().Error
	return
}

// GetAssociationParentIDs get parent ids, and changes in parents, sons must change
// example: common.GetAssociationParentIDs(r.TableName(), ids)
func (orm *ORM) GetAssociationParentIDs(ctx context.Context, modelName string, ids *[]uint) (err error) {
	err = orm.DB.Raw(fmt.Sprintf(`
		WITH RECURSIVE res as (
			select br.id from %s br where br.id in (?)
			UNION
			select br2.id from %s br2 INNER JOIN res s ON s.id = br2.parent_id
		) SELECT id FROM res;`, modelName, modelName), *ids).Pluck("id", ids).Error
	return
}

// GetCascadeAssociationIDs Get data without cascading relationship, and control them.
// example:  err = common.GetCascadeAssociationIDs(ids, "blog_user_role", "role_id", &msg);
func (orm *ORM) GetCascadeAssociationIDs(ctx context.Context, ids *[]uint, tableName string, whoIDs string, msg *string) (err error) {
	var count []uint
	var ctlRoleID []uint
	var passRoleID []uint
	for _, id := range *ids {
		if err = orm.DB.Table(tableName).Select("count(1) as count").Where(fmt.Sprintf("%s = %d", whoIDs, id)).Pluck("count", &count).Error; err != nil {
			return
		}
		if count[0] == 0 {
			ctlRoleID = append(ctlRoleID, id)
		} else {
			passRoleID = append(passRoleID, id)
		}
	}
	*ids = ctlRoleID
	if len(passRoleID) > 0 {
		*msg += fmt.Sprintf("Role id %v can't delete, It had user used.", passRoleID)
	}
	return
}

// AssociationDelete delete asssociation
func (orm *ORM) AssociationDelete(ctx context.Context, ids *[]uint, associationModel, whatIDs string) (err error) {
	if ids != nil && len(*ids) > 0 {
		if err = orm.DB.Exec(fmt.Sprintf(`DELETE FROM %s WHERE (%s IN (?))`, associationModel, whatIDs), *ids).Error; err != nil {
			return
		}
	}
	return
}

// SimpleHandleAssociation simple handle association function
func (orm *ORM) SimpleHandleAssociation(model interface{}, ids *[]uint, action func(interface{}) *gorm.Association) (err error) {
	modelVal := reflect.ValueOf(model).Elem().FieldByName("ID")
	for _, id := range *ids {
		modelVal.Set(reflect.ValueOf(id))
		if err = action(model).Error; err != nil {
			return
		}
	}
	modelVal.Set(reflect.ValueOf(uint(0))) // set nil id
	return
}

// SimpleHanleSelfUpdate simple handle self update
func (orm *ORM) SimpleHanleSelfUpdate(ctx context.Context, modelName string, model interface{}, ids *[]uint, form map[string]interface{}) (count int64, err error) {
	count, err = 0, nil
	if ids != nil && len(*ids) > 0 {
		count, err = orm.Update(ctx, modelName, form, orm.SingleWhereForm("id in (?)", *ids))
	}
	return
}

// SimpleHanleSelfDelete simple handle self Delete
func (orm *ORM) SimpleHanleSelfDelete(ctx context.Context, modelName string, model interface{}, ids *[]uint, form map[string]interface{}) (count int64, err error) {
	count, err = 0, nil
	if ids != nil && len(*ids) > 0 {
		count, err = orm.UpdateColums(ctx, modelName, form, orm.SingleWhereForm("id in (?)", *ids))
	}
	return
}

// SimpleHandleSelfRemove simple handle self remove
func (orm *ORM) SimpleHandleSelfRemove(ctx context.Context, modelName string, model interface{}, ids *[]uint, form map[string]interface{}) (count int64, err error) {
	count, err = 0, nil
	if ids != nil && len(*ids) > 0 {
		count, err = orm.PhysicalDeleted(ctx, model, orm.SingleWhereForm("id in (?)", *ids))
	}
	return
}

// SimpleUpdateFunc get HandleAssociationFunc and HandleSelfFunc
func (orm *ORM) SimpleUpdateFunc(ctx context.Context) HandleFunc {
	return HandleFunc{
		HandleAssociationFunc: func(ctx context.Context, model interface{}, ids *[]uint, column string, associationData interface{}) (err error) {
			return orm.SimpleHandleAssociation(model, ids,
				func(tmp interface{}) *gorm.Association {
					return orm.DB.Model(tmp).Association(column).Replace(associationData)
				})
		},
		HandleSelfFunc: orm.SimpleHanleSelfUpdate,
	}
}

// SimpleDeleteFunc get HandleAssociationFunc and HandleSelfFunc
func (orm *ORM) SimpleDeleteFunc(ctx context.Context) HandleFunc {
	return HandleFunc{
		HandleAssociationFunc: func(ctx context.Context, model interface{}, ids *[]uint, column string, associationData interface{}) (err error) {
			return orm.SimpleHandleAssociation(model, ids,
				func(tmp interface{}) *gorm.Association {
					return orm.DB.Model(tmp).Association(column).Clear()
				})
		},
		HandleSelfFunc: orm.SimpleHandleSelfRemove,
	}
}

// SimpleAssociationUpdates simple association updates
func (orm *ORM) SimpleAssociationUpdates(
	ctx context.Context,
	modelName string, model interface{},
	ids *[]uint, form, association map[string]interface{}) (count int64, err error) {
	return orm.HandleAssociation(ctx, modelName, model, ids, form, association, orm.SimpleUpdateFunc(ctx))
}

// SimpleAssociationDelete simple association Delete.
func (orm *ORM) SimpleAssociationDelete(
	ctx context.Context,
	modelName string, model interface{},
	ids *[]uint, association map[string]interface{}) (count int64, err error) {
	return orm.HandleAssociation(ctx, modelName, model, ids, nil, association, orm.SimpleDeleteFunc(ctx))
}
