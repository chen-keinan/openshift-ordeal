package commands

import (
	"fmt"
	"github.com/chen-keinan/go-command-eval/eval"
	"github.com/chen-keinan/openshift-ordeal/internal/logger"
	"github.com/chen-keinan/openshift-ordeal/internal/models"
	"github.com/chen-keinan/openshift-ordeal/internal/reports"
	"github.com/chen-keinan/openshift-ordeal/internal/startup"
	"github.com/chen-keinan/openshift-ordeal/pkg/filters"
	m2 "github.com/chen-keinan/openshift-ordeal/pkg/models"
	"github.com/chen-keinan/openshift-ordeal/pkg/utils"
	"github.com/chen-keinan/openshift-ordeal/ui"
	"github.com/mitchellh/colorstring"
	"github.com/olekukonko/tablewriter"
	"os"
)

//OpenshiftAudit openshift benchmark object
type OpenshiftAudit struct {
	ResultProcessor ResultProcessor
	OutputGenerator ui.OutputGenerator
	FileLoader      TestLoader
	PredicateChain  []filters.Predicate
	PredicateParams []string
	PlChan          chan m2.OpenshiftAuditResults
	CompletedChan   chan bool
	FilesInfo       []utils.FilesInfo
	Evaluator       eval.CmdEvaluator
	log             *logger.OsOrdealLogger
}

// ResultProcessor process audit results
type ResultProcessor func(at *models.AuditBench, isSucceeded bool) []*models.AuditBench

// ConsoleOutputGenerator print audit tests to stdout
var ConsoleOutputGenerator ui.OutputGenerator = func(at []*models.SubCategory, log *logger.OsOrdealLogger) {
	grandTotal := make([]models.AuditTestTotals, 0)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Category", "Status", "Type", "Audit Test Description"})
	table.SetAutoWrapText(false)
	table.SetBorder(true) // Set
	for _, a := range at {
		categoryTotal := printTestResults(a.AuditTests, table, a.Name)
		grandTotal = append(grandTotal, categoryTotal)
	}
	table.SetAutoMergeCellsByColumnIndex([]int{0})
	table.SetRowLine(true)
	table.Render()
	log.Console(printFinalResults(grandTotal))
}

// ClassicOutputGenerator print audit tests to stdout in classic view
var ClassicOutputGenerator ui.OutputGenerator = func(at []*models.SubCategory, log *logger.OsOrdealLogger) {
	grandTotal := make([]models.AuditTestTotals, 0)
	for _, a := range at {
		log.Console(fmt.Sprintf("%s %s\n", "[Category]", a.Name))
		categoryTotal := printClassicTestResults(a.AuditTests, log)
		grandTotal = append(grandTotal, categoryTotal)
	}
	log.Console(printFinalResults(grandTotal))
}

func printFinalResults(grandTotal []models.AuditTestTotals) string {
	finalTotal := calculateFinalTotal(grandTotal)
	passTest := colorstring.Color("[green]Pass:")
	failTest := colorstring.Color("[red]Fail:")
	warnTest := colorstring.Color("[yellow]Warn:")
	title := "Test Result Total:   "
	return fmt.Sprintf("%s %s %d , %s %d , %s %d ", title, passTest, finalTotal.Pass, warnTest, finalTotal.Warn, failTest, finalTotal.Fail)
}

func calculateFinalTotal(granTotal []models.AuditTestTotals) models.AuditTestTotals {
	var (
		warn int
		fail int
		pass int
	)
	for _, total := range granTotal {
		warn = warn + total.Warn
		fail = fail + total.Fail
		pass = pass + total.Pass
	}
	return models.AuditTestTotals{Pass: pass, Fail: fail, Warn: warn}
}

// ReportOutputGenerator print failed audit test to human report
var ReportOutputGenerator ui.OutputGenerator = func(at []*models.SubCategory, log *logger.OsOrdealLogger) {
	for _, a := range at {
		log.Table(reports.GenerateAuditReport(a.AuditTests))
	}
}

