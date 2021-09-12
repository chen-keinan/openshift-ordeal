package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chen-keinan/openshift-scrutiny/pkg/models"
	"net/http"
	"strings"
)

//openshiftBenchAuditResultHook this plugin method accept openshift audit bench results
//event include test data , description , audit, remediation and result
func openshiftBenchAuditResultHook(openshiftAuditResults models.OpenshiftAuditResults) error {
	var sb = new(bytes.Buffer)
	err := json.NewEncoder(sb).Encode(openshiftAuditResults)
	fmt.Print(openshiftAuditResults)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", "http://localhost:8090/audit-results", strings.NewReader(sb.String()))
	if err != nil {
		return err
	}
	client := http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}
