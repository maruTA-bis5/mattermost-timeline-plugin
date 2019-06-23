module github.com/maruTA-bis5/mattermost-timeline-plugin

go 1.12

require (
	github.com/deckarep/golang-set v1.7.1
	github.com/go-ldap/ldap v3.0.3+incompatible // indirect
	github.com/mattermost/mattermost-server v0.0.0-20190614191352-bfd66aa445a2
	github.com/nicksnyder/go-i18n v1.10.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.3.0
	google.golang.org/genproto v0.0.0-20190404172233-64821d5d2107 // indirect
	google.golang.org/grpc v1.20.0 // indirect
)

// Workaround for https://github.com/golang/go/issues/30831 and fallout.
replace github.com/golang/lint => github.com/golang/lint v0.0.0-20190227174305-8f45f776aaf1
