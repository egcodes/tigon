package main

// Path configs
type Path struct {
	Raw    string
	Parsed string
	Backup string
	Config string
}

// CustomFileExtention use with tigon
type CustomFileExtention struct {
	ParsedFileExt        string
	OracleControlFileExt string
}

// Concurency limits
type Concurency struct {
	Uncompress int
	Transform  int
	Load       int
}

// Scheduler running time
type Scheduler struct {
	Every uint64
}

// Configuration all
type Configuration struct {
	Path                Path
	CustomFileExtention CustomFileExtention
	Concurency          Concurency
	Scheduler           Scheduler
}
