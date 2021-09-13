package cli

import (
	"github.com/chen-keinan/go-command-eval/eval"
	"github.com/chen-keinan/openshift-ordeal/internal/cli/commands"
	"github.com/chen-keinan/openshift-ordeal/internal/cli/mocks"
	"github.com/chen-keinan/openshift-ordeal/internal/common"
	m3 "github.com/chen-keinan/openshift-ordeal/internal/mocks"
	"github.com/chen-keinan/openshift-ordeal/internal/models"
	m2 "github.com/chen-keinan/openshift-ordeal/pkg/models"
	"github.com/chen-keinan/openshift-ordeal/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

//Test_StartCli tests
func Test_StartCli(t *testing.T) {
	fm := utils.NewKFolder()
	initBenchmarkSpecData(fm, ArgsData{SpecType: "openshift", SpecVersion: "v1.0.0"})
	files, err := utils.GetopenshiftBenchAuditFiles("openshift", "v1.0.0", fm)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(files), 1)
	assert.Equal(t, files[0].Name, common.FilesystemConfiguration)
}

func Test_ArgsSanitizer(t *testing.T) {
	args := []string{"--a", "-b"}
	ad := ArgsSanitizer(args)
	assert.Equal(t, ad.Filters[0], "a")
	assert.Equal(t, ad.Filters[1], "b")
	assert.False(t, ad.Help)
	args = []string{}
	ad = ArgsSanitizer(args)
	assert.True(t, ad.Filters[0] == "")
	args = []string{"--help"}
	ad = ArgsSanitizer(args)
	assert.True(t, ad.Help)
}

//Test_openshiftProbeHelpFunc test
func Test_openshiftProbeHelpFunc(t *testing.T) {
	cm := make(map[string]cli.CommandFactory)
	bhf := openshiftProbeHelpFunc(common.OpenshiftordealCli)
	helpFile := bhf(cm)
	assert.True(t, strings.Contains(helpFile, "Available commands are:"))
	assert.True(t, strings.Contains(helpFile, "Usage: openshift-ordeal [--version] [--help] <command> [<args>]"))
}

//Test_createCliBuilderData test
func Test_createCliBuilderData(t *testing.T) {
	cmdArgs := []string{"a"}
	ad := ArgsSanitizer(os.Args[1:])
	cmdArgs = append(cmdArgs, ad.Filters...)
	cmds := make([]cli.Command, 0)
	completedChan := make(chan bool)
	plChan := make(chan m2.OpenshiftAuditResults)
	// invoke cli
	cmds = append(cmds, commands.NewopenshiftAudit(ad.Filters, plChan, completedChan, []utils.FilesInfo{}, eval.NewEvalCmd()))
	c := createCliBuilderData(cmdArgs, cmds)
	_, ok := c["a"]
	assert.True(t, ok)

}

//Test_InvokeCli test
func Test_InvokeCli(t *testing.T) {
	ab := &models.AuditBench{}
	ab.AuditCommand = []string{"aaa"}
	ab.EvalExpr = "'$0' != '';"
	ab.CommandParams = map[int][]string{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	evalCmd := mocks.NewMockCmdEvaluator(ctrl)
	evalCmd.EXPECT().EvalCommand([]string{"aaa"}, ab.EvalExpr).Return(eval.CmdEvalResult{Match: true}).Times(1)
	completedChan := make(chan bool)
	plChan := make(chan m2.OpenshiftAuditResults)
	tl := m3.NewMockTestLoader(ctrl)
	tl.EXPECT().LoadAuditTests(nil).Return([]*models.SubCategory{{Name: "te", AuditTests: []*models.AuditBench{ab}}})
	go func() {
		<-plChan
		completedChan <- true
	}()
	kb := &commands.OpenshiftAudit{Evaluator: evalCmd, ResultProcessor: commands.GetResultProcessingFunction([]string{}), FileLoader: tl, OutputGenerator: commands.ConsoleOutputGenerator, PlChan: plChan, CompletedChan: completedChan}
	cmdArgs := []string{"a"}
	cmds := make([]cli.Command, 0)
	// invoke cli
	cmds = append(cmds, kb)
	c := createCliBuilderData(cmdArgs, cmds)
	a, err := invokeCommandCli(cmdArgs, c)
	assert.NoError(t, err)
	assert.True(t, a == 0)
}

func Test_InitPluginFolder(t *testing.T) {
	fm := utils.NewKFolder()
	initPluginFolders(fm)
}

func Test_InitPluginWorker(t *testing.T) {
	completedChan := make(chan bool)
	plChan := make(chan m2.OpenshiftAuditResults)
	go func() {
		plChan <- m2.OpenshiftAuditResults{}
		completedChan <- true
	}()
	initPluginWorker(plChan, completedChan)

}
