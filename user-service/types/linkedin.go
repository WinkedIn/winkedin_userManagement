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
	ID                 string         `json:"id"`
	LocalizedFirstName string         `json:"localizedFirstName"`
	LocalizedLastName  string         `json:"localizedLastName"`
	FirstName          FirstName      `json:"firstName"`
	LastName           LastName       `json:"lastName"`
	ProfilePicture     ProfilePicture `json:"profilePicture"`
	Headline           Headline       `json:"headline"`
	LocalizedHeadline  string         `json:"localizedHeadline"`
	VanityName         string         `json:"vanityName"`
}

// FirstName - Struct for LinkedIn API first name response
type FirstName struct {
	Localized struct {
		EnUS string `json:"en_US"`
	} `json:"localized"`
	PreferredLocale struct {
		Country  string `json:"country"`
		Language string `json:"language"`
	} `json:"preferredLocale"`
}

// LastName - Struct for LinkedIn API last name response
type LastName struct {
	Localized struct {
		EnUS string `json:"en_US"`
	} `json:"localized"`
	PreferredLocale struct {
		Country  string `json:"country"`
		Language string `json:"language"`
	} `json:"preferredLocale"`
}

// ProfilePicture - Struct for LinkedIn API profile picture response
type ProfilePicture struct {
	DisplayImage     string `json:"displayImage"`
	DisplayImageFull struct {
		Paging struct {
			Count int   `json:"count"`
			Start int   `json:"start"`
			Links []any `json:"links"`
		} `json:"paging"`
		Elements []ProfilePictureElements `json:"elements"`
	} `json:"displayImage~"`
}

// Headline - Struct for LinkedIn API headline response
type Headline struct {
	Localized struct {
		EnUS string `json:"en_US"`
	} `json:"localized"`
	PreferredLocale struct {
		Country  string `json:"country"`
		Language string `json:"language"`
	} `json:"preferredLocale"`
}

// ProfilePictureElements - Struct for LinkedIn API profile picture elements response
type ProfilePictureElements struct {
	Artifact            string                     `json:"artifact"`
	AuthorizationMethod string                     `json:"authorizationMethod"`
	Data                ProfilePictureElementsData `json:"data"`
	Identifiers         []struct {
		Identifier                 string `json:"identifier"`
		Index                      int    `json:"index"`
		MediaType                  string `json:"mediaType"`
		File                       string `json:"file"`
		IdentifierType             string `json:"identifierType"`
		IdentifierExpiresInSeconds int    `json:"identifierExpiresInSeconds"`
	} `json:"identifiers"`
}

// ProfilePictureElementsData - Struct for LinkedIn API profile picture elements data response
type ProfilePictureElementsData struct {
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
}
