module github.com/containerssh/backend/v2

go 1.16

require (
	github.com/aws/aws-sdk-go v1.38.45 // indirect
	github.com/containerssh/backend v1.0.0
	github.com/containerssh/configuration/v2 v2.0.1
	github.com/containerssh/docker/v2 v2.0.1
	github.com/containerssh/geoip v1.0.0
	github.com/containerssh/kubernetes/v2 v2.0.1
	github.com/containerssh/log v1.1.6
	github.com/containerssh/metrics v1.0.0
	github.com/containerssh/security v1.0.0
	github.com/containerssh/service v1.0.0
	github.com/containerssh/sshproxy v1.0.0
	github.com/containerssh/sshserver v1.0.0
	github.com/containerssh/structutils v1.0.0
	github.com/stretchr/testify v1.8.0
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
	golang.org/x/net v0.0.0-20210521195947-fe42d452be8f // indirect
	golang.org/x/sys v0.0.0-20210521203332-0cec03c779c1 // indirect
	google.golang.org/genproto v0.0.0-20210521181308-5ccab8a35a9a // indirect
	k8s.io/utils v0.0.0-20210521133846-da695404a2bc // indirect
)

// Fixes CVE-2020-9283
replace (
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2 => golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550 => golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83
	golang.org/x/crypto v0.0.0-20200220183623-bac4c82f6975 => golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 => golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83
)

// Fixes CVE-2020-14040
replace (
	golang.org/x/text v0.3.0 => golang.org/x/text v0.3.3
	golang.org/x/text v0.3.1 => golang.org/x/text v0.3.3
	golang.org/x/text v0.3.2 => golang.org/x/text v0.3.3
)

// Fixes CVE-2019-11254
replace (
	gopkg.in/yaml.v2 v2.2.0 => gopkg.in/yaml.v2 v2.2.8
	gopkg.in/yaml.v2 v2.2.1 => gopkg.in/yaml.v2 v2.2.8
	gopkg.in/yaml.v2 v2.2.2 => gopkg.in/yaml.v2 v2.2.8
	gopkg.in/yaml.v2 v2.2.3 => gopkg.in/yaml.v2 v2.2.8
	gopkg.in/yaml.v2 v2.2.4 => gopkg.in/yaml.v2 v2.2.8
	gopkg.in/yaml.v2 v2.2.5 => gopkg.in/yaml.v2 v2.2.8
	gopkg.in/yaml.v2 v2.2.6 => gopkg.in/yaml.v2 v2.2.8
	gopkg.in/yaml.v2 v2.2.7 => gopkg.in/yaml.v2 v2.2.8
)
