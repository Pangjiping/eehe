package env

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/Pangjiping/eehe/framework/contract"
	"io"
	"os"
	"path"
	"strings"
)

// EeheEnvService implements eehe env.
type EeheEnvService struct {
	// folder is the directory where the .env file is located.
	folder string
	// maps is used to save all environment variables.
	maps map[string]string
}

// NewEeheEnvService receives directory where the .env file is located,
// example: NewEeheEnvService("/envfolder/") will load file /envfolder/.env
// The data of the .env file is written in the form of key-value, example: FOO_ENV=BAR
func NewEeheEnvService(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("NewEeheEnvService params error")
	}

	folder := params[0].(string)
	eeheEnv := &EeheEnvService{
		folder: folder,
		maps:   map[string]string{"APP_ENV": contract.EnvDevelopment},
	}

	// Parse .env file.
	file := path.Join(folder, ".env")

	fi, err := os.Open(file)
	if err == nil {
		defer fi.Close()

		// read file
		reader := bufio.NewReader(fi)
		for {
			line, _, c := reader.ReadLine()
			if c == io.EOF {
				break
			}

			// parse by "="
			s := bytes.SplitN(line, []byte{'='}, 2)

			// passed if not in compliance
			if len(s) < 2 {
				continue
			}

			// saved to maps
			key := string(s[0])
			value := string(s[1])
			eeheEnv.maps[key] = value
		}
	}

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) < 2 {
			continue
		}
		eeheEnv.maps[pair[0]] = pair[1]
	}
	return eeheEnv, nil
}

// AppEnv gets current app env key APP_ENV.
func (svc *EeheEnvService) AppEnv() string {
	return svc.Get("APP_ENV")
}

// IsExist determines if an environment variable has been set.
func (svc *EeheEnvService) IsExist(key string) bool {
	_, ok := svc.maps[key]
	return ok
}

// Get gets given environment variable.
func (svc *EeheEnvService) Get(key string) string {
	if val, ok := svc.maps[key]; ok {
		return val
	}
	return ""
}

// All gets all env data.
func (svc *EeheEnvService) All() map[string]string {
	return svc.maps
}
