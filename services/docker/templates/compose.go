package dockertmpl

import (
	"log"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/system/fs"
	homedir "github.com/mitchellh/go-homedir"
)

func calcPhpVersionString(version string) string {
	var v string
	switch version {
	case "7.1":
		v = "71"
	case "7.2":
		v = "72"
	default:
		log.Fatalf("PHP version %s is not supported. Must use '7.1' or '7.2'", version)
	}

	return v
}

// DrupalSettings ...
func DrupalSettings(drupalRoot string, projectName string) []byte {
	if conf.GetConfig().Global.Syncservice == "fusion" {
		return []byte(`services:
    fpm:
      environment:
        DRUPAL_ROOT: ` + drupalRoot + `
    nginx:
      environment:
        DRUPAL_ROOT: ` + drupalRoot + `
    drush:
      environment:
        DRUPAL_ROOT: ` + drupalRoot + `
        PROJECT_NAME: ` + projectName + `
    kishu:
      environment:
        DRUPAL_ROOT: ` + drupalRoot)
	}
	// Return without kishu for Docker Volume mounts
	return []byte(`services:
    fpm:
      environment:
        DRUPAL_ROOT: ` + drupalRoot + `
    nginx:
      environment:
        DRUPAL_ROOT: ` + drupalRoot + `
    drush:
      environment:
        DRUPAL_ROOT: ` + drupalRoot + `
        PROJECT_NAME: ` + projectName)
}

// StabilityLevel ...
func StabilityLevel(phpVersion, stability string) []byte {
	v := calcPhpVersionString(phpVersion)

	if conf.GetConfig().Global.Syncservice == "fusion" {
		return []byte(`services:
  sync:
    image: tokaido/sync:` + stability + `
  syslog:
    image: tokaido/syslog:` + stability + `
  haproxy:
    image: tokaido/haproxy:` + stability + `
  varnish:
    image: tokaido/varnish:` + stability + `
  nginx:
    image: tokaido/nginx:` + stability + `
  fpm:
    image: tokaido/php` + v + `-fpm:` + stability + `
  drush:
    image: tokaido/admin` + v + `-heavy:` + stability + ``)
	}

	return []byte(`services:
  syslog:
    image: tokaido/syslog:` + stability + `
  haproxy:
    image: tokaido/haproxy:` + stability + `
  varnish:
    image: tokaido/varnish:` + stability + `
  nginx:
    image: tokaido/nginx:` + stability + `
  fpm:
    image: tokaido/php` + v + `-fpm:` + stability + `
  drush:
    image: tokaido/admin` + v + `-heavy:` + stability + ``)
}

// EnableSolr ...
func EnableSolr(version string) []byte {
	return []byte(`services:
  solr:
    image: tokaido/solr:` + version + `
    ports:
      - "8983"
    entrypoint:
      - solr-precreate
      - drupal
      - /opt/solr/server/solr/configsets/search-api-solr/
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name)
}

// EnableRedis ...
func EnableRedis(version string) []byte {
	return []byte(`services:
  redis:
    image: redis:` + version + `
    ports:
      - "6379"
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name)
}

// EnableMailhog ...
func EnableMailhog(version string) []byte {
	return []byte(`services:
  mailhog:
    image: mailhog/mailhog:` + version + `
    ports:
      - "1025"
      - "8025"
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name)
}

// EnableAdminer ...
func EnableAdminer(version string) []byte {
	return []byte(`services:
  adminer:
    image: adminer:` + version + `
    ports:
      - "8080"
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name)
}

// EnableMemcache ...
func EnableMemcache(version string) []byte {
	return []byte(`services:
  memcache:
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
    image: memcached:` + version)
}

// ExternalVolumeDeclare ...
func ExternalVolumeDeclare(name string) []byte {
	return []byte(`volumes:
  ` + name + `:
    external: true
`)
}

// InternalVolumeDeclare ...
func InternalVolumeDeclare(name string) []byte {
	return []byte(`volumes:
  ` + name + `:
    external: false
`)
}

// MysqlVolumeAttach ...
func MysqlVolumeAttach(name string) []byte {
	return []byte(`services:
  mysql:
    volumes:
      - ` + name + `:/var/lib/mysql
`)
}

// SetDatabase sets the database engine and configuration
func SetDatabase(image, version string) []byte {
	return []byte(`services:
  mysql:
    image: ` + image + `:` + version + `
`)
}

// TokaidoFusionSiteVolumeAttach ...
func TokaidoFusionSiteVolumeAttach(path, name string) []byte {
	return []byte(`services:
  sync:
    volumes:
      - ` + path + `:/tokaido/host-volume
      - ` + name + `:/tokaido/site
  drush:
    volumes:
      - ` + name + `:/tokaido/site
      - tok_composer_cache:/home/tok/.composer/cache
  nginx:
    volumes:
      - ` + name + `:/tokaido/site
  fpm:
    volumes:
      - ` + name + `:/tokaido/site
  kishu:
    volumes:
      - ` + name + `:/tokaido/site
`)
}

