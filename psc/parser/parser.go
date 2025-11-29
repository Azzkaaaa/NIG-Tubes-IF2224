package parser

import (
	"fmt"
	"strings"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

type Parser struct {
	buffer []dt.Token
	pos    int
}

type ParseError struct {
	buffer   []dt.Token
	Line     int
	Col      int
	Tips     string
	Got      *dt.Token
	Expected []dt.TokenType
}

func reconstruct_source(tokens []dt.Token, line int) string {
	var sb strings.Builder
	currentLine := max(1, line-1)
	currenCol := 1
	foundLine := false
	for _, t := range tokens {
		if t.Line < max(1, line-1) {
			continue
		}
		if t.Line > line {
			break
		}
		if !foundLine && t.Line == currentLine {
			sb.WriteString(fmt.Sprintf("\n%d: ", t.Line))
			foundLine = true
		} else if t.Line > currentLine {
			sb.WriteString(fmt.Sprintf("\n%d: ", t.Line))
			currenCol = 1
			currentLine = t.Line
		}
		sb.WriteString(strings.Repeat(" ", max(t.Col-currenCol, 0)))
		sb.WriteString(t.Lexeme)
		currenCol += max(t.Col-currenCol, 0) + len(t.Lexeme)

	}
	return sb.String()
}

func (e *ParseError) Error() string {
	expected := ""
	if len(e.Expected) > 0 {
		for i, t := range e.Expected {
			if i > 0 {
				expected += ", "
			}
			expected += t.String()
		}
		expected = fmt.Sprintf("(expected: %s)", expected)
	}

	if e.Tips != "" {
		return fmt.Sprintf(`
Line %d, Col %d
%s
%s
SyntaxError: Unexpected token '%s' %s
%s
`,
			e.Line, e.Col, reconstruct_source(e.buffer, e.Line), strings.Repeat(" ", e.Col-1+3)+"^", e.Got.Lexeme, expected, e.Tips)
	} else {
		return fmt.Sprintf(`
Line %d, Col %d
%s
%s
SyntaxError: Unexpected token '%s' %s
`,
			e.Line, e.Col, reconstruct_source(e.buffer, e.Line), strings.Repeat(" ", e.Col-1+3)+"^", e.Got.Lexeme, expected)
	}
}

func (p *Parser) createParseError(expectedType dt.TokenType, tips string) error {
	curr := p.buffer[p.pos]
	return &ParseError{
		buffer:   p.buffer,
		Line:     curr.Line,
		Col:      curr.Col,
		Tips:     tips,
		Expected: []dt.TokenType{expectedType},
		Got:      &curr,
	}
}

func (p *Parser) createParseErrorMany(expectedType []dt.TokenType, tips string) error {
	curr := p.buffer[p.pos]
	return &ParseError{
		buffer:   p.buffer,
		Line:     curr.Line,
		Col:      curr.Col,
		Tips:     tips,
		Expected: expectedType,
		Got:      &curr,
	}
}

func New(tokens []dt.Token) *Parser {
	return &Parser{
		buffer: tokens,
		pos:    0,
	}
}

func (p *Parser) peek() *dt.Token {
	return &p.buffer[p.pos]
}

func (p *Parser) consume(expectedType dt.TokenType) *dt.Token {
	curr := p.peek()
	if curr.Type == expectedType {
		p.pos++
		return curr
	}
	return nil
}

func (p *Parser) consumeMany(expectedType []dt.TokenType) *dt.Token {
	curr := p.peek()
	for _, tokenType := range expectedType {
		if curr.Type == tokenType {
			p.pos++
			return curr
		}
	}
	return nil
}

func (p *Parser) consumeExact(expectedType dt.TokenType, expectedLexeme string) *dt.Token {
	curr := p.peek()
	if curr.Type == expectedType && curr.Lexeme == expectedLexeme {
		p.pos++
		return curr
	}

	return nil
}

func (p *Parser) match(expectedType dt.TokenType) bool {
	curr := p.peek()
	return curr.Type == expectedType
}

func (p *Parser) matchExact(expectedType dt.TokenType, expectedLexeme string) bool {
	curr := p.peek()
	return curr.Type == expectedType && curr.Lexeme == expectedLexeme
}

func (p *Parser) Parse() (*dt.ParseTree, error) {
	tree, err := p.parseProgram()
	return tree, err
}

func (p *Parser) parseProgram() (*dt.ParseTree, error) {
	headerTree, err := p.parseProgramHeader()

	if err != nil {
		return nil, err
	}

	declarationTree, err := p.parseDeclarationPart()

	if err != nil {
		return nil, err
	}

	compoundTree, err := p.parseCompoundStatement()
	if err != nil {
		return nil, err
	}

	dotToken := p.consume(dt.DOT)
	if dotToken == nil {
		return nil, p.createParseError(dt.DOT, "program must end with a dot (.)")
	}

	if p.pos < len(p.buffer) {
		curr := p.peek()
		return nil, &ParseError{
			buffer: p.buffer,
			Line:   curr.Line,
			Col:    curr.Col,
			Tips:   "unexpected token after program end (.)",
			Got:    curr,
		}
	}

	programTree := dt.ParseTree{
		RootType:   dt.PROGRAM_NODE,
		TokenValue: nil,
		Children: []dt.ParseTree{
			*headerTree,
			*declarationTree,
			*compoundTree,
			{
				RootType:   dt.TOKEN_NODE,
				TokenValue: dotToken,
				Children:   make([]dt.ParseTree, 0),
			},
		},
	}

	return &programTree, nil
}

func (p *Parser) parseProgramHeader() (*dt.ParseTree, error) {

	if p.consumeExact(dt.KEYWORD, "program") == nil {
		return nil, p.createParseError(dt.KEYWORD, "all programs must start with program keyword")
	}

	if p.consume(dt.IDENTIFIER) == nil {
		return nil, p.createParseError(dt.IDENTIFIER, "program name must only use alphanumerical characters and underscores")
	}

	if p.consume(dt.SEMICOLON) == nil {
		return nil, p.createParseError(dt.SEMICOLON, "program name must be a single word and strictly end with ;")
	}

	headerTree := dt.ParseTree{
		RootType:   dt.PROGRAM_HEADER_NODE,
		TokenValue: nil,
		Children: []dt.ParseTree{{
			RootType:   dt.TOKEN_NODE,
			TokenValue: &p.buffer[0],
			Children:   make([]dt.ParseTree, 0),
		}, {
			RootType:   dt.TOKEN_NODE,
			TokenValue: &p.buffer[1],
			Children:   make([]dt.ParseTree, 0),
		}, {
			RootType:   dt.TOKEN_NODE,
			TokenValue: &p.buffer[2],
			Children:   make([]dt.ParseTree, 0),
		}},
	}

	return &headerTree, nil
}

func (p *Parser) parseDeclarationPart() (*dt.ParseTree, error) {
	var err error

	declarationTree := dt.ParseTree{
		RootType:   dt.DECLARATION_PART_NODE,
		TokenValue: nil,
		Children:   make([]dt.ParseTree, 0),
	}

	for err == nil {
		constDeclaration, newErr := p.parseConstDeclarationPart()
		if constDeclaration == nil && newErr == nil {
			break
		}
		err = newErr
		if err == nil {
			declarationTree.Children = append(declarationTree.Children, *constDeclaration)
		}
	}

	for err == nil {
		typeDeclaration, newErr := p.parseTypeDeclarationPart()
		if typeDeclaration == nil && newErr == nil {
			break
		}
		err = newErr
		if err == nil {
			declarationTree.Children = append(declarationTree.Children, *typeDeclaration)
		}
	}

	for err == nil {
		varDeclaration, newErr := p.parseVarDeclarationPart()
		if varDeclaration == nil && newErr == nil {
			break
		}
		err = newErr
		if err == nil {
			declarationTree.Children = append(declarationTree.Children, *varDeclaration)
		}
	}

	for err == nil {
		subprogramDeclaration, newErr := p.parseSubprogramDeclaration()

		if subprogramDeclaration == nil && newErr == nil {
			break
		}

		err = newErr
		if err == nil {
			declarationTree.Children = append(declarationTree.Children, *subprogramDeclaration)
		}
	}

	return &declarationTree, err
}

func (p *Parser) parseConstDeclarationPart() (*dt.ParseTree, error) {
	expectedConst := p.consumeExact(dt.KEYWORD, "konstanta")
	if expectedConst == nil {
		return nil, nil
	}

	constDeclaration := dt.ParseTree{
		RootType:   dt.CONST_DECLARATION_PART_NODE,
		TokenValue: nil,
		Children: []dt.ParseTree{
			{
				RootType:   dt.TOKEN_NODE,
				TokenValue: expectedConst,
				Children:   nil,
			},
		},
	}

	if !p.match(dt.IDENTIFIER) {
		return nil, p.createParseError(dt.IDENTIFIER, "expected at least one const definition")
	}

	for p.match(dt.IDENTIFIER) {
		constDefintion, err := p.parseConstDeclaration()

		if err != nil {
			return nil, err
		}

		constDeclaration.Children = append(constDeclaration.Children,
			*constDefintion,
		)
	}
	return &constDeclaration, nil
}

func (p *Parser) parseConstDeclaration() (*dt.ParseTree, error) {

	identifier := p.consume(dt.IDENTIFIER)

	if identifier == nil {
		return nil, p.createParseError(dt.IDENTIFIER, "expected identifier on const definition")
	}

	expectedEqual := p.consumeExact(dt.RELATIONAL_OPERATOR, "=")
	if expectedEqual == nil {
		return nil, p.createParseError(dt.RELATIONAL_OPERATOR, "expected = after const identifier")
	}

	value := p.consumeMany([]dt.TokenType{dt.CHAR_LITERAL, dt.STRING_LITERAL, dt.NUMBER})
	if value == nil {
		return nil, p.createParseErrorMany([]dt.TokenType{dt.CHAR_LITERAL, dt.STRING_LITERAL, dt.NUMBER}, "expected literal value for const definition")
	}

	expectedSemicolon := p.consume(dt.SEMICOLON)
	if expectedSemicolon == nil {
		return nil, p.createParseError(dt.SEMICOLON, "const definition must end with semicolon")
	}

	constDefintion := dt.ParseTree{
		RootType:   dt.CONST_DECLARATION_NODE,
		TokenValue: nil,
		Children: []dt.ParseTree{
			{
				RootType:   dt.TOKEN_NODE,
				TokenValue: identifier,
				Children:   nil,
			},
			{
				RootType:   dt.TOKEN_NODE,
				TokenValue: expectedEqual,
				Children:   nil,
			},
			{
				RootType:   dt.TOKEN_NODE,
				TokenValue: value,
				Children:   nil,
			},
			{
				RootType:   dt.TOKEN_NODE,
				TokenValue: expectedSemicolon,
				Children:   nil,
			},
		},
	}

	return &constDefintion, nil
}

func (p *Parser) parseTypeDeclarationPart() (*dt.ParseTree, error) {
	expectedType := p.consumeExact(dt.KEYWORD, "tipe")
	if expectedType == nil {
		return nil, nil
	}

	typeDeclarationPart := dt.ParseTree{
		RootType:   dt.TYPE_DECLARATION_PART_NODE,
		TokenValue: nil,
		Children: []dt.ParseTree{
			{
				RootType:   dt.TOKEN_NODE,
				TokenValue: expectedType,
				Children:   nil,
			},
		},
	}

	if !p.match(dt.IDENTIFIER) {
		return nil, p.createParseError(dt.IDENTIFIER, "expected at least one type declaration")
	}

	for p.match(dt.IDENTIFIER) {
		typeDeclaration, err := p.parseTypeDeclaration()

		if err != nil {
			return nil, err
		}

		typeDeclarationPart.Children = append(typeDeclarationPart.Children,
			*typeDeclaration,
		)
	}
	return &typeDeclarationPart, nil
}

func (p *Parser) parseTypeDeclaration() (*dt.ParseTree, error) {

	identifier := p.consume(dt.IDENTIFIER)

	if identifier == nil {
		return nil, p.createParseError(dt.IDENTIFIER, "expected identifier on type declaration")
	}

	expectedEqual := p.consumeExact(dt.RELATIONAL_OPERATOR, "=")
	if expectedEqual == nil {
		return nil, p.createParseError(dt.RELATIONAL_OPERATOR, "expected = after type identifier")
	}

	var parsedType *dt.ParseTree
	var err error

	if p.matchExact(dt.KEYWORD, "rekaman") {
		parsedType, err = p.parseRecordType()
	} else {
		parsedType, err = p.parseType()
	}

	if err != nil {
		return nil, err
	}

	expectedSemicolon := p.consume(dt.SEMICOLON)
	if expectedSemicolon == nil {
		return nil, p.createParseError(dt.SEMICOLON, "type declaration must end with semicolon")
	}

	typeDeclaration := dt.ParseTree{
		RootType:   dt.TYPE_DECLARATION_NODE,
		TokenValue: nil,
		Children: []dt.ParseTree{
			{
				RootType:   dt.TOKEN_NODE,
				TokenValue: identifier,
				Children:   nil,
			},
			{
				RootType:   dt.TOKEN_NODE,
				TokenValue: expectedEqual,
				Children:   nil,
			},
			*parsedType,
			{
				RootType:   dt.TOKEN_NODE,
				TokenValue: expectedSemicolon,
				Children:   nil,
			},
		},
	}

	return &typeDeclaration, nil
}

func (p *Parser) parseVarDeclarationPart() (*dt.ParseTree, error) {
	expectedVar := p.consumeExact(dt.KEYWORD, "variabel")
	if expectedVar == nil {
		return nil, nil
	}

	varDeclaration := dt.ParseTree{
		RootType:   dt.VAR_DECLARATION_PART_NODE,
		TokenValue: nil,
		Children: []dt.ParseTree{
			{
				RootType:   dt.TOKEN_NODE,
				TokenValue: expectedVar,
				Children:   nil,
			},
		},
	}

	if !p.match(dt.IDENTIFIER) {
		return nil, p.createParseError(dt.IDENTIFIER, "expected at least one variable declaration")
	}

	for p.match(dt.IDENTIFIER) {
		variableDeclaration, err := p.parseVarDeclaration()

		if err != nil {
			return nil, err
		}

		varDeclaration.Children = append(varDeclaration.Children,
			*variableDeclaration,
		)
	}
	return &varDeclaration, nil
}

func (p *Parser) parseVarDeclaration() (*dt.ParseTree, error) {

	identifier, err := p.parseIdentifierList()

	if err != nil {
		return nil, err
	}

	expectedColon := p.consume(dt.COLON)
	if expectedColon == nil {
		return nil, p.createParseError(dt.COLON, "")
	}

	parsedType, err := p.parseType()

	if err != nil {
		return nil, err
	}

	expectedSemicolon := p.consume(dt.SEMICOLON)
	if expectedSemicolon == nil {
		return nil, p.createParseError(dt.SEMICOLON, "variable declaration must end with semicolon")
	}

	variableDeclaration := dt.ParseTree{
		RootType:   dt.VAR_DECLARATION_NODE,
		TokenValue: nil,
		Children: []dt.ParseTree{
			*identifier,
			{
				RootType:   dt.TOKEN_NODE,
				TokenValue: expectedColon,
				Children:   nil,
			},
			*parsedType,
			{
				RootType:   dt.TOKEN_NODE,
				TokenValue: expectedSemicolon,
				Children:   nil,
			},
		},
	}

	return &variableDeclaration, nil
}

func (p *Parser) parseIdentifierList() (*dt.ParseTree, error) {

	expectedIdentifier := p.consume(dt.IDENTIFIER)
	if expectedIdentifier == nil {
		return nil, p.createParseError(dt.IDENTIFIER, "")
	}

	identifierListTree := dt.ParseTree{
		RootType:   dt.IDENTIFIER_LIST_NODE,
		TokenValue: nil,
		Children: []dt.ParseTree{{
			RootType:   dt.TOKEN_NODE,
			TokenValue: expectedIdentifier,
			Children:   make([]dt.ParseTree, 0),
		}},
	}

	for {
		expectedComma := p.consume(dt.COMMA)
		if expectedComma == nil {
			break
		}
		expectedIdentifier := p.consume(dt.IDENTIFIER)
		if expectedIdentifier == nil {
			return nil, p.createParseError(dt.IDENTIFIER, "expected identifier after comma")
		}
		identifierListTree.Children = append(identifierListTree.Children,
			dt.ParseTree{
				RootType:   dt.TOKEN_NODE,
				TokenValue: expectedComma,
				Children:   make([]dt.ParseTree, 0),
			},
			dt.ParseTree{
				RootType:   dt.TOKEN_NODE,
				TokenValue: expectedIdentifier,
				Children:   make([]dt.ParseTree, 0),
			},
		)
	}

	return &identifierListTree, nil
}

func (p *Parser) parseType() (*dt.ParseTree, error) {
	typeTree := dt.ParseTree{
		RootType:   dt.TYPE_NODE,
		TokenValue: nil,
		Children:   make([]dt.ParseTree, 1),
	}

	if p.match(dt.IDENTIFIER) {
		typeTree.Children[0] = dt.ParseTree{
			RootType:   dt.TOKEN_NODE,
			TokenValue: p.consume(dt.IDENTIFIER),
			Children:   make([]dt.ParseTree, 0),
		}
	} else {
		if !p.match(dt.KEYWORD) {
			return nil, p.createParseError(dt.KEYWORD, "expected type")
		}

		switch p.peek().Lexeme {
		case "integer":
			fallthrough
		case "real":
			fallthrough
		case "boolean":
			fallthrough
		case "char":
			typeTree.Children[0] = dt.ParseTree{
				RootType:   dt.TOKEN_NODE,
				TokenValue: p.consume(dt.KEYWORD),
				Children:   make([]dt.ParseTree, 0),
			}
		case "larik":
			fallthrough
		case "array":
			arrayTypeTree, err := p.parseArrayType()

			if err != nil {
				return nil, err
			}

			typeTree.Children[0] = *arrayTypeTree
		}
	}

	return &typeTree, nil
}

func (p *Parser) parseArrayType() (*dt.ParseTree, error) {

	expectedLarik := p.consumeExact(dt.KEYWORD, "larik")
	if expectedLarik == nil {
		return nil, p.createParseError(dt.KEYWORD, "expected larik keyword")
	}

	expectedLB := p.consume(dt.LBRACKET)
	if expectedLB == nil {
		return nil, p.createParseError(dt.LBRACKET, "expected [ after larik")
	}

	rangeTree, err := p.parseRange()

	if err != nil {
		return nil, err
	}

	expectedRB := p.consume(dt.RBRACKET)
	if expectedRB == nil {
		return nil, p.createParseError(dt.RBRACKET, "expected ] after larik range")
	}

	expectedDari := p.consumeExact(dt.KEYWORD, "dari")
	if expectedDari == nil {
		return nil, p.createParseError(dt.KEYWORD, "expected 'dari' after ]")
	}

	typeTree, err := p.parseType()

	if err != nil {
		return nil, err
	}

	arrayTypeTree := dt.ParseTree{
		RootType:   dt.ARRAY_TYPE_NODE,
		TokenValue: nil,
		Children: []dt.ParseTree{{
			RootType:   dt.TOKEN_NODE,
			TokenValue: expectedLarik,
			Children:   make([]dt.ParseTree, 0),
		}, {
			RootType:   dt.TOKEN_NODE,
			TokenValue: expectedLB,
			Children:   make([]dt.ParseTree, 0),
		},
			*rangeTree,
			{
				RootType:   dt.TOKEN_NODE,
				TokenValue: expectedRB,
				Children:   make([]dt.ParseTree, 0),
			}, {
				RootType:   dt.TOKEN_NODE,
				TokenValue: expectedDari,
				Children:   make([]dt.ParseTree, 0),
			},
			*typeTree,
		},
	}

	return &arrayTypeTree, nil
}

func (p *Parser) parseRange() (*dt.ParseTree, error) {
	expression1, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	rangeOperator := p.consume(dt.RANGE_OPERATOR)
	if rangeOperator == nil {
		return nil, p.createParseError(dt.RANGE_OPERATOR, "expected '..' for range")
	}

	expression2, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	rangeTree := dt.ParseTree{
		RootType:   dt.RANGE_NODE,
		TokenValue: nil,
		Children: []dt.ParseTree{
			*expression1,
			{
				RootType:   dt.TOKEN_NODE,
				TokenValue: rangeOperator,
				Children:   make([]dt.ParseTree, 0),
			},
			*expression2,
		},
	}

	return &rangeTree, nil
}

func (p *Parser) parseStatement() (*dt.ParseTree, error) {
	if p.matchExact(dt.KEYWORD, "mulai") {
		return p.parseCompoundStatement()
	}
	if p.matchExact(dt.KEYWORD, "jika") {
		return p.parseIfStatement()
	}
	if p.matchExact(dt.KEYWORD, "selama") {
		return p.parseWhileStatement()
	}
	if p.matchExact(dt.KEYWORD, "untuk") {
		return p.parseForStatement()
	}
	if p.match(dt.IDENTIFIER) {
		if p.pos+1 < len(p.buffer) && (p.buffer[p.pos+1].Type == dt.ASSIGN_OPERATOR || p.buffer[p.pos+1].Type == dt.LBRACKET) {
			return p.parseAssignmentStatement()
		} else {
			return p.parseSubprogramCall()
		}
	}
	return nil, p.createParseErrorMany(
		[]dt.TokenType{dt.KEYWORD, dt.IDENTIFIER},
		"expected a statement (jika, selama, untuk, mulai, or identifier)",
	)
}

func (p *Parser) parseSubprogramDeclaration() (*dt.ParseTree, error) {
	procTree, err := p.parseProcedureDeclaration()
	if err != nil {
		return nil, err
	}
	if procTree != nil {
		return &dt.ParseTree{
			RootType: dt.SUBPROGRAM_DECLARATION_NODE,
			Children: []dt.ParseTree{*procTree},
		}, nil
	}

	funcTree, err := p.parseFunctionDeclaration()
	if err != nil {
		return nil, err
	}
	if funcTree != nil {
		return &dt.ParseTree{
			RootType: dt.SUBPROGRAM_DECLARATION_NODE,
			Children: []dt.ParseTree{*funcTree},
		}, nil
	}

	return nil, nil
}

func (p *Parser) parseProcedureDeclaration() (*dt.ParseTree, error) {
	prosedurToken := p.consumeExact(dt.KEYWORD, "prosedur")
	if prosedurToken == nil {
		return nil, nil
	}

	identifier := p.consume(dt.IDENTIFIER)
	if identifier == nil {
		return nil, p.createParseError(dt.IDENTIFIER, "expected procedure name")
	}

	procTree := dt.ParseTree{
		RootType: dt.PROCEDURE_DECLARATION_NODE,
		Children: []dt.ParseTree{
			{RootType: dt.TOKEN_NODE, TokenValue: prosedurToken},
			{RootType: dt.TOKEN_NODE, TokenValue: identifier},
		},
	}

	if p.match(dt.LPARENTHESIS) {
		paramList, err := p.parseFormalParameterList()
		if err != nil {
			return nil, err
		}
		procTree.Children = append(procTree.Children, *paramList)
	}

	semicolon1 := p.consume(dt.SEMICOLON)
	if semicolon1 == nil {
		return nil, p.createParseError(dt.SEMICOLON, "expected ';' after procedure header")
	}
	procTree.Children = append(procTree.Children, dt.ParseTree{RootType: dt.TOKEN_NODE, TokenValue: semicolon1})

	declarations, err := p.parseDeclarationPart()
	if err != nil {
		return nil, err
	}
	procTree.Children = append(procTree.Children, *declarations)

	compoundStmt, err := p.parseCompoundStatement()
	if err != nil {
		return nil, err
	}
	procTree.Children = append(procTree.Children, *compoundStmt)

	semicolon2 := p.consume(dt.SEMICOLON)
	if semicolon2 == nil {
		return nil, p.createParseError(dt.SEMICOLON, "expected ';' after procedure block")
	}
	procTree.Children = append(procTree.Children, dt.ParseTree{RootType: dt.TOKEN_NODE, TokenValue: semicolon2})

	return &procTree, nil
}

func (p *Parser) parseFunctionDeclaration() (*dt.ParseTree, error) {
	fungsiToken := p.consumeExact(dt.KEYWORD, "fungsi")
	if fungsiToken == nil {
		return nil, nil
	}

	identifier := p.consume(dt.IDENTIFIER)
	if identifier == nil {
		return nil, p.createParseError(dt.IDENTIFIER, "expected function name")
	}

	funcTree := dt.ParseTree{
		RootType: dt.FUNCTION_DECLARATION_NODE,
		Children: []dt.ParseTree{
			{RootType: dt.TOKEN_NODE, TokenValue: fungsiToken},
			{RootType: dt.TOKEN_NODE, TokenValue: identifier},
		},
	}

	if p.match(dt.LPARENTHESIS) {
		paramList, err := p.parseFormalParameterList()
		if err != nil {
			return nil, err
		}
		funcTree.Children = append(funcTree.Children, *paramList)
	}

	colonToken := p.consume(dt.COLON)
	if colonToken == nil {
		return nil, p.createParseError(dt.COLON, "expected ':' for function return type")
	}

	returnType, err := p.parseType()
	if err != nil {
		return nil, err
	}

	semicolon1 := p.consume(dt.SEMICOLON)
	if semicolon1 == nil {
		return nil, p.createParseError(dt.SEMICOLON, "expected ';' after function header")
	}

	funcTree.Children = append(funcTree.Children,
		dt.ParseTree{RootType: dt.TOKEN_NODE, TokenValue: colonToken},
		*returnType,
		dt.ParseTree{RootType: dt.TOKEN_NODE, TokenValue: semicolon1},
	)

	declarations, err := p.parseDeclarationPart()
	if err != nil {
		return nil, err
	}
	funcTree.Children = append(funcTree.Children, *declarations)

	compoundStmt, err := p.parseCompoundStatement()
	if err != nil {
		return nil, err
	}
	funcTree.Children = append(funcTree.Children, *compoundStmt)

	semicolon2 := p.consume(dt.SEMICOLON)
	if semicolon2 == nil {
		return nil, p.createParseError(dt.SEMICOLON, "expected ';' after function block")
	}
	funcTree.Children = append(funcTree.Children, dt.ParseTree{RootType: dt.TOKEN_NODE, TokenValue: semicolon2})

	return &funcTree, nil
}

func (p *Parser) parseFormalParameterList() (*dt.ParseTree, error) {
	lpToken := p.consume(dt.LPARENTHESIS)
	if lpToken == nil {
		return nil, p.createParseError(dt.LPARENTHESIS, "expected '(' to start parameter list")
	}

	paramListTree := dt.ParseTree{
		RootType: dt.FORMAL_PARAMETER_LIST_NODE,
		Children: []dt.ParseTree{
			{RootType: dt.TOKEN_NODE, TokenValue: lpToken},
		},
	}
	if !p.match(dt.RPARENTHESIS) {
		idList, err := p.parseIdentifierList()
		if err != nil {
			return nil, err
		}

		colonToken := p.consume(dt.COLON)
		if colonToken == nil {
			return nil, p.createParseError(dt.COLON, "expected ':' after identifier list in parameters")
		}

		typeNode, err := p.parseType()
		if err != nil {
			return nil, err
		}

		paramListTree.Children = append(paramListTree.Children,
			*idList,
			dt.ParseTree{RootType: dt.TOKEN_NODE, TokenValue: colonToken},
			*typeNode,
		)
		for p.match(dt.SEMICOLON) {
			semicolonToken := p.consume(dt.SEMICOLON)
			nextIdList, err := p.parseIdentifierList()
			if err != nil {
				return nil, p.createParseError(dt.IDENTIFIER, "expected identifier list after ';'")
			}

			nextColonToken := p.consume(dt.COLON)
			if nextColonToken == nil {
				return nil, p.createParseError(dt.COLON, "expected ':' after identifier list in parameters")
			}

			nextTypeNode, err := p.parseType()
			if err != nil {
				return nil, err
			}

			paramListTree.Children = append(paramListTree.Children,
				dt.ParseTree{RootType: dt.TOKEN_NODE, TokenValue: semicolonToken},
				*nextIdList,
				dt.ParseTree{RootType: dt.TOKEN_NODE, TokenValue: nextColonToken},
				*nextTypeNode,
			)
		}
	}

	rpToken := p.consume(dt.RPARENTHESIS)
	if rpToken == nil {
		return nil, p.createParseError(dt.RPARENTHESIS, "expected ')' to end parameter list")
	}
	paramListTree.Children = append(paramListTree.Children,
		dt.ParseTree{RootType: dt.TOKEN_NODE, TokenValue: rpToken},
	)

	return &paramListTree, nil
}

func (p *Parser) parseCompoundStatement() (*dt.ParseTree, error) {
	mulaiToken := p.consumeExact(dt.KEYWORD, "mulai")
	if mulaiToken == nil {
		return nil, p.createParseError(dt.KEYWORD, "expected 'mulai' keyword")
	}

	stmtList, err := p.parseStatementList()
	if err != nil {
		return nil, err
	}

	selesaiToken := p.consumeExact(dt.KEYWORD, "selesai")
	if selesaiToken == nil {
		return nil, p.createParseError(dt.KEYWORD, "expected 'selesai' keyword to end compound statement")
	}

	compoundTree := dt.ParseTree{
		RootType: dt.COMPOUND_STATEMENT_NODE,
		Children: []dt.ParseTree{
			{RootType: dt.TOKEN_NODE, TokenValue: mulaiToken},
			*stmtList,
			{RootType: dt.TOKEN_NODE, TokenValue: selesaiToken},
		},
	}
	return &compoundTree, nil
}

func (p *Parser) parseStatementList() (*dt.ParseTree, error) {
	stmt, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	stmtListTree := dt.ParseTree{
		RootType: dt.STATEMENT_LIST_NODE,
		Children: []dt.ParseTree{*stmt},
	}

	for p.match(dt.SEMICOLON) {
		semicolon := p.consume(dt.SEMICOLON)

		if p.matchExact(dt.KEYWORD, "selesai") {
			stmtListTree.Children = append(stmtListTree.Children, dt.ParseTree{
				RootType:   dt.TOKEN_NODE,
				TokenValue: semicolon,
			})
			break
		}

		nextStmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		stmtListTree.Children = append(stmtListTree.Children,
			dt.ParseTree{
				RootType:   dt.TOKEN_NODE,
				TokenValue: semicolon,
			},
			*nextStmt,
		)
	}
	return &stmtListTree, nil
}

func (p *Parser) parseAssignmentStatement() (*dt.ParseTree, error) {
	var lvalue dt.ParseTree

	p.pos++
	if p.match(dt.LBRACKET) {
		p.pos--

		plvalue, err := p.parseArrayAccess()

		if err != nil {
			return nil, err
		}

		lvalue = *plvalue
	} else {
		p.pos--

		identifier := p.consume(dt.IDENTIFIER)
		if identifier == nil {
			return nil, p.createParseError(dt.IDENTIFIER, "expected identifier for assignment")
		}

		lvalue = dt.ParseTree{
			RootType:   dt.TOKEN_NODE,
			TokenValue: identifier,
		}
	}

	assignOp := p.consume(dt.ASSIGN_OPERATOR)
	if assignOp == nil {
		return nil, p.createParseError(dt.ASSIGN_OPERATOR, "expected ':=' after identifier")
	}

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	assignTree := dt.ParseTree{
		RootType: dt.ASSIGNMENT_STATEMENT_NODE,
		Children: []dt.ParseTree{
			lvalue,
			{RootType: dt.TOKEN_NODE, TokenValue: assignOp},
			*expr,
		},
	}
	return &assignTree, nil
}

func (p *Parser) parseIfStatement() (*dt.ParseTree, error) {
	jikaToken := p.consumeExact(dt.KEYWORD, "jika")
	if jikaToken == nil {
		return nil, p.createParseError(dt.KEYWORD, "expected 'jika'")
	}

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	makaToken := p.consumeExact(dt.KEYWORD, "maka")
	if makaToken == nil {
		return nil, p.createParseError(dt.KEYWORD, "expected 'maka' after if-expression")
	}

	thenStmt, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	ifTree := dt.ParseTree{
		RootType: dt.IF_STATEMENT_NODE,
		Children: []dt.ParseTree{
			{RootType: dt.TOKEN_NODE, TokenValue: jikaToken},
			*expr,
			{RootType: dt.TOKEN_NODE, TokenValue: makaToken},
			*thenStmt,
		},
	}

	if p.matchExact(dt.KEYWORD, "selain_itu") {
		selainItuToken := p.consumeExact(dt.KEYWORD, "selain_itu")

		elseStmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		ifTree.Children = append(ifTree.Children,
			dt.ParseTree{RootType: dt.TOKEN_NODE, TokenValue: selainItuToken},
			*elseStmt,
		)
	}

	return &ifTree, nil
}

func (p *Parser) parseWhileStatement() (*dt.ParseTree, error) {
	selamaToken := p.consumeExact(dt.KEYWORD, "selama")
	if selamaToken == nil {
		return nil, p.createParseError(dt.KEYWORD, "expected 'selama'")
	}

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	lakukanToken := p.consumeExact(dt.KEYWORD, "lakukan")
	if lakukanToken == nil {
		return nil, p.createParseError(dt.KEYWORD, "expected 'lakukan' after while-expression")
	}

	stmt, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	whileTree := dt.ParseTree{
		RootType: dt.WHILE_STATEMENT_NODE,
		Children: []dt.ParseTree{
			{RootType: dt.TOKEN_NODE, TokenValue: selamaToken},
			*expr,
			{RootType: dt.TOKEN_NODE, TokenValue: lakukanToken},
			*stmt,
		},
	}
	return &whileTree, nil
}

func (p *Parser) parseForStatement() (*dt.ParseTree, error) {
	untukToken := p.consumeExact(dt.KEYWORD, "untuk")
	if untukToken == nil {
		return nil, p.createParseError(dt.KEYWORD, "expected 'untuk'")
	}

	identifier := p.consume(dt.IDENTIFIER)
	if identifier == nil {
		return nil, p.createParseError(dt.IDENTIFIER, "expected counter identifier after 'untuk'")
	}

	assignOp := p.consume(dt.ASSIGN_OPERATOR)
	if assignOp == nil {
		return nil, p.createParseError(dt.ASSIGN_OPERATOR, "expected ':=' after for-identifier")
	}

	startExpr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	// (KEYWORD(ke)|KEYWORD(turun-ke))
	var directionToken *dt.Token
	if p.matchExact(dt.KEYWORD, "ke") {
		directionToken = p.consumeExact(dt.KEYWORD, "ke")
	} else if p.matchExact(dt.KEYWORD, "turun_ke") {
		directionToken = p.consumeExact(dt.KEYWORD, "turun_ke")
	}

	if directionToken == nil {
		return nil, p.createParseError(dt.KEYWORD, "expected 'ke' or 'turun_ke' in for loop")
	}

	endExpr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	lakukanToken := p.consumeExact(dt.KEYWORD, "lakukan")
	if lakukanToken == nil {
		return nil, p.createParseError(dt.KEYWORD, "expected 'lakukan' in for loop")
	}

	stmt, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	forTree := dt.ParseTree{
		RootType: dt.FOR_STATEMENT_NODE,
		Children: []dt.ParseTree{
			{RootType: dt.TOKEN_NODE, TokenValue: untukToken},
			{RootType: dt.TOKEN_NODE, TokenValue: identifier},
			{RootType: dt.TOKEN_NODE, TokenValue: assignOp},
			*startExpr,
			{RootType: dt.TOKEN_NODE, TokenValue: directionToken},
			*endExpr,
			{RootType: dt.TOKEN_NODE, TokenValue: lakukanToken},
			*stmt,
		},
	}
	return &forTree, nil
}

func (p *Parser) parseSubprogramCall() (*dt.ParseTree, error) {
	identifier := p.consume(dt.IDENTIFIER)
	if identifier == nil {
		return nil, p.createParseError(dt.RANGE_OPERATOR, "expected function/procedure identifier")
	}

	subprogramCall := dt.ParseTree{
		RootType:   dt.SUBPROGRAM_CALL_NODE,
		TokenValue: nil,
		Children: []dt.ParseTree{
			{
				RootType:   dt.TOKEN_NODE,
				TokenValue: identifier,
				Children:   make([]dt.ParseTree, 0),
			},
		},
	}

	if p.match(dt.LPARENTHESIS) {
		lp := p.consume(dt.LPARENTHESIS)
		params, err := p.parseParameterList()
		if err != nil {
			return nil, err
		}
		rp := p.consume(dt.RPARENTHESIS)

		if rp == nil {
			return nil, p.createParseError(dt.RPARENTHESIS, "expected ) after function call")
		}

		subprogramCall.Children = append(subprogramCall.Children,
			dt.ParseTree{
				RootType:   dt.TOKEN_NODE,
				TokenValue: lp,
				Children:   nil,
			},
			*params,
			dt.ParseTree{
				RootType:   dt.TOKEN_NODE,
				TokenValue: rp,
				Children:   nil,
			},
		)
	}

	return &subprogramCall, nil
}

func (p *Parser) parseParameterList() (*dt.ParseTree, error) {
	expression, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	parameterListTree := dt.ParseTree{
		RootType:   dt.PARAMETER_LIST_NODE,
		TokenValue: nil,
		Children: []dt.ParseTree{
			*expression,
		},
	}

	for {
		expectedComma := p.consume(dt.COMMA)
		if expectedComma == nil {
			break
		}
		nextExpression, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		parameterListTree.Children = append(parameterListTree.Children,
			dt.ParseTree{
				RootType:   dt.TOKEN_NODE,
				TokenValue: expectedComma,
				Children:   make([]dt.ParseTree, 0),
			},
			*nextExpression,
		)
	}

	return &parameterListTree, nil
}

func (p *Parser) parseExpression() (*dt.ParseTree, error) {
	simpleExpression, err := p.parseSimpleExpression()

	if err != nil {
		return nil, err
	}

	expressionTree := dt.ParseTree{
		RootType:   dt.EXPRESSION_NODE,
		TokenValue: nil,
		Children: []dt.ParseTree{
			*simpleExpression,
		},
	}

	if p.match(dt.RELATIONAL_OPERATOR) {
		relationalOperator, err := p.parseRelationalOperator()

		if err != nil {
			return nil, p.createParseError(dt.RELATIONAL_OPERATOR, "expected relational opeartor after simple expression")
		}

		simpleExpression2, err := p.parseSimpleExpression()

		if err != nil {
			return nil, err
		}

		expressionTree.Children = append(expressionTree.Children,
			*relationalOperator,
			*simpleExpression2,
		)
	}

	return &expressionTree, nil
}

func (p *Parser) parseSimpleExpression() (*dt.ParseTree, error) {
	var expectedArithmeticOperator *dt.Token
	if p.matchExact(dt.ARITHMETIC_OPERATOR, "+") {
		expectedArithmeticOperator = p.consumeExact(dt.ARITHMETIC_OPERATOR, "+")
	} else if p.matchExact(dt.ARITHMETIC_OPERATOR, "-") {
		expectedArithmeticOperator = p.consumeExact(dt.ARITHMETIC_OPERATOR, "-")
	} else if p.match(dt.ARITHMETIC_OPERATOR) { // arithmetic operator selain +,-
		return nil, p.createParseError(dt.ARITHMETIC_OPERATOR, "only + and - operator are allowed before term")
	}
	simpleExpressionTree := dt.ParseTree{
		RootType:   dt.SIMPLE_EXPRESSION_NODE,
		TokenValue: nil,
		Children:   make([]dt.ParseTree, 0),
	}
	if expectedArithmeticOperator != nil {
		simpleExpressionTree.Children = append(simpleExpressionTree.Children,
			dt.ParseTree{
				RootType:   dt.TOKEN_NODE,
				TokenValue: expectedArithmeticOperator,
				Children:   make([]dt.ParseTree, 0),
			},
		)
	}

	term, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	simpleExpressionTree.Children = append(simpleExpressionTree.Children,
		*term,
	)
	for {
		additiveOperator, err := p.parseAdditiveOperator()
		if err != nil {
			break
		}

		nextTerm, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		simpleExpressionTree.Children = append(simpleExpressionTree.Children,
			*additiveOperator,
			*nextTerm,
		)
	}
	return &simpleExpressionTree, nil
}

func (p *Parser) parseTerm() (*dt.ParseTree, error) {
	factor, err := p.parseFactor()
	if err != nil {
		return nil, err
	}
	termTree := dt.ParseTree{
		RootType:   dt.TERM_NODE,
		TokenValue: nil,
		Children: []dt.ParseTree{
			*factor,
		},
	}
	for {
		multiplicativeOperator, err := p.parseMultiplicativeOperator()
		if err != nil {
			break
		}

		nextFactor, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		termTree.Children = append(termTree.Children,
			*multiplicativeOperator,
			*nextFactor,
		)
	}
	return &termTree, nil
}

func (p *Parser) parseFactor() (*dt.ParseTree, error) {
	factorTree := dt.ParseTree{
		RootType:   dt.FACTOR_NODE,
		TokenValue: nil,
		Children:   make([]dt.ParseTree, 0),
	}
	if p.match(dt.IDENTIFIER) {
		identifier, err := p.parseAccess()

		if err != nil {
			return nil, err
		}

		factorTree.Children = append(factorTree.Children, *identifier)
	} else if p.match(dt.NUMBER) {
		factor := p.consume(dt.NUMBER)
		factorTree.Children = append(factorTree.Children,
			dt.ParseTree{
				RootType:   dt.TOKEN_NODE,
				TokenValue: factor,
				Children:   nil,
			},
		)
	} else if p.match(dt.CHAR_LITERAL) {
		factor := p.consume(dt.CHAR_LITERAL)
		factorTree.Children = append(factorTree.Children,
			dt.ParseTree{
				RootType:   dt.TOKEN_NODE,
				TokenValue: factor,
				Children:   nil,
			},
		)
	} else if p.match(dt.STRING_LITERAL) {
		factor := p.consume(dt.STRING_LITERAL)
		factorTree.Children = append(factorTree.Children,
			dt.ParseTree{
				RootType:   dt.TOKEN_NODE,
				TokenValue: factor,
				Children:   nil,
			},
		)
	} else if p.match(dt.LPARENTHESIS) {
		lp := p.consume(dt.LPARENTHESIS)
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		rp := p.consume(dt.RPARENTHESIS)
		if rp == nil {
			return nil, p.createParseError(dt.RPARENTHESIS, "closing ) not found after expression")
		}
		factorTree.Children = append(factorTree.Children,
			dt.ParseTree{
				RootType:   dt.TOKEN_NODE,
				TokenValue: lp,
				Children:   nil,
			},
			*expr,
			dt.ParseTree{
				RootType:   dt.TOKEN_NODE,
				TokenValue: rp,
				Children:   nil,
			},
		)
	} else if p.matchExact(dt.LOGICAL_OPERATOR, "tidak") {
		expectedNot := p.consumeExact(dt.LOGICAL_OPERATOR, "tidak")
		factor, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		factorTree.Children = append(factorTree.Children,
			dt.ParseTree{
				RootType:   dt.TOKEN_NODE,
				TokenValue: expectedNot,
				Children:   nil,
			},
			*factor,
		)
	} else {
		return nil, p.createParseErrorMany([]dt.TokenType{dt.IDENTIFIER, dt.NUMBER, dt.CHAR_LITERAL, dt.STRING_LITERAL, dt.LPARENTHESIS}, "cannot parse factor")
	}

	return &factorTree, nil
}

func (p *Parser) parseRelationalOperator() (*dt.ParseTree, error) {
	expectedRelationalOperator := p.consume(dt.RELATIONAL_OPERATOR)

	if expectedRelationalOperator == nil {
		return nil, p.createParseError(dt.RELATIONAL_OPERATOR, "")
	}

	relationalOperator := dt.ParseTree{
		RootType:   dt.RELATIONAL_OPERATOR_NODE,
		TokenValue: nil,
		Children: []dt.ParseTree{
			{
				RootType:   dt.TOKEN_NODE,
				TokenValue: expectedRelationalOperator,
				Children:   make([]dt.ParseTree, 0),
			},
		},
	}

	return &relationalOperator, nil
}

func (p *Parser) parseAdditiveOperator() (*dt.ParseTree, error) {
	var expectedAdditionOperator *dt.Token
	if p.matchExact(dt.ARITHMETIC_OPERATOR, "+") {
		expectedAdditionOperator = p.consumeExact(dt.ARITHMETIC_OPERATOR, "+")
	} else if p.matchExact(dt.ARITHMETIC_OPERATOR, "-") {
		expectedAdditionOperator = p.consumeExact(dt.ARITHMETIC_OPERATOR, "-")
	} else if p.matchExact(dt.LOGICAL_OPERATOR, "atau") {
		expectedAdditionOperator = p.consumeExact(dt.LOGICAL_OPERATOR, "atau")
	} else {
		return nil, p.createParseErrorMany([]dt.TokenType{dt.ARITHMETIC_OPERATOR, dt.LOGICAL_OPERATOR}, "additive operator not found")
	}

	additiveOperator := dt.ParseTree{
		RootType:   dt.ADDITIVE_OPERATOR_NODE,
		TokenValue: nil,
		Children: []dt.ParseTree{
			{
				RootType:   dt.TOKEN_NODE,
				TokenValue: expectedAdditionOperator,
				Children:   make([]dt.ParseTree, 0),
			},
		},
	}

	return &additiveOperator, nil
}

func (p *Parser) parseMultiplicativeOperator() (*dt.ParseTree, error) {
	var expectedMultiplicativeOperator *dt.Token
	if p.matchExact(dt.ARITHMETIC_OPERATOR, "*") {
		expectedMultiplicativeOperator = p.consumeExact(dt.ARITHMETIC_OPERATOR, "*")
	} else if p.matchExact(dt.ARITHMETIC_OPERATOR, "/") {
		expectedMultiplicativeOperator = p.consumeExact(dt.ARITHMETIC_OPERATOR, "/")
	} else if p.matchExact(dt.ARITHMETIC_OPERATOR, "bagi") {
		expectedMultiplicativeOperator = p.consumeExact(dt.ARITHMETIC_OPERATOR, "bagi")
	} else if p.matchExact(dt.ARITHMETIC_OPERATOR, "mod") {
		expectedMultiplicativeOperator = p.consumeExact(dt.ARITHMETIC_OPERATOR, "mod")
	} else if p.matchExact(dt.LOGICAL_OPERATOR, "dan") {
		expectedMultiplicativeOperator = p.consumeExact(dt.LOGICAL_OPERATOR, "dan")
	} else {
		return nil, p.createParseErrorMany([]dt.TokenType{dt.ARITHMETIC_OPERATOR, dt.LOGICAL_OPERATOR}, "multiplicative operator not found")
	}

	multiplicativeOperator := dt.ParseTree{
		RootType:   dt.MULTIPLICATIVE_OPERATOR_NODE,
		TokenValue: nil,
		Children: []dt.ParseTree{
			{
				RootType:   dt.TOKEN_NODE,
				TokenValue: expectedMultiplicativeOperator,
				Children:   make([]dt.ParseTree, 0),
			},
		},
	}

	return &multiplicativeOperator, nil
}

func (p *Parser) parseAccess() (*dt.ParseTree, error) {
	access := dt.ParseTree{
		RootType:   dt.ACCESS_NODE,
		TokenValue: nil,
		Children:   make([]dt.ParseTree, 0),
	}

	if !p.match(dt.IDENTIFIER) {
		return nil, p.createParseError(dt.IDENTIFIER, "expected identifier")
	}

	var node *dt.ParseTree
	var err error

	if p.pos < len(p.buffer)-1 {
		switch p.buffer[p.pos+1].Type {
		case dt.LPARENTHESIS:
			node, err = p.parseSubprogramCall()
		case dt.LBRACKET:
			node, err = p.parseArrayAccess()
		default:
			node = &dt.ParseTree{
				RootType:   dt.TOKEN_NODE,
				TokenValue: p.consume(dt.IDENTIFIER),
			}
		}
	} else {
		node = &dt.ParseTree{
			RootType:   dt.TOKEN_NODE,
			TokenValue: p.consume(dt.IDENTIFIER),
		}
	}

	if err != nil {
		return nil, err
	}

	access.Children = append(access.Children, *node)

	for p.match(dt.DOT) {
		access.Children = append(access.Children, dt.ParseTree{
			RootType:   dt.TOKEN_NODE,
			TokenValue: p.consume(dt.DOT),
		})

		if p.pos < len(p.buffer)-1 {
			switch p.buffer[p.pos+1].Type {
			case dt.LPARENTHESIS:
				node, err = p.parseSubprogramCall()
			case dt.LBRACKET:
				node, err = p.parseArrayAccess()
			default:
				node = &dt.ParseTree{
					RootType:   dt.TOKEN_NODE,
					TokenValue: p.consume(dt.IDENTIFIER),
				}
			}
		} else {
			node = &dt.ParseTree{
				RootType:   dt.TOKEN_NODE,
				TokenValue: p.consume(dt.IDENTIFIER),
			}
		}

		if err != nil {
			return nil, err
		}

		access.Children = append(access.Children, *node)
	}

	return &access, nil
}

func (p *Parser) parseArrayAccess() (*dt.ParseTree, error) {
	arrayAccess := dt.ParseTree{
		RootType:   dt.ARRAY_ACCESS_NODE,
		TokenValue: nil,
		Children:   make([]dt.ParseTree, 0, 4),
	}

	if !p.match(dt.IDENTIFIER) {
		return nil, p.createParseError(dt.IDENTIFIER, "expected array identifier")
	}

	arrayAccess.Children = append(arrayAccess.Children, dt.ParseTree{
		RootType:   dt.TOKEN_NODE,
		TokenValue: p.consume(dt.IDENTIFIER),
		Children:   make([]dt.ParseTree, 0),
	})

	if !p.match(dt.LBRACKET) {
		return nil, p.createParseError(dt.LBRACKET, "expected [")
	}

	arrayAccess.Children = append(arrayAccess.Children, dt.ParseTree{
		RootType:   dt.TOKEN_NODE,
		TokenValue: p.consume(dt.LBRACKET),
		Children:   make([]dt.ParseTree, 0),
	})

	expression, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	arrayAccess.Children = append(arrayAccess.Children,
		*expression,
	)

	if !p.match(dt.RBRACKET) {
		return nil, p.createParseError(dt.RBRACKET, "expected ]")
	}

	arrayAccess.Children = append(arrayAccess.Children, dt.ParseTree{
		RootType:   dt.TOKEN_NODE,
		TokenValue: p.consume(dt.RBRACKET),
		Children:   make([]dt.ParseTree, 0),
	})

	for p.match(dt.LBRACKET) {
		arrayAccess.Children = append(arrayAccess.Children, dt.ParseTree{
			RootType:   dt.TOKEN_NODE,
			TokenValue: p.consume(dt.LBRACKET),
			Children:   make([]dt.ParseTree, 0),
		})

		expression, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		arrayAccess.Children = append(arrayAccess.Children,
			*expression,
		)

		if !p.match(dt.RBRACKET) {
			return nil, p.createParseError(dt.RBRACKET, "expected ]")
		}

		arrayAccess.Children = append(arrayAccess.Children, dt.ParseTree{
			RootType:   dt.TOKEN_NODE,
			TokenValue: p.consume(dt.RBRACKET),
			Children:   make([]dt.ParseTree, 0),
		})
	}

	return &arrayAccess, nil
}

func (p *Parser) parseRecordType() (*dt.ParseTree, error) {
	record := p.consumeExact(dt.KEYWORD, "rekaman")

	if record == nil {
		return nil, p.createParseError(dt.KEYWORD, "expected rekaman keyword")
	}

	recordType := dt.ParseTree{
		RootType:   dt.VAR_DECLARATION_PART_NODE,
		TokenValue: nil,
		Children: []dt.ParseTree{
			{
				RootType:   dt.TOKEN_NODE,
				TokenValue: record,
				Children:   nil,
			},
		},
	}

	if !p.match(dt.IDENTIFIER) {
		return nil, p.createParseError(dt.IDENTIFIER, "expected at least one field declaration")
	}

	for p.match(dt.IDENTIFIER) {
		variableDeclaration, err := p.parseVarDeclaration()

		if err != nil {
			return nil, err
		}

		recordType.Children = append(recordType.Children,
			*variableDeclaration,
		)
	}
	return &recordType, nil
}
