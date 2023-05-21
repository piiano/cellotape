package utils

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type Answer struct {
	OneAnswer      *string
	ManyAnswers    *[]string
	UltimateAnswer *int
}

func TestBadMultiType(t *testing.T) {
	testCases := []struct {
		name string
		mt   multiType
	}{
		{
			name: "non struct type",
			mt:   &MultiType[string]{},
		},
		{
			name: "non pointer field",
			mt: &MultiType[struct {
				A string
				B *int
			}]{},
		},
		{
			name: "empty struct",
			mt:   &MultiType[struct{}]{},
		},
		{
			name: "non unique field type",
			mt: &MultiType[struct {
				A *string
				B *string
			}]{},
		},
		{
			name: "non unique field type with anonymous structs",
			mt: &MultiType[struct {
				A *struct {
					A string
				}
				B *struct {
					A string
				}
			}]{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mtType := reflect.TypeOf(testCase.mt)
			require.True(t, IsMultiType(mtType))

			_, err := ExtractMultiTypeTypes(mtType)
			require.ErrorIs(t, err, ErrInvalidUseOfMultiType)

			err = json.Unmarshal([]byte(`""`), &testCase.mt)
			require.ErrorIs(t, err, ErrInvalidUseOfMultiType)

			_, err = json.Marshal(testCase.mt)
			require.ErrorIs(t, err, ErrInvalidUseOfMultiType)
		})
	}
}

func TestMultiTypeTypes(t *testing.T) {
	testCases := []struct {
		name  string
		mt    multiType
		types []reflect.Type
	}{
		{
			name: "answer",
			mt:   &MultiType[Answer]{},
			types: []reflect.Type{
				GetType[*string](),
				GetType[*[]string](),
				GetType[*int](),
			},
		},
		{
			name: "anonymous struct",
			mt: &MultiType[struct {
				A *bool
				B *string
				C *struct {
					C1 []string
				}
			}]{},
			types: []reflect.Type{
				GetType[*bool](),
				GetType[*string](),
				GetType[*struct{ C1 []string }](),
			},
		},
		{
			name: "array of single of same struct",
			mt: &MultiType[struct {
				A *struct {
					A string
				}
				B *[]struct {
					A string
				}
			}]{},
			types: []reflect.Type{GetType[*struct {
				A string
			}](), GetType[*[]struct {
				A string
			}]()},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mtType := reflect.TypeOf(testCase.mt)
			require.True(t, IsMultiType(mtType))
			types, err := ExtractMultiTypeTypes(mtType)
			require.NoError(t, err)
			require.ElementsMatch(t, testCase.types, types)
		})
	}
}

func TestMarshalMultiType(t *testing.T) {
	testCases := []struct {
		name   string
		input  MultiType[Answer]
		output string
	}{
		{
			name: "OneAnswer",
			input: MultiType[Answer]{
				Values: Answer{
					OneAnswer: Ptr("foo"),
				},
			},
			output: `"foo"`,
		},
		{
			name: "ManyAnswers",
			input: MultiType[Answer]{
				Values: Answer{
					ManyAnswers: &[]string{"foo", "bar"},
				},
			},
			output: `["foo", "bar"]`,
		},
		{
			name: "UltimateAnswer",
			input: MultiType[Answer]{
				Values: Answer{
					UltimateAnswer: Ptr(42),
				},
			},
			output: `42`,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			bytes, err := json.Marshal(&test.input)
			require.NoError(t, err)
			require.JSONEq(t, test.output, string(bytes))
		})
	}
}

func TestMarshalMultiTypeError(t *testing.T) {
	multi := &MultiType[Answer]{}

	_, err := json.Marshal(multi)
	require.Error(t, err)
	require.IsType(t, &json.MarshalerError{}, err)

	multi = &MultiType[Answer]{
		Values: Answer{
			UltimateAnswer: Ptr(42),
			OneAnswer:      Ptr("foo"),
		},
	}

	_, err = json.Marshal(multi)
	require.Error(t, err)
	require.IsType(t, &json.MarshalerError{}, err)
}

func TestUnmarshalMultiType(t *testing.T) {
	multi := &MultiType[Answer]{}

	err := json.Unmarshal([]byte(`["foo", "bar"]`), multi)
	require.NoError(t, err)
	require.NotNil(t, multi.Values.ManyAnswers)
	require.ElementsMatch(t, []string{"foo", "bar"}, *multi.Values.ManyAnswers)
	require.Nil(t, multi.Values.OneAnswer)
	require.Nil(t, multi.Values.UltimateAnswer)

	multi = &MultiType[Answer]{}

	err = json.Unmarshal([]byte(`"foo"`), multi)
	require.NoError(t, err)
	require.NotNil(t, multi.Values.OneAnswer)
	require.Equal(t, "foo", *multi.Values.OneAnswer)
	require.Nil(t, multi.Values.ManyAnswers)
	require.Nil(t, multi.Values.UltimateAnswer)

	multi = &MultiType[Answer]{}

	err = json.Unmarshal([]byte(`42`), multi)
	require.NoError(t, err)
	require.NotNil(t, multi.Values.UltimateAnswer)
	require.Equal(t, 42, *multi.Values.UltimateAnswer)
	require.Nil(t, multi.Values.ManyAnswers)
	require.Nil(t, multi.Values.OneAnswer)

	multi = &MultiType[Answer]{}

	err = json.Unmarshal([]byte(`true`), multi)
	require.Error(t, err)
	require.IsType(t, &json.UnmarshalTypeError{}, err)
	require.Equal(t, &json.UnmarshalTypeError{
		Value:  "bool",
		Type:   reflect.TypeOf(multi),
		Offset: 4,
	}, err)

	multi = &MultiType[Answer]{}

	err = json.Unmarshal([]byte(`}`), multi)
	require.Error(t, err)
	require.IsType(t, &json.SyntaxError{}, err)
}
