/* osxuser.go */

package osxuser

import (
	"fmt"
	"os/user"
	"strconv"
	"strings"
)

type OSXUser struct {
	Name          string
	UID           int
	Groups        []string
	Password      string
	Fullname      string
	HomeDirectory string
	Shell         string
	IsHidden      bool
}

func NewFromString(entry string) (user *OSXUser) {
	// User info
	user = new(OSXUser)
	user.UID = -1
	// Split
	fields := strings.Split(strings.TrimSpace(entry), ":")
	// Field Map
	fieldsMap := []*string{
		&user.Name, &user.Password, nil, nil,
		&user.Fullname, &user.HomeDirectory, &user.Shell, nil}
	// For fields
	for index, field := range fields {
		//
		field = strings.TrimSpace(field)
		//
		if len(field) <= 0 || field == "*" {
			continue
		}
		//
		switch pField := fieldsMap[index]; true {
		case pField != nil:
			*pField = field
		case index == 2:
			if uid, err := strconv.Atoi(field); err == nil {
				user.UID = uid
			} else {
				user.UID = -1
			}
		case index == 3:
			//
			groups := []string{}
			//
			for _, group := range strings.Split(field, ",") {
				group = strings.TrimSpace(group)
				if g := strings.TrimSpace(group); len(g) > 0 && g != "*" {
					groups = append(groups, g)
				}
			}
			//
			user.Groups = groups
		case index == 7:
			user.IsHidden = strings.ToLower(strings.TrimSpace(field)) == "yes"
		}
	}
	// Set defaults
	if len(user.Name) > 0 {
		if len(user.Groups) <= 0 {
			user.Groups = []string{"staff"}
		}
		if len(user.Fullname) <= 0 {
			user.Fullname = user.Name
		}
		if len(user.HomeDirectory) <= 0 {
			user.HomeDirectory = fmt.Sprintf("/Users/%s", user.Name)
		}
		if len(user.Shell) <= 0 {
			user.Shell = "/bin/bash"
		}
	} else {
		user = nil
	}
	// end
	return user
}

func (self *OSXUser) UidFor(stringValue string) int {
	if len(strings.TrimSpace(stringValue)) > 0 {
		if uid, err := strconv.Atoi(stringValue); err == nil {
			self.UID = uid
		} else {
			self.UID = -1
		}
	}
	return self.UID
}

func (self *OSXUser) GroupsFor(stringValue string) []string {
	if len(strings.TrimSpace(stringValue)) > 0 {
		//
		groups := []string{}
		//
		for _, group := range strings.Split(
			strings.TrimSpace(stringValue), ",") {
			group = strings.TrimSpace(group)
			if g := strings.TrimSpace(group); len(g) > 0 && g != "*" {
				groups = append(groups, g)
			}
		}
		//
		self.Groups = groups
	}
	return self.Groups
}

func (self *OSXUser) Exists() (exists bool) {
	exists = false
	if _, err := user.Lookup(self.Name); err == nil {
		exists = true
	} else if self.UID >= 0 {
		if _, err := user.LookupId(fmt.Sprintf("%d", self.UID)); err == nil {
			exists = true
		}
	}
	return exists
}
