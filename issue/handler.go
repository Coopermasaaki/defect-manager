package issue

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/opensourceways/robot-gitee-lib/client"
	"github.com/opensourceways/server-common-lib/utils"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/opensourceways/defect-manager/defect/app"
)

type EventHandler interface {
	HandleIssueEvent(e *sdk.IssueEvent) error
	HandleNoteEvent(e *sdk.NoteEvent) error
}

type iClient interface {
	CreateIssueComment(org, repo string, number string, comment string) error
	ListIssueComments(org, repo, number string) ([]sdk.Note, error)
	CloseIssue(owner, repo string, number string) error
}

func NewEventHandler(c *Config, s app.DefectService) *eventHandler {
	cli := client.NewClient(func() []byte {
		return []byte(c.RobotToken)
	})
	return &eventHandler{
		cfg:     c,
		cli:     cli,
		service: s,
	}
}

type eventHandler struct {
	cfg     *Config
	cli     iClient
	service app.DefectService
}

func (impl eventHandler) HandleIssueEvent(e *sdk.IssueEvent) error {
	if e.Issue.State != sdk.StatusOpen || e.Issue.TypeName != impl.cfg.IssueType {
		return nil
	}

	if _, err := impl.parseIssue(e.Issue.Body); err != nil {
		return impl.cli.CreateIssueComment(e.Project.Namespace,
			e.Project.Name, e.Issue.Number, strings.Replace(err.Error(), ". ", "\n\n", -1),
		)
	}

	return nil
}

func (impl eventHandler) HandleNoteEvent(e *sdk.NoteEvent) error {
	if !e.IsIssue() || e.Issue.TypeName != impl.cfg.IssueType || e.Issue.State == sdk.StatusClosed {
		return nil
	}

	commentIssue := func(content string) error {
		return impl.cli.CreateIssueComment(e.Project.Namespace,
			e.Project.Name, e.Issue.Number, content,
		)
	}

	if !impl.isValidCmd(e.Comment.Body) {
		if !strings.Contains(e.Comment.Body, "受影响版本排查") {
			return nil
		}

		if _, _, err := impl.parseComment(e.Comment.Body); err != nil {
			return commentIssue(err.Error())
		}

		return nil
	}

	issueInfo, err := impl.parseIssue(e.Issue.Body)
	if err != nil {
		return commentIssue(strings.Replace(err.Error(), ". ", "\n\n", -1))
	}

	comment := impl.approveCmdReplyToComment(e)
	if comment == "" {
		return nil
	}

	version, abi, err := impl.parseComment(comment)
	if err != nil {
		return commentIssue(err.Error())
	}

	if err = impl.checkRelatedPR(e, version); err != nil {
		return commentIssue(err.Error())
	}

	if err = impl.cli.CloseIssue(e.Project.Namespace, e.Project.Name, e.Issue.Number); err != nil {
		return fmt.Errorf("close pr error: %s", err.Error())
	}

	issueInfo[itemAbi] = strings.Join(abi, ",")
	err = impl.service.SaveDefect(impl.toCmd(e, issueInfo, version))
	if err == nil {
		return commentIssue("Your issue is accepted, thank you")
	}

	return err
}

// the content of the comment of the newest /approve reply to
func (impl eventHandler) approveCmdReplyToComment(e *sdk.NoteEvent) string {
	comments, err := impl.cli.ListIssueComments(e.Project.Namespace, e.Project.Name, e.Issue.Number)
	if err != nil {
		logrus.Errorf("get comments error: %s", err.Error())

		return ""
	}

	var id int32
	// Iterate from the end to get the latest approve command
	for i := len(comments) - 1; i >= 0; i-- {
		if strings.Contains(comments[i].Body, cmdApprove) &&
			comments[i].User.Login == e.Issue.User.Login { // approve from the author of the issue is valid
			id = comments[i].InReplyToId
			break
		}
	}
	if id == 0 {
		return ""
	}

	for _, v := range comments {
		if v.Id == id {
			return v.Body
		}
	}

	return ""
}

func (impl eventHandler) toCmd(e *sdk.NoteEvent, issueInfo map[string]string, version []string) app.CmdToHandleDefect {
	return app.CmdToHandleDefect{}
}

func (impl eventHandler) checkRelatedPR(e *sdk.NoteEvent, versions []string) error {
	endpoint := fmt.Sprintf("https://gitee.com/api/v5/repos/%v/issues/%v/pull_requests?access_token=%s&repo=%s",
		e.Project.Namespace, e.Issue.Number, impl.cfg.RobotToken, e.Project.Name,
	)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	var prs []sdk.PullRequest
	cli := utils.NewHttpClient(3)
	bytes, _, err := cli.Download(req)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &prs); err != nil {
		return err
	}

	mergedVersion := sets.NewString()
	for _, pr := range prs {
		if pr.State == sdk.StatusMerged {
			mergedVersion.Insert(pr.Base.Ref)
		}
	}

	var relatedPRNotMerged []string
	for _, v := range versions {
		if !mergedVersion.Has(v) {
			relatedPRNotMerged = append(relatedPRNotMerged, v)
		}
	}

	if len(relatedPRNotMerged) != 0 {
		return fmt.Errorf("受影响分支关联pr未合入: %s", strings.Join(relatedPRNotMerged, ","))
	}

	return nil
}