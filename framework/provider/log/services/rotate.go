package services

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Pangjiping/eehe/framework"
	"github.com/Pangjiping/eehe/framework/contract"
	"github.com/Pangjiping/eehe/framework/util"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
)

type EeheRotateLog struct {
	EeheLog
	folder string // log file folder
	file   string // log file location
}

func NewEeheRotateLog(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)

	appService := c.MustMake(contract.AppKey).(contract.App)
	configService := c.MustMake(contract.ConfigKey).(contract.Config)

	folder := appService.LogFolder()
	if configService.IsExist("log.folder") {
		folder = configService.GetString("log.folder")
	}

	if !util.FileExist(folder) {
		os.MkdirAll(folder, os.ModePerm)
	}

	file := "eehe.log"
	if configService.IsExist("log.file") {
		file = configService.GetString("log.file")
	}

	dateFormat := "%Y%m%d%H"
	if configService.IsExist("log.date_format") {
		dateFormat = configService.GetString("log.date_format")
	}

	linkName := rotatelogs.WithLinkName(filepath.Join(folder, file))
	options := []rotatelogs.Option{linkName}

	if configService.IsExist("log.rotate_count") {
		rotateCount := configService.GetInt("log.rotate_count")
		options = append(options, rotatelogs.WithRotationCount(uint(rotateCount)))
	}

	if configService.IsExist("log.rotate_size") {
		rotateSize := configService.GetInt("log.rotate_size")
		options = append(options, rotatelogs.WithRotationSize(int64(rotateSize)))
	}

	if configService.IsExist("log.max_age") {
		if maxAgeParse, err := time.ParseDuration(configService.GetString("log.max_age")); err == nil {
			options = append(options, rotatelogs.WithMaxAge(maxAgeParse))
		}
	}

	if configService.IsExist("log.rotate_time") {
		if rotateTimeParse, err := time.ParseDuration(configService.GetString("log.rotate_time")); err == nil {
			options = append(options, rotatelogs.WithRotationTime(rotateTimeParse))
		}
	}

	log := &EeheRotateLog{}
	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)
	log.folder = folder
	log.file = file

	w, err := rotatelogs.New(fmt.Sprintf("%s.%s", filepath.Join(log.folder, log.file), dateFormat), options...)
	if err != nil {
		return nil, errors.Wrap(err, "new rotatelog error")
	}

	log.SetOutput(w)
	log.c = c
	return log, nil
}
