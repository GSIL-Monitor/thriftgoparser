package proto

import "fmt"

const (
	INVALID = 0
	STILL_CANDIDATE = 1     // the text may or may not be the program doctext
	ALREADY_PROCESSED = 2    // doctext has been used and is no longer available
	ABSOLUTELY_SURE = 3     // this is the program doctext
	NO_PROGRAM_DOCTEXT = 4    // there is no program doctext
)

var(
	G_doctext_lineno int32
	G_program_doctext_lineno int32
	G_program_doctext_status int32

)

type TDoc struct{
	Doc_ string
	HasDoc_ bool
}

func NewDoc(doc string) *TDoc{
	ret := new(TDoc)
	ret.SetDoc(doc)
	return ret
}

func (t *TDoc)SetDoc(doc string) {
	t.Doc_ = doc;
	t.HasDoc_ = true
	if((G_program_doctext_lineno == G_doctext_lineno) &&  (G_program_doctext_status == STILL_CANDIDATE)) {
		G_program_doctext_status = ALREADY_PROCESSED;
		fmt.Printf("%s\n","program doctext set to ALREADY_PROCESSED");
	}
}

func (t *TDoc)HasDoc() bool{
	return t.HasDoc_;
}

func (t *TDoc) GetDoc() string {
	return t.Doc_
}

func (t *TDoc)SetCafDoc(cafDoc ICafDoc){
	t.Doc_ = cafDoc.GetDoc()
	t.HasDoc_ = cafDoc.HasDoc()
}


type ICafDoc interface{
	GetDoc() string
	HasDoc() bool
}