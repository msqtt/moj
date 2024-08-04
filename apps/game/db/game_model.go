package db

import (
	"moj/domain/game"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameQuestion struct {
	QuestionID string `bson:"question_id,omitempty"`
	Score      int    `bson:"score"`
}

type SignUpAccount struct {
	AcountID   string    `bson:"account_id,omitempty"`
	SignUpTime time.Time `bson:"sign_up_time"`
}

type GameModel struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	AccountID         string             `bson:"account_id,omitempty"`
	Title             string             `bson:"title"`
	Description       string             `bson:"description"`
	StartTime         time.Time          `bson:"start_time"`
	EndTime           time.Time          `bson:"end_time"`
	QuestionList      []GameQuestion     `bson:"question_list"`
	SignUpAccountList []SignUpAccount    `bson:"sign_up_account_list"`
	CreateTime        time.Time          `bson:"create_time"`
}

func FromQuestionList(questionList []game.GameQuestion) []GameQuestion {
	gameQuestionList := make([]GameQuestion, 0)

	for _, question := range questionList {
		gameQuestionList = append(gameQuestionList, GameQuestion{
			QuestionID: question.QuestionID,
			Score:      question.Score,
		})
	}

	return gameQuestionList
}

func FromSignUpAccountList(signUpAccountList []game.SignUpAccount) []SignUpAccount {
	signUpAccountModelList := make([]SignUpAccount, 0)

	for _, signUpAccount := range signUpAccountList {
		signUpAccountModelList = append(signUpAccountModelList, SignUpAccount{
			AcountID:   signUpAccount.AccountID,
			SignUpTime: time.Unix(signUpAccount.SignUpTime, 0),
		})
	}

	return signUpAccountModelList
}

func NewFromAggreation(g *game.Game) *GameModel {
	id, _ := primitive.ObjectIDFromHex(g.GameID)
	gameQuestionList := FromQuestionList(g.QuestionList)
	signUpAccountList := FromSignUpAccountList(g.SignUpUserList)
	return &GameModel{
		ID:                id,
		AccountID:         g.AccountID,
		Title:             g.Title,
		Description:       g.Description,
		StartTime:         time.Unix(g.StartTime, 0),
		EndTime:           time.Unix(g.EndTime, 0),
		QuestionList:      gameQuestionList,
		SignUpAccountList: signUpAccountList,
		CreateTime:        time.Unix(g.CreateTime, 0),
	}
}

func (g *GameModel) ToAggreation() *game.Game {
	gameQuestionList := ToQuestionList(g.QuestionList)
	signUpAccountList := ToSignUpAccountList(g.SignUpAccountList)
	return &game.Game{
		GameID:         g.ID.Hex(),
		AccountID:      g.AccountID,
		Title:          g.Title,
		Description:    g.Description,
		StartTime:      g.StartTime.Unix(),
		EndTime:        g.EndTime.Unix(),
		QuestionList:   gameQuestionList,
		SignUpUserList: signUpAccountList,
		CreateTime:     g.CreateTime.Unix(),
	}
}

func ToQuestionList(gameQuestionList []GameQuestion) []game.GameQuestion {
	questionList := make([]game.GameQuestion, 0)

	for _, gameQuestion := range gameQuestionList {
		questionList = append(questionList, game.GameQuestion{
			QuestionID: gameQuestion.QuestionID,
			Score:      gameQuestion.Score,
		})
	}

	return questionList
}

func ToSignUpAccountList(signUpAccountList []SignUpAccount) []game.SignUpAccount {
	signUpAccountListModel := make([]game.SignUpAccount, 0)

	for _, signUpAccount := range signUpAccountList {
		signUpAccountListModel = append(signUpAccountListModel, game.SignUpAccount{
			AccountID:  signUpAccount.AcountID,
			SignUpTime: signUpAccount.SignUpTime.Unix(),
		})
	}

	return signUpAccountListModel
}
