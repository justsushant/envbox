package utils

import (
	"github.com/justsushant/envbox/types"
)

var IMAGE_MIG_TABLE = []types.ImageMigration{
	{
		Label: "Jupyter Notebook",
		Name:  "envbox/jupyter:0.1",
		Path:  "./migration/DockerFiles/Dockerfile_jupyter_0.1",
	},
	{
		Label: "VS Code (Go)",
		Name:  "envbox/vscode_go:0.1",
		Path:  "./migration/DockerFiles/Dockerfile_vscode_go_0.1",
	},
}

var CREATE_TABLES_STMT = `
	CREATE TABLE IF NOT EXISTS mst_images (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		path TEXT
	);

	CREATE TABLE IF NOT EXISTS containers_running (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		imageID INTEGER,
		containerID TEXT,
		accessLink TEXT,
		active BOOLEAN NOT NULL DEFAULT 1 CHECK (active IN (0, 1)),
		createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS nginx_cfg (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		containerID INTEGER,
		name TEXT NOT NULL,
		addr TEXT NOT NULL,
		isRewrite BOOLEAN NOT NULL
	);
`
