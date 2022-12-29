package backup

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/JamesStewy/go-mysqldump"
	"github.com/MrBoombastic/GhostBackupper/pkg/logs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mholt/archiver/v4"
	"github.com/t3rm1n4l/go-mega"
	"github.com/urfave/cli/v2"
	"os"
)

func Create(ctx *cli.Context) error {
	logs.Info("Dumping database")
	database, err := dumpDatabase(ctx)
	if err != nil {
		return err
	}
	logs.Info("Removing logs")
	err = os.RemoveAll(fmt.Sprintf("%v/logs/", ctx.String("content")))
	if err != nil {
		return err
	}
	logs.Info("Archiving files")
	logs.Info(database)
	files, err := archiver.FilesFromDisk(nil, map[string]string{
		ctx.String("content"): "GhostContent",
		database:              "MySQLDatabase",
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
	logs.Info("Success!")
	err = Upload(ctx)
	if err != nil {
		return err
	}
	return nil
}

func Upload(ctx *cli.Context) error {
	logs.Info("Logging into Mega")
	client := mega.New()
	err := client.Login(ctx.String("mega_login"), ctx.String("mega_password"))
	if err != nil {
		return err
	}
	logs.Info("Uploading to Mega")
	_, err = client.UploadFile(ctx.String("output"), client.FS.GetRoot(), "", nil)
	if err != nil {
		return err
	}
	logs.Info("Upload done")
	return nil
}

// username string, password string, hostname string, port uint, dbname string, dumpDir string
func dumpDatabase(ctx *cli.Context) (dumpFile string, err error) {
	dumpFilename := fmt.Sprintf("%s-20060102T150405", ctx.String("database"))
	logs.Info("Connecting to MySQL")
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", ctx.String("username"), ctx.String("password"), ctx.String("host"), ctx.Uint("port"), ctx.String("database")))
	if err != nil {
		return "", err
	}
	dumper, err := mysqldump.Register(db, "./", dumpFilename)
	if err != nil {
		return "", err
	}
	logs.Info(fmt.Sprintf("Dumping database to %v", dumpFilename))
	resultFilename, err := dumper.Dump()
	if err != nil {
		return "", err
	}
	// Close dumper and connected database
	dumper.Close()
	logs.Info("Database dumping done")
	return resultFilename, nil
}
