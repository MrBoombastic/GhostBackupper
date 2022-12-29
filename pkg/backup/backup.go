package backup

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mholt/archiver/v4"
	"github.com/urfave/cli/v2"
	"os"

	"github.com/JamesStewy/go-mysqldump"
	_ "github.com/go-sql-driver/mysql"
)

func Create(ctx *cli.Context) error {
	tempDir, err := os.MkdirTemp("", "ghostbackupper-*")
	if err != nil {
		return err
	}
	database, err := dumpDatabase(ctx, tempDir)
	if err != nil {
		return err
	}
	files, err := archiver.FilesFromDisk(nil, map[string]string{
		fmt.Sprintf("%s", ctx.String("content")): "GhostContent",
		fmt.Sprintf("%s/%s", tempDir, database):  "MySQLDatabase",
	})
	out, err := os.Create(ctx.String("output"))
	if err != nil {
		return err
	}
	defer out.Close()
	format := archiver.CompressedArchive{
		Archival:    archiver.Tar{},
		Compression: archiver.Gz{},
	}

	err = format.Archive(context.Background(), out, files)
	if err != nil {
		return err
	}
	return nil
}

// username string, password string, hostname string, port uint, dbname string, dumpDir string
func dumpDatabase(ctx *cli.Context, dumpDir string) (dumpFile string, err error) {
	dumpFilenameFormat := fmt.Sprintf("%s-20060102T150405", ctx.String("database"))
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", ctx.String("username"), ctx.String("password"), ctx.String("hostname"), ctx.Uint("port"), ctx.String("database")))
	if err != nil {
		return "", err
	}
	dumper, err := mysqldump.Register(db, dumpDir, dumpFilenameFormat)
	if err != nil {
		return "", err
	}
	// Dump database to file
	resultFilename, err := dumper.Dump()
	if err != nil {
		return "", err
	}
	// Close dumper and connected database
	dumper.Close()
	return resultFilename, nil
}
