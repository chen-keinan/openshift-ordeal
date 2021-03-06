package startup

import (
	"github.com/chen-keinan/openshift-ordeal/internal/common"
	"github.com/chen-keinan/openshift-ordeal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

//Test_CreateopenshiftBenchmarkFilesIfNotExist test
func Test_CreateopenshiftBenchmarkFilesIfNotExist(t *testing.T) {
	bFiles, err := GenerateOpenshiftBenchmarkFiles()
	if err != nil {
		t.Fatal(err)
	}
	// generate test with packr
	assert.Equal(t, bFiles[0].Name, common.FilesystemConfiguration)
	fm := utils.NewKFolder()
	err = utils.CreateBenchmarkFolderIfNotExist("openshift", "v1.0.0", fm)
	assert.NoError(t, err)
	// save benchmark files to folder
	err = SaveBenchmarkFilesIfNotExist("openshift", "v1.0.0", bFiles)
	assert.NoError(t, err)
	// fetch files from benchmark folder
	bFiles, err = utils.GetopenshiftBenchAuditFiles("openshift", "v1.0.0", fm)
	assert.Equal(t, bFiles[0].Name, common.FilesystemConfiguration)
	assert.NoError(t, err)
	err = os.RemoveAll(utils.GetHomeFolder())
	assert.NoError(t, err)
}

//Test_GetHelpSynopsis test
func Test_GetHelpSynopsis(t *testing.T) {
	hs := GetHelpSynopsis()
	assert.True(t, len(hs) != 0)
}

//Test_SaveBenchmarkFilesIfNotExist test
func Test_SaveBenchmarkFilesIfNotExist(t *testing.T) {
	fm := utils.NewKFolder()
	folder, err2 := utils.GetBenchmarkFolder("openshift", "v1.6.0", fm)
	assert.NoError(t, err2)
	err := os.RemoveAll(folder)
	assert.NoError(t, err)
	err = os.RemoveAll(utils.GetHomeFolder())
	assert.NoError(t, err)
}
