package cmd

import (
	"context"
	"fmt"

	"github.com/src-d/sourced-ce/cmd/sourced/compose"
	"github.com/src-d/sourced-ce/cmd/sourced/compose/workdir"
)

type pruneCmd struct {
	Command `name:"prune" short-description:"Stop and remove components and resources" long-description:"Stops containers and removes containers, networks, volumes and configuration created by 'init' for the current working directory.\nTo delete resources for all working directories pass --all flag.\nImages are not deleted unless you specify the --images flag."`

	All    bool `short:"a" long:"all" description:"Remove containers and resources for all working directories. It will be ignored if a list of workdirs is passed."`
	Images bool `long:"images" description:"Remove docker images"`
	Args   struct {
		Workdirs []string `positional-arg-name:"workdir" description:"Workdir to be pruned"`
	} `positional-args:"yes"`
}

func (c *pruneCmd) Execute(args []string) error {
	current, _ := workdir.Active()
	var dirs []string
	var err error
	if len(c.Args.Workdirs) > 0 {
		for _, dir := range c.Args.Workdirs {
			err = workdir.ValidatePath(dir)
			if err != nil {
				fmt.Printf("ignored workdir '%s': %s\n", dir, err)
				continue
			}

			dirs = append(dirs, dir)
		}
	} else if c.All {
		dirs, err = workdir.ListPaths()
		if err != nil {
			return err
		}
	} else {
		return c.pruneActive()
	}

	for _, dir := range dirs {
		fmt.Println(dir)
		continue
		if err := workdir.SetActivePath(dir); err != nil {
			return err
		}

		if err = c.pruneActive(); err != nil {
			return err
		}

		if current == dir {
			current = ""
		}
	}

	return nil
	if current != "" {
		workdir.SetActivePath(current)
	}

	return nil
}

func (c *pruneCmd) pruneActive() error {
	a := []string{"down", "--volumes"}
	if c.Images {
		a = append(a, "--rmi", "all")
	}

	if err := compose.Run(context.Background(), a...); err != nil {
		return err
	}

	dir, err := workdir.ActivePath()
	if err != nil {
		return err
	}

	if err := workdir.RemovePath(dir); err != nil {
		return err
	}

	return workdir.UnsetActive()
}

func init() {
	rootCmd.AddCommand(&pruneCmd{})
}
