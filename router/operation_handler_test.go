package router

import (
	"bytes"
	"io"
)

//
//func TestOperationFuncTypeExtraction(t *testing.T) {
//	fn := operationFunc[Nil, Nil, Nil, Nil](func(r Request[Nil, Nil, Nil]) (Response[Nil], error) { return Response[Nil]{}, nil })
//	types := fn.requestTypes()
//	assert.Equal(t, types.requestBody, nilType)
//	assert.Equal(t, types.pathParams, nilType)
//	assert.Equal(t, types.queryParams, nilType)
//}

//func TestOperationFuncAsGinHandler(t *testing.T) {
//	type responses struct {
//		Answer int `Status:"200"`
//	}
//	fn := operationFunc[Nil, Nil, Nil, responses](func(r Request[Nil, Nil, Nil]) (Response[responses], error) {
//		return Send(200, responses{Answer: 42})
//	})
//	spec, err := NewSpecFromData([]byte(`
//   { "paths": { "/abc": { "post": {
//     "operationId": "id",
//     "summary": "the ultimate answer to life the universe and everything",
//     "responses":{ "200": { "content": { "application/json": { "schema": { "type": "integer" } } } } }
//   } } } }`))
//	require.Nil(t, err)
//	router := NewOpenAPIRouter(spec).WithOperation("id", fn)
//	handler, err := router.AsHandler()
//	ts := httptest.NewServer(handler)
//	defer ts.Close()
//	log.Println(ts.URL)
//	// test valid httpRequest
//	resp, err := http.Post(fmt.Sprintf("%s/abc", ts.URL), "application/json", nil)
//	require.Nil(t, err)
//	res, err := toString(resp.Body)
//	require.Nil(t, err)
//	assert.Equal(t, "42", res)
//	assert.Equal(t, 200, resp.StatusCode)
//
//	// test bad httpRequest
//	resp, err = http.Post(ts.URL, "application/json", bytes.NewBufferString("{}"))
//	require.Nil(t, err)
//	assert.Equal(t, 400, resp.StatusCode)
//	res, err = toString(resp.Body)
//	require.Nil(t, err)
//	assert.Equal(t, res, `{"error":"expected httpRequest with no Body payload"}`)
//
//}

func toString(reader io.Reader) (string, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
