package kiln

import (
	"errors"

	"github.com/pivotal-cf/jhanda/flags"
)

type ArgParser struct{}

func NewArgParser() ArgParser {
	return ArgParser{}
}

func (a ArgParser) Parse(args []string) (ApplicationConfig, error) {
	cfg := ApplicationConfig{}

	args, err := flags.Parse(&cfg, args)
	if err != nil {
		panic(err)
	}

	if len(cfg.ReleaseTarballs) == 0 {
		return cfg, errors.New("Please specify at least one release tarball with the --release-tarball parameter")
	}

	if cfg.StemcellTarball == "" {
		return cfg, errors.New("--stemcell-tarball is a required parameter")
	}

	if cfg.Handcraft == "" {
		return cfg, errors.New("--handcraft is a required parameter")
	}

	if cfg.Version == "" {
		return cfg, errors.New("--version is a required parameter")
	}

	if cfg.FinalVersion == "" {
		return cfg, errors.New("--final-version is a required parameter")
	}

	if cfg.ProductName == "" {
		return cfg, errors.New("--product-name is a required parameter")
	}

	if cfg.FilenamePrefix == "" {
		return cfg, errors.New("--filename-prefix is a required parameter")
	}

	if cfg.OutputDir == "" {
		return cfg, errors.New("--output-dir is a required parameter")
	}

	if len(cfg.Migrations) > 0 && len(cfg.ContentMigrations) > 0 {
		return cfg, errors.New("cannot build a tile with content migrations and migrations")
	}

	if len(cfg.ContentMigrations) > 0 && cfg.BaseContentMigration == "" {
		return cfg, errors.New("base content migration is required when content migrations are provided")
	}

	if len(cfg.Migrations) > 0 && cfg.BaseContentMigration != "" {
		return cfg, errors.New("cannot build a tile with a base content migration and migrations")
	}

	return cfg, nil
}
