package kyklos

import (
    "log"
    "io"
    "os"
)

var (
	Debug   *log.Logger
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
    TwoPC   *log.Logger
    Consis  *log.Logger
)

func Init(
	debugHandle io.Writer,
    infoHandle io.Writer,
    warningHandle io.Writer,
    errorHandle io.Writer) {

	Debug = log.New(debugHandle,
        "DEBUG: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Info = log.New(infoHandle,
        "INFO: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Warning = log.New(warningHandle,
        "WARNING: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Error = log.New(errorHandle,
        "ERROR: ",
        log.Ldate|log.Ltime|log.Lshortfile)
}

func InitFileLogs(pcHandle io.Writer,
    csHandle io.Writer) {
    TwoPC = log.New(pcHandle,
        "2PC: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Consis = log.New(csHandle,
        "Consistency: ",
        log.Ldate|log.Ltime|log.Lshortfile)
}

func CheckError(err error) {
    if err != nil {
        Error.Println(err.Error())
        os.Exit(1)
    }
}