/* DsclCmd.go */

package dscl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

// Cmd ...
type Cmd struct {
	path       string
	datasource string
	user       string
	password   string
}

// New ...
func New() (dsclCmd *Cmd, err error) {
	return NewWithDatasource(".")
}

// NewWithDatasource ...
func NewWithDatasource(datasource string) (dsclCmd *Cmd, err error) {
	// Looup "dscl" path
	path := ""
	path, err = exec.LookPath("dscl")
	// Has error
	if err != nil {
		return nil, err
	}
	// Datasource
	if len(datasource) <= 0 {
		datasource = "."
	}
	// end
	return &Cmd{path, datasource, "", ""}, err
}

func dsclCreateCommand(dsclCmd *Cmd, command string, args ...string) (cmd *exec.Cmd, err error) {
	// Check
	if dsclCmd == nil || len(dsclCmd.path) <= 0 || len(command) <= 0 || len(args) <= 0 {
		return nil, os.ErrInvalid
	}
	// Command
	if strings.HasPrefix(command, "-") == false {
		command = "-" + command
	}
	// Default options
	dsclCmdArgs := []string{"-q", "-url"}
	// Authentication user
	if len(dsclCmd.user) > 0 {
		dsclCmdArgs = append(dsclCmdArgs, "-u", dsclCmd.user)
	}
	// Authentication password
	if len(dsclCmd.password) > 0 {
		dsclCmdArgs = append(dsclCmdArgs, "-P", dsclCmd.password)
	}
	// Datasource
	dsclCmdArgs = append(dsclCmdArgs, dsclCmd.datasource, command)
	// DsclCmd Args
	dsclCmdArgs = append(dsclCmdArgs, args...)
	// Create command
	cmd = exec.Command(dsclCmd.path, dsclCmdArgs...)
	// Check
	if cmd == nil {
		return nil, exec.ErrNotFound
	}
	// end
	return cmd, nil
}

func dsclExecute(dsclCmd *Cmd, command string, args ...string) (output *bytes.Buffer, err error) {
	// DSCL Command
	var cmd *exec.Cmd
	// New DsclCmd command
	if cmd, err = dsclCreateCommand(dsclCmd, command, args...); err != nil {
		return nil, err
	}
	// Stdout and Stderr pipe
	var stdout, stderr io.ReadCloser
	// Open stdout pipe
	if stdout, err = cmd.StdoutPipe(); err != nil {
		return nil, err
	}
	// Open stderr pipe
	if stderr, err = cmd.StderrPipe(); err != nil {
		return nil, err
	}
	// Bytes buffer
	var stdOutBuff, stdErrBuff bytes.Buffer
	// Start
	if err = cmd.Start(); err != nil {
		return nil, err
	}
	// WaitGroup
	var wg sync.WaitGroup
	// Register WaitGroup
	wg.Add(2)
	// Stdout
	go func() {
		// Copy
		io.Copy(&stdOutBuff, stdout)
		// Close stdout
		stdout.Close()
		// Done
		wg.Done()
		// end
		return
	}()
	// Stderr
	go func() {
		// Copy
		io.Copy(&stdErrBuff, stderr)
		// Close stdout
		stdout.Close()
		// Done
		wg.Done()
		// end
		return
	}()
	// Waiting goroutines
	wg.Wait()
	// Waiting command
	if err = cmd.Wait(); err != nil {
		// Exit Status
		exitStatus := 0
		// Type Assertion
		if ee, ok := err.(*exec.ExitError); ok {
			if es, ok := ee.Sys().(syscall.WaitStatus); ok {
				exitStatus = es.ExitStatus()
			}
		}
		// Setup Error
		err = newDsclError(exitStatus, stdErrBuff.Bytes())
		// end
		return nil, err
	}
	// Output
	output = &stdOutBuff
	// end
	return output, err
}

// Create ...
func (dsclCmd *Cmd) Create(path string) (err error) {
	if dsclCmd == nil || len(path) <= 0 {
		return syscall.EINVAL
	}
	// end
	return dsclCmd.CreateWithProperties(path, nil)
}

// CreateWithProperties ...
func (dsclCmd *Cmd) CreateWithProperties(path string, props Properties) (err error) {
	if dsclCmd == nil || len(path) <= 0 {
		return os.ErrInvalid
	}
	// New DsclCmd command
	if _, err = dsclExecute(dsclCmd, "-create", path); err != nil {
		return err
	}
	// Create properties
	for key, propValue := range props {
		// Values
		values := propValue.Strings()
		// Items
		if len(values) <= 0 {
			continue
		}
		// New DsclCmd command
		if _, err = dsclExecute(dsclCmd, "-create", path, key, values[0]); err != nil {
			continue
		}
		// Append values
		if len(values) <= 1 {
			continue
		}
		// Append values
		if err = dsclCmd.Append(path, key, &Value{values[1:]}); err != nil {
			continue
		}
	}
	// end
	return err
}

