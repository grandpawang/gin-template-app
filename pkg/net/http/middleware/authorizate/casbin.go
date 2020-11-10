package authorizate

import (
	"fmt"
	"strconv"

	"github.com/casbin/casbin"
)

const (
	// PrefixUserID user id prefix
	PrefixUserID = "u"
	// PrefixRoleID role id prefix
	PrefixRoleID = "r"
	// SimplePermission simple permission
	SimplePermission = "GET|POST"
)

// Enforcer casbin Enforcer
var Enforcer *casbin.Enforcer

// RBAC role base access control
var RBAC = casbin.NewModel(`
	[request_definition]
	r = sub, obj, act

	[policy_definition]
	p = sub, obj, act

	[role_definition]
	g = _, _

	[policy_effect]
	e = some(where (p.eft == allow))

	[matchers]
	m = g(r.sub, p.sub) == true \
		&& keyMatch2(r.obj, p.obj) == true \
		&& regexMatch(r.act, p.act) == true \
		|| r.sub == "root"
	`)

// SetupCasbinEnforcer setup cashbin enforcer
func SetupCasbinEnforcer(policy map[uint][]string) (err error) {
	// If the configuration file does not exist, an error is thrown
	Enforcer, err = casbin.NewEnforcerSafe(RBAC)
	if err != nil {
		return
	}
	// add casbin role permission
	for roleid, menuAccess := range policy {
		setRolePermission(roleid, menuAccess)
	}
	fmt.Println(Enforcer.GetPolicy())
	return
}

// setRolePermission set Role
func setRolePermission(roleid uint, urls []string) {
	for _, url := range urls {
		Enforcer.AddPermissionForUser(PrefixRoleID+strconv.FormatInt(int64(roleid), 10), url, SimplePermission)
	}
}

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=permission-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=

// CsbinCheckPermission casbin check permission
func CsbinCheckPermission(userID, url, methodtype string) (bool, error) {
	return Enforcer.EnforceSafe(PrefixUserID+userID, url, methodtype)
}

// CsbinUpdatePermission casbin update permission
func CsbinUpdatePermission(ourl, url string, roleIDs []uint) {
	if Enforcer == nil {
		return
	}

	for _, roleid := range roleIDs {
		rid := PrefixRoleID + strconv.FormatInt(int64(roleid), 10)
		Enforcer.DeletePermissionForUser(rid, ourl, SimplePermission)
		Enforcer.AddPermissionForUser(rid, url, SimplePermission)
	}
	return
}

// CsbinRestorePermission casbin restore permission
func CsbinRestorePermission(res map[string][]uint) {
	if Enforcer == nil {
		return
	}
	for url, roleIDs := range res {
		for _, roleid := range roleIDs {
			rid := PrefixRoleID + strconv.FormatInt(int64(roleid), 10)
			Enforcer.AddPermissionForUser(rid, url, SimplePermission)
		}
	}
}

// CsbinDeletePermission casbin delete permiison when the menu delete
func CsbinDeletePermission(urls []string) {
	if Enforcer == nil {
		return
	}
	for _, url := range urls {
		Enforcer.DeletePermission(url, SimplePermission)
	}
}

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=role-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=

// CsbinSetAPermissionsForRole casbin set role permission
func CsbinSetAPermissionsForRole(rolesid []uint, urls []string) {
	if Enforcer == nil {
		return
	}
	for _, roleid := range rolesid {
		Enforcer.DeletePermissionsForUser(PrefixRoleID + strconv.FormatInt(int64(roleid), 10))
		setRolePermission(roleid, urls)
	}
}

// CsbinRestoreRole casbin restore role
func CsbinRestoreRole(roleid uint, userids []uint, urls []string) {
	if Enforcer == nil {
		return
	}
	rid := PrefixRoleID + strconv.FormatInt(int64(roleid), 10)
	for _, url := range urls {
		Enforcer.AddPermissionForUser(rid, url, SimplePermission)
	}
	for _, userid := range userids {
		uid := PrefixUserID + strconv.FormatInt(int64(userid), 10)
		Enforcer.AddRoleForUser(uid, rid)
	}
}

// CsbinDeleteRole casbin delete role
func CsbinDeleteRole(roleids []uint) {
	if Enforcer == nil {
		return
	}
	for _, rid := range roleids {
		Enforcer.DeleteRole(PrefixRoleID + strconv.FormatInt(int64(rid), 10))
	}
}

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=user-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=

// CsbinSetARolesForUser csbin handle role for user
func CsbinSetARolesForUser(userids []uint, roleids []uint) (err error) {
	if Enforcer == nil {
		return
	}
	for _, userid := range userids {
		uid := PrefixUserID + strconv.FormatInt(int64(userid), 10)
		for _, roleid := range roleids {
			rid := PrefixRoleID + strconv.FormatInt(int64(roleid), 10)
			Enforcer.AddRoleForUser(uid, rid)
		}
	}
	return
}

// CsbinRestoreUser casbin restore user
func CsbinRestoreUser(res map[uint][]uint) {
	if Enforcer == nil {
		return
	}
	for userid, roleids := range res {
		uid := PrefixUserID + strconv.FormatInt(int64(userid), 10)
		for _, roleid := range roleids {
			rid := PrefixRoleID + strconv.FormatInt(int64(roleid), 10)
			Enforcer.AddRoleForUser(uid, rid)
		}

	}

}

// CsbinDeleteUser csbin delete user
func CsbinDeleteUser(userids []uint) {
	if Enforcer == nil {
		return
	}
	for _, userid := range userids {
		uid := PrefixUserID + strconv.FormatInt(int64(userid), 10)
		Enforcer.DeleteUser(uid)
	}
	return
}
