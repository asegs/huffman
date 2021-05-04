package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"time"
)

type SortLetter struct {
	val rune
	count int
}

type LetterNode struct {
	letter SortLetter
	left *LetterNode
	right *LetterNode
}

const byteSize = 2048

var letterMap = make(map[rune]int)

func profile(t time.Time,m string){
	d := time.Now().Sub(t)
	fmt.Printf("%v\n%v\n",m,d)
}


func readByteByByteIntoLetterMap(filename string){
	defer profile(time.Now(),"Read to byte array")
	f,err := os.Open(filename)
	if err != nil{
		fmt.Println("File error: "+err.Error())
		return
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	buf := make([]byte,byteSize)
	for {
		_,err = reader.Read(buf)
		if err != nil{
			if err != io.EOF{
				fmt.Println("Error reading file: "+err.Error())
				return
			}
			break
		}
		for _,letter := range buf{
			letterMap[rune(letter)]++
		}
	}
}

func createSortLetterList()[]SortLetter{
	result := make([]SortLetter,len(letterMap))
	count := 0
	for key,value := range letterMap{
		result[count] = SortLetter{
			val:   key,
			count: value,
		}
		count++
	}
	return result
}

func sortLetterList(lst []SortLetter)[]SortLetter{
	sort.Slice(lst, func(i, j int) bool {
		return lst[i].count < lst[j].count
	})
	return lst
}


//fails when new element is greatest in list, inserts in second place
func insertAtIndex(slice []LetterNode,toInsert LetterNode,index int)[]LetterNode{
	newArr := make([]LetterNode,len(slice)+1)
	for i:=0;i<index;i++{
		newArr[i] = slice[i]
	}
	newArr[index] = toInsert
	for i:=index+1;i<len(newArr);i++{
		newArr[i] = slice[i-1]
	}
	return newArr
}


func insertAtHighestIndexLessThan(slice []LetterNode,lowerBound int,upperBound int,toInsert LetterNode,currentGreatest int,greatestIndex int)[]LetterNode{
	if upperBound-lowerBound <= 1 || slice[lowerBound].letter.count > toInsert.letter.count{
		return insertAtIndex(slice,toInsert,greatestIndex+1)
	}
	mid := (upperBound-lowerBound)/2 + lowerBound
	val := slice[mid].letter.count
	if val<=toInsert.letter.count{
		if val > currentGreatest{
			currentGreatest = val
			greatestIndex = mid
		}
		return insertAtHighestIndexLessThan(slice,mid,upperBound,toInsert,currentGreatest,greatestIndex)
	} else{
		return insertAtHighestIndexLessThan(slice,lowerBound,mid,toInsert,currentGreatest,greatestIndex)
	}
}






func main(){
	readByteByByteIntoLetterMap("files/sample.txt")
	result := createSortLetterList()
	sorted := sortLetterList(result)
	newStructs := make([]LetterNode,len(sorted))
	for i:=0;i<len(newStructs);i++{
		newStructs[i] = LetterNode{
			letter: sorted[i],
			left:   nil,
			right:  nil,
		}
	}
	toInsert := LetterNode{
		letter: SortLetter{
			val:   '~',
			count: 70000,
		},
		left:   nil,
		right:  nil,
	}
	insertAt := insertAtHighestIndexLessThan(newStructs,0,len(newStructs)-1,toInsert,-1,-1)
	for i,s := range insertAt{
		fmt.Printf("Index: %d  Char: %c  Count: %d\n",i,s.letter.val,s.letter.count)
	}
}