package remote

// FIXME: Maybe create separate "static" repos
/*func Test_integration_Tags(t *testing.T) {
	gitRemoteService := NewGitRemoteService()
	tags, err := gitRemoteService.Tags("https://github.com/AmadlaOrg/QAFixturesEntityMultipleTagVersion")
	if err != nil {
		t.Errorf("Tags returned an error: %s", err)
	}

	expectedTags := []string{
		"v1.0.0",
		"v1.0.0-alpha.2",
		"v1.0.0-beta.1",
		"v2.0.0",
	}

	sort.Strings(tags)
	sort.Strings(expectedTags)

	if !reflect.DeepEqual(tags, expectedTags) {
		t.Errorf("Tags do not match expected values.\nGot: %v\nExpected: %v", tags, expectedTags)
	}

	// Optionally log the sorted tags for debugging purposes
	for _, tag := range tags {
		t.Logf("Retrieved tag: %s", tag)
	}
}*/

// FIXME: Maybe create separate "static" repos
/*func Test_integration_CommitHeadHash(t *testing.T) {
	gitRemoteService := NewGitRemoteService()
	hash, err := gitRemoteService.CommitHeadHash("https://github.com/AmadlaOrg/QAFixturesEntityMultipleTagVersion")
	if err != nil {
		t.Errorf("CommitHeadHash returned an error: %s", err)
	}

	assert.Equal(t, hash, "c351cf75321ae8a7676b8bef6837b67a60cabdbc")
}*/
