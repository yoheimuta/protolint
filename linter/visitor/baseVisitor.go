package visitor

import (
	"github.com/yoheimuta/go-protoparser/v4/parser"
)

// BaseVisitor represents a base visitor with noop logic.
type BaseVisitor struct{}

// OnStart works noop.
func (BaseVisitor) OnStart(*parser.Proto) error { return nil }

// Finally works noop.
func (BaseVisitor) Finally() error { return nil }

// VisitComment works noop.
func (BaseVisitor) VisitComment(*parser.Comment) {}

// VisitDeclaration works noop.
func (BaseVisitor) VisitDeclaration(*parser.Declaration) (next bool) { return true }

// VisitEdition works noop.
func (BaseVisitor) VisitEdition(*parser.Edition) (next bool) { return true }

// VisitEmptyStatement works noop.
func (BaseVisitor) VisitEmptyStatement(*parser.EmptyStatement) (next bool) { return true }

// VisitEnum works noop.
func (BaseVisitor) VisitEnum(*parser.Enum) (next bool) { return true }

// VisitEnumField works noop.
func (BaseVisitor) VisitEnumField(*parser.EnumField) (next bool) { return true }

// VisitExtensions works noop.
func (BaseVisitor) VisitExtensions(*parser.Extensions) bool { return true }

// VisitExtend works noop.
func (BaseVisitor) VisitExtend(*parser.Extend) (next bool) { return true }

// VisitField works noop.
func (BaseVisitor) VisitField(*parser.Field) (next bool) { return true }

// VisitGroupField works noop.
func (BaseVisitor) VisitGroupField(*parser.GroupField) bool { return true }

// VisitImport works noop.
func (BaseVisitor) VisitImport(*parser.Import) (next bool) { return true }

// VisitMapField works noop.
func (BaseVisitor) VisitMapField(*parser.MapField) (next bool) { return true }

// VisitMessage works noop.
func (BaseVisitor) VisitMessage(*parser.Message) (next bool) { return true }

// VisitOneof works noop.
func (BaseVisitor) VisitOneof(*parser.Oneof) (next bool) { return true }

// VisitOneofField works noop.
func (BaseVisitor) VisitOneofField(*parser.OneofField) (next bool) { return true }

// VisitOption works noop.
func (BaseVisitor) VisitOption(*parser.Option) (next bool) { return true }

// VisitPackage works noop.
func (BaseVisitor) VisitPackage(*parser.Package) (next bool) { return true }

// VisitReserved works noop.
func (BaseVisitor) VisitReserved(*parser.Reserved) (next bool) { return true }

// VisitRPC works noop.
func (BaseVisitor) VisitRPC(*parser.RPC) (next bool) { return true }

// VisitService works noop.
func (BaseVisitor) VisitService(*parser.Service) (next bool) { return true }

// VisitSyntax works noop.
func (BaseVisitor) VisitSyntax(*parser.Syntax) (next bool) { return true }
