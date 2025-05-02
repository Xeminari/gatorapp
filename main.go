package main

import (
	"context"
	"database/sql"
	"gatorapp/internal/config"
	"gatorapp/internal/database"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	programState := &state{
		cfg: &cfg,
		db:  dbQueries,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.Register("login", handlerLogin)
	cmds.Register("register", handlerRegister)
	cmds.Register("reset", handlerReset)
	cmds.Register("users", handlerUsers)
	cmds.Register("agg", middlewareLoggedIn(handlerAgg))
	cmds.Register("addfeed", middlewareLoggedIn(handlerAddfeed))
	cmds.Register("feeds", handlerFeeds)
	cmds.Register("follow", middlewareLoggedIn(handlerFollow))
	cmds.Register("following", middlewareLoggedIn(handlerFollowing))
	cmds.Register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.Register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
		return
	}
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	cmd := command{
		Name: cmdName,
		Args: cmdArgs,
	}
	err = cmds.Run(programState, cmd)
	if err != nil {
		log.Fatal(err)
	}

}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
