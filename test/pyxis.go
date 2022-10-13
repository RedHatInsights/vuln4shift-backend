package test

const PyxisAPIReposResp = `{
    "data": [
      {
        "_id": "57ea8cd79c624c035f96f327",
        "_links": {
          "certification_project": {
            "href": "/v1/repositories/registry/registry.access.redhat.com/repository/rhel7.1/projects/certification"
          },
          "images": {
            "href": "/v1/repositories/registry/registry.access.redhat.com/repository/rhel7.1/images"
          },
          "operator_bundles": {
            "href": "/v1/repositories/registry/registry.access.redhat.com/repository/rhel7.1/operators/bundles"
          },
          "product_listings": {
            "href": "/v1/repositories/registry/registry.access.redhat.com/repository/rhel7.1/product-listings"
          },
          "replaced_by_repository": {
            "href": "/v1/repositories/registry/registry.access.redhat.com/repository/rhel7/rhel"
          },
          "vendor": {
            "href": "/v1/vendors/label/redhat"
          }
        },
        "application_categories": [
          "Operating System"
        ],
        "beta": false,
        "build_categories": [
          "Standalone image"
        ],
        "can_auto_release_cve_rebuild": false,
        "content_sets": [],
        "creation_date": "2016-09-27T15:14:31.512000+00:00",
        "deprecated": true,
        "description": "This platform image provides a minimal runtime to build, run and deploy Red Hat Enterprise Linux 7.1 applications as a container on a Red Hat Enterprise Linux 7 and Red Hat Enterprise Linux 7 Atomic host.",
        "display_data": {
          "long_description": "The Red Hat Enterprise Linux Base image is designed to be a fully supported foundation for your containerized applications. This base image provides your operations and application teams with the packages, language runtimes and tools necessary to run, maintain, and troubleshoot all of your applications. This image is maintained by Red Hat and updated regularly. It is designed and engineered to be the base layer for all of your containerized applications, middleware and utilities. When used as the source for all of your containers, only one copy will ever be downloaded and cached in your production environment. Use this image just like you would a regular Red Hat Enterprise Linux distribution. Tools like yum, gzip, and bash are provided by default. For further information on how this image was built look at the /root/anacanda-ks.cfg file.",
          "long_description_markdown": "The Red Hat Enterprise Linux Base image is designed to be a fully supported foundation for your containerized applications. This base image provides your operations and application teams with the packages, language runtimes and tools necessary to run, maintain, and troubleshoot all of your applications. This image is maintained by Red Hat and updated regularly. It is designed and engineered to be the base layer for all of your containerized applications, middleware and utilities. When used as the source for all of your containers, only one copy will ever be downloaded and cached in your production environment. Use this image just like you would a regular Red Hat Enterprise Linux distribution. Tools like yum, gzip, and bash are provided by default. For further information on how this image was built look at the /root/anacanda-ks.cfg file.",
          "name": "Red Hat Enterprise Linux 7.1",
          "openshift_tags": "base rhel7.1",
          "short_description": "Provides Red Hat Enterprise Linux 7.1 in a fully featured and supported base image."
        },
        "documentation_links": [
          {
            "title": "UNKNOWN",
            "type": "Documentation",
            "url": "https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux_atomic_host/7/html-single/getting_started_with_containers/#using_standard_rhel_base_images_rhel6_and_rhel7"
          },
          {
            "title": "UNKNOWN",
            "type": "Documentation",
            "url": "https://access.redhat.com/documentation/en/red-hat-enterprise-linux-atomic-host/version-7/getting-started-with-containers/#get_started_with_docker_formatted_container_images"
          }
        ],
        "eol_date": "2015-11-18T00:00:00+00:00",
        "freshness_grades_unknown_until_date": null,
        "label_override": {
          "description": "The Red Hat Enterprise Linux Base image is designed to be a fully supported foundation for your containerized applications. This base image provides your operations and application teams with the packages, language runtimes and tools necessary to run, maintain, and troubleshoot all of your applications. This image is maintained by Red Hat and updated regularly. It is designed and engineered to be the base layer for all of your containerized applications, middleware and utilities. When used as the source for all of your containers, only one copy will ever be downloaded and cached in your production environment. Use this image just like you would a regular Red Hat Enterprise Linux distribution. Tools like yum, gzip, and bash are provided by default. For further information on how this image was built look at the /root/anacanda-ks.cfg file.",
          "io_k8s_displayName": "Red Hat Enterprise Linux 7.1",
          "io_openshift_tags": "base rhel7.1",
          "summary": "Provides Red Hat Enterprise Linux 7.1 in a fully featured and supported base image."
        },
        "last_update_date": "2022-05-12T11:30:50.822000+00:00",
        "metrics": {
          "last_update_date": "2018-03-15T07:00:58.864000+00:00",
          "pulls_in_last_30_days": 421
        },
        "namespace": "rhel7.1",
        "non_production_only": false,
        "object_type": "containerRepository",
        "privileged_images_allowed": false,
        "product_id": "RedHatEnterpriseLinux",
        "product_listings": [
          "5eed1bf53eda773b377f4909"
        ],
        "product_versions": [
          "7"
        ],
        "protected_for_pull": false,
        "protected_for_search": false,
        "published": true,
        "registry": "registry.access.redhat.com",
        "registry_target": "pulp",
        "release_categories": [
          "Deprecated"
        ],
        "replaced_by_repository_name": "rhel7/rhel",
        "repository": "rhel7.1",
        "requires_terms": false,
        "runs_on": {},
        "tech_preview": false,
        "total_size_bytes": 337111006,
        "total_uncompressed_size_bytes": 0,
        "vendor_label": "redhat"
      }
    ],
    "page": 0,
    "page_size": 1,
    "total": 8031
  }`

