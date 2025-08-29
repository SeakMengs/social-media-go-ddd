package dto

type (
	NewUser struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	UserLogin struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	CommonUserAggregate struct {
		Followed       bool `json:"followed"`
		FollowerCount  int  `json:"followerCount"`
		FollowingCount int  `json:"followingCount"`
	}
)
