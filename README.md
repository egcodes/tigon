# Tigon
Tigon is a parser tool. Simply if files compressed, uncompress the files, parsing them according to the configs given, and load them to the database. Three operations run concurency. The concurency settings are given in config.toml

### Tigon currently supports the following stuffs
- [x] Uncompress files | zip, tar.gz, tgz, tar, tar.bz2, tar.xz, rar, 7z, gz
- [x] Transform from csv, txt, xls, xlsx
- [x] Load to Oracle (Sqlldr)

### config.toml default settings
    [path]
    raw = "workspace/files_raw"
    parsed = "workspace/files_parsed"
    backup = "workspace/files_backup"
    config = "config"
    
    [customFileExtention]
    parsedFileExt = ".parsed"
    oracleControlFileExt = ".ctl"
    
    [concurency]
    uncompress = 8
    transform = 16
    load = 8

### Build & Usage
Before you start, run following commands. These commands create a necessary folders and your folder name in workspace according to config.toml and your given argument.
> FolderName is a seperator for your raw files.
```sh
$ dep ensure
$ go build
$ ./tigon <FolderName>
```

Then copy your files to workspace/files_raw/<folderName>. You need to create a config toml file like below in your config directory (workspace/files_config/<folderName>.toml) for the settings of the files in the folder you are giving it. Then start tigon following commands. 

#### Sample config file for txt files
> Below ".test" keyword be the same as your in folderName file. Each file must have an transform, load setting in its name. The same setting file is used for the files in the same directory.
###
    [transform]
        [transform.test]
            parseColumns = [0,3,4]
            parseDataStartIndex = 1
            parseDataEndIndex = -1
            fileSplitChar = ""
            fileRegexStr = "[^?!(\\s)]+"
            outputSplitChar = "|"
    
    [load]
        [load.test]
            db = "oracle"
            loadControlFile = """
                                LOAD DATA
                                INFILE 'workspace/files_parsed/sample_txt/test.parsed'
                                BADFILE 'workspace/files_parsed/sample_txt/test.bad'
                                DISCARDFILE 'workspace/files_parsed/sample_txt/test.dsc'
                                APPEND INTO TABLE NORTHI_PARSER_SETTINGS.TEST
                                Fields terminated by "|" Optionally enclosed by '"'
                                (
                                DATA_DATE DATE "YYYY-MM-DD HH24:MI" NULLIF (DATA_DATE="NULL"),
                                COLUMN1,
                                COLUMN2
                                )
                            """

```sh
$ ./tigon <FolderName>
```

### To Do
- [ ] Scheduler running
- [ ] Add Mysql, Postgresql, Elasticsearch loader
- [ ] XLS transformer is not always working properly

License
----
MIT
