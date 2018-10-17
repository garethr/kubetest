package kubetest

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/garethr/skyhook"
	"gopkg.in/yaml.v2"

	"github.com/garethr/kubetest/assert"
)

func listTests(testDir string) []string {
	fullTestDir, err := filepath.Abs(testDir)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(fullTestDir); os.IsNotExist(err) {
		log.Fatal(fmt.Sprintf("Unable to find test directory: %s", fullTestDir))
	}
	var files []string
	filepath.Walk(fullTestDir, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(".sky$", f.Name())
			if err == nil && r {
				files = append(files, f.Name())
			}
		}
		return nil
	})
	return files
}

func Run(config []byte, filePath string, fileName string) bool {
	var spec interface{}
	yaml.Unmarshal(config, &spec)

	// A file with all commented out content will otherwise
	// panic, and we can't make assertions against a blank data
	// structure anyhow. But there are valid usecases for commented
	// out files so we warn without throwing an error
	if spec == nil {
		log.Warn("The document " + fileName + " does not contain any content")
		return true
	}

	sky := skyhook.New([]string{filePath})
	globals := map[string]interface{}{
		"file_name":           fileName,
		"spec":                spec,
		"assert_equal":        assert.Equal,
		"assert_contains":     assert.Contains,
		"assert_not_contains": assert.NotContains,
		"assert_not_equal":    assert.NotEqual,
		"assert_nil":          assert.Nil,
		"assert_not_nil":      assert.NotNil,
		"fail":                assert.Fail,
		"fail_now":            assert.FailNow,
		"assert_empty":        assert.Empty,
		"assert_not_empty":    assert.NotEmpty,
		"assert_true":         assert.True,
		"assert_false":        assert.False,
	}

	tests := listTests(filePath)
	for _, test := range tests {
		_, err := sky.Run(test, globals)
		if err != nil {
			log.Fatal(err)
		}
	}

	success := true

	for _, result := range assert.Results {
		message := fmt.Sprintf("%s %s", fileName, result.Message)
		if result.Kind == assert.AssertionError {
			log.Error(message)
		} else if result.Kind == assert.AssertionFailure {
			log.Warn(message)
			success = false
		} else if result.Kind == assert.AssertionSuccess {
			log.Info(message)
		}
	}

	assert.Results = nil

	return success
}

// detectLineBreak returns the relevant platform specific line ending
func detectLineBreak(haystack []byte) string {
	windowsLineEnding := bytes.Contains(haystack, []byte("\r\n"))
	if windowsLineEnding && runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

func Runs(config []byte, filePath string, fileName string) bool {

	if len(config) == 0 {
		log.Error("The document " + fileName + " appears to be empty")
	}

	bits := bytes.Split(config, []byte("---"+detectLineBreak(config)))

	results := make([]bool, 0)
	for _, element := range bits {
		if len(element) > 0 {
			result := Run(element, filePath, fileName)
			results = append(results, result)
		}
	}

	for _, value := range results {
		if value == false {
			return false
		}
	}
	return true
}
