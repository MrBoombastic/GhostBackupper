package backup

import (
	"context"
	"fmt"
	"github.com/MrBoombastic/GhostBackupper/pkg/logs"
	"github.com/jarvanstack/mysqldump"
	"github.com/mholt/archiver/v4"
	"github.com/t3rm1n4l/go-mega"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

func Create(ctx *cli.Context) error {
	// check if directory exists
	if _, err := os.Stat(ctx.String("content")); os.IsNotExist(err) {
		logs.Error("Content directory does not exist")
		return nil
	}

	logs.Info("Dumping database")
	database, err := dumpDatabase(ctx)
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	logs.Info("Removing logs")
	err = os.RemoveAll(fmt.Sprintf("%v/logs/", ctx.String("content")))
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	logs.Info("Archiving files")

	files, err := archiver.FilesFromDisk(nil, map[string]string{
		ctx.String("content"): "GhostContent",
		database:              database,
	})

	filename := fmt.Sprintf("%v-%v", time.Now().UnixMilli(), ctx.String("output"))
	out, err := os.Create(filename)
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
		logs.Error(err.Error())
		return nil
	}

	// remove database dump
	err = os.Remove(database)
	if err != nil {
		logs.Error("failed to delete sql dump temp file:" + err.Error())
		return nil
	}

	logs.Info("Backed up successfully!")
	if ctx.String("mega_login") != "" && ctx.String("mega_password") != "" {
		err = uploadToMega(ctx, filename)
		if err != nil {
			logs.Error(err.Error())
			return nil
		}
	}
	return nil
}

func uploadToMega(ctx *cli.Context, filename string) error {
	logs.Info("Logging into Mega")
	client := mega.New()
	err := client.Login(ctx.String("mega_login"), ctx.String("mega_password"))
	if err != nil {
		return err
	}
	logs.Info("Uploading to Mega")
	_, err = client.UploadFile(filename, client.FS.GetRoot(), "", nil)
	if err != nil {
		return err
	}
	logs.Info("Uploading done successfully!")
	return nil
}

// username string, password string, hostname string, port uint, dbname string, dumpDir string
func dumpDatabase(ctx *cli.Context) (filename string, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=true&loc=UTC", ctx.String("db_user"), ctx.String("db_password"), ctx.String("db_host"), ctx.Uint("db_port"), ctx.String("db_database"))
	filename = fmt.Sprintf("%v.sql", ctx.String("db_database"))
	f, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("couldn't create dump file: %v", err)
	}
	logs.Info("Exporting MySQL...")

	err = mysqldump.Dump(
		dsn,
		mysqldump.WithData(),
		mysqldump.WithWriter(f),
	)

	if err != nil {
		return "", err
	}

	logs.Info("Database dumping done")
	return filename, nil
}
