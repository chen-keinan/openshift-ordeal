package startup

import (
	"fmt"
	"github.com/chen-keinan/openshift-ordeal/internal/common"
	"github.com/chen-keinan/openshift-ordeal/pkg/utils"
	"github.com/gobuffalo/packr"
	"os"
	"path/filepath"
)

//GenerateOpenshiftBenchmarkFiles use packr to load benchmark audit test yaml
//nolint:gocyclo
func GenerateOpenshiftBenchmarkFiles() ([]utils.FilesInfo, error) {
	fileInfo := make([]utils.FilesInfo, 0)
	box := packr.NewBox("./../benchmark/openshift/v1.0.0/")
	// Add Master Node Configuration tests
	//1
	mnc, err := box.FindString(common.FilesystemConfiguration)
	if err != nil {
		return []utils.FilesInfo{}, fmt.Errorf("faild to load openshift benchmarks audit tests %s  %s", common.FilesystemConfiguration, err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.FilesystemConfiguration, Data: mnc})
	return fileInfo, nil
}

//GetHelpSynopsis get help synopsis file
func GetHelpSynopsis() string {
	box := packr.NewBox("./../cli/commands/help/")
	// Add Master Node Configuration tests
	hs, err := box.FindString(common.Synopsis)
	if err != nil {
		panic(fmt.Sprintf("faild to load cli help synopsis %s", err.Error()))
	}
	return hs
}

//SaveBenchmarkFilesIfNotExist create benchmark audit file if not exist
func SaveBenchmarkFilesIfNotExist(spec, version string, filesData []utils.FilesInfo) error {
	fm := utils.NewKFolder()
	folder, err := utils.GetBenchmarkFolder(spec, version, fm)
	if err != nil {
		return err
	}
	for _, fileData := range filesData {
		filePath := filepath.Join(folder, fileData.Name)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			f, err := os.Create(filePath)
			if err != nil {
				return fmt.Errorf(err.Error())
			}
			_, err = f.WriteString(fileData.Data)
			if err != nil {
				return fmt.Errorf("failed to write benchmark file")
			}
			err = f.Close()
			if err != nil {
				return fmt.Errorf("faild to close file %s", filePath)
			}
		}
	}
	return nil
}