// TokaidoDockerSiteVolumeAttach ...
func TokaidoDockerSiteVolumeAttach(path string) []byte {
	h, err := homedir.Dir()
	if err != nil {
		log.Fatalf("Could not resolve your home directory: %v", err)
	}

	vols := `services:
  nginx:
    volumes:
      - ` + path + `:/tokaido/site
  fpm:
    volumes:
      - ` + path + `:/tokaido/site
  drush:
    volumes:
      - ` + path + `:/tokaido/site
      - tok_composer_cache:/home/tok/.composer/cache`

	// We'll mount the .gitconfig and .drush paths if they exist
	gp := h + "/.gitconfig"
	dp := h + "/.drush:/home/tok/.drush"

	if fs.CheckExists(gp) {
		vols = vols + `
      - ` + h + `/.gitconfig:/home/tok/.gitconfig`
	}

	if fs.CheckExists(dp) {
		vols = vols + `
      - ` + h + `/.gitconfig:/home/tok/.drush`
	}

	return []byte(vols)
}

// ModWarning - Displayed at the top of `docker-compose.tok.yml`
var ModWarning = []byte(`
# WARNING: THIS FILE IS MANAGED DIRECTLY BY TOKAIDO.
# DO NOT MAKE MODIFICATIONS HERE, THEY WILL BE OVERWRITTEN

`)

// ComposeTokDefaultsFusionSync - Template byte array for `docker-compose.tok.yml`
var ComposeTokDefaultsFusionSync = []byte(`
version: "2"
services:
  sync:
    image: tokaido/sync:stable
    volumes:
      - waiting
    environment:
      AUTO_SYNC: "true"
    restart: unless-stopped
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  syslog:
    image: tokaido/syslog:stable
    volumes:
      - /tokaido/logs
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  haproxy:
    user: "1005"
    image: tokaido/haproxy:stable
    ports:
      - "8080"
      - "8443"
    depends_on:
      - varnish
      - nginx
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  varnish:
    user: "1004"
    image: tokaido/varnish:stable
    ports:
      - "8081"
    depends_on:
      - nginx
    volumes_from:
      - syslog
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  nginx:
    user: "1002"
    image: tokaido/nginx:stable
    volumes:
      - waiting
    volumes_from:
      - syslog
    depends_on:
      - fpm
    ports:
      - "8082"
    environment:
      DRUPAL_ROOT: docroot
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  fpm:
    user: "1001"
    image: tokaido/php71-fpm:stable
    working_dir: /tokaido/site/
    volumes:
      - waiting
    volumes_from:
      - syslog
    depends_on:
      - syslog
    ports:
      - "9000"
    environment:
      PHP_DISPLAY_ERRORS: "yes"
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  mysql:
    image: mysql:5.7
    volumes_from:
      - syslog
    volumes:
      - waiting
    ports:
      - "3306"
    command: --max_allowed_packet=1073741824 --ignore-db-dir=lost+found
    environment:
      MYSQL_DATABASE: tokaido
      MYSQL_USER: tokaido
      MYSQL_PASSWORD: tokaido
      MYSQL_ROOT_PASSWORD: tokaido
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  drush:
    image: tokaido/admin71-heavy:stable
    hostname: 'tokaido'
    ports:
      - "22"
    working_dir: /tokaido/site
    volumes:
      - waiting
    volumes_from:
      - syslog
    environment:
      SSH_AUTH_SOCK: /ssh/auth/sock
      APP_ENV: local
      PROJECT_NAME: tokaido
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  kishu:
    image: tokaido/kishu:stable
    volumes:
      - waiting
    environment:
      DRUPAL_ROOT: docroot
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
`)

// ComposeTokDefaultsDockerVolume - Template byte array for `docker-compose.tok.yml`
var ComposeTokDefaultsDockerVolume = []byte(`
version: "2"
services:
  syslog:
    image: tokaido/syslog:stable
    volumes:
      - /tokaido/logs
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  haproxy:
    user: "1005"
    image: tokaido/haproxy:stable
    ports:
      - "8080"
      - "8443"
    depends_on:
      - varnish
      - nginx
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  varnish:
    user: "1004"
    image: tokaido/varnish:stable
    ports:
      - "8081"
    depends_on:
      - nginx
    volumes_from:
      - syslog
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  nginx:
    user: "1002"
    image: tokaido/nginx:stable
    volumes:
      - waiting
    volumes_from:
      - syslog
    depends_on:
      - fpm
    ports:
      - "8082"
    environment:
      DRUPAL_ROOT: docroot
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  fpm:
    user: "1001"
    image: tokaido/php71-fpm:stable
    working_dir: /tokaido/site/
    volumes:
      - waiting
    volumes_from:
      - syslog
    depends_on:
      - syslog
    ports:
      - "9000"
    environment:
      PHP_DISPLAY_ERRORS: "yes"
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  mysql:
    image: mysql:5.7
    volumes_from:
      - syslog
    volumes:
      - waiting
    ports:
      - "3306"
    command: --max_allowed_packet=1073741824 --ignore-db-dir=lost+found
    environment:
      MYSQL_DATABASE: tokaido
      MYSQL_USER: tokaido
      MYSQL_PASSWORD: tokaido
      MYSQL_ROOT_PASSWORD: tokaido
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  drush:
    image: tokaido/admin71-heavy:stable
    hostname: 'tokaido'
    ports:
      - "22"
    working_dir: /tokaido/site
    volumes:
      - waiting
    volumes_from:
      - syslog
    environment:
      SSH_AUTH_SOCK: /ssh/auth/sock
      APP_ENV: local
      PROJECT_NAME: tokaido
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
`)
