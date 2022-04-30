package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Pangjiping/eehe/framework"
	"github.com/Pangjiping/eehe/framework/contract"
	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"gopkg.in/yaml.v2"
)

type EeheConfig struct {
	container  framework.Container    // 服务容器
	folder     string                 // 文件夹
	keyDelim   string                 // 路径的分隔符，默认为.
	mu         sync.RWMutex           // 配置文件读写锁
	envMaps    map[string]string      // 所有的环境变量
	configMaps map[string]interface{} // 配置文件结构，key为文件名
	configRaws map[string][]byte      // 配置文件原始信息
}

func (eehe *EeheConfig) loadConfigFile(folder string, file string) error {
	eehe.mu.Lock()
	defer eehe.mu.Unlock()

	// 判断文件后缀
	s := strings.Split(file, ".")
	if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
		name := s[0]

		// 读取文件内容
		bf, err := ioutil.ReadFile(filepath.Join(folder, file))
		if err != nil {
			return err
		}

		// 直接针对文本做环境变量的替换
		bf = replace(bf, eehe.envMaps)

		// 解析对应的文件
		c := map[string]interface{}{}
		if err := yaml.Unmarshal(bf, &c); err != nil {
			return err
		}
		eehe.configMaps[name] = c
		eehe.configRaws[name] = bf

		// 读取app.path中的信息，更新app对应的folder
		if name == "app" && eehe.container.IsBind(contract.AppKey) {
			if p, ok := c["path"]; ok {
				appService := eehe.container.MustMake(contract.AppKey).(contract.App)
				appService.LoadAppConfig(cast.ToStringMapString(p))
			}
		}
	}
	return nil
}

// removeConfigFile deletes config file.
func (eehe *EeheConfig) removeConfigFile(folder string, file string) error {
	eehe.mu.Lock()
	defer eehe.mu.Unlock()

	s := strings.Split(file, ".")

	// .yaml and .yml is expected
	if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
		name := s[0]
		delete(eehe.configRaws, name)
		delete(eehe.configMaps, name)
	}
	return nil
}

func NewEeheConfig(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	envFolder := params[1].(string)
	envMaps := params[2].(map[string]string)

	// instance
	eeheConfig := &EeheConfig{
		container:  container,
		folder:     envFolder,
		envMaps:    envMaps,
		configMaps: map[string]interface{}{},
		configRaws: map[string][]byte{},
		keyDelim:   ".",
		mu:         sync.RWMutex{},
	}

	// check folder
	if _, err := os.Stat(envFolder); os.IsNotExist(err) {
		return eeheConfig, nil
	}

	// read env files
	files, err := ioutil.ReadDir(envFolder)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, file := range files {
		fileName := file.Name()
		err := eeheConfig.loadConfigFile(envFolder, fileName)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	// watch files
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	err = watcher.Add(envFolder)
	if err != nil {
		return nil, err
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()

		for {
			select {
			case ev := <-watcher.Events:
				{
					// event type
					// create/write/remove

					path, _ := filepath.Abs(ev.Name)
					index := strings.LastIndex(path, string(os.PathSeparator))
					folder := path[:index]
					fileName := path[index+1:]

					if ev.Op&fsnotify.Create == fsnotify.Create {
						log.Println("create file: ", ev.Name)
						eeheConfig.loadConfigFile(folder, fileName)
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						log.Println("write file: ", ev.Name)
						eeheConfig.loadConfigFile(folder, fileName)
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						log.Println("remove file: ", ev.Name)
						eeheConfig.removeConfigFile(folder, fileName)
					}

				}

			case err := <-watcher.Errors:
				log.Println("error: ", err)
				return
			}
		}
	}()

	return eeheConfig, nil
}

// replace replaces env config.
func replace(content []byte, maps map[string]string) []byte {
	if maps == nil {
		return content
	}

	for key, val := range maps {
		reKey := "env(" + key + ")"
		content = bytes.ReplaceAll(content, []byte(reKey), []byte(val))
	}
	return content
}

// searchMap searches for path env config.
func searchMap(source map[string]interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}

	// has next path?
	next, ok := source[path[0]]
	if ok {
		if len(path) == 1 {
			return next
		}

		switch next.(type) {
		case map[interface{}]interface{}:
			return searchMap(cast.ToStringMap(next), path[1:])
		case map[string]interface{}:
			return searchMap(next.(map[string]interface{}), path[1:])
		default:
			return nil
		}
	}
	return nil
}

// find searches a env info by key.
func (config *EeheConfig) find(key string) interface{} {
	config.mu.Lock()
	defer config.mu.Unlock()
	return searchMap(config.configMaps, strings.Split(key, config.keyDelim))
}

func (config *EeheConfig) IsExist(key string) bool {
	return config.find(key) != nil
}

func (config *EeheConfig) Get(key string) interface{} {
	return config.find(key)
}

func (config *EeheConfig) GetBool(key string) bool {
	return cast.ToBool(config.find(key))
}

func (config *EeheConfig) GetInt(key string) int {
	return cast.ToInt(config.find(key))
}

func (config *EeheConfig) GetFloat64(key string) float64 {
	return cast.ToFloat64(config.find(key))
}

func (config *EeheConfig) GetTime(key string) time.Time {
	return cast.ToTime(config.find(key))
}

func (config *EeheConfig) GetString(key string) string {
	return cast.ToString(config.find(key))
}

func (config *EeheConfig) GetIntSlice(key string) []int {
	return cast.ToIntSlice(config.find(key))
}

func (config *EeheConfig) GetStringSlice(key string) []string {
	return cast.ToStringSlice(config.find(key))
}

func (config *EeheConfig) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(config.find(key))
}

func (config *EeheConfig) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(config.find(key))
}

func (config *EeheConfig) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(config.find(key))
}

// Load a config to a struct, val should be an pointer.
func (config *EeheConfig) Load(key string, val interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "yaml",
		Result:  val,
	})
	if err != nil {
		return err
	}

	return decoder.Decode(config.find(key))
}