const PyxisAPIReposRespJBoss = `{
      "_id": "57ea8cd99c624c035f96f332",
      "_links": {
        "certification_project": {
          "href": "/v1/repositories/registry/registry.access.redhat.com/repository/jboss-fuse-6/fis-java-openshift/projects/certification"
        },
        "images": {
          "href": "/v1/repositories/registry/registry.access.redhat.com/repository/jboss-fuse-6/fis-java-openshift/images"
        },
        "operator_bundles": {
          "href": "/v1/repositories/registry/registry.access.redhat.com/repository/jboss-fuse-6/fis-java-openshift/operators/bundles"
        },
        "product_listings": {
          "href": "/v1/repositories/registry/registry.access.redhat.com/repository/jboss-fuse-6/fis-java-openshift/product-listings"
        },
        "vendor": {
          "href": "/v1/vendors/label/redhat"
        }
      },
      "application_categories": [
        "Integration",
        "Middleware",
        "Web Services"
      ],
      "architectures": [],
      "auto_rebuild_tags": [
        "latest"
      ],
      "beta": false,
      "build_categories": [
        "Intermediate image"
      ],
      "can_auto_release_cve_rebuild": true,
      "content_sets": [],
      "content_stream_grades": [
        {
          "grade": "C",
          "tag": "latest"
        }
      ],
      "content_stream_tags": [
        "latest"
      ],
      "creation_date": "2016-09-27T15:14:33.219000+00:00",
      "deprecated": false,
      "description": "Fuse Integration Services, karaf and java builders",
      "display_data": {
        "long_description": "<p>Platform for building and running JBoss Fuse Integration Services</p>\n",
        "long_description_markdown": "Platform for building and running JBoss Fuse Integration Services",
        "name": "JBoss Fuse for OpenShift",
        "openshift_tags": "builder java fuse openshift",
        "short_description": "JBoss Fuse for OpenShift"
      },
      "documentation_links": [
        {
          "title": "UNKNOWN",
          "type": "Documentation",
          "url": "https://access.redhat.com/documentation/en/red-hat-jboss-middleware-for-openshift/3/single/red-hat-jboss-fuse-integration-services-20-for-openshift/"
        },
        {
          "title": "UNKNOWN",
          "type": "Documentation",
          "url": "https://access.redhat.com/documentation/en/red-hat-jboss-fuse"
        }
      ],
      "freshness_grades_unknown_until_date": null,
      "includes_multiple_content_streams": false,
      "label_override": {
        "description": "Platform for building and running JBoss Fuse Integration Services",
        "io_k8s_displayName": "JBoss Fuse for OpenShift",
        "io_openshift_tags": "builder java fuse openshift",
        "summary": "JBoss Fuse for OpenShift"
      },
      "last_update_date": "2022-09-23T00:00:23.595000+00:00",
      "metrics": {
        "last_update_date": "2018-03-15T07:00:26.811000+00:00",
        "pulls_in_last_30_days": 599941
      },
      "namespace": "jboss-fuse-6",
      "non_production_only": false,
      "object_type": "containerRepository",
      "privileged_images_allowed": false,
      "product_id": "JbossFuse",
      "product_listings": [
        "5ec5304978e79e6a879fa269"
      ],
      "protected_for_pull": false,
      "protected_for_search": false,
      "published": true,
      "registry": "registry.access.redhat.com",
      "registry_target": "pulp",
      "release_categories": [
        "Generally Available"
      ],
      "repository": "jboss-fuse-6/fis-java-openshift",
      "requires_terms": false,
      "runs_on": {},
      "tech_preview": false,
      "total_size_bytes": 9158901478,
      "total_uncompressed_size_bytes": 21436136350,
      "vendor_label": "redhat"
    }`

