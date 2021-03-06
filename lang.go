//line lang.y:2
package gohaml

import "fmt"

var Output inode

//line lang.y:9
type yySymType struct {
	yys int
	n   inode
	s   string
	i   interface{}
	c   icodenode
}

const IDENT = 57346
const ATOM = 57347
const FOR = 57348
const RANGE = 57349

var yyToknames = []string{
	"IDENT",
	"ATOM",
	"FOR",
	"RANGE",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line lang.y:66

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 7
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 19

var yyAct = []int{

	13, 14, 15, 7, 12, 5, 6, 17, 3, 19,
	2, 11, 10, 16, 8, 4, 9, 18, 1,
}
var yyPact = []int{

	4, -1000, 11, -4, -2, -7, 10, 7, -5, -1000,
	-1000, -10, -8, -1000, 9, 0, -10, 5, -1000, -1000,
}
var yyPgo = []int{

	0, 18, 16, 0,
}
var yyR1 = []int{

	0, 1, 1, 2, 2, 3, 3,
}
var yyR2 = []int{

	0, 8, 4, 1, 2, 3, 0,
}
var yyChk = []int{

	-1000, -1, 6, 4, 4, 9, 8, 10, 4, -2,
	5, 4, 9, -3, 11, 10, 4, 7, -3, 4,
}
var yyDef = []int{

	0, -2, 0, 0, 0, 0, 0, 0, 0, 2,
	3, 6, 0, 4, 0, 0, 6, 0, 5, 1,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 8, 3, 11, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 9, 3,
	3, 10,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7,
}
var yyTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c > 0 && c <= len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return fmt.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return fmt.Sprintf("state-%v", s)
}

func yylex1(lex yyLexer, lval *yySymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		c = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			c = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		c = yyTok3[i+0]
		if c == char {
			c = yyTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		fmt.Printf("lex %U %s\n", uint(char), yyTokname(c))
	}
	return c
}

func yyParse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		fmt.Printf("char %v in %v\n", yyTokname(yychar), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yychar < 0 {
		yychar = yylex1(yylex, &yylval)
	}
	yyn += yychar
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yychar { /* valid shift */
		yychar = -1
		yyVAL = yylval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yychar < 0 {
			yychar = yylex1(yylex, &yylval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yychar {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error("syntax error")
			Nerrs++
			if yyDebug >= 1 {
				fmt.Printf("%s", yyStatname(yystate))
				fmt.Printf("saw %s\n", yyTokname(yychar))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					fmt.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				fmt.Printf("error recovery discards %s\n", yyTokname(yychar))
			}
			if yychar == yyEofCode {
				goto ret1
			}
			yychar = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		fmt.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		//line lang.y:25
		{
			rn := new(rangenode)
			rn._lhs1 = yyS[yypt-6].s
			rn._lhs2 = yyS[yypt-4].s
			rn._rhs = res{yyS[yypt-0].s, true}
			yyVAL.n = rn
			Output = yyVAL.n
		}
	case 2:
		//line lang.y:34
		{
			yyS[yypt-0].c.setLHS(yyS[yypt-3].s)
			yyVAL.n = yyS[yypt-0].c
			Output = yyVAL.n
		}
	case 3:
		//line lang.y:42
		{
			dan := new(declassnode)
			dan._rhs = yyS[yypt-0].i
			yyVAL.c = dan
		}
	case 4:
		//line lang.y:48
		{
			dan := new(vdeclassnode)
			dan._rhs.value = yyS[yypt-1].s + yyS[yypt-0].s
			dan._rhs.needsResolution = true
			yyVAL.c = dan
		}
	case 5:
		//line lang.y:57
		{
			yyVAL.s = fmt.Sprintf(".%s%s", yyS[yypt-1].s, yyS[yypt-0].s)
		}
	case 6:
		//line lang.y:61
		{
			yyVAL.s = ""
		}
	}
	goto yystack /* stack new state and value */
}
