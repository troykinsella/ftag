package main

import (
	"errors"
	"fmt"
	"github.com/troykinsella/ftag/tagmap"
	"github.com/urfave/cli"
	"os"
	"path/filepath"
)

const (
	AppName = "ftag"

	optTagMap     = "m"
	optTagMapLong = "tag-map"

	defaultTagMap = ".ftag"
)

var (
	AppVersion = "0.0.0-dev.0"
)

func getTagMapPath(c *cli.Context) (string, error) {
	p := c.GlobalString(optTagMap)
	return resolvePath(p)
}

func resolvePath(p string) (string, error) {
	if !filepath.IsAbs(p) {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}

		p = filepath.Join(cwd, p)
	}
	return p, nil
}

func createFTag(c *cli.Context) (*FTag, error) {

	tagMapPath, err := getTagMapPath(c)
	if err != nil {
		return nil, err
	}

	tagMapStore := tagmap.NewJSONFileStore(tagMapPath)

	ftag := New(tagMapStore)

	err = ftag.LoadTagMap()
	if err != nil {
		return nil, err
	}

	return ftag, nil
}

func getFileArg(c *cli.Context) (string, error) {
	f := c.Args().First()
	if f == "" {
		cli.ShowSubcommandHelp(c)
		return "", errors.New("must supply a file argument")
	}
	return f, nil
}

func getTagArgs(c *cli.Context) ([]string, error) {
	tags := c.Args()[1:]
	if len(tags) < 1 {
		cli.ShowSubcommandHelp(c)
		return nil, errors.New("must supply a tag")
	}
	return tags, nil
}

func commandAdd(c *cli.Context) error {

	f, err := getFileArg(c)
	if err != nil {
		return err
	}

	tags, err := getTagArgs(c)
	if err != nil {
		return err
	}

	ftag, err := createFTag(c)
	if err != nil {
		return err
	}

	err = ftag.Add(f, tags...)
	if err != nil {
		return err
	}

	err = ftag.StoreTagMap()
	if err != nil {
		return err
	}

	return nil
}

func commandCheck(c *cli.Context) error {
	ftag, err := createFTag(c)
	if err != nil {
		return err
	}

	errs := ftag.Check()
	if len(errs) > 0 {
		return cli.NewMultiError(errs...)
	}

	return nil
}

func commandClear(c *cli.Context) error {
	ftag, err := createFTag(c)
	if err != nil {
		return err
	}

	ftag.Clear(c.Args()...)

	err = ftag.StoreTagMap()
	if err != nil {
		return err
	}

	return nil
}

func commandFind(c *cli.Context) error {
	tags := c.Args()
	if len(tags) == 0 {
		cli.ShowSubcommandHelp(c)
		return errors.New("must supply a tag expression")
	}

	ftag, err := createFTag(c)
	if err != nil {
		return err
	}

	files := ftag.Find(tags...)
	for _, f := range files {
		fmt.Println(f)
	}

	return nil
}

func commandList(c *cli.Context) error {

	ftag, err := createFTag(c)
	if err != nil {
		return err
	}

	tags := ftag.List(c.Args())
	for _, t := range tags {
		fmt.Println(t)
	}

	return nil
}

func commandMove(c *cli.Context) error {
	from := c.Args().First()
	if from == "" {
		cli.ShowSubcommandHelp(c)
		return errors.New("must supply a 'from' file argument")
	}

	to := c.Args().Get(1)
	if to == "" {
		cli.ShowSubcommandHelp(c)
		return errors.New("must supply a 'to' file argument")
	}

	ftag, err := createFTag(c)
	if err != nil {
		return err
	}

	err = ftag.Move(from, to)
	if err != nil {
		return err
	}

	err = ftag.StoreTagMap()
	if err != nil {
		return err
	}

	return nil
}

func commandRemove(c *cli.Context) error {

	f, err := getFileArg(c)
	if err != nil {
		return err
	}

	tags, err := getTagArgs(c)
	if err != nil {
		return err
	}

	ftag, err := createFTag(c)
	if err != nil {
		return err
	}

	ftag.Remove(f, tags...)

	err = ftag.StoreTagMap()
	if err != nil {
		return err
	}

	return nil
}

func newCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = AppName
	app.Version = AppVersion
	app.Usage = "\"file tag\" - run "
	app.Author = "Troy Kinsella (troy.kinsella@startmail.com)"

	app.EnableBashCompletion = true

	app.Commands = []cli.Command{
		{
			Name:      "add",
			Aliases:   []string{"a"},
			Usage:     "Add one ore more tags to a file",
			UsageText: AppName + " add <file> <tag> [tag...]",
			Action:    commandAdd,
		},
		{
			Name:   "check",
			Usage:  "Verify that files referenced in the tag mapping exist",
			Action: commandCheck,
		},
		{
			Name:      "clear",
			Aliases:   []string{"clr"},
			Usage:     "Clear all tags associated with the given files",
			UsageText: AppName + " clear <file> [file...]",
			Action:    commandClear,
		},
		{
			Name:      "find",
			Aliases:   []string{"f"},
			Usage:     "Lookup files associated with the given tags",
			UsageText: AppName + " find <tag> [tag...]",
			Action:    commandFind,
		},
		{
			Name:      "list",
			Aliases:   []string{"ls"},
			Usage:     "List tags associated with the given files",
			UsageText: AppName + " list [file...]",
			Action:    commandList,
		},
		{
			Name:      "move",
			Aliases:   []string{"mv"},
			Usage:     "Tell " + AppName + " about a moved file so it can update the tag mapping",
			UsageText: AppName + " move <from> <to>",
			Action:    commandMove,
		},
		{
			Name:      "remove",
			Aliases:   []string{"rm"},
			Usage:     "Remove one or more tags from a file",
			UsageText: AppName + " remove <file> <tag> [tag...]",
			Action:    commandRemove,
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  optTagMap + ", " + optTagMapLong,
			Value: defaultTagMap,
			Usage: "",
		},
	}

	return app
}

func main() {
	app := newCliApp()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
