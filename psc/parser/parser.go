package parser

import (
	"errors"
	"fmt"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

type Parser struct {
	buffer []dt.Token
	pos    int
}

func New(tokens []dt.Token) *Parser {
	return &Parser{
		buffer: tokens,
		pos:    0,
	}
}

func (p *Parser) Consume(expectedType dt.TokenType) error {
	curr := p.buffer[p.pos]
	if curr.Type == expectedType {
		p.pos++
		return nil
	}
	return fmt.Errorf("Line %d, Col %d: Unexpected token %s (expected: %s)\n", curr.Line, curr.Col, curr.Lexeme, expectedType)
}

func (p *Parser) ConsumeExact(expectedType dt.TokenType, expectedLexeme string) error {
	curr := p.buffer[p.pos]
	if curr.Type == expectedType && curr.Lexeme == expectedLexeme {
		p.pos++
		return nil
	}
	return fmt.Errorf("Line %d, Col %d: Unexpected token %s (expected: %s)\n", curr.Line, curr.Col, curr.Lexeme, expectedLexeme)
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

	// declarationTree, remainder, err := parseDeclarationPart(remainder)

	// if err != nil {
	// 	return nil, nil, err
	// }
	programTree := dt.ParseTree{
		RootType:   dt.PROGRAM_NODE,
		TokenValue: nil,
		Children: []dt.ParseTree{
			*headerTree,
			{
				RootType:   dt.TOKEN_NODE,
				TokenValue: &p.buffer[p.pos],
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

	err := p.ConsumeExact(dt.KEYWORD, "program")
	if err != nil {
		return nil, err
	}

	err = p.Consume(dt.IDENTIFIER)
	if err != nil {
		return nil, err
	}

	err = p.Consume(dt.SEMICOLON)
	if err != nil {
		return nil, err
	}

	// if tokens[0].Type != dt.KEYWORD || tokens[0].Lexeme != "program" {
	// 	return nil, nil, errors.New("all programs must start with program keyword")
	// }

	// if tokens[1].Type != dt.IDENTIFIER {
	// 	return nil, nil, errors.New("program name must only use alphanumerical characters and underscores")
	// }

	// if tokens[2].Type != dt.SEMICOLON {
	// 	return nil, nil, errors.New("program name must be a single word and strictly end with ;")
	// }

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

func parseDeclarationPart(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	remainder := tokens
	var err error

	declarationTree := dt.ParseTree{
		RootType:   dt.DECLARATION_PART_NODE,
		TokenValue: nil,
		Children:   make([]dt.ParseTree, 0),
	}

	// for err == nil {
	// 	constDeclaration, newRemainder, newErr := parseConstDeclaration(remainder)
	// 	err = newErr
	// 	if err == nil {
	// 		remainder = newRemainder
	// 		declarationTree.Children = append(declarationTree.Children, *constDeclaration)
	// 	}
	// }

	// for err == nil {
	// 	typeDeclaration, newRemainder, newErr := parseTypeDeclaration(remainder)
	// 	err = newErr
	// 	if err == nil {
	// 		remainder = newRemainder
	// 		declarationTree.Children = append(declarationTree.Children, *typeDeclaration)
	// 	}
	// }

	for err == nil {
		varDeclaration, newRemainder, newErr := parseVarDeclaration(remainder)
		err = newErr
		if err == nil {
			remainder = newRemainder
			declarationTree.Children = append(declarationTree.Children, *varDeclaration)
		}
	}

	// fmt.Println("4")
	// for err == nil {
	// 	subprogramDeclaration, newRemainder, newErr := parseSubprogramDeclaration(remainder)
	// 	err = newErr
	// 	if err == nil {
	// 		remainder = newRemainder
	// 		declarationTree.Children = append(declarationTree.Children, *subprogramDeclaration)
	// 	}
	// }

	return &declarationTree, nil, nil
}

func parseConstDeclaration(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	return &dt.ParseTree{}, nil, nil
}

func parseTypeDeclaration(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	return &dt.ParseTree{}, nil, nil
}

func parseVarDeclaration(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	return &dt.ParseTree{}, nil, nil
}

func parseIdentifierList(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	if len(tokens) < 1 {
		return nil, nil, errors.New("did not find identifier")
	}

	if tokens[0].Type != dt.IDENTIFIER {
		return nil, nil, errors.New("expected identifier")
	}

	identifierListTree := dt.ParseTree{
		RootType:   dt.IDENTIFIER_LIST_NODE,
		TokenValue: nil,
		Children: []dt.ParseTree{{
			RootType:   dt.TOKEN_NODE,
			TokenValue: &tokens[0],
			Children:   make([]dt.ParseTree, 0),
		}},
	}

	remainder := tokens[1:]

	for len(remainder) > 1 && remainder[0].Type == dt.COMMA && remainder[1].Type == dt.IDENTIFIER {
		identifierListTree.Children = append(identifierListTree.Children,
			dt.ParseTree{
				RootType:   dt.TOKEN_NODE,
				TokenValue: &remainder[0],
				Children:   make([]dt.ParseTree, 0),
			},
			dt.ParseTree{
				RootType:   dt.TOKEN_NODE,
				TokenValue: &remainder[1],
				Children:   make([]dt.ParseTree, 0),
			},
		)

		remainder = remainder[2:]
	}

	return &identifierListTree, remainder, nil
}

func parseType(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	if len(tokens) < 1 {
		return nil, nil, errors.New("type keyword not found")
	}

	if tokens[0].Type != dt.KEYWORD {
		return nil, nil, errors.New("keyword not found")
	}

	typeTree := dt.ParseTree{
		RootType:   dt.TYPE_NODE,
		TokenValue: nil,
		Children:   make([]dt.ParseTree, 1),
	}

	var remainder []dt.Token

	switch tokens[0].Lexeme {
	case "integer":
		fallthrough
	case "real":
		fallthrough
	case "boolean":
		fallthrough
	case "char":
		remainder = tokens[1:]
		typeTree.Children[0] = dt.ParseTree{
			RootType:   dt.TOKEN_NODE,
			TokenValue: &tokens[0],
			Children:   make([]dt.ParseTree, 0),
		}
	default:
		arrayTypeTree, newRemainder, err := parseArrayType(tokens)

		if err != nil {
			return nil, nil, err
		}

		typeTree.Children[0] = *arrayTypeTree
		remainder = newRemainder
	}

	return &typeTree, remainder, nil
}

func parseArrayType(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	if len(tokens) < 6 {
		return nil, nil, errors.New("insufficient tokens to construct array type")
	}

	if tokens[0].Type != dt.KEYWORD || tokens[0].Lexeme != "larik" {
		return nil, nil, errors.New("expected keyword larik")
	}

	if tokens[1].Type != dt.LBRACKET {
		return nil, nil, errors.New("expected '['")
	}

	rangeTree, rangeRemainder, err := parseRange(tokens[2:])

	if err != nil {
		return nil, nil, err
	}

	if rangeRemainder[0].Type != dt.RBRACKET {
		return nil, nil, errors.New("expected ']'")
	}

	if rangeRemainder[1].Type != dt.KEYWORD || rangeRemainder[1].Lexeme != "dari" {
		return nil, nil, errors.New("expected 'dari'")
	}

	typeTree, remainder, err := parseType(rangeRemainder[2:])

	if err != nil {
		return nil, nil, err
	}

	arrayTypeTree := dt.ParseTree{
		RootType:   dt.ARRAY_TYPE_NODE,
		TokenValue: nil,
		Children: []dt.ParseTree{{
			RootType:   dt.TOKEN_NODE,
			TokenValue: &tokens[0],
			Children:   make([]dt.ParseTree, 0),
		}, {
			RootType:   dt.TOKEN_NODE,
			TokenValue: &tokens[1],
			Children:   make([]dt.ParseTree, 0),
		}, *rangeTree, {
			RootType:   dt.TOKEN_NODE,
			TokenValue: &rangeRemainder[0],
			Children:   make([]dt.ParseTree, 0),
		}, {
			RootType:   dt.TOKEN_NODE,
			TokenValue: &rangeRemainder[1],
			Children:   make([]dt.ParseTree, 0),
		}, *typeTree,
		},
	}

	return &arrayTypeTree, remainder, nil
}

func parseRange(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	return &dt.ParseTree{}, nil, nil
}

func parseSubprogramDeclaration(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	return &dt.ParseTree{}, nil, nil
}

func parseProcedureDeclaration(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	return &dt.ParseTree{}, nil, nil
}

func parseFunctionDeclaration(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	return &dt.ParseTree{}, nil, nil
}

func parseFormalParameterList(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	return &dt.ParseTree{}, nil, nil
}

func parseCompoundStatement(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	return &dt.ParseTree{}, nil, nil
}

func parseStatementList(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	return &dt.ParseTree{}, nil, nil
}

func parseAssignmentStatement(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	return &dt.ParseTree{}, nil, nil
}

func parseIfStatement(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	return &dt.ParseTree{}, nil, nil
}

func parseWhileStatement(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	return &dt.ParseTree{}, nil, nil
}

func parseForStatement(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	return &dt.ParseTree{}, nil, nil
}

func parseSubprogramCall(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	return &dt.ParseTree{}, nil, nil
}

func parseParameterList(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	return &dt.ParseTree{}, nil, nil
}

func parseExpression(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	return &dt.ParseTree{}, nil, nil
}

func parseSimpleExpression(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	return &dt.ParseTree{}, nil, nil
}

func parseTerm(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	return &dt.ParseTree{}, nil, nil
}

func parseFactor(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	return &dt.ParseTree{}, nil, nil
}

func parseRelationalOperator(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	return &dt.ParseTree{}, nil, nil
}

func parseAdditiveOperator(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	return &dt.ParseTree{}, nil, nil
}

func parseMultiplicativeOperator(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	return &dt.ParseTree{}, nil, nil
}
