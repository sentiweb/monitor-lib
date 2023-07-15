package formatter

import(
	"github.com/sentiweb/monitor-lib/notify/types"	
)

type DefaultFormatterFactory struct {
	
}

// Get provides formatter for a given notifier
// Default formatter provides the same text regardless the notifier
func (f *DefaultFormatterFactory) Get(notifierName string) types.Formatter {
	return genericFormatter
}

var (
	genericFormatter types.Formatter = &GenericFormatter{}
	defaultFormatterFactory types.FormatterFactory = &DefaultFormatterFactory{}
)

// Change the formatter Factory
func SetDefaultFactory(factory types.FormatterFactory) {
	defaultFormatterFactory = factory
}

func Get(notifierName string) types.Formatter {
	return defaultFormatterFactory.Get(notifierName)
}
