package argocd

type ApplicationManifestQuery struct {
	Name                 *string  `protobuf:"bytes,1,req,name=name" json:"name,omitempty"`
	Revision             *string  `protobuf:"bytes,2,opt,name=revision" json:"revision,omitempty"`
	AppNamespace         *string  `protobuf:"bytes,3,opt,name=appNamespace" json:"appNamespace,omitempty"`
	Project              *string  `protobuf:"bytes,4,opt,name=project" json:"project,omitempty"`
	SourcePositions      []int64  `protobuf:"varint,5,rep,name=sourcePositions" json:"sourcePositions,omitempty"`
	Revisions            []string `protobuf:"bytes,6,rep,name=revisions" json:"revisions,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}
