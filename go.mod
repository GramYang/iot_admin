module iot_admin

go 1.15

require (
	apis v0.0.0
	github.com/GramYang/gylog v0.0.0-20201207013425-49af352abc7f
	github.com/gin-gonic/gin v1.6.3
	github.com/go-sql-driver/mysql v1.5.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
)

replace apis => ./apis
