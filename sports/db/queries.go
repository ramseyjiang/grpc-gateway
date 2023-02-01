package db

const (
	sportsList = "list"
)

func getSportQueries() map[string]string {
	return map[string]string{
		sportsList: `
			SELECT 
				id,
				name,
				visible, 
				result,
				start_time, 
				end_time,
				advertised_start_time 
			FROM sports
		`,
	}
}
