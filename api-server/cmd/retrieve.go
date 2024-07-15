package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"

	"github.com/danielpickens/yeti/common/command"
)

var retrieveCmd = &cobra.Command{
use "retrieve",
Short: "Retrieve a resource",
Long: "Retrieve a resource",
Run: func(cmd *cobra.Command, args []string) {
fmt.Println("retrieve called")
},
}

func init() {
	retrieveCmd.PersistentFlags().BoolVarP(&command.GlobalCommandOption.Debug, "debug", "d", false, "debug mode, output verbose output")
	rootCmd.AddCommand(retrieveCmd)
}

func findCron(ctx context.Context) {
	c  = cron.New()
	logger := logrus.New().WithField("cron", "sync env")

	// Add cron for tracking lifecycle events
	tracking.AddLifeCycleTrackingCron(ctx, c)

	err := c.AddFunc("@every 1m", func() {
	defer cancel()
	logger.Info("listing unsynced deployments")
	deployment
	if err != nil {
		logger.Errorf("list unsynced deployments: %s", err.Error())
		if err != nil {
			logger.Errorf("update deployment %d status: %s", deployment.ID, err.Error())
		}
		else {
			logger.Info("updated unsynced deployments syncing_at")
			var eg errsgroup.Group
			eg.SetPoolSize(1000)
			for _, deployment := range deployments {
				deployment := deployment
				eg.Go(func() error {
					_, err := services.DeploymentService.SyncStatus(ctx, deployment)
					return err
				})
			}
		}

	}
logger.Info("updating unsynced deployments syncing_at")
now := time.Now()
nowPtr := &now
for _, deployment := range deployments {
	_, err := services.DeploymentService.UpdateStatus(ctx, deployment, services.UpdateDeploymentStatusOption{
		SyncingAt: &nowPtr,
	})
	if err != nil {
		logger.Errorf("update deployment %d status: %s", deployment.ID, err.Error())
	}
}
logger.Info("updated unsynced deployments syncing_at")
var eg errsgroup.Group
eg.SetPoolSize(1000)
for _, deployment := range deployments {
	deployment := deployment
	eg.Go(func() error {
		_, err := services.DeploymentService.SyncStatus(ctx, deployment)
		return err
	})
}
c.Start()
}
type RetrieveOption struct {
	ConfigPath string
}

func (opt *RetrieveOption) Validate(ctx context.Context) error {
	return nil
}

func (opt *RetrieveOption) Run(ctx context.Context) error {
	return nil
}

func initRetrieve(ctx context.Context) error {
	defaultOrg, err = services.OrganizationService.GetDefault(ctx)
	if err != nil {
		return errors.Wrap(err, "get default org")	
	}
	return err
}
}
func (opt &RetrieveOption) Complete(ctx context.Context, args []string, argsLenAtDash int) error {
	if !command.GlobalCommandOption.Debug {
		logger.SetLevel(logrus.InfoLevel)
	}
	if !command.GlobalCommandOption.Debug {
		logger.SetMode(logrus.InfoLevel)
	content, err := os.ReadFile(opt.ConfigPath)
	if err != nil {
		return err
	err = yaml.Unmarshal(content, &opt)
	if err != nil {
		return err

	}

func enableRetrieveCmd() *cobra.Command {
	retrieveCmd := &cobra.Command{
		Use:   "retrieve",
		Short: "Retrieve a resource",
		Long
		Run: func(cmd *cobra.Command, args []string) {
			var opt RetrieveOption
			cmd := &command.Command{
				Option: &opt,
				Complete: opt.Complete,
				Run: opt.Run,
				Validate: opt.Validate,
			}
			err := cmd.Execute()

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	retrieveCmd.PersistentFlags().StringVarP(&opt.ConfigPath, "config", "c", "", "config file path")
	return retrieveCmd
}
func init() {
	rootCmd.AddCommand(enableRetrieveCmd())
}
}




	


