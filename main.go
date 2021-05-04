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
	if toInsert.letter.count>slice[len(slice)-1].letter.count{
		return append(slice,toInsert)
	}
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


func insertAtHighestIndexLessThan(slice []LetterNode,lowerBound int,upperBound int,toInsert LetterNode,currentGreatest int,greatestIndex int)int{
	if upperBound-lowerBound <= 1 || slice[lowerBound].letter.count > toInsert.letter.count{
		return greatestIndex+1
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


func toTree(sorted []LetterNode)LetterNode{
	if len(sorted)==1{
		return sorted[0]
	}
	lowest := sorted[0]
	lower := sorted[1]
	newNode := LetterNode{
		letter: SortLetter{
			val:   0,
			count: lowest.letter.count+lower.letter.count,
		},
		left:   &lowest,
		right:  &lower,
	}
	location := insertAtHighestIndexLessThan(sorted,0,len(sorted)-1,newNode,-1,-1)
	sorted = insertAtIndex(sorted,newNode,location)
	return toTree(sorted[2:])
}

func letterToNode(s []SortLetter)[]LetterNode{
	nodeSorted := make([]LetterNode,len(s))
	for i,node := range s{
		nodeSorted[i] = LetterNode{
			letter: node,
			left:   nil,
			right:  nil,
		}
	}
	return nodeSorted
}


func main(){
	readByteByByteIntoLetterMap("files/sample.txt")
	result := createSortLetterList()
	sorted := sortLetterList(result)
	nodeSorted := letterToNode(sorted)
	treeNode := toTree(nodeSorted)
}