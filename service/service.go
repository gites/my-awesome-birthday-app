package service

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// BDayService is the struct for service config
type BDayService struct {
	DbHost        string `envconfig:"DB_HOST" required:"true"`
	DbUser        string `envconfig:"DB_USER" required:"true"`
	DbPass        string `envconfig:"DB_PASS" required:"true"`
	DbPort        string `envconfig:"DB_PORT" required:"true"`
	DbName        string `envconfig:"DB_NAME" required:"true"`
	HttpPort      string `envconfig:"HTTP_PORT" default:"8080"`
	MetricsPort   string `envconfig:"METRICS_PORT" default:"9090"`
	MigrationsDir string `envconfig:"MIGRATIONS_DIR" default:"/migrations"`
	ShutdownWait  int    `envconfig:"HTTP_SHUTDOWN_WAIT" default:"1"`
	Debug         bool   `envconfig:"DEBUG" default:"false"`
	db            *sql.DB
	isReady       bool
}

// NewBDayService creates new instance of the services
func NewBDayService() (*BDayService, error) {
	newSvc := &BDayService{
		isReady: false,
		Debug:   true,
	}
	if err := envconfig.Process("", newSvc); err != nil {
		return nil, err
	}
	log.SetFormatter(&log.JSONFormatter{})
	if newSvc.Debug {
		log.SetLevel(logrus.DebugLevel)
	}
	log.WithField("context", "NewBDayService").Infof("initializing...")

	return newSvc, nil
}

// ConnectDB connects to db and store db connection in service object
func (b *BDayService) ConnectDB() error {

	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		b.DbHost, b.DbUser, b.DbPass, b.DbName)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return err
	}
	b.db = db
	return nil
}

// Run runs the services
func (b *BDayService) Run() {
	log.WithField("context", "Run").Infof("running...")

	defer b.db.Close()

	// start metrics / health server
	go b.runMetricsServer()

	// start http server and process
	b.runServer()
}

// waitForExitSignal waits for SIGINT, SIGTERM signals and shutodown http servers
func (b *BDayService) waitForExitSignal(server *http.Server, name string) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-signalChan:
			log.WithField("context", name).Infof("shutting down...")
			b.isReady = false
			time.Sleep(time.Duration(b.ShutdownWait) * time.Second)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				log.WithField("context", name).Errorf("%s", err)
				os.Exit(1)
			}
			os.Exit(0)
		}
	}
}

// RunMigrations runs db migrations
func RunMigrations(svc *BDayService) error {
	src := fmt.Sprintf("file://%s", svc.MigrationsDir)
	dst := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		svc.DbUser, svc.DbPass, svc.DbHost, svc.DbPort, svc.DbName)
	m, err := migrate.New(src, dst)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		return err
	}
	return nil
}
