package nginxcfg

import (
	"database/sql"
	"os"
	"slices"
	"text/template"

	"github.com/justsushant/envbox/config"
	"github.com/justsushant/envbox/types"
)

type NginxCfgWriter struct {
	tmpl     *template.Template
	confPath string
	config   []types.NginxUpstreamConfig
	db       *sql.DB
}

func NewNginxCfgWriter(db *sql.DB, tmplPath, confPath string) (*NginxCfgWriter, error) {
	// parse the template file
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return nil, err
	}

	// create object based on input
	n := &NginxCfgWriter{
		config:   make([]types.NginxUpstreamConfig, 0),
		tmpl:     tmpl,
		confPath: confPath,
		db:       db,
	}

	// hit db, get details and render the nginx config
	err = n.startupAndRender()
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (n *NginxCfgWriter) startupAndRender() error {
	c, err := n.getAllNginxConfigs()
	if err != nil {
		return err
	}

	n.config = c

	err = n.render()
	return err
}

func (n *NginxCfgWriter) AddUpstreamAndRender(containerID, name, addr string, isRewrite bool) error {
	// insert into db
	_, err := n.db.Exec(
		"INSERT INTO nginx_cfg (containerID, name, addr, isRewrite) VALUES (?, ?, ?, ?)",
		containerID, name, addr, isRewrite,
	)
	if err != nil {
		return err
	}

	// add to local config object
	n.config = append(n.config, types.NginxUpstreamConfig{
		ContainerID: containerID,
		Name:      name,
		Address:   addr,
		IsRewrite: isRewrite,
	})

	// write to file
	err = n.render()
	if err != nil {
		return err
	}

	return nil
}

func (n *NginxCfgWriter) RemoveUpstreamAndRender(containerID string) error {
	// remove from db
	_, err := n.db.Exec("DELETE FROM nginx_cfg WHERE containerID = ?", containerID)
	if err != nil {
		return err
	}

	// remove from local config object
	n.config = slices.DeleteFunc(n.config, func(u types.NginxUpstreamConfig) bool {
        return u.ContainerID == containerID
    })

	// write to file
	err = n.render()
	if err != nil {
		return err
	}

	return nil
}

func (n *NginxCfgWriter) render() error {
	file, err := os.OpenFile(n.confPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)	
	if err != nil {
		return err
	}
	defer file.Close()

	err = n.tmpl.Execute(file, map[string]any{
		"config": n.config,
		"envs": map[string]string{
			"host": config.Envs.Host,
			"port": config.Envs.Port,
			"public": config.Envs.Public,
		},
	})
	return err
}

func (n *NginxCfgWriter) getAllNginxConfigs() ([]types.NginxUpstreamConfig, error) {
	rows, err := n.db.Query("SELECT name, addr, isRewrite FROM nginx_cfg")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	defer rows.Close()

	var configs []types.NginxUpstreamConfig
	for rows.Next() {
		var config types.NginxUpstreamConfig
		err := rows.Scan(&config.Name, &config.Address, &config.IsRewrite)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil
}
