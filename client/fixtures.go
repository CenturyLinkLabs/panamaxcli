package client

var (
	templateSearchJSON = `{
		"q": "redis",
		"templates": [
				{
						"created_at": "2015-03-25T20:38:35.198Z",
						"description": "This template is Gitlab with separate containers for PostgreSQL and Redis. This is based on an image from sameersbn / gitlab. You need at least 2 GB RAM and 2 Cores to run this template!",
						"documentation": "GitLab with PostgrSQL + Redis\n=======================\n\n##System requirements\nRecommend 2GB of RAM for your Host and 2 Cores for best performance!\n\n##Setup\nThe values of environemental variables for both PostgreSQL and GitLab need to match. The Keys cannot change. Also, the aliases used in the links to Redis and Posatgres from GHitLab cannot change.\n\nIf the GitLab service does not start up, try the **Rebuild App** function on the application details page to kick start it. Watch the journal for output.\n\nTo view the GUI after launching the template, browse to http://panamax.local:10080.\n\n##Running\n__NOTE__: Please allow a few minutes for the GitLab service to start. Watch the journal output for the message:\n\ndocker 127.0.0.1 - - [DATE/TIME] \"GET /api/v3/internal/check HTTP/1.1\" 200 68 \"-\" \"Ruby\"\n\nLogin using the default username and password:\n\nusername: **root**\n\npassword: **5iveL!fe**\n",
						"icon_src": "/assets/type_icons/default.svg",
						"id": 24,
						"image_count": 3,
						"image_count_label": "Images",
						"images": [
								{
										"category": "Web",
										"command": null,
										"description": null,
										"environment": [
												{
														"value": "10080",
														"variable": "GITLAB_PORT"
												},
												{
														"value": "22",
														"variable": "GITLAB_SSH_PORT"
												},
												{
														"value": "gitlab",
														"variable": "DB_NAME"
												},
												{
														"value": "gitlabuser",
														"variable": "DB_USER"
												},
												{
														"value": "password",
														"variable": "DB_PASS"
												}
										],
										"expose": [],
										"id": 44,
										"links": [
												{
														"alias": "postgresql",
														"service": "PostgreSQL"
												},
												{
														"alias": "redisio",
														"service": "Redis"
												}
										],
										"name": "GitLab",
										"ports": [
												{
														"container_port": "80",
														"host_port": "10080",
														"proto": "TCP"
												},
												{
														"container_port": "22",
														"host_port": "10022",
														"proto": "TCP"
												}
										],
										"source": "centurylink/gitlab:7.1.1",
										"type": "Default",
										"volumes": [],
										"volumes_from": []
								},
								{
										"category": "DB",
										"command": null,
										"description": null,
										"environment": [
												{
														"value": "gitlab",
														"variable": "DB"
												},
												{
														"value": "password",
														"variable": "PASS"
												},
												{
														"value": "gitlabuser",
														"variable": "USER"
												}
										],
										"expose": [
												"5432"
										],
										"id": 45,
										"links": [],
										"name": "PostgreSQL",
										"ports": [
												{
														"container_port": "5432",
														"host_port": "5432"
												}
										],
										"source": "centurylink/postgresql:9.3",
										"type": "Default",
										"volumes": [],
										"volumes_from": []
								},
								{
										"category": "DB",
										"command": null,
										"description": null,
										"environment": [],
										"expose": [
												"6379"
										],
										"id": 46,
										"links": [],
										"name": "Redis",
										"ports": [
												{
														"container_port": "6379",
														"host_port": "6379"
												}
										],
										"source": "sameersbn/redis:latest",
										"type": "Default",
										"volumes": [],
										"volumes_from": []
								}
						],
						"keywords": "gitlab, postgresql, redis, nginx, openssh, public",
						"last_updated_on": "March 25th, 2015 20:38",
						"name": "GitLab 7.1.1 with PostgreSQL and Redis",
						"short_description": "This template is Gitlab with separate containers for PostgreSQL and Redis. This is based on an image from sameersbn /...",
						"source": "centurylinklabs/panamax-public-templates",
						"type": "Default",
						"updated_at": "2015-03-25T20:38:35.217Z"
				}
		]`

	appJSON = `{
			"categories": [
			{
				"id": 1,
				"name": "Web Tier",
				"position": null
			},
			{
				"id": 2,
				"name": "DB Tier",
				"position": null
			}
			],
			"documentation": "Wordpress with MySQL\nFoo",
			"errors": null,
			"from": "Template: Wordpress with MySQL",
			"id": 1,
			"name": "Wordpress with MySQL"
	}`

	appsJSON = `[
		{
			"categories": [
			{
				"id": 1,
				"name": "Web Tier",
				"position": null
			},
			{
				"id": 2,
				"name": "DB Tier",
				"position": null
			}
			],
			"documentation": "Wordpress with MySQL\nFoo",
			"errors": null,
			"from": "Template: Wordpress with MySQL",
			"id": 1,
			"name": "Wordpress with MySQL"
		}
	]`
)
