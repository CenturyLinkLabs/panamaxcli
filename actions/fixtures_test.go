package actions

var wordpressTemplate = `
name: Wordpress with MySQL
description: Wordpress container linked to a MySQL container
keywords: wordpress, mysql, public
recommended: true
documentation: |
  Wordpress with MySQL
  ============================
  The alias for the link between Wordpress and MySQL needs to be _DB_1_. If this is changed, the template will not work.

  Also, the password can be changed, but the environmental variables need to be changed on both services.

  To view the GUI after launching the template go to http://10.0.0.200:8080 or http://panamax.local:8080 in a browser 

authors:
- 'ctl-labs-futuretech@savvis.com'
type: wordpress
images:
- name: WP
  source: centurylink/wordpress:3.9.1
  description: Wordpress
  environment:
    - variable: DB_PASSWORD
      value: pass@word01
    - variable: DB_NAME
      value: wordpress
  links:
  - service: DB
    alias: DB_1
  ports:
  - host_port: 8080
    container_port: 80
  category: Web Tier
  type: wordpress
- name: DB
  source: centurylink/mysql:5.5
  description: MySQL
  environment:
    - variable: MYSQL_ROOT_PASSWORD
      value: pass@word01
  ports:
  - host_port: 3306
    container_port: 3306
  category: DB Tier
  type: mysql
`
