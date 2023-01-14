package temporal

import (
	"demo/internal/cli"
	"fmt"
	"github.com/spf13/cobra"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/fx"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"reflect"
	"runtime"
	"strings"
)

type WorkerRegistry interface {
	Get(queueName string) worker.Worker
}

type TemporalRegistry struct {
}

func NewTemporalRegistry(client client.Client) *TemporalRegistry {
	return &TemporalRegistry{}
}

func (r *TemporalRegistry) Get(queueName string) worker.Worker {

	return nil
}

type Workflow interface {
	Name() string
	Do(ctx workflow.Context) error
}

type Activity interface {
	Prefix() string
}

var Prefixes = map[string]string{}

func GetActivityName(activity interface{}) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(activity).Pointer()).Name()
	elements := strings.Split(fullName, ".")
	structName := strings.Trim(elements[len(elements)-2], "(*)")
	fnName := strings.TrimSuffix(elements[len(elements)-1], "-fm")
	lower := cases.Title(language.English, cases.NoLower)

	return fmt.Sprintf("%s.%s", Prefixes[structName], lower.String(fnName))
}

func ExecuteActivity(ctx workflow.Context, activity interface{}, args ...interface{}) workflow.Future {
	activityName := GetActivityName(activity)

	return workflow.ExecuteActivity(ctx, activityName, args...)
}

func ProvideActivity[T Activity](activity interface{}) fx.Option {
	return fx.Provide(
		activity,
		fx.Annotate(
			func(a T) T { return a },
			fx.As(new(Activity)),
			fx.ResultTags(`group:"temporal.activities"`),
		),
	)
}

func ProvideWorkflow[T Workflow](workflow interface{}) fx.Option {
	return fx.Provide(
		workflow,
		fx.Annotate(
			func(a T) T { return a },
			fx.As(new(Workflow)),
			fx.ResultTags(`group:"temporal.workflows"`),
		),
	)
}

func Module() fx.Option {
	return fx.Module(
		"temporal",
		fx.Provide(
			NewConfig,
			NewLogger,

			func(config *Config, logger *Logger) (client.Client, error) {
				opts := client.Options{
					HostPort: config.HostPort,
					Logger:   logger,
				}

				return client.NewLazyClient(opts)
			},

			NewWorker,

			cli.ProvideCommand(
				func(worker *Worker) *cobra.Command {
					root := &cobra.Command{Use: "temporal"}

					root.AddCommand(worker.Command())

					return root
				},
			),
		),
	)
}
