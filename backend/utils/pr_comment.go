package utils

import (
	"fmt"
	"log"
	"strconv"

	comment_utils "github.com/diggerhq/digger/libs/comment_utils/utils"
	"github.com/diggerhq/digger/libs/orchestrator"
	github2 "github.com/diggerhq/digger/libs/orchestrator/github"
)

type CommentReporter struct {
	PrNumber  int
	PrService *github2.GithubService
	CommentId int64
}

func InitCommentReporter(prService *github2.GithubService, prNumber int, commentMessage string) (*CommentReporter, error) {
	comment, err := prService.PublishComment(prNumber, commentMessage)
	if err != nil {
		return nil, fmt.Errorf("count not initialize comment reporter: %v", err)
	}
	commentId, err := strconv.ParseInt(fmt.Sprintf("%v", comment.Id), 10, 64)
	if err != nil {
		log.Printf("could not convert to int64, %v", err)
		return nil, fmt.Errorf("could not convert to int64, %v", err)
	}
	return &CommentReporter{
		PrNumber:  prNumber,
		PrService: prService,
		CommentId: commentId,
	}, nil
}

func ReportInitialJobsStatus(cr *CommentReporter, jobs []orchestrator.Job) error {
	prNumber := cr.PrNumber
	prService := cr.PrService
	commentId := cr.CommentId

	message := ""
	if len(jobs) > 0 {
		headers := []string{"Project", "Status"}
		message = comment_utils.CreateTableComment[orchestrator.Job](headers, jobs, func(_ int, job orchestrator.Job) []string {
			return []string{
				fmt.Sprintf(":clock11: **%v**", job.ProjectName),
				"pending...",
			}
		})
	} else {
		message = ":construction_worker: No projects impacted"
	}

	err := prService.EditComment(prNumber, commentId, message)
	return err
}

func ReportNoProjectsImpacted(cr *CommentReporter, jobs []orchestrator.Job) error {
	prNumber := cr.PrNumber
	prService := cr.PrService
	commentId := cr.CommentId
	message := "" +
		":construction_worker: The following projects are impacted\n\n"
	for _, job := range jobs {
		message = message + fmt.Sprintf(""+
			"<!-- PROJECTHOLDER %v -->\n"+
			":clock11: **%v**: pending...\n"+
			"<!-- PROJECTHOLDEREND %v -->\n"+
			"", job.ProjectName, job.ProjectName, job.ProjectName)
	}
	err := prService.EditComment(prNumber, commentId, message)
	return err
}
