package table

var (
	Power = &PowerDataTable{file: "power.csv"}
)

var tableList = []iTable{
	Power,
}
