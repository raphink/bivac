package providers

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
  "github.com/camptocamp/conplicity/handler"
)

type OpenLDAPProvider struct {
  handler   *handler.Conplicity
  vol       *docker.Volume
  backupDir string
}

func (p *OpenLDAPProvider) PrepareBackup() (err error) {
  log.Infof("PG_VERSION file found, this should be a OpenLDAP datadir")
  log.Infof("Searching postgres container using this volume...")
  c := p.handler.Client
  vol := p.vol
  log.Infof("DB_CONFIG file found, this should be and OpenLDAP datadir")
  log.Infof("Searching OpenLDAP container using this volume...")
  containers, _ := c.ListContainers(docker.ListContainersOptions{})
  for _, container := range containers {
    for _, mount := range container.Mounts {
      if mount.Name == vol.Name {
        log.Infof("Volume %v is used by container %v", vol.Name, container.ID)
        log.Infof("Launch slapcat in container %v...", container.ID)
        exec, err := c.CreateExec(
          docker.CreateExecOptions{
            Container: container.ID,
            Cmd: []string{
              "sh",
              "-c",
              "mkdir -p " + mount.Destination + "/backups && slapcat > " + mount.Destination + "/backups/all.ldif",
            },
          },
        )

        checkErr(err, "Failed to create exec", 1)

        err = c.StartExec(
          exec.ID,
          docker.StartExecOptions{},
        )

        checkErr(err, "Failed to create exec", 1)

        p.backupDir = "backups"
      }
    }
  }
  return
}

func (p *OpenLDAPProvider) BackupVolume() (err error) {
  return
}