const PyxisAPIRepoImagesResp = `{
    "data": [
        {
            "_id": "57ea8d0d9c624c035f96f45e",
            "_links": {
                "artifacts": {
                    "href": "/v1/images/id/57ea8d0d9c624c035f96f45e/artifacts"
                },
                "requests": {
                    "href": "/v1/images/id/57ea8d0d9c624c035f96f45e/requests"
                },
                "rpm_manifest": {
                    "href": "/v1/images/id/57ea8d0d9c624c035f96f45e/rpm-manifest"
                },
                "test_results": {
                    "href": "/v1/images/id/57ea8d0d9c624c035f96f45e/test-results"
                },
                "vulnerabilities": {
                    "href": "/v1/images/id/57ea8d0d9c624c035f96f45e/vulnerabilities"
                }
            },
            "architecture": "amd64",
            "brew": {
                "build": "rhel-server-docker-7.1-4",
                "completion_date": "2015-02-05T10:07:08+00:00",
                "nvra": "rhel-server-docker-7.1-4.amd64",
                "package": "rhel-server-docker"
            },
            "certified": false,
            "content_sets": [],
            "cpe_ids": [
                "cpe:/o:redhat:rhel_aus:7.6::server",
                "cpe:/a:redhat:satellite:6.0::el7"
            ],
            "creation_date": "2016-09-27T15:15:25.810000+00:00",
            "docker_image_digest": "temp:sha256:3817ddfacc32be3501dce396efcbf864ec68c3d9794a38d0c959377fca65e881",
            "docker_image_id": "sha256:3817ddfacc32be3501dce396efcbf864ec68c3d9794a38d0c959377fca65e881",
            "freshness_grades": [
                {
                    "creation_date": "2019-04-20T06:02:02.865000+00:00",
                    "end_date": "2016-07-22T18:43:00+00:00",
                    "grade": "E",
                    "start_date": "2016-02-02T12:12:00+00:00"
                },
                {
                    "creation_date": "2019-04-20T06:02:02.865000+00:00",
                    "grade": "F",
                    "start_date": "2016-07-22T18:43:00+00:00"
                }
            ],
            "image_id": "10acc31def5d6f249b548e01e8ffbaccfd61af0240c17315a7ad393d022c5ca2",
            "last_update_date": "2022-10-07T01:51:21.689000+00:00",
            "object_type": "containerImage",
            "parsed_data": {
                "architecture": "amd64",
                "command": "[u'/usr/bin/bash']",
                "comment": "Imported from -",
                "created": "2015-02-05T15:04:11.445171245Z",
                "docker_version": "1.4.1-dev",
                "env_variables": [
                    "container=docker"
                ],
                "labels": [
                    {
                        "name": "Vendor",
                        "value": "Red Hat, Inc."
                    },
                    {
                        "name": "Name",
                        "value": "rhel-server-docker"
                    },
                    {
                        "name": "Build_Host",
                        "value": "rcm-img04.build.eng.bos.redhat.com"
                    },
                    {
                        "name": "Version",
                        "value": "7.1"
                    },
                    {
                        "name": "Architecture",
                        "value": "x86_64"
                    },
                    {
                        "name": "Release",
                        "value": "4"
                    },
                    {
                        "name": "BZComponent",
                        "value": "rhel-server-docker"
                    }
                ],
                "layers": [
                    "sha256:605b0bf9a50a8777656bf7559c3778009656c079ac35c1d5aa85c3e809ece617"
                ],
                "os": "linux",
                "size": 0,
                "uncompressed_layer_sizes": [],
                "uncompressed_size_bytes": 0,
                "user": ""
            },
            "raw_config": "{\"architecture\": \"amd64\", \"comment\": \"Imported from -\", \"config\": {\"Hostname\": \"\", \"Domainname\": \"\", \"User\": \"\", \"AttachStdin\": false, \"AttachStdout\": false, \"AttachStderr\": false, \"Tty\": false, \"OpenStdin\": false, \"StdinOnce\": false, \"Env\": [\"container=docker\"], \"Cmd\": [\"/usr/bin/bash\"], \"Image\": \"\", \"Volumes\": null, \"WorkingDir\": \"\", \"Entrypoint\": null, \"OnBuild\": null, \"Labels\": {\"Architecture\": \"x86_64\", \"BZComponent\": \"rhel-server-docker\", \"Build_Host\": \"rcm-img04.build.eng.bos.redhat.com\", \"Name\": \"rhel-server-docker\", \"Release\": \"4\", \"Vendor\": \"Red Hat, Inc.\", \"Version\": \"7.1\"}}, \"container_config\": {\"Hostname\": \"\", \"Domainname\": \"\", \"User\": \"\", \"AttachStdin\": false, \"AttachStdout\": false, \"AttachStderr\": false, \"Tty\": false, \"OpenStdin\": false, \"StdinOnce\": false, \"Env\": null, \"Cmd\": null, \"Image\": \"\", \"Volumes\": null, \"WorkingDir\": \"\", \"Entrypoint\": null, \"OnBuild\": null, \"Labels\": null}, \"created\": \"2015-02-05T15:04:11.445171245Z\", \"docker_version\": \"1.4.1-dev\", \"id\": \"df95058f898d22aca1f4b7ed5c06a0871151d53cb9073de12e6db748a6765c48\", \"os\": \"linux\"}",
            "repositories": [
                {
                    "_links": {
                        "image_advisory": {
                            "href": "/v1/advisories/redhat/id/RHEA-2015:0615"
                        },
                        "repository": {
                            "href": "/v1/repositories/registry/registry.access.redhat.com/repository/rhel"
                        }
                    },
                    "comparison": {
                        "advisory_rpm_mapping": [
                            {
                                "advisory_ids": [
                                    "RHBA-2015:0526"
                                ],
                                "nvra": "rpm-4.11.1-25.el7.x86_64"
                            },
                            {
                                "advisory_ids": [
                                    "RHBA-2015:0501"
                                ],
                                "nvra": "libblkid-2.23.2-21.el7.x86_64"
                            }
                        ],
                        "reason": "OK",
                        "reason_text": "No error",
                        "rpms": {
                            "downgrade": [],
                            "new": [],
                            "remove": [],
                            "upgrade": [
                                "nspr-4.10.6-3.el7.x86_64",
                                "redhat-release-server-7.1-1.el7.x86_64"
                            ]
                        },
                        "with_nvr": "rhel-server-docker-7.0-27"
                    },
                    "content_advisory_ids": [
                        "RHBA-2015:0501",
                        "RHEA-2015:0503"
                    ],
                    "image_advisory_id": "RHEA-2015:0615",
                    "published": true,
                    "published_date": "2022-10-06T17:46:49+00:00",
                    "push_date": "2015-03-05T01:54:47+00:00",
                    "registry": "registry.access.redhat.com",
                    "repository": "rhel",
                    "signatures": [],
                    "tags": [
                        {
                            "_links": {
                                "tag_history": {
                                    "href": "/v1/tag-history/registry/registry.access.redhat.com/repository/rhel/tag/7.1-4"
                                }
                            },
                            "added_date": "2022-10-07T01:51:21.689000+00:00",
                            "manifest_schema1_digest": "sha256:527f5859cd42547c19e52a6f3713112db73241dfcbb11a88b20a613ff6998615",
                            "name": "7.1-4"
                        }
                    ]
                },
                {
                    "_links": {
                        "image_advisory": {
                            "href": "/v1/advisories/redhat/id/RHEA-2015:0615"
                        },
                        "repository": {
                            "href": "/v1/repositories/registry/registry.access.redhat.com/repository/rhel7"
                        }
                    },
                    "comparison": {
                        "advisory_rpm_mapping": [
                            {
                                "advisory_ids": [
                                    "RHBA-2015:0364"
                                ],
                                "nvra": "nspr-4.10.6-3.el7.x86_64"
                            },
                            {
                                "advisory_ids": [
                                    "RHBA-2015:0530"
                                ],
                                "nvra": "libss-1.42.9-7.el7.x86_64"
                            }
                        ],
                        "reason": "OK",
                        "reason_text": "No error",
                        "rpms": {
                            "downgrade": [],
                            "new": [],
                            "remove": [],
                            "upgrade": [
                                "subscription-manager-1.13.19-1.el7.x86_64",
                                "nspr-4.10.6-3.el7.x86_64",
                                "libss-1.42.9-7.el7.x86_64"
                            ]
                        },
                        "with_nvr": "rhel-server-docker-7.0-27"
                    },
                    "content_advisory_ids": [
                        "RHSA-2015:0327",
                        "RHBA-2015:0339"
                    ],
                    "image_advisory_id": "RHEA-2015:0615",
                    "published": true,
                    "published_date": "2022-10-06T17:46:54+00:00",
                    "push_date": "2015-03-05T01:54:47+00:00",
                    "registry": "registry.access.redhat.com",
                    "repository": "rhel7",
                    "signatures": [],
                    "tags": [
                        {
                            "_links": {
                                "tag_history": {
                                    "href": "/v1/tag-history/registry/registry.access.redhat.com/repository/rhel7/tag/7.1-4"
                                }
                            },
                            "added_date": "2022-10-06T17:50:49.428000+00:00",
                            "manifest_schema1_digest": "sha256:c1daa29d08e6c5cd5d4018ec8661ec097de62a57249a34f74af4b25aa45c22e3",
                            "name": "7.1-4"
                        }
                    ]
                },
                {
                    "_links": {
                        "image_advisory": {
                            "href": "/v1/advisories/redhat/id/RHEA-2015:0615"
                        },
                        "repository": {
                            "href": "/v1/repositories/registry/registry.access.redhat.com/repository/rhel7/rhel"
                        }
                    },
                    "comparison": {
                        "advisory_rpm_mapping": [
                            {
                                "advisory_ids": [
                                    "RHBA-2015:0364"
                                ],
                                "nvra": "nspr-4.10.6-3.el7.x86_64"
                            },
                            {
                                "advisory_ids": [
                                    "RHEA-2015:0524"
                                ],
                                "nvra": "redhat-release-server-7.1-1.el7.x86_64"
                            }
                        ],
                        "reason": "OK",
                        "reason_text": "No error",
                        "rpms": {
                            "downgrade": [],
                            "new": [],
                            "remove": [],
                            "upgrade": [
                                "nspr-4.10.6-3.el7.x86_64",
                                "redhat-release-server-7.1-1.el7.x86_64"
                            ]
                        },
                        "with_nvr": "rhel-server-docker-7.0-27"
                    },
                    "content_advisory_ids": [
                        "RHBA-2015:0501",
                        "RHEA-2015:0503"
                    ],
                    "image_advisory_id": "RHEA-2015:0615",
                    "published": true,
                    "published_date": "2022-10-06T17:46:50+00:00",
                    "push_date": "2015-03-05T01:54:47+00:00",
                    "registry": "registry.access.redhat.com",
                    "repository": "rhel7/rhel",
                    "signatures": [],
                    "tags": [
                        {
                            "_links": {
                                "tag_history": {
                                    "href": "/v1/tag-history/registry/registry.access.redhat.com/repository/rhel7/rhel/tag/7.1-4"
                                }
                            },
                            "added_date": "2022-10-06T17:50:43.759000+00:00",
                            "manifest_schema1_digest": "sha256:15255251dfe39bac073e4247b1cf0946d854b0fdeaf3171a10efe006ea5273f9",
                            "name": "7.1-4"
                        }
                    ]
                },
                {
                    "_links": {
                        "image_advisory": {
                            "href": "/v1/advisories/redhat/id/RHEA-2015:0615"
                        },
                        "repository": {
                            "href": "/v1/repositories/registry/registry.access.redhat.com/repository/rhel7.1"
                        }
                    },
                    "comparison": {
                        "reason": "NO_PREVIOUS_IMAGE",
                        "reason_text": "No previous image for comparison found in given repository",
                        "rpms": {
                            "downgrade": [],
                            "new": [],
                            "remove": [],
                            "upgrade": []
                        }
                    },
                    "image_advisory_id": "RHEA-2015:0615",
                    "published": true,
                    "published_date": "2018-04-18T21:19:21+00:00",
                    "push_date": "2015-03-05T01:54:47+00:00",
                    "registry": "registry.access.redhat.com",
                    "repository": "rhel7.1",
                    "signatures": [],
                    "tags": [
                        {
                            "_links": {
                                "tag_history": {
                                    "href": "/v1/tag-history/registry/registry.access.redhat.com/repository/rhel7.1/tag/7.1-4"
                                }
                            },
                            "added_date": "2018-05-27T02:21:41.805000+00:00",
                            "manifest_schema1_digest": "sha256:04dc92b8fa9ba3fb940c17b938275819e83bbdf20a50bdc454f72457027688d8",
                            "name": "7.1-4"
                        }
                    ]
                }
            ],
            "sum_layer_size_bytes": 55046581,
            "top_layer_id": "sha256:605b0bf9a50a8777656bf7559c3778009656c079ac35c1d5aa85c3e809ece617"
        }
    ],
    "page": 0,
    "page_size": 1,
    "total": 6
}`

const PyxisAPIImageCvesResp = `{
  "data": [
    {
      "_id": "59e8d76769aea312bbaec686",
      "_links": {
        "advisory": {
          "href": "/v1/advisories/redhat/id/2016:1940"
        }
      },
      "advisory_id": "2016:1940",
      "advisory_type": "RHSA",
      "creation_date": "2017-10-19T16:48:39.140000+00:00",
      "cve_id": "CVE-2016-2180",
      "last_update_date": "2019-04-15T14:21:53.032000+00:00",
      "object_type": "containerImageVulnerability",
      "packages": [
        {
          "rpm_nvra": [
            "openssl-libs-1.0.1e-42.el7.x86_64"
          ],
          "srpm_nevra": "openssl-1:1.0.1e-51.el7_2.7.src"
        }
      ],
      "public_date": "20160927T11:50:00.000+0000",
      "severity": "Low"
    }
  ],
  "page": 0,
  "page_size": 1,
  "total": 225
}
`