// Delete ...
func (dsclCmd *Cmd) Delete(path string, keys ...string) (props Properties, err error) {
	if dsclCmd == nil || len(path) <= 0 {
		return nil, os.ErrInvalid
	}
	// Property
	var tempProps Properties
	// Read
	if tempProps, err = dsclCmd.Read(path, keys...); err != nil {
		return nil, err
	}
	//
	if len(keys) <= 0 {
		if _, err = dsclExecute(dsclCmd, "-delete", path); err != nil {
			return nil, err
		}
		//
		props = tempProps
	} else {
		for _, key := range keys {
			//
			if len(key) <= 0 {
				continue
			}
			if v, ok := tempProps[key]; ok {
				props[key] = v
			}
			//
			if _, ok := props[key]; ok {
				if _, err = dsclExecute(dsclCmd, "-delete", path, key); err != nil {
					return nil, err
				}
			}
		}
	}
	// end
	return props, err
}

// Read ...
func (dsclCmd *Cmd) Read(path string, keys ...string) (props Properties, err error) {
	if dsclCmd == nil || len(path) <= 0 {
		return nil, os.ErrInvalid
	}
	// Output buffer
	var output *bytes.Buffer
	// Execute DSCL command
	if output, err = dsclExecute(dsclCmd, "-read", path); err != nil {
		return nil, err
	}
	// New scanner
	scanner := bufio.NewScanner(output)
	// Scanner
	if scanner != nil {
		return nil, syscall.EINVAL
	}
	// Key and Value
	key := ""
	value := ""
	// Scan
	for scanner.Scan() {
		// Scaning text
		line := scanner.Text()
		// line
		if len(line) <= 0 {
			continue
		}
		// Key and value
		if strings.HasPrefix(line, " ") && len(key) > 0 {
			value = line[1:]
		} else if index := strings.Index(line, ": "); index > 0 {
			key = line[0:index]
			value = line[index+1:]
		} else if strings.HasSuffix(line, ":") {
			key = line[0 : len(line)-1]
			value = ""
		}
		// Decode value
		decodeValue := ""
		// Decode
		if decodeValue, err = url.QueryUnescape(value); err == nil {
			value = decodeValue
		}
		//
		if len(keys) > 0 {
			//
			keysFound := false
			//
			for _, keyItem := range keys {
				if key == keyItem {
					keysFound = true
					break
				}
			}
			//
			if keysFound {
				continue
			}
		}
		//
		if pval, ok := props[key]; ok {
			//
			pval.value = pval.String() + value
			//
			props[key] = pval
		} else {
			props[key] = &Value{value}
		}
	}
	// Check error
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	// end
	return props, err
}

// Change ...
func (dsclCmd *Cmd) Change(path string, key string, value *Value) (err error) {
	if dsclCmd == nil || len(path) <= 0 {
		return syscall.EINVAL
	}
	// end
	return dsclCmd.ChangeAtIndex(path, key, value, 0)
}

// ChangeAtIndex ...
func (dsclCmd *Cmd) ChangeAtIndex(path string, key string, value *Value, index int) (err error) {
	if dsclCmd == nil || len(path) <= 0 {
		return os.ErrInvalid
	}
	// Key and value
	if len(key) <= 0 || value == nil {
		return os.ErrInvalid
	}
	// Current property
	var props Properties
	// Read property
	if props, err = dsclCmd.Read(path, key); err != nil {
		return err
	} else if props == nil {
		return syscall.EINVAL
	}
	// Arguments
	dsclArgs := append([]string{path, key}, fmt.Sprintf("%d", index), value.String())
	// Append
	if _, err = dsclExecute(dsclCmd, "-changei", dsclArgs...); err != nil {
		return err
	}
	// end
	return err
}

// Append ...
func (dsclCmd *Cmd) Append(path string, key string, value *Value) (err error) {
	if dsclCmd == nil || len(path) <= 0 {
		return os.ErrInvalid
	}
	// Key and value
	if len(key) <= 0 || value == nil {
		return os.ErrInvalid
	}
	// Current property
	var props Properties
	// Read property
	if props, err = dsclCmd.Read(path, key); err != nil {
		return err
	} else if props == nil {
		return syscall.EINVAL
	}
	// Arguments
	dsclArgs := append([]string{path, key}, (value.Strings())...)
	// Append
	if _, err = dsclExecute(dsclCmd, "-append", dsclArgs...); err != nil {
		return err
	}
	// end
	return err
}

// List ...
func (dsclCmd *Cmd) List(path string) (list []string, err error) {
	if dsclCmd == nil || len(path) <= 0 {
		return nil, os.ErrInvalid
	}
	// path
	if len(path) > 0 {
		return nil, os.ErrInvalid
	}
	// Output
	var output *bytes.Buffer
	// Execute list
	if output, err = dsclExecute(dsclCmd, "-list", path); err != nil {
		return nil, err
	}
	// Create a buffered reader
	if reader := bufio.NewReader(output); reader != nil {
		//
		for {
			//
			var line []byte
			//
			line, _, err = reader.ReadLine()
			//
			if err == io.EOF {
				break
			} else if err != nil {
				break
			} else if len(line) <= 0 {
				break
			}
			// Append list
			list = append(list, string(line))
		}
	}
	// end
	return list, err
}
