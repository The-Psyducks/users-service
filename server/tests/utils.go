package tests

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"users-service/src/router"
// )

// func postValidUser(router *router.Router, snap User) (int, GetSnapResponse, error) {
// 	marshalledSnap, err := json.Marshal(snap)

// 	if err != nil {
// 		return 0, GetSnapResponse{}, err
// 	}

// 	req, err := http.NewRequest("POST", "/snaps", bytes.NewReader(marshalledSnap))

// 	if err != nil {
// 		return 0, GetSnapResponse{}, err
// 	}

// 	req.Header.Add("content-type", "application/json")
// 	recorder := httptest.NewRecorder()
// 	router.Engine.ServeHTTP(recorder, req)
// 	result := GetSnapResponse{}
// 	err = json.Unmarshal(recorder.Body.Bytes(), &result)

// 	if err != nil {
// 		return 0, GetSnapResponse{}, err
// 	}

// 	return recorder.Code, result, nil
// }

// func getSnapsInOrder(router *router.Router) (int, GetAllSnapsResponse, error) {
// 	req, err := http.NewRequest("GET", "/snaps", &bytes.Reader{})

// 	if err != nil {
// 		return 0, GetAllSnapsResponse{}, err
// 	}

// 	recorder := httptest.NewRecorder()
// 	router.Engine.ServeHTTP(recorder, req)
// 	result := GetAllSnapsResponse{}
// 	err = json.Unmarshal(recorder.Body.Bytes(), &result)

// 	if err != nil {
// 		return 0, GetAllSnapsResponse{}, err
// 	}
// 	return recorder.Code, result, nil
// }

// func getExistingSnapById(router *router.Router, id string) (int, GetSnapResponse, error) {
// 	req, err := http.NewRequest("GET", "/snaps/"+id, &bytes.Reader{})

// 	if err != nil {
// 		return 0, GetSnapResponse{}, err
// 	}

// 	recorder := httptest.NewRecorder()
// 	router.Engine.ServeHTTP(recorder, req)
// 	result := GetSnapResponse{}
// 	err = json.Unmarshal(recorder.Body.Bytes(), &result)

// 	if err != nil {
// 		return 0, GetSnapResponse{}, err
// 	}
// 	return recorder.Code, result, nil
// }

// func deleteExistingSnapById(router *router.Router, id string) (int, error) {
// 	req, err := http.NewRequest("DELETE", "/snaps/"+id, &bytes.Reader{})

// 	if err != nil {
// 		return 0, err
// 	}

// 	recorder := httptest.NewRecorder()
// 	router.Engine.ServeHTTP(recorder, req)

// 	return recorder.Code, nil
// }