package configuration

import (
	"os"
	"strconv"
	"strings"

	ini "gopkg.in/ini.v1"
)

const (
	configFilename = "config.ini"
)

var configgg *ConfigInfo

//ConfigInfo configuration file information
type ConfigInfo struct {
	DbAddress   string //database address; e.g. localhost
	DbName      string //database name
	DbUsername  string //database username
	DbPassword  string //database password
	DbPort      int    //database port number
	DbInitTable bool   //flag; create basic datatable if not found

	// EmailServer     string //SMTP email server address
	// EmailPortNumber int    //SMTP email server port number
	// EMailUsername   string //SMTP email username
	// EmailPassword   string //SMTP email password

	PortNumber  int    //web server listen port number
	APIStartURL string //ApiStartURL API starting URL path
	LogicDir    string //directory where store logical physical files; e.g. pay-slip.pdf
	StaticDir   string //directory where store direct access physical files; e.g. index.html

	AdminUsername string //AdminUsername system admin username
	AdminPassword string //AdminPassword system admin password

	DevStartURL string //DevStartURL development starting URL
	DevEnable   bool   //DevEnable enable development mode
}

//GetConfig get configutation setting
func GetConfig() *ConfigInfo { return configgg }

//LoadINIConfigFile init .ini file
func LoadINIConfigFile() error {
	//check INI file exists or not; otherwise create one
	if !isFileExists(configFilename) {
		cfg := ini.Empty()
		sec, err := cfg.NewSection("database")
		if err != nil {
			return err
		}
		if _, err = sec.NewKey("dbserver", "localhost"); err != nil {
			return err
		}
		if _, err = sec.NewKey("dbname", ""); err != nil {
			return err
		}
		if _, err = sec.NewKey("dbusername", "root"); err != nil {
			return err
		}
		if _, err = sec.NewKey("dbpassword", ""); err != nil {
			return err
		}
		if _, err = sec.NewKey("dbport", "3306"); err != nil {
			return err
		}
		if _, err := sec.NewKey("db_init_table", "false"); err != nil {
			return err
		}

		sec, err = cfg.NewSection("http")
		if _, err = sec.NewKey("portnumber", "8888"); err != nil {
			return err
		}
		// if _, err = sec.NewKey("api_start_url", "api"); err != nil {
		// 	return err
		// }
		if _, err = sec.NewKey("static_dir", "static-files"); err != nil {
			return err
		}
		if _, err = sec.NewKey("logical_dir", "logical-files"); err != nil {
			return err
		}

		sec, err = cfg.NewSection("admin")
		if _, err = sec.NewKey("username", "admin"); err != nil {
			return err
		}
		if _, err = sec.NewKey("password", "123456789"); err != nil {
			return err
		}

		sec, err = cfg.NewSection("development")
		if _, err = sec.NewKey("enable", "true"); err != nil {
			return err
		}
		if _, err = sec.NewKey("start_url", "dev"); err != nil {
			return err
		}

		//save to physical INI file
		if err = cfg.SaveTo(configFilename); err != nil {
			return err
		}
	}

	tmpConfig, tmpErr := loadConfiguration()
	if tmpErr == nil {
		configgg = tmpConfig
	} else {
		configgg = nil
	}

	return tmpErr
}

//loadConfiguration load .ini file
func loadConfiguration() (*ConfigInfo, error) {
	cfg, err := ini.InsensitiveLoad(configFilename) //ignore capital letter key, all keys is small letter

	//save configuration to physical INI file before exit
	defer cfg.SaveTo(configFilename)

	if err != nil {
		return nil, err
	}

	config := ConfigInfo{}

	dbSection, err := cfg.GetSection("database")
	if err != nil {
		return nil, err
	}
	if config.DbAddress, err = getConfigString(dbSection, "dbserver", "localhost"); err != nil {
		return nil, err
	}
	if config.DbName, err = getConfigString(dbSection, "dbname", ""); err != nil {
		return nil, err
	}
	if config.DbUsername, err = getConfigString(dbSection, "dbusername", "root"); err != nil {
		return nil, err
	}
	if config.DbPassword, err = getConfigString(dbSection, "dbpassword", ""); err != nil {
		return nil, err
	}
	if config.DbPort, err = getConfigInt(dbSection, "dbport", 3306); err != nil {
		return nil, err
	}
	tmp, err := getConfigString(dbSection, "db_init_table", "false")
	if err != nil {
		return nil, err
	}
	if strings.Compare(strings.ToLower(tmp), "true") == 0 {
		config.DbInitTable = true
	} else {
		config.DbInitTable = false
	}

	httpSection, err := cfg.GetSection("http")
	if err != nil {
		return nil, err
	}
	if config.PortNumber, err = getConfigInt(httpSection, "portnumber", 80); err != nil {
		return nil, err
	}
	// if config.APIStartURL, err = getConfigString(httpSection, "api_start_url", "api"); err != nil {
	// 	return nil, err
	// }
	if config.LogicDir, err = getConfigString(httpSection, "logical_dir", "logical-files"); err != nil {
		return nil, err
	}
	if config.StaticDir, err = getConfigString(httpSection, "static_dir", "static-files"); err != nil {
		return nil, err
	}

	adminSection, adminErr := cfg.GetSection("admin")
	if adminErr != nil {
		return nil, adminErr
	}
	if config.AdminUsername, err = getConfigString(adminSection, "username", "admin"); err != nil {
		return nil, err
	}
	if config.AdminPassword, err = getConfigString(adminSection, "password", "123456789"); err != nil {
		return nil, err
	}

	devSection, devErr := cfg.GetSection("development")
	if devErr != nil {
		return nil, devErr
	}
	devModeRaw, ModeErr := getConfigString(devSection, "enable", "true")
	if ModeErr != nil {
		return nil, err
	}
	if config.DevStartURL, err = getConfigString(devSection, "start_url", "dev"); err != nil {
		return nil, err
	}
	config.DevEnable = strings.Compare(strings.ToLower(devModeRaw), "true") == 0

	return &config, nil
}

func getConfigString(section *ini.Section, key string, defaultValue string) (string, error) {
	if section.Haskey(key) {
		return section.Key(key).String(), nil
	}

	section.NewKey(key, defaultValue)
	return defaultValue, nil
}

func getConfigInt(section *ini.Section, key string, defaultValue int) (int, error) {
	if section.Haskey(key) {
		return section.Key(key).Int()
	}

	section.NewKey(key, strconv.Itoa(defaultValue))
	return defaultValue, nil
}

func isFileExists(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false //file not found
		}

		return false //stat command error
	}

	return true //file exists
}

func isDirectoryExists(directoryName string) (bool, error) {
	stat, err := os.Stat(directoryName)

	if err != nil {
		return false, nil //other errors
	}

	return stat.IsDir(), nil
}
