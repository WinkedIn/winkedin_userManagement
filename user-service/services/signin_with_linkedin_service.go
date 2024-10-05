package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/winkedin/user-service/constants"
	"github.com/winkedin/user-service/interfaces"
	"github.com/winkedin/user-service/logger"
	"github.com/winkedin/user-service/models"
	"github.com/winkedin/user-service/store/user"
	"github.com/winkedin/user-service/types"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

const (
	// EndpointProfile is the endpoint for profile api.
	EndpointProfile      = "https://api.linkedin.com/v2/me?projection=(id,firstName,lastName,vanityName,localizedHeadline,localizedFirstName,localizedLastName,localizedHeadline,headline,profilePicture(displayImage~:playableStreams))"
	EndpointEmailAddress = "https://api.linkedin.com/v2/clientAwareMemberHandles?q=members&projection=(elements*(primary,type,handle~))"
)

var (
	oauth2AuthURL  = "https://www.linkedin.com/oauth/v2/authorization"
	oauth2TokenURL = "https://www.linkedin.com/oauth/v2/accessToken"
)

type SignInWithLinkedInServiceImpl struct {
	rdb          *redis.Client
	userStore    user.UserStore
	loginSvc     interfaces.LoginService
	ClientId     string
	ClientSecret string
	RedirectURL  string
}

type HttpClient struct {
	http.Client
	OAuth2AccessToken oauth2.Token
}

type LinkedInClient struct {
	// ClientID is the api key client's ID.
	clientID string
	// ClientSecret is the api key client's secret.
	clientSecret string
	// Scopes is the list of scopes that the client will request.
	scopes []string
	// redirectURL is the URL that the user will be redirected to after
	// authenticating with LinkedIn in the GetAuthURL url response.
	redirectURL string
}

func NewSignInWithLinkedInService(db *gorm.DB, rdb *redis.Client, loginSvc interfaces.LoginService) interfaces.SignInWithLinkedInService {
	v := GetConfig(*ConfigFilePath)
	return &SignInWithLinkedInServiceImpl{
		rdb:          rdb,
		userStore:    user.NewUserStore(db),
		loginSvc:     loginSvc,
		ClientId:     v.GetString("linkedin.client_id"),
		ClientSecret: v.GetString("linkedin.client_secret"),
		RedirectURL:  v.GetString("linkedin.redirect_url"),
	}

}

// NewLinkedInClient creates a new LinkedInClient.
func NewLinkedInClient(clientID, clientSecret string, scopes []string, redirectURL string) *LinkedInClient {
	return &LinkedInClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		scopes:       scopes,
		redirectURL:  redirectURL,
	}
}

// getOAuth2Config returns the oauth2 config for the LinkedIn client.
func (c *LinkedInClient) getOauth2Config(oAuth2Endpoint oauth2.Endpoint) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     c.clientID,
		ClientSecret: c.clientSecret,
		Scopes:       c.scopes,
		Endpoint:     oAuth2Endpoint,
		RedirectURL:  c.redirectURL,
	}
}

// GetAuthURL returns the URL to the LinkedIn login page
// The state is a string that will be returned to the redirect URL, so it can be used to prevent CSRF attacks
func (c *LinkedInClient) GetAuthURL() string {
	oauth2Config := c.getOauth2Config(oauth2.Endpoint{
		AuthURL:   oauth2AuthURL,
		TokenURL:  oauth2TokenURL,
		AuthStyle: oauth2.AuthStyleInParams,
	})

	// Generate a random state string for CSRF protection
	state := GenerateOTP(4)
	url := oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return url
}

// GetClient will exchange the code for an access token and
// use it to create a new client with an authorized http client
func (c *LinkedInClient) GetClient(ctx context.Context, code string) (*HttpClient, error) {
	logger.LogFunctionPointWithContext(ctx, constants.LogFunctionEntry)
	defer logger.LogFunctionPointWithContext(ctx, constants.LogFunctionExit)
	oauth2Config := c.getOauth2Config(oauth2.Endpoint{
		AuthURL:   oauth2AuthURL,
		TokenURL:  oauth2TokenURL,
		AuthStyle: oauth2.AuthStyleInParams,
	})

	token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	client := oauth2Config.Client(ctx, token)
	newHttpClient := &HttpClient{*client, *token}
	return newHttpClient, nil
}

