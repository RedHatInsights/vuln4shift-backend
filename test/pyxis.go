package test

const PyxisAPIReposResp = `{
	  "data": [
		{
		  "_id": "57ea8cd79c624c035f96f327",
		  "last_update_date": "2022-05-12T11:30:50.822000+00:00",
		  "registry": "registry.access.redhat.com",
		  "repository": "rhel7.1"
		}
	  ],
	  "page": 0,
	  "page_size": 1,
	  "total": 8031
	}`

const PyxisAPIReposRespSync = `{
	  "data": [
		{
		  "_id": "57ea8cd79c624c035f962327",
		  "registry": "registry.access.redhat.com",
		  "last_update_date": "2022-05-12T11:30:50.822000+00:00",
		  "repository": "ubi7/ubi-minimal"
		},
		{
		  "_id": "57ea8cd79c624c035f96f444",
		  "registry": "registry.access.redhat.com",
		  "last_update_date": "2022-05-12T11:40:50.822000+00:00",
		  "repository": "ubi8/ubi-micro"
		},
		{
		  "_id": "57ea8cd89c624c035f96f330",
		  "registry": "registry.access.redhat.com",
		  "last_update_date": "2022-05-12T11:40:50.822000+00:00",
		  "repository": "rhel7.1"
		},
		{
		  "_id": "57ea8cd99c624c035f96f332",
		  "registry": "registry.access.redhat.com",
		  "last_update_date": "2022-05-12T11:40:50.822000+00:00",
		  "repository": "rhel7/sadc"
		},
		{
		  "_id": "57ea8cd99c624c035f96f332",
		  "registry": "registry.access.redhat.com",
		  "last_update_date": "2023-05-12T11:40:50.822000+00:00",
		  "repository": "rhel7/sadc"
		}
	  ],
	  "page": 0,
	  "page_size": 3,
	  "total": 8031
	}`

const PyxisAPIReposRespJBoss = `{
	  "_id": "57ea8cd99c624c035f96f332",
	  "last_update_date": "2022-09-23T00:00:23.595000+00:00",
	  "registry": "registry.access.redhat.com",
	  "repository": "jboss-fuse-8/fis-java-openshift"
	}`

const PyxisAPIRepoImagesResp = `{
	  "data": [
		{
		  "_id": "57ea8d0d9c624c035f96f45e",
		  "architecture": "amd64",
		  "docker_image_digest": "temp:sha256:3817ddfacc32be3501dce396efcbf864ec68c3d9794a38d0c959377fca65e881",
		  "last_update_date": "2022-10-07T01:51:21.689000+00:00",
		  "repositories": [
			{
				"manifest_list_digest": "test_manifest_list_digest",
				"manifest_schema2_digest": "test_manifest_schema2_digest"
			}
		  ]
		}
	  ],
	  "page": 0,
	  "page_size": 1,
	  "total": 6
	}`

const PyxisAPIRepoImagesNewResp = `{
	  "data": [
		{
		  "_id": "57ea8d0d9c624c488f96f45e",
		  "architecture": "amd64",
		  "docker_image_digest": "temp:sha256:3817ddfacc32be3501dce396efcbf864ec68c3d9794a38d0c959377fca65e881",
		  "last_update_date": "2022-10-07T01:51:21.689000+00:00",
		  "repositories": [
			{
				"manifest_list_digest": "test_manifest_list_digest",
				"manifest_schema2_digest": "test_manifest_schema2_digest"
			}
		  ]
		}
	  ],
	  "page": 0,
	  "page_size": 1,
	  "total": 6
	}`

const PyxisAPIRepoImagesRespNoDigest = `{
	  "data": [
		{
		  "_id": "57ea8d0d9c624c035f96f45e",
		  "architecture": "amd64",
		  "docker_image_digest": "temp:sha256:3817ddfacc32be3501dce396efcbf864ec68c3d9794a38d0c959377fca65e881",
		  "last_update_date": "2022-10-07T01:51:21.689000+00:00",
		  "repositories": [
			{}
		  ]
		}
	  ],
	  "page": 0,
	  "page_size": 1,
	  "total": 6
	}`

const PyxisAPIRepoImagesResp2 = `{
	  "data": [
		{
		  "_id": "57ea8d0d9c624c035f96f45e",
		  "architecture": "amd64",
		  "docker_image_digest": "temp:sha256:3817ddfacc32be3501dce396efcbf864ec68c3d9794a38d0c959377fca65e881",
		  "last_update_date": "2022-10-07T01:51:21.689000+00:00",
		  "repositories": [
			{
				"manifest_list_digest": "test_manifest_list_digest",
				"manifest_schema2_digest": "test_manifest_schema2_digest"
			}
		  ]
		}
	  ],
	  "page": 0,
	  "page_size": 1,
	  "total": 6
	}`

const PyxisAPIImageCvesResp = `{
	  "data": [
		{
		  "cve_id": "CVE-2016-2180"
		}
	  ],
	  "page": 0,
	  "page_size": 1,
	  "total": 225
	}`
