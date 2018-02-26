package htmlLoader

import "testing"

func TestParser(t *testing.T) {
	cases := []struct{ src, exp string }{
		{
			"<foo/>",
			"<foo/>",
		},
		{
			"<meta foo>",
			"<meta foo/>",
		},
		{
			"<a href='foo'>link</a>",
			"<a href='foo'>link</a>",
		},
		{
			"<    a    href   =   'foo'    > link        <   /   a   >",
			"<a href='foo'> link </a>",
		},
		{
			"<!DOCTYPE html><html/>",
			"<!doctype html/><html/>",
		},
		{
			"<html><head/><body/></html>",
			"<html><head/><body/></html>",
		},
		{
			"<ThIsIsFINe></THISISFINE>",
			"<thisisfine/>",
		},
	}

	for _, c := range cases {
		toks := lex([]byte(c.src))
		ast := parseTokens(toks)
		text := printAst(ast)
		if text != c.exp {
			t.Error(ast)
			t.Errorf("Expected %s, got %s", c.exp, text)
		}
	}
}