func (s *SignInWithLinkedInServiceImpl) GetLinkedInProfileAndLogin(ctx context.Context, code string) (string, error) {
	logger.LogFunctionPointWithContext(ctx, constants.LogFunctionEntry)

	// check if LinkedIn client ID, secret, redirect URL are present in config
	if s.ClientId == "" || s.ClientSecret == "" || s.RedirectURL == "" {
		logger.LogErrorWithContext(ctx, "LinkedIn client ID, secret, redirect URL not found in config")
		return "", fmt.Errorf("LinkedIn client ID, secret, redirect URL not found in config")
	}

	// new LinkedIn client
	client := NewLinkedInClient(s.ClientId, s.ClientSecret, []string{"r_liteprofile", "r_basicprofile", "r_emailaddress"}, s.RedirectURL)

	// new http client to make requests
	httpClient, err := client.GetClient(ctx, code)
	if err != nil {
		logger.LogErrorWithContext(ctx, err.Error())
		return "", fmt.Errorf("failed to create http client from LinkedIn client: %v", err)
	}

	// TODO - following code will fetch user's email address -> check if user exists -> if not, create user (saving all data possible to get from LinkedIn) -> login user ->
	// TODO - check background info and compare db data and new data -> update db data if necessary -> login user
	// fetch LinkedIn email address
	emailResp, err := httpClient.Get(EndpointEmailAddress)
	if err != nil {
		logger.LogErrorWithContext(ctx, err.Error())
		return "", fmt.Errorf("failed to get LinkedIn email address: %v", err)
	}
	defer emailResp.Body.Close()

	emailStruct := types.LinkedInEmailAddress{}
	err = json.NewDecoder(emailResp.Body).Decode(&emailStruct)
	if err != nil {
		logger.LogErrorWithContext(ctx, err.Error())
		return "", fmt.Errorf("failed to decode LinkedIn email address response: %v", err)
	}

	emailAddress := emailStruct.Elements[0].HandleTilde.EmailAddress

	// fetch LinkedIn profile
	profileResp, err := httpClient.Get(EndpointProfile)
	if err != nil {
		logger.LogErrorWithContext(ctx, err.Error())
		return "", fmt.Errorf("failed to get LinkedIn profile: %v", err)
	}
	defer profileResp.Body.Close()

	// decode LinkedIn profile response
	profileStruct := types.LinkedInProfile{}
	err = json.NewDecoder(profileResp.Body).Decode(&profileStruct)
	if err != nil {
		logger.LogErrorWithContext(ctx, err.Error())
		return "", fmt.Errorf("failed to decode LinkedIn profile response: %v", err)
	}

	// check if user exists
	_, userExists, err := s.userStore.GetUserByEmail(ctx, emailAddress)
	if err != nil {
		logger.LogErrorWithContext(ctx, err.Error())
		return "", fmt.Errorf("failed to check if user exists: %v", err)
	}

	// create user, login and return if user doesn't exist else update user and login
	if !userExists {
		_, err = s.userStore.CreateUserFromLinkedInProfile(ctx, &profileStruct, emailAddress)
		if err != nil {
			logger.LogErrorWithContext(ctx, err.Error())
			return "", fmt.Errorf("failed to create user from LinkedIn profile: %v", err)
		}

		// login user
		token, err := s.loginSvc.Login(ctx, emailAddress, httpClient.OAuth2AccessToken.AccessToken)
		if err != nil {
			logger.LogErrorWithContext(ctx, err.Error())
			return "", fmt.Errorf("failed to login user: %v", err)
		}

		return token, nil
	}

	// compare db data and new data
	// fetch user from db
	user, _, err := s.userStore.GetUserByEmail(ctx, emailAddress)
	if err != nil {
		logger.LogErrorWithContext(ctx, err.Error())
		return "", fmt.Errorf("failed to get user by email: %v", err)
	}

	// hardcoded user check for now TODO - replace with linkedin data once we have access to paid endpoints
	dummyUser := &models.User{
		ID:          "123",
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "johndoe@linkedin.com",
		DateOfBirth: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		Gender:      "male",
		LinkedInID:  "123",
		CompanyName: "LinkedIn",
		JobTitle:    "Software Engineer",
		Education:   "Bachelors",
	}

	// compare user data
	if dummyUser.Education != user.Education {
		// update user
		updatedUser := &models.User{
			ID:             user.ID,
			FirstName:      profileStruct.LocalizedFirstName,
			LastName:       profileStruct.LocalizedLastName,
			Email:          emailAddress,
			DateOfBirth:    dummyUser.DateOfBirth,
			Gender:         dummyUser.Gender,
			LinkedInID:     profileStruct.ID,
			CompanyName:    dummyUser.CompanyName,
			JobTitle:       dummyUser.JobTitle,
			Education:      dummyUser.Education,
			ProfilePicture: profileStruct.ProfilePicture.DisplayImage,
		}
		err = s.userStore.UpdateUser(ctx, updatedUser)
		if err != nil {
			logger.LogErrorWithContext(ctx, err.Error())
			return "", fmt.Errorf("failed to update user from LinkedIn profile: %v", err)
		}
	}

	// login user
	token, err := s.loginSvc.Login(ctx, emailAddress, httpClient.OAuth2AccessToken.AccessToken)
	if err != nil {
		logger.LogErrorWithContext(ctx, err.Error())
		return "", fmt.Errorf("failed to login user: %v", err)
	}

	defer logger.LogFunctionPointWithContext(ctx, constants.LogFunctionExit)
	return token, nil
}

// once user signs in using linkedin share the token with frontend and Fe will monitor the TTL for user login
// use context in Db queries
// remove the login service
