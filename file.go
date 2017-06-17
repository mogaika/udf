package udf

import (
	"io"
	"os"
	"time"
)

type File struct {
	Udf               *Udf
	Fid               *FileIdentifierDescriptor
	fe                *FileEntry
	FileEntryPosition uint64
}

func (f *File) FileEntry() *FileEntry {
	if f.fe == nil {
		f.FileEntryPosition = f.Fid.ICB.Location
		f.fe = NewFileEntry(f.Udf.ReadSector(f.Udf.PartitionStart() + f.FileEntryPosition))
	}
	return f.fe
}

func (f *File) NewReader() *io.SectionReader {
	fe := f.FileEntry()
	return io.NewSectionReader(f.Udf.r,
		int64(f.Udf.PartitionStart()+f.Fid.ICB.Location),
		int64(fe.InformationLength))
}

func (f *File) Name() string {
	return f.Fid.FileIdentifier
}

func (f *File) Mode() os.FileMode {
	var mode os.FileMode

	perms := os.FileMode(f.FileEntry().Permissions)
	mode |= ((perms >> 0) & 7) << 0
	mode |= ((perms >> 5) & 7) << 3
	mode |= ((perms >> 10) & 7) << 6

	if f.IsDir() {
		mode |= os.ModeDir
	}

	return mode
}

func (f *File) Size() int64 {
	return int64(f.FileEntry().InformationLength)
}

func (f *File) ModTime() time.Time {
	return f.FileEntry().ModificationTime
}

func (f *File) IsDir() bool {
	// TODO :Fix! This field always 0 :(
	return f.FileEntry().ICBTag.FileType == 4
}

func (f *File) Sys() interface{} {
	return f.Fid
}

func (f *File) ReadDir() []File {
	return f.Udf.ReadDir(f.FileEntry())
}
