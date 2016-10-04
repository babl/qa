package main

import "github.com/urfave/cli"

func configureCli() (app *cli.App) {
	app = cli.NewApp()
	app.Usage = "Babl Quality Assurance"
	app.Version = Version
	app.Action = func(c *cli.Context) {
		httpServerBind := c.String("listen")
		kafkaBrokers := c.String("kafka-brokers")
		if len(httpServerBind) == 0 {
			httpServerBind = ":8888"
		}
		if len(kafkaBrokers) == 0 {
			kafkaBrokers = "queue.babl.sh:9092"
		}
		run(httpServerBind, kafkaBrokers, c.GlobalBool("debug"))
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "listen, l",
			Usage: "Host & port to listen to",
			Value: ":8888",
		},
		cli.StringFlag{
			Name:  "kafka-brokers, kb",
			Usage: "Comma separated list of kafka brokers",
			Value: "queue.babl.sh:9092",
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "Enable debug mode & verbose logging",
			EnvVar: "BABL_DEBUG",
		},
	}
	return
}
