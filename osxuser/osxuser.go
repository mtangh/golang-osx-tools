/* osxuser.go */

package osxuser

import (
	"fmt"
	"os/user"
	"strconv"
	"strings"
)

// OSXUser ...
type OSXUser struct {
	Name          string
	Password      string
	UID           int
	Groups        []string
	Fullname      string
	HomeDirectory string
	Shell         string
	IsHidden      bool
}

// Lookup ...
func Lookup(name string) (osxuser *OSXUser, err error) {
	// Lookup
	u, err := user.Lookup(name)
	// No found
	if err != nil {
		return nil, err
	}
	// Found, setup OSXUser
	osxuser.Name = u.Username
	osxuser.Password = "*"
	// UID
	if uid, err := strconv.Atoi(u.Uid); err == nil {
		osxuser.UID = uid
	} else {
		osxuser.UID = -1
	}
	// Groups
	osxuser.Groups = []string{}
	//
	if groups, err := u.GroupIds(); err == nil {
		for _, gid := range groups {
			if group, err := user.LookupGroupId(gid); err == nil {
				osxuser.Groups = append(osxuser.Groups, group.Name)
			}
		}
	}
	//
	osxuser.Fullname = u.Name
	osxuser.HomeDirectory = u.HomeDir
	osxuser.Shell = ""
	osxuser.IsHidden = false
	// end
	return osxuser, nil
}

// NewFromString ...
func NewFromString(entry string) (osxuser *OSXUser) {
	// Split
	fields := strings.Split(strings.TrimSpace(entry), ":")
	// No fields
	if len(fields) <= 0 {
		return nil
	}
	// User info
	osxuser = new(OSXUser)
	osxuser.UID = -1
	// Field Map
	fieldsMap := [](interface{}){
		&osxuser.Name, &osxuser.Password, &osxuser.UID, &osxuser.Groups,
		&osxuser.Fullname, &osxuser.HomeDirectory, &osxuser.Shell,
		&osxuser.IsHidden}
	// For fields
	for index, field := range fields {
		//
		field = strings.TrimSpace(field)
		//
		if len(field) <= 0 || field == "*" {
			continue
		}
		//
		switch p := (fieldsMap[index]).(type) {
		case (*string):
			*p = field
		case (*[]string):
			//
			*p = []string{}
			//
			for _, value := range strings.Split(field, ",") {
				if v := strings.TrimSpace(value); len(v) > 0 && v != "*" {
					*p = append(*p, v)
				}
			}
		case (*int):
			if intValue, err := strconv.Atoi(field); err == nil {
				*p = intValue
			} else {
				*p = -1
			}
		case (*[]int):
			//
			*p = []int{}
			//
			for _, value := range strings.Split(field, ",") {
				if intValue, err := strconv.Atoi(value); err == nil {
					*p = append(*p, intValue)
				}
			}
		case (*bool):
			*p = strings.ToLower(field) == "true" ||
				strings.ToLower(field) == "yes"
		}
	}
	// Set defaults
	if len(osxuser.Name) > 0 {
		if len(osxuser.Groups) <= 0 {
			osxuser.Groups = []string{"staff"}
		}
		if len(osxuser.Fullname) <= 0 {
			osxuser.Fullname = osxuser.Name
		}
		if len(osxuser.HomeDirectory) <= 0 {
			osxuser.HomeDirectory = fmt.Sprintf("/Users/%s", osxuser.Name)
		}
		if len(osxuser.Shell) <= 0 {
			osxuser.Shell = "/bin/bash"
		}
	} else {
		osxuser = nil
	}
	// end
	return osxuser
}

// UIDFor ...
func (osxuser *OSXUser) UIDFor(stringValue string) int {
	if osxuser == nil {
		return -1
	}
	//
	if len(strings.TrimSpace(stringValue)) > 0 {
		if uid, err := strconv.Atoi(stringValue); err == nil {
			osxuser.UID = uid
		} else {
			osxuser.UID = -1
		}
	}
	// end
	return osxuser.UID
}

// GroupsFor ...
func (osxuser *OSXUser) GroupsFor(stringValue string) []string {
	if osxuser == nil {
		return nil
	}
	//
	if len(strings.TrimSpace(stringValue)) > 0 {
		//
		groups := []string{}
		//
		for _, group := range strings.Split(
			strings.TrimSpace(stringValue), ",") {
			if g := strings.TrimSpace(group); len(g) > 0 && g != "*" {
				groups = append(groups, g)
			}
		}
		//
		osxuser.Groups = groups
	}
	// end
	return osxuser.Groups
}

// Exists ...
func (osxuser *OSXUser) Exists() (exists bool) {
	if osxuser == nil {
		return false
	}
	//
	exists = false
	//
	if _, err := user.Lookup(osxuser.Name); err == nil {
		exists = true
	} else if osxuser.UID >= 0 {
		if _, err := user.LookupId(fmt.Sprintf("%d", osxuser.UID)); err == nil {
			exists = true
		}
	}
	// end
	return exists
}
