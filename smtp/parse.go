package smtp

import (
	"fmt"
	"io"
	"os"
	"time"
)

var (
	filename = "/var/log/mail.log"
	//filename    = "mail.log"
	numLines    = 25
	filewritten = opt.Receiver + "_log.txt"
)

// writes the log to a new file
func writeLog() {
	f, err := os.Create(filewritten)
	if err != nil {
		fmt.Println(err)
		return
	}
	text := getLog()
	l, err := f.WriteString(text)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func getLog() string {
	fmt.Printf("Before getting the log parsed, wait for 100 seconds")
	time.Sleep(10 * time.Second)
	text, err := GoTail(filename, numLines)
	if err != nil {
		fmt.Println(err)
	}
	return text
}

// GoTail function reads from the last line of the file that is parsed.
// Thanks to: https://stackoverflow.com/questions/52247755/remove-first-n-lines-of-file
func GoTail(filename string, numLines int) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	//SEEK BACKWARD CHARACTER BY CHARACTER ADDING UP NEW LINES
	//offset must start at "-1" otherwise we are already at the EOF
	//"-1" from numLines since we ignore "last" newline in a file
	numNewLines := 0
	var offset int64 = -1
	var finalReadStartPos int64
	for numNewLines <= numLines-1 {
		//seek to new position in file
		startPos, err := file.Seek(offset, 2)
		if err != nil {
			return "", err
		}

		//make sure start position can never be less than 0
		//aka, you cannot read from before the file starts
		if startPos == 0 {
			//set to -1 since we +1 to this below
			//the position will then start from the first character
			finalReadStartPos = -1
			break
		}

		//read the character at this position
		b := make([]byte, 1)
		_, err = file.ReadAt(b, startPos)
		if err != nil {
			return "", err
		}

		//ignore if first character being read is a newline
		if offset == int64(-1) && string(b) == "\n" {
			offset--
			continue
		}

		//if the character is a newline
		//add this to the number of lines read
		//and remember position in case we have reached our target number of lines
		if string(b) == "\n" {
			numNewLines++
			finalReadStartPos = startPos
		}

		//decrease offset for reading next character
		//remember, we are reading backward!
		offset--
	}

	//READ TO END OF FILE
	//add "1" here to move offset from the newline position to first character in line of text
	//this position should be the first character in the "first" line of data we want
	endPos, err := file.Seek(int64(-1), 2)
	if err != nil {
		return "", err
	}
	b := make([]byte, (endPos+1)-finalReadStartPos)
	_, err = file.ReadAt(b, finalReadStartPos+1)
	if err == io.EOF {
		return string(b), nil
	} else if err != nil {
		return "", err
	}

	//special case
	//if text is read, then err == io.EOF should hit
	//there should *never* not be an error above
	//so this line should never return
	return "**No error but no text read.**", nil
}