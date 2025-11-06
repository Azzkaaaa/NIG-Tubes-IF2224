package parser

import (
	"errors"

	dt "github.com/Azzkaaaa/NIG-Tubes-IF2224/psc/datatype"
)

func Parse(tokens []dt.Token) (*dt.ParseTree, error) {
	tree, _, err := parseProgram(tokens)
	return tree, err
}

func parseProgram(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	headerTree, remainder, err := parseProgramHeader(tokens)

	if err != nil {
		return nil, nil, err
	}

	declarationTree, remainder, err := parseDeclarationPart(remainder)

	if err != nil {
		return nil, nil, err
	}

	statementTree, remainder, err := parseCompoundStatement(remainder)

	if err != nil {
		return nil, nil, err
	}

	if len(remainder) == 0 {
		return nil, nil, errors.New("missing dot at eof")
	}

	if remainder[0].Type != dt.DOT {
		return nil, nil, errors.New("program does not end at dot")
	}

	if len(remainder) > 1 {
		return nil, nil, errors.New("program should end after dot")
	}

	programTree := dt.ParseTree{
		RootType:   dt.PROGRAM_NODE,
		TokenValue: nil,
		Children: []dt.ParseTree{
			*headerTree,
			*declarationTree,
			*statementTree,
			{
				RootType:   dt.TOKEN_NODE,
				TokenValue: &remainder[1],
				Children:   make([]dt.ParseTree, 0),
			},
		},
	}

	return &programTree, make([]dt.Token, 0), nil
}

func parseProgramHeader(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	if len(tokens) != 3 {
		return nil, nil, errors.New("not enough tokens to form program header")
	}

	if tokens[0].Type != dt.KEYWORD || tokens[0].Lexeme != "program" {
		return nil, nil, errors.New("all programs must start with program keyword")
	}

	if tokens[1].Type != dt.IDENTIFIER {
		return nil, nil, errors.New("program name must only use alphanumerical characters and underscores")
	}

	if tokens[2].Type != dt.SEMICOLON {
		return nil, nil, errors.New("program name must be a single word and strictly end with ;")
	}

	headerTree := dt.ParseTree{
		RootType:   dt.PROGRAM_HEADER_NODE,
		TokenValue: nil,
		Children: []dt.ParseTree{{
			RootType:   dt.TOKEN_NODE,
			TokenValue: &tokens[0],
			Children:   make([]dt.ParseTree, 0),
		}, {
			RootType:   dt.TOKEN_NODE,
			TokenValue: &tokens[1],
			Children:   make([]dt.ParseTree, 0),
		}, {
			RootType:   dt.TOKEN_NODE,
			TokenValue: &tokens[2],
			Children:   make([]dt.ParseTree, 0),
		}},
	}

	return &headerTree, tokens[3:], nil
}

func parseDeclarationPart(tokens []dt.Token) (*dt.ParseTree, []dt.Token, error) {
	remainder := tokens
	var err error

	declarationTree := dt.ParseTree{
		RootType:   dt.DECLARATION_PART_NODE,
		TokenValue: nil,
		Children:   make([]dt.ParseTree, 0),
	}

	for err == nil {
		constDeclaration, newRemainder, newErr := parseConstDeclaration(remainder)
		err = newErr
		if err == nil {
			remainder = newRemainder
			declarationTree.Children = append(declarationTree.Children, *constDeclaration)
		}
	}

	for err == nil {
		typeDeclaration, newRemainder, newErr := parseTypeDeclaration(remainder)
		err = newErr
		if err == nil {
			remainder = newRemainder
			declarationTree.Children = append(declarationTree.Children, *typeDeclaration)
		}
	}

	for err == nil {
		varDeclaration, newRemainder, newErr := parseVarDeclaration(remainder)
		err = newErr
		if err == nil {
			remainder = newRemainder
			declarationTree.Children = append(declarationTree.Children, *varDeclaration)
		}
	}

	for err == nil {
		subprogramDeclaration, newRemainder, newErr := parseSubprogramDeclaration(remainder)
		err = newErr
		if err == nil {
			remainder = newRemainder
			declarationTree.Children = append(declarationTree.Children, *subprogramDeclaration)
		}
	}

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
