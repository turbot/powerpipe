package intern

// Pre-interned resource type strings.
// These are commonly repeated across all resources.
var (
	TypeDashboard          string
	TypeBenchmark          string
	TypeControl            string
	TypeQuery              string
	TypeCard               string
	TypeChart              string
	TypeContainer          string
	TypeFlow               string
	TypeGraph              string
	TypeHierarchy          string
	TypeImage              string
	TypeInput              string
	TypeNode               string
	TypeEdge               string
	TypeTable              string
	TypeText               string
	TypeCategory           string
	TypeDetection          string
	TypeDetectionBenchmark string
	TypeVariable           string
	TypeWith               string
	TypeLocals             string
)

// Pre-interned benchmark type strings.
var (
	BenchmarkTypeControl   string
	BenchmarkTypeDetection string
)

// Pre-interned common tag keys.
var (
	TagService  string
	TagCategory string
	TagType     string
	TagPlugin   string
	TagClass    string
)

// Pre-interned common attribute names.
var (
	AttrTitle       string
	AttrDescription string
	AttrSQL         string
	AttrQuery       string
	AttrArgs        string
	AttrWidth       string
	AttrBase        string
	AttrType        string
)

func init() {
	preInternCommonStrings(DefaultInterner)
}

// preInternCommonStrings pre-interns commonly repeated strings into an interner.
func preInternCommonStrings(i *StringInterner) {
	// Resource types
	TypeDashboard = i.Intern("dashboard")
	TypeBenchmark = i.Intern("benchmark")
	TypeControl = i.Intern("control")
	TypeQuery = i.Intern("query")
	TypeCard = i.Intern("card")
	TypeChart = i.Intern("chart")
	TypeContainer = i.Intern("container")
	TypeFlow = i.Intern("flow")
	TypeGraph = i.Intern("graph")
	TypeHierarchy = i.Intern("hierarchy")
	TypeImage = i.Intern("image")
	TypeInput = i.Intern("input")
	TypeNode = i.Intern("node")
	TypeEdge = i.Intern("edge")
	TypeTable = i.Intern("table")
	TypeText = i.Intern("text")
	TypeCategory = i.Intern("category")
	TypeDetection = i.Intern("detection")
	TypeDetectionBenchmark = i.Intern("detection_benchmark")
	TypeVariable = i.Intern("variable")
	TypeWith = i.Intern("with")
	TypeLocals = i.Intern("locals")

	// Benchmark types
	BenchmarkTypeControl = i.Intern("control")
	BenchmarkTypeDetection = i.Intern("detection")

	// Common tag keys
	TagService = i.Intern("service")
	TagCategory = i.Intern("category")
	TagType = i.Intern("type")
	TagPlugin = i.Intern("plugin")
	TagClass = i.Intern("class")

	// Common attribute names
	AttrTitle = i.Intern("title")
	AttrDescription = i.Intern("description")
	AttrSQL = i.Intern("sql")
	AttrQuery = i.Intern("query")
	AttrArgs = i.Intern("args")
	AttrWidth = i.Intern("width")
	AttrBase = i.Intern("base")
	AttrType = i.Intern("type")
}
