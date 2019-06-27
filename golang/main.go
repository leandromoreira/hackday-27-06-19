package main

import 	(
			"fmt"

			"github.com/giorgisio/goav/avformat"
			"github.com/giorgisio/goav/avcodec"
		)

func main() {
	var packet avcodec.Packet
	inputFilename := "chunk.ts"
	outputFilename := "chunk.mp4"

	avformat.AvRegisterAll()

	ctxInput := avformat.AvformatAllocContext()
	ctxOutput := avformat.AvformatAllocContext()

	// Open video file
	if avformat.AvformatOpenInput(&ctxInput, inputFilename, nil, nil) != 0 {
		fmt.Println("Error: Couldn't open file.")
		return
	}

	// Retrieve stream information
	if ctxInput.AvformatFindStreamInfo(nil) < 0 {
		fmt.Println("Error: Couldn't find input stream information.")
		ctxInput.AvformatCloseInput()
		return
	}

	avformat.AvformatAllocOutputContext2(&ctxOutput, nil, "mp4", outputFilename)

	for range(ctxInput.Streams()) {
		outStream := ctxOutput.AvformatNewStream(ctxInput.VideoCodec())
		if outStream == nil {
			fmt.Printf("Failed allocating output stream\n");
			return
		}
	}

	fmt.Println("===== INPUT DUMP =====")
	ctxInput.AvDumpFormat(0, inputFilename, 0)
	fmt.Printf("\n===== OUTPUT DUMP =====\n")
	ctxOutput.AvDumpFormat(0, outputFilename, 1)

	_, err := avformat.AvIOOpen(outputFilename, avformat.AVIO_FLAG_WRITE)
	if err != nil {
		fmt.Println("Erro no AVIO")
		return
	}

	opts := &avformat.Dictionary{}


	ret := ctxOutput.AvformatWriteHeader(&opts)
	if ret < 0 {
		fmt.Println("Error ocurred when opening output file")
		return
	}

	for {
		ret := ctxInput.AvReadFrame(&packet)
		if ret < 0 {
			break
		}

		inputStream := ctxInput.Streams()[packet.StreamIndex()]
		outputStream := ctxOutput.Streams()[packet.StreamIndex()]

		fmt.Println("Error ocurred when opening output file")

		packet.AvPacketRescaleTs(inputStream.TimeBase(), outputStream.TimeBase())
		ret = ctxOutput.AvInterleavedWriteFrame(&packet)

		if ret < 0 {
			panic("Interleaved fail")
		}
	}
}
