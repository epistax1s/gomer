package report

import "github.com/epistax1s/gomer/internal/server"

func publishReport(server *server.Server, report []string) error {
	groups, err := server.GroupService.FindAll()
	if err != nil {
		return err
	}

	for _, group := range groups {
		for _, message := range report {
			server.Gomer.SendMessage(group.GroupID, message)
		}
	}

	return nil
}
