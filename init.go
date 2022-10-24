package main

import (
	"os"
	"strings"
	"time"

	"github.com/codingconcepts/env"
	"github.com/joho/godotenv"

	"github.com/gofrs/uuid"

	"github.com/sirupsen/logrus"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logger.Warn(".env file not found")
	}

	if err := env.Set(&cfg); err != nil {
		logger.Fatalf("%+v", err)
	}

	localTime = time.Now()

	if tzLocation, err := time.LoadLocation(cfg.Timezone); err == nil {
		localTime = localTime.In(tzLocation)
	}
	cfg.Timezone, tzOffset = localTime.Zone()

	if _, instGUIDFound := os.LookupEnv("CF_INSTANCE_GUID"); !instGUIDFound {
		cfg.InstanceGUID = uuid.Must(uuid.NewV4()).String()
	}

	switch strings.ToUpper(cfg.LogLevel) {
	case "PANIC":
		logger.Level = logrus.PanicLevel
	case "FATAL":
		logger.Level = logrus.FatalLevel
	case "ERROR":
		logger.Level = logrus.ErrorLevel
	case "WARNING":
		logger.Level = logrus.WarnLevel
	case "INFO":
		logger.Level = logrus.InfoLevel
	case "DEBUG":
		logger.Level = logrus.DebugLevel
	case "TRACE":
		logger.Level = logrus.TraceLevel
	default:
		logger.Level = logrus.InfoLevel
	}

	logFullTimestamp := false
	if cfg.DebugLevel > 0 {
		logFullTimestamp = true
	}

	switch strings.ToUpper(cfg.LogAs) {
	case "TEXT":
		logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: logFullTimestamp})
	case "JSON":
		logger.SetFormatter(&logrus.JSONFormatter{})
	default:
		logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: logFullTimestamp})
	}
	logger.Out = os.Stdout

	if cfg.LogDetails {
		log = logger.WithFields(logrus.Fields{"app": cfg.ApplicationName, "version": cfg.ApplicationVersion, "guid": cfg.InstanceGUID})
	}

	if _, bucketNameFound := os.LookupEnv("BUCKET_NAME"); !bucketNameFound {
		logger.Fatal("BUCKET_NAME env not set")
	}

	log.WithFields(logrus.Fields{
		"app":        cfg.ApplicationName,
		"version":    cfg.ApplicationVersion,
		"guid":       cfg.InstanceGUID,
		"logLevel":   logger.Level,
		"debugLevel": cfg.DebugLevel,
		"logAs":      cfg.LogAs,
		"logDetails": cfg.LogDetails,
		"bucketName": cfg.BucketName,
		"numObjects": cfg.NumObjects,
	}).Infof("%s %s initialized", appName, appVersion)

	log.WithFields(logrus.Fields{
		"tzName": cfg.Timezone, "tzOffset": tzOffset,
	}).Debug("timezone")
}
