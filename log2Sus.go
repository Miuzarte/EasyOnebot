package EasyOnebot

import "fmt"

type log2Sus struct {
	Send func(any) error
}

func (b *Bot) wrapLogFunc() func(any) error {
	return func(msg any) error {
		for _, su := range b.sus {
			_, err := b.Call().Std.SendPrivateMsg(su, msg)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func (l *log2Sus) Error(s ...any) error {
	ss := fmt.Sprint(s...)
	return l.Send("[Error] " + ss)
}

func (l *log2Sus) Errorf(format string, a ...any) error {
	return l.Error(fmt.Sprintf(format, a...))
}

func (l *log2Sus) Warn(s ...any) error {
	ss := fmt.Sprint(s...)
	return l.Send("[Warn] " + ss)
}

func (l *log2Sus) Warnf(format string, a ...any) error {
	return l.Warn(fmt.Sprintf(format, a...))
}

func (l *log2Sus) Info(s ...any) error {
	ss := fmt.Sprint(s...)
	return l.Send("[Info] " + ss)
}

func (l *log2Sus) Infof(format string, a ...any) error {
	return l.Info(fmt.Sprintf(format, a...))
}

func (l *log2Sus) Debug(s ...any) error {
	ss := fmt.Sprint(s...)
	return l.Send("[Debug] " + ss)
}

func (l *log2Sus) Debugf(format string, a ...any) error {
	return l.Debug(fmt.Sprintf(format, a...))
}

func (l *log2Sus) Trace(s ...any) error {
	ss := fmt.Sprint(s...)
	return l.Send("[Trace] " + ss)
}

func (l *log2Sus) Tracef(format string, a ...any) error {
	return l.Trace(fmt.Sprintf(format, a...))
}

func (l *log2Sus) Print(s ...any) error {
	ss := fmt.Sprint(s...)
	return l.Send(ss)
}

func (l *log2Sus) Printf(format string, a ...any) error {
	return l.Print(fmt.Sprintf(format, a...))
}
