package goDotEnvLight

import (
	"bufio"
	"os"
	"strings"
)

// 1. export
// 2. exportLatest

type EnvInfo struct {
	Key   string
	Value string
	Valid bool
}

type Envs struct {
	Infos []EnvInfo
}

func Export(overwrite bool, envFilePaths ...string) (succeeds map[string]string, fails map[string]string, err error) {

	succeeds = make(map[string]string)
	fails = make(map[string]string)

	envs := Envs{Infos: make([]EnvInfo, 0)}

	for _, envFilePath := range envFilePaths {
		// Read .env(s)
		if err = envs.readDotEnvFile(envFilePath); err != nil {
			return
		}
	}

	// Check Prev Env
	if !overwrite {
		if err = envs.checkPrevEnvs(); err != nil {
			return
		}
	}

	// Set Env
	envs.setEnvs()

	for _, envInfo := range envs.Infos {
		switch envInfo.Valid {
		case true:
			succeeds[envInfo.Key] = envInfo.Value
		case false:
			fails[envInfo.Key] = envInfo.Value
		}
	}

	return
}

func (env *Envs) readDotEnvFile(dotEnvFilePath string) (err error) {
	dotEnvFile, err := os.Open(dotEnvFilePath)
	if err != nil {
		return err
	}
	defer dotEnvFile.Close()

	scanner := bufio.NewScanner(dotEnvFile)
	for scanner.Scan() {
		tLine := scanner.Text()
		tLineSplit := strings.Split(tLine, "=")
		if len(tLineSplit) < 2 {
			continue
		}
		tKey := tLineSplit[0]
		tValue := strings.Join(tLineSplit[1:], "=")
		env.Infos = append(env.Infos, EnvInfo{Key: tKey, Value: tValue, Valid: true})
	}

	return
}

func (env *Envs) checkPrevEnvs() (err error) {

	prevEnvKeys := map[string]bool{}
	prevRawEnvs := os.Environ()
	for _, prevRawEnvLine := range prevRawEnvs {
		key := strings.Split(prevRawEnvLine, "=")[0]
		prevEnvKeys[key] = true
	}

	for i, envInfo := range env.Infos {
		if prevEnvKeys[envInfo.Key] {
			envInfo.Valid = false
			env.Infos[i] = envInfo
		}
	}

	return
}

func (env *Envs) setEnvs() {
	for i, envInfo := range env.Infos {
		if !envInfo.Valid {
			continue
		}
		if err := os.Setenv(envInfo.Key, envInfo.Value); err != nil {
			envInfo.Valid = false
			env.Infos[i] = envInfo
		}
	}
}
