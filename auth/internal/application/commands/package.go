package commands

import "github.com/samber/do/v2"

var Package = do.Package(
	do.Lazy(NewCreateAccountCommand),
	do.Lazy(NewIssueJWTCommand),
	do.Lazy(NewRefreshJWTCommand),
)
