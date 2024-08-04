package game

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type MockEventQueue struct {
	event any
	err   error
}

func (m *MockEventQueue) EnQueue(event any) error {
	m.event = event
	return m.err
}

func (m *MockEventQueue) Queue() []any {
	return nil
}

func TestCalculateScore(t *testing.T) {
	// Create a mock EventQueue
	mockQueue := MockEventQueue{}

	// Create a Game instance
	g := &Game{
		GameID:         "1",
		SignUpUserList: []SignUpAccount{{AccountID: "1"}},
		QuestionList:   []GameQuestion{{QuestionID: "1", Score: 100}},
	}

	// Test case 1: Account does not exist
	cmd1 := CalculateScoreCmd{
		AccountID: "2",
	}
	err1 := g.calculate(&mockQueue, cmd1)
	require.ErrorIs(t, err1, ErrAccountNotFound, "Error: %w", err1)

	// Test case 2: Question does not exist
	cmd2 := CalculateScoreCmd{
		AccountID:  "1",
		QuestionID: "2",
	}
	err2 := g.calculate(&mockQueue, cmd2)
	require.ErrorIs(t, err2, ErrQuestionNotFound, "Error: %w", err2)

	// Test case 3: Normal calculation
	cmd3 := CalculateScoreCmd{
		AccountID:          "1",
		QuestionID:         "1",
		NumberFinishedAt:   10,
		LastMostFinishedAt: 5,
		TotalQuestion:      20,
	}
	err3 := g.calculate(&mockQueue, cmd3)
	require.NoError(t, err3, "Error: %w", err3)

	// Verify that the event was enqueued correctly
	expectedEvent := CalculateScoreEvent{
		GameID:     g.GameID,
		AccountID:  cmd3.AccountID,
		QuestionID: cmd3.QuestionID,
		Score:      25,
	}

	require.Equal(t, expectedEvent, mockQueue.event, "Unexpected event: %v, Expected: %v",
		mockQueue.event, expectedEvent)
}
