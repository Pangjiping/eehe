//go:build !windows
// +build !windows

// Package distributed NOTICE: local distributed service is not supported in WINDOWS system.
package distributed

import (
	"errors"
	"github.com/Pangjiping/eehe/framework"
	"github.com/Pangjiping/eehe/framework/contract"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// LocalDistributedService represents local distributed engine of eehe.
type LocalDistributedService struct {
	container framework.Container
}

// NewLocalDistributedService will create a local distributed service object.
func NewLocalDistributedService(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("params error")
	}

	container := params[0].(framework.Container)
	return &LocalDistributedService{
		container: container}, nil
}

// Select will select an app instance to provide service.
// It's a locally distributed selector by using local file lock.
func (s LocalDistributedService) Select(serviceName string, appID string, holdTime time.Duration) (string, error) {
	appService := s.container.MustMake(contract.AppKey).(contract.App)
	runtimeFolder := appService.RuntimeFolder()
	lockFile := filepath.Join(runtimeFolder, "distribute_"+serviceName)

	// Create file mutex lock.
	lock, err := os.OpenFile(lockFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}

	// Try to preempt this file lock.
	err = syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)

	// Return existed AppID if this lock has been preempted.
	if err != nil {
		selectAppIDByte, err := ioutil.ReadAll(lock)
		if err != nil {
			return "", err
		}
		return string(selectAppIDByte), err
	}

	// If an instance preempts the lock, it will hold the lock for a period of time.
	// During this time, other instances are not allowed to preempt.
	go func() {
		defer func() {
			// Release file lock.
			syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
			// Release file resource.
			lock.Close()
			// Delete corresponded file.
			os.Remove(lockFile)
		}()

		// Select timer.
		timer := time.NewTimer(holdTime)
		<-timer.C
	}()

	// Got this file lock!
	// Write current AppID to file.
	if _, err := lock.WriteString(appID); err != nil {
		return "", err
	}

	return appID, nil
}
