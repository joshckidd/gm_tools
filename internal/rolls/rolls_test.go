package rolls

import "testing"

func TestParseRoll(t *testing.T) {
	cases := []struct {
		input    string
		expected []RollType
	}{
		{
			input: "3d8e",
			expected: []RollType{
				{
					Number:    3,
					Dice:      8,
					Aggregate: Sum,
					Signum:    Positive,
					Exploding: true,
				},
			},
		},
		{
			input: "2d6-1",
			expected: []RollType{
				{
					Number:    2,
					Dice:      6,
					Aggregate: Sum,
					Signum:    Positive,
					Exploding: false,
				},
				{
					Number:    1,
					Dice:      1,
					Aggregate: Sum,
					Signum:    Negative,
					Exploding: false,
				},
			},
		},
		{
			input: "min2d20",
			expected: []RollType{
				{
					Number:    2,
					Dice:      20,
					Aggregate: Min,
					Signum:    Positive,
					Exploding: false,
				},
			},
		},
		{
			input: "3d8e-2d6+1d4",
			expected: []RollType{
				{
					Number:    3,
					Dice:      8,
					Aggregate: Sum,
					Signum:    Positive,
					Exploding: true,
				},
				{
					Number:    2,
					Dice:      6,
					Aggregate: Sum,
					Signum:    Negative,
					Exploding: false,
				},
				{
					Number:    1,
					Dice:      4,
					Aggregate: Sum,
					Signum:    Positive,
					Exploding: false,
				},
			},
		},
	}

	for _, c := range cases {
		actual := ParseRoll(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("Expected slice of length %d, got slice of length %d", len(c.expected), len(actual))
		}

		for i := range actual {
			roll := actual[i]
			expectedRoll := c.expected[i]

			if roll != expectedRoll {
				t.Errorf("Expected roll %v, got roll %v", expectedRoll, roll)
			}
		}
	}

}
