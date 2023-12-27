package api

import (
	"codesignal/datatypes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type TakeTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func TestTakeTestSuite(t *testing.T) {
	suite.Run(t, &TakeTestSuite{})
}

func (s *TakeTestSuite) SetupSuite() {
	s.router = gin.Default()
	limiter := datatypes.LoadFromConfig("config.json")
	s.router.GET(
		"/take",
		func(context *gin.Context) {
			CheckStatus(context, limiter)
		},
	)
	s.router.GET(
		"/user/:id",
		limiter.BuildFor(http.MethodGet, "/user/:id"),
		TakeApiGetUserId,
	)
	s.router.PATCH(
		"/user/:id",
		limiter.BuildFor(http.MethodPatch, "/user/:id"),
		PatchUserId,
	)
	s.router.POST(
		"/userinfo",
		limiter.BuildFor(http.MethodPost, "/userinfo"),
		PostUserInfo,
	)
}

func (s *TakeTestSuite) TestTakeGetRespondsWithHelloWorldAndHttpStatus200() {

	// 1 call one endpoint
	{
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/user/1234", nil)

		s.router.ServeHTTP(w, req)

		s.Assert().Equal(http.StatusOK, w.Result().StatusCode)

		responseBody, err := ioutil.ReadAll(w.Result().Body)
		s.Assert().NoError(err)
		s.Assert().Equal("Hello World with user id", string(responseBody))
	}
	// check status
	{
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/take?method=GET&route=/user/:id", nil)

		s.router.ServeHTTP(w, req)

		s.Assert().Equal(http.StatusOK, w.Result().StatusCode)

		responseBody, err := ioutil.ReadAll(w.Result().Body)
		s.Assert().NoError(err)
		s.Assert().Equal(`{"allow":true,"available_tokens":5}`, string(responseBody))
	}
}

func (s *TakeTestSuite) TestGetUserWIthIdOnce() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/user/1234", nil)

	s.router.ServeHTTP(w, req)

	s.Assert().Equal(http.StatusOK, w.Result().StatusCode)

	responseBody, err := ioutil.ReadAll(w.Result().Body)
	s.Assert().NoError(err)
	s.Assert().Equal("Hello World with user id", string(responseBody))

}

func (s *TakeTestSuite) TestGetUserWIthIdOnceWithRefill() {
	{
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPatch, "/user/1234", nil)

		s.router.ServeHTTP(w, req)

		s.Assert().Equal(http.StatusOK, w.Result().StatusCode)

		responseBody, err := ioutil.ReadAll(w.Result().Body)
		s.Assert().NoError(err)
		s.Assert().Equal("this is patch user id", string(responseBody))
	}
	// 6 req/min
	time.Sleep(11 * time.Second)
	{
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPatch, "/user/1234", nil)

		s.router.ServeHTTP(w, req)

		s.Assert().Equal(http.StatusOK, w.Result().StatusCode)

		responseBody, err := ioutil.ReadAll(w.Result().Body)
		s.Assert().NoError(err)
		s.Assert().Equal("this is patch user id", string(responseBody))
	}
	// check status
	{
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/take?method=PATCH&route=/user/:id", nil)

		s.router.ServeHTTP(w, req)

		s.Assert().Equal(http.StatusOK, w.Result().StatusCode)

		responseBody, err := ioutil.ReadAll(w.Result().Body)
		s.Assert().NoError(err)
		s.Assert().Contains(string(responseBody), `{"allow":true,"available_tokens":`)
	}
}
