module git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api

require (
	git.m2mfacil.com.br/golang/go-logging-package-level/pkg/logging v1.0.0
	github.com/BurntSushi/toml v0.0.0-20170626110600-a368813c5e64 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-sql-driver/mysql v0.0.0-20180719071942-99ff426eb706
	github.com/hazelcast/hazelcast-go-client v0.0.0-20180329145339-da1bb5ba3442
	github.com/jinzhu/configor v0.0.0-20180308034956-ec34f328956f
	github.com/julienschmidt/httprouter v0.0.0-20180222160526-d18983907793
	github.com/natefinch/lumberjack v0.0.0-20150411233054-6d54cbc97d7e // indirect
	github.com/rileyr/middleware v0.0.0-20171218202914-adaf4755ca16
	github.com/streadway/amqp v0.0.0-20180315184602-8e4aba63da9f // indirect
	google.golang.org/appengine v1.1.0 // indirect
	gopkg.in/mgo.v2 v2.0.0-20160818020120-3f83fa500528
	gopkg.in/natefinch/lumberjack.v2 v2.0.0-20170531160350-a96e63847dc3
	gopkg.in/yaml.v2 v2.2.0 // indirect
)

//replace (
//	git.m2mfacil.com.br/golang/go-logging-package-level/pkg/logging => ../_go-logging-package-level/pkg/logging
//)

//replace git.m2mfacil.com.br/golang/go-logging-package-level/pkg/logging => ../go-logging-package-level/pkg/logging
replace git.m2mfacil.com.br/golang/go-logging-package-level/pkg/logging => ../go-logging-package-level/pkg/logging
