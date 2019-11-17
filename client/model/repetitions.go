package model

import "github.com/mohae/deepcopy"

// RepetitionRange represents a single, inclusive range of repetition numbers,
// e.g. 1-4.
// An RepetitionRange can also represent a single ending number, e.g. 1-1.
type RepetitionRange struct {
	First int32
	Last  int32
}

// OnRepetitions wraps an event (something that implements ScoreUpdate) in order
// to specify that it should only occur on certain repetition numbers through a
// repeated pattern.
type OnRepetitions struct {
	Repetitions []RepetitionRange
	Event       ScoreUpdate
}

// AppliesTo returns true if a particular repetition number belongs to one of
// the specified repetition ranges.
func (or OnRepetitions) AppliesTo(repetition int32) bool {
	for _, r := range or.Repetitions {
		if r.First <= repetition && repetition <= r.Last {
			return true
		}
	}

	return false
}

// UpdateScore implements ScoreUpdate.UpdateScore by either updating the score
// with the event or doing nothing, depending on whether or not we are currently
// on a relevant repetition.
func (or OnRepetitions) UpdateScore(score *Score) error {
	for _, part := range score.CurrentParts {
		if or.AppliesTo(part.CurrentRepetition) {
			if err := or.Event.UpdateScore(score); err != nil {
				return err
			}
		}
	}

	return nil
}

// DurationMs implements ScoreUpdate.DurationMs by returning the duration of the
// event on the current repetition.
func (or OnRepetitions) DurationMs(part *Part) float32 {
	if or.AppliesTo(part.CurrentRepetition) {
		return or.Event.DurationMs(part)
	}

	return 0
}

// VariableValue implements ScoreUpdate.VariableValue by returning a version of
// the OnRepetitions where the value of the event is captured.
func (or OnRepetitions) VariableValue(score *Score) (ScoreUpdate, error) {
	result := deepcopy.Copy(or).(OnRepetitions)

	eventValue, err := or.Event.VariableValue(score)

	if err != nil {
		return nil, err
	}

	result.Event = eventValue

	return result, nil
}
