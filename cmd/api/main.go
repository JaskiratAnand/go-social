package main

import (
	"expvar"
	"runtime"
	"time"

	"github.com/JaskiratAnand/go-social/internal/auth"
	"github.com/JaskiratAnand/go-social/internal/db"
	"github.com/JaskiratAnand/go-social/internal/env"
	"github.com/JaskiratAnand/go-social/internal/mailer"
	"github.com/JaskiratAnand/go-social/internal/ratelimiter"
	"github.com/JaskiratAnand/go-social/internal/store"
	"github.com/JaskiratAnand/go-social/internal/store/cache"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const version = "0.0.1"
const QueryTimeoutDuration = time.Second * 5

//	@title			GoSocial API
//	@description	API for GoSocial, a social networking application.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/v1

//	@securityDefinitions.basic	BasicAuth
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization

func main() {
	cfg := config{
		addr:        env.GetString("ADDR", ":8080"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "https://localhost:4000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/go-social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		redisCfg: redisConfig{
			addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
			pw:      env.GetString("REDIS_PW", ""),
			db:      env.GetInt("Redis_DB", 0),
			enabled: env.GetBool("REDIS_ENABLED", false),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp:       (24 * time.Hour),
			emailAddr: env.GetString("EMAIL_ADDR", ""),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
		},
		auth: authConfig{
			basic: basicConfig{
				user: env.GetString("AUTH_BASIC_USER", "admin"),
				pass: env.GetString("AUTH_BASIC_PASS", "admin"),
			},
			token: tokenConfig{
				secret: env.GetString("JWT_AUTH_SECRET", "secret"),
				exp:    time.Hour * 24 * 3,
				iss:    "GoSocial",
			},
		},
		ratelimiter: ratelimiter.Config{
			RequestPerTimeFrame: env.GetInt("RATE_LIMITER_REQUESTS_COUNT", 20),
			TimeFrame:           time.Second * 5,
			Enabled:             env.GetBool("RATE_LIMITER_ENABLED", true),
		},
	}

	// logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// database
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infow("database connection pool established")

	defer db.Close()

	// redis cache
	var rdb *redis.Client
	if cfg.redisCfg.enabled {
		logger.Infow("redis cache connection established")
		rdb = cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.pw, cfg.redisCfg.db)
	}

	// rate limiter
	ratelimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.ratelimiter.RequestPerTimeFrame,
		cfg.ratelimiter.TimeFrame,
	)

	// stores
	store := store.New(db)
	cacheStorage := cache.NewRedisStorage(rdb)

	// mailer
	mailer := mailer.NewSendGrid(cfg.mail.sendGrid.apiKey, cfg.mail.emailAddr)

	// auth
	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.token.secret,
		cfg.auth.token.iss,
		cfg.auth.token.iss,
	)

	app := &application{
		config:        cfg,
		store:         store,
		cacheStorage:  cacheStorage,
		logger:        logger,
		mailer:        mailer,
		authenticator: jwtAuthenticator,
		rateLimiter:   ratelimiter,
	}

	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux := app.mount()

	logger.Fatal(app.run(mux))
}

func If[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}
