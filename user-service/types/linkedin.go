package types

// ErrorResponse - Common struct for error response in LinkedIn API
type ErrorResponse struct {
	ServiceErrorCode int    `json:"serviceErrorCode"`
	Message          string `json:"message"`
	Status           int    `json:"status"`
}

// LinkedInEmailAddress - Struct for LinkedIn API email address response
type LinkedInEmailAddress struct {
	ErrorResponse
	Elements []struct {
		Handle      string `json:"handle"`
		HandleTilde struct {
			EmailAddress string `json:"emailAddress"`
		} `json:"handle~"`
	} `json:"elements"`
}

// LinkedInProfile - Struct for LinkedIn API profile response
type LinkedInProfile struct {
	ErrorResponse
	ID                 string `json:"id"`
	LocalizedFirstName string `json:"localizedFirstName"`
	LocalizedLastName  string `json:"localizedLastName"`
	FirstName          struct {
		Localized struct {
			EnUS string `json:"en_US"`
		} `json:"localized"`
		PreferredLocale struct {
			Country  string `json:"country"`
			Language string `json:"language"`
		} `json:"preferredLocale"`
	} `json:"firstName"`
	LastName struct {
		Localized struct {
			EnUS string `json:"en_US"`
		} `json:"localized"`
		PreferredLocale struct {
			Country  string `json:"country"`
			Language string `json:"language"`
		} `json:"preferredLocale"`
	} `json:"lastName"`
	ProfilePicture struct {
		DisplayImage     string `json:"displayImage"`
		DisplayImageFull struct {
			Paging struct {
				Count int   `json:"count"`
				Start int   `json:"start"`
				Links []any `json:"links"`
			} `json:"paging"`
			Elements []struct {
				Artifact            string `json:"artifact"`
				AuthorizationMethod string `json:"authorizationMethod"`
				Data                struct {
					ComLinkedinDigitalmediaMediaartifactStillImage struct {
						MediaType    string `json:"mediaType"`
						RawCodecSpec struct {
							Name string `json:"name"`
							Type string `json:"type"`
						} `json:"rawCodecSpec"`
						DisplaySize struct {
							Width  float64 `json:"width"`
							Uom    string  `json:"uom"`
							Height float64 `json:"height"`
						} `json:"displaySize"`
						StorageSize struct {
							Width  int `json:"width"`
							Height int `json:"height"`
						} `json:"storageSize"`
						StorageAspectRatio struct {
							WidthAspect  float64 `json:"widthAspect"`
							HeightAspect float64 `json:"heightAspect"`
							Formatted    string  `json:"formatted"`
						} `json:"storageAspectRatio"`
						DisplayAspectRatio struct {
							WidthAspect  float64 `json:"widthAspect"`
							HeightAspect float64 `json:"heightAspect"`
							Formatted    string  `json:"formatted"`
						} `json:"displayAspectRatio"`
					} `json:"com.linkedin.digitalmedia.mediaartifact.StillImage"`
				} `json:"data"`
				Identifiers []struct {
					Identifier                 string `json:"identifier"`
					Index                      int    `json:"index"`
					MediaType                  string `json:"mediaType"`
					File                       string `json:"file"`
					IdentifierType             string `json:"identifierType"`
					IdentifierExpiresInSeconds int    `json:"identifierExpiresInSeconds"`
				} `json:"identifiers"`
			} `json:"elements"`
		} `json:"displayImage~"`
	} `json:"profilePicture"`
	Headline struct {
		Localized struct {
			EnUS string `json:"en_US"`
		} `json:"localized"`
		PreferredLocale struct {
			Country  string `json:"country"`
			Language string `json:"language"`
		} `json:"preferredLocale"`
	} `json:"headline"`
	LocalizedHeadline string `json:"localizedHeadline"`
	VanityName        string `json:"vanityName"`
}
