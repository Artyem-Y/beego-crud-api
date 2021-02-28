package models

import (
	"errors"
	"fmt"
	_ "github.com/prometheus/common/log"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Users struct {
	ID                 int64     `orm:"column(id);pk;auto" json:"id"`
	Name               string    `orm:"size(128)" json:"name"`
	Email              string    `orm:"size(128);unique" json:"email"`
	Phone              string    `orm:"size(128);unique" json:"phone"`
	Password           string    `orm:"size(128)" json:"password"`
	AccessToken        string    `orm:"size(128)" json:"access_token"`
	Role               string    `orm:"size(128)" json:"role"`
	CreatedAt          time.Time `orm:"column(created_at);auto_now_add;type(timestamp with time zone);null" json:"created_at"`
	UpdatedAt          time.Time `orm:"column(updated_at);auto_now;type(timestamp with time zone);null" json:"updated_at"`
	RecentLogin        time.Time `orm:"column(recent_login);type(timestamp with time zone);null" json:"recent_login"`
	ValidationCodeSent time.Time `orm:"column(validation_code_sent);type(timestamp with time zone);null" json:"validation_code_sent"`
	EmailConfirmed     bool      `orm:"size(128)" json:"email_confirmed"`
}

func init() {
	orm.RegisterModel(new(Users))
}

// AddUsers insert a new User into database and returns
// last inserted Id on success.
func AddUsers(m *Users) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetUsersById retrieves User by Id. Returns error if
// Id doesn't exist
func GetUsersById(id int64) (v *Users, err error) {
	o := orm.NewOrm()
	v = &Users{ID: id}

	if err = o.QueryTable(new(Users)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetUsersByEmail retrieves Customer by Email. Returns error if
// Id doesn't exist
func GetUsersByEmail(email string) (v *Users, err error) {
	o := orm.NewOrm()
	v = &Users{Email: email}

	if err = o.QueryTable(new(Users)).Filter("Email", email).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetUsersByEmail retrieves Customer by Email. Returns error if
// Id doesn't exist
func GetUsersByPhoneNumber(phone string) (v *Users, err error) {
	o := orm.NewOrm()
	v = &Users{Phone: phone}

	if err = o.QueryTable(new(Users)).Filter("Phone", phone).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllUsers retrieves all Users matches certain condition. Returns empty list if
// no records exist
func GetAllUsers(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Users))

	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string

	if len(sortby) != 0 {

		if len(sortby) == len(order) {

			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""

				if order[i] == "desc" {
					orderby = "-" + v

				} else if order[i] == "asc" {
					orderby = v

				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)

		} else if len(sortby) != len(order) && len(order) == 1 {

			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""

				if order[0] == "desc" {
					orderby = "-" + v

				} else if order[0] == "asc" {
					orderby = v

				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}

		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}

	} else {

		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Users
	qs = qs.OrderBy(sortFields...).RelatedSel()
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateUsers updates Users by Id and returns error if
// the record to be updated doesn't exist
func UpdateUsersById(m *Users) (err error) {
	o := orm.NewOrm()
	v := Users{ID: m.ID}

	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64

		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUser deletes User by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUsers(id int64) (err error) {
	o := orm.NewOrm()
	v := Users{ID: id}

	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64

		if num, err = o.Delete(&Users{ID: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
