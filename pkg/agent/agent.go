package agent

import (
	"errors"
)

type block struct {
	typ  rune
	b    rune
	num  int
	numd float64
}

type Result struct {
	Answer int
	Err    error
}

func delSpace(buff []block, i int, v int) ([]block, int) {
	for j := v; j < i-1; j++ {
		buff[j] = buff[j+1]
	}
	i--
	return buff, i
}

func solveRec(buff []block, i int, s, e int) int {
	if e-s == 1 || i == 1 {
		return buff[s].num
	}

	for j := s; j < e; j++ {
		if buff[j].b == '(' {
			scCnt := 1
			w := j + 1
			for w > e || scCnt > 0 {
				if buff[w].b == '(' {
					scCnt++
				} else if buff[w].b == ')' {
					scCnt--
				}

				w++
			}
			w--

			buff[j].typ = 'd'
			buff[j].num = solveRec(buff, i, j+1, w)
			e -= (w - j)

			buff, i = delSpace(buff, i, j+1)
			buff, i = delSpace(buff, i, j+1)

		}
	}

	for j := s; j < e; j++ {
		if buff[j].typ == 'c' {

			if j+1 < i {

				if buff[j].b == '*' && buff[j-1].typ == 'd' && buff[j+1].typ == 'd' {
					buff[j-1].num *= buff[j+1].num
					buff, i = delSpace(buff, i, j)
					buff, i = delSpace(buff, i, j)
					j = s
					e -= 2
				}
				if buff[j].b == '/' && buff[j-1].typ == 'd' && buff[j+1].typ == 'd' {
					buff[j-1].num /= buff[j+1].num
					buff, i = delSpace(buff, i, j)
					buff, i = delSpace(buff, i, j)
					j = s
					e -= 2
				}
			}
		}
	}
	for j := s; j < e; j++ {
		if buff[j].typ == 'c' {
			if j+1 < i {
				if buff[j].b == '+' && buff[j-1].typ == 'd' && buff[j+1].typ == 'd' {
					buff[j-1].num += buff[j+1].num
					buff, i = delSpace(buff, i, j)
					buff, i = delSpace(buff, i, j)
					j = s
					e -= 2
				}
				if buff[j].b == '-' && buff[j-1].typ == 'd' && buff[j+1].typ == 'd' {
					buff[j-1].num -= buff[j+1].num
					buff, i = delSpace(buff, i, j)
					buff, i = delSpace(buff, i, j)
					j = s
					e -= 2
				}
			}
		}
	}
	return solveRec(buff, i, s, e)

}

func solve(buff []block) Result {

	for j := 0; j < len(buff); j++ {
		if buff[j].b >= '0' && buff[j].b <= '9' {
			buff[j].typ = 'd'
			buff[j].num = int(buff[j].b - '0')
		}

	}
	i := len(buff)

	for j := 0; j < i; j++ {
		if buff[j].typ == 'c' {
			if buff[j].b == ' ' {
				buff, i = delSpace(buff, i, j)
				j--
			}
		}
		if j+1 < i {
			if buff[j].typ == 'd' && buff[j+1].typ == 'd' {
				buff[j].num = buff[j].num*10 + buff[j+1].num
				buff[j+1].typ = 'c'
				buff[j+1].b = ' '
			}
		}
	}
	for i, e := range buff {
		if e.typ == 'c' && (e.b != '(' && e.b != ')' && e.b != ' ') {
			if i < 1 || i >= len(buff)-1 {
				return Result{0, errors.New("error while parsing!")}
			} else {
				if buff[i-1].typ == 'c' || buff[i+1].typ == 'c' {
					return Result{0, errors.New("error while parsing!")}
				}
			}
		}
	}

	return Result{solveRec(buff, i, 0, i), nil}

}

func StartCalculating(s string, res chan<- Result) {
	var buff []block
	var i int
	for i < len(s) {
		if s[i] == '\n' {
			break
		}
		buff = append(buff, block{typ: 'c', b: rune(s[i])})
		i++
	}
	ans := solve(buff)
	res <- ans

}