// simpleResultProcessor process audit results to stdout print only
var simpleResultProcessor ResultProcessor = func(at *models.AuditBench, isSucceeded bool) []*models.AuditBench {
	return AddAllMessages(at, isSucceeded)
}

// ResultProcessor process audit results to std out and failure results
var reportResultProcessor ResultProcessor = func(at *models.AuditBench, isSucceeded bool) []*models.AuditBench {
	// append failed messages
	return AddFailedMessages(at, isSucceeded)
}

//CmdEvaluator interface expose one method to evaluate command with evalExpr
//openshift-audit.go
//go:generate mockgen -destination=../mocks/mock_CmdEvaluator.go -package=mocks . CmdEvaluator
type CmdEvaluator interface {
	EvalCommand(commands []string, evalExpr string) eval.CmdEvalResult
}

//NewopenshiftAudit new audit object
func NewopenshiftAudit(filters []string, plChan chan m2.OpenshiftAuditResults, completedChan chan bool, fi []utils.FilesInfo, evaluator CmdEvaluator) *OpenshiftAudit {
	return &OpenshiftAudit{
		PredicateChain:  buildPredicateChain(filters),
		PredicateParams: buildPredicateChainParams(filters),
		ResultProcessor: GetResultProcessingFunction(filters),
		OutputGenerator: getOutputGeneratorFunction(filters),
		FileLoader:      NewFileLoader(),
		PlChan:          plChan,
		FilesInfo:       fi,
		Evaluator:       evaluator,
		CompletedChan:   completedChan}
}

//Help return benchmark command help
func (ldx OpenshiftAudit) Help() string {
	return startup.GetHelpSynopsis()
}

//Run execute the full openshift benchmark
func (ldx *OpenshiftAudit) Run(args []string) int {
	// load audit tests fro benchmark folder
	auditTests := ldx.FileLoader.LoadAuditTests(ldx.FilesInfo)
	// filter tests by cmd criteria
	ft := filteredAuditBenchTests(auditTests, ldx.PredicateChain, ldx.PredicateParams)
	//execute audit tests and show it in progress bar
	completedTest := executeTests(ft, ldx.runAuditTest, ldx.log)
	// generate output data
	ui.PrintOutput(completedTest, ldx.OutputGenerator, ldx.log)
	// send test results to plugin
	sendResultToPlugin(ldx.PlChan, ldx.CompletedChan, completedTest)
	return 0
}

func sendResultToPlugin(plChan chan m2.OpenshiftAuditResults, completedChan chan bool, auditTests []*models.SubCategory) {
	ka := m2.OpenshiftAuditResults{BenchmarkType: "openshift", Categories: make([]m2.AuditBenchResult, 0)}
	for _, at := range auditTests {
		for _, ab := range at.AuditTests {
			var testResult = "FAIL"
			if ab.TestSucceed {
				testResult = "PASS"
			}
			abr := m2.AuditBenchResult{Category: at.Name, ProfileApplicability: ab.ProfileApplicability, Description: ab.Description, AuditCommand: ab.AuditCommand, Remediation: ab.Remediation, Impact: ab.Impact, AdditionalInfo: ab.AdditionalInfo, References: ab.References, TestResult: testResult}
			ka.Categories = append(ka.Categories, abr)
		}
	}
	plChan <- ka
	<-completedChan
}

// runAuditTest execute category of audit tests
func (ldx *OpenshiftAudit) runAuditTest(at *models.AuditBench) []*models.AuditBench {
	auditRes := make([]*models.AuditBench, 0)
	if at.NonApplicable {
		auditRes = append(auditRes, at)
		return auditRes
	}
	// execute audit test command
	cmdEvalResult := ldx.Evaluator.EvalCommand(at.AuditCommand, at.EvalExpr)
	// continue with result processing
	auditRes = append(auditRes, ldx.ResultProcessor(at, cmdEvalResult.Match)...)
	return auditRes
}

//Synopsis for help
func (ldx *OpenshiftAudit) Synopsis() string {
	return ldx.Help()
}
