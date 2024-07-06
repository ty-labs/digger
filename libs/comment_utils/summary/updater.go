package comment_updater

import (
	"fmt"
	"log"
	"strconv"

	comment_utils "github.com/diggerhq/digger/libs/comment_utils/utils"
	"github.com/diggerhq/digger/libs/orchestrator"
	"github.com/diggerhq/digger/libs/orchestrator/scheduler"
)

type CommentUpdater interface {
	UpdateComment(jobs []scheduler.SerializedJob, prNumber int, prService orchestrator.PullRequestService, prCommentId int64) error
}

type BasicCommentUpdater struct{}

func (b BasicCommentUpdater) UpdateComment(jobs []scheduler.SerializedJob, prNumber int, prService orchestrator.PullRequestService, prCommentId int64) error {
	jobSpecs, err := scheduler.GetJobSpecs(jobs)
	if err != nil {
		log.Printf("could not get jobspecs: %v", err)
		return err
	}
	firstJobSpec := jobSpecs[0]
	isPlan := firstJobSpec.IsPlan()

	headers := []string{}
	if isPlan {
		headers = append(headers, "Project", "Status", "Plan", "+", "~", "-")
	} else {
		headers = append(headers, "Project", "Status", "Apply")
	}
	message := comment_utils.CreateTableComment[scheduler.SerializedJob](headers, jobs, func(index int, job scheduler.SerializedJob) []string {
		jobSpec := jobSpecs[index]
		prCommentUrl := job.PRCommentUrl
		if isPlan {
			return []string{
				fmt.Sprintf("%v **%v**", job.Status.ToEmoji(), jobSpec.ProjectName),
				fmt.Sprintf("<a href='%v'>%v</a>", *job.WorkflowRunUrl, job.Status.ToString()),
				fmt.Sprintf("<a href='%v'>plan</a>", prCommentUrl),
				strconv.FormatInt(int64(job.ResourcesCreated), 10),
				strconv.FormatInt(int64(job.ResourcesUpdated), 10),
				strconv.FormatInt(int64(job.ResourcesDeleted), 10),
			}
			return fmt.Sprintf("|%v **%v** |<a href='%v'>%v</a> | <a href='%v'>plan</a> | %v | %v | %v|\n", job.Status.ToEmoji(), jobSpec.ProjectName, *job.WorkflowRunUrl, job.Status.ToString(), prCommentUrl, job.ResourcesCreated, job.ResourcesUpdated, job.ResourcesDeleted)
		}
		return []string{
			fmt.Sprintf("%v **%v**", job.Status.ToEmoji(), jobSpec.ProjectName),
			fmt.Sprintf("<a href='%v'>%v</a>", *job.WorkflowRunUrl, job.Status.ToString()),
			fmt.Sprintf("<a href='%v'>apply</a>", prCommentUrl),
		}
		return fmt.Sprintf("|%v **%v** |<a href='%v'>%v</a> | <a href='%v'>apply</a> |\n", job.Status.ToEmoji(), jobSpec.ProjectName, *job.WorkflowRunUrl, job.Status.ToString(), prCommentUrl)
	})

	prService.EditComment(prNumber, prCommentId, message)
	return nil
}

type NoopCommentUpdater struct{}

func (b NoopCommentUpdater) UpdateComment(jobs []scheduler.SerializedJob, prNumber int, prService orchestrator.PullRequestService, prCommentId int64) error {
	return nil
}
