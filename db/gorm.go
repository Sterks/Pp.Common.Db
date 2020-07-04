package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Sterks/Pp.Common.Db/models"

	"github.com/Sterks/fReader/config"
	"github.com/Sterks/fReader/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //....
)

//Database ...
type Database struct {
	Config   *config.Config
	Database *gorm.DB
	Logger   *logger.Logger
}

const (
	host     = "localhost"
	port     = 5432
	user     = "user_ro"
	password = "4r2w3e1q"
	dbname   = "freader"
)

// OpenDatabase ...
func (d *Database) OpenDatabase() *gorm.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		log.Printf("Соединиться не удалось - %s", err)
	}
	if err2 := db.DB().Ping(); err2 != nil {
		log.Printf("База не отвечает - %v", err2)
	}
	d.Database = db
	return d.Database
}

// CreateInfoFile ...
func (d *Database) CreateInfoFile(info os.FileInfo, region string, hash string, fullpath string, typeFile string, file string) int {
	// d.database.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&files)
	// filesTypes := d.database.Table("FileType")

	d.Database.LogMode(true)

	var gf models.SourceRegions
	d.Database.Table("SourceRegions").Where("r_name = ?", region).Find(&gf)

	var sr models.SourceResources
	d.Database.Table("SourceResources").Where("sr_name = ?", file).Find(&sr)

	checker := d.CheckExistFileDb(info, hash)
	if checker != 0 {
		var lf models.File
		d.Database.Table("Files").Where("f_id = ?", checker).Find(&lf)
		lf.TDateLastCheck = time.Now()
		d.Database.Save(&lf)
		log.Printf("Дата успешно обновлена - %v", lf.TDateLastCheck.String())
	}
	if checker == 0 {

		ext := filepath.Ext(info.Name())

		var fileType models.FileType
		var lastID models.File
		d.Database.Table("FilesTypes").Where("ft_ext = ?", ext).Find(&fileType)

		d.Database.Table("Files")
		d.Database.Create(&models.File{
			TName:                 info.Name(),
			TArea:                 gf.RID,
			FileType:              fileType,
			TType:                 fileType.FTID,
			THash:                 hash,
			TSize:                 info.Size(),
			CreatedAt:             time.Now(),
			TDateCreateFromSource: info.ModTime(),
			TDateLastCheck:        time.Now(),
			TFullpath:             fullpath,
			TSourceResources:      sr.SRID,
		}).Scan(&lastID)
		log.Printf("Файл успешно добавлен - %v", lastID.TName)
		return lastID.TID
	} else {
		log.Printf("Файл существует - %v\n", info.Name())
		return 0
	}
}

// LastID ...
func (d *Database) LastID() int {
	var ff models.File
	d.Database.Table("Files").Last(&ff)
	return ff.TID
}

// QuantityTypeDoc ...
func (d *Database) QuantityTypeDoc(typeFile string) int {
	var sr models.SourceResources
	d.Database.Table("SourceResources").Where("sr_name = ?", typeFile).Find(&sr)

	var ff []models.File
	var total int
	d.Database.Model(&models.File{}).Related(&models.FileType{})
	d.Database.Table("Files").Where("f_source_resources_id = ?", sr.SRID).Find(&ff).Count(&total)
	return total
}

// CheckerExistFileDBNotHash ...
func (d *Database) CheckerExistFileDBNotHash(file os.FileInfo) (int, string) {
	var ff models.File
	fmt.Printf("%v - %v - %v", file.Size(), file.Name(), file.ModTime())
	d.Database.Table("Files").Where("f_size = ? and f_name = ? and f_date_create_from_source = ?", file.Size(), file.Name(), file.ModTime()).Find(&ff)
	return ff.TID, ff.THash
}

// CheckExistFileDb ...
func (d *Database) CheckExistFileDb(file os.FileInfo, hash string) int {
	var ff models.File
	d.Database.Table("Files").Where("f_hash = ? and f_size = ? and f_name = ?", hash, file.Size(), file.Name()).Find(&ff)
	return ff.TID
}

//CheckRegionsDb Проверка существует ли регион в базе данных
func (d *Database) CheckRegionsDb(region string) int {
	var reg models.SourceRegions
	d.Database.Table("SourceRegions").Where("r_name = ?", region).First(&reg)
	return reg.RID
}

// CheckSourceResourcesDb - Вернуть ID ресурса
func (d *Database) CheckSourceResourcesDb(resource string) int {
	var res models.SourceResources
	d.Database.Table("SourceResources").Where("sr_name = ?", resource).First(&res)
	return res.SRID
}

//ReaderRegionsDb Все регионы из базы
func (d *Database) ReaderRegionsDb() []models.SourceRegions {
	var regions []models.SourceRegions
	d.Database.Table("SourceRegions").Find(&regions)
	return regions
}

//AddRegionsDb ...
func (d *Database) AddRegionsDb(region string, law string) {
	var typeLaw models.FederalLaw
	d.Database.Table("FederalLaw").Where("fl_name_law = ?", law).First(&typeLaw)
	var reg models.SourceRegions
	reg.RName = region
	reg.RDateCreate = time.Now()
	reg.RDateUpdate = time.Now()
	reg.RFZLaw = typeLaw.RLID
	d.Database.Table("SourceRegions").Create(&reg)
}

// FirstOrCreate Создать или получить
func (d *Database) FirstOrCreate(region string, law string) models.SourceRegions {
	var typeLaw models.FederalLaw
	d.Database.Table("FederalLaw").Where("fl_name_law = ?", law).First(&typeLaw)
	var reg models.SourceRegions
	reg.RName = region
	reg.RDateCreate = time.Now()
	reg.RDateUpdate = time.Now()
	reg.RFZLaw = typeLaw.RLID
	d.Database.Table("SourceRegions").Where("r_name = ?", region).FirstOrCreate(&reg)
	return reg
}

//CreateTask - дабавление заданий
func (d *Database) CreateTask(tsName string, tsDataStart time.Time, tsRunTimes int, tsComment string) {
	d.Database.Table("Tasks").Create(&models.Task{
		TSName:      tsName,
		TSDataStart: tsDataStart,
		TSRunTimes:  tsRunTimes,
		TSComment:   tsComment,
	})
}
