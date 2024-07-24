package game

import (
	"errors"

	domain_err "moj/domain/pkg/error"
	"moj/domain/pkg/queue"
)

type GameQuestion struct {
	QuestionID string
	Score      int
}

type SignUpAccount struct {
	AccountID  string
	SignUpTime int64
}

var (
	ErrGameNotFound        = errors.New("game not found")
	ErrAccountNotFound     = errors.New("sign up account not found")
	ErrQuestionNotFound    = errors.New("collect question not found")
	ErrAccountAlreadyExist = errors.Join(domain_err.ErrDuplicated, errors.New("account already sign up"))
	ErrInvalidTimeRange    = errors.Join(domain_err.ErrInValided, errors.New("invalid time range"))
)

type Game struct {
	GameID         string
	AccountID      string
	Title          string
	Description    string
	CreateTime     int64
	StartTime      int64
	EndTime        int64
	QuestionList   []GameQuestion
	SignUpUserList []SignUpAccount
}

func NewGame(userID string, title, desc string,
	time, startTime, endTime int64, ques []GameQuestion) (*Game, error) {

	if !isValidTimeRange(startTime, endTime) {
		return nil, ErrInvalidTimeRange
	}

	return &Game{
		AccountID:    userID,
		Title:        title,
		Description:  desc,
		CreateTime:   time,
		StartTime:    startTime,
		EndTime:      endTime,
		QuestionList: ques,
	}, nil
}

func isValidTimeRange(startTime, endTime int64) bool {
	return startTime < endTime
}

func (g *Game) modify(cmd ModifyGameCmd) error {
	if !isValidTimeRange(cmd.StartTime, cmd.EndTime) {
		return ErrInvalidTimeRange
	}

	g.Title = cmd.Title
	g.Description = cmd.Descirption
	g.StartTime = cmd.StartTime
	g.EndTime = cmd.EndTime
	g.QuestionList = cmd.QuestionList

	return nil
}

func (g *Game) findSignedUp(accountID string) int {
	for id, user := range g.SignUpUserList {
		if user.AccountID == accountID {
			return id
		}
	}
	return -1
}

func (g *Game) findQuestion(questionID string) int {
	for id, ques := range g.QuestionList {
		if ques.QuestionID == questionID {
			return id
		}
	}
	return -1
}

func (g *Game) calculate(queue queue.EventQueue, cmd CalculateScoreCmd) error {
	if g.findSignedUp(cmd.AccountID) == -1 {
		return ErrAccountNotFound
	}

	queId := g.findQuestion(cmd.QuestionID)
	if queId == -1 {
		return ErrQuestionNotFound
	}

	gross := g.QuestionList[queId].Score
	num := cmd.NumberFinishedAt - cmd.LastFinishedAt
	deno := cmd.TotalQuestion

	score := getScore(num, deno, gross)

	event := CalculateScoreEvent{
		GameID:     g.GameID,
		AccountID:  cmd.AccountID,
		QuestionID: cmd.QuestionID,
		Score:      score,
	}
	return queue.EnQueue(event)
}

func getScore(num, deno, gross int) int {
	return num * gross / deno
}

func (g *Game) signUp(queue queue.EventQueue, insertArrayFn func(gid, aid string, time int64) error,
	cmd SignUpGameCmd) error {
	if g.findSignedUp(cmd.AccountID) > -1 {
		return ErrAccountAlreadyExist
	}

	err := insertArrayFn(g.GameID, cmd.AccountID, cmd.Time)
	if err != nil {
		return err
	}

	event := SignUpGameEvent{
		GameID:     g.GameID,
		AccountID:  cmd.AccountID,
		SignUpTime: cmd.Time,
	}
	return queue.EnQueue(event)
}

func (g *Game) cancelSignUp(queue queue.EventQueue, deleteArrayFn func(gid, aid string) error,
	cmd CancelSignUpGameCmd) error {
	accId := g.findSignedUp(cmd.AccountID)
	if accId == -1 {
		return ErrAccountNotFound
	}

	err := deleteArrayFn(g.GameID, cmd.AccountID)
	if err != nil {
		return err
	}

	event := CancelSignUpGameEvent{
		GameID:     g.GameID,
		AccountID:  cmd.AccountID,
		CancelTime: cmd.Time,
	}
	return queue.EnQueue(event)
}
