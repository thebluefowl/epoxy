package processor

var processors = make(map[uint8]Processor)

func RegisterProcessor(prefix uint8, processor Processor) {
	processors[prefix] = processor
}

func GetProcessor(prefix uint8) Processor {
	return processors[prefix]
}
