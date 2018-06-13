package main

// BuildInfo - contains full build information.
type BuildInfo struct {
	UUID    string `json:"uuid"`
	Succeed bool   `json:"succeed"`
	Score   int    `json:"score"`
	Details string `json:"details"`
}

// RegisterBuild - contains information required to register new build
// Language - either "c++" or "pascal"
type RegisterBuild struct {
	UUID       string `json:"uuid"`
	Source     string `json:"source"`
	WebHookURL string `json:"web_hook_url"`
	Language   string `json:"language"`
}

func getBuildInfo(c APIContext) error {
	key := c.Vars()["uuid"]

	db, err := c.ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()
	repo := NewRepository(db)
	row, err := repo.GetBuildInfo(key)
	if err != nil {
		return err
	}

	info := &BuildInfo{
		UUID:    key,
		Succeed: row.Succeed,
		Score:   row.Score,
		Details: row.Report,
	}
	return c.WriteJSON(info)
}

func createBuild(c APIContext) error {
	var params RegisterBuild
	err := c.ReadJSON(params)
	if err != nil {
		return err
	}

	db, err := c.ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	repo := NewRepository(db)
	err = repo.RegisterBuild(RegisterBuildParams{
		Key:        params.UUID,
		Source:     params.Source,
		WebHookURL: params.WebHookURL,
	})

	return err
}
