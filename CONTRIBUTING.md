# Contributing

- Fork it
- Create your feature branch: git checkout -b your-new-feature
- Commit changes: git commit -m 'Add your feature'
- Pass all tests
- Push to the branch: git push origin your-new-feature
- Submit a pull request


## Publish the Maven artifacts

Once a release is done, artifacts will need to be promoted in Maven Central to make them generally available.

1. Head over to https://s01.oss.sonatype.org/#stagingRepositories
1. Verify the contents of the staging repo and close it
1. After successful closing (test suite is run), release the repo
