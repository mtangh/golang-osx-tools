/* osx-useradd */

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"./osxuser"
)

var (
	cmdname        = filepath.Base(os.Args[0])
	optIsStdinFlag = flag.Bool("stdin", false,
		"read user entry from stdin")
	optIsUsageFlag = flag.Bool("help", false,
		"Show this help message")
)

func main() {
	//
	var (
		optUID      string
		optGroups   string
		optPassword string
		optComment  string
		optHomeDir  string
		optShell    string
		optIsHidden bool
	)
	//
	flag.StringVar(&optUID, "u", "",
		"user ID of the new account")
	flag.StringVar(&optGroups, "G", "",
		"list of supplementary groups of the new account")
	flag.StringVar(&optPassword, "p", "",
		"encrypted password of the new account")
	flag.StringVar(&optComment, "c", "",
		"GECOS field of the new account")
	flag.StringVar(&optHomeDir, "d", "",
		"home directory of the new account")
	flag.StringVar(&optShell, "s", "",
		"login shell of the new account")
	flag.BoolVar(&optIsHidden, "H", false,
		"set hidden flag")
	// Parse cmd options
	flag.Parse()
	// Help ?
	if *optIsUsageFlag {
		flag.Usage()
		os.Exit(0)
	}
	// User info
	var users []*osxuser.OSXUser
	// No options,
	if *optIsStdinFlag {
		// Read from stdin
		users = readFrom(os.Stdin)
	} else if flag.NArg() >= 1 {
		//
		var user *osxuser.OSXUser
		//
		isHidden := ""
		//
		if optIsHidden {
			isHidden = "yes"
		}
		//
		user = osxuser.NewFromString(fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s:%s",
			strings.TrimSpace(flag.Arg(0)),
			optUID, optGroups, optPassword, optComment, optHomeDir, optShell,
			isHidden))
		//
		users = append(users, user)
	} else {
		flag.Usage()
	}
	// Print user info
	for _, user := range users {
		if err := createOsxUser(user); err != nil {
			os.Exit(1)
		}
	}
	// Exit
	os.Exit(0)
}

func readFrom(fp *os.File) (users []*osxuser.OSXUser) {
	// New reader
	reader := bufio.NewReaderSize(fp, 4096)
	// Read record
	for {

		// Read record
		record, err := reader.ReadString('\n')

		// EOF
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		// Trim
		record = strings.TrimSpace(record)

		// Check record
		if len(record) <= 0 {
			continue
		} else if strings.HasPrefix(record, "#") {
			continue
		}

		// New From String
		user := osxuser.NewFromString(record)
		// Check
		if user != nil {
			users = append(users, user)
		}

	}

	// End
	return users
}

func createOsxUser(user *osxuser.OSXUser) error {
	if user.Exists() {
		fmt.Fprintf(os.Stderr, "%s: User '%s' exists.\n",
			cmdname, user.Name)
	} else {
		fmt.Printf("%T: %s:%d:%s:%s:%s:%s:%s:%v\n",
			user,
			user.Name, user.UID, user.Groups, user.Password,
			user.Fullname, user.HomeDirectory, user.Shell,
			user.IsHidden)
	}
	return nil
}
