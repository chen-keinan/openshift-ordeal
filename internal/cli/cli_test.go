package cli

import (
	"github.com/chen-keinan/go-command-eval/eval"
	"github.com/chen-keinan/openshift-scrutiny/internal/cli/commands"
	"github.com/chen-keinan/openshift-scrutiny/internal/cli/mocks"
	"github.com/chen-keinan/openshift-scrutiny/internal/common"
	m3 "github.com/chen-keinan/openshift-scrutiny/internal/mocks"
	"github.com/chen-keinan/openshift-scrutiny/internal/models"
	m2 "github.com/chen-keinan/openshift-scrutiny/pkg/models"
	"github.com/chen-keinan/openshift-scrutiny/pkg/utils"
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
	assert.Equal(t, len(files), 27)
	assert.Equal(t, files[0].Name, common.FilesystemConfiguration)
	assert.Equal(t, files[1].Name, common.ConfigureSoftwareUpdates)
	assert.Equal(t, files[2].Name, common.ConfigureSudo)
	assert.Equal(t, files[3].Name, common.FilesystemIntegrityChecking)
	assert.Equal(t, files[4].Name, common.AdditionalProcessHardening)
	assert.Equal(t, files[5].Name, common.MandatoryAccessControl)
	assert.Equal(t, files[6].Name, common.WarningBanners)
	assert.Equal(t, files[7].Name, common.EnsureUpdates)
	assert.Equal(t, files[8].Name, common.InetdServices)
	assert.Equal(t, files[9].Name, common.SpecialPurposeServices)
	assert.Equal(t, files[10].Name, common.ServiceClients)
	assert.Equal(t, files[11].Name, common.NonessentialServices)
	assert.Equal(t, files[12].Name, common.NetworkParameters)
	assert.Equal(t, files[13].Name, common.NetworkParametersHost)
	assert.Equal(t, files[14].Name, common.TCPWrappers)
	assert.Equal(t, files[15].Name, common.FirewallConfiguration)
	assert.Equal(t, files[16].Name, common.ConfigureLogging)
	assert.Equal(t, files[17].Name, common.EnsureLogrotateConfigured)
	assert.Equal(t, files[18].Name, common.EnsureLogrotateAssignsAppropriatePermissions)
	assert.Equal(t, files[19].Name, common.ConfigureCron)
	assert.Equal(t, files[20].Name, common.SSHServerConfiguration)
	assert.Equal(t, files[21].Name, common.ConfigurePam)
	assert.Equal(t, files[22].Name, common.UserAccountsAndEnvironment)
	assert.Equal(t, files[23].Name, common.RootLoginRestrictedSystemConsole)
	assert.Equal(t, files[24].Name, common.EnsureAccessSuCommandRestricted)
	assert.Equal(t, files[25].Name, common.SystemFilePermissions)
	assert.Equal(t, files[26].Name, common.UserAndGroupSettings)
}

func Test_ArgsSanitizer(t *testing.T) {
	os.Args = append(os.Args,"--report")
		os.Args = append(os.Args,"--exclude=1.1.10")
	os.Args = append(os.Args,"--help")
	ad := ArgsSanitizer()
	assert.Equal(t, ad.Filters[0], "report")
	assert.Equal(t, ad.Filters[1], "include=")
	assert.Equal(t, ad.Filters[2], "exclude=1.1.10")
  	assert.True(t, ad.Help)
}

//Test_openshiftProbeHelpFunc test
func Test_openshiftProbeHelpFunc(t *testing.T) {
	cm := make(map[string]cli.CommandFactory)
	bhf := openshiftProbeHelpFunc(common.openshiftProbe)
	helpFile := bhf(cm)
	assert.True(t, strings.Contains(helpFile, "Available commands are:"))
	assert.True(t, strings.Contains(helpFile, "Usage: openshift-scrutiny [--version] [--help] <command> [<args>]"))
}

//Test_createCliBuilderData test
func Test_createCliBuilderData(t *testing.T) {
 	cmdArgs := []string{"a"}
	ad := ArgsSanitizer()
	cmdArgs = append(cmdArgs, ad.Filters...)
	cmds := make([]cli.Command, 0)
	completedChan := make(chan bool)
	plChan := make(chan m2.openshiftAuditResults)
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
	plChan := make(chan m2.openshiftAuditResults)
	tl := m3.NewMockTestLoader(ctrl)
	tl.EXPECT().LoadAuditTests(nil).Return([]*models.SubCategory{{Name: "te", AuditTests: []*models.AuditBench{ab}}})
	go func() {
		<-plChan
		completedChan <- true
	}()
	kb := &commands.openshiftAudit{Evaluator: evalCmd, ResultProcessor: commands.GetResultProcessingFunction([]string{}), FileLoader: tl, OutputGenerator: commands.ConsoleOutputGenerator, PlChan: plChan, CompletedChan: completedChan}
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
	plChan := make(chan m2.openshiftAuditResults)
	go func() {
		plChan <- m2.openshiftAuditResults{}
		completedChan <- true
	}()
	initPluginWorker(plChan, completedChan)

}
