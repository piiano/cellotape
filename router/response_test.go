package router

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

//func TestExtractResponses(t *testing.T) {
//	tableTest(t, extractResponses, []testCase[reflect.Type, map[int]httpResponse]{
//		{in: nilType, out: map[int]httpResponse{}},
//		{in: reflect.TypeOf(struct{}{}), out: map[int]httpResponse{}},
//		{in: reflect.TypeOf(struct{ ok string }{}), out: map[int]httpResponse{}},
//		{in: reflect.TypeOf(struct {
//			OK string `Status:"200"`
//		}{}), out: map[int]httpResponse{200: {
//			declaredAt: reflect.TypeOf(struct {
//				OK string `Status:"200"`
//			}{}),
//			status:       200,
//			responseType: reflect.TypeOf(""),
//			fieldIndex:   []int{0},
//			isNilType:    false,
//		}}},
//		{in: reflect.TypeOf(struct {
//			OK     string `Status:"200"`
//			Teapot bool   `Status:"418"`
//		}{}), out: map[int]httpResponse{
//			200: {
//				status:       200,
//				responseType: reflect.TypeOf(""),
//				fieldIndex:   []int{0},
//			},
//			418: {
//				status:       418,
//				responseType: reflect.TypeOf(true),
//				fieldIndex:   []int{1},
//			},
//		}},
//		{in: reflect.TypeOf(struct{ OK string }{}), out: map[int]httpResponse{}},
//		{in: reflect.TypeOf(struct {
//			NoDashSupport int `Status:"-"`
//		}{}), out: map[int]httpResponse{}},
//	})
//}

func TestParseStatus(t *testing.T) {
	tableTestWithErr(t, parseStatus, []testCase[string, int]{
		{in: "100", out: 100},
		{in: "200", out: 200},
		{in: "204", out: 204},
		{in: "400", out: 400},
		{in: "404", out: 404},
		{in: "500", out: 500},
		{in: "599", out: 599},
		{in: "99", shouldErr: true},
		{in: "600", shouldErr: true},
		{in: "2e2", shouldErr: true},
		{in: "2E2", shouldErr: true},
		{in: "200.0", shouldErr: true},
		{in: "OK", shouldErr: true},
		{in: "2 00", shouldErr: true},
		{in: "-200", shouldErr: true},
		{in: "200 OK", shouldErr: true},
	})
}

type testCase[I any, O any] struct {
	in        I
	out       O
	shouldErr bool
}

func tableTestWithErr[I any, O any](t *testing.T, function func(I) (O, error), testCases []testCase[I, O]) {
	for _, test := range testCases {
		should := "pass"
		if test.shouldErr {
			should = "fail"
		}
		t.Run(fmt.Sprintf("should %s for input: %v", should, test.in), func(t *testing.T) {
			out, err := function(test.in)
			if test.shouldErr {
				require.NotNil(t, err)
			} else {
				assert.Equal(t, test.out, out)
				require.Nil(t, err)
			}
		})
	}
}

func tableTest[I any, O any](t *testing.T, function func(I) O, testCases []testCase[I, O]) {
	for _, test := range testCases {
		should := "pass"
		if test.shouldErr {
			should = "fail"
		}
		t.Run(fmt.Sprintf("should %s for input: %v", should, test.in), func(t *testing.T) {
			out := function(test.in)
			assert.Equal(t, test.out, out)
		})
	}
}
