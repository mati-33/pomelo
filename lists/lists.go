package lists

var (
	lists []List = []List{
		{
			Name:     "pomelo project",
			Created:  "12:17 AM",
			Modified: "-",
			Tasks: []Task{
				{Name: "tui"},
				{Name: "backend"},
			},
		},
		{
			Name:     "06-10-2026",
			Created:  "11:11 AM",
			Modified: "12:12",
			Tasks:    []Task{},
		},
		{
			Name:     "07-10-2026",
			Created:  "10:10 PM",
			Modified: "-",
			Tasks:    []Task{},
		},
		{
			Name:     "terminal typing app",
			Created:  "04:27 AM",
			Modified: "-",
			Tasks:    []Task{},
		},
		{
			Name:     "terminal chat app",
			Created:  "08:00 PM",
			Modified: "09:30",
			Tasks:    []Task{},
		},
	}
)

type Task struct {
	Name   string
	IsDone bool
}

type List struct {
	Name     string
	Created  string
	Modified string
	Tasks    []Task
}

func GetAllLists() []List {
	return lists
}

func AddList(l List) {
	lists = append(lists, l)
}
