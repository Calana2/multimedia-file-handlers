package jpg

import "encoding/binary"

func parseSegment(data []byte, size int) Segment {
	var s Segment
	s.Marker = [2]byte{data[0], data[1]}
	s.Length = uint16(size)
	s.Data = data[4:]
	return s
}

func parseAPP(data []byte, size int) APPSegment {
	var s APPSegment
	s.Marker = [2]byte{data[0], data[1]}
	s.Length = uint16(size)
	s.Identifier = string(data[4:9])
	s.Version = string(data[9]+0x30) + "." + string(data[10]+0x30)
	s.DensityUnits = uint8(data[11])
	s.Xdensity = binary.BigEndian.Uint16(data[12:14])
	s.Ydensity = binary.BigEndian.Uint16(data[14:16])
	s.XThumbnail = uint8(data[16])
	s.XThumbnail = uint8(data[17])
	if s.XThumbnail != 0 {
		copy(s.ThumbnailData, data[18:])
	}
	return s
}

func parseEXIF(data []byte, size int) EXIFSegment {
	var s EXIFSegment
	s.Marker = [2]byte{data[0], data[1]}
	s.Length = uint16(size)
	s.Identifier = string(data[4:10])
	s.TIFFHeader.Alignment = string(data[10:12])
	s.TIFFHeader.FixedBytes = [2]byte{data[12], data[13]}
	s.TIFFHeader.IFDOffset = binary.BigEndian.Uint32(data[14:18])
	return s
}

func parseCOM(data []byte, size int) COMSegment {
	var s COMSegment
	s.Marker = [2]byte{data[0], data[1]}
	s.Length = uint16(size)
	s.Data = string(data[4:])
	return s
}

func parseDQT(data []byte, size int) DQTSegment {
	var s DQTSegment
	s.Marker = [2]byte{data[0], data[1]}
	s.Length = uint16(size)
	s.Destination = uint8(data[5])
	s.Data = data[4:]
	return s
}

func parseSOF(data []byte, size int) SOFSegment {
	var s SOFSegment
	s.Marker = [2]byte{data[0], data[1]}
	s.Length = uint16(size)
	s.Precision = data[4]
	s.LineNB = binary.BigEndian.Uint16(data[5:7])
	s.Samples_line = binary.BigEndian.Uint16(data[7:9])
	s.Components = uint8(data[9])
	offset := 10
	for i := 0; i < int(s.Components); i++ {
		var d = &SOFComponentInfo{
			Component:     data[offset],
			Sampling_X__Y: data[offset+1],
			Quantization:  data[offset+2],
		}
		s.SOFComponents = append(s.SOFComponents, *d)
		offset += 3
	}
	return s
}

func parseDHT(data []byte, size int) DHTSegment {
	var s DHTSegment
	s.Marker = [2]byte{data[0], data[1]}
	s.Length = uint16(size)
	s.Class__Idx = uint8(data[4])
	offset := 5
	for i := 0; offset < 21; i++ {
		s.Bit_Codes[i] = uint8(data[offset])
		offset++
	}
	s.Real_Huffman_Codes = data[offset:]
	return s
}

func parseSOS(data []byte, size int) SOSSegment {
	var s SOSSegment
	s.Marker = [2]byte{data[0], data[1]}
	s.Length = uint16(size)
	s.Components = uint8(data[4])
	offset := 0
	switch s.Components {
	case 1:
		s.YIndex = uint8(data[5])
		s.Y_AC__DC = uint8(data[6])
		offset = 7
	case 2:
		s.YIndex = uint8(data[5])
		s.Y_AC__DC = uint8(data[6])
		s.CbIndex = uint8(data[7])
		s.Cb_AC__DC = uint8(data[8])
		offset = 9
	case 3:
		s.YIndex = uint8(data[5])
		s.Y_AC__DC = uint8(data[6])
		s.CbIndex = uint8(data[7])
		s.Cb_AC__DC = uint8(data[8])
		s.CrIndex = uint8(data[9])
		s.Cr_AC__DC = uint8(data[10])
		offset = 11
	}
	s.SS_Start = uint8(data[offset])
	s.SS_End = uint8(data[offset+1])
	s.Sucessive_approx = uint8(data[offset+2])
	return s
}
