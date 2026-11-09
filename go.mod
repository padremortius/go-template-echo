module github.com/padremortius/go-template-echo

go 1.26.0

require (
	github.com/VictoriaMetrics/metrics v1.45.0
	github.com/go-playground/validator/v10 v10.30.6
	github.com/goccy/go-json v0.12.0
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/labstack/echo/v5 v5.4.4
	github.com/mvrilo/go-redoc v0.1.5
	github.com/mvrilo/go-redoc/echo v0.0.0-20250209151614-3a15e2c08553
	github.com/robfig/cron/v3 v3.0.1
	github.com/rs/zerolog v1.35.1
	github.com/ziflex/lecho/v4 v4.0.0
	gorm.io/driver/sqlite v1.6.0
	gorm.io/gorm v1.31.2
)

require (
	github.com/BurntSushi/toml v1.6.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.16 // indirect
	github.com/go-playground/locales v0.14.2 // indirect
	github.com/go-playground/universal-translator v0.18.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/leodido/go-urn v1.5.0 // indirect
	github.com/mattn/go-colorable v0.1.16 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/mattn/go-sqlite3 v1.14.53 // indirect
	github.com/stretchr/testify v1.12.1 // indirect
	github.com/valyala/fastrand v1.1.0 // indirect
	github.com/valyala/histogram v1.2.0 // indirect
	golang.org/x/crypto v0.58.0 // indirect
	golang.org/x/net v0.62.0 // indirect
	golang.org/x/sys v0.50.0 // indirect
	golang.org/x/text v0.44.0 // indirect
	golang.org/x/time v0.16.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace github.com/mvrilo/go-redoc => ./replaced/mvrilo/go-redoc

replace github.com/mvrilo/go-redoc/echo => ./replaced/mvrilo/go-redoc/echo
