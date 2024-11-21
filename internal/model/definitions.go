package model


// AtomicType represents basic types like Integer, String, etc.
type AtomicType string

const (
	Integer AtomicType = "Integer"
	Double  AtomicType = "Double"
	String  AtomicType = "String"
	Boolean AtomicType = "Boolean"
)

// AtomicType represents basic types like Integer, String, etc.
type Difficulty string

const (
	Hard   Difficulty = "Hard"
	Medium Difficulty = "Medium"
	Easy   Difficulty = "Easy"
)

// CompositeType represents composite types like TreeNode, Array, etc.
type CompositeType string

const (
	GraphNode CompositeType = "GraphNode"
	TreeNode  CompositeType = "TreeNode"
	ListNode  CompositeType = "ListNode"
	Array     CompositeType = "Array"
	Matrix    CompositeType = "Matrix"
)

type PredefinedSupportedLanguage string

const (
	JavaScript PredefinedSupportedLanguage = "JavaScript"
	Python     PredefinedSupportedLanguage = "Python"
)

// PredefinedCategory represents categories for questions.
type PredefinedCategory string

const (
	ArrayCategory              PredefinedCategory = "Array"
	GraphCategory              PredefinedCategory = "Graph"
	StringCategory             PredefinedCategory = "String"
	TreeCategory               PredefinedCategory = "Tree"
	DynamicProgrammingCategory PredefinedCategory = "DynamicProgramming"
	LinkedListCategory         PredefinedCategory = "LinkedList"
	MatrixCategory             PredefinedCategory = "Matrix"
)
