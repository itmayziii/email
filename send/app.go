package send

import (
	"gocloud.dev/blob"
	"gocloud.dev/blob/memblob"
	"io"
	"log"
)

// App defines the dependencies the application uses.
type App struct {
	// flusher provides an opportunity to flush any buffers prior to the [EmailEvent] function ending.
	flusher Flusher
	// infoLogger is meant to log "info" severity related events. This is a log.Logger instance because nobody can
	// agree on what a logging interface should look like so the easiest decision for this package to make is to rely
	// on the [Go standard logger].
	//
	// [Go standard logger]: https://pkg.go.dev/log#Logger
	infoLogger *log.Logger
	// infoLogger is meant to log "error" severity related events. This is a log.Logger instance because nobody can
	// agree on what a logging interface should look like so the easiest decision for this package to make is to rely
	// on the [Go standard logger].
	//
	// [Go standard logger]: https://pkg.go.dev/log#Logger
	errorLogger *log.Logger
	// fileStorage is a [blob.Bucket] which could represent a cloud storage service or a local filesystem.
	// [blob.Bucket]: https://gocloud.dev/howto/blob/
	fileStorage *blob.Bucket
	// domainSenders maps domains to Sender implementations. This could be sendgrid, mailgun, etc...
	// This allows flexibility to choose how to send emails per domain. A [Sender] is chosen from this map based on
	// the [Sender] email address. i.e. no-reply@google.com -> google.com is the domain.
	domainSenders map[string]Sender
}

// NewApp is a constructor for [App] which utilizes the [options pattern].
//
// [options pattern]: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
func NewApp(opts ...AppOption) *App {
	app := &App{
		domainSenders: make(map[string]Sender),
	}

	for _, opt := range opts {
		opt(app)
	}

	if app.flusher == nil {
		app.flusher = noopFlusher{}
	}

	if app.infoLogger == nil {
		noopLogger := log.New(io.Discard, "", 0)
		app.infoLogger = noopLogger
	}
	if app.errorLogger == nil {
		noopLogger := log.New(io.Discard, "", 0)
		app.errorLogger = noopLogger
	}

	if app.fileStorage == nil {
		memblob.OpenBucket(nil)
	}

	return app
}

type AppOption func(*App)

func AppWithFlusher(flusher Flusher) AppOption {
	return func(app *App) {
		app.flusher = flusher
	}
}

// AppWithLogger is a shortcut for [AppWithInfoLogger] and [AppWithErrorLogger] when the logger is the same between
// the two.
func AppWithLogger(logger *log.Logger) AppOption {
	return func(app *App) {
		app.infoLogger = logger
		app.errorLogger = logger
	}
}

// AppWithInfoLogger provides an option to supply an info severity logger.
func AppWithInfoLogger(logger *log.Logger) AppOption {
	return func(app *App) {
		app.infoLogger = logger
	}
}

// AppWithErrorLogger provides an option to supply an error severity logger.
func AppWithErrorLogger(logger *log.Logger) AppOption {
	return func(app *App) {
		app.errorLogger = logger
	}
}

// AppWithFileStorage provides an option to specify where email templates should be read from. This can be a
// [local file system], [GCS], [S3] or any of the other storages supported by the [blob package].
//
// [local file system]: https://gocloud.dev/howto/blob/#local
// [GCS]: https://gocloud.dev/howto/blob/#gcs
// [S3]: https://gocloud.dev/howto/blob/#s3
// [blob package]: https://gocloud.dev/howto/blob/
func AppWithFileStorage(fileStorage *blob.Bucket) AppOption {
	return func(app *App) {
		app.fileStorage = fileStorage
	}
}

// AppWithDomainSender associates a domain with a [Sender]. Domains will be matched with event supplied
// [EventData.Sender] i.e. Sender = no-reply@tommymay.dev: domain = tommymay.dev. The matching sender will be
// used to send the email.
func AppWithDomainSender(domain string, sender Sender) AppOption {
	return func(app *App) {
		app.domainSenders[domain] = sender
	}
}
