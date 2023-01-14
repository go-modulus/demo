package temporal

import (
	"context"
	"demo/internal/logger"
	"fmt"
	"github.com/spf13/cobra"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/fx"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"reflect"
)

type Worker struct {
	logger     logger.Logger
	temporal   client.Client
	activities []Activity
	workflows  []Workflow
}

type WorkersParams struct {
	fx.In

	Logger     logger.Logger
	Temporal   client.Client
	Activities []Activity `group:"temporal.activities"`
	Workflows  []Workflow `group:"temporal.workflows"`
}

func NewWorker(params WorkersParams) *Worker {
	return &Worker{
		logger:     params.Logger,
		temporal:   params.Temporal,
		activities: params.Activities,
		workflows:  params.Workflows,
	}
}

func (w *Worker) Command() *cobra.Command {
	return &cobra.Command{
		Use:  "worker",
		Args: cobra.NoArgs,
		Run:  w.Run,
	}
}

func (w *Worker) Run(cmd *cobra.Command, args []string) {
	tw := worker.New(w.temporal, "rate", worker.Options{})

	lower := cases.Title(language.English, cases.NoLower)

	for _, a := range w.activities {
		r := reflect.TypeOf(a)

		Prefixes[r.Elem().Name()] = a.Prefix()

		// iterate over all methods of the activity struct
		for i := 0; i < r.NumMethod(); i++ {
			m := r.Method(i)

			if m.Name == "Prefix" {
				continue
			}

			name := fmt.Sprintf("%s.%s", a.Prefix(), lower.String(m.Name))

			tw.RegisterActivityWithOptions(m.Func.Interface(), activity.RegisterOptions{Name: name})
		}
	}

	for _, f := range w.workflows {
		opts := workflow.RegisterOptions{
			Name: f.Name(),
		}

		tw.RegisterWorkflowWithOptions(f.Do, opts)
	}

	err := tw.Run(worker.InterruptCh())

	if err != nil {
		w.logger.Error(context.Background(), "temporal error", logger.Field("err", err))
	}
}
