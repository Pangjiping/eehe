package config

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Pangjiping/eehe/framework"
	"github.com/Pangjiping/eehe/framework/contract"
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

func (eehe *EeheConfig) removeConfigFile(folder string, file string) error {
	eehe.mu.Lock()
	defer eehe.mu.Unlock()

	s := strings.Split(file, ".")

	// 只有yaml/yml后缀才执行
	if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
		name := s[0]
		delete(eehe.configRaws, name)
		delete(eehe.configMaps, name)
	}
	return nil
}
