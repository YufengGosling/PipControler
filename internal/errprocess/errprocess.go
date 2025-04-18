package pipctrler_sai

import (
	"log"
)

/*
Author: Yufeng Gosling
Creation time: 2025/4/6
*/

/*
ErrProcess Handle errors in a unified manner within the project.

Tips: If the passed err is nil, then the function will not run.

Param:
errorMessage - Error message format. Should follow the format used by Golang's printf. For example: Cannot do something %v\n. And if there are more parameters,% v must be placed at the end
err - Pass in a value of type error. This will be included in the output error message.
processMethod - Handling method. Have "logging"(Output error but do not terminate operation), "error"(Output error and terminate operation), and "panic"(Direct panic). If the processing method does not exist. It will output error messages but will not terminate the operation.In this case, the function will not output normally.
morePar - More error message, such as reading failed files. This parameter must be of type interface. If this parameter is present, the function will use append to merge it with theErr into one slice at runtime.

Return:
No return.

Example:
errProcess("Cannot read file %s\n Error: %v", err, "logging", fileInfo)

This annotation is translated using Baidu Translate and may be inaccurate. Please forgive us.
*/
func ErrProcess(errorMessage string, err error, processMethod string, morePar ...interface{}) {
	if err != nil {
    	haveMorePar := morePar != nil
    	logPar := append(morePar, theErr)
    	switch processMethod {
    	case "logging":
    		if haveMorePar {
				log.Printf(errorMessage, logPar...)
			} else {
		    	log.Printf(errorMessage, theErr)
			}
		case "error":
			if haveMorePar {
				log.Fatalf(errorMessage, logPar...)
			} else {
				log.Fatalf(errorMessage, theErr)
			}
		case "panic":
			if haveMorePar {
				log.Panicf(errorMessage, logPar...)
			} else {
				log.Panicf(errorMessage, theErr)
			}
		default:
			log.Printf("The process %s method not found\n", processMethod)
		}
	}
}

