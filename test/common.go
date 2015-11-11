/*
test package contains all test utilities, mocks and stubs
*/
package test

import (
	"io"
	"log"

	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/Symantec/image-lifecycle-manager/pkg/builder"
	"github.com/Symantec/image-lifecycle-manager/pkg/common"
	check "gopkg.in/check.v1"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// TestArtifact dump artifact for testing with in memory content of 3 bytes.
type TestArtifact struct{}

func (ta *TestArtifact) GetName() string {
	return "TestArtifactName"
}

func (ta *TestArtifact) GetType() builder.ArtifactType {
	return builder.ArtifactLog
}

func (ta *TestArtifact) GetContent() io.ReadCloser {
	return common.NewClosableReader([]byte{1, 2, 3})
}

type artifactContentChecker struct {
	*check.CheckerInfo
}

var (
	// ArtifactContentChecker checks if content matches regex
	ArtifactContentChecker check.Checker = &artifactContentChecker{}
)

func (*artifactContentChecker) Info() *check.CheckerInfo {
	return &check.CheckerInfo{Name: "ArtifactContentChecker",
		Params: []string{"value", "regex"}}
}

func (*artifactContentChecker) Check(params []interface{}, names []string) (result bool, error string) {
	reStr, ok := params[1].(string)
	if !ok {
		return false, "Regex must be a string"
	}

	artifact, ok := params[0].(builder.Artifact)
	if !ok {
		return false, "Value should be builder.Artifact type"
	}

	if artifact.GetContent() == nil {
		return false, "No content"
	}

	content := artifact.GetContent()
	defer content.Close()
	data, err := ioutil.ReadAll(content)

	if err != nil {
		return false, fmt.Sprintf("Couldn't read artifact content %s", err)
	}

	strContent := fmt.Sprintf("%s", data)
	matches, err := regexp.MatchString(reStr, strContent)
	if err != nil {
		return false, "Can't compile regex: " + err.Error()
	}
	return matches, ""
}

type TestNotifier struct {
	notifications []string
}

func (n *TestNotifier) Notify(source string, status string, payload string) {
	n.notifications = append(n.notifications, fmt.Sprintf("SOURCE:%s STATUS:%s PAYLOAD:%s", source, status, payload))
}

func (n *TestNotifier) GetNotifications() []string {
	return n.notifications
}

func (n *TestNotifier) Clean() {
	n.notifications = n.notifications[0:0]
}
