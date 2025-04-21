package errcode

import (
	"github.com/sohaha/zlsgo/zerror"
)

type ErrCode zerror.ErrCode

func (code ErrCode) New() error {
	return ErrorMsg(code, "")
}

func (code ErrCode) WrapText(msg string, err ...error) error {
	return ErrorMsg(code, msg, err...)
}

func (code ErrCode) WrapErr(err error) error {
	return ErrorMsg(code, "", err)
}

func ErrorMsg(code ErrCode, text string, err ...error) error {
	if text == "" {
		text, _ = GetI18n(code)
	}

	var tags []zerror.External

	// 指定 HTTP 状态码
	switch code {
	case Unauthorized, AuthorizedExpires:
		tags = append(tags, zerror.WrapTag(zerror.Unauthorized))
	case PermissionDenied:
		tags = append(tags, zerror.WrapTag(zerror.PermissionDenied))
	case InvalidInput:
		tags = append(tags, zerror.WrapTag(zerror.InvalidInput))
	}

	if len(err) > 0 {
		return zerror.Wrap(err[0], zerror.ErrCode(code), text, tags...)
	}

	return zerror.New(zerror.ErrCode(code), text, tags...)
}

func SetI18n(n map[ErrCode]string, lang ...string) {
	l := DefaultLang
	if len(lang) > 0 {
		l = lang[0]
	}

	for c, v := range n {
		if _, ok := i18n[l]; !ok {
			i18n[l] = map[ErrCode]string{}
		}
		i18n[l][c] = v
	}
}

func GetI18n(n ErrCode, lang ...string) (string, bool) {
	l := DefaultLang
	if len(lang) > 0 {
		l = lang[0]
	}
	if _, ok := i18n[l]; !ok {
		i18n[l] = map[ErrCode]string{}
	}

	t, ok := i18n[l][n]
	if !ok {
		t = "no defined"
	}
	return t, ok
}
