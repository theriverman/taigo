module github.com/theriverman/taigo/cli

replace github.com/theriverman/taigo => ../

replace github.com/theriverman/taigo/cli/passwordbasedencryption => ./passwordbasedencryption

go 1.19

require (
	github.com/denisbrodbeck/machineid v1.0.1
	github.com/theriverman/taigo v1.6.1
	github.com/theriverman/taigo/cli/passwordbasedencryption v0.0.0-00010101000000-000000000000
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/google/go-querystring v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/text v0.3.8 // indirect
)
