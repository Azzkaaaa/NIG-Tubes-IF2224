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

	programTree := dt.ParseTree{
		RootType:   dt.PROGRAM_NODE,
		TokenValue: nil,
		Children: []dt.ParseTree{
			*headerTree,
			*declarationTree,
			{
				RootType:   dt.TOKEN_NODE,
				TokenValue: p.peek(),
				Children:   make([]dt.ParseTree, 0),
			},
		},
	}

	return &programTree, nil

	// statementTree, remainder, err := parseCompoundStatement(remainder)

	// if err != nil {
	// 	return nil, nil, err
	// }

	// if len(remainder) == 0 {
	// 	return nil, nil, errors.New("missing dot at eof")
	// }

	// if remainder[0].Type != dt.DOT {
	// 	return nil, nil, errors.New("program does not end at dot")
	// }

	// if len(remainder) > 1 {
	// 	return nil, nil, errors.New("program should end after dot")
	// }

	// programTree := dt.ParseTree{
	// 	RootType:   dt.PROGRAM_NODE,
	// 	TokenValue: nil,
	// 	Children: []dt.ParseTree{
	// 		*headerTree,
	// 		*declarationTree,
	// 		*statementTree,
	// 		{
	// 			RootType:   dt.TOKEN_NODE,
	// 			TokenValue: &remainder[1],
	// 			Children:   make([]dt.ParseTree, 0),
	// 		},
	// 	},
	// }

	// return &programTree, make([]dt.Token, 0), nil
}

func (p *Parser) parseProgramHeader() (*dt.ParseTree, error) {
	// if len(tokens) < 3 {
	// 	return nil, nil, errors.New("not enough tokens to form program header")
	// }

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

	// for err == nil {
	// 	subprogramDeclaration, newRemainder, newErr := parseSubprogramDeclaration(remainder)
	// 	err = newErr
	// 	if err == nil {
	// 		remainder = newRemainder
	// 		declarationTree.Children = append(declarationTree.Children, *subprogramDeclaration)
	// 	}
	// }

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

	parsedType, err := p.parseType()
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
	// if len(tokens) < 1 {
	// 	return nil, errors.New("did not find identifier")
	// }

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
	// if len(tokens) < 1 {
	// 	return nil, errors.New("type keyword not found")
	// }

	if !p.match(dt.KEYWORD) {
		return nil, p.createParseError(dt.KEYWORD, "keyword not found")
	}

	typeTree := dt.ParseTree{
		RootType:   dt.TYPE_NODE,
		TokenValue: nil,
		Children:   make([]dt.ParseTree, 1),
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

	return &typeTree, nil
}

func (p *Parser) parseArrayType() (*dt.ParseTree, error) {
	// if len(tokens) < 6 {
	// 	return nil, errors.New("insufficient tokens to construct array type")
	// }

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
	return &dt.ParseTree{}, nil
}

func (p *Parser) parseSubprogramDeclaration() (*dt.ParseTree, error) {
	return &dt.ParseTree{}, nil
}

func (p *Parser) parseProcedureDeclaration() (*dt.ParseTree, error) {
	return &dt.ParseTree{}, nil
}

func (p *Parser) parseFunctionDeclaration() (*dt.ParseTree, error) {
	return &dt.ParseTree{}, nil
}

func (p *Parser) parseFormalParameterList() (*dt.ParseTree, error) {
	return &dt.ParseTree{}, nil
}

func (p *Parser) parseCompoundStatement() (*dt.ParseTree, error) {
	return &dt.ParseTree{}, nil
}

func (p *Parser) parseStatementList() (*dt.ParseTree, error) {
	return &dt.ParseTree{}, nil
}

func (p *Parser) parseAssignmentStatement() (*dt.ParseTree, error) {
	return &dt.ParseTree{}, nil
}

func (p *Parser) parseIfStatement() (*dt.ParseTree, error) {
	return &dt.ParseTree{}, nil
}

func (p *Parser) parseWhileStatement() (*dt.ParseTree, error) {
	return &dt.ParseTree{}, nil
}

func (p *Parser) parseForStatement() (*dt.ParseTree, error) {
	return &dt.ParseTree{}, nil
}

func (p *Parser) parseSubprogramCall() (*dt.ParseTree, error) {
	return &dt.ParseTree{}, nil
}

func (p *Parser) parseParameterList() (*dt.ParseTree, error) {
	return &dt.ParseTree{}, nil
}

func (p *Parser) parseExpression() (*dt.ParseTree, error) {
	return &dt.ParseTree{}, nil
}

func (p *Parser) parseSimpleExpression() (*dt.ParseTree, error) {
	return &dt.ParseTree{}, nil
}

func (p *Parser) parseTerm() (*dt.ParseTree, error) {
	return &dt.ParseTree{}, nil
}

func (p *Parser) parseFactor() (*dt.ParseTree, error) {
	return &dt.ParseTree{}, nil
}

func (p *Parser) parseRelationalOperator() (*dt.ParseTree, error) {
	return &dt.ParseTree{}, nil
}

func (p *Parser) parseAdditiveOperator() (*dt.ParseTree, error) {
	return &dt.ParseTree{}, nil
}

func (p *Parser) parseMultiplicativeOperator() (*dt.ParseTree, error) {
	return &dt.ParseTree{}, nil
}
