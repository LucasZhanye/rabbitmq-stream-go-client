package stream

import (
	"bufio"
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Compression algorithms", func() {
	var entries *subEntries

	BeforeEach(func() {
		messagePayload := make([]byte, 4096)
		for i := range messagePayload {
			messagePayload[i] = 99
		}

		message := &messageSequence{
			messageBytes: messagePayload,
			publishingId: 0,
		}

		entries = &subEntries{
			items: []*subEntry{{
				messages:         []*messageSequence{message},
				publishingId:     0,
				unCompressedSize: len(messagePayload) + 4,
				sizeInBytes:      0,
				dataInBytes:      nil,
			}},
			totalSizeInBytes: 0,
		}
	})

	It("NONE", func() {
		err := compressNONE{}.Compress(entries)
		Expect(err).NotTo(HaveOccurred())
		Expect(entries.totalSizeInBytes).To(Equal(entries.items[0].sizeInBytes))
		Expect(entries.totalSizeInBytes).To(Equal(entries.items[0].unCompressedSize))
	})

	It("GZIP", func() {
		gzip := compressGZIP{}
		err := gzip.Compress(entries)
		Expect(err).NotTo(HaveOccurred())
		verifyCompression(gzip, entries)
	})

	It("SNAPPY", func() {
		snappy := compressSnappy{}
		err := snappy.Compress(entries)
		Expect(err).NotTo(HaveOccurred())
		verifyCompression(snappy, entries)
	})

	It("LZ4", func() {
		lz4 := compressLZ4{}
		err := lz4.Compress(entries)
		Expect(err).NotTo(HaveOccurred())
		verifyCompression(lz4, entries)
	})

	It("ZSTD", func() {
		zstd := compressZSTD{}
		err := zstd.Compress(entries)
		Expect(err).NotTo(HaveOccurred())
		verifyCompression(zstd, entries)
	})
})

func verifyCompression(algo iCompress, subEntries *subEntries) {
	Expect(subEntries.totalSizeInBytes).To(SatisfyAll(BeNumerically("<", subEntries.items[0].unCompressedSize)))
	Expect(subEntries.totalSizeInBytes).To(Equal(subEntries.items[0].sizeInBytes))

	bufferReader := bytes.NewReader(subEntries.items[0].dataInBytes)
	_, err := algo.UnCompress(bufio.NewReader(bufferReader),
		uint32(subEntries.totalSizeInBytes), uint32(subEntries.items[0].unCompressedSize))

	Expect(err).NotTo(HaveOccurred())
}
