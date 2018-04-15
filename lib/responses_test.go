package travis

var buildResponse = []byte(`
{
  "@type": "builds",
  "@href": "/repo/kevinburke%2Ftravis/builds?branch.name=master",
  "@representation": "standard",
  "@pagination": {
    "limit": 25,
    "offset": 0,
    "count": 2,
    "is_first": true,
    "is_last": true,
    "next": null,
    "prev": null,
    "first": {
      "@href": "/repo/kevinburke%2Ftravis/builds?branch.name=master",
      "offset": 0,
      "limit": 25
    },
    "last": {
      "@href": "/repo/kevinburke%2Ftravis/builds?branch.name=master",
      "offset": 0,
      "limit": 25
    }
  },
  "builds": [
    {
      "@type": "build",
      "@href": "/build/366635873",
      "@representation": "standard",
      "@permissions": {
        "read": true,
        "cancel": true,
        "restart": true
      },
      "id": 366635873,
      "number": "2",
      "state": "failed",
      "duration": 36,
      "event_type": "push",
      "previous_state": "passed",
      "pull_request_title": null,
      "pull_request_number": null,
      "started_at": "2018-04-14T23:03:34Z",
      "finished_at": "2018-04-14T23:04:10Z",
      "repository": {
        "@type": "repository",
        "@href": "/repo/18699435",
        "@representation": "minimal",
        "id": 18699435,
        "name": "travis",
        "slug": "kevinburke/travis"
      },
      "branch": {
        "@type": "branch",
        "@href": "/repo/18699435/branch/master",
        "@representation": "minimal",
        "name": "master"
      },
      "tag": null,
      "commit": {
        "@type": "commit",
        "@representation": "minimal",
        "id": 109419530,
        "sha": "4208212faa90e5ce7cf7c41bcc17ce091ab56137",
        "ref": "refs/heads/master",
        "message": "initial work on open command",
        "compare_url": "https://github.com/kevinburke/travis/compare/a0890b604823...4208212faa90",
        "committed_at": "2018-04-14T23:03:23Z"
      },
      "jobs": [
        {
          "@type": "job",
          "@href": "/job/366635874",
          "@representation": "minimal",
          "id": 366635874
        }
      ],
      "stages": [

      ],
      "created_by": {
        "@type": "user",
        "@href": "/user/6151",
        "@representation": "minimal",
        "id": 6151,
        "login": "kevinburke"
      },
      "updated_at": "2018-04-14T23:04:10.260Z"
    },
    {
      "@type": "build",
      "@href": "/build/366632359",
      "@representation": "standard",
      "@permissions": {
        "read": true,
        "cancel": true,
        "restart": true
      },
      "id": 366632359,
      "number": "1",
      "state": "passed",
      "duration": 56,
      "event_type": "api",
      "previous_state": null,
      "pull_request_title": null,
      "pull_request_number": null,
      "started_at": "2018-04-14T22:44:31Z",
      "finished_at": "2018-04-14T22:45:27Z",
      "repository": {
        "@type": "repository",
        "@href": "/repo/18699435",
        "@representation": "minimal",
        "id": 18699435,
        "name": "travis",
        "slug": "kevinburke/travis"
      },
      "branch": {
        "@type": "branch",
        "@href": "/repo/18699435/branch/master",
        "@representation": "minimal",
        "name": "master"
      },
      "tag": null,
      "commit": {
        "@type": "commit",
        "@representation": "minimal",
        "id": 109418312,
        "sha": "a0890b60482341a26e549f3c005405d6da0517bc",
        "ref": null,
        "message": "Initial commit",
        "compare_url": "https://github.com/kevinburke/travis/commit/a0890b60",
        "committed_at": "2018-04-14T22:42:14Z"
      },
      "jobs": [
        {
          "@type": "job",
          "@href": "/job/366632360",
          "@representation": "minimal",
          "id": 366632360
        }
      ],
      "stages": [

      ],
      "created_by": {
        "@type": "user",
        "@href": "/user/6151",
        "@representation": "minimal",
        "id": 6151,
        "login": "kevinburke"
      },
      "updated_at": "2018-04-14T22:45:27.750Z"
    }
  ]
}`)
